package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"music-recommender-backend/internal/entity"
)

type ArtistRepository struct {
	db *sql.DB
}

func NewArtistRepository(db *sql.DB) *ArtistRepository {
	return &ArtistRepository{db: db}
}

func (r *ArtistRepository) ListArtists(ctx context.Context) ([]entity.Artist, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("artist repository database is nil")
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, image_url, bio, spotify_url, created_at, updated_at
		FROM artists
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list artists: %w", err)
	}
	defer rows.Close()

	var artists []entity.Artist
	for rows.Next() {
		var artist entity.Artist
		var imageURL sql.NullString
		var bio sql.NullString
		var spotifyURL sql.NullString

		if err := rows.Scan(
			&artist.ID,
			&artist.Name,
			&imageURL,
			&bio,
			&spotifyURL,
			&artist.CreatedAt,
			&artist.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan artist: %w", err)
		}

		artist.ImageURL = nullableString(imageURL)
		artist.Bio = nullableString(bio)
		artist.SpotifyURL = nullableString(spotifyURL)

		artists = append(artists, artist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate artists: %w", err)
	}

	return artists, nil
}
