package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
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
	validators    chainStateValidators
	validMoves    []Move
	turnNumber    int
}

func NewState(board *PieceBoard) *State {
	players := NewPlayers(board)
	state := &State{
		board:         board,
		players:       players,
		currentPlayer: players[color.White],
		record:        []RecordedMove{},
		isRecording:   true,
		turnNumber:    1,
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
	return g.OpponentTo(g.currentPlayer)
}

func (g *State) OpponentTo(player *Player) *Player {
	var opponentsColor color.Color
	switch player.Color() {
	case color.White:
		opponentsColor = color.Black
	case color.Black:
		opponentsColor = color.White
	}
	return g.Player(opponentsColor)
}

func (g *State) EndTurn() {
	g.currentPlayer = g.CurrentOpponent()
	g.turnNumber += 1
}

type StateValidator func(*State, *Move) bool
type chainStateValidators []StateValidator

func (validators chainStateValidators) Validate(state *State, move *Move) bool {
	for _, validator := range validators {
		if !validator(state, move) {
			return false
		}
	}
	return true
}

func (g *State) AddStateValidator(validator StateValidator) {
	g.validators = append(g.validators, validator)
}

func (g *State) ValidMoves() []Move {
	if g.validMoves == nil {
		g.generateValidMoves()
	}
	return g.validMoves
}

func (g *State) generateValidMoves() {
	result := make([]Move, 0)
	moves := g.currentPlayer.moves()
	for _, move := range moves {
		err := move.PerformWithoutAction()

		isValid := false
		if err != nil {
			fmt.Printf("error performing move for validation: %v", err)
		} else {
			isValid = g.validators.Validate(g, &move)
		}
		g.UndoTurn()

		if isValid {
			result = append(result, move)
			fmt.Printf("DEBUG: generated move: %v\n", &move)
		}
	}
	g.validMoves = result
}

func (g *State) Handle(event event.Event) {
	if !g.isRecording {
		return
	}

	switch e := event.(type) {
	case PieceMoved:
		g.record = append(g.record, RecordedMove{
			Piece:      e.Piece,
			From:       e.From,
			To:         e.To,
			TurnNumber: g.turnNumber,
			Captures:   map[*Piece]struct{}{},
		})
	case PieceCaptured:
		g.recordCapture(&e)
	}

	g.validMoves = nil
}

func (g *State) recordCapture(event *PieceCaptured) {
	if len(g.record) == 0 {
		panic(fmt.Errorf("tried to record a capture, but no moved was recorded ealier"))
	}
	lastMove := g.record[len(g.record)-1]
	lastMove.Captures[event.Piece] = struct{}{}
}

func (g *State) UndoTurn() {
	for len(g.record) > 0 && g.record[len(g.record)-1].TurnNumber == g.turnNumber {
		g.undoOne()
	}
}

func (g *State) undoOne() {
	if len(g.record) == 0 {
		return
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
}

func (g *State) Record() []RecordedMove {
	return g.record
}

type RecordedMove struct {
	Piece      *Piece
	From       board.Square
	To         board.Square
	TurnNumber int
	Captures   map[*Piece]struct{}
}

type Controller interface {
	PickWinner(state *State) (bool, *Player)
}
