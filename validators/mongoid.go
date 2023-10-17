package validators

import (
	"regexp"

	"github.com/invopop/validation"
)

type ObjectID string

var (
	ObjectIDRuleErr = validation.NewError("validation_is_mongoid", "Invalid ObjectID")
	ObjectIDRule    = validation.NewStringRuleWithError(IsObjectId, ObjectIDRuleErr)

	objectIdRegex = regexp.MustCompile("(?i)^[0-9a-f]{24}$")
)

func IsObjectId(val string) bool {
	return objectIdRegex.MatchString(val)
}

func (o ObjectID) Validate() error {
	return validation.Validate(string(o), ObjectIDRule)
}
