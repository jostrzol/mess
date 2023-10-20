package rules

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/mess"
)

func DecodeRules(filename string, placePieces bool) (*mess.Game, error) {
	ctx := InitialEvalContext

	rules, err := decodeRules(filename, ctx)
	if err != nil {
		return nil, fmt.Errorf("decoding rules: %w", err)
	}

	game, err := rules.toEmptyGameState(ctx)
	if err != nil {
		return nil, fmt.Errorf("initializing game from rules: %w", err)
	}

	if placePieces {
		err = rules.placePieces(game.State)
		if err != nil {
			return nil, fmt.Errorf("placing initial pieces: %w", err)
		}
	}

	return game, nil
}
