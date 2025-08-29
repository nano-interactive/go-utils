package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalJSONPreservesOrder(t *testing.T) {
	t.Parallel()

	om := NewOrderedMap[string, any](10)
	om.Set("id", 123)
	om.Set("name", "Anamarija")
	om.Set("city", "Belgrade")

	data, err := json.Marshal(om)
	require.NoError(t, err)

	require.JSONEq(t,
		`{"id":123,"name":"Anamarija","city":"Belgrade"}`,
		string(data),
	)
}

func TestUnmarshalJSONPopulatesValuesAndOrder(t *testing.T) {
	t.Parallel()

	input := []byte(`{"z":1,"a":2,"x":5}`)
	om := NewOrderedMap[string, int](10)
	err := json.Unmarshal(input, &om)
	require.NoError(t, err)

	require.Len(t, om.values, 3)
	require.Contains(t, om.values, "z")
	require.Contains(t, om.values, "a")
	require.Contains(t, om.values, "x")

	require.Equal(t, []string{"z", "a", "x"}, om.insertionOrder)
}

func TestMarshalJSONWithEmptyMap(t *testing.T) {
	t.Parallel()

	om := NewOrderedMap[string, int](0)
	data, err := json.Marshal(om)
	require.NoError(t, err)
	require.Equal(t, "{}", string(data))
}

func TestMarshalJSONWithCustomSort(t *testing.T) {
	t.Parallel()
	om := NewOrderedMapWithCompareFunc[string, int](10, func(i, j string) int {
		if i < j {
			return -1
		}
		if i > j {
			return 1
		}
		return 0
	})

	om.Set("b", 2)
	om.Set("a", 1)

	data, err := json.Marshal(om)
	require.NoError(t, err)

	require.JSONEq(t, `{"a":1,"b":2}`, string(data))
}
