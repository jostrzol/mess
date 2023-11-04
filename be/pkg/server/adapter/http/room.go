package http

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jostrzol/mess/pkg/server/adapter/httpschema"
	"github.com/jostrzol/mess/pkg/server/adapter/inmem"
	"github.com/jostrzol/mess/pkg/server/adapter/ws"
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

type RoomHandler struct {
	service    *room.Service       `container:"type"`
	logger     *zap.Logger         `container:"type"`
	websockets *inmem.WsRepository `container:"type"`
}

func CreateRoom(h *RoomHandler, g *gin.Engine) {
	g.POST("/rooms", func(c *gin.Context) {
		session := GetSessionData(sessions.Default(c))

		room, err := h.service.CreateRoom(session.ID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, httpschema.NewRoom(room))
	})
}

func GetRoom(h *RoomHandler, g *gin.Engine) {
	g.GET("/rooms/:id", func(c *gin.Context) {
		id := c.Param("id")

		roomID, err := parseUUID(id)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		room, err := h.service.GetRoom(roomID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, httpschema.NewRoom(room))
	})
}

func JoinRoom(h *RoomHandler, g *gin.Engine) {
	g.PUT("/rooms/:id/players", func(c *gin.Context) {
		session := GetSessionData(sessions.Default(c))
		id := c.Param("id")

		roomID, err := parseUUID(id)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		room, err := h.service.JoinRoom(session.ID, roomID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		for _, playerID := range room.Players() {
			if playerID != session.ID {
				err := h.websockets.Send(playerID, &ws.RoomChanged{})
				if err != nil {
					h.logger.Error("sending room changed event",
						zap.Stringer("target", playerID),
						zap.Error(err))
					return
				}
			}
		}

		c.JSON(http.StatusOK, httpschema.NewRoom(room))
	})
}

func HandleWebsocket(h *RoomHandler, g *gin.Engine) {
	g.GET("/rooms/:id/websocket", func(c *gin.Context) {
		session := GetSessionData(sessions.Default(c))
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			AbortWithError(c, err)
			return
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
			bytes, err := ws.Marshal(event)
			if err != nil {
				h.logger.Error("marshaling websocket message", zap.Error(err))
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				h.logger.Error("writing websocket message", zap.Error(err))
				return
			}
		}
	})
}

func init() {
	ioc.MustHandlerFill[RoomHandler](
		CreateRoom,
		GetRoom,
		JoinRoom,
		HandleWebsocket,
	)
}
