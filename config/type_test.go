package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConfigType(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	data := []struct {
		input    string
		expected Type
	}{
		{"json", JSON},
		{"yaml", YAML},
		{"", YAML},
		{"toml", TOML},
	}

	for _, item := range data {
		t.Run("Parsing Config Type: "+item.input, func(t *testing.T) {
			configType, err := ParseType(item.input)
			assert.NoError(err)
			assert.Equal(item.expected, configType)
		})
	}
}

func TestParseConfigType_Error(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	_, err := ParseType("invalid_type")
	assert.Error(err)
}
