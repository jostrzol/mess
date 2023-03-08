package config

import (
	"fmt"
	"log"

	"github.com/jostrzol/mess/pkg/mess"
)

func DecodeConfig(filename string) (*mess.State, mess.Controller, error) {
	state := new(mess.State)
	ctx := newEvalContext(state)

	config, err := decodeConfig(filename, ctx, state)
	if err != nil {
		return nil, nil, fmt.Errorf("decoding config: %w", err)
	}
	log.Printf("Loaded config: %#v", config)

	initializedState, err := config.toGameState()
	if err != nil {
		return nil, nil, fmt.Errorf("initializing game state: %w", err)
	}
	*state = *initializedState

	controller := config.toController()

	return state, controller, nil
}
