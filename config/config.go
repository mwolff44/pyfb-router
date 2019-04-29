package config

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error
	config = viper.New()
	config.SetDefault("server.port", "127.0.0.1:8081")
	config.SetDefault("GIN_MODE", "debug")
	config.SetDefault("POSTGRES_HOST", "postgres")
	config.SetDefault("POSTGRES_PORT", "5432")
	config.SetDefault("POSTGRES_DB", "pyfreebilling")
	config.SetDefault("POSTGRES_USER", "pyfreebilling")
	config.SetDefault("POSTGRES_PASSWORD", "secret")
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	config.AddConfigPath(".")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
}

func relativePath(basedir string, path *string) {
	p := *path
	if len(p) > 0 && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

// GetConfig exports configuration settings
func GetConfig() *viper.Viper {
	return config
}
