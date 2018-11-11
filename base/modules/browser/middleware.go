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
func EncryptChannel(Header *c.Header, RQ []byte, P *c.Packet, K *k.Keys) chan []byte {
	EC := make(chan []byte)
	//if RQ, ok := <-DC; ok {
	var err error
	fmt.Println(RQ)
	EncryptedRQ, err := K.Encrypt(RQ)
	fmt.Println(EncryptedRQ)
	e.ErrorHandler(err)
	RQPacket, err := c.GeneratePacket(Header, K.PublicKey, EncryptedRQ)
	e.ErrorHandler(err)
	fmt.Println(RQPacket)
	EC <- RQPacket
	return EC
	//}
	//return nil
}

/*DecryptChannel : decrypt channel using RSA Private Key Passed*/
func DecryptChannel(EC chan []byte, P *c.Packet, K *k.Keys) chan []byte {
	DC := make(chan []byte)

	if RS, ok := <-EC; ok {
		var err error
		_, _, P.Body, err = c.ParserPacketBytes(RS)
		e.ErrorHandler(err)
		DecryptedRS, err := K.Decrypt(P.Body)
		e.ErrorHandler(err)
		DC <- DecryptedRS
		return DC
	}
	return nil
}
