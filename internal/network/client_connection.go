package network

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"time"

	"github.com/finallly/go-client-server/internal/encryption"
	"github.com/finallly/go-client-server/pkg/helpers"
	"github.com/finallly/go-client-server/pkg/parser"

	"github.com/charmbracelet/log"
)

func StartClientConnection() error {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	connection, err := net.Dial("tcp", getAddress())

	if err != nil {
		return err
	}

	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			logger.Warn(`error closing connection.`, `error`, err.Error())
		}
	}(connection)

	keyPair, err := encryption.GenerateRsaKeyPair()

	if err != nil {
		return err
	}

	publicKeyBytes, err := json.Marshal(keyPair.Public)

	if err != nil {
		return err
	}

	err = helpers.WriteMessage(connection, publicKeyBytes)

	if err != nil {
		return err
	}

	secretKey, err := helpers.ReadMessage(connection)
	if err != nil {
		logger.Fatalf("failed to read: %s", err.Error())
	}
	secretKey, err = keyPair.DecryptWithPrivateKey(secretKey)

	log.Info(`secret key from server.`, `key`, secretKey)

	if err != nil {
		return err
	}

	// at this point client and server both have trusted arbiter and can communicate via aes encrypted messages
	arbiter, err := encryption.NewArbiter(secretKey)

	if err != nil {
		return err
	}

	for {
		reader := bufio.NewReader(os.Stdin)

		message, err := reader.ReadBytes('\n')

		if err != nil {
			return err
		}

		message, err = arbiter.Encrypt(message[:len(message)-1])

		if err != nil {
			logger.Error(`error while encrypting message`)

			return err
		}

		err = helpers.WriteMessage(connection, message)
		if err != nil {
			return err
		}

		message, err = helpers.ReadMessage(connection)

		if err != nil {
			return err
		}

		message, err = arbiter.Decrypt(message)

		if err != nil {
			logger.Error(`error while decrypting message`)

			return err
		}

		log.Info(`message received from client.`, `message`, string(message))
	}
}

func getAddress() string {
	host := parser.GetDataFromConfig("host")
	port := parser.GetDataFromConfig("port")

	return host + ":" + port
}
