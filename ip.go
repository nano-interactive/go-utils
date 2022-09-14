package utils

import (
	"bytes"
	"crypto/rand"
	"errors"
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

// AnonymizeIp slices ip address byte array and writes 0 to last octet
// example 192.172.90.70 -> 192.172.90.0
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

// GetRequestId returns hashed hostname trimmed to first 6 characters
func GetRequestId() ([32]byte, error) {
	var bytes [32]byte
	n, err := rand.Read(bytes[:])
	if err != nil {
		return [32]byte{}, err
	}
	if n != 32 {
		return [32]byte{}, errors.New("not enough bytes")
	}

	return bytes, nil
}
