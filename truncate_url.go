package utils

import "bytes"

var (
	GooglePlay  = []byte("play.google")
	ITunesApple = []byte("itunes.apple")

	HTTPS         = []byte("https://")
	HTTP          = []byte("http://")
	HTTPSProtocol = []byte("https:")
	HTTPProtocol  = []byte("http:")
)

func TruncateUrl(value []byte) []byte {
	if bytes.Contains(value, GooglePlay) || bytes.Contains(value, ITunesApple) {
		return value
	}

	n := bytes.IndexAny(value, "?#")
	if n != -1 {
		value = value[0:n]
	}

	if value[len(value)-1] != byte('/') {
		value = append(value, byte('/'))
	}

	if bytes.Compare(value[0:4], []byte("http")) != 0 {
		value = append(HTTPS, value...)
	}

	if bytes.Compare(value[0:5], HTTPProtocol) == 0 {
		value = append(HTTPSProtocol, value[5:]...)
	}

	return value
}
