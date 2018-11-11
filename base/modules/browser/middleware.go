package browser

import (
	c "../connection"
	e "../error"
	k "../keys"
	"crypto/rsa"
	"fmt"
)

/*InitPacketHeader : intialize packet header*/
func InitPacketHeader(b *BrowserConn, Version int, PK *rsa.PublicKey) *c.Header {
	header := c.Header{Version: Version, REQ: true}
	return &header
}

/*EncryptChannel : encrypt channel using the RSA Public Key passed */
func EncryptChannel(Header *c.Header, RQ []byte, P *c.Packet, K *k.Keys) []byte {
	var err error
	fmt.Println(RQ)
	EncryptedRQ, err := K.Encrypt(RQ)
	e.ErrorHandler(err)
	RQPacket, err := c.GeneratePacket(Header, K.PublicKey, EncryptedRQ)
	e.ErrorHandler(err)
	return RQPacket
}

/*DecryptChannel : decrypt channel using RSA Private Key Passed*/
func DecryptChannel(RS []byte, P *c.Packet, K *k.Keys) []byte {
	var err error
	_, _, P.Body, err = c.ParserPacketBytes(RS)
	e.ErrorHandler(err)
	DecryptedRS, err := K.Decrypt(P.Body)
	e.ErrorHandler(err)
	return DecryptedRS
}
