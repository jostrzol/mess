package repository

import (
	"fmt"

	"github.com/golobby/container/v3"
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/core/room"
)

type RoomRepository struct {
	rooms map[uuid.UUID]*room.Room
}

func New() *RoomRepository {
	return &RoomRepository{rooms: make(map[uuid.UUID]*room.Room)}
}

func init() {
	container.MustSingleton(container.Global, func() room.Repository {
		return New()
	})
}

func (r *RoomRepository) Save(room *room.Room) error {
	r.rooms[room.ID] = room
	return nil
}

func (r *RoomRepository) Get(roomID uuid.UUID) (*room.Room, error) {
	room, ok := r.rooms[roomID]
	if !ok {
		return nil, ErrNotFound
	}
	return room, nil
}

var ErrNotFound = fmt.Errorf("room not found")
