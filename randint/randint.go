package randint

import (
	"crypto/rand"
	"encoding/binary"
)

const (
	uint8Len  = 1
	uint16Len = 2
	uint32Len = 4
	uint64Len = 8
)

func generateBytes(length int) ([]byte, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

func Uint8() (uint8, error) {
	b, err := generateBytes(uint8Len)
	if err != nil {
		return 0, err
	}

	return uint8(b[0]), nil
}

func Uint16() (uint16, error) {
	b, err := generateBytes(uint16Len)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint16(b), nil
}

func Uint32() (uint32, error) {
	b, err := generateBytes(uint32Len)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(b), nil
}

func Uint64() (uint64, error) {
	b, err := generateBytes(uint64Len)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(b), nil
}
