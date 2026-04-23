package app

import (
	"strings"
	"testing"
)

func TestNewPostgresDBInvalidPGURL(t *testing.T) {
	t.Parallel()

	_, err := newPostgresDB("postgres://postgres:postgres@127.0.0.1:1/music?sslmode=disable&connect_timeout=1")
	if err == nil {
		t.Fatal("newPostgresDB() error = nil, want connection error")
	}

	if !strings.Contains(err.Error(), "connect to postgres from PG_URL") {
		t.Fatalf("error = %q, want clear PG_URL bootstrap error", err.Error())
	}
}
