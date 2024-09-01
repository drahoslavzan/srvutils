package logging

import (
	"time"

	"github.com/drahoslavzan/srvutils/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	var config zap.Config
	if env.IsDevelopment() {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger := zap.Must(config.Build())

	zap.ReplaceGlobals(logger)

	return logger
}
