package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
)

const Version = 1

type Keys struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func NewKeys() *Keys {
	return &Keys{}
}

func (keys *Keys) CreateRSAPair() error {
	var err error
	keys.PrivateKey, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}
	keys.PublicKey = &keys.PrivateKey.PublicKey
	return nil
}

func (keys *Keys) ExportPublicKey() ([]byte, error) {
	return json.Marshal(keys.PrivateKey.PublicKey)
}

func (keys *Keys) Encrypt(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		message,
		[]byte(""),
	)
}

func (keys *Keys) Decrypt(cipherText []byte) ([]byte, error) {
	return rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		keys.PrivateKey,
		cipherText,
		[]byte(""),
	)
}
