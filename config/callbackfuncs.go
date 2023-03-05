package config

import (
	"fmt"
	"log"

	"github.com/jostrzol/mess/config/ctyconv"
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type callbackFunctionsConfig struct {
	DecideWinnerFunc function.Function            `mapstructure:"decide_winner"`
	CustomFuncs      map[string]function.Function `mapstructure:",remain"`
}

func (c *callbackFunctionsConfig) DecideWinner(state *mess.State) *mess.Player {
	ctyState := ctyconv.GameStateToCty(state)
	ctyWinnerColor, err := c.DecideWinnerFunc.Call([]cty.Value{ctyState})
	if err != nil {
		log.Printf("calling user-defined function: %v", err)
		return nil
	}

	color, err := ctyconv.ColorFromCty(ctyWinnerColor)
	if err != nil {
		log.Printf("parsing winner color: %v", err)
		return nil
	}

	if color == nil {
		return nil
	}
	return state.Player(*color)
}

func (c *callbackFunctionsConfig) GetCustomFuncAsGenerator(name string) (mess.MotionGenerator, error) {
	funcCty, ok := c.CustomFuncs[name]
	if !ok {
		return nil, fmt.Errorf("user function %q not found", name)
	}

	return mess.FuncMotionGenerator(func(piece *mess.Piece) []board.Square {
		pieceCty := ctyconv.PieceToCty(piece)
		squareCty := ctyconv.SquareToCty(piece.Square())
		result, err := funcCty.Call([]cty.Value{squareCty, pieceCty})
		if err != nil {
			log.Printf("calling motion generator for %v at %v: %v", piece, piece.Square(), err)
			return make([]board.Square, 0)
		}

		squares, err := ctyconv.SquaresFromCty(result)
		if err != nil {
			log.Printf("parsing motion generator result: %v", err)
		}
		return squares
	}), nil
}
