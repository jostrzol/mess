package room

import (
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/jostrzol/mess/pkg/rules"
	"github.com/jostrzol/mess/pkg/server/core/event"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
)

const PlayersNeeded = 2

type Room struct {
	id        id.Room
	players   [2]id.Session
	nPlayers  int
	RulesFile *rules.File
	game      id.Game
	mutex     sync.Mutex
}

func New() *Room {
	return &Room{
		id:        id.New[id.Room](),
		RulesFile: defaultRulesFile(),
	}
}

func defaultRulesFile() *rules.File {
	filename := "./rules/chess.hcl"
	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return &rules.File{
		Src:      src,
		Filename: filepath.Base(filename),
	}
}

func (r *Room) ID() id.Room {
	return r.id
}

func (r *Room) AddPlayer(sessionID id.Session) (event.Event, error) {
	r.mutex.Lock()
	defer func() { r.mutex.Unlock() }()
	if slices.Contains(r.players[:r.nPlayers], sessionID) {
		return nil, ErrAlreadyInRoom
	}
	if r.nPlayers >= PlayersNeeded {
		return nil, ErrRoomFull
	}
	r.players[r.nPlayers] = sessionID
	r.nPlayers++
	return &event.PlayerJoined{
		RoomID:   r.id,
		PlayerID: sessionID,
	}, nil
}

func (r *Room) Players() []id.Session {
	return r.players[:r.nPlayers]
}

func (r *Room) IsStarted() bool {
	return r.game != id.Game{}
}

func (r *Room) IsStartable() bool {
	return r.assertStartable() == nil
}

func (r *Room) assertStartable() error {
	switch {
	case r.nPlayers != PlayersNeeded:
		return ErrNotEnoughPlayers
	case r.IsStarted():
		return ErrAlreadyStarted
	default:
		return nil
	}
}

func (r *Room) Rules() *rules.File {
	return r.RulesFile
}

func (r *Room) UpdateRules(session id.Session, filename string, data []byte) (event.Event, error) {
	if filename == "" {
		return nil, usrerr.Errorf("filename cannot be empty")
	}

	r.RulesFile = &rules.File{Filename: filename, Src: data}

	return &event.RoomRulesChanged{RoomID: r.id, By: session}, nil
}

func (r *Room) StartGame(sessionID id.Session) (event.Event, error) {
	r.mutex.Lock()
	defer func() { r.mutex.Unlock() }()
	if err := r.assertStartable(); err != nil {
		return nil, err
	}
	r.game = id.New[id.Game]()
	return &event.GameStarted{
		GameID:  r.game,
		RoomID:  r.id,
		Players: r.players,
		Rules:   r.RulesFile,
		By:      sessionID,
	}, nil
}

func (r *Room) Game() id.Game {
	return r.game
}

var ErrRoomFull = usrerr.Errorf("room full")
var ErrNoRules = usrerr.Errorf("no rules file")
var ErrNotEnoughPlayers = usrerr.Errorf("not enough players")
var ErrAlreadyStarted = usrerr.Errorf("game is already started")
var ErrAlreadyInRoom = usrerr.Errorf("player already in room")
