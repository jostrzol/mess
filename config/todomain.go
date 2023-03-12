package config

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
)

func (c *config) toGameState() (*mess.State, error) {
	board, err := mess.NewPieceBoard(int(c.Board.Width), int(c.Board.Height))
	if err != nil {
		return nil, fmt.Errorf("creating board: %w", err)
	}
	state := mess.NewState(board)

	pieceTypes := make(map[string]*mess.PieceType, len(c.PieceTypes.PieceTypes))
	for _, pieceTypeConfig := range c.PieceTypes.PieceTypes {
		pieceType := mess.NewPieceType(pieceTypeConfig.Name)
		for _, motionConfig := range pieceTypeConfig.Motions {
			moveGenerator, err := c.Functions.GetCustomFuncAsGenerator(motionConfig.GeneratorName)
			if err != nil {
				return nil, err
			}
			pieceType.AddMoveGenerator(moveGenerator)
		}
		pieceTypes[pieceType.Name()] = pieceType
	}

	err = placePieces(state, c.InitialState.Pieces, pieceTypes)
	if err != nil {
		return nil, fmt.Errorf("placing initial pieces: %w", err)
	}
	return state, nil
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

func (c *config) toController() mess.Controller {
	return c.Functions
}
