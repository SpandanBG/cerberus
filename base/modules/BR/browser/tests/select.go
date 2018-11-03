package main

import (
	"fmt"
)

func GetResponse(Pmsg chan []byte) chan []byte {
	if RP, ok := <-Pmsg; ok {
		Rmsg := make(chan []byte, 100)
		Rmsg <- RP
		return Rmsg
	}
	return nil
}

func main() {
	Pmsg := make(chan []byte, 100)
	Rmsg := make(chan []byte, 100)
	RQ := []byte("Hello")

	for {
		select {
		case Pmsg <- RQ:
			if len(RQ) > 0 {
				fmt.Println("Written to Channel")
				Rmsg = GetResponse(Pmsg)
				RQ = make([]byte, 0)
			}
		case RP := <-Rmsg:
			fmt.Println(string(RP))
			return
		}
	}
}
