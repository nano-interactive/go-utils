package logging

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewSuccess(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Act
	logger := New(zerolog.InfoLevel.String(), false)

	// Assert
	assert.NotNil(logger)
}

func TestNewPrettyPrintSuccess(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Act
	logger := New(zerolog.InfoLevel.String(), true)

	// Assert
	assert.NotNil(logger)
}

func TestNewWrongLevel(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Assert
	assert.PanicsWithValue("Failed to parse logging level: wrong", func() {
		// Act
		New("wrong", false)
	}, "Code panics")
}

func TestConfigureDefaultLoggerPanics(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Assert
	assert.PanicsWithValue("Failed to parse logging level: wrong", func() {
		// Act
		ConfigureDefaultLogger("wrong", false)
	}, "Code panics")
}

func TestConfigureDefaultLoggerSuccessWithPrettyPrint(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Assert
	assert.NotPanics(func() {
		// Act
		ConfigureDefaultLogger(zerolog.InfoLevel.String(), true)
	})
}

func TestConfigureDefaultLoggerSuccessWithoutPrettyPrint(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Assert
	assert.NotPanics(func() {
		// Act
		ConfigureDefaultLogger(zerolog.InfoLevel.String(), false)
	})
}
