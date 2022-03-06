package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetEnv(key string) string {
	val := os.Getenv(key)
	if len(val) < 1 {
		panic(fmt.Errorf("missing env variable: %s", key))
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

func GetEnvDef(key string, def string) string {
	if v := GetEnvOpt(key); v == nil {
		return def
	} else {
		return *v
	}
}

func GetIntEnv(key string) int {
	val := GetEnv(key)
	num, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("invalid integer value for env variable %s: %v", key, val))
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

func GetIntEnvDef(key string, def int) int {
	if v := GetIntEnvOpt(key); v == nil {
		return def
	} else {
		return *v
	}
}

func GetEnvBool(key string) bool {
	v := GetEnv(key)
	return isTrue(v)
}

func GetEnvBoolDef(key string, def bool) bool {
	if v := GetEnvOpt(key); v == nil {
		return def
	} else {
		return isTrue(*v)
	}
}

func isTrue(v string) bool {
	return len(v) > 0 && v != "0" && strings.ToLower(v) != "false"
}
