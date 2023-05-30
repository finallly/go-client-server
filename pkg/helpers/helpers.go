package helpers

func ByteArrayModification(slc []byte, modification string) []byte {
	return []byte(string(slc) + modification)
}
