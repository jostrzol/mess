package mess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func trueStateValidator(*State, *GeneratedMove) bool  { return true }
func falseStateValidator(*State, *GeneratedMove) bool { return false }

func TestChainStateValidator(t *testing.T) {
	tests := []struct {
		name       string
		validators []StateValidator
		expected   bool
	}{
		{
			name:       "Empty",
			validators: []StateValidator{},
			expected:   true,
		},
		{
			name:       "OneTrue",
			validators: []StateValidator{trueStateValidator},
			expected:   true,
		},
		{
			name:       "OneFalse",
			validators: []StateValidator{falseStateValidator},
			expected:   false,
		},
		{
			name:       "OneFalseOneTrue",
			validators: []StateValidator{falseStateValidator, trueStateValidator},
			expected:   false,
		},
		{
			name:       "TwoFalse",
			validators: []StateValidator{falseStateValidator, falseStateValidator},
			expected:   false,
		},
		{
			name:       "TwoTrue",
			validators: []StateValidator{trueStateValidator, trueStateValidator},
			expected:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validators := chainStateValidators(tt.validators)
			isValid := validators.Validate(nil, nil)
			assert.Equal(t, tt.expected, isValid)
		})
	}
}
