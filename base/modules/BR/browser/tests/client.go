package main

import (
	e "../../error"
	"bufio"
	"fmt"
	"net"
	"net/http"
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
func ConnectionHandler(conn net.Conn, c chan *http.Response) {
	go func() {
		connReader := bufio.NewReader(conn)
		if request, err := http.ReadRequest(connReader); err == nil {
			//request.Body.Close()
			fmt.Println(c, "->", request, "\n")
			if response, err := http.DefaultClient.Get(request.URL.String()); err == nil {
				c <- response
			}
			/*if response, err := ReadFile("resp"); err == nil {
				c <- response
			}*/
		} else {
			return
		}
	}()

	/*go func() {
		connWriter := bufio.NewWriter(conn)
		if response, ok := <-c; ok {
			fmt.Println(string(response))
			e.ErrorHandler(err)
			connWriter.Write([]byte(response))
		} else {
			return
		}
	}()*/
	go func() {
		if response, ok := <-c; ok {
			fmt.Println("Found!")
			fmt.Println(c, "->", response, " \n")
			response.Write(conn)
			conn.Close()
			response.Body.Close()
			return

			// if response.StatusCode == http.StatusOK {
			// defer response.Body.Close()
			// scanner := bufio.NewScanner(response.Body)
			// scanner.Split(bufio.ScanBytes)
			// for scanner.Scan() {
			// fmt.Print(scanner.Text())
			// }
			// }

		} else {
			fmt.Println("Not found!")
			return
		}
	}()
}

func main() {
	tcpLn, err := net.Listen("tcp4", PORT)
	e.ErrorHandler(err)
	for {
		tcpConn, _ := tcpLn.Accept()
		e.ErrorHandler(err)
		c := make(chan *http.Response)
		go ConnectionHandler(tcpConn, c)
	}
}
