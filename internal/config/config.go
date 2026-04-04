package config

import (
	"os"
	"time"
)

type Config struct {
	Port         string
	LogLevel     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() *Config {
	return &Config{
		Port:         GetEnv("APP_PORT", "8080"),
		LogLevel:     GetEnv("LOG_LEVEL", "info"),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func GetEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
