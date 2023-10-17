package validators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsObjectId(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	objectId := "652e77eeadd7d603e4420c3d"

	assert.NoError(ObjectID(objectId).Validate())
}

func BenchmarkIsObjectId(b *testing.B) {
	b.Run("Good", func(b *testing.B) {
		objectId := "652e77eeadd7d603e4420c3d"

		for i := 0; i < b.N; i++ {
			IsObjectId(objectId)
		}
	})

	b.Run("NotEvenLength", func(b *testing.B) {
		objectId := "652e77eeadd7d603e4420c3d3"

		for i := 0; i < b.N; i++ {
			IsObjectId(objectId)
		}
	})

	b.Run("EventButLessThan24", func(b *testing.B) {
		objectId := "652e77eeadd7d603e4420c3"

		for i := 0; i < b.N; i++ {
			IsObjectId(objectId)
		}
	})
}
