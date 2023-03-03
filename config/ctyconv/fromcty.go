package ctyconv

import (
	"fmt"

	"github.com/jostrzol/mess/game"
	"github.com/jostrzol/mess/game/board"
	"github.com/jostrzol/mess/game/piece/color"
	plr "github.com/jostrzol/mess/game/player"
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

func PlayerFromCty(state *game.State, player cty.Value) (*plr.Player, error) {
	if player.IsNull() {
		return nil, nil
	}
	winnerColorStr := player.GetAttr("color").AsString()
	winnerColor, err := color.ColorString(winnerColorStr)
	if err != nil {
		return nil, fmt.Errorf("parsing player color: %w", err)
	}
	winner := state.GetPlayer(winnerColor)
	return winner, nil
}
