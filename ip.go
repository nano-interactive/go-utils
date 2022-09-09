package utils

import (
	"net"
	"sync"
)

var (
	ip      = ""
	ips     = make([]string, 0, 5)
	onceIp  = &sync.Once{}
	onceIps = &sync.Once{}
)

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
