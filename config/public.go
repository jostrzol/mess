package config

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/mess"
)

type Interactor interface {
	Choose(options []string) int
}

func DecodeConfig(filename string, interactor Interactor, placePieces bool) (*mess.Game, error) {
	ctx := InitialEvalContext

	config, err := decodeConfig(filename, ctx)
	if err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}

	game, err := config.toEmptyGameState(ctx, interactor)
	if err != nil {
		return nil, fmt.Errorf("initializing game from config: %w", err)
	}

	if placePieces {
		err = config.placePieces(game.State)
		if err != nil {
			return nil, fmt.Errorf("placing initial pieces: %w", err)
		}
	}

	return game, nil
}
