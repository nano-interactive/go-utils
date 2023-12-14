package utils

import (
	"crypto/rand"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/nano-interactive/go-utils/v2/__mocks__/io"
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

func TestGetRequestId_ReturnError_On_Crypto_Read_Call(t *testing.T) {
	assert := require.New(t)

	reader := rand.Reader

	t.Cleanup(func() {
		rand.Reader = reader
	})

	mockReader := &io.MockReader{}

	rand.Reader = mockReader

	mockReader.On("Read", mock.Anything).Return(0, errors.New("failed to read enough bytes"))

	requestId, err := GetRequestId()

	assert.EqualError(err, "failed to read enough bytes")
	assert.Empty(requestId)

	mockReader.AssertNumberOfCalls(t, "Read", 1)
	mockReader.AssertExpectations(t)
}

func BenchmarkGetRequestId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetRequestId()
	}
}
