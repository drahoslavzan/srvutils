package logging

import (
	"github.com/ThreeDotsLabs/watermill"
	"go.uber.org/zap"
)

type (
	wmZapLoggerAdapter struct {
		logger *zap.Logger
	}
)

func WatermillFromZap(logger *zap.Logger) watermill.LoggerAdapter {
	return wmZapLoggerAdapter{logger}
}

func (m wmZapLoggerAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return WatermillFromZap(m.logger.With(m.fields(fields)...))
}

func (m wmZapLoggerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	m.logger.With(m.fields(fields)...).Error(msg, zap.Error(err))
}

func (m wmZapLoggerAdapter) Info(msg string, fields watermill.LogFields) {
	m.logger.With(m.fields(fields)...).Info(msg)
}

func (m wmZapLoggerAdapter) Debug(msg string, fields watermill.LogFields) {
	m.logger.With(m.fields(fields)...).Debug(msg)
}

func (m wmZapLoggerAdapter) Trace(msg string, fields watermill.LogFields) {
	m.logger.With(m.fields(fields)...).Debug(msg)
}

func (m wmZapLoggerAdapter) fields(fields watermill.LogFields) []zap.Field {
	ret := []zap.Field{}
	for k, v := range fields {
		ret = append(ret, zap.Any(k, v))
	}

	return ret
}
