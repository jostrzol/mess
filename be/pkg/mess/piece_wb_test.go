package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/jostrzol/mess/pkg/board/boardtest"
	"github.com/stretchr/testify/assert"
)

func staticMoveGenerator(t *testing.T, strings ...string) Motion {
	t.Helper()
	return Motion{
		Name: "test_generator",
		MoveGenerator: func(piece *Piece) []brd.Square {
			destinations := make([]brd.Square, 0, len(strings))
			for _, squareStr := range strings {
				square, err := brd.NewSquare(squareStr)
				assert.NoError(t, err)
				destinations = append(destinations, square)
			}
			return destinations
		},
	}
}

func TestChainMoveGenerators(t *testing.T) {
	tests := []struct {
		name       string
		generators []Motion
		expected   []string
	}{
		{
			name:       "Empty",
			generators: []Motion{},
			expected:   []string{},
		},
		{
			name: "One",
			generators: []Motion{
				staticMoveGenerator(t, "A1"),
			},
			expected: []string{"A1"},
		},
		{
			name: "Two",
			generators: []Motion{
				staticMoveGenerator(t, "A1"),
				staticMoveGenerator(t, "B1"),
			},
			expected: []string{"A1", "B1"},
		},
		{
			name: "TwoOverlapping",
			generators: []Motion{
				staticMoveGenerator(t, "A1"),
				staticMoveGenerator(t, "A1"),
			},
			expected: []string{"A1"},
		},
		{
			name: "TwoOverlapping",
			generators: []Motion{
				staticMoveGenerator(t, "A1", "B2"),
				staticMoveGenerator(t, "C5"),
				staticMoveGenerator(t, "B2", "D4", "C5"),
			},
			expected: []string{"A1", "B2", "C5", "D4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generators := chainMotions(tt.generators)
			destinations := make([]board.Square, 0)
			for _, generated := range generators.Generate(&Piece{square: boardtest.NewSquare("A1")}) {
				destinations = append(destinations, generated.To)
			}
			assertSquaresMatch(t, destinations, tt.expected...)
		})
	}
}

func assertSquaresMatch(t *testing.T, actual []brd.Square, expected ...string) {
	assert.Len(t, actual, len(expected))
	for _, str := range expected {
		square, err := brd.NewSquare(str)
		assert.NoError(t, err)
		assert.Containsf(t, actual, square, "%v doesnt contain square %v", actual, square)
	}
}
