package room

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/ioc"
)

type Service struct {
	repository Repository `container:"type"`
}

func init() {
	ioc.MustSingletonFill[Service]()
}

func (s *Service) CreateRoom(sessionID uuid.UUID) (*Room, error) {
	room := New()
	err := room.AddPlayer(sessionID)
	if err != nil {
		return nil, fmt.Errorf("adding a player: %w", err)
	}
	err = s.repository.Save(room)
	if err != nil {
		return nil, fmt.Errorf("saving new room: %w", err)
	}
	return room, nil
}

func (s *Service) JoinRoom(sessionID uuid.UUID, roomID uuid.UUID) (*Room, error) {
	room, err := s.repository.Get(roomID)
	if err != nil {
		return nil, fmt.Errorf("getting room %v: %w", roomID, err)
	}
	err = room.AddPlayer(sessionID)
	if err != nil {
		return room, fmt.Errorf("adding a player: %w", err)
	}
	err = s.repository.Save(room)
	if err != nil {
		return room, fmt.Errorf("saving new room: %w", err)
	}
	return room, nil
}

func (s *Service) GetRoom(roomID uuid.UUID) (*Room, error) {
	room, err := s.repository.Get(roomID)
	if err != nil {
		return nil, fmt.Errorf("getting room %v: %w", roomID, err)
	}
	return room, nil
}

func (s *Service) StartGame(roomID uuid.UUID) (*Room, error) {
	room, err := s.repository.Get(roomID)
	if err != nil {
		return nil, fmt.Errorf("getting room %v: %w", roomID, err)
	}
	err = room.StartGame()
	if err != nil {
		return room, err
	}
	return room, nil
}

func (s *Service) GetGameState(roomID uuid.UUID, sessionID uuid.UUID) (*State, error) {
	room, err := s.repository.Get(roomID)
	if err != nil {
		return nil, fmt.Errorf("getting room %v: %w", roomID, err)
	}

	game, err := room.Game()
	if err != nil {
		return nil, fmt.Errorf("getting game: %w", err)
	}

	optionTree, err := game.TurnOptions()
	if err != nil {
		return nil, fmt.Errorf("generating turn options: %w", err)
	}

	return &State{
		Board:      game.Board(),
		OptionTree: optionTree,
	}, nil
}
