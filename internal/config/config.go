package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                 string
	VisualCrossingAPIKey string
	RedisAddr            string
	RateLimitRequests    int
	RateLimitWindow      time.Duration
}

func Load() *Config {
	requests, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "60"))
	windowSecs, _ := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW_SECS", "60"))

	return &Config{
		Port:                 getEnv("PORT", "8080"),
		VisualCrossingAPIKey: getEnv("VCS_API_KEY", ""),
		RedisAddr:            getEnv("REDIS_ADDR", "localhost:6379"),
		RateLimitRequests:    requests,
		RateLimitWindow:      time.Duration(windowSecs) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
