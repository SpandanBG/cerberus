package bin

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"hash"
)

var (
	// PrivateKey holds the private rsa key
	PrivateKey *rsa.PrivateKey

	// PublicKey holds the public rsa key
	PublicKey *rsa.PublicKey

	// SHA2 holds the sha256 hash function for RSA signature
	SHA2 hash.Hash
)

// InitRSAKeys inits the RSA keys
func InitRSAKeys() {
	var pErr error
	PrivateKey, pErr = rsa.GenerateKey(rand.Reader, 2048)
	if pErr == nil {
		panic("RSA Keys Gen Error")
	}
	PublicKey = &PrivateKey.PublicKey
	SHA2 = sha256.New()
}

// Encrypt encrypts the message
func Encrypt(msg []byte) ([]byte, error) {
	return rsa.EncryptOAEP(SHA2, rand.Reader, PublicKey, msg, nil)
}

// Decrypt encrypts the message
func Decrypt(msg []byte) ([]byte, error) {
	return rsa.DecryptOAEP(SHA2, rand.Reader, PrivateKey, msg, nil)
}
