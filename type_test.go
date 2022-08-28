package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConfigType(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	type test struct {
		input    string
		expected Type
	}

	data := []test{
		{"json", JSON},
		{"yaml", YAML},
		{"", YAML},
		{"toml", TOML},
	}

	for _, item := range data {
		configType, err := ParseConfigType(item.input)
		assert.NoError(err)
		assert.Equal(item.expected, configType)
	}
}

func TestParseConfigType_Error(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	_, err := ParseConfigType("invalid_type")
	assert.Error(err)
}
