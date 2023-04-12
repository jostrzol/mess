package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	_, _, err := DecodeConfig("../rules.hcl")
	assert.NoError(t, err)
}