package trafficaggrigate

import (
	"fmt"
	"time"
	"wireguard-web-ui/trafficcapture"
)

const (
	SourceGroup = "source-group"
)

type IDataWriter interface {
	WriteSourceGroup(map[string]*trafficcapture.SourceGroup)
}

func New(aggName string, writeInterval int, dataWriter IDataWriter) (trafficcapture.IPacketProcessor, error) {
	if aggName == SourceGroup {
		sourceGroupAgg := sourceIPs{}
		sourceGroupAgg.sourceIPTable = make(map[string]*trafficcapture.SourceGroup)
		sourceGroupAgg.aggrigatorQueue = make(chan trafficcapture.IPPacket, 100000)
		sourceGroupAgg.dataWriter = dataWriter
		go sourceGroupAgg.createSourceTables()
		go func() {
			for {
				time.Sleep(time.Duration(writeInterval) * time.Second)
				sourceGroupAgg.writeAndFlush()
			}
		}()
		return &sourceGroupAgg, nil
	} else {
		return nil, fmt.Errorf("aggrigatore not found")
	}
}
