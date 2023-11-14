package handler

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/core/room"
	"github.com/jostrzol/mess/pkg/server/ioc"
)

type RoomHandler struct {
	service   *room.Service `container:"type"`
	wsHandler *WsHandler    `container:"type"`
}

func CreateRoom(h *RoomHandler, g *gin.Engine) {
	g.POST("/rooms", func(c *gin.Context) {
		session := GetSessionData(sessions.Default(c))

		room, err := h.service.CreateRoom(session.ID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, schema.RoomFromDomain(room))
	})
}

func GetRoom(h *RoomHandler, g *gin.Engine) {
	g.GET("/rooms/:id", func(c *gin.Context) {
		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		room, err := h.service.GetRoom(roomID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, schema.RoomFromDomain(room))
	})
}

func JoinRoom(h *RoomHandler, g *gin.Engine) {
	g.PUT("/rooms/:id/players", func(c *gin.Context) {
		session := GetSessionData(sessions.Default(c))

		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		r, err := h.service.JoinRoom(session.ID, roomID)
		switch {
		case errors.Is(err, room.ErrAlreadyInRoom):
			break
		case err != nil:
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, schema.RoomFromDomain(r))
	})
}

func StartGame(h *RoomHandler, g *gin.Engine) {
	g.PUT("/rooms/:id/game", func(c *gin.Context) {
		session := GetSessionData(sessions.Default(c))

		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		r, err := h.service.StartGame(session.ID, roomID)
		switch {
		case errors.Is(err, room.ErrAlreadyStarted):
			break
		case err != nil:
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, schema.RoomFromDomain(r))
	})
}

func HandleWebsocket(h *RoomHandler, g *gin.Engine) {
	g.GET("/websocket", func(c *gin.Context) {
		err := h.wsHandler.handle(c)
		if err != nil {
			AbortWithError(c, err)
			return
		}
	})
}

func init() {
	ioc.MustHandlerFill[RoomHandler](
		CreateRoom,
		GetRoom,
		JoinRoom,
		StartGame,
		HandleWebsocket,
	)
}
