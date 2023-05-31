package helpers

import "strings"

func ByteArrayModification(slc []byte, modification string) []byte {
	return []byte(string(slc) + modification)
}

func TrimByteArray(message []byte) []byte {
	return []byte(strings.Trim(string(message), "\n"))
}
