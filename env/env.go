package env

import (
	"fmt"
	"os"
	"strconv"
)

func IsProduction() bool {
	isProd := os.Getenv("PRODUCTION")
	return isProd == "true"
}

func IsDevelopment() bool {
	return !IsProduction()
}

func EnvType() string {
	if IsProduction() {
		return "production"
	}

	return "development"
}

func GetEnv(key string) string {
	val := os.Getenv(key)
	if len(val) < 1 {
		panic(fmt.Errorf("missing env value for '%s'", key))
	}
	return val
}

func GetIntEnv(key string) int {
	val := GetEnv(key)
	num, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("invalid env value for '%s': %v", key, val))
	}
	return num
}
