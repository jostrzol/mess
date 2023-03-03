package piecetest

import (
	"testing"

	"github.com/jostrzol/mess/game/board"
	brd "github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
	"github.com/stretchr/testify/assert"
)

func NewStaticMotionGenerator(t *testing.T, strings ...string) piece.MotionGenerator {
	t.Helper()
	return piece.FuncMotionGenerator(func(piece *piece.Piece) []brd.Square {
		destinations := make([]brd.Square, 0, len(strings))
		for _, squareStr := range strings {
			square, err := brd.NewSquare(squareStr)
			assert.NoError(t, err)
			destinations = append(destinations, *square)
		}
		return destinations
	})
}

func NewOffsetMotionGenerator(t *testing.T, offsets ...board.Offset) piece.MotionGenerator {
	t.Helper()
	return piece.FuncMotionGenerator(func(piece *piece.Piece) []brd.Square {
		destinations := make([]brd.Square, 0, len(offsets))
		for _, offset := range offsets {
			square := piece.Square.Offset(offset)
			if piece.Board.Contains(square) {
				destinations = append(destinations, *square)
			}
		}
		return destinations
	})
}
