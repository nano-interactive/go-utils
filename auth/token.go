package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/nano-interactive/go-utils"
)

const (
	AuthTypeBearer = "bearer"
)

func ParseTokenFromHeader(header string) ([]byte, error) {
	components := strings.SplitN(header, " ", 2)

	if len(components) < 2 {
		return nil, fmt.Errorf("Invalid header")
	}

	authType := components[0]
	token := components[1]

	if !strings.EqualFold(authType, AuthTypeBearer) {
		return nil, fmt.Errorf("Invalid authorization type: %s", authType)
	}

	var tokenBytes [32]byte

	if base64.RawURLEncoding.DecodedLen(len(token)) != 32 {
		return nil, fmt.Errorf("Invalid token")
	}

	_, err := base64.RawURLEncoding.Decode(tokenBytes[:], utils.UnsafeBytes(token))
	if err != nil {
		return nil, fmt.Errorf("Invalid token")
	}

	return tokenBytes[:], nil
}

func CompareTokens(token1 []byte, token2 []byte) int {
	// Now, you may wonder why we need this function if all ti does
	// is call a library function. Well my dear reader, do you think
	// that everyone will remember to compare tokens in constant time?
	//
	// Correct, neither do I.
	return subtle.ConstantTimeCompare(token1, token2)
}
