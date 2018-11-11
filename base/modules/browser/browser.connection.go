package browser

import (
	e "../error"
	i "../init"
	k "../keys"
	u "../utils"
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
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
	//fmt.Println(RAddr)
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
		if nn, err := Reader.Read(Response); err == nil {
			fmt.Println(nn, "Bytes Read!")
			if len(Response) > 0 {
				return Response, nil
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
func BrowserConnHandler(B BrowserConn, K *k.Keys, C i.Config) {
	var RQ, RP []byte
	var Request *http.Request
	var temp *http.Request
	var err error
	RemoteAddr := C.RemoteAddr

	TCReader := bufio.NewReader(B.Conn)
	Request, err = http.ReadRequest(TCReader)

	RQBytes, err := httputil.DumpRequest(Request, true)
	e.ErrorHandler(err)
	RQ, err = u.GOBEncode(RQBytes)
	e.ErrorHandler(err)
	err = u.GOBDecode(RQ, &temp)
	RQHeader := InitPacketHeader(&B, C.Version, K.PublicKey)
	ReqBytes := EncryptChannel(RQHeader, RQ, B.RQPacket, K)
	RP, err = RemoteDial(ReqBytes, RemoteAddr)
	e.ErrorHandler(err)
	Response := DecryptChannel(RP, B.RSPacket, K)
	fmt.Println(string(Response))
	B.CloseBrowserConn()
	return
}
