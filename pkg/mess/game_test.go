package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
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
	rookSquare := *boardtest.NewSquare("A1")
	knightSquare := *boardtest.NewSquare("A2")
	kingSquare := *boardtest.NewSquare("A3")
	emptySquare1 := *boardtest.NewSquare("B1")
	emptySquare2 := *boardtest.NewSquare("B2")

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
					Move: Move{
						Piece: &rook,
						From:  rookSquare,
						To:    emptySquare1,
					},
					Captures: map[*Piece]struct{}{},
				},
			},
		},
		{
			name:  "Two",
			moves: []move{{&rook, emptySquare1}, {&rook, emptySquare2}},
			expected: []RecordedMove{
				{
					Move: Move{
						Piece: &rook,
						From:  rookSquare,
						To:    emptySquare1,
					},
					Captures: map[*Piece]struct{}{},
				},
				{
					Move: Move{
						Piece: &rook,
						From:  emptySquare1,
						To:    emptySquare2,
					},
					Captures: map[*Piece]struct{}{},
				},
			},
		},
		{
			name:  "Capture",
			moves: []move{{&rook, knightSquare}},
			expected: []RecordedMove{
				{
					Move: Move{
						Piece: &rook,
						From:  rookSquare,
						To:    knightSquare,
					},
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
					Move: Move{
						Piece: &rook,
						From:  rookSquare,
						To:    knightSquare,
					},
					Captures: map[*Piece]struct{}{
						&knight: {},
					},
				},
				{
					Move: Move{
						Piece: &rook,
						From:  knightSquare,
						To:    kingSquare,
					},
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
			rook.PlaceOn(s.game.board, &rookSquare)
			knight = *NewPiece(Knight(s.T()), s.game.currentPlayer)
			knight.PlaceOn(s.game.board, &knightSquare)
			king = *NewPiece(Knight(s.T()), s.game.currentPlayer)
			king.PlaceOn(s.game.board, &kingSquare)

			for _, move := range tt.moves {
				err := move.Piece.MoveTo(&move.Square)
				s.NoError(err)
			}

			s.Equal(s.game.record, tt.expected)
		})
	}
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}
