package config

import (
	"fmt"

	"github.com/jostrzol/mess/game"
	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/player"
)

func (c *config) ToGame() (*game.State, error) {
	board, err := board.NewBoard[*piece.Piece](int(c.Board.Width), int(c.Board.Height))
	if err != nil {
		return nil, fmt.Errorf("creating board: %w", err)
	}

	pieceTypes := make(map[string]*piece.Type, len(c.PieceTypes.PieceTypes))
	for _, pieceType := range c.PieceTypes.PieceTypes {
		pieceTypes[pieceType.Name] = &piece.Type{
			Name: pieceType.Name,
		}
	}

	players := player.NewPlayers()

	state := &game.State{
		Board:   board,
		Players: players,
	}

	err = placePieces(state, c.InitialState.Pieces, pieceTypes)
	if err != nil {
		return nil, fmt.Errorf("placing initial pieces: %w", err)
	}

	return state, nil
}

func placePieces(state *game.State, pieces []piecesConfig, pieceTypes map[string]*piece.Type) error {
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

			piece := piece.NewPiece(pieceType, player)

			err = piece.PlaceOn(state.Board, square)
			if err != nil {
				return fmt.Errorf("placing a piece: %w", err)
			}
		}
	}
	return nil
}

func (c *config) ToController() game.Controller {
	return c.Functions
}
