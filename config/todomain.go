package config

import (
	"fmt"

	"github.com/jostrzol/mess/game"
	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/piece/color"
)

func (c *config) toGameState(state *game.State) error {
	board, err := board.NewBoard[*piece.Piece](int(c.Board.Width), int(c.Board.Height))
	if err != nil {
		return fmt.Errorf("creating board: %w", err)
	}
	*state = *game.NewState(board)

	pieceTypes := make(map[string]*piece.Type, len(c.PieceTypes.PieceTypes))
	for _, pieceTypeConfig := range c.PieceTypes.PieceTypes {
		pieceType := piece.NewType(pieceTypeConfig.Name)
		for _, motionConfig := range pieceTypeConfig.Motions {
			motionGenerator, err := c.Functions.GetCustomFuncAsGenerator(motionConfig.GeneratorName)
			if err != nil {
				return err
			}
			pieceType.AddMotionGenerator(motionGenerator)
		}
		pieceTypes[pieceType.Name] = pieceType
	}

	err = placePieces(state, c.InitialState.Pieces, pieceTypes)
	if err != nil {
		return fmt.Errorf("placing initial pieces: %w", err)
	}
	return nil
}

func placePieces(state *game.State, pieces []piecesConfig, pieceTypes map[string]*piece.Type) error {
	for _, pieces := range pieces {
		color, err := color.ColorString(pieces.PlayerColor)
		if err != nil {
			return fmt.Errorf("parsing player color: %w", err)
		}
		player := state.GetPlayer(color)

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

			piece := piece.NewPiece(pieceType, player)

			err = piece.PlaceOn(state.Board, square)
			if err != nil {
				return fmt.Errorf("placing a piece: %w", err)
			}
		}
	}
	return nil
}

func (c *config) toController() game.Controller {
	return c.Functions
}
