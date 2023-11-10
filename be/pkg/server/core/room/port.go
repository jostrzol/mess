package room

import (
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
)

type Repository interface {
	Save(room *Room) error
	Get(roomID id.Room) (*Room, error)
}

var ErrNotFound = usrerr.Errorf("room not found")
