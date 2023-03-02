package piece_test

import (
	"testing"

	brd "github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
	"github.com/jostrzol/mess/game/piece/piecetest"
	"github.com/stretchr/testify/assert"
)

func assertSquaresMatch(t *testing.T, actual []brd.Square, expected ...string) {
	assert.Len(t, actual, len(expected))
	for _, str := range expected {
		square, err := brd.NewSquare(str)
		assert.NoError(t, err)
		assert.Containsf(t, actual, square, "%v doesnt contain square %v", actual, square)
	}
}

func TestGenerateMotionsZero(t *testing.T) {
	generators := piece.MotionGenerators([]piece.MotionGenerator{})
	motions := generators.GenerateMotions(nil)
	assert.Empty(t, motions)
}

func TestGenerateMotionsOne(t *testing.T) {
	generators := piece.MotionGenerators([]piece.MotionGenerator{
		piecetest.NewStaticMotionGenerator(t, "A1"),
	})
	motions := generators.GenerateMotions(nil)
	assertSquaresMatch(t, motions, "A1")
}

func TestGenerateMotionsTwo(t *testing.T) {
	generators := piece.MotionGenerators([]piece.MotionGenerator{
		piecetest.NewStaticMotionGenerator(t, "A1"),
		piecetest.NewStaticMotionGenerator(t, "B1"),
	})
	motions := generators.GenerateMotions(nil)
	assertSquaresMatch(t, motions, "A1", "B1")
}

func TestGenerateMotionsTwoOverlapping(t *testing.T) {
	generators := piece.MotionGenerators([]piece.MotionGenerator{
		piecetest.NewStaticMotionGenerator(t, "A1"),
		piecetest.NewStaticMotionGenerator(t, "A1"),
	})
	motions := generators.GenerateMotions(nil)
	assertSquaresMatch(t, motions, "A1")
}

func TestGenerateMotionsMany(t *testing.T) {
	generators := piece.MotionGenerators([]piece.MotionGenerator{
		piecetest.NewStaticMotionGenerator(t, "A1", "B2"),
		piecetest.NewStaticMotionGenerator(t, "C5"),
		piecetest.NewStaticMotionGenerator(t, "B2", "D4", "C5"),
	})
	motions := generators.GenerateMotions(nil)
	assertSquaresMatch(t, motions, "A1", "B2", "C5", "D4")
}
