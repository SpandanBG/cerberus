package keys

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/json"
)

const (
	MAXMSG    = 86
	MAXCIPHER = 128
)

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
	var buffer bytes.Buffer
	msgBuffer := bytes.NewBuffer(message)
	for msgBuffer.Len() > 0 {
		cipherText, err := makeCipherText(publicKey, msgBuffer.Next(MAXMSG))
		if err != nil {
			return []byte{}, err
		}
		buffer.Write(cipherText)
	}
	return buffer.Bytes(), nil
}

func (keys *Keys) Decrypt(cipherText []byte) ([]byte, error) {
	var buffer bytes.Buffer
	cipherBuffer := bytes.NewBuffer(cipherText)
	for cipherBuffer.Len() > 0 {
		plainText, err := solveCipherText(keys.PrivateKey, cipherBuffer.Next(MAXCIPHER))
		if err != nil {
			return []byte{}, err
		}
		buffer.Write(plainText)
	}
	return buffer.Bytes(), nil
}

func makeCipherText(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptOAEP(
		sha1.New(),
		rand.Reader,
		publicKey,
		message,
		[]byte(""),
	)
}

func solveCipherText(privateKey *rsa.PrivateKey, cipherText []byte) ([]byte, error) {
	return rsa.DecryptOAEP(
		sha1.New(),
		rand.Reader,
		privateKey,
		cipherText,
		[]byte(""),
	)
}
