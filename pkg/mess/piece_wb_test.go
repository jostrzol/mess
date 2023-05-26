package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/stretchr/testify/assert"
)

func staticMoveGenerator(t *testing.T, strings ...string) MoveGenerator {
	t.Helper()
	return func(piece *Piece) ([]brd.Square, MoveAction) {
		destinations := make([]brd.Square, 0, len(strings))
		for _, squareStr := range strings {
			square, err := brd.NewSquare(squareStr)
			assert.NoError(t, err)
			destinations = append(destinations, square)
		}
		return destinations, nil
	}
}

func TestChainMoveGenerators(t *testing.T) {
	tests := []struct {
		name       string
		generators []MoveGenerator
		expected   []string
	}{
		{
			name:       "Empty",
			generators: []MoveGenerator{},
			expected:   []string{},
		},
		{
			name: "One",
			generators: []MoveGenerator{
				staticMoveGenerator(t, "A1"),
			},
			expected: []string{"A1"},
		},
		{
			name: "Two",
			generators: []MoveGenerator{
				staticMoveGenerator(t, "A1"),
				staticMoveGenerator(t, "B1"),
			},
			expected: []string{"A1", "B1"},
		},
		{
			name: "TwoOverlapping",
			generators: []MoveGenerator{
				staticMoveGenerator(t, "A1"),
				staticMoveGenerator(t, "A1"),
			},
			expected: []string{"A1"},
		},
		{
			name: "TwoOverlapping",
			generators: []MoveGenerator{
				staticMoveGenerator(t, "A1", "B2"),
				staticMoveGenerator(t, "C5"),
				staticMoveGenerator(t, "B2", "D4", "C5"),
			},
			expected: []string{"A1", "B2", "C5", "D4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generators := chainMoveGenerators(tt.generators)
			destinations := make([]board.Square, 0)
			for _, generated := range generators.Generate(nil) {
				destinations = append(destinations, generated.destination)
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
