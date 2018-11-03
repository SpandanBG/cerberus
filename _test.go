package main

import (
	"fmt"
	"./bin/keys"
)

func main() {
	hash := make([]byte,32)
	pkey:= make([]byte, 128)

	for i:=0; i<128; i++ {
		if i<32 {
			hash[i] = 'a'
		}
		pkey[i] = 'a'
	}

	pack := keys.NewExchangePacket()
	pack.ValidsFound = int8(^uint8(0)>>1)
	pack.KeyHash = string(hash)
	pack.PublicKey = string(pkey)

	json, _ := pack.GenerateJSONPacket()
	jsonByte := []byte(json)

	fmt.Println(json, len(jsonByte))
}