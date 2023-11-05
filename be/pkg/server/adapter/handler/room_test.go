package handler_test

import (
	"testing"

	"github.com/jostrzol/mess/pkg/server/adapter/handler/handlertest"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/stretchr/testify/suite"
)

type JSON = handlertest.JSON
type RoomSuite struct{ handlertest.HandlerSuite }

func (s *RoomSuite) TestCreateRoom() {
	s.createRoom()
}

func (s *RoomSuite) TestJoinRoom() {
	s.createRoom()
}

func (s *RoomSuite) createRoom() (room schema.Room) {
	s.ServeHTTPOkAs("POST", "/rooms", nil, &room)
	return
}

func TestRoomSuite(t *testing.T) {
	suite.Run(t, new(RoomSuite))
}
