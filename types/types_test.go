package types

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

type X struct {
	Active NullBool `json:"active"`
}

func TestNullBool_UnmarshalJSON(t *testing.T) {
	assert := require.New(t)
	data := []byte(`{"active": true}`)
	x := X{}

	err := json.Unmarshal(data, &x)
	assert.NoError(err)
}

func TestNullBool_UnmarshalJSONActiveNull(t *testing.T) {
	assert := require.New(t)
	data := []byte(`{"active": null}`)
	x := X{}

	err := json.Unmarshal(data, &x)
	assert.NoError(err)
	assert.False(x.Active.Valid)
}
