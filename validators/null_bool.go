package validators

import "errors"

// NullBoolValidator validates if a value is boolean pointer.
// Useful for checking if a variable exists in configuration or not.
type NullBoolValidator struct{}

func (r NullBoolValidator) Validate(value any) error {
	ptr := value.(*bool)
	if ptr == nil {
		return errors.New("boolean value is required")
	}

	return nil
}
