package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
	"github.com/jostrzol/mess/pkg/iter"
)

type State struct {
	board         *PieceBoard
	players       map[color.Color]*Player
	currentPlayer *Player
	record        []RecordedMove
	isRecording   bool
}

func NewState(board *PieceBoard) *State {
	players := NewPlayers(board)
	state := &State{
		board:         board,
		players:       players,
		currentPlayer: players[color.White],
		record:        []RecordedMove{},
		isRecording:   true,
	}
	board.Observe(state)
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

func (g *State) Handle(event event.Event) {
	if !g.isRecording {
		return
	}

	switch e := event.(type) {
	case PieceMoved:
		g.record = append(g.record, RecordedMove{
			Move:     Move(e),
			Captures: map[*Piece]struct{}{},
		})
	case PieceCaptured:
		g.recordCapture(&e)
	}
}

func (g *State) recordCapture(event *PieceCaptured) {
	if len(g.record) == 0 {
		panic(fmt.Errorf("tried to record a capture, but no moved was recorded ealier"))
	}
	lastMove := g.record[len(g.record)-1]
	lastMove.Captures[event.Piece] = struct{}{}
}

func (g *State) Undo() *RecordedMove {
	if len(g.record) == 0 {
		return nil
	}
	lastMove := g.record[len(g.record)-1]

	g.isRecording = false
	defer func() { g.isRecording = true }()

	err := lastMove.Piece.MoveTo(lastMove.From)
	if err != nil {
		panic(err)
	}

	for c := range lastMove.Captures {
		err := g.board.Place(c, c.Square())
		if err != nil {
			panic(err)
		}
	}
	g.record = g.record[:len(g.record)-1]
	return &lastMove
}

type RecordedMove struct {
	Move
	Captures map[*Piece]struct{}
}

type Controller interface {
	PickWinner(state *State) (bool, *Player)
}
