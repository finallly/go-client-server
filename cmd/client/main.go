package main

import (
	"github.com/charmbracelet/log"
	"github.com/finallly/go-client-server/internal/network"
	"github.com/finallly/go-client-server/pkg/parser"
)

const configFile = `client_config`

func main() {
	if err := parser.ParseConfig(configFile); err != nil {
		log.Error(`error reading config from file.`, `error`, err.Error())
	}

	if err := network.StartClientConnection(); err != nil {
		log.Error(`error connection to server`, `error`, err.Error())
	}
}
