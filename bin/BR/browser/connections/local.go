package connections

import (
	e "../../error"
	b "../init"
	"net"
)

/*LocalListen : Listens on the configured Local IP and Port
for all incoming requests from the browser. HTTP_PROXY.env is set to
the configured Host+Port in order to re-route all the HTTP browser requests
to a that Address*/
func LocalListen(config *b.CONFIG) {
	LAddr := config.Host + ":" + config.Port
	Ln, err := net.Listen("tcp", LAddr)
	e.ErrorHandler(err)
	for {
		Conn, err := Ln.Accept()
		e.ErrorHandler(err)
		go BrowserConnHandler(Conn)
	}
}
