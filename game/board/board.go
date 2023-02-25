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

type PieceOnSquare struct {
	Piece  *piece.Piece
	Square *Square
}

func (b Board) PieceOn(square *Square) (PieceOnSquare, error) {
	if !b.contains(square) {
		err := fmt.Errorf("square %s out of board's bound", square)
		return PieceOnSquare{Square: square}, err
	}
	x, y := square.toCoords()
	piece := b[y][x]
	return PieceOnSquare{Piece: piece, Square: square}, nil
}

func (b Board) contains(square *Square) bool {
	x, y := square.toCoords()
	width, height := b.Size()
	return x < width && y < height
}

func (b Board) Place(piece *piece.Piece, square *Square) error {
	old, err := b.PieceOn(square)
	if err != nil {
		return fmt.Errorf("retrieving piece: %w", err)
	}
	if old.Piece != nil {
		log.Printf("replacing piece %q on %s with %q", old.Piece, square, piece)
	}
	x, y := square.toCoords()
	b[y][x] = piece
	return nil
}

func (b Board) AllPieces() []PieceOnSquare {
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
