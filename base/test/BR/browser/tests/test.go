package main

import (
	e "../../error"
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	//PORT to be used
	PORT = ":8085"
)

var response, request []byte
var err error

/*ReadFile : To read the response from file "response.txt" on localhost.
  Serve the response to the requesting sockets.
**/
func ReadFile(fileName string) (resp []byte, err error) {
	file, err := os.Open(fileName)
	e.ErrorHandler(err)

	resp = make([]byte, 2048)
	_, err = file.Read(resp)
	return
}

/*ConnectionHandler : Thread to handle each connection separately
  Reads requests from each incoming connections
  Reads the response from the response file "response.txt"
  Sends the response to the connections.
**/
func ConnectionHandler(conn net.Conn, c chan []byte) {
	go func() {
		connReader := bufio.NewReader(conn)
		for {
			if request, err := connReader.ReadBytes('\n'); err == nil {
				fmt.Println("Reader Request from", conn.RemoteAddr(), " : ", c, " -- ", string(request))
				if response, err := ReadFile("resp"); err == nil {
					c <- response
				}
			} else {
				return
			}
		}
	}()

	go func() {
		connWriter := bufio.NewWriter(conn)
		for {
			if response, ok := <-c; ok {
				fmt.Println(string(response))
				e.ErrorHandler(err)
				connWriter.Write([]byte(response))
			} else {
				return
			}
		}
	}()
}

/*func main() {
	tcpLn, err := net.Listen("tcp4", PORT)
	e.ErrorHandler(err)
	for {
		tcpConn, _ := tcpLn.Accept()
		e.ErrorHandler(err)
		c := make(chan []byte)
		go ConnectionHandler(tcpConn, c)
	}
}*/
