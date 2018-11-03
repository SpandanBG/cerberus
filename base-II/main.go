package main

import (
	"fmt"
	"log"

	conn "./modules/connection"
)

const (
	BrtIP = "192.168.1.255"
	IP    = "192.168.1.104"
	Port  = 4123
)

var RouterAddr *conn.Address

func init() {
	// load rsa keys
	go loadSearchRequestReceiver()
}

func main() {
	var a int
	fmt.Scan(&a)
	fmt.Println(a)
}

func loadSearchRequestReceiver() {
	req := make([]byte, 5620)
	udpServer := conn.NewConnection(conn.Address{IP: IP, Port: Port})
	err := udpServer.OpenUDPPort()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tracer Listening Server Started")
	for {
		if n, remoteAddr, err := udpServer.UDP.ReadFromUDP(req); err == nil {
			fmt.Println(remoteAddr)
			if conn.TracerRequestValid(req[:n]) {
				if res, err := conn.CreateTracerRespnsePacket(); err == nil {
					udpServer.UDP.WriteToUDP(res, remoteAddr)
				}
			}
		}
	}
}
