package utils

import (
	"errors"
	"os"
)

func GetSignal(s string) (os.Signal, error) {
	val, ok := signals[s]

	if !ok {
		return nil, errors.New("Cannot find signal " + s)
	}

	return val, nil
}
