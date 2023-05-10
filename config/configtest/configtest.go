package configtest

import "math/rand"

type RandomInteractor struct{}

func (r RandomInteractor) Choose(options []string) int {
	return rand.Int() % len(options)
}
