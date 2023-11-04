package schema

import (
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/core/room"
)

type Room struct {
	ID          uuid.UUID
	Players     int
	IsStartable bool
	IsStarted   bool
}

func NewRoom(room *room.Room) *Room {
	return &Room{
		ID:          room.ID(),
		Players:     len(room.Players()),
		IsStartable: room.IsStartable(),
		IsStarted:   room.IsStarted(),
	}
}
