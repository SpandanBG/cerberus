package connection

import (
	"net"
)

type Connection struct {
	Addr string
	TCP  *net.TCPListener
	UDP  *net.UDPConn
}

func NewConnection(addr string) *Connection {
	return &Connection{Addr: addr}
}

func (conn *Connection) OpenTCPPort() (err error) {
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", conn.Addr)
	if err == nil {
		conn.TCP, err = net.ListenTCP("tcp", tcpAddr)
	}
	return
}

func (conn *Connection) OpenUDPPort() (err error) {
	var udpAddr *net.UDPAddr
	udpAddr, err = net.ResolveUDPAddr("udp", conn.Addr)
	if err == nil {
		conn.UDP, err = net.ListenUDP("udp", udpAddr)
	}
	return
}
