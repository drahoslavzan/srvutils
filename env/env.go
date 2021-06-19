package env

import (
	"fmt"
	"os"
)

func IsProduction() bool {
	isProd := os.Getenv("PRODUCTION")
	return isProd == "true"
}

func GetEnv(key string) string {
	val := os.Getenv(key)
	if len(val) < 1 {
		panic(fmt.Errorf("missing env value for '%s'", key))
	}
	return val
}
