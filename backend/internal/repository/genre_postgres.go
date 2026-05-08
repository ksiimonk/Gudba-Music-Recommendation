package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"music-recommender-backend/internal/entity"
)

type GenreRepository struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{db: db}
}

func (r *GenreRepository) ListGenres(ctx context.Context) ([]entity.Genre, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("genre repository database is nil")
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, created_at
		FROM genres
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list genres: %w", err)
	}
	defer rows.Close()

	var genres []entity.Genre
	for rows.Next() {
		var genre entity.Genre
		if err := rows.Scan(&genre.ID, &genre.Name, &genre.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan genre: %w", err)
		}

		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate genres: %w", err)
	}

	return genres, nil
}
