package connection

import (
	"fmt"
	"net"

	"../configs"
)

// RunTracerServer : Launches the tracer server on a seperate routine
func RunTracerServer(conn *net.UDPConn) {
	go func() {
		for {
			packet := make([]byte, 128)
			n, addr, err := conn.ReadFromUDP(packet)
			fmt.Println("Tracer Packet from : " + addr.String())
			if err == nil {
				go SendTracerResponse(packet[:n], conn, addr)
			}
		}
	}()
}

// SendTracerResponse : Send a response to the tracer packet
func SendTracerResponse(packet []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	reqHeader, _, _, err := ParserPacketBytes(packet)
	if err != nil {
		fmt.Println("Tracer Packet From " + addr.String() + " Error : " + err.Error())
		return
	}
	if VerifyTracerPacket(reqHeader) == false {
		fmt.Println("Tracer Packet From " + addr.String() + " Error : " + err.Error())
		return
	}
	resPacket, err := GeneratePacket(&Header{Version: configs.VERSION, DREQ: true, RES: true}, nil, []byte{})
	if err != nil {
		fmt.Println("Tracer Packet From " + addr.String() + " Error : " + err.Error())
		return
	}
	conn.WriteToUDP(resPacket, addr)
}

// VerifyTracerPacket : Verifies the validity of request tracer packet
func VerifyTracerPacket(header *Header) bool {
	required := header.Version == configs.VERSION
	required = required && header.DREQ
	notRequired := header.REQ && header.RES && header.KX
	return required && !notRequired
}
