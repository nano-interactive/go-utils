package testing

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFindFileSuccess(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Act
	file, err := FindFile(".", "README.md")

	// Assert
	assert.NotEmpty(file)
	assert.NoError(err)
}
