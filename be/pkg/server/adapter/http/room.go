package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/core/room"
	"github.com/jostrzol/mess/pkg/server/ioc"
)

type RoomHandler struct {
	service *room.Service `container:"type"`
}

func CreateRoom(h *RoomHandler, g *gin.Engine) {
	g.POST("/rooms", func(c *gin.Context) {
		var params struct {
			PlayerID string `form:"player_id" binding:"uuid"`
		}
		err := c.BindQuery(&params)
		if err != nil {
			return
		}
		playerID, err := uuid.Parse(params.PlayerID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		room, err := h.service.CreateRoom(playerID)
		if err != nil {
			AbortWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, room)
	})
}

func init() {
	ioc.MustHandler[RoomHandler](CreateRoom)
}
