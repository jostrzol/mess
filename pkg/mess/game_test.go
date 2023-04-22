package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GameSuite struct {
	suite.Suite
	game *State
}

func (s *GameSuite) SetupTest() {
	board, err := NewPieceBoard(8, 8)
	s.NoError(err)
	s.game = NewState(board)
}

func (s *GameSuite) TestGetPlayer() {
	for _, color := range color.ColorValues() {
		s.Run(color.String(), func() {
			player := s.game.Player(color)
			s.Equal(player.Color(), color)
		})
	}
}

func (s *GameSuite) TestGetPlayerNotFound() {
	s.Panics(func() {
		s.game.Player(color.Color(-1))
	})
}

func (s *GameSuite) TestEndTurn() {
	firstTurnPlayer := s.game.CurrentPlayer()
	s.Equal(s.game.Player(color.White), firstTurnPlayer)

	s.game.EndTurn()

	secondTurnPlayer := s.game.CurrentPlayer()
	s.Equal(s.game.Player(color.Black), secondTurnPlayer)
}

func (s *GameSuite) TestRecordMove() {
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
					Piece:    &rook,
					From:     rookSquare,
					To:       emptySquare1,
					Captures: map[*Piece]struct{}{},
				},
			},
		},
		{
			name:  "Two",
			moves: []move{{&rook, emptySquare1}, {&rook, emptySquare2}},
			expected: []RecordedMove{
				{
					Piece:    &rook,
					From:     rookSquare,
					To:       emptySquare1,
					Captures: map[*Piece]struct{}{},
				},
				{
					Piece:    &rook,
					From:     emptySquare1,
					To:       emptySquare2,
					Captures: map[*Piece]struct{}{},
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
				},
				{
					Piece: &rook,
					From:  knightSquare,
					To:    kingSquare,
					Captures: map[*Piece]struct{}{
						&king: {},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		s.SetupTest()
		s.Run(tt.name, func() {
			rook = *NewPiece(Rook(s.T()), s.game.currentPlayer)
			rook.PlaceOn(s.game.board, rookSquare)
			knight = *NewPiece(Knight(s.T()), s.game.currentPlayer)
			knight.PlaceOn(s.game.board, knightSquare)
			king = *NewPiece(Knight(s.T()), s.game.currentPlayer)
			king.PlaceOn(s.game.board, kingSquare)

			for _, move := range tt.moves {
				err := move.Piece.MoveTo(move.Square)
				s.NoError(err)
			}

			s.Equal(s.game.record, tt.expected)
		})
	}
}

func (s *GameSuite) TestUndoNothing() {
	s.Nil(s.game.Undo())
}

func (s *GameSuite) TestUndo() {
	rook := NewPiece(Rook(s.T()), s.game.currentPlayer)
	a1 := boardtest.NewSquare("A1")
	rook.PlaceOn(s.game.board, a1)

	a2 := boardtest.NewSquare("A2")
	err := rook.MoveTo(a2)
	s.NoError(err)

	move := s.game.Undo()

	expected := &RecordedMove{
		Piece:    rook,
		From:     a1,
		To:       a2,
		Captures: map[*Piece]struct{}{},
	}
	s.Equal(expected, move)

	pieceA1, err := s.game.Board().At(a1)
	s.NoError(err)
	s.Equal(rook, pieceA1)

	pieceA2, err := s.game.Board().At(a2)
	s.NoError(err)
	s.Nil(pieceA2)
}

func (s *GameSuite) TestUndoDoubleMove() {
	rook := NewPiece(Rook(s.T()), s.game.currentPlayer)
	a1 := boardtest.NewSquare("A1")
	rook.PlaceOn(s.game.board, a1)

	a2 := boardtest.NewSquare("A2")
	err := rook.MoveTo(a2)
	s.NoError(err)

	a3 := boardtest.NewSquare("A3")
	err = rook.MoveTo(a3)
	s.NoError(err)

	move := s.game.Undo()

	expected := &RecordedMove{
		Piece:    rook,
		From:     a2,
		To:       a3,
		Captures: map[*Piece]struct{}{},
	}
	s.Equal(expected, move)

	pieceA1, err := s.game.Board().At(a1)
	s.NoError(err)
	s.Nil(pieceA1)

	pieceA2, err := s.game.Board().At(a2)
	s.NoError(err)
	s.Equal(rook, pieceA2)

	pieceA3, err := s.game.Board().At(a3)
	s.NoError(err)
	s.Nil(pieceA3)
}

func (s *GameSuite) TestUndoCapture() {
	rook := NewPiece(Rook(s.T()), s.game.currentPlayer)
	a1 := boardtest.NewSquare("A1")
	rook.PlaceOn(s.game.board, a1)

	knight := NewPiece(Knight(s.T()), s.game.currentPlayer)
	a2 := boardtest.NewSquare("A2")
	knight.PlaceOn(s.game.board, a2)

	err := rook.MoveTo(a2)
	s.NoError(err)

	move := s.game.Undo()

	expected := &RecordedMove{
		Piece:    rook,
		From:     a1,
		To:       a2,
		Captures: map[*Piece]struct{}{knight: {}},
	}
	s.Equal(expected, move)

	pieceA1, err := s.game.Board().At(a1)
	s.NoError(err)
	s.Equal(rook, pieceA1)

	pieceA2, err := s.game.Board().At(a2)
	s.NoError(err)
	s.Equal(knight, pieceA2)
}

func (s *GameSuite) TestValidMoves() {
	king := NewPiece(King(s.T()), s.game.CurrentPlayer())
	king.PlaceOn(s.game.Board(), boardtest.NewSquare("A1"))

	moves := s.game.ValidMoves()

	movesMatch(s.T(), moves, movesMatcher(king, "A2", "B1"))
}

func (s *GameSuite) TestValidMovesWithValidator() {
	king := NewPiece(King(s.T()), s.game.CurrentPlayer())
	king.PlaceOn(s.game.Board(), boardtest.NewSquare("A1"))

	s.game.AddStateValidator(func(s *State, m *Move) bool {
		return m.To != boardtest.NewSquare("A2")
	})

	moves := s.game.ValidMoves()

	movesMatch(s.T(), moves, movesMatcher(king, "B1"))
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
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
