package auth_test

import (
	"testing"

	"github.com/nano-interactive/go-utils/auth"
	"github.com/stretchr/testify/require"
)

func TestParseTokenFromHeader(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	header := "Bearer U_lENfGpWMPoD6fP7resxu9ZH2b0B5QfpufFDMOSwgs"

	token, err := auth.ParseTokenFromHeader(header)

	assert.NoError(err)
	assert.Len(token, 32)
}

func TestParseTokenFromHeaderInvalidHeader(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	header := "Invalid"

	token, err := auth.ParseTokenFromHeader(header)

	assert.Nil(token)
	assert.Error(err)
	assert.EqualError(err, "Invalid header")
}

func TestParseTokenFromHeaderInvalidAuthType(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	header := "Basic U_lENfGpWMPoD6fP7resxu9ZH2b0B5QfpufFDMOSwgs"

	token, err := auth.ParseTokenFromHeader(header)

	assert.Nil(token)
	assert.Error(err)
	assert.EqualError(err, "Invalid authorization type: Basic")
}

func TestParseTokenFromHeaderInvalidTokenContents(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	header := "Bearer not_base64_encoded"

	token, err := auth.ParseTokenFromHeader(header)

	assert.Nil(token)
	assert.Error(err)
	assert.EqualError(err, "Invalid token")
}
