package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"music-recommender-backend/internal/entity"
)

var ErrPlaylistNotFound = errors.New("playlist not found")

type PlaylistRepository struct {
	db *sql.DB
}

func NewPlaylistRepository(db *sql.DB) *PlaylistRepository {
	return &PlaylistRepository{db: db}
}

func (r *PlaylistRepository) ListPlaylists(ctx context.Context) ([]entity.Playlist, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("playlist repository database is nil")
	}

	const query = `
		SELECT p.id, p.user_id, p.name, p.description, p.is_public, p.created_at, p.updated_at, COUNT(pt.track_id) AS track_count
		FROM playlists p
		LEFT JOIN playlist_tracks pt ON pt.playlist_id = p.id
		WHERE p.is_public = true
		GROUP BY p.id
		ORDER BY p.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list playlists: %w", err)
	}
	defer rows.Close()

	var playlists []entity.Playlist
	for rows.Next() {
		playlist, err := scanPlaylistSummary(rows)
		if err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate playlists: %w", err)
	}

	return playlists, nil
}

func (r *PlaylistRepository) GetPlaylistByID(ctx context.Context, id int64) (*entity.Playlist, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("playlist repository database is nil")
	}

	const playlistQuery = `
		SELECT p.id, p.user_id, p.name, p.description, p.is_public, p.created_at, p.updated_at, COUNT(pt.track_id) AS track_count
		FROM playlists p
		LEFT JOIN playlist_tracks pt ON pt.playlist_id = p.id
		WHERE p.id = $1 AND p.is_public = true
		GROUP BY p.id
	`

	playlist, err := scanPlaylistSummary(r.db.QueryRowContext(ctx, playlistQuery, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlaylistNotFound
		}

		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, playlistTracksSelectQuery, id)
	if err != nil {
		return nil, fmt.Errorf("list playlist tracks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		track, err := scanTrack(rows)
		if err != nil {
			return nil, err
		}

		playlist.Tracks = append(playlist.Tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate playlist tracks: %w", err)
	}

	return &playlist, nil
}

const playlistTracksSelectQuery = `
	SELECT ` + trackFields + `
	FROM playlist_tracks pt
	INNER JOIN tracks t ON t.id = pt.track_id
	INNER JOIN artists a ON a.id = t.artist_id
	LEFT JOIN track_genres tg ON tg.track_id = t.id
	LEFT JOIN genres g ON g.id = tg.genre_id
	WHERE pt.playlist_id = $1
	GROUP BY pt.position, t.id, a.id
	ORDER BY pt.position ASC, t.id ASC
`

func scanPlaylistSummary(scanner trackRowScanner) (entity.Playlist, error) {
	var playlist entity.Playlist
	var description sql.NullString

	if err := scanner.Scan(
		&playlist.ID,
		&playlist.UserID,
		&playlist.Name,
		&description,
		&playlist.IsPublic,
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
		&playlist.TrackCount,
	); err != nil {
		return entity.Playlist{}, err
	}

	playlist.Description = nullableString(description)
	return playlist, nil
}
