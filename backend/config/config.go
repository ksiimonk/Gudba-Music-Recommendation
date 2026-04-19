package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppName  string
	HTTPPort string
	LogLevel string
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		AppName:  os.Getenv("APP_NAME"),
		HTTPPort: os.Getenv("HTTP_PORT"),
		LogLevel: os.Getenv("LOG_LEVEL"),
	}

	if cfg.AppName == "" {
		return nil, fmt.Errorf("APP_NAME is required")
	}

	if cfg.HTTPPort == "" {
		return nil, fmt.Errorf("HTTP_PORT is required")
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "debug"
	}

	return cfg, nil
}
