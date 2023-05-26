package configtest

import "math/rand"

type RandomInteractor struct{}

func (i RandomInteractor) Choose(options []string) int {
	return rand.Int() % len(options)
}

type PanicInteractor struct{}

func (i PanicInteractor) Choose(options []string) int {
	panic("interaction not expected")
}
