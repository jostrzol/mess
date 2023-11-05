package handler

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/jostrzol/mess/pkg/server/core/room"
	"github.com/jostrzol/mess/pkg/server/ioc"
)

type GameHandler struct {
	service   *room.Service `container:"type"`
	wsHandler *WsHandler    `container:"type"`
}

func StartGame(h *GameHandler, g *gin.Engine) {
	g.PUT("/rooms/:id/game", func(c *gin.Context) {
		session := GetSessionData(sessions.Default(c))
		id := c.Param("id")

		roomID, err := parseUUID(id)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		r, err := h.service.StartGame(roomID)
		switch {
		case errors.Is(err, room.ErrAlreadyStarted):
			break
		case err != nil:
			AbortWithError(c, err)
			return
		default:
			h.wsHandler.sendToOpponents(r, session.ID, &schema.GameStarted{})
		}

		c.JSON(http.StatusOK, schema.RoomFromDomain(r))
	})
}

func GetGameState(h *GameHandler, g *gin.Engine) {
	g.GET("/rooms/:id/game", func(c *gin.Context) {
		id := c.Param("id")

		roomID, err := parseUUID(id)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		session := GetSessionData(sessions.Default(c))
		state, err := h.service.GetGameState(roomID, session.ID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, schema.StateFromDomain(state))
	})
}

func init() {
	ioc.MustHandlerFill[GameHandler](StartGame, GetGameState)
}
