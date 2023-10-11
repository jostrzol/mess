package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
)

func (c *rules) toEmptyGameState(ctx *hcl.EvalContext, interactor mess.Interactor) (*mess.Game, error) {
	brd, err := mess.NewPieceBoard(int(c.Board.Width), int(c.Board.Height))
	if err != nil {
		return nil, fmt.Errorf("creating new board: %w", err)
	}

	state := mess.NewState(brd)
	controller := newController(state, ctx, c)

	game := mess.NewGame(state, controller, interactor)

	stateValidators, err := controller.GetStateValidators()
	if err != nil {
		return nil, fmt.Errorf("parsing state validators: %w", err)
	}
	for _, validator := range stateValidators {
		state.AddStateValidator(validator)
	}

	for _, pieceTypeRules := range c.PieceTypes.PieceTypes {
		pieceType := mess.NewPieceType(pieceTypeRules.Name)
		for _, motionRules := range pieceTypeRules.Motions {
			moveGenerator, err := controller.GetCustomFuncAsGenerator(motionRules.GeneratorName)
			if err != nil {
				return nil, err
			}
			actions := make([]mess.MoveAction, 0, len(motionRules.ActionNames))
			for _, actionName := range motionRules.ActionNames {
				action, err := controller.GetCustomFuncAsAction(actionName)
				if err != nil {
					return nil, err
				}
				actions = append(actions, action)
			}
			var action mess.MoveAction
			if len(actions) != 0 {
				action = func(piece *mess.Piece, from board.Square, to board.Square) {
					for _, action := range actions {
						action(piece, from, to)
					}
				}
			}
			pieceType.AddMoveGenerator(
				mess.MoveGenerator{
					Name: motionRules.GeneratorName,
					Generate: func(piece *mess.Piece) ([]board.Square, mess.MoveAction) {
						return moveGenerator(piece), action
					},
				},
			)
		}
		state.AddPieceType(pieceType)
	}

	initializeContext(ctx, game)
	return game, nil
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
