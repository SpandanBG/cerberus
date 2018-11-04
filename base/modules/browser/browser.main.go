package browser

import (
	p "../connection"
	e "../error"
	i "../init"
	k "../keys"
	"net"
)

type BrowserConn struct {
	Conn               net.Conn
	RQPacket, RSPacket *p.Packet
	EChan, DChan       chan []byte
}

/*NewBrowserConn : creates a new browser connection object*/
func NewBrowserConn(Conn *net.Conn) *BrowserConn {
	return &BrowserConn{
		Conn:  *Conn,
		EChan: make(chan []byte),
		DChan: make(chan []byte),
	}
}

/*LocalListen : Listens on the configured Local IP and Port
for all incoming requests from the browser. HTTP_PROXY.env is set to
the configured Host+Port in order to re-route all the HTTP browser requests
to a that Address*/
func LocalListen(config *i.Config, keys *k.Keys) {
	LocalAddr := config.Host + ":" + config.Port
	Listener, err := net.Listen("tcp", LocalAddr)
	e.ErrorHandler(err)
	for {
		Conn, err := Listener.Accept()
		e.ErrorHandler(err)
		BC := NewBrowserConn(&Conn)
		go BrowserConnHandler(*BC, keys, *config)
	}
}

/*CloseBrowserConn: closes all channels and connections in the
BrowserConn object*/
func (B *BrowserConn) CloseBrowserConn() {
	B.Conn.Close()
	close(B.EChan)
	close(B.DChan)
}
