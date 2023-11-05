package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/adapter/handler/handlertest"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/stretchr/testify/suite"
)

type JSON = handlertest.JSON
type RoomSuite struct {
	handlertest.HandlerSuite[RoomClient]
}

func (s *RoomSuite) TestCreateRoom() {
	s.Client().createRoom()
}

func (s *RoomSuite) TestJoinRoomSameClient() {
	room := s.Client().createRoom()
	room = s.Client().joinRoom(room.ID)
	s.Equal(1, room.Players)
}

func (s *RoomSuite) TestJoinRoomDifferentClient() {
	room := s.Client().createRoom()

	c2 := s.NewClient()
	room = c2.joinRoom(room.ID)
	s.Equal(2, room.Players)
}

func (s *RoomSuite) TestGetRoom() {
	room := s.Client().createRoom()
	s.Client().getRoom(room.ID)
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

func roomURL(roomID uuid.UUID) string {
	return "/rooms/" + roomID.String()
}

func TestRoomSuite(t *testing.T) {
	suite.Run(t, new(RoomSuite))
}
