package schema

import (
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/core/room"
)

type Room struct {
	ID            uuid.UUID
	Players       int
	PlayersNeeded int
	IsStartable   bool
	IsStarted     bool
}

func NewRoom(r *room.Room) *Room {
	return &Room{
		ID:            r.ID(),
		Players:       len(r.Players()),
		PlayersNeeded: room.PlayersNeeded,
		IsStartable:   r.IsStartable(),
		IsStarted:     r.IsStarted(),
	}
}
