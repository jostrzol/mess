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

func (s *RoomSuite) TestGetRules() {
	// given
	room := s.Client().createRoom()

	// expect
	s.Client().getRules(room.ID)
}

func (s *RoomSuite) TestSetRules() {
	// given
	room := s.Client().createRoom()

	// when
	s.Client().setRules(room.ID, "rules.hcl", "board { width = 2; height = 2 }")

	// then
	rules := s.Client().getRules(room.ID)
	s.Equal(rules, "board { width = 2; height = 2 }")

	// and
	room = s.Client().getRoom(room.ID)
	s.Equal(room.RulesFilename, "rules.hcl")
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
	c.ServeJSONOkAs("POST", "/rooms", nil, &room)
	return
}

func (c *RoomClient) joinRoom(roomID uuid.UUID) (room schema.Room) {
	c.ServeJSONOkAs("PUT", roomURL(roomID)+"/players", nil, &room)
	return
}

func (c *RoomClient) getRoom(roomID uuid.UUID) (room schema.Room) {
	c.ServeJSONOkAs("GET", roomURL(roomID), nil, &room)
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

func (c *RoomClient) getRules(roomID uuid.UUID) string {
	res := c.ServeJSONOk("GET", roomURL(roomID)+"/rules", nil)
	return res.Body.String()
}

func (c *RoomClient) setRules(roomID uuid.UUID, filename string, data string) {
	c.ServeOk("PUT", roomURL(roomID)+"/rules/"+filename, []byte(data))
}

func (c *RoomClient) startGame(roomID uuid.UUID) (room schema.Room) {
	c.ServeJSONOkAs("PUT", roomURL(roomID)+"/game", nil, &room)
	return
}

func roomURL(roomID uuid.UUID) string {
	return "/rooms/" + roomID.String()
}

func TestRoomSuite(t *testing.T) {
	suite.Run(t, new(RoomSuite))
}
