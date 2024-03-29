package handler_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/adapter/handler/handlertest"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/stretchr/testify/suite"
)

type GameSuite struct {
	handlertest.HandlerSuite[GameClient]
}

func (s *GameSuite) TestGetGameStaticData() {
	// given
	room := s.Client().createStartedRoom()

	// when
	staticData := s.Client().getStaticData(room.ID)

	// expect
	s.NotZero(staticData)
}

func (s *GameSuite) TestGetGameState() {
	// given
	room := s.Client().createFilledRoom()

	// and
	room = s.Client().startGame(room.ID)

	// expect
	s.Client().getGameState(room.ID)
}

func (s *GameSuite) TestGetTurnOptions() {
	// given
	room := s.Client().createFilledRoom()

	// and
	room = s.Client().startGame(room.ID)

	// expect
	s.Client().getTurnOptions(room.ID)
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

func (s *GameSuite) TestGetResolution() {
	// given
	room := s.Client().createStartedRoom()

	// when
	resolution := s.Client().getResolution(room.ID)

	// expect
	s.NotZero(resolution)
}

func (s *GameSuite) TestGetAsset() {
	// given
	room := s.Client().createStartedRoom()

	// expect
	s.Client().getAsset(room.ID, "/piece_types/king.svg")
}

type GameClient struct{ RoomClient }

func (c *GameClient) getStaticData(roomID uuid.UUID) (staticData schema.StaticData) {
	c.ServeJSONOkAs("GET", roomURL(roomID)+"/game/static", nil, &staticData)
	return
}

func (c *GameClient) getGameState(roomID uuid.UUID) (state schema.State) {
	c.ServeJSONOkAs("GET", roomURL(roomID)+"/game/state", nil, &state)
	return
}

type OptionNode struct {
	Type    string
	Message string
	Data    []struct {
		Option   any
		Children []*OptionNode
	}
}

func (c *GameClient) createStartedRoom() (room schema.Room) {
	room = c.createFilledRoom()
	room = c.startGame(room.ID)
	return
}

func (c *GameClient) getTurnOptions(roomID uuid.UUID) (optionTree *OptionNode) {
	c.ServeJSONOkAs("GET", roomURL(roomID)+"/game/options", nil, &optionTree)
	return
}

func (c *GameClient) chooseTurnOpionRoute(roomID uuid.UUID, turn int, route any) (state schema.State) {
	c.ServeJSONOkAs("PUT", roomURL(roomID)+fmt.Sprintf("/game/turns/%v", turn), route, &state)
	return
}

func (c *GameClient) getResolution(roomID uuid.UUID) (resolution schema.Resolution) {
	c.ServeJSONOkAs("GET", roomURL(roomID)+"/game/resolution", nil, &resolution)
	return
}

func (c *GameClient) getAsset(roomID uuid.UUID, assetKey string) []byte {
	res := c.ServeJSONOk("GET", roomURL(roomID)+"/game/assets"+assetKey, nil)
	bytes, err := io.ReadAll(res.Body)
	c.NoError(err)
	return bytes
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}
