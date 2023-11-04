package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
)

func AbortWithError(c *gin.Context, err error) {
	_ = c.Error(err)
	schemaError := schema.NewError(err)
	c.JSON(schemaError.Status, schemaError)
	c.Abort()
}
