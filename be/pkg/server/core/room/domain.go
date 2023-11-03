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
	id      uuid.UUID
	players map[color.Color]uuid.UUID
	rules   *string
	game    *mess.Game
}

func New() *Room {
	return &Room{id: uuid.New(), players: make(map[color.Color]uuid.UUID)}
}

func (r *Room) ID() uuid.UUID {
	return r.id
}

func (r *Room) AddPlayer(sessionID uuid.UUID) error {
	for _, color := range color.ColorValues() {
		playerID, present := r.players[color]
		if playerID == sessionID {
			return nil
		} else if !present {
			r.players[color] = sessionID
			return nil
		}
	}
	return ErrRoomFull
}

func (r *Room) Players() int {
	return len(r.players)
}

func (r *Room) Start() error {
	if r.rules == nil {
		return ErrNoRules
	}
	game, err := rules.DecodeRules(*r.rules, true)
	if err != nil {
		return fmt.Errorf("decoding rules: %w", err)
	}
	r.game = game
	return nil
}

var ErrRoomFull = usrerr.Errorf("room full")
var ErrNoRules = usrerr.Errorf("no rules file")
