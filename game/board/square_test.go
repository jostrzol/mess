package board_test

import (
	"fmt"
	"testing"

	"github.com/jostrzol/mess/game/board"
	"github.com/stretchr/testify/assert"
)

func TestNewSquare(t *testing.T) {
	tests := []struct {
		str  string
		file int
		rank int
	}{
		{"A1", 1, 1},
		{"B1", 2, 1},
		{"A2", 1, 2},
		{"B2", 2, 2},
		{"H6", 8, 6},
		{"Z9", 26, 9},
		{"a1", 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			square, err := board.NewSquare(tt.str)
			assert.NoError(t, err)
			assert.Equal(t, square.File, tt.file)
			assert.Equal(t, square.Rank, tt.rank)
		})
	}
}

func TestNewSquareMalformed(t *testing.T) {
	tests := []string{"A10", "Ż1", "A", "1", "AB1", "hello", "-", " ", "", " A1", "A1 ", " A1 "}
	for _, str := range tests {
		t.Run(str, func(t *testing.T) {
			_, err := board.NewSquare(str)
			assert.Error(t, err)
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"A1", "A1"},
		{"B1", "B1"},
		{"A2", "A2"},
		{"B2", "B2"},
		{"H6", "H6"},
		{"Z9", "Z9"},
		{"a1", "A1"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			square, err := board.NewSquare(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, square.String())
		})
	}
}

func TestOffset(t *testing.T) {
	tests := []struct {
		input    string
		x        int
		y        int
		expected string
	}{
		{"B2", 0, 0, "B2"},
		{"B2", 1, 0, "C2"},
		{"B2", 1, 1, "C3"},
		{"B2", 0, 1, "B3"},
		{"B2", -1, 1, "A3"},
		{"B2", -1, 0, "A2"},
		{"B2", -1, -1, "A1"},
		{"B2", 0, -1, "B1"},
		{"B2", 1, -1, "C1"},
		{"G5", -3, 4, "D9"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s->%s", tt.input, tt.expected), func(t *testing.T) {
			square, err := board.NewSquare(tt.input)
			assert.NoError(t, err)
			expectedSquare, err := board.NewSquare(tt.expected)
			assert.NoError(t, err)

			newSquare := square.Offset(board.Offset{X: tt.x, Y: tt.y})
			assert.Equal(t, expectedSquare, newSquare)
		})
	}
}
