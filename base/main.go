package main

import (
	b "./modules/browser"
	c "./modules/connection"
	e "./modules/error"
	i "./modules/init"
	k "./modules/keys"
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

/*DoConfig : Configuration Initialization*/
func DoConfig() {
	Config = i.NewConfig()
	err = Config.LoadConfig()
	e.ErrorHandler(err)
}

/*DoKeys : RSA Keys Initialization*/
func DoKeys() {
	Keys = k.NewKeys()
	err = Keys.CreateRSAPair()
	e.ErrorHandler(err)
	SearchProxyRouter()
	e.ErrorHandler(err)

	Keys.GetRemotePublicKey(Conn)
}

/*SearchProxyRouter : Establish connection to Proxy Router*/
func SearchProxyRouter() {
	Conn = c.NewConnection(Config.Host + Config.Port)
	err = Conn.OpenUDPPort()
	e.ErrorHandler(err)

	//fmt.Println(Config.BCast)
	Conn.RemoteAddr, err = Conn.LaunchUDPTracer(Config.IP, Config.BCast+":"+Config.Port)
	Config.RemoteAddr = Conn.RemoteAddr.String()
	e.ErrorHandler(err)
}

func main() {
	DoConfig()
	DoKeys()
	b.LocalListen(Config, Keys)
}
