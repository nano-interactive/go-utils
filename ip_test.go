package utils_test

import (
	"net"
	"testing"

	"github.com/nano-interactive/go-utils"
	"github.com/stretchr/testify/require"
)

func TestGetLocalIP(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	addrs, err := net.InterfaceAddrs()
	assert.NoError(err)
	var ip string

	for _, address := range addrs {
		ipnet, ok := address.(*net.IPNet)

		// check the address type and if it is not a loopback the display it
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ip = ipnet.IP.String()
			break
		}
	}

	localIp := utils.GetLocalIP()

	assert.Equal(ip, localIp)

	// Testing the cached version
	localIp = utils.GetLocalIP()

	assert.Equal(ip, localIp)

}

func TestGetLocalIPs(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	addrs, err := net.InterfaceAddrs()
	assert.NoError(err)
	var ips []string

	for _, address := range addrs {
		ipnet, ok := address.(*net.IPNet)

		// check the address type and if it is not a loopback the display it
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	localIps := utils.GetLocalIPs()

	assert.EqualValues(ips, localIps)

	// Testing the cached version
	localIps = utils.GetLocalIPs()

	assert.EqualValues(ips, localIps)
}
