package game

import (
	"fmt"

	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/player"
)

type PieceBoard = board.Board[*piece.Piece]

type State struct {
	Board   PieceBoard
	Players map[player.Color]*player.Player
}

func (g *State) GetPlayer(color player.Color) (*player.Player, error) {
	player, ok := g.Players[color]
	if !ok {
		return nil, fmt.Errorf("player of color %q not found", color)
	}
	return player, nil
}

func (g *State) PiecesPerPlayer() map[*player.Player][]*piece.Piece {
	pieces := g.Board.AllItems()
	perPlayer := make(map[*player.Player][]*piece.Piece, len(pieces))
	for _, player := range g.Players {
		perPlayer[player] = make([]*piece.Piece, 0)
	}
	for _, piece := range pieces {
		owner := piece.Owner
		perPlayer[owner] = append(perPlayer[owner], piece)
	}
	return perPlayer
}

type Controller interface {
	DecideWinner(state *State) (*player.Player, error)
}
