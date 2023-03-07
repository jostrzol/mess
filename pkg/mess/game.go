package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/gen"
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
	return fmt.Sprintf("Board:\n%v\nCurrent player: %v", g.board, g.currentPlayer)
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

func (g *State) CurrentPlayer() *Player {
	return g.currentPlayer
}

func (g *State) EndTurn() {
	g.currentPlayer = g.otherPlayer(g.currentPlayer)
}

func (g *State) otherPlayer(player *Player) *Player {
	var otherColor color.Color
	switch player.Color() {
	case color.White:
		otherColor = color.Black
	case color.Black:
		otherColor = color.White
	}
	return g.Player(otherColor)
}

type Controller interface {
	DecideWinner(state *State) *Player
}
