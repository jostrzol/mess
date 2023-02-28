package piecetest

import (
	"testing"

	"github.com/jostrzol/mess/game/piece"
)

func Rook(t *testing.T) *piece.Type {
	t.Helper()
	pieceType := &piece.Type{
		Name: "rook",
	}
	pieceType.AddMotionGenerator(NewStubMotionGenerator(t, "A1", "A2"))
	return pieceType
}

func Knight(t *testing.T) *piece.Type {
	t.Helper()
	pieceType := &piece.Type{
		Name: "knight",
	}
	pieceType.AddMotionGenerator(NewStubMotionGenerator(t, "B1", "B2"))
	return pieceType
}

func Noones(pieceType *piece.Type) *piece.Piece {
	return &piece.Piece{
		Type:  pieceType,
		Owner: nil,
	}
}
