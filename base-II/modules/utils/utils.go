package utils

import (
	"bytes"
	"net"
	"runtime"
	"strings"

	"github.com/ugorji/go/codec"
)

// GetWiFiIPAddr : Retives the IP address of the connected WiFi network
//
// Returns :
//		string : The IP address as string
//		error  : The error occured during retival
func GetWiFiIPAddr() (addr string, err error) {
	var platform string
	var netInterface *net.Interface
	var interfaceAddrs []net.Addr
	switch runtime.GOOS {
	case "windows":
		platform = "Wi-Fi"
	case "linux":
		platform = "wlan0"
	}
	netInterface, err = net.InterfaceByName(platform)
	if err == nil {
		interfaceAddrs, err = netInterface.Addrs()
		if err == nil {
			firstAddr := interfaceAddrs[0]
			addr = strings.Split(firstAddr.String(), "/")[0]
		}
	}
	return
}

// CBOREncode : Encodes the input interface to bytes
func CBOREncode(input interface{}) ([]byte, error) {
	var ch codec.CborHandle
	var buffer bytes.Buffer
	enc := codec.NewEncoder(&buffer, &ch)
	err := enc.Encode(input)
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// CBORDecode : Decodes the input bytes to the interface type
func CBORDecode(input []byte, output interface{}) error {
	var ch codec.CborHandle
	buffer := bytes.NewBuffer(input)
	dec := codec.NewDecoder(buffer, &ch)
	return dec.Decode(output)
}
