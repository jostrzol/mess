package rulestest

import (
	"fmt"
	"math/rand"

	"github.com/jostrzol/mess/pkg/mess"
)

type RandomInteractor struct{}

func (i RandomInteractor) PreTurn(_ *mess.State) {
}

func (i RandomInteractor) ChooseOption(options []string) int {
	return rand.Int() % len(options)
}

func (i RandomInteractor) ChooseMove(_ *mess.State, validMoves []mess.GeneratedMove) (*mess.GeneratedMove, error) {
	idx := rand.Int() % len(validMoves)
	return &validMoves[idx], nil
}

type PanicInteractor struct{}

func (i PanicInteractor) PreTurn(_ *mess.State) {
}

func (i PanicInteractor) ChooseOption(_ []string) int {
	panic("option choosing not expected")
}

func (i PanicInteractor) ChooseMove(_ *mess.State, _ []mess.GeneratedMove) (*mess.GeneratedMove, error) {
	panic("move choosing not expected")
}

type ConstOptionInteractor struct {
	Option string
}

func (i ConstOptionInteractor) PreTurn(_ *mess.State) {
}

func (i ConstOptionInteractor) ChooseOption(options []string) int {
	for n, option := range options {
		if option == i.Option {
			return n
		}
	}
	panic(fmt.Sprintf("expected option %q not found in %v", i.Option, options))
}

func (i ConstOptionInteractor) ChooseMove(_ *mess.State, _ []mess.GeneratedMove) (*mess.GeneratedMove, error) {
	panic("move choosing not expected")
}
