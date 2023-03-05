package boardtest

import "github.com/jostrzol/mess/pkg/board"

func NewSquare(text string) *board.Square {
	square, err := board.NewSquare(text)
	if err != nil {
		panic(err)
	}
	return square
}
