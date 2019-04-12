package server

import (
	"github.com/mwolff44/pyfb-router/config"
)

// Init intialise a server
func Init() {
	config := config.GetConfig()
	r := NewRouter()
	r.Run(config.GetString("server.port"))
}
