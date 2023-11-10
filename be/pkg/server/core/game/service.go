package game

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/server/core/event"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/ioc"
	"go.uber.org/zap"
)

type Service struct {
	events     *event.Broker `container:"type"`
	repository Repository    `container:"type"`
	logger     *zap.Logger   `container:"type"`
}

func init() {
	ioc.MustSingletonObserverFill[Service]()
}

func (s *Service) GetGameState(sessionID id.Session, roomID id.Room) (*State, error) {
	game, err := s.repository.GetForRoom(roomID)
	if err != nil {
		return nil, fmt.Errorf("getting game %v: %w", roomID, err)
	}

	state, err := game.State()
	if err != nil {
		return nil, fmt.Errorf("calculating game state: %w", err)
	}

	return state, nil
}

func (s *Service) PlayTurn(sessionID id.Session, roomID id.Room, turn int, route mess.Route) (*State, error) {
	game, err := s.repository.GetForRoom(roomID)
	if err != nil {
		return nil, fmt.Errorf("getting room %v: %w", roomID, err)
	}

	ev, err := game.PlayTurn(sessionID, turn, route)
	if err != nil {
		return nil, fmt.Errorf("playing turn: %w", err)
	}
	err = s.repository.Save(game)
	if err != nil {
		return nil, fmt.Errorf("saving game: %w", err)
	}
	s.events.Notify(ev)

	state, err := game.State()
	if err != nil {
		return nil, fmt.Errorf("calculating game state: %w", err)
	}

	return state, nil
}

func (s *Service) Handle(evnt event.Event) {
	switch ev := evnt.(type) {
	case *event.GameStarted:
		game, err := New(ev)
		if err != nil {
			s.logger.Error("creating game", zap.Error(err))
			return
		}
		err = s.repository.Save(game)
		if err != nil {
			s.logger.Error("saving game", zap.Error(err))
			return
		}
	}
}
