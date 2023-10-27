package room

import "github.com/google/uuid"

type Repository interface {
	Save(room *Room) error
	Get(roomID uuid.UUID) (*Room, error)
}
