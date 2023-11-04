package http

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/pkg/server/adapter/httpschema"
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

		h.wsHandler.sendToOpponents(room, session.ID, &httpschema.RoomChanged{})

		c.JSON(http.StatusOK, httpschema.NewRoom(room))
	})
}

func StartGame(h *RoomHandler, g *gin.Engine) {
	g.PUT("/rooms/:id/game", func(c *gin.Context) {
		session := GetSessionData(sessions.Default(c))
		id := c.Param("id")

		roomID, err := parseUUID(id)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		room, err := h.service.StartGame(roomID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		h.wsHandler.sendToOpponents(room, session.ID, &httpschema.GameStarted{})

		c.JSON(http.StatusOK, httpschema.NewRoom(room))
	})
}

func HandleWebsocket(h *RoomHandler, g *gin.Engine) {
	g.GET("/rooms/:id/websocket", func(c *gin.Context) {
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
		HandleWebsocket,
		StartGame,
	)
}
