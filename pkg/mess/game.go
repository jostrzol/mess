package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/iter"
)

type State struct {
	board         *PieceBoard
	players       map[color.Color]*Player
	currentPlayer *Player
}

func NewState(board *PieceBoard) *State {
	players := NewPlayers(board)
	state := &State{
		board:         board,
		players:       players,
		currentPlayer: players[color.White],
	}
	return state
}

func (g *State) String() string {
	return fmt.Sprintf("Board:\n%v\nCurrent player: %v\n", g.board, g.currentPlayer)
}

func (g *State) PrettyString() string {
	return fmt.Sprintf("%v\nCurrent player: %v", g.board.PrettyString(), g.currentPlayer)
}

func (g *State) Board() *PieceBoard {
	return g.board
}

func (g *State) Players() <-chan *Player {
	return iter.FromValues(g.players)
}

func (g *State) Player(color color.Color) *Player {
	player, ok := g.players[color]
	if !ok {
		panic(fmt.Errorf("player of color %s not found", color))
	}
	return player
}

func (g *State) CurrentPlayer() *Player {
	return g.currentPlayer
}

func (g *State) CurrentOpponent() *Player {
	var opponentsColor color.Color
	switch g.CurrentPlayer().Color() {
	case color.White:
		opponentsColor = color.Black
	case color.Black:
		opponentsColor = color.White
	}
	return g.Player(opponentsColor)
}

func (g *State) EndTurn() {
	g.currentPlayer = g.CurrentOpponent()
}

type Controller interface {
	PickWinner(state *State) (bool, *Player)
}
