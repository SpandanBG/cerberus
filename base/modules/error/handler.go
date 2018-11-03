package error

import (
	"fmt"
	"os"
)

/*ErrorHandler : Handles all possible errors that occur
  Responds accordingly.
  Exits if err exists else returns
**/
func ErrorHandler(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(0)
	}
}
