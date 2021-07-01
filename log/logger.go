package log

import (
	"path"
	"runtime"
	"strings"

	"github.com/drahoslavzan/srvutils/env"
	"github.com/sirupsen/logrus"
)

type LoggerOpts struct {
	FuncName string
	Fields   map[string]interface{}
}

type serviceLogger struct {
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
		logrus.SetLevel(logrus.WarnLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func GetLogger(opts ...LoggerOpts) *serviceLogger {
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

	return &serviceLogger{logrus.WithFields(fs)}
}

func getFields(funcName string) logrus.Fields {
	ci := retrieveCallInfo()

	if len(funcName) < 1 {
		funcName = ci.funcName
	}

	return logrus.Fields{
		"function": funcName,
		"package":  ci.packageName,
		"service":  env.GetEnv("SERVICE_NAME"),
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
