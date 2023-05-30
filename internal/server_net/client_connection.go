package server_net

import (
	"bufio"
	"encoding/json"
	"net"
	"os"

	"github.com/finallly/go-client-server/internal/encryption"
	"github.com/finallly/go-client-server/pkg/helpers"
	"github.com/finallly/go-client-server/pkg/parser"

	"github.com/charmbracelet/log"
)

func StartClientConnection() error {
	connection, err := net.Dial("tcp", getAddress())

	if err != nil {
		return err
	}

	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			log.Debug(`error closing connection: `, err.Error())
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

	_, err = connection.Write(helpers.ByteArrayModification(publicKeyBytes, "\n"))

	if err != nil {
		return err
	}

	secretKey, _ := bufio.NewReader(connection).ReadBytes('\n')
	secretKey, err = keyPair.DecryptWithPrivateKey(secretKey)

	log.Info(`secret key from server: `, secretKey)

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

		message, err = arbiter.Encrypt(message)

		if err != nil {
			log.Error(`error while encrypting message`)

			return err
		}

		_, err = connection.Write(helpers.ByteArrayModification(message, "\n"))

		if err != nil {
			return err
		}

		message, _ = bufio.NewReader(connection).ReadBytes('\n')

		message, err = arbiter.Decrypt(message)

		if err != nil {
			log.Error(`error while decrypting message`)

			return err
		}

		log.Info(`message received from client: `, string(message))
	}
}

func getAddress() string {
	host := parser.GetDataFromConfig("host")
	port := parser.GetDataFromConfig("port")

	return host + ":" + port
}
