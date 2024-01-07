package board_test

import (
	"fmt"
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestNewBoard(t *testing.T) {
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
			board, err := board.NewBoard[int](tt.x, tt.y)
			assert.NoError(t, err)
			assert.Len(t, board, tt.y)
			for _, row := range board {
				assert.Len(t, row, tt.x)
				for _, value := range row {
					assert.Zero(t, value)
				}
			}
		})
	}
}

func TestNewBoardNotPositive(t *testing.T) {
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
			_, err := board.NewBoard[int](tt.x, tt.y)
			assert.Error(t, err)
		})
	}
}

type BoardSuite struct {
	suite.Suite
	board board.Board[int]
}

func (s *BoardSuite) SetupTest() {
	board, err := board.NewBoard[int](6, 8)
	s.NoError(err)
	s.board = board
}

func (s *BoardSuite) TestSize() {
	x, y := s.board.Size()
	s.Equal(6, x)
	s.Equal(8, y)
}

func (s *BoardSuite) TestAtEmpty() {
	tests := []string{"A1", "B1", "A2", "B2", "F8"}
	for _, squareStr := range tests {
		s.Run(squareStr, func() {
			item, err := s.board.At(boardtest.NewSquare(squareStr))
			s.NoError(err)
			s.Zero(item)
		})
	}
}

func (s *BoardSuite) TestAtOutOfBound() {
	tests := []string{"G8", "F9", "G9"}
	for _, squareStr := range tests {
		s.Run(squareStr, func() {
			_, err := s.board.At(boardtest.NewSquare(squareStr))
			s.Error(err)
		})
	}
}

func (s *BoardSuite) TestPlace() {
	square := boardtest.NewSquare("B3")

	old, err := s.board.Place(1, square)
	s.Zero(old)
	s.NoError(err)

	item, _ := s.board.At(square)
	s.Equal(1, item)
}

func (s *BoardSuite) TestPlaceReplace() {
	square := boardtest.NewSquare("B3")

	_, err := s.board.Place(1, square)
	s.NoError(err)
	old, err := s.board.Place(2, square)
	s.Equal(1, old)
	s.NoError(err)

	item, _ := s.board.At(square)
	s.Equal(2, item)
}

func (s *BoardSuite) TestAllPieces() {
	square1 := boardtest.NewSquare("B3")
	square2 := boardtest.NewSquare("D6")
	_, err := s.board.Place(1, square1)
	s.NoError(err)
	_, err = s.board.Place(2, square2)
	s.NoError(err)

	items := s.board.AllItems()

	s.ElementsMatch(items, []int{1, 2})
}

func TestBoardSuite(t *testing.T) {
	suite.Run(t, new(BoardSuite))
}
