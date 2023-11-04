package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/pkg/server/adapter/httpschema"
)

func AbortWithError(c *gin.Context, err error) {
	_ = c.Error(err)
	schemaError := httpschema.NewError(err)
	c.JSON(schemaError.Status, schemaError)
	c.Abort()
}
