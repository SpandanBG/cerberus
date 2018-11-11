package utils

import (
	"bytes"
	codec "github.com/ugorji/go/codec"
	"net"
	"runtime"
	"strings"
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

// GOBEncode : Encodes the input interface to bytes
func GOBEncode(input interface{}) ([]byte, error) {
	/*var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(input)
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil*/
	var ch codec.CborHandle
	var buffer bytes.Buffer
	encoder := codec.NewEncoder(&buffer, &ch)
	err := encoder.Encode(input)
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

// GOBDecode : Decodes the input bytes to the interface type
func GOBDecode(input []byte, output interface{}) error {
	/*var buffer bytes.Buffer
	buffer.Write(input)
	dec := gob.NewDecoder(&buffer)
	return dec.Decode(output)*/
	var ch codec.CborHandle
	buffer := bytes.NewBuffer(input)
	decoder := codec.NewDecoder(buffer, &ch)
	err := decoder.Decode(&output)
	return err
}
