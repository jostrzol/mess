package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StateSuite struct {
	suite.Suite
	state *State
}

func (s *StateSuite) SetupTest() {
	board, err := NewPieceBoard(8, 8)
	s.NoError(err)
	s.state = NewState(board)
}

func (s *StateSuite) TestGetPlayer() {
	for _, color := range color.ColorValues() {
		s.Run(color.String(), func() {
			player := s.state.Player(color)
			s.Equal(player.Color(), color)
		})
	}
}

func (s *StateSuite) TestGetPlayerNotFound() {
	s.Panics(func() {
		s.state.Player(color.Color(-1))
	})
}

func (s *StateSuite) TestEndTurn() {
	firstTurnPlayer := s.state.CurrentPlayer()
	s.Equal(s.state.Player(color.White), firstTurnPlayer)

	s.state.EndTurn()

	secondTurnPlayer := s.state.CurrentPlayer()
	s.Equal(s.state.Player(color.Black), secondTurnPlayer)
}

func (s *StateSuite) TestRecordMove() {
	var rook, knight, king Piece
	rookSquare := boardtest.NewSquare("A1")
	knightSquare := boardtest.NewSquare("A2")
	kingSquare := boardtest.NewSquare("A3")
	emptySquare1 := boardtest.NewSquare("B1")
	emptySquare2 := boardtest.NewSquare("B2")

	type move struct {
		*Piece
		board.Square
	}
	tests := []struct {
		name     string
		moves    []move
		expected []RecordedMove
	}{
		{
			name:  "One",
			moves: []move{{&rook, emptySquare1}},
			expected: []RecordedMove{
				{
					Piece:      &rook,
					From:       rookSquare,
					To:         emptySquare1,
					Captures:   map[*Piece]struct{}{},
					TurnNumber: 1,
				},
			},
		},
		{
			name:  "Two",
			moves: []move{{&rook, emptySquare1}, {&rook, emptySquare2}},
			expected: []RecordedMove{
				{
					Piece:      &rook,
					From:       rookSquare,
					To:         emptySquare1,
					Captures:   map[*Piece]struct{}{},
					TurnNumber: 1,
				},
				{
					Piece:      &rook,
					From:       emptySquare1,
					To:         emptySquare2,
					Captures:   map[*Piece]struct{}{},
					TurnNumber: 1,
				},
			},
		},
		{
			name:  "Capture",
			moves: []move{{&rook, knightSquare}},
			expected: []RecordedMove{
				{
					Piece: &rook,
					From:  rookSquare,
					To:    knightSquare,
					Captures: map[*Piece]struct{}{
						&knight: {},
					},
					TurnNumber: 1,
				},
			},
		},
		{
			name:  "DoubleCapture",
			moves: []move{{&rook, knightSquare}, {&rook, kingSquare}},
			expected: []RecordedMove{
				{
					Piece: &rook,
					From:  rookSquare,
					To:    knightSquare,
					Captures: map[*Piece]struct{}{
						&knight: {},
					},
					TurnNumber: 1,
				},
				{
					Piece: &rook,
					From:  knightSquare,
					To:    kingSquare,
					Captures: map[*Piece]struct{}{
						&king: {},
					},
					TurnNumber: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		s.SetupTest()
		s.Run(tt.name, func() {
			rook = *NewPiece(Rook(s.T()), s.state.currentPlayer)
			rook.PlaceOn(s.state.board, rookSquare)
			knight = *NewPiece(Knight(s.T()), s.state.currentPlayer)
			knight.PlaceOn(s.state.board, knightSquare)
			king = *NewPiece(Knight(s.T()), s.state.currentPlayer)
			king.PlaceOn(s.state.board, kingSquare)

			for _, move := range tt.moves {
				err := move.Piece.MoveTo(move.Square)
				s.NoError(err)
			}

			s.Equal(s.state.record, tt.expected)
		})
	}
}

func (s *StateSuite) TestRecordTurnNumber() {
	rook := *NewPiece(Rook(s.T()), s.state.currentPlayer)
	rook.PlaceOn(s.state.board, boardtest.NewSquare("A1"))

	err := rook.MoveTo(boardtest.NewSquare("A2"))
	s.NoError(err)

	s.state.EndTurn()

	err = rook.MoveTo(boardtest.NewSquare("A3"))
	s.NoError(err)

	err = rook.MoveTo(boardtest.NewSquare("B3"))
	s.NoError(err)

	s.state.EndTurn()

	err = rook.MoveTo(boardtest.NewSquare("A4"))
	s.NoError(err)

	s.Equal(s.state.record[0].TurnNumber, 1)
	s.Equal(s.state.record[1].TurnNumber, 2)
	s.Equal(s.state.record[2].TurnNumber, 2)
	s.Equal(s.state.record[3].TurnNumber, 3)
}

func (s *StateSuite) TestUndoNothing() {
	s.state.UndoTurn()
}

func (s *StateSuite) TestUndo() {
	rook := NewPiece(Rook(s.T()), s.state.currentPlayer)
	a1 := boardtest.NewSquare("A1")
	rook.PlaceOn(s.state.board, a1)

	a2 := boardtest.NewSquare("A2")
	err := rook.MoveTo(a2)
	s.NoError(err)

	s.state.UndoTurn()

	pieceA1, err := s.state.Board().At(a1)
	s.NoError(err)
	s.Equal(rook, pieceA1)

	pieceA2, err := s.state.Board().At(a2)
	s.NoError(err)
	s.Nil(pieceA2)
}

func (s *StateSuite) TestUndoDoubleMove() {
	rook := NewPiece(Rook(s.T()), s.state.currentPlayer)
	a1 := boardtest.NewSquare("A1")
	rook.PlaceOn(s.state.board, a1)

	a2 := boardtest.NewSquare("A2")
	err := rook.MoveTo(a2)
	s.NoError(err)

	s.state.EndTurn()

	a3 := boardtest.NewSquare("A3")
	err = rook.MoveTo(a3)
	s.NoError(err)

	s.state.UndoTurn()

	pieceA1, err := s.state.Board().At(a1)
	s.NoError(err)
	s.Nil(pieceA1)

	pieceA2, err := s.state.Board().At(a2)
	s.NoError(err)
	s.Equal(rook, pieceA2)

	pieceA3, err := s.state.Board().At(a3)
	s.NoError(err)
	s.Nil(pieceA3)
}

func (s *StateSuite) TestUndoCapture() {
	rook := NewPiece(Rook(s.T()), s.state.currentPlayer)
	a1 := boardtest.NewSquare("A1")
	rook.PlaceOn(s.state.board, a1)

	knight := NewPiece(Knight(s.T()), s.state.currentPlayer)
	a2 := boardtest.NewSquare("A2")
	knight.PlaceOn(s.state.board, a2)

	err := rook.MoveTo(a2)
	s.NoError(err)

	s.state.UndoTurn()

	pieceA1, err := s.state.Board().At(a1)
	s.NoError(err)
	s.Equal(rook, pieceA1)

	pieceA2, err := s.state.Board().At(a2)
	s.NoError(err)
	s.Equal(knight, pieceA2)
}

func (s *StateSuite) TestValidMoves() {
	king := NewPiece(King(s.T()), s.state.CurrentPlayer())
	king.PlaceOn(s.state.Board(), boardtest.NewSquare("A1"))

	moves := s.state.ValidMoves()

	movesMatch(s.T(), moves, movesMatcher(king, "A2", "B1"))
}

func (s *StateSuite) TestValidMovesWithValidator() {
	king := NewPiece(King(s.T()), s.state.CurrentPlayer())
	king.PlaceOn(s.state.Board(), boardtest.NewSquare("A1"))

	s.state.AddStateValidator(func(s *State, m *Move) bool {
		return m.To != boardtest.NewSquare("A2")
	})

	moves := s.state.ValidMoves()

	movesMatch(s.T(), moves, movesMatcher(king, "B1"))
}

func TestStateSuite(t *testing.T) {
	suite.Run(t, new(StateSuite))
}

func trueStateValidator(*State, *Move) bool  { return true }
func falseStateValidator(*State, *Move) bool { return false }

func TestChainStateValidator(t *testing.T) {
	tests := []struct {
		name       string
		validators []StateValidator
		expected   bool
	}{
		{
			name:       "Empty",
			validators: []StateValidator{},
			expected:   true,
		},
		{
			name:       "OneTrue",
			validators: []StateValidator{trueStateValidator},
			expected:   true,
		},
		{
			name:       "OneFalse",
			validators: []StateValidator{falseStateValidator},
			expected:   false,
		},
		{
			name:       "OneFalseOneTrue",
			validators: []StateValidator{falseStateValidator, trueStateValidator},
			expected:   false,
		},
		{
			name:       "TwoFalse",
			validators: []StateValidator{falseStateValidator, falseStateValidator},
			expected:   false,
		},
		{
			name:       "TwoTrue",
			validators: []StateValidator{trueStateValidator, trueStateValidator},
			expected:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validators := chainStateValidators(tt.validators)
			isValid := validators.Validate(nil, nil)
			assert.Equal(t, tt.expected, isValid)
		})
	}
}
