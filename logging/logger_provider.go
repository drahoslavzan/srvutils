package logging

import (
	"github.com/drahoslavzan/srvutils/env"
	"go.uber.org/zap"
)

type LoggerProvider struct {
	loggers []*zap.Logger
}

func NewLoggerProvider() *LoggerProvider {
	return &LoggerProvider{}
}

func (m *LoggerProvider) Sync() {
	for _, l := range m.loggers {
		l.Sync()
	}
}

func (m *LoggerProvider) GetLogger() *zap.Logger {
	if len(m.loggers) > 0 {
		return m.loggers[0]
	}

	var logger *zap.Logger
	var err error
	if env.IsDevelopment() {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}

	m.loggers = append(m.loggers, logger)
	return logger
}

func (m *LoggerProvider) PkgLogger(pkg string) *zap.Logger {
	logger := m.GetLogger().With(zap.String("package", pkg))
	m.loggers = append(m.loggers, logger)
	return logger
}
