package configs

import (
	"../keys"
)

const (
	// VERSION :
	VERSION = 1
	// PORT :
	PORT = "4123"
	// CERBERUSHEADERSIZE : Cerberus packet size
	CERBERUSHEADERSIZE = 129 //bytes
	// HTTPHEADERSIZE : Max HTTP header size
	HTTPHEADERSIZE = 16384 //bytes
)

var (
	// RSAKeys :
	RSAKeys *keys.Keys
)