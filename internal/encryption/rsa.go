package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

type KeyPair struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

func GenerateRsaKeyPair() (*KeyPair, error) {
	private, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Private: private,
		Public:  &private.PublicKey,
	}, nil
}

func (kp *KeyPair) EncryptWithPublicKey(text []byte) ([]byte, error) {
	hash := sha512.New()

	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, kp.Public, text, nil)

	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func (kp *KeyPair) DecryptWithPrivateKey(text []byte) ([]byte, error) {
	hash := sha512.New()

	ciphertext, err := rsa.DecryptOAEP(hash, rand.Reader, kp.Private, text, nil)

	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}
