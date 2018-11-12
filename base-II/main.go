package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
			go proxyHTTPConnection(incoming)
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

func proxyHTTPConnection(conn net.Conn) {
	rAddr := conn.RemoteAddr().String()
	reqPacket, err := proxyRequestReader(&conn)
	if err != nil {
		fmt.Println("HTTP Proxy REQ Error :", err.Error())
		return
	}
	resPacket, err := connection.ProxyHandler(reqPacket, rAddr)
	if err != nil {
		fmt.Println("HTTP Proxy RES Error :", err.Error())
		return
	}
	n, err := proxyResponseWriter(&conn, resPacket)
	if err != nil {
		fmt.Println("Socket Write Error", err, "Written", n)
	}
	fmt.Println("Responded", rAddr, "with", n, "bytes")
}

func proxyRequestReader(conn *net.Conn) ([]byte, error) {
	reader := bufio.NewReader(*conn)
	var readerBuffer bytes.Buffer
	for {
		reqPacket := make([]byte, 20)
		n, err := reader.Read(reqPacket)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return []byte{}, err
			}
		}
		readerBuffer.Write(reqPacket[:n])
	}
	return readerBuffer.Bytes(), nil
}

func proxyResponseWriter(conn *net.Conn, data []byte) (n int, err error) {
	writer := bufio.NewWriter(*conn)
	n, err = writer.Write(data)
	if err != nil {
		return n, err
	}
	writer.Flush()
	return n, nil
}
