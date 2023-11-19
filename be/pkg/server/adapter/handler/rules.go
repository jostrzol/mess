package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/jostrzol/mess/pkg/server/core/room"
	"github.com/jostrzol/mess/pkg/server/ioc"
)

const HclContent = "application/hcl"

type RulesHandler struct {
	service *room.Service `container:"type"`
}

func FormatRules(h *RulesHandler, g *gin.Engine) {
	g.PUT("/rules/format", func(c *gin.Context) {
		src, err := io.ReadAll(c.Request.Body)
		if err != nil {
			AbortWithError(c, err)
			return
		}
		out := hclwrite.Format(src)
		c.Data(http.StatusOK, HclContent, out)
	})
}

func init() {
	ioc.MustHandlerFill[RulesHandler](FormatRules)
}
