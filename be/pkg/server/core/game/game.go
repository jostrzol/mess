package game

import (
	"fmt"
	"sync"

	"github.com/jostrzol/mess/pkg/color"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/rules"
	"github.com/jostrzol/mess/pkg/server/core/event"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
	"golang.org/x/exp/maps"
)

type Game struct {
	id      id.Game
	room    id.Room
	players map[color.Color]id.Session
	mutex   sync.Mutex
	game    *mess.Game
}

type State struct {
	ID         id.Game
	TurnNumber int
	Board      *mess.PieceBoard
	OptionTree *mess.OptionNode
	State      *mess.State
	IsMyTurn   bool
}

func New(event *event.GameStarted) (*Game, error) {
	game, err := rules.DecodeRules(event.Rules, true)
	if err != nil {
		return nil, fmt.Errorf("decoding rules: %w", err)
	}
	return &Game{
		id:   event.GameID,
		room: event.RoomID,
		players: map[color.Color]id.Session{
			color.White: event.Players[color.White],
			color.Black: event.Players[color.Black],
		},
		game: game,
	}, nil
}

func (g *Game) ID() id.Game {
	return g.id
}

func (g *Game) RoomID() id.Room {
	return g.room
}

func (g *Game) Players() []id.Session {
	return maps.Values(g.players)
}

func (g *Game) IsCurrentPlayer(session id.Session) bool {
	currentPlayerSession := g.players[g.game.CurrentPlayer().Color()]
	return currentPlayerSession == session
}

func (g *Game) State(session id.Session) (*State, error) {
	g.mutex.Lock()
	defer func() { g.mutex.Unlock() }()

	optionTree, err := g.game.TurnOptions()
	if err != nil {
		return nil, fmt.Errorf("generating turn options: %w", err)
	}

	return &State{
		ID:         g.id,
		TurnNumber: g.game.TurnNumber(),
		Board:      g.game.Board(),
		OptionTree: optionTree,
		State:      g.game.State,
		IsMyTurn:   g.IsCurrentPlayer(session),
	}, nil
}

func (g *Game) PlayTurn(session id.Session, turn int, route mess.Route) (event.Event, error) {
	g.mutex.Lock()
	defer func() { g.mutex.Unlock() }()

	if !g.IsCurrentPlayer(session) {
		return nil, ErrNotYourTurn
	}

	currentTurn := g.game.State.TurnNumber()
	if turn < currentTurn {
		// TODO: make idempotent -- check if options match
		// if match => noop
		// if not => error
		return nil, ErrTurnTooSmall
	} else if turn > currentTurn {
		return nil, ErrTurnTooBig
	}

	err := g.game.Turn(route)
	if err != nil {
		return nil, usrerr.Errorf("choosing turn options: %w", err)
	}
	return &event.GameChanged{
		GameID: g.id,
		By:     session,
	}, nil
}

var ErrTurnTooSmall = usrerr.Errorf("the selected turn has already been played")
var ErrTurnTooBig = usrerr.Errorf("the selected turn hasn't started yet")
var ErrNotYourTurn = usrerr.Errorf("it's not your turn")
