package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

const commandTimeout = 30 * time.Second

type migration struct {
	Version string
	Path    string
}

func main() {
	_ = godotenv.Load(".env", filepath.Join("backend", ".env"))

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	pgURL := strings.TrimSpace(os.Getenv("PG_URL"))
	if pgURL == "" {
		log.Fatal("PG_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	switch command {
	case "up":
		if err := ensureDatabase(ctx, pgURL); err != nil {
			log.Fatalf("ensure database: %v", err)
		}

		db, err := openPostgres(ctx, pgURL)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if err := applyUpMigrations(ctx, db); err != nil {
			log.Fatalf("apply migrations: %v", err)
		}
	case "reset":
		if err := ensureDatabase(ctx, pgURL); err != nil {
			log.Fatalf("ensure database: %v", err)
		}

		db, err := openPostgres(ctx, pgURL)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if err := resetSchema(ctx, db); err != nil {
			log.Fatalf("reset schema: %v", err)
		}

		if err := applyUpMigrations(ctx, db); err != nil {
			log.Fatalf("apply migrations: %v", err)
		}
	case "status":
		db, err := openPostgres(ctx, pgURL)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if err := printMigrationStatus(ctx, db); err != nil {
			log.Fatalf("migration status: %v", err)
		}
	default:
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/migrate [up|reset|status]")
		os.Exit(2)
	}
}

func openPostgres(ctx context.Context, pgURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	return db, nil
}

func ensureDatabase(ctx context.Context, pgURL string) error {
	targetURL, err := url.Parse(pgURL)
	if err != nil {
		return fmt.Errorf("parse PG_URL: %w", err)
	}

	databaseName := strings.TrimPrefix(targetURL.Path, "/")
	if databaseName == "" || databaseName == "postgres" {
		return nil
	}

	adminURL := *targetURL
	adminURL.Path = "/postgres"

	db, err := openPostgres(ctx, adminURL.String())
	if err != nil {
		return err
	}
	defer db.Close()

	var exists bool
	if err := db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM pg_database
			WHERE datname = $1
		)
	`, databaseName).Scan(&exists); err != nil {
		return fmt.Errorf("check database %q: %w", databaseName, err)
	}

	if exists {
		return nil
	}

	if _, err := db.ExecContext(ctx, "CREATE DATABASE "+pq.QuoteIdentifier(databaseName)); err != nil {
		return fmt.Errorf("create database %q: %w", databaseName, err)
	}

	fmt.Printf("created database %s\n", databaseName)
	return nil
}

func resetSchema(ctx context.Context, db *sql.DB) error {
	statements := []string{
		"DROP SCHEMA IF EXISTS public CASCADE",
		"CREATE SCHEMA public",
		"GRANT ALL ON SCHEMA public TO PUBLIC",
	}

	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("execute %q: %w", statement, err)
		}
	}

	fmt.Println("reset public schema")
	return nil
}

func applyUpMigrations(ctx context.Context, db *sql.DB) error {
	if err := ensureSchemaMigrations(ctx, db); err != nil {
		return err
	}

	migrations, err := listUpMigrations()
	if err != nil {
		return err
	}

	appliedCount := 0
	for _, migration := range migrations {
		applied, err := isMigrationApplied(ctx, db, migration.Version)
		if err != nil {
			return err
		}

		if applied {
			fmt.Printf("skip %s\n", migration.Version)
			continue
		}

		if err := applyMigration(ctx, db, migration); err != nil {
			return err
		}

		appliedCount++
		fmt.Printf("applied %s\n", migration.Version)
	}

	fmt.Printf("done: %d migration(s) applied\n", appliedCount)
	return nil
}

func printMigrationStatus(ctx context.Context, db *sql.DB) error {
	if err := ensureSchemaMigrations(ctx, db); err != nil {
		return err
	}

	migrations, err := listUpMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		status := "pending"
		applied, err := isMigrationApplied(ctx, db, migration.Version)
		if err != nil {
			return err
		}

		if applied {
			status = "applied"
		}

		fmt.Printf("%-8s %s\n", status, migration.Version)
	}

	return nil
}

func listUpMigrations() ([]migration, error) {
	migrationsDir, err := findMigrationsDir()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("read migrations directory: %w", err)
	}

	var migrations []migration
	for _, file := range files {
		name := file.Name()
		if file.IsDir() || !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		migrations = append(migrations, migration{
			Version: name,
			Path:    filepath.Join(migrationsDir, name),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func findMigrationsDir() (string, error) {
	candidates := []string{
		"migrations",
		filepath.Join("backend", "migrations"),
	}

	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate, nil
		}
	}

	return "", errors.New("migrations directory was not found")
}

func ensureSchemaMigrations(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	return nil
}

func isMigrationApplied(ctx context.Context, db *sql.DB, version string) (bool, error) {
	var exists bool
	if err := db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM schema_migrations
			WHERE version = $1
		)
	`, version).Scan(&exists); err != nil {
		return false, fmt.Errorf("check migration %s: %w", version, err)
	}

	return exists, nil
}

func applyMigration(ctx context.Context, db *sql.DB, migration migration) error {
	payload, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf("read migration %s: %w", migration.Version, err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start migration %s transaction: %w", migration.Version, err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, string(payload)); err != nil {
		return fmt.Errorf("execute migration %s: %w", migration.Version, err)
	}

	if err := recordMigrationTx(ctx, tx, migration.Version); err != nil {
		return fmt.Errorf("record migration %s: %w", migration.Version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", migration.Version, err)
	}

	return nil
}

func recordMigrationTx(ctx context.Context, tx *sql.Tx, version string) error {
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO schema_migrations (version)
		VALUES ($1)
		ON CONFLICT (version) DO NOTHING
	`, version); err != nil {
		return err
	}

	return nil
}
