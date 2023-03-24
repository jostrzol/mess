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
		IsFinished     bool
		WinnerColorCty cty.Value
	}
	if err = gocty.FromCtyValue(resultCty, &result); err != nil {
		log.Printf("parsing pick_winner user-defined function's result: %v", err)
		return false, nil
	}

	if !result.IsFinished || result.WinnerColorCty.IsNull() {
		return result.IsFinished, nil
	}

	var winnerColor string
	if err = gocty.FromCtyValue(result.WinnerColorCty, &winnerColor); err != nil {
		log.Printf("parsing pick_winner user-defined function's result: winner color: %v", err)
		return false, nil
	}

	color, err := color.ColorString(winnerColor)
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
			log.Printf("calling motion generator %q for %v at %v: %v", name, piece, piece.Square(), err)
			return make([]board.Square, 0)
		}

		squares, err := ctymess.SquaresFromCty(result)
		if err != nil {
			log.Printf("parsing motion generator result: %v", err)
		}
		return squares
	}, nil
}

func (c *callbackFunctionsConfig) GetStateValidators() ([]mess.StateValidator, error) {
	validators := make([]mess.StateValidator, 0, len(c.StateValidators))

	for validatorName, validatorCty := range c.StateValidators {
		validator := func(state *mess.State, move *mess.Move) bool {
			stateCty := ctymess.GameStateToCty(state)
			moveCty := ctymess.MoveToCty(move)
			c.refreshGameStateInContext()
			resultCty, err := validatorCty.Call([]cty.Value{stateCty, moveCty})
			if err != nil {
				log.Printf("calling state validator %q for move %v: %v", validatorName, move, err)
				return false
			}
			var result bool
			err = gocty.FromCtyValue(resultCty, &result)
			if err != nil {
				log.Printf("parsing state validator result: %v", err)
			}
			return result
		}
		validators = append(validators, validator)
	}
	return validators, nil
}

func (c *callbackFunctionsConfig) refreshGameStateInContext() {
	newGame := ctymess.GameStateToCty(c.State)
	c.EvalContext.Variables["game"] = newGame
}
