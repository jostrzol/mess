package piece_test

import (
	"testing"

	brd "github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/board/boardtest"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/piece/piecetest"
	"github.com/stretchr/testify/suite"
)

type PieceSuite struct {
	suite.Suite
	board piece.Board
}

func (s *PieceSuite) SetupTest() {
	board, err := brd.NewBoard[*piece.Piece](4, 4)
	s.NoError(err)

	s.board = board
}

func (s *PieceSuite) TestGenerateMotionsRook() {
	rook := piecetest.Noones(piecetest.Rook(s.T()))
	rook.PlaceOn(s.board, boardtest.NewSquare("B2"))

	motions := rook.GenerateMotions()
	s.ElementsMatch(motions, []brd.Square{
		boardtest.NewSquare("B1"),
		boardtest.NewSquare("B3"),
		boardtest.NewSquare("B4"),
		boardtest.NewSquare("A2"),
		boardtest.NewSquare("C2"),
		boardtest.NewSquare("D2"),
	})
}

func (s *PieceSuite) TestGenerateMotionsKnight() {
	knight := piecetest.Noones(piecetest.Knight(s.T()))
	knight.PlaceOn(s.board, boardtest.NewSquare("B2"))

	motions := knight.GenerateMotions()
	s.ElementsMatch(motions, []brd.Square{
		boardtest.NewSquare("A4"),
		boardtest.NewSquare("C4"),
		boardtest.NewSquare("D1"),
		boardtest.NewSquare("D3"),
	})
}

func TestPieceSuite(t *testing.T) {
	suite.Run(t, new(PieceSuite))
}
