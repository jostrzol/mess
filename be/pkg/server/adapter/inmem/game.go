package inmem

import (
	"github.com/golobby/container/v3"
	"github.com/jostrzol/mess/pkg/server/core/game"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/core/room"
)

type GameRepository struct {
	games map[id.Game]*game.Game
	rooms room.Repository `container:"type"`
}

func NewGameRepository() *GameRepository {
	return &GameRepository{games: make(map[id.Game]*game.Game)}
}

func init() {
	container.MustSingletonLazy(container.Global, func() game.Repository {
		repo := NewGameRepository()
		container.MustFill(container.Global, repo)
		return repo
	})
}

func (r *GameRepository) Save(game *game.Game) error {
	r.games[game.ID()] = game
	return nil
}

func (r *GameRepository) Get(gameID id.Game) (*game.Game, error) {
	result, ok := r.games[gameID]
	if !ok {
		return nil, game.ErrNotFound
	}
	return result, nil
}

func (r *GameRepository) GetForRoom(roomID id.Room) (*game.Game, error) {
	room, err := r.rooms.Get(roomID)
	if err != nil {
		return nil, err
	}
	result, ok := r.games[room.Game()]
	if !ok {
		return nil, game.ErrNotFound
	}
	return result, nil
}
