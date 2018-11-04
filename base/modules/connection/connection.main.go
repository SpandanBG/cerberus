package connection

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type Connection struct {
	LocalAddr  *net.UDPAddr
	RemoteAddr *net.UDPAddr
	UDP        *net.UDPConn
}

func NewConnection(addr string) *Connection {
	Local, _ := net.ResolveUDPAddr("udp", addr)
	return &Connection{LocalAddr: Local}
}

func (conn *Connection) OpenUDPPort() error {
	udpConn, err := net.ListenUDP("udp", conn.LocalAddr)
	if err != nil {
		return err
	}
	conn.UDP = udpConn
	return nil
}

func (conn Connection) LaunchUDPTracer(tracingAddr string) (*net.UDPAddr, error) {
	BCast, err := net.ResolveUDPAddr("udp", tracingAddr)
	if err != nil {
		fmt.Println(err)
	}
	if conn.UDP == nil {
		return nil, errors.New("No UDP Server Setup Found")
	}
	requestPacket, err := CreateTracerPacket()
	responsePacket := make([]byte, len(requestPacket)+10)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		timeout := time.Now().Add(time.Second * 5)
		for time.Now().Before(timeout) {
			fmt.Println("Tracing... Attempt", i)
			conn.UDP.WriteToUDP(requestPacket, BCast)
			fmt.Println("Written!", BCast)
			for {
				fmt.Println("Written!")
				_, addr, _ := conn.UDP.ReadFromUDP(responsePacket)
				fmt.Println("Written!")
				fmt.Println(addr)
				if addr.IP.String() != conn.LocalAddr.IP.String() {
					if TracerResponseValid(responsePacket) == true {
						return addr, nil
					}
					break
				}
			}
		}
	}
	return nil, errors.New("Tracer Returned Nothing")
}

func TracerResponseValid(response []byte) bool {
	head, _, _, err := ParserPacketBytes(response)
	if err != nil {
		return false
	}
	if head.DREQ && head.RES {
		return true
	}
	return false
}

func CreateTracerPacket() ([]byte, error) {
	head := &Header{DREQ: true}
	body := []byte("")
	return GeneratePacket(head, nil, body)
}
