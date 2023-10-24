package rules

import (
	"fmt"
	"unicode/utf8"

	"github.com/hashicorp/hcl/v2"
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
)

func (c *rules) toEmptyGameState(ctx *hcl.EvalContext) (*mess.Game, error) {
	brd, err := mess.NewPieceBoard(int(c.Board.Width), int(c.Board.Height))
	if err != nil {
		return nil, fmt.Errorf("creating new board: %w", err)
	}

	state := mess.NewState(brd)
	controller := newController(state, ctx, c)

	game := mess.NewGame(state, controller)

	stateValidators, err := controller.GetStateValidators()
	if err != nil {
		return nil, fmt.Errorf("parsing state validators: %w", err)
	}
	for _, validator := range stateValidators {
		state.AddStateValidator(validator)
	}

	for _, pieceTypeRules := range c.PieceTypes.PieceTypes {
		pieceType, err := decodePieceType(controller, pieceTypeRules)
		if err != nil {
			return nil, fmt.Errorf("decoding piece type %q: %v", pieceTypeRules.Name, err)
		}
		state.AddPieceType(pieceType)
	}

	initializeContext(ctx, game)
	return game, nil
}

func decodePieceType(controller *controller, pieceTypeRules pieceTypeRules) (*mess.PieceType, error) {
	pieceType := mess.NewPieceType(pieceTypeRules.Name)
	for _, motionRules := range pieceTypeRules.Motions {
		moveGenerator, err := controller.GetCustomFuncAsGenerator(motionRules.GeneratorName)
		if err != nil {
			return nil, err
		}
		var action mess.MoveActionFunc
		if motionRules.ActionName != "" {
			action, err = controller.GetCustomFuncAsAction(motionRules.ActionName)
			if err != nil {
				return nil, err
			}
		}
		choiceGenerators := make([]mess.MoveChoiceGeneratorFunc, 0, len(motionRules.ChoiceGeneratorNames))
		for _, choiceGeneratorName := range motionRules.ChoiceGeneratorNames {
			choiceGenerator, err := controller.GetCustomFuncAsChoiceGenerator(choiceGeneratorName)
			if err != nil {
				return nil, err
			}
			choiceGenerators = append(choiceGenerators, choiceGenerator)
		}
		if motionRules.ActionName != "" {
		}
		pieceType.AddMotion(
			mess.Motion{
				Name:             motionRules.GeneratorName,
				MoveGenerator:    moveGenerator,
				ChoiceGenerators: choiceGenerators,
				Action:           action,
			},
		)
	}
	if pieceTypeRules.Symbols != nil {
		symbolWhite, err := decodeSymbol(pieceTypeRules.Symbols.White)
		if err != nil {
			return nil, err
		}
		symbolBlack, err := decodeSymbol(pieceTypeRules.Symbols.Black)
		if err != nil {
			return nil, err
		}
		pieceType.SetSymbols(symbolWhite, symbolBlack)
	}
	return pieceType, nil
}

func decodeSymbol(symbol string) (rune, error) {
	r, n := utf8.DecodeRuneInString(symbol)
	if n == 0 {
		return 0, fmt.Errorf("symbol cannot be empty")
	} else if r == utf8.RuneError {
		return 0, fmt.Errorf("symbol not an utf-8 character")
	} else if n != len(symbol) {
		return 0, fmt.Errorf("symbol too long (must be exactly one utf-8 character)")
	}
	return r, nil
}

func (c *rules) placePieces(state *mess.State) error {
	placementRules := map[color.Color]map[string]string{
		color.White: c.InitialState.WhitePieces,
		color.Black: c.InitialState.BlackPieces,
	}
	for color, pieces := range placementRules {
		player := state.Player(color)

		for squareString, pieceTypeName := range pieces {
			square, err := board.NewSquare(squareString)
			if err != nil {
				return fmt.Errorf("parsing square: %w", err)
			}

			pieceType, err := state.GetPieceType(pieceTypeName)
			if err != nil {
				return fmt.Errorf("getting piece type: %w", err)
			}

			piece := mess.NewPiece(pieceType, player)

			err = piece.PlaceOn(state.Board(), square)
			if err != nil {
				return fmt.Errorf("placing a piece: %w", err)
			}
		}
	}

	return nil
}
