package room

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/rules"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
)

type Room struct {
	ID      uuid.UUID
	players map[color.Color]uuid.UUID
	Rules   *string
	Game    *mess.Game
}

func New() *Room {
	return &Room{ID: uuid.New(), players: make(map[color.Color]uuid.UUID)}
}

func (r *Room) AddPlayer(sessionID uuid.UUID) error {
	for _, color := range color.ColorValues() {
		if _, present := r.players[color]; !present {
			r.players[color] = sessionID
			return nil
		}
	}
	return ErrRoomFull
}

func (r *Room) Start() error {
	if r.Rules == nil {
		return ErrNoRules
	}
	game, err := rules.DecodeRules(*r.Rules, true)
	if err != nil {
		return fmt.Errorf("decoding rules: %w", err)
	}
	r.Game = game
	return nil
}

var ErrRoomFull = usrerr.Errorf("room full")
var ErrNoRules = usrerr.Errorf("no rules file")
