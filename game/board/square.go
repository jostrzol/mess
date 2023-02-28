package board

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Square struct {
	File int
	Rank int
}

func NewSquare(text string) (Square, error) {
	var zero Square
	if len(text) != 2 {
		return zero, errors.New("malformed position: expected 2 characters")
	}

	fileRune := int(strings.ToUpper(text)[0])
	if fileRune < 'A' || fileRune > 'Z' {
		return zero, fmt.Errorf("expected letter, not %q", fileRune)
	}
	file := fileRune - 'A' + 1

	rank, err := strconv.Atoi(string(text[1]))
	if err != nil {
		return zero, fmt.Errorf("parsing rank: %v", err)
	}
	if rank <= 0 {
		return zero, fmt.Errorf("rank not positive: %d", rank)
	}

	return Square{File: file, Rank: rank}, nil
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

func fromCoords(x int, y int) *Square {
	return &Square{
		File: x + 1,
		Rank: y + 1,
	}
}
