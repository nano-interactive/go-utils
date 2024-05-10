package randint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUint8(t *testing.T) {
	t.Parallel()

	// arrange
	assert := require.New(t)

	// act
	_, err := Uint8()

	// assert
	assert.NoError(err)
}

func TestUint16(t *testing.T) {
	t.Parallel()

	// arrange
	assert := require.New(t)

	// act
	_, err := Uint16()

	// assert
	assert.NoError(err)
}

func TestUint32(t *testing.T) {
	t.Parallel()

	// arrange
	assert := require.New(t)

	// act
	_, err := Uint32()

	// assert
	assert.NoError(err)
}

func TestUint64(t *testing.T) {
	t.Parallel()

	// arrange
	assert := require.New(t)

	// act
	_, err := Uint64()

	// assert
	assert.NoError(err)
}
