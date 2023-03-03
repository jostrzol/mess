package ctyconv

import (
	"strings"
)

type manyErrors []error

func (errors manyErrors) Error() string {
	var b strings.Builder
	b.WriteString("[\n")
	for _, err := range errors {
		b.WriteByte('\t')
		b.WriteString(err.Error())
		b.WriteByte('\n')
	}
	return b.String()
}
