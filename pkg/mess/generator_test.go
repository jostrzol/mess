package mess

import (
	"testing"

	"github.com/jostrzol/mess/pkg/board"
	brd "github.com/jostrzol/mess/pkg/board"
	"github.com/stretchr/testify/assert"
)

func staticMotionGenerator(t *testing.T, strings ...string) MotionGenerator {
	t.Helper()
	return FuncMotionGenerator(func(piece *Piece) []brd.Square {
		destinations := make([]brd.Square, 0, len(strings))
		for _, squareStr := range strings {
			square, err := brd.NewSquare(squareStr)
			assert.NoError(t, err)
			destinations = append(destinations, *square)
		}
		return destinations
	})
}

func offsetMotionGenerator(t *testing.T, offsets ...board.Offset) MotionGenerator {
	t.Helper()
	return FuncMotionGenerator(func(piece *Piece) []brd.Square {
		destinations := make([]brd.Square, 0, len(offsets))
		for _, offset := range offsets {
			square := piece.Square().Offset(offset)
			if piece.Board().Contains(square) {
				destinations = append(destinations, *square)
			}
		}
		return destinations
	})
}

func assertSquaresMatch(t *testing.T, actual []brd.Square, expected ...string) {
	assert.Len(t, actual, len(expected))
	for _, str := range expected {
		square, err := brd.NewSquare(str)
		assert.NoError(t, err)
		assert.Containsf(t, actual, *square, "%v doesnt contain square %v", actual, square)
	}
}

func TestGenerateMotionsZero(t *testing.T) {
	generators := MotionGenerators([]MotionGenerator{})
	motions := generators.GenerateMotions(nil)
	assert.Empty(t, motions)
}

func TestGenerateMotionsOne(t *testing.T) {
	generators := MotionGenerators([]MotionGenerator{
		staticMotionGenerator(t, "A1"),
	})
	motions := generators.GenerateMotions(nil)
	assertSquaresMatch(t, motions, "A1")
}

func TestGenerateMotionsTwo(t *testing.T) {
	generators := MotionGenerators([]MotionGenerator{
		staticMotionGenerator(t, "A1"),
		staticMotionGenerator(t, "B1"),
	})
	motions := generators.GenerateMotions(nil)
	assertSquaresMatch(t, motions, "A1", "B1")
}

func TestGenerateMotionsTwoOverlapping(t *testing.T) {
	generators := MotionGenerators([]MotionGenerator{
		staticMotionGenerator(t, "A1"),
		staticMotionGenerator(t, "A1"),
	})
	motions := generators.GenerateMotions(nil)
	assertSquaresMatch(t, motions, "A1")
}

func TestGenerateMotionsMany(t *testing.T) {
	generators := MotionGenerators([]MotionGenerator{
		staticMotionGenerator(t, "A1", "B2"),
		staticMotionGenerator(t, "C5"),
		staticMotionGenerator(t, "B2", "D4", "C5"),
	})
	motions := generators.GenerateMotions(nil)
	assertSquaresMatch(t, motions, "A1", "B2", "C5", "D4")
}
