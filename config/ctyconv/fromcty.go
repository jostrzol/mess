package ctyconv

import (
	"fmt"

	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece/color"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func SquareFromCty(value cty.Value) (*board.Square, error) {
	var squareStr string
	if err := gocty.FromCtyValue(value, &squareStr); err != nil {
		return nil, err
	}
	square, err := board.NewSquare(squareStr)
	if err != nil {
		return nil, fmt.Errorf("parsing square %q: %w", squareStr, err)
	}
	return square, nil
}

func SquaresFromCty(value cty.Value) ([]board.Square, error) {
	var squareStrs []string
	if err := gocty.FromCtyValue(value, &squareStrs); err != nil {
		return nil, err
	}
	errors := make(manyErrors, 0)
	result := make([]board.Square, 0, len(squareStrs))
	for _, squareStr := range squareStrs {
		square, err := board.NewSquare(squareStr)
		if err != nil {
			errors = append(errors, fmt.Errorf("parsing square %q: %w", squareStr, err))
		} else {
			result = append(result, *square)
		}
	}
	if len(errors) > 0 {
		return result, errors
	}
	return result, nil
}

func ColorFromCty(colorCty cty.Value) (*color.Color, error) {
	if colorCty.IsNull() {
		return nil, nil
	}
	var colorStr string
	var err error
	if err := gocty.FromCtyValue(colorCty, &colorStr); err != nil {
		return nil, err
	}
	color, err := color.ColorString(colorStr)
	if err != nil {
		return nil, err
	}
	return &color, nil
}
