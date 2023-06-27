package utils_test

import (
	"net"
	"testing"

	"github.com/valyala/fasthttp"

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

func TestAnonymizeIpSuccess(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	ip := []byte("192.172.90.70")
	expectedIp := []byte("192.172.90.0")

	// Act
	anonymizedIp := utils.AnonymizeIp(ip)

	// Assert
	assert.NotEqual(ip, anonymizedIp)
	assert.Equal(expectedIp, anonymizedIp)
}

func TestAnonymizeIpInvalidIp(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	ip := []byte("invalid_ip")

	// Act
	anonymizedIp := utils.AnonymizeIp(ip)

	// Assert
	assert.Equal([]byte(utils.UnknownIp), anonymizedIp)
}

func TestRealIpReturnNilEmptyHeader(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	header := fasthttp.ResponseHeader{}

	// Act
	ip := utils.RealIp(&header)

	// Assert
	assert.Nil(ip)
}

func TestRealIpReturnIpNormalIp(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	header := fasthttp.ResponseHeader{}
	ipString := "10.20.30.40"
	header.Set(utils.HeaderXForwardedFor, ipString)

	// Act
	ip := utils.RealIp(&header)

	// Assert
	assert.NotNil(ip)
	assert.Equal([]byte(ipString), ip)
}

func TestRealIpReturnIpCommaSeperated(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	header := fasthttp.ResponseHeader{}
	expected := "10.20.30.40"
	ipString := expected + ",50.60.70.80"
	header.Set(utils.HeaderXForwardedFor, ipString)

	// Act
	ip := utils.RealIp(&header)

	// Assert
	assert.NotNil(ip)
	assert.Equal([]byte(expected), ip)
}
