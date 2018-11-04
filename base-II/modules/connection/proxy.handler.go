package connection

import (
	"../configs"
	"../utils"
)

// ProxyHandler : Handles request packets
func ProxyHandler(reqRaw []byte) (resRaw []byte, err error) {
	var packet Packet
	err = utils.GOBDecode(reqRaw, packet)
	if err == nil {

	}
	return
}

// IsHttpRequestPacket : Checks if the packet is a HTTP request packet
func IsHttpRequestPacket(header *Header) bool {
	required := header.Version == configs.VERSION
	required = required && header.REQ
	notRequired := header.DREQ && header.RES && header.KX
	return required && !notRequired
}

// IsKeyExchangePacket : Checks if the packet is a key exchange packet
func IsKeyExchangePacket(header *Header) bool {
	required := header.Version == configs.VERSION
	required = required && header.REQ && header.KX
	notRequired := header.DREQ && header.RES
	return required && !notRequired
}
