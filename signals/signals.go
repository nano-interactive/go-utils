package signals

import (
	"errors"
	"os"
)

// Returns an os.Signal instance for given signal
// Example: SIGHUP -> syscall.SIGHUP
func GetSignal(s string) (os.Signal, error) {
	val, ok := signals[s]

	if !ok {
		return nil, errors.New("Cannot find signal " + s)
	}

	return val, nil
}
