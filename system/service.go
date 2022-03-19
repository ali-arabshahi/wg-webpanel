package system

import (
	"fmt"
	"net"
)

type Isystem interface {
	NetworkInterfaces() ([]NetworkInterfacs, error)
}

type systemService struct{}

func New() *systemService {
	return &systemService{}
}

func (sy *systemService) NetworkInterfaces() ([]NetworkInterfacs, error) {
	netInterfaces := []NetworkInterfacs{}
	ifaces, err := net.Interfaces()
	if err != nil {
		return netInterfaces, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				// continue
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			netInterfaces = append(netInterfaces, NetworkInterfacs{
				Text:      i.Name,
				IP:        ip.String(),
				Interface: i.Name,
			})

		}
	}
	return netInterfaces, nil
}
