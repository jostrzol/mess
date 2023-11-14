package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/golobby/container/v3"
	"github.com/jostrzol/mess/configs/serverconfig"
	"github.com/jostrzol/mess/pkg/server/ioc"
	"go.uber.org/zap"
)

func init() {
	container.MustSingletonLazy(container.Global, func(
		logger *zap.Logger,
		config *serverconfig.Config,
	) *gin.Engine {
		mode := gin.DebugMode
		if config.IsProduction {
			mode = gin.ReleaseMode
		}
		gin.SetMode(mode)

		g := gin.New()
		c := cors.DefaultConfig()
		c.AllowOrigins = []string{config.IncomingOrigin}
		c.AllowCredentials = true
		g.Use(cors.New(c))
		g.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		g.Use(ginzap.RecoveryWithZap(logger, true))
		store := memstore.NewStore([]byte(config.SessionSecret))
		g.Use(sessions.Sessions(SessionKey, store))

		for _, initializer := range ioc.HandlerInitializers {
			initializer(g)
		}

		return g
	})
}
