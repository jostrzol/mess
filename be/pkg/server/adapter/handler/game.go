package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/jostrzol/mess/pkg/server/core/game"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/ioc"
)

type GameHandler struct {
	service   *game.Service `container:"type"`
	wsHandler *WsHandler    `container:"type"`
}

func GetGameState(h *GameHandler, g *gin.Engine) {
	g.GET("/rooms/:id/game", func(c *gin.Context) {
		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		session := GetSessionData(sessions.Default(c))
		state, err := h.service.GetGameState(session.ID, roomID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, schema.StateFromDomain(state))
	})
}

func PlayTurn(h *GameHandler, g *gin.Engine) {
	g.PUT("/rooms/:id/game/turns/:turn", func(c *gin.Context) {
		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		turn, err := strconv.Atoi(c.Param("turn"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		session := GetSessionData(sessions.Default(c))

		state, err := h.service.GetGameState(session.ID, roomID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		var routeDto schema.Route
		err = c.ShouldBindWith(&routeDto, schema.RouteBinding{})
		if err != nil {
			AbortWithError(c, err)
			return
		}

		route, err := routeDto.ToDomain(state.State)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		fmt.Printf("route: %#v\n", route)

		state, err = h.service.PlayTurn(session.ID, roomID, turn, route)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, schema.StateFromDomain(state))
	})
}

func init() {
	ioc.MustHandlerFill[GameHandler](
		GetGameState,
		PlayTurn,
	)
}
