package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
	"github.com/jostrzol/mess/pkg/iter"
)

type State struct {
	board             *PieceBoard
	players           map[color.Color]*Player
	currentPlayer     *Player
	record            []RecordedMove
	isRecording       bool
	validators        chainStateValidators
	validMoves        []Move
	turnNumber        int
	isGeneratingMoves bool
	pieceTypes        []*PieceType
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

func (s *State) String() string {
	return fmt.Sprintf("Board:\n%v\nCurrent player: %v\n", s.board, s.currentPlayer)
}

func (s *State) PrettyString() string {
	return fmt.Sprintf("%v\nCurrent player: %v", s.board.PrettyString(), s.currentPlayer)
}

func (s *State) Board() *PieceBoard {
	return s.board
}

func (s *State) Players() <-chan *Player {
	return iter.FromValues(s.players)
}

func (s *State) Player(color color.Color) *Player {
	player, ok := s.players[color]
	if !ok {
		panic(fmt.Errorf("player of color %s not found", color))
	}
	return player
}

func (s *State) CurrentPlayer() *Player {
	return s.currentPlayer
}

func (s *State) CurrentOpponent() *Player {
	return s.OpponentTo(s.currentPlayer)
}

func (s *State) OpponentTo(player *Player) *Player {
	var opponentsColor color.Color
	switch player.Color() {
	case color.White:
		opponentsColor = color.Black
	case color.Black:
		opponentsColor = color.White
	}
	return s.Player(opponentsColor)
}

func (s *State) EndTurn() {
	s.currentPlayer = s.CurrentOpponent()
	s.turnNumber++
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

func (s *State) AddStateValidator(validator StateValidator) {
	s.validators = append(s.validators, validator)
}

func (s *State) ValidMoves() []Move {
	if s.validMoves == nil {
		s.generateValidMoves()
	}
	return s.validMoves
}

func (s *State) generateValidMoves() {
	result := make([]Move, 0)
	s.isGeneratingMoves = true
	defer func() { s.isGeneratingMoves = false }()

	moves := s.currentPlayer.moves()
	for _, move := range moves {
		err := move.PerformWithoutAction()

		isValid := false
		if err != nil {
			fmt.Printf("error performing move for validation: %v", err)
		} else {
			isValid = s.validators.Validate(s, &move)
		}
		s.UndoTurn()

		if isValid {
			result = append(result, move)
			fmt.Printf("DEBUG: generated move: %v\n", &move)
		}
	}
	s.validMoves = result
}

func (s *State) Handle(event event.Event) {
	if !s.isRecording {
		return
	}

	switch e := event.(type) {
	case PieceMoved:
		s.record = append(s.record, RecordedMove{
			Piece:      e.Piece,
			From:       e.From,
			To:         e.To,
			TurnNumber: s.turnNumber,
			Captures:   map[*Piece]struct{}{},
		})
	case PieceCaptured:
		s.recordCapture(&e)
	}

	s.validMoves = nil
}

func (s *State) recordCapture(event *PieceCaptured) {
	if len(s.record) == 0 {
		panic(fmt.Errorf("tried to record a capture, but no moved was recorded ealier"))
	}
	lastMove := s.record[len(s.record)-1]
	lastMove.Captures[event.Piece] = struct{}{}
}

func (s *State) UndoTurn() {
	for len(s.record) > 0 && s.record[len(s.record)-1].TurnNumber == s.turnNumber {
		s.undoOne()
	}
}

func (s *State) undoOne() {
	if len(s.record) == 0 {
		return
	}
	lastMove := s.record[len(s.record)-1]

	s.isRecording = false
	defer func() { s.isRecording = true }()

	err := lastMove.Piece.MoveTo(lastMove.From)
	if err != nil {
		panic(err)
	}

	for c := range lastMove.Captures {
		err := s.board.Place(c, c.Square())
		if err != nil {
			panic(err)
		}
	}
	s.record = s.record[:len(s.record)-1]
}

func (s *State) Record() []RecordedMove {
	return s.record
}

type RecordedMove struct {
	Piece      *Piece
	From       board.Square
	To         board.Square
	TurnNumber int
	Captures   map[*Piece]struct{}
}

func (s *State) IsGeneratingMoves() bool {
	return s.isGeneratingMoves
}

func (s *State) AddPieceType(pieceType *PieceType) {
	for _, pt := range s.pieceTypes {
		if pt == pieceType {
			return
		}
	}
	s.pieceTypes = append(s.pieceTypes, pieceType)
}

func (s *State) PieceTypes() []*PieceType {
	return s.pieceTypes
}
