package init

import (
	e "../../error"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

//SYSTEM CONSTANTS
const (
	//ConfigPath : Path to configuration file
	ConfigPath = "router/config.json"
)

/*CONFIG : loads the required configuration for the system
Several information about the system is stored in this structure
Any changes to be made to the configuration needs to be made in
file named as "config.json"
*/
type CONFIG struct {
	Version     string `json:"version"`
	Name        string `json:"name"`
	Protocol    string `json:"protocol"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Platform    string `json:"platform"`
	Environment string `json:"environment"`
}

/*LoadConfig : To load the cofiguration file into the structure CONFIG*/
func LoadConfig() (*CONFIG, error) {
	var config *CONFIG
	jsonFile, err := os.OpenFile(ConfigPath, os.O_RDONLY, 0666)
	e.ErrorHandler(err)
	defer jsonFile.Close()

	jsonValue, err := ioutil.ReadAll(jsonFile)
	e.ErrorHandler(err)

	if err = json.Unmarshal(jsonValue, &config); err == nil {
		config.Platform = runtime.GOOS
		DisplayConfigDetails(config)
		return config, nil
	}
	return nil, err
}

/*DisplayConfigDetails : Displays details of the configuration*/
func DisplayConfigDetails(C *CONFIG) {
	fmt.Println("=> Booting Cerberus")
	fmt.Println("=> Cerberus loading in the single mode ...")
	fmt.Println("* Version :", C.Version, ", Codename :", C.Name)
	fmt.Println("* Platform :", C.Platform, ", Environment :", C.Environment)
	Addr := C.Protocol + "://" + C.Host + ":" + C.Port
	fmt.Println("Listening on ", Addr)
	fmt.Println("Use Ctrl-C to stop")
}
