package randint

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestUint8(t *testing.T) {
	t.Parallel()

	// arrange
	assert := require.New(t)

	// act
	v, err := Uint8()

	// assert
	assert.NoError(err)
	assert.Equal(uint8Len, int(unsafe.Sizeof(v)))
}

func TestUint16(t *testing.T) {
	t.Parallel()

	// arrange
	assert := require.New(t)

	// act
	v, err := Uint16()

	// assert
	assert.NoError(err)
	assert.Equal(uint16Len, int(unsafe.Sizeof(v)))
}

func TestUint32(t *testing.T) {
	t.Parallel()

	// arrange
	assert := require.New(t)

	// act
	v, err := Uint32()

	// assert
	assert.NoError(err)
	assert.Equal(uint32Len, int(unsafe.Sizeof(v)))
}

func TestUint64(t *testing.T) {
	t.Parallel()

	// arrange
	assert := require.New(t)

	// act
	v, err := Uint64()

	// assert
	assert.NoError(err)
	assert.Equal(uint64Len, int(unsafe.Sizeof(v)))
}
