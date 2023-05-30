package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type Arbiter struct {
	Key []byte
	GCM cipher.AEAD
}

func NewArbiter(key []byte) (*Arbiter, error) {
	arbiter := &Arbiter{}
	_, err := arbiter.generateSecretKey(key)

	if err != nil {
		return nil, err
	}

	return arbiter, nil
}

func (a *Arbiter) generateSecretKey(key []byte) ([]byte, error) {
	if key == nil {
		key = make([]byte, 32)

		_, err := rand.Read(key)

		if err != nil {
			return nil, err
		}
	}

	c, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, err
	}

	a.Key = key
	a.GCM = gcm

	return key, nil
}

func (a *Arbiter) Encrypt(text []byte) ([]byte, error) {

	nonce := make([]byte, a.GCM.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return a.GCM.Seal(nonce, nonce, text, nil), nil
}

func (a *Arbiter) Decrypt(text []byte) ([]byte, error) {
	nonceSize := a.GCM.NonceSize()

	nonce, ciphertext := text[:nonceSize], text[nonceSize:]

	plaintext, err := a.GCM.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
