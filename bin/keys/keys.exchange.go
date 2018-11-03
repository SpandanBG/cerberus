package keys

/*
	Algorithm of the Key Exchange

	The struct of the packet is as such:
	{
		type: "REQ", 			// or RES
		hashkey: "xxxx"			// or "0000" if no hashkey
		[, publickey: "asdf"]	// optional entry
	}

	Initiator Side Steps:
		1. Send REQ packet.
			The packet should contain the packet type, the hashkey and possibly the public key.
			Initially no public key is to be appended.
		2. Receive RES packet.
		3. Prepare a new REQ packet
		4. If RES packet `hashkey` IS NOT VALID:
			4-1. Add `publickey` to new REQ packet
			4-2. Mark new REQ packet to be sent.
		5. If RES packet HAS `publickey`:
			5-1. Update public key to database.
			5-2. Create hash key of the public key.
			5-3. Add new hash key to the new REQ packet.
			5-4. Mark new REQ packet to be sent.
		6. If new REQ packet marked to be sent:
			6-1. Go to 1.
		7. End.

	Acknowlegor Side Steps:
		1. Receive REQ packet.
		2. Prepares a new RES packet.
		3. If REQ packet `hashkey` IS NOT VALID:
			3-1. Add `publickey` to new RES packet.
		4. If REQ packet HAS `publickey`:
			4-1. Update public key to database.
			4-2. Create hash key of the public key.
			4-3. Add hash key to the new RES packet.
		5. Send RES packet.
		6. End.

*/

// REQ and RES -> [strings] : Specifies request or response packet
// dHASH -> [string] : The default hashkey if not present
const (
	REQ   = "REQ"
	RES   = "RES"
	dHASH = "0000"
)

// JSONPacketStruct - The structure of the key exchange packet
//
// Type -> [string] : Describes type of the packet. REQ or RES.
// Hashkey -> [string] : Contains the hashkey to be verified.
// Publickey -> [string] : Contains the publickey of self if required.
type JSONPacketStruct struct {
	Type      string `json:"TYP"`
	Hashkey   string `json:"HK"`
	Publickey string `json:"PK,omitempty"`
}

// KeyCommunication - The methods available to perfom the key exchange.
type KeyCommunication interface {
	NewExchangePacket() *JSONPacketStruct
	ProcessPacket() (*JSONPacketStruct, error)
}

// NewExchangePacket - Creates a new packet with the passed values.
//
// Type -> [string] : The type of the packet to be sent. If "" then defaults to REQ.
// Hashkey -> [string] : The hashkey to be verified. If "" then defaults to 0000.
// Publickey -> [string] : The publickey to be sent. "" if need not be sent.
func NewExchangePacket(Type string, Hashkey string, Publickey string) *JSONPacketStruct {
	if Type == "" {
		Type = REQ
	}
	if Hashkey == "" {
		Hashkey = dHASH
	}
	return &JSONPacketStruct{
		Type:      Type,
		Hashkey:   Hashkey,
		Publickey: Publickey,
	}
}

// ProcessPacket - Performs processing of the packet. Returns a new packet if new exchange is required.
//
// * -> [*JSONPacketStruct] : Nil if no new exchange is to be done. Otherwise contain a new packet.
// * -> [error] : Provides the error occured during the processing. Nil if packet was correct.
func (packet *JSONPacketStruct) ProcessPacket() (*JSONPacketStruct, error) {
	var err error
	newPacket := &JSONPacketStruct{}
	newPacket.setAltREQorRES(packet.Type)

	return newPacket, err
}

// setAltREQorRES - Sets alternate of REQ or RES from the given type in the parameter
//
// Type -> [string] : The type from which to alternate
func (packet *JSONPacketStruct) setAltREQorRES(Type string) {
	if Type == REQ {
		packet.Type = RES
	} else if Type == RES {
		packet.Type = RES
	}
}
