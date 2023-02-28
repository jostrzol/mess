package board

import (
	"errors"
	"fmt"
	"log"
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

func (b Board[T]) Size() (int, int) {
	row := b[0]
	return len(row), len(b)
}

func (b Board[T]) At(square *Square) (T, error) {
	var zero T
	if !b.contains(square) {
		err := fmt.Errorf("square %s out of board's bound", square)
		return zero, err
	}
	x, y := square.toCoords()
	item := b[y][x]
	return item, nil
}

func (b Board[T]) contains(square *Square) bool {
	x, y := square.toCoords()
	width, height := b.Size()
	return x < width && y < height
}

func (b Board[T]) Place(item T, square *Square) error {
	var zero T
	old, err := b.At(square)
	if err != nil {
		return fmt.Errorf("retrieving item: %w", err)
	}
	if old != zero {
		log.Printf("replacing item '%v' on %s with '%v'", old, square, item)
	}
	x, y := square.toCoords()
	b[y][x] = item
	return nil
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
