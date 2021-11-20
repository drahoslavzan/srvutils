package log

import (
	"path"
	"runtime"
	"strings"

	"github.com/drahoslavzan/srvutils/env"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Fatal(err error)
	Panic(err error)
	Error(err error)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)

	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type LoggerOpts struct {
	FuncName string
	Fields   map[string]interface{}
}

type logger struct {
	*logrus.Entry
}

type callInfo struct {
	packageName string
	fileName    string
	funcName    string
	line        int
}

func init() {
	if env.IsProduction() {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func GetLogger(opts ...LoggerOpts) Logger {
	var o LoggerOpts
	if len(opts) > 0 {
		o = opts[0]
	} else {
		o = LoggerOpts{}
	}

	fs := getFields(o.FuncName)
	for k, v := range o.Fields {
		fs[k] = v
	}

	return &logger{logrus.WithFields(fs)}
}

func (m *logger) Fatal(err error) {
	m.Entry.Fatal(err)
}

func (m *logger) Panic(err error) {
	m.Entry.Panic(err)
}

func (m *logger) Error(err error) {
	m.Entry.Error(err)
}

func (m *logger) Warn(msg string) {
	m.Entry.Warn(msg)
}

func (m *logger) Info(msg string) {
	m.Entry.Info(msg)
}

func (m *logger) Debug(msg string) {
	m.Entry.Debug(msg)
}

func (m *logger) Fatalf(format string, args ...interface{}) {
	m.Entry.Fatalf(format, args...)
}

func (m *logger) Panicf(format string, args ...interface{}) {
	m.Entry.Panicf(format, args...)
}

func (m *logger) Errorf(format string, args ...interface{}) {
	m.Entry.Errorf(format, args...)
}

func (m *logger) Warnf(format string, args ...interface{}) {
	m.Entry.Warnf(format, args...)
}

func (m *logger) Infof(format string, args ...interface{}) {
	m.Entry.Infof(format, args...)
}

func (m *logger) Debugf(format string, args ...interface{}) {
	m.Entry.Debugf(format, args...)
}

func getFields(funcName string) logrus.Fields {
	ci := retrieveCallInfo()

	if len(funcName) < 1 {
		funcName = ci.funcName
	}

	return logrus.Fields{
		"function": funcName,
		"package":  ci.packageName,
	}
}

func retrieveCallInfo() *callInfo {
	pc, file, line, _ := runtime.Caller(3)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return &callInfo{
		packageName: packageName,
		fileName:    fileName,
		funcName:    funcName,
		line:        line,
	}
}
