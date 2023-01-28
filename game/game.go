package game

import "fmt"

type GameState struct {
	Board   Board
	Players map[Color]*Player
}

func (g *GameState) GetPlayer(color Color) (*Player, error) {
	player, ok := g.Players[color]
	if !ok {
		return nil, fmt.Errorf("player of color %q not found", color)
	}
	return player, nil
}

func (g *GameState) PiecesPerPlayer() map[*Player][]*PieceOnSquare {
	pieces := g.Board.Pieces()
	perPlayer := make(map[*Player][]*PieceOnSquare, len(pieces))
	for _, player := range g.Players {
		perPlayer[player] = make([]*PieceOnSquare, 0)
	}
	for _, pieceOnSquare := range pieces {
		owner := pieceOnSquare.Piece.Owner
		perPlayer[owner] = append(perPlayer[owner], &pieceOnSquare)
	}
	return perPlayer
}

type GameController interface {
	DecideWinner(state *GameState) (*Player, error)
}
