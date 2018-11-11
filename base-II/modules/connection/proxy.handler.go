package connection

import (
	"crypto/rsa"
	"fmt"
	"net/http"

	"../configs"
	"../utils"
)

// ProxyHandler : Handles request packets
func ProxyHandler(reqRaw []byte, rAddr string) (resRaw []byte, err error) {
	header, pubkey, body, err := ParserPacketBytes(reqRaw)
	if err == nil {
		if IsHTTPRequestPacket(header) {
			fmt.Println("Proxy Request From :", rAddr)
			httpResponse, err := MakeHTTPRequest(body)
			if err == nil {
				resRaw, err = GenerateHTTPPacket(httpResponse, pubkey)
			}
		} else if IsKeyExchangePacket(header) {
			fmt.Println("Public Key Request From :", rAddr)
			resRaw, err = GeneratePublicKeyPacket()
		}
	}
	return
}

// MakeHTTPRequest : Makes the http request and returns the response
func MakeHTTPRequest(body []byte) (resRaw []byte, err error) {
	var reqHeader *http.Request
	var response *http.Response
	var presponse *PseudoResponse
	reqHeader, err = BodyToHTTPRequest(body)
	if err == nil {
		response, err = http.DefaultClient.Do(reqHeader)
		if err == nil {
			presponse, err = ResponseToPseudoResponse(response)
			if err == nil {
				resRaw, err = utils.CBOREncode(presponse)
			}
		}
	}
	return
}

// BodyToHTTPRequest : Converts packet body to http.Request
func BodyToHTTPRequest(body []byte) (*http.Request, error) {
	requestRaw, err := configs.RSAKeys.Decrypt(body)
	if err != nil {
		return nil, err
	}
	var prequest PseudoRequest
	err = utils.CBORDecode(requestRaw, &prequest)
	return PseudoRequestToRequest(&prequest)
}

// GeneratePublicKeyPacket : Creates the public key packet
func GeneratePublicKeyPacket() ([]byte, error) {
	header := &Header{Version: configs.VERSION, REQ: true, RES: true, KX: true}
	key := configs.RSAKeys.PublicKey
	body := []byte{}
	return GeneratePacket(header, key, body)
}

// GenerateHTTPPacket : Creates the Http response packet
func GenerateHTTPPacket(body []byte, pubkey *rsa.PublicKey) ([]byte, error) {
	header := &Header{Version: configs.VERSION, REQ: true, RES: true}
	encBody, err := configs.RSAKeys.Encrypt(pubkey, body)
	if err != nil {
		return []byte{}, err
	}
	return GeneratePacket(header, nil, encBody)
}

// IsHTTPRequestPacket : Checks if the packet is a HTTP request packet
func IsHTTPRequestPacket(header *Header) bool {
	required := header.Version == configs.VERSION
	required = required && header.REQ
	notRequired := header.DREQ || header.RES || header.KX
	return required && !notRequired
}

// IsKeyExchangePacket : Checks if the packet is a key exchange packet
func IsKeyExchangePacket(header *Header) bool {
	required := header.Version == configs.VERSION
	required = required && header.REQ && header.KX
	notRequired := header.DREQ || header.RES
	return required && !notRequired
}
