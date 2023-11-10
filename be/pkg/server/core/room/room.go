package room

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/rules"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
	"golang.org/x/exp/maps"
)

// TODO: Make this dynamic.
const rulesFile = "./rules/chess.hcl"

const PlayersNeeded = 2

type Room struct {
	id      uuid.UUID
	players map[color.Color]uuid.UUID
	game    *mess.Game
	mutex   sync.Mutex
}

func New() *Room {
	return &Room{id: uuid.New(), players: make(map[color.Color]uuid.UUID)}
}

func (r *Room) ID() uuid.UUID {
	return r.id
}

func (r *Room) AddPlayer(sessionID uuid.UUID) error {
	r.mutex.Lock()
	defer func() { r.mutex.Unlock() }()
	for _, color := range color.ColorValues() {
		playerID, present := r.players[color]
		if playerID == sessionID {
			return ErrAlreadyInRoom
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
	case len(r.players) != PlayersNeeded:
		return ErrNotEnoughPlayers
	case r.IsStarted():
		return ErrAlreadyStarted
	default:
		return nil
	}
}

func (r *Room) StartGame() error {
	r.mutex.Lock()
	defer func() { r.mutex.Unlock() }()
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

func (r *Room) GameState() (*State, error) {
	if !r.IsStarted() {
		return nil, ErrNotStarted
	}

	r.mutex.Lock()
	defer func() { r.mutex.Unlock() }()

	optionTree, err := r.game.TurnOptions()
	if err != nil {
		return nil, fmt.Errorf("generating turn options: %w", err)
	}

	return &State{
		TurnNumber: r.game.TurnNumber(),
		Board:      r.game.Board(),
		OptionTree: optionTree,
	}, nil
}

func (r *Room) PlayTurn(turn int, route mess.Route) error {
	if !r.IsStarted() {
		return ErrNotStarted
	}

	r.mutex.Lock()
	defer func() { r.mutex.Unlock() }()

	currentTurn := r.game.State.TurnNumber()
	if turn < currentTurn {
		// TODO: make idempotent -- check if options match
		// if match => noop
		// if not => error
		return ErrTurnTooSmall
	} else if turn > currentTurn {
		return ErrTurnTooBig
	}

	err := r.game.Turn(route)
	if err != nil {
		return fmt.Errorf("choosing turn options: %w", err)
	}
	return nil
}

var ErrRoomFull = usrerr.Errorf("room full")
var ErrNoRules = usrerr.Errorf("no rules file")
var ErrNotEnoughPlayers = usrerr.Errorf("not enough players")
var ErrAlreadyStarted = usrerr.Errorf("game is already started")
var ErrAlreadyInRoom = usrerr.Errorf("player already in room")
var ErrNotStarted = usrerr.Errorf("game not started")
var ErrTurnTooSmall = usrerr.Errorf("the selected turn has already been played")
var ErrTurnTooBig = usrerr.Errorf("the selected turn hasn't started yet")
