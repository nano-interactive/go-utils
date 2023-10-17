package validators

import (
	"encoding/hex"

	"github.com/invopop/validation"
	"github.com/nano-interactive/go-utils"
)

type ObjectID string

var (
	ObjectIDRuleErr = validation.NewError("validation_is_mongoid", "Invalid ObjectID")
	ObjectIDRule    = validation.NewStringRuleWithError(IsObjectId, ObjectIDRuleErr)
)

func IsObjectId(val string) bool {
	l := len(val)

	if l != 24 {
		return false
	}

	var data [12]byte

	_, err := hex.Decode(data[:], utils.UnsafeBytes(val))
	return err == nil
}

func (o ObjectID) Validate() error {
	return validation.Validate(string(o), ObjectIDRule)
}
