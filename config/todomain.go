package config

import (
	"fmt"

	"github.com/jostrzol/mess/game"
)

func (c *config) ToGame() (*game.GameState, error) {
	board, err := game.NewBoard(int(c.Board.Width), int(c.Board.Height))
	if err != nil {
		return nil, fmt.Errorf("creating board: %w", err)
	}

	pieceTypes := make(map[string]*game.PieceType, len(c.PieceTypes.PieceTypes))
	for _, pieceType := range c.PieceTypes.PieceTypes {
		pieceTypes[pieceType.Name] = &game.PieceType{
			Name: pieceType.Name,
		}
	}

	players := game.NewPlayers()

	state := &game.GameState{
		Board:   board,
		Players: players,
	}

	err = placePieces(state, c.InitialState.Pieces, pieceTypes)
	if err != nil {
		return nil, fmt.Errorf("placing initial pieces: %w", err)
	}

	return state, nil
}

func placePieces(state *game.GameState, pieces []piecesConfig, pieceTypes map[string]*game.PieceType) error {
	for _, pieces := range pieces {
		color, err := game.ColorString(pieces.PlayerColor)
		if err != nil {
			return fmt.Errorf("parsing player color: %w", err)
		}
		player, err := state.GetPlayer(color)
		if err != nil {
			return fmt.Errorf("getting player: %w", err)
		}

		for _, piecePlacement := range pieces.Placements {
			squareString := piecePlacement.Name
			square, err := game.ParseSquare(squareString)
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

			piece := &game.Piece{
				Type:  pieceType,
				Owner: player,
			}

			err = state.Board.Place(piece, square)
			if err != nil {
				return fmt.Errorf("placing a piece: %w", err)
			}
		}
	}
	return nil
}

func (c *config) ToController() game.GameController {
	return c.Functions
}
