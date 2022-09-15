package utils

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	RequestIdLength = 32
	MaxEncodedLength = (RequestIdLength*8 + 5) / 6
)
// GetRequestId returns random byte slice of length 32
func GetRequestId() (string, error) {
	var bytes [RequestIdLength]byte

	_, err := rand.Read(bytes[:])

	if err != nil {
		return "", err
	}

	var encodedBytes [MaxEncodedLength]byte

	base64.RawStdEncoding.Encode(encodedBytes[:], bytes[:])

	return string(encodedBytes[:]), nil
}
