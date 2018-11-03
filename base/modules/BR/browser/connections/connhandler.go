package connections

import (
	e "../../error"
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
)

/*RemoteDial : Dialing a remote connection to the
router in order to fetch the response. Encrypted request sent to
the router and an encrypted response is received.
Request and Response are both in the form of []bytes*/
func RemoteDial(Request []byte) ([]byte, error) {
	var Response []byte
	var err error
	RIP := "192.168.43.67"
	RPort := "8080"
	TAddr := RIP + ":" + RPort
	if TConn, err := net.Dial("tcp", TAddr); err == nil {
		defer TConn.Close()
		if Response, err = GetRemoteResponse(TConn, Request); err == nil {
			return Response, nil
		}
		return nil, err
	}
	return nil, err
}

/*GetRemoteResponse : Fetches Remote Response from the Router
in an encrypted []byte response*/
func GetRemoteResponse(RConn net.Conn, RQ []byte) ([]byte, error) {
	var err error
	RQ = append(RQ, '\n')

	if _, err := RConn.Write(RQ); err == nil {
		fmt.Println("Written!")
		for {
			if RP, err := ioutil.ReadAll(RConn); err == nil {
				fmt.Println(RConn.LocalAddr(), "Read!", string(RP))
				if len(RP) > 0 {
					return RP, nil
				}
			}
		}
	}
	return nil, err
}

/*BrowserConnHandler : Read Write Requests and Responses
Browser send requests to the TCP Connections and the module
builds a goroutine for each such request. It, then, dials
a connection to the Proxy Router and sends an encrypted request
to the Router.
The encrypted Response received from the Router which is then
decrypted by the Middleware and served to the browser.
*/
func BrowserConnHandler(TConn net.Conn) {
	var RQ, RP []byte
	var err error
	DChan := make(chan []byte, 2048)
	EChan := make(chan []byte, 2048)

	TCReader := bufio.NewReader(TConn)
	//TCWriter := bufio.NewWriter(TConn)

	var RQJson *http.Request

	if RQJson, err = http.ReadRequest(TCReader); err == nil {
		//fmt.Println(TConn.RemoteAddr(), "->", RQJson)
		RQ, _ = httputil.DumpRequest(RQJson, true)
		for {
			select {
			case DChan <- RQ:
				if len(RQ) > 0 {
					RQ = make([]byte, 0)
					fmt.Println(TConn.RemoteAddr(), "->", RQJson)
					EChan = Middleware(DChan)
					if Req, ok := <-EChan; ok {
						//fmt.Println("Encrypted Request : ", string(Req))
						RP, err = RemoteDial(Req)
						e.ErrorHandler(err)
						//fmt.Println("Encrypted Response : ", RP)
						EChan <- RP
					}
				}
			default:
				DChan = Middleware(EChan)
				if RP, ok := <-DChan; ok {
					fmt.Println(TConn.RemoteAddr(), "<-", string(RP))
					n, err := TConn.Write(RP)
					fmt.Println("Bytes : ", n, ":", err)
					TConn.Close()
					close(DChan)
					close(EChan)
					return
				}
			}
		}
	}
}
