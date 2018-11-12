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
	CERBERUSHEADERSIZE = 244 //bytes
	// HTTPHEADERSIZE : Max HTTP header size
	HTTPHEADERSIZE = 16384 //bytes
	// ART :
	ART = `                         __    _
                    _wr""        "-q__
                 _dP                 9m_
               _#P                     9#_
              d#@                       9#m
             d##                         ###
            J###                         ###L
            {###K                       J###K
            ]####K      ___aaa___      J####F
        __gmM######_  w#P""   ""9#m  _d#####Mmw__
     _g##############mZ_         __g##############m_
   _d####M@PPPP@@M#######Mmp gm#########@@PPP9@M####m_
  a###""          ,Z"#####@" '######"\g          ""M##m
 J#@"             0L  "*##     ##@"  J#              *#K
 #"               ` + "`#    \"_gmwgm_~    dF               `" + `#_
7F                 "#_   ]#####F   _dK                 JE
]                    *m__ ##### __g@"                   F
                       "PJ#####LP"
 ` + "`" + `                       0######_                      '
                       _0########_
     .               _d#####^#####m__              ,
      "*w_________am#####P"   ~9#####mw_________w*"
          ""9@#####@M""           ""P@#####@M""                    `
)

var (
	// RSAKeys :
	RSAKeys *keys.Keys
)
