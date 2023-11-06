package validators

import (
	"testing"

	"github.com/invopop/validation"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type request struct {
	ID primitive.ObjectID
}

func (r request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, NewObjectIDRule()))
}

func TestReq(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	id, err := primitive.ObjectIDFromHex("invalid id")
	assert.Error(err)
	req := request{
		ID: id,
	}

	err = req.Validate()

	assert.Error(err)
	assert.EqualError(err, "ID: Invalid ObjectID.")
}

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
			_ = NewObjectIDRule().Validate(objectId)
		}
	})

	b.Run("NotEvenLength", func(b *testing.B) {
		objectId := "652e77eeadd7d603e4420c3d3"

		for i := 0; i < b.N; i++ {
			_ = NewObjectIDRule().Validate(objectId)
		}
	})

	b.Run("EventButLessThan24", func(b *testing.B) {
		objectId := "652e77eeadd7d603e4420c3"

		for i := 0; i < b.N; i++ {
			NewObjectIDRule().Validate(objectId)
		}
	})
}
