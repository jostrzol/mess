package configtest

import (
	"fmt"
	"math/rand"
)

type RandomInteractor struct{}

func (i RandomInteractor) Choose(options []string) int {
	return rand.Int() % len(options)
}

type PanicInteractor struct{}

func (i PanicInteractor) Choose(options []string) int {
	panic("interaction not expected")
}

type ConstInteractor struct {
	Option string
}

func (i ConstInteractor) Choose(options []string) int {
	for n, option := range options {
		if option == i.Option {
			return n
		}
	}
	panic(fmt.Sprintf("expected option %q not found in %v", i.Option, options))
}
