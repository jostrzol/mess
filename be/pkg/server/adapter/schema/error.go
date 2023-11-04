package schema

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
)

type Error struct {
	Status     int
	Message    string
	Validation []ValidationError `json:",omitempty"`
}

type ValidationError struct {
	Field   string
	Message string
}

func NewError(err error) *Error {
	var uerr usrerr.UserError
	var verrs validator.ValidationErrors
	switch {
	case errors.As(err, &uerr):
		return &Error{
			Status:  http.StatusBadRequest,
			Message: uerr.UserError(),
		}
	case errors.As(err, &verrs):
		validation := make([]ValidationError, 0, len(verrs))
		for _, verr := range verrs {
			validation = append(validation, ValidationError{
				Field:   verr.Field(),
				Message: verr.Error(),
			})
		}
		return &Error{
			Status:     http.StatusUnprocessableEntity,
			Message:    "unprocessable entity",
			Validation: validation,
		}
	default:
		return &Error{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
}
