package utils

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	// RequestIdLength
	RequestIdLength = 32

	// MaxEncodedLength
	MaxEncodedLength = (RequestIdLength*8 + 5) / 6
)

// Returns 32 character string of random bytes
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
