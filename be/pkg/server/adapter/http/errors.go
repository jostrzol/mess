package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
)

func AbortWithError(c *gin.Context, err error) {
	_ = c.Error(err)
	errCurr := err
	for errCurr != nil {
		errUser, ok := errCurr.(usrerr.UserError)
		if ok {
			c.Set("user-error", errUser)
			status := http.StatusBadRequest
			c.JSON(status, &httpError{Status: status, Message: errUser.UserError()})
			c.Abort()
			return
		}
		errCurr = errors.Unwrap(errCurr)
	}
	status := http.StatusInternalServerError
	c.JSON(status, &httpError{Status: status, Message: "internal server error"})
	c.Abort()
}

type httpError struct {
	Status  int
	Message string
}
