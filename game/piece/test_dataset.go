package piece

import "testing"

func Rook(t *testing.T) *Piece {
	t.Helper()
	return &Piece{
		Type: &PieceType{
			Name: "rook",
		},
		Owner: nil,
	}
}

func Knight(t *testing.T) *Piece {
	t.Helper()
	return &Piece{
		Type: &PieceType{
			Name: "knight",
		},
		Owner: nil,
	}
}
