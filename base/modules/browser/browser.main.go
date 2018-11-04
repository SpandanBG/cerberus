package browser

import (
	e "../../error"
	p "../connection"
	b "../init"
	"crypto/rsa"
	"net"
)

type BrowserConn struct {
	Conn   *net.Conn
	Packet *p.Packet
	EChan  chan []byte
	DChan  chan []byte
}

/*NewBrowserConn : creates a new browser connection object*/
func NewBrowserConn(Conn *net.Conn) *BrowserConn {
	return &BrowserConn{
		Conn:  Conn,
		EChan: make(chan []byte, 2048),
		DChan: make(chan []byte, 2048)
	}
}

/*LocalListen : Listens on the configured Local IP and Port
for all incoming requests from the browser. HTTP_PROXY.env is set to
the configured Host+Port in order to re-route all the HTTP browser requests
to a that Address*/
func LocalListen(config *b.CONFIG, keys *Keys) {
	LocalAddr := config.Host + ":" + config.Port
	RemoteAddr := config.RemoteAddr + ":" + config.Port
	Listener, err := net.Listen("tcp", LocalAddr)
	e.ErrorHandler(err)
	for {
		Conn, err := Listener.Accept()
		e.ErrorHandler(err)
		BC := NewBrowserConn(Conn)
		BC.InitPacketHeader(config.Version, keys.PublicKey)
		go BrowserConnHandler(BC)
	}
}

/*InitPacketHeader : intialize packet header*/
func (b *BrowserConn) InitPacketHeader(Version int, PK *rsa.PublicKey) {
	header := p.Header{Version: Version, REQ: true}
	b.Packet.Header = header
	b.Packet.Key = PK
}
