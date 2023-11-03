package room

import (
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
)

type Repository interface {
	Save(room *Room) error
	Get(roomID uuid.UUID) (*Room, error)
}

var ErrNotFound = usrerr.Errorf("room not found")
