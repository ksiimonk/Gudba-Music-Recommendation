package config

import (
	"strings"
	"testing"
)

func TestNewConfigSuccess(t *testing.T) {
	t.Setenv("APP_NAME", "music-recommender")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("PG_URL", "postgres://user:pass@localhost:5432/music")

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}

	if cfg.AppName != "music-recommender" {
		t.Fatalf("AppName = %q, want %q", cfg.AppName, "music-recommender")
	}

	if cfg.HTTPPort != "8080" {
		t.Fatalf("HTTPPort = %q, want %q", cfg.HTTPPort, "8080")
	}

	if cfg.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}

	if cfg.PGURL != "postgres://user:pass@localhost:5432/music" {
		t.Fatalf("PGURL = %q, want %q", cfg.PGURL, "postgres://user:pass@localhost:5432/music")
	}
}

func TestNewConfigMissingRequiredEnvVars(t *testing.T) {
	t.Setenv("APP_NAME", "music-recommender")
	t.Setenv("HTTP_PORT", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("PG_URL", "")

	_, err := NewConfig()
	if err == nil {
		t.Fatal("NewConfig() error = nil, want missing env vars error")
	}

	want := "missing required env vars: HTTP_PORT, LOG_LEVEL, PG_URL"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestNewConfigInvalidHTTPPort(t *testing.T) {
	t.Setenv("APP_NAME", "music-recommender")
	t.Setenv("HTTP_PORT", "abc")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("PG_URL", "postgres://user:pass@localhost:5432/music")

	_, err := NewConfig()
	if err == nil {
		t.Fatal("NewConfig() error = nil, want invalid HTTP_PORT error")
	}

	if !strings.Contains(err.Error(), `invalid HTTP_PORT "abc": must be a valid number`) {
		t.Fatalf("error = %q, want invalid HTTP_PORT message", err.Error())
	}
}
