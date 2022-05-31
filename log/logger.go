package log

import (
	"github.com/drahoslavzan/srvutils/env"
	"go.uber.org/zap"
)

type Logger interface {
	Panic(err error)
	Error(err error)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)

	Panicw(err error, keysAndValues ...interface{})
	Errorw(err error, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}

type logger struct {
	z *zap.SugaredLogger
}

func NewLogger(keysAndValues ...interface{}) Logger {
	var err error
	var zl *zap.Logger
	if env.IsProduction() {
		zl, err = zap.NewProduction()
	} else {
		zl, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	z := zl.Sugar().With(keysAndValues...)
	return &logger{z}
}

func (m *logger) Panic(err error) {
	m.z.Panic(err)
}

func (m *logger) Error(err error) {
	m.z.Error(err)
}

func (m *logger) Warn(msg string) {
	m.z.Warn(msg)
}

func (m *logger) Info(msg string) {
	m.z.Info(msg)
}

func (m *logger) Debug(msg string) {
	m.z.Debug(msg)
}

func (m *logger) Panicw(err error, keysAndValues ...interface{}) {
	m.z.Panicw(err.Error(), keysAndValues...)
}

func (m *logger) Errorw(err error, keysAndValues ...interface{}) {
	m.z.Errorw(err.Error(), keysAndValues...)
}

func (m *logger) Warnw(msg string, keysAndValues ...interface{}) {
	m.z.Warnw(msg, keysAndValues...)
}

func (m *logger) Infow(msg string, keysAndValues ...interface{}) {
	m.z.Infow(msg, keysAndValues...)
}

func (m *logger) Debugw(msg string, keysAndValues ...interface{}) {
	m.z.Debugw(msg, keysAndValues...)
}
