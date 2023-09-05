package mess_test

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/event"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/mess/messtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestNewPlayers(t *testing.T) {
	board := event.NewSubject()

	players := mess.NewPlayers(board)
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
	players map[color.Color]*mess.Player
	white   *mess.Player
	black   *mess.Player
}

func (s *PlayerSuiteMockedBoard) SetupTest() {
	s.board = event.NewSubject()
	s.players = mess.NewPlayers(s.board)
	s.white = s.players[color.White]
	s.black = s.players[color.Black]
}

func (s *PlayerSuiteMockedBoard) TestString() {
	s.Equal(s.white.String(), color.White.String())
}

func (s *PlayerSuiteMockedBoard) TestCapturesEmpty() {
	s.Empty(s.white.Captures())
}

func (s *PlayerSuiteMockedBoard) TestCapturesCapture() {
	knight := mess.NewPiece(Knight(s.T()), s.black)
	s.board.Notify(mess.PieceCaptured{
		Piece:      knight,
		CapturedBy: s.white,
	})

	s.Len(s.white.Captures(), 1)
	s.Contains(s.white.Captures(), knight)

	s.Empty(s.black.Captures())
}

func (s *PlayerSuiteMockedBoard) TestCapturesRelease() {
	knight := mess.NewPiece(Knight(s.T()), s.black)
	s.board.Notify(mess.PieceCaptured{
		Piece:      knight,
		CapturedBy: s.white,
	})

	s.board.Notify(mess.PiecePlaced{
		Piece: knight,
	})

	s.Empty(s.white.Captures())
	s.Empty(s.black.Captures())
}

func (s *PlayerSuiteMockedBoard) TestPiecesEmpty() {
	s.Empty(s.white.Pieces())
}

func (s *PlayerSuiteMockedBoard) TestPiecePlaced() {
	knight := mess.NewPiece(Knight(s.T()), s.white)
	s.board.Notify(mess.PiecePlaced{Piece: knight})

	s.Len(s.white.Pieces(), 1)
	s.Contains(s.white.Pieces(), knight)

	s.Empty(s.black.Pieces())
}

func (s *PlayerSuiteMockedBoard) TestPieceRemoved() {
	knight := mess.NewPiece(Knight(s.T()), s.white)
	s.board.Notify(mess.PiecePlaced{Piece: knight})
	s.board.Notify(mess.PieceRemoved{Piece: knight})

	s.Empty(s.white.Pieces())
	s.Empty(s.black.Pieces())
}

func TestPlayerSourceMockedBoard(t *testing.T) {
	suite.Run(t, new(PlayerSuiteMockedBoard))
}

type PlayerSuiteRealBoard struct {
	suite.Suite
	board   *mess.PieceBoard
	players map[color.Color]*mess.Player
	white   *mess.Player
	black   *mess.Player
}

func (s *PlayerSuiteRealBoard) SetupTest() {
	board, err := mess.NewPieceBoard(3, 3)
	s.NoError(err)

	s.board = board
	s.players = mess.NewPlayers(s.board)
	s.white = s.players[color.White]
	s.black = s.players[color.Black]
}

func (s *PlayerSuiteRealBoard) TestMovesNone() {
	s.Empty(s.white.Moves())
}

func (s *PlayerSuiteRealBoard) TestMovesOnePiece() {
	king := mess.NewPiece(King(s.T()), s.white)
	s.board.Place(king, boardtest.NewSquare("A1"))

	moves := s.white.Moves()
	messtest.MovesMatch(s.T(), moves, messtest.MovesMatcher(king, "A2", "B1"))
}

func (s *PlayerSuiteRealBoard) TestMovesOnePieceOneEnemy() {
	kingW := mess.NewPiece(King(s.T()), s.white)
	s.board.Place(kingW, boardtest.NewSquare("A1"))

	kingB := mess.NewPiece(King(s.T()), s.black)
	s.board.Place(kingB, boardtest.NewSquare("A3"))

	moves := s.white.Moves()
	messtest.MovesMatch(s.T(), moves, messtest.MovesMatcher(kingW, "A2", "B1"))
}

func (s *PlayerSuiteRealBoard) TestMovesTwoPieces() {
	kingW1 := mess.NewPiece(King(s.T()), s.white)
	s.board.Place(kingW1, boardtest.NewSquare("A1"))

	kingW2 := mess.NewPiece(King(s.T()), s.white)
	s.board.Place(kingW2, boardtest.NewSquare("A3"))

	moves := s.white.Moves()
	messtest.MovesMatch(s.T(), moves,
		messtest.MovesMatcher(kingW1, "A2", "B1"),
		messtest.MovesMatcher(kingW2, "A2", "B3"),
	)
}

func TestPlayerSourceRealBoard(t *testing.T) {
	suite.Run(t, new(PlayerSuiteRealBoard))
}
