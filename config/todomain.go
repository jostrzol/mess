package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
)

func (c *config) toGameState(ctx *hcl.EvalContext) (*mess.Game, error) {
	brd, err := mess.NewPieceBoard(int(c.Board.Width), int(c.Board.Height))
	if err != nil {
		return nil, fmt.Errorf("creating new board: %w", err)
	}

	state := mess.NewState(brd)
	populateContextWithState(ctx, state)

	controller := controller{
		ctx:    ctx,
		config: &c.Functions,
	}

	stateValidators, err := controller.GetStateValidators(state)
	if err != nil {
		return nil, fmt.Errorf("parsing state validators: %w", err)
	}
	for _, validator := range stateValidators {
		state.AddStateValidator(validator)
	}

	pieceTypes := make(map[string]*mess.PieceType, len(c.PieceTypes.PieceTypes))
	for _, pieceTypeConfig := range c.PieceTypes.PieceTypes {
		pieceType := mess.NewPieceType(pieceTypeConfig.Name)
		for _, motionConfig := range pieceTypeConfig.Motions {
			moveGenerator, err := controller.GetCustomFuncAsGenerator(motionConfig.GeneratorName, state)
			if err != nil {
				return nil, err
			}
			actions := make([]mess.MoveAction, 0, len(motionConfig.ActionNames))
			for _, actionName := range motionConfig.ActionNames {
				action, err := controller.GetCustomFuncAsAction(actionName, state)
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
			pieceType.AddMoveGenerator(func(piece *mess.Piece) ([]board.Square, mess.MoveAction) {
				return moveGenerator(piece), action
			})
		}
		pieceTypes[pieceType.Name()] = pieceType
	}

	err = placePieces(state, c.InitialState.Pieces, pieceTypes)
	if err != nil {
		return nil, fmt.Errorf("placing initial pieces: %w", err)
	}

	return mess.NewGame(state, &controller), nil
}

func placePieces(state *mess.State, pieces []piecesConfig, pieceTypes map[string]*mess.PieceType) error {
	for _, pieces := range pieces {
		color, err := color.ColorString(pieces.PlayerColor)
		if err != nil {
			return fmt.Errorf("parsing player color: %w", err)
		}
		player := state.Player(color)

		for _, piecePlacement := range pieces.Placements {
			squareString := piecePlacement.Name
			square, err := board.NewSquare(squareString)
			if err != nil {
				return fmt.Errorf("parsing square: %w", err)
			}

			pieceTypeName, diags := piecePlacement.Expr.Value(nil)
			if diags.HasErrors() {
				return fmt.Errorf("parsing piece type: %w", diags)
			}
			pieceType := pieceTypes[pieceTypeName.AsString()]
			if pieceType == nil {
				return fmt.Errorf("piece type %q not defined", pieceTypeName.AsString())
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
