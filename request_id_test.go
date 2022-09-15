package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRequestId(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Act
	requestId, err := GetRequestId()

	// Assert
	assert.Len(requestId, MaxEncodedLength)
	assert.Nil(err)
}

func BenchmarkGetRequestId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetRequestId()
	}
}
