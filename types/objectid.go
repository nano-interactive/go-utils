package types

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nano-interactive/go-utils"
)

type ObjectID struct {
	primitive.ObjectID
}

func (o *ObjectID) IsNull() bool {
	return o.ObjectID.IsZero()
}

func (o *ObjectID) MarshalJSON() ([]byte, error) {
	if o.IsNull() {
		return utils.UnsafeBytes("null"), nil
	}

	return o.ObjectID.MarshalJSON()
}

func (o *ObjectID) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.ObjectID = primitive.NilObjectID
	case bson.TypeObjectID:
		return bson.UnmarshalValue(t, bytes, &o.ObjectID)
	case bson.TypeString:
		var str string
		if err := bson.UnmarshalValue(t, bytes, &str); err != nil {
			return err
		}
		data, err := primitive.ObjectIDFromHex(str)
		if err != nil {
			return err
		}
		o.ObjectID = data
	default:
		return errors.New("not a object ID")
	}

	return nil
}
func (o *ObjectID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.IsNull() {
		return bson.MarshalValue(nil)
	}

	return bson.MarshalValue(o.ObjectID)
}
