package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/nano-interactive/go-utils"
)

func ParseTokenFromHeader(header string, authType AuthType) ([]byte, error) {
	components := strings.SplitN(header, " ", 2)

	if len(components) < 2 {
		return nil, errors.New("Invalid header")
	}

	headerAuthType := components[0]
	headerToken := components[1]

	if !strings.EqualFold(headerAuthType, string(authType)) {
		return nil, errors.New("Invalid authorization type: " + headerAuthType)
	}

	var tokenBytes [32]byte

	if base64.RawURLEncoding.DecodedLen(len(headerToken)) != 32 {
		return nil, errors.New("Invalid token")
	}

	_, err := base64.RawURLEncoding.Decode(tokenBytes[:], utils.UnsafeBytes(headerToken))
	if err != nil {
		return nil, errors.New("Invalid token")
	}

	return tokenBytes[:], nil
}

func CompareTokens(token1 []byte, token2 []byte) int {
	// Now, you may wonder why we need this function if all it does
	// is call a library function. Well my dear reader, do you think
	// that everyone will remember to compare tokens in constant time?
	//
	// Correct, neither do I.
	return subtle.ConstantTimeCompare(token1, token2)
}
