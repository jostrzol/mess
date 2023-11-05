package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/adapter/handler/handlertest"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/stretchr/testify/suite"
)

type GameSuite struct {
	handlertest.HandlerSuite[GameClient]
}

func (s *GameSuite) TestStartGame() {
	// given
	room := s.Client().createFilledRoom()

	// when
	room = s.Client().startGame(room.ID)

	// then
	s.True(room.IsStarted)
}

func (s *GameSuite) TestGetGameState() {
	// given
	room := s.Client().createFilledRoom()

	// and
	room = s.Client().startGame(room.ID)

	// expect
	s.Client().getGameState(room.ID)
}

type GameClient struct{ RoomClient }

func (c *GameClient) createFilledRoom() (room schema.Room) {
	room = c.createRoom()
	for !room.IsStartable {
		c2 := handlertest.CloneWithEmptyJar(c)
		room = c2.joinRoom(room.ID)
	}
	return
}

func (c *GameClient) startGame(roomID uuid.UUID) (room schema.Room) {
	c.ServeHTTPOkAs("PUT", roomURL(roomID)+"/game", nil, &room)
	return
}

func (c *GameClient) getGameState(roomID uuid.UUID) (state schema.State) {
	c.ServeHTTPOkAs("GET", roomURL(roomID)+"/game", nil, &state)
	return
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}
