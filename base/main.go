package main

import (
	c "./modules/connection"
	e "./modules/error"
	i "./modules/init"
	k "./modules/keys"
	"fmt"
	// "net"
)

//SYSTEM GLOBAL VARIABLES
var (
	/*Config : variable that stores the configurations*/
	Config *i.Config

	/*Keys : variable that store the private and remote public keys*/
	Keys *k.Keys

	/*Conn : variable that stores the connection*/
	Conn *c.Connection

	//error variable to store all errors
	err error
)

func main() {

	/*Loaded the Configuration*/
	Config = i.NewConfig()
	err = Config.LoadConfig()
	e.ErrorHandler(err)

	/*Generate the Keys*/
	Keys = k.NewKeys()
	err = Keys.CreateRSAPair()
	e.ErrorHandler(err)

	/*Search Proxy Router*/
	Conn = c.NewConnection(Config.IP)
	err = Conn.OpenUDPPort()
	e.ErrorHandler(err)
	fmt.Println(Config.BCast)
	Conn.RemoteAddr, err = Conn.LaunchUDPTracer(Config.BCast)
	e.ErrorHandler(err)

	//c.LocalListen(config)
}

// "log"
// conn "./modules/connection"

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
