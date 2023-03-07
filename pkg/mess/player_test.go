package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
	"github.com/jostrzol/mess/pkg/genassert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestNewPlayers(t *testing.T) {
	board := event.NewSubject()

	players := NewPlayers(board)
	assert.Len(t, players, 2)

	for _, color := range color.ColorValues() {
		assert.Contains(t, players, color)
		player := players[color]
		assert.Equal(t, player.Color(), color)
		assert.True(t, board.IsObserving(player))
	}
}

type PlayerSuiteMockedBoard struct {
	suite.Suite
	board   event.Subject
	players map[color.Color]*Player
	white   *Player
	black   *Player
}

func (s *PlayerSuiteMockedBoard) SetupTest() {
	s.board = event.NewSubject()
	s.players = NewPlayers(s.board)
	s.white = s.players[color.White]
	s.black = s.players[color.Black]
}

func (s *PlayerSuiteMockedBoard) TestString() {
	s.Equal(s.white.String(), color.White.String())
}

func (s *PlayerSuiteMockedBoard) TestPrisonersEmpty() {
	genassert.Empty(s.T(), s.white.Prisoners())
}

func (s *PlayerSuiteMockedBoard) TestPrisonersCapture() {
	knight := NewPiece(Knight(s.T()), s.black)
	s.board.Notify(PieceCaptured{
		Piece:      knight,
		CapturedBy: s.white,
	})

	genassert.Len(s.T(), s.white.Prisoners(), 1)
	genassert.Contains(s.T(), s.white.Prisoners(), knight)

	genassert.Empty(s.T(), s.black.Prisoners())
}

func (s *PlayerSuiteMockedBoard) TestPiecesEmpty() {
	genassert.Empty(s.T(), s.white.Prisoners())
}

func (s *PlayerSuiteMockedBoard) TestPiecePlaced() {
	knight := NewPiece(Knight(s.T()), s.white)
	s.board.Notify(PiecePlaced{Piece: knight})

	genassert.Len(s.T(), s.white.Pieces(), 1)
	genassert.Contains(s.T(), s.white.Pieces(), knight)

	genassert.Empty(s.T(), s.black.Pieces())
}

func (s *PlayerSuiteMockedBoard) TestPieceRemoved() {
	knight := NewPiece(Knight(s.T()), s.white)
	s.board.Notify(PiecePlaced{Piece: knight})
	s.board.Notify(PieceRemoved{Piece: knight})

	genassert.Empty(s.T(), s.white.Pieces())
	genassert.Empty(s.T(), s.black.Pieces())
}

func TestPlayerSourceMockedBoard(t *testing.T) {
	suite.Run(t, new(PlayerSuiteMockedBoard))
}

type PlayerSuiteRealBoard struct {
	suite.Suite
	board   *PieceBoard
	players map[color.Color]*Player
	white   *Player
	black   *Player
}

func (s *PlayerSuiteRealBoard) SetupTest() {
	board, err := NewPieceBoard(3, 3)
	s.NoError(err)

	s.board = board
	s.players = NewPlayers(s.board)
	s.white = s.players[color.White]
	s.black = s.players[color.Black]
}

func (s *PlayerSuiteRealBoard) TestGenerateMotionsNone() {
	s.Empty(s.white.GenerateMotions())
}

func (s *PlayerSuiteRealBoard) TestGenerateMotionsOnePiece() {
	king := NewPiece(King(s.T()), s.white)
	s.board.Place(king, boardtest.NewSquare("A1"))

	motions := s.white.GenerateMotions()
	s.ElementsMatch(motions, []Motion{
		{Piece: king,
			From: *boardtest.NewSquare("A1"),
			To:   *boardtest.NewSquare("A2")},
		{Piece: king,
			From: *boardtest.NewSquare("A1"),
			To:   *boardtest.NewSquare("B1")},
	})
}

func (s *PlayerSuiteRealBoard) TestGenerateMotionsOnePieceOneEnemy() {
	kingW := NewPiece(King(s.T()), s.white)
	s.board.Place(kingW, boardtest.NewSquare("A1"))

	kingB := NewPiece(King(s.T()), s.black)
	s.board.Place(kingB, boardtest.NewSquare("A3"))

	motions := s.white.GenerateMotions()
	s.ElementsMatch(motions, []Motion{
		{Piece: kingW,
			From: *boardtest.NewSquare("A1"),
			To:   *boardtest.NewSquare("A2")},
		{Piece: kingW,
			From: *boardtest.NewSquare("A1"),
			To:   *boardtest.NewSquare("B1")},
	})
}

func (s *PlayerSuiteRealBoard) TestGenerateMotionsTwoPieces() {
	kingW1 := NewPiece(King(s.T()), s.white)
	s.board.Place(kingW1, boardtest.NewSquare("A1"))

	kingW2 := NewPiece(King(s.T()), s.white)
	s.board.Place(kingW2, boardtest.NewSquare("A3"))

	motions := s.white.GenerateMotions()
	s.ElementsMatch(motions, []Motion{
		{Piece: kingW1,
			From: *boardtest.NewSquare("A1"),
			To:   *boardtest.NewSquare("A2")},
		{Piece: kingW1,
			From: *boardtest.NewSquare("A1"),
			To:   *boardtest.NewSquare("B1")},
		{Piece: kingW2,
			From: *boardtest.NewSquare("A3"),
			To:   *boardtest.NewSquare("A2")},
		{Piece: kingW2,
			From: *boardtest.NewSquare("A3"),
			To:   *boardtest.NewSquare("B3")},
	})
}

func TestPlayerSourceRealBoard(t *testing.T) {
	suite.Run(t, new(PlayerSuiteRealBoard))
}
