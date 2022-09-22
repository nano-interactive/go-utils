package utils

import (
	"bytes"
	"net"
	"sync"
)

var (
	ip      = ""
	ips     = make([]string, 0, 5)
	onceIp  = &sync.Once{}
	onceIps = &sync.Once{}
)

const (
	UnknownIp = "UNKNOWN IP"
)

// Returns IP address of local machine, empty string if fails
func GetLocalIP() string {
	onceIp.Do(func() {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			ip = ""
			return
		}

		for _, address := range addrs {
			ipnet, ok := address.(*net.IPNet)

			// check the address type and if it is not a loopback the display it
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	})

	return ip
}

// Returns strinngs slice of IP found on local machine
func GetLocalIPs() []string {
	onceIps.Do(func() {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			ips = nil
			return
		}

		for _, address := range addrs {
			ipnet, ok := address.(*net.IPNet)

			// check the address type and if it is not a loopback the display it
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	})

	return ips
}

// Slices IP address byte slice and writes 0 to last octet
// Example: AnonymizeIp([]byte{'1','9','2','.','1','7','2','.','9','0','.','7','0'})
// Returns []byte{'1','9','2','.','1','7','2','.','9','0','.','0'}
func AnonymizeIp(ip []byte) []byte {
	position := bytes.LastIndexByte(ip, '.')

	if position == -1 {
		return []byte(UnknownIp)
	}

	firstPart := ip[:position]
	anonymous := make([]byte, len(firstPart), len(firstPart)+2)

	copy(anonymous, firstPart)

	anonymous = append(anonymous, '.', '0')

	return anonymous
}
