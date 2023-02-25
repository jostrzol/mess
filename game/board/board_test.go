package board

import (
	"fmt"
	"testing"

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
	s.Equal(x, 6)
	s.Equal(y, 8)
}

func TestBoardSuite(t *testing.T) {
	suite.Run(t, new(BoardSuite))
}
