package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/gen"
)

type State struct {
	board   *PieceBoard
	players map[color.Color]*Player
}

func NewState(board *PieceBoard) *State {
	state := &State{
		board:   board,
		players: NewPlayers(board),
	}
	return state
}

func (g *State) String() string {
	return fmt.Sprintf("Board:\n%v", g.board)
}

func (g *State) Board() *PieceBoard {
	return g.board
}

func (g *State) Players() <-chan *Player {
	return gen.FromValues(g.players)
}

func (g *State) Player(color color.Color) *Player {
	player, ok := g.players[color]
	if !ok {
		panic(fmt.Errorf("player of color %s not found", color))
	}
	return player
}

func (g *State) PiecesPerPlayer() map[*Player][]*Piece {
	pieces := g.board.AllPieces()
	perPlayer := make(map[*Player][]*Piece, len(pieces))
	for _, player := range g.players {
		perPlayer[player] = make([]*Piece, 0)
	}
	for _, piece := range pieces {
		owner := g.Player(piece.Color())
		perPlayer[owner] = append(perPlayer[owner], piece)
	}
	return perPlayer
}

type Controller interface {
	DecideWinner(state *State) *Player
}
