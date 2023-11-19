package handler

import (
	"errors"
	"io"
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

func GetRules(h *RoomHandler, g *gin.Engine) {
	g.GET("/rooms/:id/rules", func(c *gin.Context) {
		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		rules, err := h.service.GetRules(roomID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.Data(http.StatusOK, HclContent, rules.Src)
	})
}

func SetRules(h *RoomHandler, g *gin.Engine) {
	g.PUT("/rooms/:id/rules/:filename", func(c *gin.Context) {
		session := GetSessionData(sessions.Default(c))

		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		filename := c.Param("filename")
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		err = h.service.SetRules(session.ID, roomID, filename, data)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
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
		GetRules,
		SetRules,
		StartGame,
		HandleWebsocket,
	)
}
