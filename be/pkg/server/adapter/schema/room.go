package schema

import (
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/core/room"
)

type Room struct {
	ID            uuid.UUID
	Players       int
	NeededPlayers int
	IsStartable   bool
	IsStarted     bool
}

func NewRoom(r *room.Room) *Room {
	return &Room{
		ID:            r.ID(),
		Players:       len(r.Players()),
		NeededPlayers: room.NeededPlayers,
		IsStartable:   r.IsStartable(),
		IsStarted:     r.IsStarted(),
	}
}
