package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/pkg/server/adapter/httpschema"
	"github.com/jostrzol/mess/pkg/server/core/room"
	"github.com/jostrzol/mess/pkg/server/ioc"
)

type RoomHandler struct {
	service *room.Service `container:"type"`
}

func CreateRoom(h *RoomHandler, g *gin.Engine) {
	g.POST("/rooms", func(c *gin.Context) {
		session := GetSessionData(c)

		room, err := h.service.CreateRoom(session.ID)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, httpschema.NewRoom(room))
	})
}

func init() {
	ioc.MustHandlerFill[RoomHandler](CreateRoom)
}
