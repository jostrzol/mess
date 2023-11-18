package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	_, err := DecodeRulesFromOs("../../rules/chess.hcl", true)
	assert.NoError(t, err)
}
