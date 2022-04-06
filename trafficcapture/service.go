package trafficcapture

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type ITrafficAnalyze interface {
	InterfacesInfo()
	StartCapturing()
	AddProcessor(ProcessorName string, aggrigator IPacketProcessor)
}

type IPacketProcessor interface {
	Aggrigate(IPPacket)
}

type ILogger interface {
	Infolog(msg string)
	Warnlog(msg string)
	ErrorLog(msg string)
	Debuglog(msg string)
	FatalLog(msg string)
}

//-----------------------------------------------------

//analyzer : TrafficAnalyze implementation
type analyzer struct {
	handle           *pcap.Handle
	snapshotLen      int32
	networkInterface string
	promiscuous      bool
	errorChan        chan error
	packetQueue      chan IPPacket
	processors       map[string]IPacketProcessor
	logger           ILogger
}

//-----------------------------------------------------

func New(nic string, promiscuous bool, snapshotLen int32, logger ILogger) (*analyzer, error) {
	analyz := analyzer{}
	analyz.logger = logger
	var err error
	analyz.handle, err = pcap.OpenLive(nic, snapshotLen, promiscuous, pcap.BlockForever)
	if err != nil {
		analyz.logger.ErrorLog(err.Error())
	}
	analyz.promiscuous = promiscuous
	analyz.networkInterface = nic
	analyz.snapshotLen = snapshotLen
	analyz.errorChan = make(chan error)
	analyz.packetQueue = make(chan IPPacket, 100000)
	analyz.processors = make(map[string]IPacketProcessor)
	return &analyz, nil
}

func (an *analyzer) InterfacesInfo() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	// Print device information
	fmt.Println("Devices found:")
	for _, device := range devices {
		fmt.Println("\nName: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices addresses: ", device.Description)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}
}

func (an *analyzer) StartCapturing() {
	an.checkCaptureHandlerError()
	go an.readingPacket()
	go an.sendToProcessor()
}

func (an *analyzer) AddProcessor(name string, agg IPacketProcessor) {
	an.processors[name] = agg
}

//-------------------------------------------------
//                Private functions
//-------------------------------------------------

// reload : reload network capturing
func (an *analyzer) reload() error {
	var err error
	an.handle, err = pcap.OpenLive(an.networkInterface, an.snapshotLen, an.promiscuous, pcap.BlockForever)
	if err != nil {
		return err
	}
	return nil
}

// checkCaptureHandlerError : check network capture handler error
func (an *analyzer) checkCaptureHandlerError() {
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			select {
			case handlerError := <-an.errorChan:
				an.logger.ErrorLog("error recieve : " + handlerError.Error())
				notStart := true
				for notStart {
					time.Sleep(500 * time.Millisecond)
					err := an.reload()
					if err == nil {
						an.logger.Infolog("reload successfull")
						go an.readingPacket()
						notStart = false
					}
				}
			default:
				continue
			}
		}
	}()
}

// readingPacket : read packet
func (an *analyzer) readingPacket() {
	an.logger.Infolog("start reading packet")
	for an.handle == nil {
		an.logger.Warnlog(fmt.Sprintf("interface %v not found", an.networkInterface))
		time.Sleep(2 * time.Second)
		an.reload()
	}
	packetSource := gopacket.NewPacketSource(an.handle, an.handle.LinkType())
	for {
		packet, err := packetSource.NextPacket()
		if err != nil {
			an.errorChan <- err
			break
		}
		if packet.NetworkLayer() != nil {
			protocol := ""
			srcPort := 0
			dstPort := 0
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				protocol = "TCP"
				srcPort = int(tcp.SrcPort)
				dstPort = int(tcp.DstPort)
			}
			if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
				udp, _ := udpLayer.(*layers.UDP)
				protocol = "UDP"
				srcPort = int(udp.SrcPort)
				dstPort = int(udp.DstPort)
			}
			var transHash uint64
			if packet.TransportLayer() != nil {
				transHash = packet.TransportLayer().TransportFlow().FastHash()
			}
			size := int64(packet.Metadata().CaptureInfo.CaptureLength)
			netFlow := packet.NetworkLayer().NetworkFlow()
			src, dst := netFlow.Endpoints()
			hash := netFlow.FastHash()
			ipPack := IPPacket{
				Hash:            hash,
				PacketSize:      size,
				Src:             src.String(),
				Dst:             dst.String(),
				Protocol:        protocol,
				SrcPort:         srcPort,
				DstPort:         dstPort,
				TransParentHash: transHash,
			}
			an.packetQueue <- ipPack
		}
	}
}

func (an *analyzer) sendToProcessor() {
	for packet := range an.packetQueue {
		for _, agg := range an.processors {
			agg.Aggrigate(packet)
		}
	}
}
