package usrerr

import (
	"fmt"
)

type UserError interface {
	UserError() string
}

func Errorf(format string, a ...any) error {
	return &simpleUserError{fmt.Errorf(format, a...)}
}

type simpleUserError struct {
	error
}

func (e *simpleUserError) UserError() string {
	return e.Error()
}
