package room

import (
	"fmt"

	"github.com/jostrzol/mess/pkg/server/core/event"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/ioc"
)

type Service struct {
	events     *event.Broker `container:"type"`
	repository Repository    `container:"type"`
}

func init() {
	ioc.MustSingletonFill[Service]()
}

func (s *Service) CreateRoom(sessionID id.Session) (*Room, error) {
	room := New()
	ev, err := room.AddPlayer(sessionID)
	if err != nil {
		return nil, fmt.Errorf("adding a player: %w", err)
	}
	err = s.repository.Save(room)
	if err != nil {
		return nil, fmt.Errorf("saving new room: %w", err)
	}
	s.events.Notify(ev)
	return room, nil
}

func (s *Service) JoinRoom(sessionID id.Session, roomID id.Room) (*Room, error) {
	room, err := s.repository.Get(roomID)
	if err != nil {
		return nil, fmt.Errorf("getting room %v: %w", roomID, err)
	}
	ev, err := room.AddPlayer(sessionID)
	if err != nil {
		return room, fmt.Errorf("adding a player: %w", err)
	}
	err = s.repository.Save(room)
	if err != nil {
		return room, fmt.Errorf("saving new room: %w", err)
	}
	s.events.Notify(ev)
	return room, nil
}

func (s *Service) GetRoom(roomID id.Room) (*Room, error) {
	room, err := s.repository.Get(roomID)
	if err != nil {
		return nil, fmt.Errorf("getting room %v: %w", roomID, err)
	}
	return room, nil
}

func (s *Service) StartGame(sessionID id.Session, roomID id.Room) (*Room, error) {
	room, err := s.repository.Get(roomID)
	if err != nil {
		return nil, fmt.Errorf("getting room %v: %w", roomID, err)
	}

	ev, err := room.StartGame(sessionID)
	if err != nil {
		return nil, fmt.Errorf("starting game: %w", err)
	}
	err = s.repository.Save(room)
	if err != nil {
		return room, fmt.Errorf("saving new room: %w", err)
	}
	s.events.Notify(ev)
	return room, nil
}
