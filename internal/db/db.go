package db

import (
	"os"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/log15adapter"
	log "gopkg.in/inconshreveable/log15.v2"
)

// db for pgx Postgresql driver
var (
	pool             *pgx.ConnPool
	Debug            = true
	maxDBConnections = 5
)

// Init configures PG setup
func Init() {
	//c := config.GetConfig()
	logger := log15adapter.NewLogger(log.New("module", "pgx"))

	var err error
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "127.0.0.1",
			User:     "uNTVzQUOEszJcfxfaxxLjCSDYYueaOFt",
			Password: "XwUHE4OfKk0DAAW0zbkSC54IK21oVSz7pZpQZZt6OdVdWZ6gFsNWWcnk9ERxgVqj",
			Database: "pyfreebilling",
			Port:     5433,
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
