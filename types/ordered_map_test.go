package types_test

import (
	"testing"

	"github.com/nano-interactive/go-utils/v2/types"
	"github.com/stretchr/testify/require"
)

func TestOrderedMapLen(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	emptyMap := types.NewOrderedMap[string, string](0)
	filledMap := types.NewOrderedMap[string, string](1)
	filledMap.Set("Opeth", "Isolation Years")
	cases := []struct {
		Name        string
		ExpectedLen int
		Input       types.OrderedMap[string, string]
	}{
		{
			Name:        "has no values",
			ExpectedLen: 0,
			Input:       emptyMap,
		},
		{
			Name:        "has some values",
			ExpectedLen: 1,
			Input:       filledMap,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(testCase.ExpectedLen, testCase.Input.Len())
		})
	}
}

func TestOrderedMapGetAndSet(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	om := types.NewOrderedMap[string, int](3)

	v, exists := om.Get("invalid-key")

	assert.Equal(0, v)
	assert.False(exists)

	om.Set("foo", 2)

	v, exists = om.Get("foo")

	assert.Equal(2, v)
	assert.True(exists)
}

func TestOrderedMapUnset(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	om := types.NewOrderedMap[string, int](3)
	om.Set("foo", -1)
	om.Unset("foo")

	v, exists := om.Get("foo")

	assert.Equal(0, v)
	assert.False(exists)
}

func TestOrderedMapReset(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	om := types.NewOrderedMap[string, int](3)
	om.Set("foo", 1)
	om.Set("bar", 2)
	om.Set("baz", 3)

	om.Reset()

	actualValues := make([]int, 0)

	for _, v := range om.Iter {
		actualValues = append(actualValues, v)
	}

	assert.Empty(actualValues)
}

func TestOrderedMapIter(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	om := types.NewOrderedMap[string, int](3)
	om.Set("foo", 78)
	om.Set("var", 90)
	om.Set("bar", 81)

	expectedKeys := []string{"foo", "var", "bar"}
	actualKeys := make([]string, 0, len(expectedKeys))

	expectedValues := []int{78, 90, 81}
	actualValues := make([]int, 0, len(expectedValues))

	for k, v := range om.Iter {
		actualKeys = append(actualKeys, k)
		actualValues = append(actualValues, v)
	}

	assert.Equal(expectedKeys, actualKeys)
	assert.Equal(expectedValues, actualValues)
}

func TestOrderedMapIterWithCustomOrder(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	om := types.NewOrderedMapWithCompareFunc[string, int](
		3,
		func(a, b string) int {
			if a == b {
				return 0
			}

			if a < b {
				return -1
			}

			return 1
		},
	)

	om.Set("foo", 78)
	om.Set("var", 90)
	om.Set("bar", 81)

	expectedKeys := []string{"bar", "foo", "var"}
	actualKeys := make([]string, 0, len(expectedKeys))

	expectedValues := []int{81, 78, 90}
	actualValues := make([]int, 0, len(expectedValues))

	for k, v := range om.Iter {
		actualKeys = append(actualKeys, k)
		actualValues = append(actualValues, v)
	}

	assert.Equal(expectedKeys, actualKeys)
	assert.Equal(expectedValues, actualValues)
}
