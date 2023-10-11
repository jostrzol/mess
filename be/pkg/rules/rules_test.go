package rules

import (
	"testing"

	"github.com/jostrzol/mess/pkg/rules/rulestest"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	_, err := DecodeRules("../../rules/chess.hcl", rulestest.RandomInteractor{}, true)
	assert.NoError(t, err)
}
