package main

import (
	c "./browser/connections"
	i "./browser/init"
	e "./error"
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
