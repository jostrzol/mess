package handler_test

import (
	"fmt"
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

func (s *GameSuite) TestChooseGameTurnOptionRoute() {
	// given
	room := s.Client().createStartedRoom()

	// when
	state := s.Client().chooseTurnOpionRoute(room.ID, 0, []any{
		map[string]any{
			"Type": "Move",
			"From": []any{0, 1},
			"To":   []any{0, 2},
		},
	})

	// expect
	s.NotZero(state)
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

func (c *GameClient) createStartedRoom() (room schema.Room) {
	room = c.createFilledRoom()
	room = c.startGame(room.ID)
	return
}

func (c *GameClient) chooseTurnOpionRoute(roomID uuid.UUID, turn int, route any) (state schema.State) {
	c.ServeHTTPOkAs("PUT", roomURL(roomID)+fmt.Sprintf("/game/turns/%v", turn), route, &state)
	return
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}
