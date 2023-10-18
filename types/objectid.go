package types

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nano-interactive/go-utils"
)

var NilObjectID = ObjectID{
	ObjectID: primitive.NilObjectID,
}

type ObjectID struct {
	primitive.ObjectID
}

func NewObjectID() ObjectID {
	return ObjectID{
		ObjectID: primitive.NewObjectID(),
	}
}

func NewObjectIDFromTimestamp(t time.Time) ObjectID {
	return ObjectID{
		ObjectID: primitive.NewObjectIDFromTimestamp(t),
	}
}

func ObjectIDFromHex(hex string) (ObjectID, error) {
	val, err := primitive.ObjectIDFromHex(hex)

	if err != nil {
		return ObjectID{}, err
	}

	return ObjectID{
		ObjectID: val,
	}, nil
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
		*o = NilObjectID
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

func (o *ObjectID) MarshalBSON() ([]byte, error) {
	if !o.IsNull() {
		return bson.Marshal(nil)
	}

	return bson.Marshal(o.ObjectID)
}

func (o *ObjectID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.IsNull() {
		return bson.MarshalValue(nil)
	}

	return bson.MarshalValue(o.ObjectID)
}
