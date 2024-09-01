package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func String(key string) string {
	val := os.Getenv(key)
	if len(val) < 1 {
		panic(fmt.Errorf("missing env variable: %s", key))
	}

	return val
}

func StringOpt(key string) *string {
	val := os.Getenv(key)
	if len(val) < 1 {
		return nil
	}

	return &val
}

func StringDef(key string, def string) string {
	if v := StringOpt(key); v != nil {
		return *v
	}

	return def
}

func Int(key string) int {
	val := String(key)
	num, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("invalid integer value for env variable %s: %s", key, val))
	}

	return num
}

func IntOpt(key string) *int {
	if StringOpt(key) != nil {
		val := Int(key)
		return &val
	}

	return nil
}

func IntDef(key string, def int) int {
	if v := IntOpt(key); v != nil {
		return *v
	}

	return def
}

// The unit (seconds, minutes, ...) of the duration returned is based on the name of the key.
func Duration(key string) time.Duration {
	d := durationForKey(key)

	return d * time.Duration(Int(key))
}

// The unit (seconds, minutes, ...) of the duration returned is based on the name of the key.
func DurationDef(key string, def int) time.Duration {
	d := durationForKey(key)
	if v := IntOpt(key); v != nil {
		return d * time.Duration(*v)
	}

	return d * time.Duration(def)
}

func Bool(key string) bool {
	return isTrue(key, String(key))
}

func BoolDef(key string, def bool) bool {
	if v := StringOpt(key); v != nil {
		return isTrue(key, *v)
	}

	return def
}

func isTrue(k, v string) bool {
	switch strings.ToLower(v) {
	case "false", "no", "off", "0":
		return false
	case "true", "yes", "on", "1":
		return true
	}

	panic(fmt.Errorf("invalid boolean value for env variable %s: %s", k, v))
}

func durationForKey(key string) time.Duration {
	idx := strings.LastIndex(key, "_")
	if idx < 0 {
		panic(fmt.Errorf("invalid key: %s", key))
	}

	s := strings.ToLower(key[idx+1:])
	switch s {
	case "s", "secs", "seconds":
		return time.Second
	case "m", "mins", "minutes":
		return time.Minute
	case "h", "hrs", "hours":
		return time.Hour
	case "d", "days":
		return 24 * time.Hour
	}

	panic(fmt.Errorf("invalid duration suffix for key %s: %s", key, s))
}
