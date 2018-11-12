package init

import (
	e "../error"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

//SYSTEM CONSTANTS
const (
	//ConfigPath : Path to configuration file
	ConfigPath = "modules/config.json"

	//CerberuArt
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

/*Config : loads the required configuration for the system
Several information about the system is stored in this structure
Any changes to be made to the configuration needs to be made in
file named as "config.json"
*/
type Config struct {
	Version     int    `json:"version"`     //Version of Protocol
	Name        string `json:"name"`        //Name of Protocol
	Protocol    string `json:"protocol"`    //Underlying Protocol
	Host        string `json:"host"`        //Localhost Loopback Address
	Port        string `json:"port"`        //Local Port
	Platform    string `json:"platform"`    //OS Platform
	Environment string `json:"environment"` //Environment : Development, Test, Production, etc.
	IP          string `json:"ip"`          //IP Address of the Local Machine over the Network
	BCast       string `json:"bcast"`       //Broadcast Address of Network
	RemoteAddr  string `json:"raddr"`       //Remote Address of Router
}

/*NewConfig : a Config struct is returned*/
func NewConfig() *Config {
	return &Config{}
}

/*LoadConfig : To load the cofiguration file into the structure Config*/
func (config *Config) LoadConfig() error {
	jsonFile, err := os.OpenFile(ConfigPath, os.O_RDONLY, 0666)
	e.ErrorHandler(err)
	defer jsonFile.Close()

	jsonValue, err := ioutil.ReadAll(jsonFile)
	e.ErrorHandler(err)

	if err = json.Unmarshal(jsonValue, &config); err == nil {
		config.Platform = runtime.GOOS
		if err = config.GetWLANInterface(config.Platform); err == nil {
			DisplayConfigDetails(config)
			return nil
		}
	}
	return err
}

/*DisplayConfigDetails : Displays details of the configuration*/
func DisplayConfigDetails(C *Config) {
	fmt.Println("=> Booting Cerberus")
	fmt.Println("=> Cerberus loading in the single mode ...")
	fmt.Println("* Version :", C.Version, ", Codename :", C.Name)
	fmt.Println("* Platform :", C.Platform, ", Environment :", C.Environment)
	Addr := C.Protocol + "://" + C.Host + ":" + C.Port
	fmt.Println("Listening on ", Addr)
	fmt.Println("IP Address: ", C.IP, " Broadcasting to ", C.BCast)
	fmt.Println("Use Ctrl-C to stop")
	fmt.Println("---------------------------------")
	fmt.Println(ART, "\n")
}
