package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/configs/serverconfig"
	"github.com/jostrzol/mess/pkg/mess"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/jostrzol/mess/pkg/server/core/game"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/ioc"
)

const GameURL = "/rooms/:id/game"

type GameHandler struct {
	service *game.Service        `container:"type"`
	config  *serverconfig.Config `container:"type"`
}

func GetGameStaticData(h *GameHandler, g *gin.Engine) {
	g.GET(GameURL+"/static", func(c *gin.Context) {
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
	g.GET(GameURL+"/state", func(c *gin.Context) {
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
	g.GET(GameURL+"/options", func(c *gin.Context) {
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
	g.PUT(GameURL+"/turns/:turn", func(c *gin.Context) {
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

		state, err = h.service.PlayTurn(session.ID, roomID, turn, route)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, schema.StateFromDomain(session.ID, state))
	})
}

func GetResolution(h *GameHandler, g *gin.Engine) {
	g.GET(GameURL+"/resolution", func(c *gin.Context) {
		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		resolution, err := h.service.GetResolution(roomID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		session := GetSessionData(sessions.Default(c))
		c.JSON(http.StatusOK, schema.ResolutionFromDomain(session.ID, resolution))
	})
}

func GetAsset(h *GameHandler, g *gin.Engine) {
	g.GET(GameURL+"/assets/*key", func(c *gin.Context) {
		roomID, err := parseUUID[id.Room](c.Param("id"))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		key := c.Param("key")

		data, err := h.service.GetAsset(roomID, mess.NewAssetKey(key))
		if err != nil {
			AbortWithError(c, err)
			return
		}

		mimetype := mimetype.Detect(data)

		c.Header("Cache-Control", fmt.Sprintf("max-age=%v", h.config.AssetsCacheMaxAge))
		c.Data(http.StatusOK, mimetype.String(), data)
	})
}

func init() {
	ioc.MustHandlerFill[GameHandler](
		GetGameStaticData,
		GetGameState,
		GetTurnOptions,
		PlayTurn,
		GetResolution,
		GetAsset,
	)
}
