package log

import "go.uber.org/zap"

func Service(s string) any {
	return zap.String("service", s)
}

func Module(m string) any {
	return zap.String("module", m)
}
