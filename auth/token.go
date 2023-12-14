package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/nano-interactive/go-utils"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrInvalidAuthType = errors.New("invalid authorization type")
)

func ExtractTokenFromHeader(header string, authType AuthType) (string, error) {
	components := strings.SplitN(header, " ", 2)

	if len(components) < 2 {
		return "", ErrInvalidToken
	}

	headerAuthType := components[0]
	headerToken := components[1]

	if !strings.EqualFold(headerAuthType, string(authType)) {
		return "", ErrInvalidAuthType
	}

	return headerToken, nil
}

func ParseTokenFromHeader(header string, authType AuthType) ([]byte, error) {
	headerToken, err := ExtractTokenFromHeader(header, authType)
	if err != nil {
		return nil, err
	}

	var tokenBytes [32]byte

	if base64.RawURLEncoding.DecodedLen(len(headerToken)) != 32 {
		return nil, ErrInvalidToken
	}

	_, err = base64.RawURLEncoding.Decode(tokenBytes[:], utils.UnsafeBytes(headerToken))
	if err != nil {
		return nil, ErrInvalidToken
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

func CompareTokenStrings(token1 string, token2 string) int {
	// Now, you may wonder why we need this function if all it does
	// is call a library function. Well my dear reader, do you think
	// that everyone will remember to compare tokens in constant time?
	//
	// Correct, neither do I.
	return subtle.ConstantTimeCompare(utils.UnsafeBytes(token1), utils.UnsafeBytes(token2))
}
