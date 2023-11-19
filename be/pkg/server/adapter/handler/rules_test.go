package handler_test

import (
	"io"
	"testing"

	"github.com/jostrzol/mess/pkg/server/adapter/handler/handlertest"
	"github.com/stretchr/testify/suite"
)

type RulesSuite struct {
	handlertest.HandlerSuite[RulesClient]
}

func (s *RulesSuite) TestFormatRules() {
	// given
	src := "a{b=   2}"

	// when
	out := s.Client().formatRules(src)

	// then
	s.Equal("a { b = 2 }", out)
}

type RulesClient struct{ *handlertest.BaseClient }

func (c *RulesClient) formatRules(src string) (out string) {
	res := c.ServeOk("PUT", "/rules/format", []byte(src))
	bytes, err := io.ReadAll(res.Body)
	c.NoError(err)
	return string(bytes)
}

func TestRulesSuite(t *testing.T) {
	suite.Run(t, new(RulesSuite))
}
