package room

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/rules"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
	"golang.org/x/exp/maps"
)

// TODO: make this dynamic.
const rulesFile = "./rules/chess.hcl"

const NeededPlayers = 2

type Room struct {
	id      uuid.UUID
	players map[color.Color]uuid.UUID
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

func (r *Room) Players() []uuid.UUID {
	return maps.Values(r.players)
}

func (r *Room) IsStarted() bool {
	return r.game != nil
}

func (r *Room) IsStartable() bool {
	return r.assertStartable() == nil
}

func (r *Room) assertStartable() error {
	switch {
	case len(r.players) != NeededPlayers:
		return ErrNotEnoughPlayers
	case r.IsStarted():
		return ErrNotEnoughPlayers
	default:
		return nil
	}
}

func (r *Room) Start() error {
	if r.IsStarted() {
		return nil
	}
	if err := r.assertStartable(); err != nil {
		return err
	}
	game, err := rules.DecodeRules(rulesFile, true)
	if err != nil {
		return fmt.Errorf("decoding rules: %w", err)
	}
	r.game = game
	return nil
}

var ErrRoomFull = usrerr.Errorf("room full")
var ErrNoRules = usrerr.Errorf("no rules file")
var ErrNotEnoughPlayers = usrerr.Errorf("not enough players")
var ErrAlreadyStarted = usrerr.Errorf("game is already started")
