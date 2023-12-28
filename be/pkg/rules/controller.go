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
	rules *rules
}

func newController(state *mess.State, ctx *hcl.EvalContext, rules *rules) *controller {
	return &controller{
		state: state,
		ctx:   ctx,
		rules: rules,
	}
}

func (c *controller) Resolution(state *mess.State) mess.Resolution {
	ctyState := c.refreshGameStateInContext()
	resultCty, err := c.rules.Functions.ResolutionFunc.Call([]cty.Value{ctyState})
	if err != nil {
		log.Printf("calling resolve user-defined function: %v", err)
		return mess.Resolution{}
	}

	var result struct {
		DidEnd         bool      `cty:"did_end"`
		WinnerColorCty cty.Value `cty:"winner"`
	}
	if err = gocty.FromCtyValue(resultCty, &result); err != nil {
		log.Printf("parsing resolve user-defined function's result: %v", err)
		return mess.Resolution{}
	}

	if !result.DidEnd || result.WinnerColorCty.IsNull() {
		// check if the current player can move - if not it's a stalemate
		if len(state.ValidMoves()) == 0 {
			result.DidEnd = true
		}
		return mess.Resolution{DidEnd: result.DidEnd}
	}

	var winnerColor string
	if err = gocty.FromCtyValue(result.WinnerColorCty, &winnerColor); err != nil {
		log.Printf("parsing resolve user-defined function's result: winner color: %v", err)
		return mess.Resolution{}
	}

	color, err := color.ColorString(winnerColor)
	if err != nil {
		log.Printf("parsing winner color: %v", err)
		return mess.Resolution{}
	}

	return mess.Resolution{DidEnd: true, Winner: state.Player(color)}
}

func (c *controller) TurnChoice(state *mess.State) (*mess.Choice, error) {
	c.refreshGameStateInContext()

	funcName := c.rules.Turn.ChoiceFunctionName
	choiceGeneratorFunc, ok := c.rules.Functions.CustomFuncs[funcName]
	if !ok {
		return nil, fmt.Errorf("user function %q not found", funcName)
	}

	choiceCty, err := choiceGeneratorFunc.Call([]cty.Value{})
	if err != nil {
		return nil, fmt.Errorf("calling turn choice function %q: %w", funcName, err)
	}

	choice, err := ctymess.ChoiceFromCty(state, choiceCty)
	if err != nil {
		return nil, fmt.Errorf("parsing choice function %q result: %w", funcName, err)
	}

	return choice, nil
}

func (c *controller) Turn(_ *mess.State, options []mess.Option) error {
	c.refreshGameStateInContext()
	optionsCty := ctymess.OptionsToCty(options)

	funcName := c.rules.Turn.ActionName
	turnFunc, ok := c.rules.Functions.CustomFuncs[funcName]
	if !ok {
		return fmt.Errorf("user function %q not found", funcName)
	}

	_, err := turnFunc.Call([]cty.Value{optionsCty})
	if err != nil {
		return fmt.Errorf("calling turn action function %q: %v", funcName, err)
	}

	return nil
}

func (c *controller) GetCustomFuncAsGenerator(name string) (mess.MoveGeneratorFunc, error) {
	funcCty, ok := c.rules.Functions.CustomFuncs[name]
	if !ok {
		return nil, fmt.Errorf("user function %q not found", name)
	}

	return func(piece *mess.Piece) []board.Square {
		pieceCty := ctymess.PieceToCty(piece)
		squareCty := ctymess.SquareToCty(piece.Square())
		c.refreshGameStateInContext()
		result, err := funcCty.Call([]cty.Value{squareCty, pieceCty})
		if err != nil {
			log.Printf("calling motion generator %q: %v", name, err)
			return make([]board.Square, 0)
		}

		squares, err := ctymess.SquaresFromCty(result)
		if err != nil {
			log.Printf("parsing motion generator %q result: %v", name, err)
		}
		return squares
	}, nil
}

func (c *controller) GetCustomFuncAsChoiceFunction(name string) (mess.MoveChoiceFunc, error) {
	funcCty, ok := c.rules.Functions.CustomFuncs[name]
	if !ok {
		return nil, fmt.Errorf("user function %q not found", name)
	}

	return func(piece *mess.Piece, from, to board.Square) *mess.Choice {
		pieceCty := ctymess.PieceToCty(piece)
		fromCty := ctymess.SquareToCty(from)
		toCty := ctymess.SquareToCty(to)

		c.refreshGameStateInContext()
		result, err := funcCty.Call([]cty.Value{pieceCty, fromCty, toCty})
		if err != nil {
			fmt.Printf("error calling choice function %q: %v\n", name, err)
			return nil
		}

		choice, err := ctymess.ChoiceFromCty(c.state, result)
		if err != nil {
			fmt.Printf("error parsing choice function %q result: %v\n", name, err)
			return nil
		}

		return choice
	}, nil
}

func (c *controller) GetCustomFuncAsAction(name string) (mess.MoveActionFunc, error) {
	funcCty, ok := c.rules.Functions.CustomFuncs[name]
	if !ok {
		return nil, fmt.Errorf("user function %q not found", name)
	}

	return func(piece *mess.Piece, from, to board.Square, optionSet []mess.Option) error {
		pieceCty := ctymess.PieceToCty(piece)
		fromCty := ctymess.SquareToCty(from)
		toCty := ctymess.SquareToCty(to)

		var err error
		c.refreshGameStateInContext()
		optionsCty := ctymess.OptionsToCty(optionSet)
		_, err = funcCty.Call([]cty.Value{pieceCty, fromCty, toCty, optionsCty})

		if err != nil {
			return fmt.Errorf("calling motion action %q: %v", name, err)
		}
		return nil
	}, nil
}

func (c *controller) GetStateValidators() ([]mess.StateValidator, error) {
	validators := make([]mess.StateValidator, 0, len(c.rules.Functions.StateValidators))

	for validatorName, validatorCty := range c.rules.Functions.StateValidators {
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
