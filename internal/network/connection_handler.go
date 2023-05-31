package network

import (
	"bufio"
	"crypto/rsa"
	"encoding/json"
	"net"

	"github.com/finallly/go-client-server/internal/encryption"
	"github.com/finallly/go-client-server/pkg/helpers"

	"github.com/charmbracelet/log"
)

func handleConnection(connection net.Conn) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error(`caught panic.`, `error`, err)
		}
	}()

	clientPublicKey, _ := bufio.NewReader(connection).ReadBytes('\n')

	publicKey := &rsa.PublicKey{}
	err := json.Unmarshal(helpers.TrimByteArray(clientPublicKey), &publicKey)

	keyPair := &encryption.KeyPair{
		Public: publicKey,
	}

	if err != nil {
		return err
	}

	log.Info(`public key received from client.`, `key`, *publicKey)

	arbiter, err := encryption.NewArbiter(nil)

	if err != nil {
		return err
	}

	encryptedKey, err := keyPair.EncryptWithPublicKey(arbiter.Key)

	if err != nil {
		return err
	}

	// at this point client and server both have trusted arbiter and can communicate via aes encrypted messages
	_, err = connection.Write(helpers.ByteArrayModification(encryptedKey, "\n"))

	go func() {
		err := handleEcho(connection, arbiter)
		if err != nil {
			log.Fatal(`error while processing echo.`, `error`, err.Error())
		}
	}()

	return nil
}

func handleEcho(connection net.Conn, arbiter *encryption.Arbiter) error {
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			log.Debug(`error closing connection.`, `error`, err.Error())
		}
	}(connection)

	for {
		message, err := bufio.NewReader(connection).ReadBytes('\n')

		if err != nil {
			log.Error(`error while reading message`)

			return err
		}

		if len(message) == 0 {
			continue
		}

		message, err = arbiter.Decrypt(helpers.TrimByteArray(message))

		if err != nil {
			log.Error(`error while decrypting message`)

			return err
		}

		log.Info(`message received from client.`, `message`, string(message))

		message, err = arbiter.Encrypt(append(message, []byte(` + server sign`)...))

		if err != nil {
			log.Error(`error while encrypting message`)

			return err
		}

		_, err = connection.Write(helpers.ByteArrayModification(message, "\n"))

		if err != nil {
			return err
		}
	}
}
