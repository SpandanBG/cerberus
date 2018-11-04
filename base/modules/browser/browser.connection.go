package browser

import (
	e "../error"
	k "../keys"
	u "../utils"
	"bufio"
	"fmt"
	"net"
	"net/http"
)

/*RemoteDial : Dialing a remote connection to the
router in order to fetch the response. Encrypted request sent to
the router and an encrypted response is received.
Request and Response are both in the form of []bytes*/
func RemoteDial(Request []byte, RAddr string) ([]byte, error) {
	var Response []byte
	var err error

	if TConn, err := net.Dial("tcp", RAddr); err == nil {
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
func GetRemoteResponse(RConn net.Conn, Request []byte) ([]byte, error) {
	var err error
	var Response []byte

	Reader := bufio.NewReader(RConn)
	Writer := bufio.NewWriter(RConn)
	if n, err := Writer.Write(Request); err == nil {
		Writer.Flush()
		fmt.Println(n, "Bytes Written!")
		for {
			if nn, err := Reader.Read(Response); err == nil {
				fmt.Println(nn, "Bytes Read!")
				if len(Response) > 0 {
					return Response, nil
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
func BrowserConnHandler(B BrowserConn, K k.Keys, RemoteAddr string) {
	var RQ, RP []byte
	var Request *http.Request
	var err error

	TCReader := bufio.NewReader(*B.Conn)
	Request, err = http.ReadRequest(TCReader)
	e.ErrorHandler(err)
	RQ, err = u.GOBEncode(Request)
	e.ErrorHandler(err)

	for {
		select {
		case B.DChan <- RQ:
			if len(RQ) > 0 {
				RQ = make([]byte, 0)
				B.EChan = EncryptChannel(B.DChan, B.RQPacket, &K)
				if ReqBytes, ok := <-B.EChan; ok {
					RP, err = RemoteDial(ReqBytes, RemoteAddr)
					e.ErrorHandler(err)
					B.EChan <- RP
				}
			}
		default:
			B.DChan = DecryptChannel(B.EChan, B.RSPacket, &K)
			if RespBytes, ok := <-B.DChan; ok {
				//n, err := TConn.Write(RP)
				B.CloseBrowserConn()
				return
			}
		}
	}
}
