package trafficcapture

//IPSessions : ip session info
type IPSessions struct {
	Hash          uint64
	Size          int64
	Send          int64
	Recieve       int64
	SizeString    string
	SendString    string
	RecieveString string
	Src           string
	Dst           string
	LastSeen      int64
}

type SourceGroup struct {
	Hash        uint64
	SourceIP    string
	PacketCount int64
	SendByte    int64
	RecieveByte int64
}

//IPPacket : ip packet info
type IPPacket struct {
	Hash            uint64
	PacketSize      int64
	Src             string
	Dst             string
	Protocol        string
	SrcPort         int
	DstPort         int
	TransParentHash uint64
}
