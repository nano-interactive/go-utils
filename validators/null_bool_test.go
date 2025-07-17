package validators_test

import (
	"errors"
	"testing"

	"github.com/nano-interactive/go-utils/v2/validators"
	"github.com/stretchr/testify/assert"
)

func TestNullBoolValidator(t *testing.T) {
	validator := validators.NullBoolValidator{}

	tests := []struct {
		name     string
		input    *bool
		expected error
	}{
		{
			name:     "Valid true value",
			input:    boolPtr(true),
			expected: nil,
		},
		{
			name:     "Valid false value",
			input:    boolPtr(false),
			expected: nil,
		},
		{
			name:     "Nil pointer",
			input:    nil,
			expected: errors.New("boolean value is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.input)
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected.Error())
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
