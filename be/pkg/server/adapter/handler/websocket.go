package handler

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jostrzol/mess/pkg/server/adapter/inmem"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/jostrzol/mess/pkg/server/core/event"
	"github.com/jostrzol/mess/pkg/server/core/game"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/core/room"
	"github.com/jostrzol/mess/pkg/server/ioc"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     websocket.IsWebSocketUpgrade,
}

const wsTimeout = time.Second

type WsHandler struct {
	logger     *zap.Logger         `container:"type"`
	websockets *inmem.WsRepository `container:"type"`
	rooms      room.Repository     `container:"type"`
	games      game.Repository     `container:"type"`
}

func init() {
	ioc.MustSingletonObserverFill[WsHandler]()
}

func (h *WsHandler) handle(c *gin.Context) error {
	session := GetSessionData(sessions.Default(c))
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return err
	}
	defer func() {
		defer conn.Close()
		h.logger.Info("closing websocket connection", zap.Stringer("session", session.ID))
		err = conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(wsTimeout))
		if err != nil {
			h.logger.Error("sending websocket close message", zap.Error(err))
			return
		}
	}()

	channel := h.websockets.New(session.ID)
	conn.SetCloseHandler(func(code int, text string) error {
		h.logger.Info("peer closed session channel, cleaning up", zap.Stringer("session", session.ID))
		h.websockets.Close(session.ID)
		return nil
	})

	for event, ok := <-channel; ok; event, ok = <-channel {
		bytes, err := schema.MarshalEvent(event)
		if err != nil {
			h.logger.Error("marshaling websocket message", zap.Error(err))
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			h.logger.Error("writing websocket message", zap.Error(err))
			continue
		}
	}
	return nil
}

func (h *WsHandler) Handle(evnt event.Event) {
	var err error
	var eventToSend schema.Event
	var players []id.Session
	var author id.Session
	switch ev := evnt.(type) {
	case *event.PlayerJoined:
		players, err = h.playersInRoom(ev.RoomID)
		author = ev.PlayerID
		eventToSend = &schema.RoomChanged{}
	case *event.RoomRulesChanged:
		players, err = h.playersInRoom(ev.RoomID)
		author = ev.By
		eventToSend = &schema.RoomChanged{}
	case *event.GameStarted:
		players, err = h.playersInRoom(ev.RoomID)
		author = ev.By
		eventToSend = &schema.RoomChanged{}
	case *event.GameChanged:
		players, err = h.playersInGame(ev.GameID)
		author = ev.By
		eventToSend = &schema.GameChanged{}
	}
	if err != nil {
		h.logger.Error("sending event", zap.Any("event", evnt), zap.Error(err))
		return
	}
	h.sendToOpponents(players, author, eventToSend)
}

func (h *WsHandler) sendToOpponents(players []id.Session, author id.Session, event schema.Event) {
	for _, player := range players {
		if player != author {
			h.logger.Debug("sending websocket message", zap.Any("event", event), zap.Stringer("session", player))
			err := h.websockets.Send(player, event)
			if err != nil {
				h.logger.Error("sending event",
					zap.Stringer("target", player),
					zap.String("eventType", event.EventType()),
					zap.Error(err),
				)
			}
		}
	}
}

func (h *WsHandler) playersInRoom(roomID id.Room) ([]id.Session, error) {
	room, err := h.rooms.Get(roomID)
	if err != nil {
		return nil, err
	}
	return room.Players(), nil
}

func (h *WsHandler) playersInGame(gameID id.Game) ([]id.Session, error) {
	game, err := h.games.Get(gameID)
	if err != nil {
		return nil, err
	}
	return game.Players(), nil
}
