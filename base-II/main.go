package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"./modules/configs"
	"./modules/connection"
	"./modules/keys"
)

// RSAKeys : The RSA key pair
var RSAKeys *keys.Keys

// Connection : The TCP and UDP connection
var Connection *connection.Connection

func init() {
	loadRSAKeyPair()
	fmt.Println("RSA Keys Pairs Generated")
	setupConnectionConfig()
	fmt.Println("Connection Configurations Setup Done")
	startTracerServer()
	fmt.Println("Tracer Server Open")
}

func main() {
	// Start TCP Server
	Connection.OpenTCPPort()
	for {
		if incoming, err := Connection.TCP.Accept(); err == nil {
			go proxyHTTPRequest(&incoming)
		}
	}
}

func loadRSAKeyPair() {
	RSAKeys = keys.NewKeys()
	if err := RSAKeys.CreateRSAPair(); err != nil {
		log.Fatal(err)
	}
}

func setupConnectionConfig() {
	addr := ":" + configs.PORT
	Connection = connection.NewConnection(addr)
}

func startTracerServer() {
	if err := Connection.OpenUDPPort(); err != nil {
		log.Fatal(err)
	}
	connection.RunTracerServer(Connection.UDP)
}

func proxyHTTPRequest(conn *net.Conn) {
	reader := bufio.NewReader(*conn)
	writer := bufio.NewWriter(*conn)
	buff := make([]byte, 1024)
	n, err := reader.Read(buff)
	if err != nil {
		fmt.Println("HTTP Proxy Error : " + err.Error())
		return
	}
	fmt.Println(buff[:n])
	writer.Write([]byte("hello"))
	writer.Flush()
}
