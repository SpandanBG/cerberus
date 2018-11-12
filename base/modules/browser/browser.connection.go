package browser

import (
	c "../connection"
	e "../error"
	i "../init"
	k "../keys"
	u "../utils"
	"bufio"
	"bytes"
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
	Reader := bufio.NewReader(RConn)
	Writer := bufio.NewWriter(RConn)
	if _, err := Writer.Write(Request); err == nil {
		Writer.Flush()
		var resBuffer bytes.Buffer
		response := make([]byte, 20)
		for {
			n, err := Reader.Read(response)
			if err != nil {
				break
			}
			resBuffer.Write(response[:n])
		}
		if resBuffer.Len() > 0 {
			return resBuffer.Bytes(), nil
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
	var output c.PseudoResponse
	var err error
	RemoteAddr := C.RemoteAddr

	TCReader := bufio.NewReader(B.Conn)
	Request, _ = http.ReadRequest(TCReader)
	PRequest, err := c.RequestToPseudoRequest(Request)
	e.ErrorHandler(err)

	RQ, err = u.GOBEncode(PRequest)
	e.ErrorHandler(err)
	RQHeader := InitPacketHeader(&B, C.Version, K.PublicKey)
	ReqBytes := EncryptChannel(RQHeader, RQ, B.RQPacket, K)
	RP, err = RemoteDial(ReqBytes, RemoteAddr)
	e.ErrorHandler(err)
	RP = DecryptChannel(RP, K)
	err = u.GOBDecode(RP, &output)
	_, err = c.PseudoResponseToResponse(&output)
	fmt.Println("Recved:", len(output.Body))
	e.ErrorHandler(err)
	B.CloseBrowserConn()
}
