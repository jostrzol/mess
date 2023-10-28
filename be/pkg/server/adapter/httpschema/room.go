package httpschema

import (
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/core/room"
)

type Room struct {
	ID uuid.UUID
}

func NewRoom(room *room.Room) *Room {
	return &Room{
		ID: room.ID,
	}
}
