package log

import (
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

type callInfo struct {
	packageName string
	fileName    string
	funcName    string
	line        int
}

func GetLogger() *log.Entry {
	ci := retrieveCallInfo()

	return log.WithFields(log.Fields{
		"package":  ci.packageName,
		"function": ci.funcName,
	})
}

func GetLambdaLogger(funcName string) *log.Entry {
	ci := retrieveCallInfo()

	return log.WithFields(log.Fields{
		"package":  ci.packageName,
		"function": funcName,
	})
}

func retrieveCallInfo() *callInfo {
	pc, file, line, _ := runtime.Caller(2)
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
