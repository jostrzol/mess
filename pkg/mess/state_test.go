package mess_test

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/mess/messtest"
	"github.com/stretchr/testify/suite"
)

type StateSuite struct {
	suite.Suite
	state *mess.State
}

func (s *StateSuite) SetupTest() {
	board, err := mess.NewPieceBoard(8, 8)
	s.NoError(err)
	s.state = mess.NewState(board)
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
	var rook, knight, king mess.Piece
	rookSquare := boardtest.NewSquare("A1")
	knightSquare := boardtest.NewSquare("A2")
	kingSquare := boardtest.NewSquare("A3")
	emptySquare1 := boardtest.NewSquare("B1")
	emptySquare2 := boardtest.NewSquare("B2")

	type move struct {
		*mess.Piece
		board.Square
	}
	tests := []struct {
		name     string
		moves    []move
		expected []mess.RecordedMove
	}{
		{
			name:  "One",
			moves: []move{{&rook, emptySquare1}},
			expected: []mess.RecordedMove{
				{
					Piece:      &rook,
					From:       rookSquare,
					To:         emptySquare1,
					Captures:   map[*mess.Piece]struct{}{},
					TurnNumber: 1,
				},
			},
		},
		{
			name:  "Two",
			moves: []move{{&rook, emptySquare1}, {&rook, emptySquare2}},
			expected: []mess.RecordedMove{
				{
					Piece:      &rook,
					From:       rookSquare,
					To:         emptySquare1,
					Captures:   map[*mess.Piece]struct{}{},
					TurnNumber: 1,
				},
				{
					Piece:      &rook,
					From:       emptySquare1,
					To:         emptySquare2,
					Captures:   map[*mess.Piece]struct{}{},
					TurnNumber: 1,
				},
			},
		},
		{
			name:  "Capture",
			moves: []move{{&rook, knightSquare}},
			expected: []mess.RecordedMove{
				{
					Piece: &rook,
					From:  rookSquare,
					To:    knightSquare,
					Captures: map[*mess.Piece]struct{}{
						&knight: {},
					},
					TurnNumber: 1,
				},
			},
		},
		{
			name:  "DoubleCapture",
			moves: []move{{&rook, knightSquare}, {&rook, kingSquare}},
			expected: []mess.RecordedMove{
				{
					Piece: &rook,
					From:  rookSquare,
					To:    knightSquare,
					Captures: map[*mess.Piece]struct{}{
						&knight: {},
					},
					TurnNumber: 1,
				},
				{
					Piece: &rook,
					From:  knightSquare,
					To:    kingSquare,
					Captures: map[*mess.Piece]struct{}{
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
			rook = *mess.NewPiece(Rook(s.T()), s.state.CurrentPlayer())
			rook.PlaceOn(s.state.Board(), rookSquare)
			knight = *mess.NewPiece(Knight(s.T()), s.state.CurrentPlayer())
			knight.PlaceOn(s.state.Board(), knightSquare)
			king = *mess.NewPiece(Knight(s.T()), s.state.CurrentPlayer())
			king.PlaceOn(s.state.Board(), kingSquare)

			for _, move := range tt.moves {
				err := move.Piece.MoveTo(move.Square)
				s.NoError(err)
			}

			s.Equal(s.state.Record(), tt.expected)
		})
	}
}

func (s *StateSuite) TestRecordTurnNumber() {
	rook := *mess.NewPiece(Rook(s.T()), s.state.CurrentPlayer())
	rook.PlaceOn(s.state.Board(), boardtest.NewSquare("A1"))

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

	s.Equal(s.state.Record()[0].TurnNumber, 1)
	s.Equal(s.state.Record()[1].TurnNumber, 2)
	s.Equal(s.state.Record()[2].TurnNumber, 2)
	s.Equal(s.state.Record()[3].TurnNumber, 3)
}

func (s *StateSuite) TestUndoNothing() {
	s.state.UndoTurn()
}

func (s *StateSuite) TestUndo() {
	rook := mess.NewPiece(Rook(s.T()), s.state.CurrentPlayer())
	a1 := boardtest.NewSquare("A1")
	rook.PlaceOn(s.state.Board(), a1)

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
	rook := mess.NewPiece(Rook(s.T()), s.state.CurrentPlayer())
	a1 := boardtest.NewSquare("A1")
	rook.PlaceOn(s.state.Board(), a1)

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
	rook := mess.NewPiece(Rook(s.T()), s.state.CurrentPlayer())
	a1 := boardtest.NewSquare("A1")
	rook.PlaceOn(s.state.Board(), a1)

	knight := mess.NewPiece(Knight(s.T()), s.state.CurrentPlayer())
	a2 := boardtest.NewSquare("A2")
	knight.PlaceOn(s.state.Board(), a2)

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
	king := mess.NewPiece(King(s.T()), s.state.CurrentPlayer())
	king.PlaceOn(s.state.Board(), boardtest.NewSquare("A1"))

	moves := s.state.ValidMoves()

	messtest.MovesMatch(s.T(), moves, messtest.MovesMatcher(king, "A2", "B1"))
}

func (s *StateSuite) TestValidMovesWithValidator() {
	king := mess.NewPiece(King(s.T()), s.state.CurrentPlayer())
	king.PlaceOn(s.state.Board(), boardtest.NewSquare("A1"))

	s.state.AddStateValidator(func(s *mess.State, m *mess.Move) bool {
		return m.To != boardtest.NewSquare("A2")
	})

	moves := s.state.ValidMoves()

	messtest.MovesMatch(s.T(), moves, messtest.MovesMatcher(king, "B1"))
}

func TestStateSuite(t *testing.T) {
	suite.Run(t, new(StateSuite))
}
