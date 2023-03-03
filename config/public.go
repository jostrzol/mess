package config

import (
	"fmt"
	"log"

	"github.com/jostrzol/mess/game"
)

func DecodeConfig(filename string) (*game.State, game.Controller, error) {
	config, err := decodeConfig(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("decoding config: %w", err)
	}
	log.Printf("Loaded config: %#v", config)

	state, err := config.toGameState()
	if err != nil {
		return nil, nil, fmt.Errorf("initializing game state: %w", err)
	}

	controller := config.toController()

	return state, controller, nil
}
