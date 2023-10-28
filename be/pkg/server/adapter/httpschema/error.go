package httpschema

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
)

type Error struct {
	Status     int               `json:"status"`
	Message    string            `json:"message"`
	Validation []ValidationError `json:"validation,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewError(err error) *Error {
	errCurr := err
	for errCurr != nil {
		switch v := errCurr.(type) {
		case usrerr.UserError:
			return &Error{
				Status:  http.StatusBadRequest,
				Message: v.UserError(),
			}
		case validator.ValidationErrors:
			validation := make([]ValidationError, 0, len(v))
			for _, e := range v {
				validation = append(validation, ValidationError{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
			return &Error{
				Status:     http.StatusUnprocessableEntity,
				Message:    "unprocessable entity",
				Validation: validation,
			}
		}
		errCurr = errors.Unwrap(errCurr)
	}
	return &Error{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}
}
