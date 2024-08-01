package validators

import (
	"encoding/hex"

	"github.com/nano-interactive/go-utils/v2"

	"github.com/nano-interactive/go-utils/v2/types"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/invopop/validation"
)

type (
	ObjectID             string
	NullableObjectID     string
	ObjectIdRule         int32
	NullableObjectIdRule int32
)

var ObjectIDRuleErr = validation.NewError("validation_is_mongoid", "Invalid ObjectID")

func NewObjectIDRule() ObjectIdRule {
	return 0
}
func NewNullableObjectIDRule() NullableObjectIdRule {
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
		if val.IsZero() {
			return ObjectIDRuleErr
		}

		return nil
	case types.ObjectID:
		if val.IsNull() {
			return ObjectIDRuleErr
		}

		return nil
	default:
		return ObjectIDRuleErr
	}
}

func (NullableObjectIdRule) Validate(v any) error {
	if v == nil {
		return nil
	}

	switch value := v.(type) {
	case string:
		if len(v.(string)) == 0 {
			return nil
		}
	case primitive.ObjectID:
		if value.IsZero() {
			return nil
		}
	case types.ObjectID:
		if value.IsNull() {
			return nil
		}
	default:
		return ObjectIDRuleErr
	}

	var rule ObjectIdRule

	return rule.Validate(v)
}

func (o ObjectID) Validate() error {
	return validation.Validate(string(o), NewObjectIDRule())
}
func (o NullableObjectID) Validate() error {
	return validation.Validate(string(o), NewNullableObjectIDRule())
}
