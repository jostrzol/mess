package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Board [][]*Piece

type Square struct {
	File int
	Rank int
}

func NewBoard(width int, height int) (Board, error) {
	if height <= 0 || width <= 0 {
		return nil, errors.New("one of board dimentions is non-positive")
	}

	board := make(Board, height)
	for i := range board {
		board[i] = make([]*Piece, width)
	}
	return board, nil
}

func (b Board) Size() (int, int) {
	row := b[0]
	return len(row), len(b)
}

func (b Board) Place(piece *Piece, square Square) error {
	width, height := b.Size()
	x, y := square.toCoords()
	if x >= width || y >= height {
		return errors.New("square out of board's bound")
	}
	b[y][x] = piece
	return nil
}

type PieceOnSquare struct {
	Piece  *Piece
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

func ParseSquare(text string) (*Square, error) {
	if len(text) != 2 {
		return nil, errors.New("malformed position: expected 2 characters")
	}

	fileRune := int(strings.ToUpper(text)[0])
	if fileRune < 'A' || fileRune > 'Z' {
		return nil, fmt.Errorf("expected letter, not %q", fileRune)
	}
	file := fileRune - 'A' + 1

	rank, err := strconv.Atoi(string(text[1]))
	if err != nil {
		return nil, fmt.Errorf("parsing rank: %v", err)
	}
	if rank <= 0 {
		return nil, fmt.Errorf("rank not positive: %d", rank)
	}

	return &Square{File: file, Rank: rank}, nil
}

func (s *Square) String() string {
	file := string(byte(s.File-1) + 'A')
	return fmt.Sprintf("%s%d", file, s.Rank)
}

func (s *Square) toCoords() (int, int) {
	x := s.File - 1
	y := s.Rank - 1
	return x, y
}

func fromCoords(x int, y int) Square {
	return Square{
		File: x + 1,
		Rank: y + 1,
	}
}
