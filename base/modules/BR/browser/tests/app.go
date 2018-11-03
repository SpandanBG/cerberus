package main

import (
	"fmt"
	"io"
	"net/http"
)

/*HandlerRR : jkasd*/
func HandlerRR(w http.ResponseWriter, req *http.Request) {
	client := &http.Client{}
	resp, err := client.Do(req)
	io.Copy(w, resp.Body)
	resp.Body.Close()
}
func main() {

}
