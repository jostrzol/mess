package game

import (
	"fmt"

	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/piece/color"
	"github.com/jostrzol/mess/game/player"
)

type State struct {
	Board   piece.Board
	Players map[color.Color]*player.Player
}

func NewState(board piece.Board) *State {
	return &State{
		Board:   board,
		Players: player.NewPlayers(),
	}
}

func (g *State) String() string {
	return fmt.Sprintf("Board:\n%v", g.Board)
}

func (g *State) Player(color color.Color) *player.Player {
	player, ok := g.Players[color]
	if !ok {
		panic(fmt.Errorf("player of color %s not found", color))
	}
	return player
}

func (g *State) PieceAt(square *board.Square) (*piece.Piece, error) {
	return g.Board.At(square)
}

func (g *State) PiecesPerPlayer() map[*player.Player][]*piece.Piece {
	pieces := g.Board.AllItems()
	perPlayer := make(map[*player.Player][]*piece.Piece, len(pieces))
	for _, player := range g.Players {
		perPlayer[player] = make([]*piece.Piece, 0)
	}
	for _, piece := range pieces {
		owner := g.Player(piece.Color())
		perPlayer[owner] = append(perPlayer[owner], piece)
	}
	return perPlayer
}

func (g *State) Move(piece *piece.Piece, square *board.Square) error {
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
	DecideWinner(state *State) *player.Player
}
