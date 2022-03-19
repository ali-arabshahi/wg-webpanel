package trafficaggrigate

import "net"

//IPRangeLookup is a method for recognizing Private and Public IPs
func ipRangePrivate(IP net.IP) bool {
	/*
		Class        Starting IPAddress    Ending IP Address    # of Hosts
		A            10.0.0.0              10.255.255.255       16,777,216
		B            172.16.0.0            172.31.255.255       1,048,576
		C            192.168.0.0           192.168.255.255      65,536
		Link-local-u 169.254.0.0           169.254.255.255      65,536
		Link-local-m 224.0.0.0             224.0.0.255          256
		Local        127.0.0.0             127.255.255.255      16777216
	*/

	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return true
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		default:
			return false
		}
	}
	return true
}
