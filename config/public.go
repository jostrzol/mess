package config

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/mess"
)

type Interactor interface {
	Choose(options []string) int
}

func DecodeConfig(filename string, interactor Interactor) (*mess.Game, error) {
	ctx := InitialEvalContext

	config, err := decodeConfig(filename, ctx)
	if err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}

	game, err := config.toGameState(ctx, interactor)
	if err != nil {
		return nil, fmt.Errorf("initializing game from config: %w", err)
	}

	return game, nil
}
