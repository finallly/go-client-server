package network

import (
	"github.com/charmbracelet/log"
	"net"
)

func ListenerHandler(netListener net.Listener) {
	defer func(netListener net.Listener) {
		_ = netListener.Close()
	}(netListener)
	for {
		connection, _ := netListener.Accept()

		log.Info(`client connected.`, `client`, connection.RemoteAddr())

		go func() {
			err := handleConnection(connection)

			if err != nil {
				log.Error(`error while handling connection.`, `error`, err.Error())
			}
		}()
	}
}
