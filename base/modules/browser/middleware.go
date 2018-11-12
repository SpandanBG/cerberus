package browser

import (
	c "../connection"
	e "../error"
	k "../keys"
	"crypto/rsa"
)

/*InitPacketHeader : intialize packet header*/
func InitPacketHeader(Version int, PK *rsa.PublicKey) *c.Header {
	header := c.Header{Version: Version, REQ: true}
	return &header
}

/*EncryptChannel : encrypt channel using the RSA Public Key passed */
func EncryptChannel(Header *c.Header, RQ []byte, K *k.Keys) []byte {
	var err error
	EncryptedRQ, err := K.Encrypt(RQ)
	e.ErrorHandler(err)
	RQPacket, err := c.GeneratePacket(Header, &K.PrivateKey.PublicKey, EncryptedRQ)
	e.ErrorHandler(err)
	return RQPacket
}

/*DecryptChannel : decrypt channel using RSA Private Key Passed*/
func DecryptChannel(RS []byte, K *k.Keys) []byte {
	var err error
	var body []byte
	_, _, body, err = c.ParserPacketBytes(RS)
	DecryptedRS, err := K.Decrypt(body)
	e.ErrorHandler(err)
	return DecryptedRS
}
