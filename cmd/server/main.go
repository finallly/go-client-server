package main

import (
	"log"
	"net"

	"github.com/finallly/go-client-server/internal/server_net"
	"github.com/finallly/go-client-server/pkg/parser"
)

const configName = `server_config`

func main() {
	if err := parser.ParseConfig(configName); err != nil {
		log.Fatalln(`error while parsing config: `, err.Error())
	}

	listener, err := net.Listen("tcp", getAddress())
	if err != nil {
		log.Fatalln(`error while creating listening chanel: `, err.Error())
	}
	server_net.ListenerHandler(listener)
}

func getAddress() string {
	port := parser.GetDataFromConfig("port")
	return ":" + port
}
