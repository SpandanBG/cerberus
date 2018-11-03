package main

import (
	"strconv"
	// "bytes"
	"fmt"
	// cbor "github.com/2tvenom/cbor"
	"net"
	"strings"
)

func main() {
	// if intr, err := net.InterfaceByName("wlo1"); err == nil {
	// 	fmt.Println(intr.Addrs())
	// }
	// os.Setenv("HTTP_PROXY", "127.0.0.1:90")
	// fmt.Println(os.Getenv("HTTP_PROXY"))
	// for _, e := range os.Environ() {
	// 	pair := strings.Split(e, "=")
	// 	fmt.Println(pair[0])

	nw := "127.0.0.1/20"
	nwsplit := strings.Split(nw, "/")
	fmt.Println(nwsplit)
	msk, _ := strconv.Atoi(nwsplit[1])
	ip := strings.Split(nwsplit[0], ".")
	ipbits := make([]byte, 4)
	for i := range ip {
		ipint, _ := strconv.Atoi(ip[i])
		ipbits[i] = byte(ipint)
	}
	mask1 := net.CIDRMask(20, 32)
	fmt.Println(mask1)
	mask := net.CIDRMask(msk, 32)
	ipmask := net.IP([]byte{ipbits[0], ipbits[1], ipbits[2], ipbits[3]})
	fmt.Println(mask, ":", ipmask)
	broadcast := net.IP(make([]byte, 4))
	for i := range ipmask {
		broadcast[i] = ipmask[i] | ^mask[i]
		fmt.Println(i)
	}
	fmt.Println(broadcast)
}
