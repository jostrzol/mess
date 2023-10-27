package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(isProduction bool) (*zap.Logger, error) {
	var config zap.Config
	if isProduction {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.DisableStacktrace = true
	return config.Build()
}
