package config

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/jostrzol/mess/config/ctymess"
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

type callbackFunctionsConfig struct {
	PickWinnerFunc  function.Function            `mapstructure:"pick_winner"`
	CustomFuncs     map[string]function.Function `mapstructure:",remain"`
	StateValidators map[string]function.Function
	EvalContext     *hcl.EvalContext
	State           *mess.State
}

func (c *callbackFunctionsConfig) PickWinner(state *mess.State) (bool, *mess.Player) {
	ctyState := ctymess.GameStateToCty(state)
	c.refreshGameStateInContext()
	resultCty, err := c.PickWinnerFunc.Call([]cty.Value{ctyState})
	if err != nil {
		log.Printf("calling pick_winner user-defined function: %v", err)
		return false, nil
	}

	var result struct {
		IsFinished  bool
		WinnerColor *string
	}
	if err = gocty.FromCtyValue(resultCty, &result); err != nil {
		log.Printf("parsing pick_winner user-defined function's result: %v", err)
		return false, nil
	}

	if !result.IsFinished || result.WinnerColor == nil {
		return result.IsFinished, nil
	}

	color, err := color.ColorString(*result.WinnerColor)
	if err != nil {
		log.Printf("parsing winner color: %v", err)
		return false, nil
	}

	return true, state.Player(color)
}

func (c *callbackFunctionsConfig) GetCustomFuncAsGenerator(name string) (mess.MoveGenerator, error) {
	funcCty, ok := c.CustomFuncs[name]
	if !ok {
		return nil, fmt.Errorf("user function %q not found", name)
	}

	return func(piece *mess.Piece) []board.Square {
		pieceCty := ctymess.PieceToCty(piece)
		squareCty := ctymess.SquareToCty(piece.Square())
		c.refreshGameStateInContext()
		result, err := funcCty.Call([]cty.Value{squareCty, pieceCty})
		if err != nil {
			log.Printf("calling motion generator for %v at %v: %v", piece, piece.Square(), err)
			return make([]board.Square, 0)
		}

		squares, err := ctymess.SquaresFromCty(result)
		if err != nil {
			log.Printf("parsing motion generator result: %v", err)
		}
		return squares
	}, nil
}

func (c *callbackFunctionsConfig) refreshGameStateInContext() {
	newGame := ctymess.GameStateToCty(c.State)
	c.EvalContext.Variables["game"] = newGame
}
