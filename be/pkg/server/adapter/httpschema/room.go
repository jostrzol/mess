package httpschema

import (
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/core/room"
)

type Room struct {
	ID      uuid.UUID
	Players int
}

func NewRoom(room *room.Room) *Room {
	return &Room{
		ID:      room.ID(),
		Players: len(room.Players()),
	}
}
