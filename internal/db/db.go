package db

import (
	"os"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/log15adapter"
	log "gopkg.in/inconshreveable/log15.v2"

	"github.com/mwolff44/pyfb-router/config"
)

// db for pgx Postgresql driver
var (
	pool             *pgx.ConnPool
	Debug            = true
	maxDBConnections = 5
)

// Init configures PG setup
func Init() {
	config := config.GetConfig()
	logger := log15adapter.NewLogger(log.New("module", "pgx"))

	var err error
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     config.GetString("POSTGRES_HOST"),
			User:     config.GetString("POSTGRES_USER"),
			Password: config.GetString("POSTGRES_PASSWORD"),
			Database: config.GetString("POSTGRES_DB"),
			Port:     uint16(config.GetInt("POSTGRES_PORT")),
			Logger:   logger,
		},
		MaxConnections: maxDBConnections,
		AfterConnect:   afterConnect,
	}
	pool, err = pgx.NewConnPool(connPoolConfig)
	if err != nil {
		log.Crit("Unable to create connection pool", "error", err)
		os.Exit(1)
	}

}

// afterConnect creates the prepared statements that this application uses
func afterConnect(conn *pgx.Conn) (err error) {
	for name, sql := range statements {
		_, err := conn.Prepare(name, sql)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetDB returns pool instance
func GetDB() *pgx.ConnPool {
	return pool
}
