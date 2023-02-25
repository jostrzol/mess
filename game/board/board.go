package board

import (
	"errors"
	"fmt"
	"log"

	"github.com/jostrzol/mess/game/piece"
)

type Board [][]*piece.Piece

func NewBoard(width int, height int) (Board, error) {
	if height <= 0 || width <= 0 {
		return nil, errors.New("one of board dimentions is non-positive")
	}

	board := make(Board, height)
	for i := range board {
		board[i] = make([]*piece.Piece, width)
	}
	return board, nil
}

func (b Board) Size() (int, int) {
	row := b[0]
	return len(row), len(b)
}

func (b Board) Place(piece *piece.Piece, square *Square) error {
	width, height := b.Size()
	x, y := square.toCoords()
	if x >= width || y >= height {
		return fmt.Errorf("square %s out of board's bound (size: %dx%d)", square, width, height)
	}
	old := b[y][x]
	if old != nil {
		log.Printf("replacing piece %q on %s with %q", old, square, piece)
	}
	b[y][x] = piece
	return nil
}

type PieceOnSquare struct {
	Piece  *piece.Piece
	Square Square
}

func (b Board) Pieces() []PieceOnSquare {
	pieces := make([]PieceOnSquare, 0)
	for y, row := range b {
		for x, piece := range row {
			if piece == nil {
				continue
			}
			pieces = append(pieces, PieceOnSquare{
				Piece:  piece,
				Square: fromCoords(x, y),
			})
		}
	}
	return pieces
}
