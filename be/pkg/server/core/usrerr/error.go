package usrerr

import (
	"fmt"
)

type UserError interface {
	UserError() string
}

func Wrap(cause error, format string, a ...any) error {
	return &wrapperUserError{cause: cause, userMessage: fmt.Sprintf(format, a...)}
}

func Errorf(format string, a ...any) error {
	return &rootUserError{fmt.Errorf(format, a...)}
}

type rootUserError struct {
	error
}

func (e *rootUserError) UserError() string {
	return e.Error()
}

type wrapperUserError struct {
	cause       error
	userMessage string
}

func (e *wrapperUserError) Error() string {
	return fmt.Sprintf("%v: %v", e.userMessage, e.cause)
}

func (e *wrapperUserError) UserError() string {
	return e.userMessage
}
