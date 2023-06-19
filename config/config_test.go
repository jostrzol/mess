package config

import (
	"testing"

	"github.com/jostrzol/mess/config/configtest"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	_, err := DecodeConfig("../rules/chess.hcl", configtest.RandomInteractor{}, true)
	assert.NoError(t, err)
}
