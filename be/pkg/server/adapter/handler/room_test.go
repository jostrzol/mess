package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/adapter/handler/handlertest"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/stretchr/testify/suite"
)

type RoomSuite struct {
	handlertest.HandlerSuite[RoomClient]
}

func (s *RoomSuite) TestCreateRoom() {
	// expect
	s.Client().createRoom()
}

func (s *RoomSuite) TestJoinRoomSameClient() {
	// given
	room := s.Client().createRoom()

	// when
	room = s.Client().joinRoom(room.ID)

	// then
	s.Equal(1, room.Players)
}

func (s *RoomSuite) TestJoinRoomDifferentClient() {
	// given
	room := s.Client().createRoom()

	// and
	c2 := s.NewClient()

	// when
	room = c2.joinRoom(room.ID)

	// then
	s.Equal(2, room.Players)
}

func (s *RoomSuite) TestGetRoom() {
	// given
	room := s.Client().createRoom()

	// expect
	s.Client().getRoom(room.ID)
}

func (s *RoomSuite) TestStartGame() {
	// given
	room := s.Client().createFilledRoom()

	// when
	room = s.Client().startGame(room.ID)

	// then
	s.True(room.IsStarted)
}

type RoomClient struct{ *handlertest.BaseClient }

func (c *RoomClient) createRoom() (room schema.Room) {
	c.ServeHTTPOkAs("POST", "/rooms", nil, &room)
	return
}

func (c *RoomClient) joinRoom(roomID uuid.UUID) (room schema.Room) {
	c.ServeHTTPOkAs("PUT", roomURL(roomID)+"/players", nil, &room)
	return
}

func (c *RoomClient) getRoom(roomID uuid.UUID) (room schema.Room) {
	c.ServeHTTPOkAs("GET", roomURL(roomID), nil, &room)
	return
}

func (c *RoomClient) createFilledRoom() (room schema.Room) {
	room = c.createRoom()
	for !room.IsStartable {
		c2 := handlertest.CloneWithEmptyJar(c)
		room = c2.joinRoom(room.ID)
	}
	return
}

func (c *RoomClient) startGame(roomID uuid.UUID) (room schema.Room) {
	c.ServeHTTPOkAs("PUT", roomURL(roomID)+"/game", nil, &room)
	return
}

func roomURL(roomID uuid.UUID) string {
	return "/rooms/" + roomID.String()
}

func TestRoomSuite(t *testing.T) {
	suite.Run(t, new(RoomSuite))
}
