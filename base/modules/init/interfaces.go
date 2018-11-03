package init

import (
	"net"
	"strconv"
	"strings"
)

/*GetWLANInterface :  Gets the IP Address and BCast Address of the
current network dependning upon the platform. Everytime the client
reconnects to a new network the application needs to be restarted
for the WLAN Interface to get updated*/
func (C *Config) GetWLANInterface(Platform string) error {
	var WLAN string
	var err error

	switch Platform {
	case "linux":
		WLAN = "wlo1"
	case "windows":
		WLAN = "Wi-Fi"
	case "osx":
		WLAN = "en0"
	}

	if inf, err := net.InterfaceByName(WLAN); err == nil {
		if n, err := inf.Addrs(); err == nil {
			network := strings.Split(n[0].String(), "/")
			IP, Mask := network[0], network[1]
			BCast := ComputeBCast(IP, Mask)
			C.IP, C.BCast = IP, BCast
			return nil
		}
	}
	return err
}

/*ComputeBCast : Computes the Broadcast IP using the Subnet Mask
and the Remote IP Address*/
func ComputeBCast(IP, M string) string {
	MInt, _ := strconv.Atoi(M)
	Mask := net.CIDRMask(MInt, 32)
	IPstring := strings.Split(IP, ".")
	IPbits := make([]byte, 4)

	for i := range IPstring {
		IPint, _ := strconv.Atoi(IPstring[i])
		IPbits[i] = byte(IPint)
	}

	IPAddr := net.IP([]byte{IPbits[0], IPbits[1], IPbits[2], IPbits[3]})
	BCast := net.IP(make([]byte, 4))
	for i := range IPAddr {
		BCast[i] = IPAddr[i] | ^Mask[i]
	}
	return BCast.String()
}
