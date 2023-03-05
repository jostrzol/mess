package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
)

type State struct {
	Board   PieceBoard
	Players map[color.Color]*Player
}

func NewState(board PieceBoard) *State {
	return &State{
		Board:   board,
		Players: NewPlayers(),
	}
}

func (g *State) String() string {
	return fmt.Sprintf("Board:\n%v", g.Board)
}

func (g *State) Player(color color.Color) *Player {
	player, ok := g.Players[color]
	if !ok {
		panic(fmt.Errorf("player of color %s not found", color))
	}
	return player
}

func (g *State) PieceAt(square *board.Square) (*Piece, error) {
	return g.Board.At(square)
}

func (g *State) PiecesPerPlayer() map[*Player][]*Piece {
	pieces := g.Board.AllItems()
	perPlayer := make(map[*Player][]*Piece, len(pieces))
	for _, player := range g.Players {
		perPlayer[player] = make([]*Piece, 0)
	}
	for _, piece := range pieces {
		owner := g.Player(piece.Color())
		perPlayer[owner] = append(perPlayer[owner], piece)
	}
	return perPlayer
}

func (g *State) Move(piece *Piece, square *board.Square) error {
	replaced, err := piece.MoveTo(square)
	if err != nil {
		return err
	}
	if replaced != nil {
		capturer := g.Player(piece.Color())
		capturer.Capture(replaced)
	}
	return nil
}

type Controller interface {
	DecideWinner(state *State) *Player
}
