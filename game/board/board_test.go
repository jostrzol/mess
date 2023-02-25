package board

import (
	"fmt"
	"testing"

	"github.com/jostrzol/mess/game/piece"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestNew(t *testing.T) {
	tests := []struct {
		x int
		y int
	}{
		{1, 1},
		{5, 5},
		{100, 100},
		{5, 2},
		{2, 5},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%dx%d", tt.x, tt.y)
		t.Run(name, func(t *testing.T) {
			board, err := NewBoard(tt.x, tt.y)
			assert.NoError(t, err)
			assert.Len(t, board, tt.y)
			for _, row := range board {
				assert.Len(t, row, tt.x)
				for _, value := range row {
					assert.Nil(t, value)
				}
			}
		})
	}
}

func TestNewNotPositive(t *testing.T) {
	tests := []struct {
		x int
		y int
	}{
		{0, 2},
		{2, 0},
		{0, 0},
		{-1, 2},
		{2, -1},
		{-1, -1},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%dx%d", tt.x, tt.y)
		t.Run(name, func(t *testing.T) {
			_, err := NewBoard(tt.x, tt.y)
			assert.Error(t, err)
		})
	}
}

type BoardSuite struct {
	suite.Suite
	board Board
}

func (s *BoardSuite) SetupTest() {
	board, err := NewBoard(6, 8)
	s.NoError(err)
	s.board = board
}

func (s *BoardSuite) TestSize() {
	x, y := s.board.Size()
	s.Equal(6, x)
	s.Equal(8, y)
}

func (s *BoardSuite) TestPieceOn() {
	tests := []string{"A1", "B1", "A2", "B2", "F8"}
	for _, squareStr := range tests {
		s.Run(squareStr, func() {
			square, _ := ParseSquare(squareStr)
			result, err := s.board.PieceOn(square)
			s.NoError(err)
			s.Nil(result.Piece)
			s.Equal(square, result.Square)
		})
	}
}

func (s *BoardSuite) TestPieceOnOutOfBound() {
	tests := []string{"G8", "F9", "G9"}
	for _, squareStr := range tests {
		s.Run(squareStr, func() {
			square, _ := ParseSquare(squareStr)
			_, err := s.board.PieceOn(square)
			s.Error(err)
		})
	}
}

func (s *BoardSuite) TestPlace() {
	rook := piece.Rook(s.T())
	square, _ := ParseSquare("B3")

	err := s.board.Place(rook, square)
	s.NoError(err)

	result, _ := s.board.PieceOn(square)
	s.Equal(rook, result.Piece)
}

func (s *BoardSuite) TestPlaceReplace() {
	rook := piece.Rook(s.T())
	knight := piece.Knight(s.T())
	square, _ := ParseSquare("B3")

	err := s.board.Place(rook, square)
	err = s.board.Place(knight, square)
	s.NoError(err)

	result, _ := s.board.PieceOn(square)
	s.Equal(knight, result.Piece)
}

func (s *BoardSuite) TestAllPieces() {
	rook := piece.Rook(s.T())
	knight := piece.Knight(s.T())
	square1, _ := ParseSquare("B3")
	square2, _ := ParseSquare("D6")
	s.board.Place(rook, square1)
	s.board.Place(knight, square2)

	pieces := s.board.AllPieces()

	s.ElementsMatch(pieces, []PieceOnSquare{
		{rook, square1},
		{knight, square2},
	})
}

func TestBoardSuite(t *testing.T) {
	suite.Run(t, new(BoardSuite))
}
