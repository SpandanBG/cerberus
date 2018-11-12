package browser

import (
	e "../error"
	i "../init"
	k "../keys"
	"io"
	"net"
	"net/http"
)

/*LocalListen : Listens on the configured Local IP and Port
for all incoming requests from the browser. HTTP_PROXY.env is set to
the configured Host+Port in order to re-route all the HTTP browser requests
to a that Address*/
func LocalListen(config *i.Config, keys *k.Keys) {
	LocalAddr := config.Host + ":" + config.Port
	Listener, err := net.Listen("tcp", LocalAddr)
	e.ErrorHandler(err)
	http.Serve(Listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodConnect {
			response := BrowserConnHandler(keys, *config, r)
			defer response.Body.Close()
			for k, vv := range response.Header {
				for _, v := range vv {
					w.Header().Add(k, v)
				}
			}
			w.WriteHeader(response.StatusCode)
			io.Copy(w, response.Body)
		}
	}))
}
