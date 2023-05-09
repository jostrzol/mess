package config

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/mess"
)

func DecodeConfig(filename string) (*mess.Game, error) {
	config, err := decodeConfig(filename)
	if err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}

	err = config.toGameState()
	if err != nil {
		return nil, fmt.Errorf("initializing game state: %w", err)
	}

	controller := config.toController()

	game := mess.NewGame(config.State, controller)
	return game, nil
}
