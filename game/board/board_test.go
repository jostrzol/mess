package board_test

import (
	"fmt"
	"testing"

	"github.com/jostrzol/mess/game/board"
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
			square, _ := board.NewSquare(squareStr)
			item, err := s.board.At(&square)
			s.NoError(err)
			s.Zero(item)
		})
	}
}

func (s *BoardSuite) TestAtOutOfBound() {
	tests := []string{"G8", "F9", "G9"}
	for _, squareStr := range tests {
		s.Run(squareStr, func() {
			square, _ := board.NewSquare(squareStr)
			_, err := s.board.At(&square)
			s.Error(err)
		})
	}
}

func (s *BoardSuite) TestPlace() {
	square, _ := board.NewSquare("B3")

	err := s.board.Place(1, &square)
	s.NoError(err)

	item, _ := s.board.At(&square)
	s.Equal(1, item)
}

func (s *BoardSuite) TestPlaceReplace() {
	square, _ := board.NewSquare("B3")

	err := s.board.Place(1, &square)
	s.NoError(err)
	err = s.board.Place(2, &square)
	s.NoError(err)

	item, _ := s.board.At(&square)
	s.Equal(2, item)
}

func (s *BoardSuite) TestAllPieces() {
	square1, _ := board.NewSquare("B3")
	square2, _ := board.NewSquare("D6")
	s.board.Place(1, &square1)
	s.board.Place(2, &square2)

	items := s.board.AllItems()

	s.ElementsMatch(items, []int{1, 2})
}

func TestBoardSuite(t *testing.T) {
	suite.Run(t, new(BoardSuite))
}
