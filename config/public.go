package config

import (
	"fmt"
	"log"

	"github.com/jostrzol/mess/pkg/mess"
)

func DecodeConfig(filename string) (*mess.State, mess.Controller, error) {
	var state mess.State
	config, err := decodeConfig(filename, &state)
	if err != nil {
		return nil, nil, fmt.Errorf("decoding config: %w", err)
	}
	log.Printf("Loaded config: %#v", config)

	err = config.toGameState(&state)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing game state: %w", err)
	}

	controller := config.toController()

	return &state, controller, nil
}
