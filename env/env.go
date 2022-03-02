package env

import (
	"fmt"
	"os"
	"strconv"
)

func IsProduction() bool {
	return os.Getenv("PRODUCTION") == "true"
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
		panic(fmt.Errorf("missing env value: %s", key))
	}
	return val
}

func GetEnvOpt(key string) *string {
	val := os.Getenv(key)
	if len(val) < 1 {
		return nil
	}
	return &val
}

func GetIntEnv(key string) int {
	val := GetEnv(key)
	num, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("invalid env value for '%s': %v", key, val))
	}
	return num
}

func GetIntEnvOpt(key string) *int {
	if GetEnvOpt(key) == nil {
		return nil
	}
	val := GetIntEnv(key)
	return &val
}
