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
	service *game.Service `container:"type"`
}

func GetGameStaticData(h *GameHandler, g *gin.Engine) {
	g.GET("/rooms/:id/game/static", func(c *gin.Context) {
		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		session := GetSessionData(sessions.Default(c))
		staticData, err := h.service.GetGameStaticData(session.ID, roomID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, schema.StaticDataFromDomain(staticData))
	})
}

func GetGameState(h *GameHandler, g *gin.Engine) {
	g.GET("/rooms/:id/game/state", func(c *gin.Context) {
		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		session := GetSessionData(sessions.Default(c))
		state, err := h.service.GetGameState(roomID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, schema.StateFromDomain(session.ID, state))
	})
}

func GetTurnOptions(h *GameHandler, g *gin.Engine) {
	g.GET("/rooms/:id/game/options", func(c *gin.Context) {
		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		optionTree, err := h.service.GetTurnOptions(roomID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, schema.OptionNodeFromDomain(optionTree))
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

		state, err := h.service.GetGameState(roomID)
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

		route, err := routeDto.ToDomain(state)
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

		c.JSON(http.StatusOK, schema.StateFromDomain(session.ID, state))
	})
}

func init() {
	ioc.MustHandlerFill[GameHandler](
		GetGameStaticData,
		GetGameState,
		GetTurnOptions,
		PlayTurn,
	)
}
