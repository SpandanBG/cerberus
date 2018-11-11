package keys

import (
	conn "../connection"
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"net"
	"strconv"
)

const (
	MAXMSG    = 62
	MAXCIPHER = 128
)

type Keys struct {
	RemotePublicKey *rsa.PublicKey
	PublicKey       *rsa.PublicKey
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

func (keys *Keys) ExportPublicKey() {
	keys.PublicKey = &keys.PrivateKey.PublicKey
	return
}

/*Needs modification*/
func (keys *Keys) GetRemotePublicKey(connect *conn.Connection) error {
	head := &conn.Header{Version: 1, REQ: true, KX: true}
	reqPacket, err := conn.GeneratePacket(head, keys.PublicKey, []byte(""))
	if err != nil {
		return err
	}
	remoteAddr := connect.RemoteAddr.IP.String() + ":" + strconv.Itoa(connect.RemoteAddr.Port)
	TConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(TConn)
	resPacket := make([]byte, 5620)
	if n, err := writer.Write(reqPacket); err == nil {
		writer.Flush()
		fmt.Println(remoteAddr, "BYtes Written : ", n)
		// udpConn.SetReadDeadline(time.Now().Add(time.Second * 5))
		if n, err := TConn.Read(resPacket); err == nil {
			fmt.Println(remoteAddr, "BYtes Written : ", n)
			_, keys.RemotePublicKey, _, err = conn.ParserPacketBytes(resPacket[:n])
			fmt.Println(keys.RemotePublicKey)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (keys *Keys) Encrypt(message []byte) ([]byte, error) {
	var buffer bytes.Buffer
	msgBuffer := bytes.NewBuffer(message)
	for msgBuffer.Len() > 0 {
		cipherText, err := makeCipherText(keys.RemotePublicKey, msgBuffer.Next(MAXMSG))
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
	//fmt.Println(len(message))
	return rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		message,
		[]byte(""),
	)
}

func solveCipherText(privateKey *rsa.PrivateKey, cipherText []byte) ([]byte, error) {
	return rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privateKey,
		cipherText,
		[]byte(""),
	)
}
