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

	Panicw(err error, keysAndValues ...any)
	Errorw(err error, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Debugw(msg string, keysAndValues ...any)
}

type logger struct {
	z *zap.SugaredLogger
}

func NewLogger(keysAndValues ...any) Logger {
	var err error
	var zl *zap.Logger

	optClrSkip := zap.AddCallerSkip(1)
	if env.IsProduction() {
		zl, err = zap.NewProduction(optClrSkip)
	} else {
		zl, err = zap.NewDevelopment(optClrSkip)
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

func (m *logger) Panicw(err error, keysAndValues ...any) {
	m.z.Panicw(err.Error(), keysAndValues...)
}

func (m *logger) Errorw(err error, keysAndValues ...any) {
	m.z.Errorw(err.Error(), keysAndValues...)
}

func (m *logger) Warnw(msg string, keysAndValues ...any) {
	m.z.Warnw(msg, keysAndValues...)
}

func (m *logger) Infow(msg string, keysAndValues ...any) {
	m.z.Infow(msg, keysAndValues...)
}

func (m *logger) Debugw(msg string, keysAndValues ...any) {
	m.z.Debugw(msg, keysAndValues...)
}
