package main

import (
	"errors"
	"log"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/golobby/container/v3"
	"github.com/jostrzol/mess/configs/serverconfig"
	"github.com/jostrzol/mess/pkg/logger"
	_ "github.com/jostrzol/mess/pkg/server/adapter/http"
	_ "github.com/jostrzol/mess/pkg/server/adapter/inmem"
	"github.com/jostrzol/mess/pkg/server/ioc"
	"go.uber.org/zap"
)

func main() {
	config, logger := loadConfigAndLogger()

	if !config.IsProduction {
		logger.Info("configuration loaded", zap.Any("config", config))
	}

	mode := gin.DebugMode
	if config.IsProduction {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	g := gin.New()
	g.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	g.Use(ginzap.RecoveryWithZap(logger, true))

	container.MustSingleton(container.Global, func() *gin.Engine { return g })
	for _, initializer := range ioc.HandlerInitializers {
		initializer(g)
	}
	log.Fatal(g.Run(":4000"))
}

func loadConfigAndLogger() (*serverconfig.Config, *zap.Logger) {
	config, errConfig := serverconfig.New()
	isProduction := false
	if config != nil {
		isProduction = config.IsProduction
	}
	logger, err := logger.New(isProduction)
	if err != nil && logger == nil {
		log.Fatal(err)
	} else if err != nil {
		logger.Fatal("", zap.Error(err))
	} else if errors.Is(errConfig, serverconfig.ErrConfigFileNotFound) {
		logger.Warn("", zap.Error(errConfig))
	} else if errConfig != nil {
		logger.Fatal("", zap.Error(errConfig))
	}

	return config, logger
}
