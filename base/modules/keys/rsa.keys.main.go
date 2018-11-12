package keys

import (
	conn "../connection"
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	MAXMSG    = 86
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
	_, err = ProxyWriter(&TConn, reqPacket)
	if err == nil {
		resPacket, err := ProxyReader(&TConn)
		_, keys.RemotePublicKey, _, err = conn.ParserPacketBytes(resPacket)
		if err != nil {
			return err
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

func ProxyWriter(conn *net.Conn, data []byte) (int, error) {
	writer := bufio.NewWriter(*conn)
	n, err := writer.Write(data)
	if err != nil {
		return 0, err
	}
	writer.Flush()
	fmt.Println("Bytes Written : ", n)
	return n, nil
}

func ProxyReader(conn *net.Conn) ([]byte, error) {
	Reader := bufio.NewReader(*conn)
	var resBuffer bytes.Buffer
	response := make([]byte, 20)
	for {
		n, err := Reader.Read(response)
		if err != nil {
			if err == io.EOF {
				break
			}
			return []byte{}, err
		}
		resBuffer.Write(response[:n])
	}
	return resBuffer.Bytes(), nil
}
