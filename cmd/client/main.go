package main

import (
	"github.com/finallly/go-client-server/pkg/parser"
	"log"

	"github.com/finallly/go-client-server/internal/network"
)

const configFile = `client_config`

func main() {
	if err := parser.ParseConfig(configFile); err != nil {
		log.Fatalln(`error reading config from file: `, err.Error())
	}

	if err := network.StartClientConnection(); err != nil {
		log.Fatalln(`error connection to server: `, err.Error())
	}
}
