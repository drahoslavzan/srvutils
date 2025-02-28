package logging

import (
	"github.com/ThreeDotsLabs/watermill"
	"go.uber.org/zap"
)

type (
	WMZapLoggerAdapter struct {
		logger *zap.Logger
	}
)

func NewWatermillZapLogger(logger *zap.Logger) *WMZapLoggerAdapter {
	return &WMZapLoggerAdapter{}
}

func (m *WMZapLoggerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	m.logger.With(m.fields(fields)...).Error(msg, zap.Error(err))
}

func (m *WMZapLoggerAdapter) Info(msg string, fields watermill.LogFields) {
	m.logger.With(m.fields(fields)...).Info(msg)
}

func (m *WMZapLoggerAdapter) Debug(msg string, fields watermill.LogFields) {
	m.logger.With(m.fields(fields)...).Debug(msg)
}

func (m *WMZapLoggerAdapter) Trace(msg string, fields watermill.LogFields) {
	m.logger.With(m.fields(fields)...).Debug(msg)
}

func (m *WMZapLoggerAdapter) fields(fields watermill.LogFields) []zap.Field {
	ret := []zap.Field{}
	for k, v := range fields {
		ret = append(ret, zap.Any(k, v))
	}

	return ret
}
