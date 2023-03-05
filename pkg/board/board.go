package board

import (
	"errors"
	"fmt"
	"strings"
)

type Board[T comparable] [][]T

func NewBoard[T comparable](width int, height int) (Board[T], error) {
	if height <= 0 || width <= 0 {
		return nil, errors.New("one of board dimentions is non-positive")
	}

	board := make(Board[T], height)
	for i := range board {
		board[i] = make([]T, width)
	}
	return board, nil
}

func (b Board[T]) String() string {
	var zero T
	var builder strings.Builder
	for y, row := range b {
		for x, item := range row {
			if item != zero {
				square := fromCoords(x, y)
				fmt.Fprintf(&builder, "%v: %v\n", square, item)
			}
		}
	}
	return builder.String()
}

func (b Board[T]) Size() (int, int) {
	row := b[0]
	return len(row), len(b)
}

func (b Board[T]) At(square *Square) (T, error) {
	var zero T
	if !b.Contains(square) {
		err := fmt.Errorf("square %s out of board's bound", square)
		return zero, err
	}
	x, y := square.toCoords()
	item := b[y][x]
	return item, nil
}

func (b Board[T]) Contains(square *Square) bool {
	x, y := square.toCoords()
	width, height := b.Size()
	return x < width && x >= 0 && y < height && y >= 0
}

func (b Board[T]) Place(item T, square *Square) (T, error) {
	var zero T
	old, err := b.At(square)
	if err != nil {
		return zero, fmt.Errorf("retrieving item: %w", err)
	}
	x, y := square.toCoords()
	b[y][x] = item
	return old, nil
}

func (b Board[T]) AllItems() []T {
	var zero T
	items := make([]T, 0)
	for _, row := range b {
		for _, item := range row {
			if item == zero {
				continue
			}
			items = append(items, item)
		}
	}
	return items
}
