package game

import (
	"fmt"

	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/player"
)

type State struct {
	Board   board.Board
	Players map[player.Color]*player.Player
}

func (g *State) GetPlayer(color player.Color) (*player.Player, error) {
	player, ok := g.Players[color]
	if !ok {
		return nil, fmt.Errorf("player of color %q not found", color)
	}
	return player, nil
}

func (g *State) PiecesPerPlayer() map[*player.Player][]board.PieceOnSquare {
	pieces := g.Board.AllPieces()
	perPlayer := make(map[*player.Player][]board.PieceOnSquare, len(pieces))
	for _, player := range g.Players {
		perPlayer[player] = make([]board.PieceOnSquare, 0)
	}
	for _, pieceOnSquare := range pieces {
		owner := pieceOnSquare.Piece.Owner
		perPlayer[owner] = append(perPlayer[owner], pieceOnSquare)
	}
	return perPlayer
}

type Controller interface {
	DecideWinner(state *State) (*player.Player, error)
}
