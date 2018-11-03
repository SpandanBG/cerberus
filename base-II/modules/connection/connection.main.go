package connection

import (
	"net"
)

type Address struct {
	IP   string
	Port int
}

type Connection struct {
	Addr Address
	TCP  *net.Conn
	UDP  *net.UDPConn
}

func NewConnection(addr Address) *Connection {
	return &Connection{Addr: addr}
}

func (conn *Connection) OpenUDPPort() error {
	udpServerAddr := TranslateAddressToUDPAddr(conn.Addr)
	udpConn, err := net.ListenUDP("udp", udpServerAddr)
	if err != nil {
		return err
	}
	conn.UDP = udpConn
	return nil
}

func TranslateAddressToUDPAddr(addr Address) *net.UDPAddr {
	return &net.UDPAddr{
		IP:   net.ParseIP(addr.IP),
		Port: addr.Port,
	}
}

func TracerRequestValid(response []byte) bool {
	head, _, err := ParserJSONPacket(response)
	if err != nil {
		return false
	}
	if head.DREQ {
		return true
	}
	return false
}

func CreateTracerRespnsePacket() ([]byte, error) {
	head := &Header{DREQ: true, RES: true}
	body := []byte("")
	return GeneratePacket(head, body)
}
