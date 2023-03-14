package board

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
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
	isFirst := true
	for y, row := range b {
		for x, item := range row {
			if item != zero {
				square := fromCoords(x, y)
				if !isFirst {
					builder.WriteRune('\n')
				}
				isFirst = false
				fmt.Fprintf(&builder, "%v: %v", square, item)
			}
		}
	}
	return builder.String()
}

func (b Board[T]) PrettyString(itemFormatter func(T) rune) string {
	var builder strings.Builder
	b.printBar(&builder)
	builder.WriteRune('\n')
	for y := len(b) - 1; y >= 0; y-- {
		row := b[y]
		b.printRow(&builder, y+1, row, itemFormatter)
		builder.WriteRune('\n')
		b.printBar(&builder)
		builder.WriteRune('\n')
	}
	b.printFileHeader(&builder)
	builder.WriteRune('\n')
	b.printBar(&builder)
	return builder.String()
}

func (b Board[T]) printBar(w io.ByteWriter) {
	width, _ := b.Size()
	widthWithHeader := width + 1
	for i := 0; i < widthWithHeader*3+1; i++ {
		w.WriteByte('-')
	}
}

func (b Board[T]) printRow(w io.Writer, rank int, row []T, itemFormatter func(T) rune) {
	fmt.Fprintf(w, "|%2d|", rank)
	for _, item := range row {
		sign := itemFormatter(item)
		bytes := make([]byte, utf8.RuneLen(sign))
		if n := utf8.EncodeRune(bytes, sign); n != len(bytes) {
			panic(fmt.Errorf("printing board row: expected to write %d bytes but wrote %d", len(bytes), n))
		}
		fmt.Fprintf(w, "%-2s|", bytes)
	}
}

func (b Board[T]) printFileHeader(w io.Writer) {
	w.Write([]byte("|  |"))
	width, _ := b.Size()
	for i := 1; i <= width; i++ {
		fmt.Fprintf(w, "%-2s|", fileString(i))
	}
}

func (b Board[T]) Size() (int, int) {
	row := b[0]
	return len(row), len(b)
}

func (b Board[T]) At(square Square) (T, error) {
	var zero T
	if !b.Contains(square) {
		err := fmt.Errorf("square %s out of board's bound", square)
		return zero, err
	}
	x, y := square.toCoords()
	item := b[y][x]
	return item, nil
}

func (b Board[T]) Contains(square Square) bool {
	x, y := square.toCoords()
	width, height := b.Size()
	return x < width && x >= 0 && y < height && y >= 0
}

func (b Board[T]) Place(item T, square Square) (T, error) {
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
