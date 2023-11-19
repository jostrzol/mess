package event

import (
	"github.com/golobby/container/v3"
	"github.com/jostrzol/mess/pkg/event"
	"github.com/jostrzol/mess/pkg/rules"
	"github.com/jostrzol/mess/pkg/server/core/id"
)

type Event = event.Event
type Subject = event.Subject
type Observer = event.Observer

type PlayerJoined struct {
	RoomID   id.Room
	PlayerID id.Session
}

type RoomRulesChanged struct {
	RoomID id.Room
	By     id.Session
}

type GameStarted struct {
	GameID  id.Game
	RoomID  id.Room
	Players [2]id.Session
	Rules   *rules.File
	By      id.Session
}

type GameChanged struct {
	GameID id.Game
	By     id.Session
}

type Broker struct {
	Subject
}

func init() {
	container.MustSingletonLazy(container.Global, func() *Broker {
		return &Broker{event.NewSubject()}
	})
}
