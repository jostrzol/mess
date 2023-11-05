package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/configs/serverconfig"
	"github.com/jostrzol/mess/pkg/logger"
	_ "github.com/jostrzol/mess/pkg/server/adapter/handler"
	_ "github.com/jostrzol/mess/pkg/server/adapter/inmem"
	"github.com/jostrzol/mess/pkg/server/ioc"
	"go.uber.org/zap"
)

func main() {
	config := loadConfig()

	g := ioc.MustResolve[*gin.Engine]()

	address := fmt.Sprintf(":%d", config.Port)
	log.Fatal(g.Run(address))
}

func loadConfig() *serverconfig.Config {
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

	if !config.IsProduction {
		logger.Info("configuration loaded", zap.Any("config", config))
	}

	ioc.MustSingleton(config)
	ioc.MustSingleton(logger)
	ioc.MustSingleton(logger.Sugar())

	return config
}
