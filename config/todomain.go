package config

import (
	"fmt"

	"github.com/jostrzol/mess/game"
	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/player"
)

func (c *config) ToGame() (*game.GameState, error) {
	board, err := board.NewBoard(int(c.Board.Width), int(c.Board.Height))
	if err != nil {
		return nil, fmt.Errorf("creating board: %w", err)
	}

	pieceTypes := make(map[string]*piece.PieceType, len(c.PieceTypes.PieceTypes))
	for _, pieceType := range c.PieceTypes.PieceTypes {
		pieceTypes[pieceType.Name] = &piece.PieceType{
			Name: pieceType.Name,
		}
	}

	players := player.NewPlayers()

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

func placePieces(state *game.GameState, pieces []piecesConfig, pieceTypes map[string]*piece.PieceType) error {
	for _, pieces := range pieces {
		color, err := player.ColorString(pieces.PlayerColor)
		if err != nil {
			return fmt.Errorf("parsing player color: %w", err)
		}
		player, err := state.GetPlayer(color)
		if err != nil {
			return fmt.Errorf("getting player: %w", err)
		}

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

			piece := &piece.Piece{
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
