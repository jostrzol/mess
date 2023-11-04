package handler

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jostrzol/mess/pkg/server/adapter/httpschema"
	"github.com/jostrzol/mess/pkg/server/adapter/inmem"
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
}

func init() {
	ioc.MustSingletonFill[WsHandler]()
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
		bytes, err := httpschema.MarshalEvent(event)
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

func (h *WsHandler) sendToOpponents(room *room.Room, mySessionID uuid.UUID, event httpschema.Event) {
	for _, playerID := range room.Players() {
		if playerID != mySessionID {
			err := h.websockets.Send(playerID, event)
			if err != nil {
				h.logger.Error("sending event",
					zap.Stringer("target", playerID),
					zap.String("eventType", event.EventType()),
					zap.Error(err),
				)
			}
		}
	}
}
