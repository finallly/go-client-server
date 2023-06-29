package helpers

import (
	"encoding/binary"
	"io"
	"net"
)

func WriteMessage(conn net.Conn, message []byte) error {
	prefix := make([]byte, 4)
	binary.BigEndian.PutUint32(prefix, uint32(len(message)))

	_, err := conn.Write(prefix)
	_, err = conn.Write(message)
	return err
}

func ReadMessage(conn net.Conn) ([]byte, error) {
	prefix := make([]byte, 4)
	_, err := io.ReadFull(conn, prefix)

	length := binary.BigEndian.Uint32(prefix)

	message := make([]byte, length)
	_, err = io.ReadFull(conn, message)
	return message, err
}
