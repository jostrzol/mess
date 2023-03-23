package config

import (
	"fmt"
	"log"

	"github.com/jostrzol/mess/pkg/mess"
)

func DecodeConfig(filename string) (*mess.State, mess.Controller, error) {
	config, err := decodeConfig(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("decoding config: %w", err)
	}
	log.Printf("Loaded config: %#v", config)

	err = config.toGameState()
	if err != nil {
		return nil, nil, fmt.Errorf("initializing game state: %w", err)
	}

	controller := config.toController()

	return config.State, controller, nil
}
