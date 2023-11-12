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
	// cachedState is a Read-Only version of the current game state,
	// not including cumputation-heavy fields. It should be calculated
	// eagerly just after game state change and be valid all the time.
	// Should be accessed through State() method.
	cachedState      *State
	cachedPieceTypes map[string]*mess.PieceType
}

type State struct {
	ID            id.Game
	TurnNumber    int
	Board         *mess.PieceBoard
	PieceTypes    map[string]*mess.PieceType
	CurrentPlayer id.Session
}

type StaticData struct {
	ID        id.Game
	BoardSize BoardSize
	MyColor   color.Color
}

type BoardSize struct {
	Width  int
	Height int
}

type Resolution struct {
	IsResolved bool
	Winner     id.Session
}

func New(event *event.GameStarted) (*Game, error) {
	game, err := rules.DecodeRules(event.Rules, true)
	if err != nil {
		return nil, fmt.Errorf("decoding rules: %w", err)
	}
	result := &Game{
		id:   event.GameID,
		room: event.RoomID,
		players: map[color.Color]id.Session{
			color.White: event.Players[color.White],
			color.Black: event.Players[color.Black],
		},
		mutex:            sync.Mutex{},
		game:             game,
		cachedPieceTypes: game.PieceTypesByName(),
	}
	result.calculateState()

	return result, nil
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

func (g *Game) CurrentPlayer() id.Session {
	return g.cachedState.CurrentPlayer
}

func (g *Game) StaticData(session id.Session) *StaticData {
	return &StaticData{
		ID:        g.id,
		BoardSize: g.boardSize(),
		MyColor:   g.playerColor(session),
	}
}

func (g *Game) boardSize() BoardSize {
	width, height := g.game.Board().Size()
	return BoardSize{
		Width:  width,
		Height: height,
	}
}

func (g *Game) playerColor(session id.Session) color.Color {
	for color, player := range g.players {
		if player == session {
			return color
		}
	}
	panic(fmt.Errorf("no color for player %v", session))
}

func (g *Game) State() *State {
	if g.cachedState == nil {
		panic(fmt.Errorf("state not calculated"))
	}
	return g.cachedState
}

func (g *Game) TurnOptions() (*mess.OptionNode, error) {
	g.mutex.Lock()
	defer func() { g.mutex.Unlock() }()

	optionTree, err := g.game.TurnOptions()
	if err != nil {
		return nil, fmt.Errorf("generating turn options: %w", err)
	}

	return optionTree, nil
}

func (g *Game) PlayTurn(session id.Session, turn int, route mess.Route) (event.Event, error) {
	g.mutex.Lock()
	defer func() { g.mutex.Unlock() }()

	if g.CurrentPlayer() != session {
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

	g.calculateState()
	return &event.GameChanged{
		GameID: g.id,
		By:     session,
	}, nil
}

// calculateState caches the current game state.
// Presumes that THE MUTEX IS LOCKED!
func (g *Game) calculateState() {
	g.cachedState = &State{
		ID:            g.id,
		TurnNumber:    g.game.TurnNumber(),
		Board:         g.game.Board().Clone(),
		CurrentPlayer: g.players[g.game.CurrentPlayer().Color()],
		PieceTypes:    g.cachedPieceTypes,
	}
}

func (g *Game) Resolution() *Resolution {
	g.mutex.Lock()
	defer func() { g.mutex.Unlock() }()

	isResolved, winner := g.game.PickWinner()
	var winnerSession id.Session
	if winner != nil {
		winnerSession = g.players[winner.Color()]
	}

	return &Resolution{
		IsResolved: isResolved,
		Winner:     winnerSession,
	}
}

var ErrTurnTooSmall = usrerr.Errorf("the selected turn has already been played")
var ErrTurnTooBig = usrerr.Errorf("the selected turn hasn't started yet")
var ErrNotYourTurn = usrerr.Errorf("it's not your turn")
