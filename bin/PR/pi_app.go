package main

import (
	e "./error"
	c "./router/connections"
	i "./router/init"
)

//SYSTEM GLOBAL VARIABLES
var (
	//config variable that stores the configurations
	config *i.CONFIG
	//error variable to store all errors
	err error
)

func main() {
	config, err = i.LoadConfig()
	e.ErrorHandler(err)
	c.LocalListen(config)
}
