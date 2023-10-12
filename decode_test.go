package utils_test

import (
	"strconv"
	"testing"

	"github.com/nano-interactive/go-utils"
	"github.com/stretchr/testify/require"
)

type item struct {
	src []byte
	exp []byte
}

func TestDecodeQuery(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	data := []item{
		{
			src: []byte("https://example.com/?hello=world"),
			exp: []byte("https://example.com/?hello=world"),
		},
		{
			src: []byte("https://example.com/?hello=world&this=works"),
			exp: []byte("https://example.com/?hello=world&this=works"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com%3Fhello%3Dworld%26this%3Dworks"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com?hello=world&this=works"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com%3Fhello%3Dworld%26this%3Dworks%26redirect%3Dhttps%253A%252F%252Fwww.third-example.com%253Fhello%253Dworld%2526this%253Dworks%2526redirect%253D"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com?hello=world&this=works&redirect=https://www.third-example.com?hello=world&this=works&redirect="),
		},
	}

	for i, item := range data {
		t.Run("DecodingURL_"+strconv.FormatInt(int64(i), 10), func(_ *testing.T) {
			dst := make([]byte, len(item.src))
			n, err := utils.DecodeQuery(dst, item.src)
			assert.NoError(err)
			assert.NotEqual(0, n)
			assert.EqualValues(string(item.exp), string(dst[:n]))
		})
	}
}

func TestDecodeQuery_WithOverlappingMemory(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	data := []item{
		{
			src: []byte("https://example.com/?hello=world"),
			exp: []byte("https://example.com/?hello=world"),
		},
		{
			src: []byte("https://example.com/?hello=world&this=works"),
			exp: []byte("https://example.com/?hello=world&this=works"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com%3Fhello%3Dworld%26this%3Dworks"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com?hello=world&this=works"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com%3Fhello%3Dworld%26this%3Dworks%26redirect%3Dhttps%253A%252F%252Fwww.third-example.com%253Fhello%253Dworld%2526this%253Dworks%2526redirect%253D"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com?hello=world&this=works&redirect=https://www.third-example.com?hello=world&this=works&redirect="),
		},
	}
	for i, item := range data {
		t.Run("DecodingURL_"+strconv.FormatInt(int64(i), 10), func(_ *testing.T) {
			n, err := utils.DecodeQuery(item.src, item.src)
			assert.NoError(err)
			assert.NotEqual(0, n)
			assert.EqualValues(item.exp, item.src[:n])
		})
	}
}

func TestDecodeQueryUnsafeString(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	// Needs seperate test because of the overlapping memory

	data := []item{
		{
			src: []byte("https://example.com/?hello=world"),
			exp: []byte("https://example.com/?hello=world"),
		},
		{
			src: []byte("https://example.com/?hello=world&this=works"),
			exp: []byte("https://example.com/?hello=world&this=works"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com%3Fhello%3Dworld%26this%3Dworks"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com?hello=world&this=works"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com%3Fhello%3Dworld%26this%3Dworks%26redirect%3Dhttps%253A%252F%252Fwww.third-example.com%253Fhello%253Dworld%2526this%253Dworks%2526redirect%253D"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com?hello=world&this=works&redirect=https://www.third-example.com?hello=world&this=works&redirect="),
		},
	}

	for i, item := range data {
		t.Run("DecodingURL_"+strconv.FormatInt(int64(i), 10), func(_ *testing.T) {
			exp, err := utils.DecodeQueryUnsafeString(item.src)
			assert.NoError(err)
			assert.EqualValues(string(item.exp), exp)
		})
	}
}

func TestDecodeQueryUnsafe(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Needs seperate test because of the overlapping memory
	data := []item{
		{
			src: []byte("https://example.com/?hello=world"),
			exp: []byte("https://example.com/?hello=world"),
		},
		{
			src: []byte("https://example.com/?hello=world&this=works"),
			exp: []byte("https://example.com/?hello=world&this=works"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com%3Fhello%3Dworld%26this%3Dworks"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com?hello=world&this=works"),
		},
		{
			src: []byte("https://example.com/?hello=world?redirect=https%3A%2F%2Fwww.new-example.com%3Fhello%3Dworld%26this%3Dworks%26redirect%3Dhttps%253A%252F%252Fwww.third-example.com%253Fhello%253Dworld%2526this%253Dworks%2526redirect%253D"),
			exp: []byte("https://example.com/?hello=world?redirect=https://www.new-example.com?hello=world&this=works&redirect=https://www.third-example.com?hello=world&this=works&redirect="),
		},
	}

	for i, item := range data {
		t.Run("DecodingURL_"+strconv.FormatInt(int64(i), 10), func(_ *testing.T) {
			exp, err := utils.DecodeQueryUnsafe(item.src)
			assert.NoError(err)
			assert.EqualValues(item.exp, exp)
		})
	}
}
