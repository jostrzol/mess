package game

import (
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
)

type Repository interface {
	Save(game *Game) error
	Get(gameID id.Game) (*Game, error)
	GetForRoom(roomID id.Room) (*Game, error)
}

var ErrNotFound = usrerr.Errorf("game not found")
