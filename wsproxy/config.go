package wsproxy

import "time"

type Config struct {
	UserAgent  string
	MaxRetries int
	Timeout    time.Duration
}

func NewConfig() *Config {
	return &Config{
		UserAgent:  "Mozilla/5.0 (X11; Linux x86_64; rv:136.0) Gecko/20100101 Firefox/136.0",
		MaxRetries: 5,
		Timeout:    15 * time.Second,
	}
}
