package rules

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/rules/ctymess"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type controller struct {
	state *mess.State
	ctx   *hcl.EvalContext
	rules *callbackFunctionsRules
}

func newController(state *mess.State, ctx *hcl.EvalContext, rules *rules) *controller {
	return &controller{
		state: state,
		ctx:   ctx,
		rules: &rules.Functions,
	}
}

func (c *controller) PickWinner(state *mess.State) (bool, *mess.Player) {
	ctyState := c.refreshGameStateInContext()
	resultCty, err := c.rules.PickWinnerFunc.Call([]cty.Value{ctyState})
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
		// check if the current player can move - if not it's a stalemate
		if len(state.ValidMoves()) == 0 {
			result.IsFinished = true
		}
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

func (c *controller) Turn(state *mess.State) error {
	c.refreshGameStateInContext()
	_, err := c.rules.TurnFunc.Call([]cty.Value{})
	return err
}

func (c *controller) GetCustomFuncAsGenerator(name string) (func(*mess.Piece) []board.Square, error) {
	funcCty, ok := c.rules.CustomFuncs[name]
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
			log.Printf("parsing motion generator %q result: %v", name, err)
		}
		return squares
	}, nil
}

func (c *controller) GetCustomFuncAsAction(name string) (mess.MoveAction, error) {
	funcCty, ok := c.rules.CustomFuncs[name]
	if !ok {
		return nil, fmt.Errorf("user function %q not found", name)
	}

	return func(piece *mess.Piece, from board.Square, to board.Square) {
		pieceCty := ctymess.PieceToCty(piece)
		fromCty := ctymess.SquareToCty(from)
		toCty := ctymess.SquareToCty(to)
		c.refreshGameStateInContext()
		_, err := funcCty.Call([]cty.Value{pieceCty, fromCty, toCty})
		if err != nil {
			log.Printf("calling motion action %q for %v %v->%v: %v", name, piece, from, to, err)
			return
		}
	}, nil
}

func (c *controller) GetStateValidators() ([]mess.StateValidator, error) {
	validators := make([]mess.StateValidator, 0, len(c.rules.StateValidators))

	for validatorName, validatorCty := range c.rules.StateValidators {
		// copy is required, because else the validator closure
		// would always take validator and name from the last c.StateValidators
		// entry (the iterator reference changes as the loop iterates)
		valNameCopy := validatorName
		valCopy := validatorCty
		validator := func(state *mess.State, move *mess.Move) bool {
			moveCty := ctymess.MoveToCty(move)
			c.refreshGameStateInContext()
			resultCty, err := valCopy.Call([]cty.Value{moveCty})
			if err != nil {
				log.Printf("calling state validator %q for move %v: %v", valNameCopy, move, err)
				return false
			}
			var result bool
			err = gocty.FromCtyValue(resultCty, &result)
			if err != nil {
				log.Printf("parsing state validator %q result: %v", valNameCopy, err)
			}
			return result
		}
		validators = append(validators, validator)
	}
	return validators, nil
}

func (c *controller) refreshGameStateInContext() cty.Value {
	newState := ctymess.StateToCty(c.state)
	c.ctx.Variables["game"] = newState
	return newState
}
