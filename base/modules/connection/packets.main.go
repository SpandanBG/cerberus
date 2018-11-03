package connection

import (
	"encoding/json"
)

type Header struct {
	Version int
	DREQ    bool
	REQ     bool
	RES     bool
	KX      bool
}

type Packet struct {
	Header byte   `json:"head"`
	Body   []byte `json:"body,omitempty"`
}

func GeneratePacket(header *Header, body []byte) ([]byte, error) {
	head := CreateHeaderByte(header)
	packet := Packet{Header: head, Body: body}
	jsonPacket, err := json.Marshal(packet)
	return jsonPacket, err
}

func ParserJSONPacket(jsonPacket []byte) (*Header, []byte, error) {
	packet := &Packet{}
	err := json.Unmarshal(jsonPacket, packet)
	head := TranslateHeaderByte(packet.Header)
	return head, packet.Body, err
}

func CreateHeaderByte(head *Header) byte {
	headByte := byte(head.Version) << 4
	if head.DREQ == true {
		headByte += byte(8)
	}
	if head.REQ == true {
		headByte += byte(4)
	}
	if head.RES == true {
		headByte += byte(2)
	}
	if head.KX == true {
		headByte += byte(1)
	}
	return headByte
}

func TranslateHeaderByte(headByte byte) *Header {
	head := &Header{}
	if int(headByte&1) == 1 {
		head.KX = true
	}
	if int(headByte&2) == 2 {
		head.RES = true
	}
	if int(headByte&4) == 4 {
		head.REQ = true
	}
	if int(headByte&8) == 8 {
		head.DREQ = true
	}
	head.Version = int(headByte >> 4)
	return head
}
