package environment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	data := []struct {
		env      string
		expected Env
	}{
		{"prod", Production},
		{"production", Production},
		{"dev", Development},
		{"development", Development},
		{"develop", Development},
		{"testing", Testing},
		{"test", Testing},
	}

	for _, item := range data {
		t.Run(fmt.Sprintf("Parsing String: %s", item.env), func(t *testing.T) {
			expected, err := Parse(item.env)

			assert.NoError(err)
			assert.Equal(item.expected, expected)
		})
	}
}

func TestParseEnvironment_Error(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	data := []struct {
		env      string
		expected string
	}{
		{"", "Invalid Environment: prod, production, dev, development, develop, testing, test, Given: "},
		{"something-else", "Invalid Environment: prod, production, dev, development, develop, testing, test, Given: something-else"},
	}

	for _, item := range data {
		t.Run(fmt.Sprintf("Parsing String: %s", item.env), func(t *testing.T) {
			_, err := Parse(item.env)

			assert.Error(err)
			assert.EqualError(err, item.expected)
		})
	}
}
