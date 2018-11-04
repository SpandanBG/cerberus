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

// Connection : The TCP and UDP connection
var Connection *connection.Connection

func init() {
	loadRSAKeyPair()
	setupConnectionConfig()
	startTracerServer()
}

func main() {
	fmt.Println("Setting Up Proxy Server")
	Connection.OpenTCPPort()
	fmt.Println("Proxy Server Open")
	fmt.Println("****************** PROTECTED BY CERBERUS ******************")
	fmt.Println(configs.ART, "\n\n")
	for {
		if incoming, err := Connection.TCP.Accept(); err == nil {
			go proxyHTTPRequest(incoming)
		}
	}
}

func loadRSAKeyPair() {
	fmt.Println("Creating RSA Keys Pairs")
	configs.RSAKeys = keys.NewKeys()
	if err := configs.RSAKeys.CreateRSAPair(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("RSA Keys Pairs Generated")
}

func setupConnectionConfig() {
	fmt.Println("Setting Up Connection Configurations")
	addr := ":" + configs.PORT
	Connection = connection.NewConnection(addr)
	fmt.Println("Connection Configurations Setup Done")
}

func startTracerServer() {
	fmt.Println("Starting Tracer Server")
	if err := Connection.OpenUDPPort(); err != nil {
		log.Fatal(err)
	}
	connection.RunTracerServer(Connection.UDP)
	fmt.Println("Tracer Server Open")
}

func proxyHTTPRequest(conn net.Conn) {
	rAddr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	reqPacket := make([]byte, configs.HTTPHEADERSIZE+configs.CERBERUSHEADERSIZE)
	n, err := reader.Read(reqPacket)
	if err != nil {
		fmt.Println("HTTP Proxy REQ Error : " + err.Error())
		return
	}
	resPacket, err := connection.ProxyHandler(reqPacket[:n], rAddr)
	if err != nil {
		fmt.Println("HTTP Proxy RES Error : " + err.Error())
		return
	}
	writer.Write(resPacket)
	writer.Flush()
}
