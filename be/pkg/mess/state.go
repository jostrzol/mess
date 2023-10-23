package mess

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
	"golang.org/x/exp/maps"
)

type State struct {
	board             *PieceBoard
	players           map[color.Color]*Player
	currentPlayer     *Player
	record            []Turn
	isRecording       bool
	validators        chainStateValidators
	validMoves        []MoveGroup
	turnNumber        int
	isGeneratingMoves bool
	pieceTypes        map[string]*PieceType
}

func NewState(board *PieceBoard) *State {
	players := NewPlayers(board)
	state := &State{
		board:         board,
		players:       players,
		currentPlayer: players[color.White],
		record:        []Turn{},
		isRecording:   true,
		turnNumber:    0,
		pieceTypes:    make(map[string]*PieceType),
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

func (s *State) Players() []*Player {
	return maps.Values(s.players)
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

func (s *State) ValidMoves() []MoveGroup {
	if s.validMoves == nil {
		s.generateValidMoves()
	}
	return s.validMoves
}

func (s *State) generateValidMoves() {
	s.isGeneratingMoves = true
	defer func() { s.isGeneratingMoves = false }()

	moveGroups := s.currentPlayer.Moves()
	result := make([]MoveGroup, 0, len(moveGroups))
	for len(moveGroups) != 0 {
		moveGroup := moveGroups[0]
		moveGroups = moveGroups[1:]

		moveGroup = moveGroup.FilterMoves(s.validateMove)

		if moveGroup.Length() != 0 {
			result = append(result, moveGroup)
		}
	}
	s.validMoves = result
}

func (s *State) validateMove(move *Move) bool {
	err := move.Perform()
	if err != nil {
		fmt.Printf("error validating move: %v\n", err)
		return false
	}
	isValid := s.validators.Validate(s, move)
	s.UndoTurn()

	return isValid
}

func (s *State) Handle(event event.Event) {
	if !s.isRecording {
		return
	}

	_, isPiecePlaced := event.(PiecePlaced)
	if s.turnNumber == 0 && len(s.record) == 0 && isPiecePlaced {
		// don't record initial setup
		return
	}

	var turn Turn
	if len(s.record) != s.turnNumber {
		// not the first move in the round -> load it
		turn = s.record[s.turnNumber]
	}

	turn = append(turn, event)

	if len(s.record) == s.turnNumber {
		// the first move in the round -> append it
		s.record = append(s.record, turn)
	} else {
		// not the first move in the round -> modify it
		s.record[s.turnNumber] = turn
	}

	s.validMoves = nil
}

func (s *State) UndoTurn() {
	if len(s.record) == 0 {
		return
	}

	turn := s.record[s.turnNumber]
	s.record = s.record[:s.turnNumber]

	s.isRecording = false
	defer func() { s.isRecording = true }()

	for i := range turn {
		iRev := len(turn) - i - 1
		event := turn[iRev]

		switch e := event.(type) {
		case PieceMoved:
			err := e.Piece.MoveTo(e.From)
			if err != nil {
				panic(err)
			}
		case PiecePlaced:
			err := e.Piece.Remove()
			if err != nil {
				panic(err)
			}
		case PieceRemoved:
			err := e.Piece.PlaceOn(s.board, e.Square)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (s *State) Record() []Turn {
	return s.record
}

type Turn []event.Event

func (t Turn) FirstMove() *PieceMoved {
	for _, event := range t {
		if e, ok := event.(PieceMoved); ok {
			return &e
		}
	}
	return nil
}

func (s *State) IsGeneratingMoves() bool {
	return s.isGeneratingMoves
}

func (s *State) AddPieceType(pieceType *PieceType) {
	s.pieceTypes[pieceType.Name()] = pieceType
}

func (s *State) GetPieceType(name string) (*PieceType, error) {
	pieceType, ok := s.pieceTypes[name]
	if !ok {
		return nil, fmt.Errorf("piece type %q not defined", name)
	}
	return pieceType, nil
}

func (s *State) PieceTypes() []*PieceType {
	result := make([]*PieceType, 0, len(s.pieceTypes))
	for _, pieceType := range s.pieceTypes {
		result = append(result, pieceType)
	}
	return result
}
