package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const postgresPingTimeout = 5 * time.Second

func newPostgresDB(pgURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection from PG_URL: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgresPingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect to postgres from PG_URL: %w", err)
	}

	return db, nil
}
