package trafficaggrigate

import (
	"log"
	"net"
	"sync"
	"wireguard-web-ui/trafficcapture"
)

type sourceIPs struct {
	aggrigatorQueue   chan trafficcapture.IPPacket
	sourceIPTable     map[string]*trafficcapture.SourceGroup
	sourceIPTableLock sync.RWMutex
	dataWriter        IDataWriter
}

func (srcIPs *sourceIPs) Aggrigate(packet trafficcapture.IPPacket) {
	srcIPs.aggrigatorQueue <- packet
}

//-------------------------------------------------
//                Private functions
//-------------------------------------------------

func (srcIPs *sourceIPs) createSourceTables() {
	for packet := range srcIPs.aggrigatorQueue {
		privateIPKey := ""
		if ipRangePrivate(net.ParseIP(packet.Src)) {
			privateIPKey = packet.Src
		} else if ipRangePrivate(net.ParseIP(packet.Dst)) {
			privateIPKey = packet.Dst
		} else {
			log.Fatalln("cant determine send or recieve : ", packet)
		}
		// update source group table
		srcIPs.sourceIPTableLock.Lock()
		srcIP, ok := srcIPs.sourceIPTable[privateIPKey]
		if ok {
			srcIPs.sourceIPTable[privateIPKey].PacketCount++
			if packet.Src == srcIP.SourceIP {
				srcIPs.sourceIPTable[privateIPKey].SendByte += packet.PacketSize
			}
			if packet.Dst == srcIP.SourceIP {
				srcIPs.sourceIPTable[privateIPKey].RecieveByte += packet.PacketSize
			}
		} else {
			if privateIPKey == packet.Src {
				srcIPs.sourceIPTable[privateIPKey] = &trafficcapture.SourceGroup{
					SourceIP:    privateIPKey,
					PacketCount: 1,
					SendByte:    packet.PacketSize,
					RecieveByte: 0,
				}
			} else if privateIPKey == packet.Dst {
				srcIPs.sourceIPTable[privateIPKey] = &trafficcapture.SourceGroup{
					SourceIP:    privateIPKey,
					PacketCount: 1,
					SendByte:    0,
					RecieveByte: packet.PacketSize,
				}
			} else {
				log.Fatalln("cant determine send or recieve : ", packet)
			}
		}
		srcIPs.sourceIPTableLock.Unlock()
	}
}

func (srcIPs *sourceIPs) writeAndFlush() {
	srcIPs.sourceIPTableLock.Lock()
	srcIPs.dataWriter.WriteSourceGroup(srcIPs.sourceIPTable)
	srcIPs.sourceIPTableLock.Unlock()
	srcIPs.sourceIPTable = make(map[string]*trafficcapture.SourceGroup)
}
