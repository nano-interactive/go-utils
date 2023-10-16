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
		{env: "staging", expected: Staging},
		{env: "stage", expected: Staging},
	}

	for _, item := range data {
		t.Run(fmt.Sprintf("Parsing String: %s", item.env), func(_ *testing.T) {
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
		{"", "invalid Environment: prod, production, dev, development, develop, testing, test"},
		{"something-else", "invalid Environment: prod, production, dev, development, develop, testing, test"},
	}

	for _, item := range data {
		t.Run(fmt.Sprintf("Parsing String: %s", item.env), func(_ *testing.T) {
			_, err := Parse(item.env)

			assert.Error(err)
			assert.EqualError(err, item.expected)
		})
	}
}

func TestMustParse_Panic(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	assert.Panics(func() {
		MustParse("something-else")
	})
}
