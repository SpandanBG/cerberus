package connections

import (
	e "../../error"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
)

/*ReadFile : To read the response from file "response.txt" on localhost.
  Serve the response to the requesting sockets.
**/
func ReadFile(fileName string) (RP []byte, err error) {
	file, err := os.Open(fileName)
	e.ErrorHandler(err)

	RP = make([]byte, 2048)
	_, err = file.Read(RP)
	return
}

/*RouterConnHandler : Read Write Requests and Responses from
and to the client. The encrypted Response received from the Client which is then
decrypted by the Middleware, read by the Router and finally
sends encrypted response to the TCP Connections and the module
builds a goroutine for each such incoming request*/
func RouterConnHandler(TConn net.Conn) {
	var RQ, RP []byte
	var err error
	DChan := make(chan []byte, 2048)
	EChan := make(chan []byte, 2048)
	TCReader := bufio.NewReader(TConn)

	var RQJson *http.Request
	var RPJson *http.Response
	if RQ, err = TCReader.ReadBytes('\n'); err == nil {
		RQ = RQ[:len(RQ)-1]
		fmt.Println(TConn.RemoteAddr(), "->", string(RQ))
		e.ErrorHandler(err)
		for {
			select {
			case EChan <- RQ:
				if len(RQ) > 0 {
					RQ = make([]byte, 0)
					DChan = Middleware(EChan)
					if Req, ok := <-DChan; ok {
						RP, err = ReadFile("router/response.txt")
						e.ErrorHandler(err)
						DChan <- RP
					} else {
						fmt.Println("Empty Channel")
					}
				}
			default:
				EChan = Middleware(DChan)
				if RP, ok := <-EChan; ok {
					fmt.Println(TConn.RemoteAddr(), "<-", string(RP))
					RP = append(RP, '\n')
					n, err := TConn.Write(RP)
					fmt.Println("BYtes Written : ", n, ":", err)
					TConn.Close()
					close(DChan)
					close(EChan)
					return
				}
			}
		}
	}
}
