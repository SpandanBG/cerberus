package browser

import (
	c "../connection"
	e "../error"
	k "../keys"
)

/*EncryptChannel : encrypt channel using the RSA Public Key passed */
func EncryptChannel(DC chan []byte, P *c.Packet, K *k.Keys) chan []byte {
	var err error
	EC := make(chan []byte)

	if RQ, ok := <-DC; ok {
		EncryptedRQ, err := K.Encrypt(RQ)
		e.ErrorHandler(err)
		RQPacket, err := c.GeneratePacket(P.Header, P.Key, EncryptedRQ)
		e.ErrorHandler(err)
		EC <- RQPacket
		return EC
	}
}

/*DecryptChannel : decrypt channel using RSA Private Key Passed*/
func DecryptChannel(EC chan []byte, P *c.Packet, K *k.Keys)  chan []byte {
	var err error
	DC := make([]byte)

	if RS, ok := <-EC; ok {
		P.Header, P.Key, P.Body, err := c.ParserPacketBytes(RS)
		e.ErrorHandler(err)
		DecryptedRS, err := K.Decrypt(P.Body)
		e.ErrorHandler(err)
		DC <- DecryptedRS
	}
}
