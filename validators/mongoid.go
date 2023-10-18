package validators

import (
	"encoding/hex"
	"github.com/nano-interactive/go-utils/types"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/invopop/validation"
	"github.com/nano-interactive/go-utils"
)

type ObjectID string

var ObjectIDRuleErr = validation.NewError("validation_is_mongoid", "Invalid ObjectID")

type ObjectIdRule int32

func NewObjectIDRule() ObjectIdRule {
	return 0
}

func (ObjectIdRule) Validate(v any) error {
	switch val := v.(type) {
	case string:
		l := len(val)

		if l != 24 {
			return ObjectIDRuleErr
		}

		var data [12]byte

		if _, err := hex.Decode(data[:], utils.UnsafeBytes(val)); err != nil {
			return ObjectIDRuleErr
		}
		return nil
	case primitive.ObjectID:
		return nil
	case types.ObjectID:
		return nil
	default:
		return ObjectIDRuleErr
	}
}

func (o ObjectID) Validate() error {
	return validation.Validate(string(o), NewObjectIDRule())
}
