package main

// "log"

// conn "./modules/connection"

const (
	BrtIP = "192.168.255.255"
	IP    = "192.168.1.105"
	Port  = 4123
)

// var RouterAddr *conn.Address

func init() {
	// searchProxyRouter()
	// Load RSA Keys
}

func main() {

}

// func searchProxyRouter() {
// 	connect := conn.NewConnection(conn.Address{IP: IP, Port: Port})
// 	err := connect.OpenUDPPort()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	RouterAddr, err = connect.LaunchUDPTracer(conn.Address{IP: BrtIP, Port: Port})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
