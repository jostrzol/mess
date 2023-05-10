package config

import (
	"testing"

	"github.com/jostrzol/mess/config/configtest"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	_, err := DecodeConfig("../rules.hcl", configtest.RandomInteractor{})
	assert.NoError(t, err)
}
