package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"net"
	"strconv"
	"time"

	conn "../connection"
)

type Keys struct {
	RemotePublicKey *rsa.PublicKey
	PrivateKey      *rsa.PrivateKey
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
	return nil
}

func (keys *Keys) ExportPublicKey() ([]byte, error) {
	return json.Marshal(keys.PrivateKey.PublicKey)
}

/*Needs modification*/
func (keys *Keys) GetRemotePublicKey(connect *conn.Connection) error {
	head := &conn.Header{Version: 1, REQ: true, KX: true}
	reqPacket, err := conn.GeneratePacket(head, []byte(""))
	if err != nil {
		return err
	}
	remoteAddr := connect.RemoteAddr.IP.String() + ":" + strconv.Itoa(connect.RemoteAddr.Port)
	udpConn, err := net.Dial("udp", remoteAddr)
	if err != nil {
		return err
	}
	resPacket := make([]byte, 5620)
	for i := 1; i <= 10; i++ {
		if _, err := udpConn.Write(reqPacket); err == nil {
			udpConn.SetReadDeadline(time.Now().Add(time.Second * 5))
			if n, err := udpConn.Read(resPacket); err == nil {
				err = json.Unmarshal(resPacket[:n], keys.RemotePublicKey)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

func (keys *Keys) Encrypt(message []byte) ([]byte, error) {
	return rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		keys.RemotePublicKey,
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
