package inmem

import (
	"github.com/golobby/container/v3"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/core/room"
)

type RoomRepository struct {
	rooms map[id.Room]*room.Room
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{rooms: make(map[id.Room]*room.Room)}
}

func init() {
	container.MustSingletonLazy(container.Global, func() room.Repository {
		return NewRoomRepository()
	})
}

func (r *RoomRepository) Save(room *room.Room) error {
	r.rooms[room.ID()] = room
	return nil
}

func (r *RoomRepository) Get(roomID id.Room) (*room.Room, error) {
	result, ok := r.rooms[roomID]
	if !ok {
		return nil, room.ErrNotFound
	}
	return result, nil
}
