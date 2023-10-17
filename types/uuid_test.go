package types

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/require"
)

func TestUUID_String(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	uuid := gocql.TimeUUID()
	newUuid := CQLUUID(uuid)

	assert.Equal(uuid.String(), newUuid.String())
}

func TestUUID_MarshalJSON(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	uuid := gocql.TimeUUID()
	newUuid := CQLUUID(uuid)

	json, err := newUuid.MarshalJSON()

	assert.NoError(err)
	assert.Equal(`"`+uuid.String()+`"`, string(json))
}

func TestUUID_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	uuid := gocql.TimeUUID()
	newUuid := CQLUUID{}

	json, err := uuid.MarshalJSON()

	assert.NoError(err)
	assert.NoError(newUuid.UnmarshalJSON(json))
	assert.Equal(uuid.String(), newUuid.String())
}
