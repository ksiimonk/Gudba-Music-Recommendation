package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppName  string
	HTTPPort string
	LogLevel string
	PGURL    string
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		AppName:  getEnv("APP_NAME"),
		HTTPPort: getEnv("HTTP_PORT"),
		LogLevel: getEnv("LOG_LEVEL"),
		PGURL:    getEnv("PG_URL"),
	}

	var missing []string
	for _, envVar := range []struct {
		name  string
		value string
	}{
		{name: "APP_NAME", value: cfg.AppName},
		{name: "HTTP_PORT", value: cfg.HTTPPort},
		{name: "LOG_LEVEL", value: cfg.LogLevel},
		{name: "PG_URL", value: cfg.PGURL},
	} {
		if envVar.value == "" {
			missing = append(missing, envVar.name)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	if _, err := strconv.Atoi(cfg.HTTPPort); err != nil {
		return nil, fmt.Errorf("invalid HTTP_PORT %q: must be a valid number", cfg.HTTPPort)
	}

	return cfg, nil
}

func (c *Config) Address() string {
	return ":" + c.HTTPPort
}

func getEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	return value
}
