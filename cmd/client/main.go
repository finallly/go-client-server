package main

import (
	"github.com/finallly/go-client-server/pkg/parser"
	"log"

	"github.com/finallly/go-client-server/internal/server_net"
)

const configFile = `client_config`

func main() {
	if err := parser.ParseConfig(configFile); err != nil {
		log.Fatalln(`error reading config from file: `, err.Error())
	}

	if err := server_net.StartClientConnection(); err != nil {
		log.Fatalln(`error connection to server: `, err.Error())
	}
}
