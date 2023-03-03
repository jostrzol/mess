package piecetest

import (
	"testing"

	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece"
)

func Rook(t *testing.T) *piece.Type {
	t.Helper()
	pieceType := &piece.Type{
		Name: "rook",
	}
	pieceType.AddMotionGenerator(piece.FuncMotionGenerator(func(piece *piece.Piece) []board.Square {
		result := make([]board.Square, 0)
		for _, offset := range []board.Offset{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			square := piece.Square.Offset(offset)
			for piece.Board.Contains(square) {
				result = append(result, *square)
				square = square.Offset(offset)
			}
		}
		return result
	}))
	return pieceType
}

func Knight(t *testing.T) *piece.Type {
	t.Helper()
	pieceType := &piece.Type{
		Name: "knight",
	}
	pieceType.AddMotionGenerator(NewOffsetMotionGenerator(t, []board.Offset{
		{X: 1, Y: 2},
		{X: 1, Y: -2},
		{X: -1, Y: 2},
		{X: -1, Y: -2},
		{X: 2, Y: 1},
		{X: 2, Y: -1},
		{X: -2, Y: 1},
		{X: -2, Y: -1},
	}...))
	return pieceType
}

func Noones(pieceType *piece.Type) *piece.Piece {
	return &piece.Piece{
		Type:  pieceType,
		Owner: nil,
	}
}
