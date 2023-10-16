package testing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindFileSuccess(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Act
	file := FindFile(t, "README.md")

	// Assert
	assert.NotEmpty(file)
}
