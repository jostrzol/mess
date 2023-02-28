package piecetest

import (
	"testing"

	brd "github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
	"github.com/stretchr/testify/assert"
)

type StubMotionGenerator struct {
	t       *testing.T
	strings []string
}

func (s StubMotionGenerator) GenerateMotions(piece *piece.Piece) []brd.Square {
	s.t.Helper()
	destinations := make([]brd.Square, 0, len(s.strings))
	for _, squareStr := range s.strings {
		square, err := brd.NewSquare(squareStr)
		assert.NoError(s.t, err)
		destinations = append(destinations, square)
	}
	return destinations
}

func NewStubMotionGenerator(t *testing.T, strings ...string) StubMotionGenerator {
	return StubMotionGenerator{t, strings}
}
