// main.go

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mwolff44/pyfb-router/config"
	"github.com/mwolff44/pyfb-router/internal/db"
	"github.com/mwolff44/pyfb-router/internal/server"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	db.Init()
	server.Init()

}
