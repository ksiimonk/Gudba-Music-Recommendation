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

func (r *PlaylistRepository) CreateFavoritesPlaylist(ctx context.Context, userID int64) (*entity.Playlist, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("playlist repository database is nil")
	}

	playlist := &entity.Playlist{
		UserID:      userID,
		Name:        "Мои любимые треки",
		Description: "Треки, которые вы отметили как понравившиеся",
		IsPublic:    false,
	}

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO playlists (user_id, name, description, is_public)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, playlist.UserID, playlist.Name, playlist.Description, playlist.IsPublic).
		Scan(&playlist.ID, &playlist.CreatedAt, &playlist.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create favorites playlist: %w", err)
	}

	return playlist, nil
}

func (r *PlaylistRepository) GetFavoritesPlaylistByUserID(ctx context.Context, userID int64) (*entity.Playlist, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("playlist repository database is nil")
	}

	const query = `
		SELECT p.id, p.user_id, p.name, p.description, p.is_public, p.created_at, p.updated_at, COUNT(pt.track_id) AS track_count
		FROM playlists p
		LEFT JOIN playlist_tracks pt ON pt.playlist_id = p.id
		WHERE p.user_id = $1 AND p.name = 'Мои любимые треки' AND p.is_public = false
		GROUP BY p.id
	`
	playlist, err := scanPlaylistSummary(r.db.QueryRowContext(ctx, query, userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlaylistNotFound
		}
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, playlistTracksSelectQuery, playlist.ID)
	if err != nil {
		return nil, fmt.Errorf("list favorites tracks: %w", err)
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
		return nil, fmt.Errorf("iterate favorites tracks: %w", err)
	}

	return &playlist, nil
}

func (r *PlaylistRepository) AddTrackToPlaylist(ctx context.Context, playlistID, trackID int64) error {
	if r == nil || r.db == nil {
		return errors.New("playlist repository database is nil")
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO playlist_tracks (playlist_id, track_id, position)
		VALUES ($1, $2, COALESCE((SELECT MAX(position) + 1 FROM playlist_tracks WHERE playlist_id = $1), 0))
		ON CONFLICT (playlist_id, track_id) DO NOTHING
	`, playlistID, trackID)
	if err != nil {
		return fmt.Errorf("add track to playlist: %w", err)
	}

	return nil
}

func (r *PlaylistRepository) RemoveTrackFromPlaylist(ctx context.Context, playlistID, trackID int64) error {
	if r == nil || r.db == nil {
		return errors.New("playlist repository database is nil")
	}

	_, err := r.db.ExecContext(ctx, `
		DELETE FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2
	`, playlistID, trackID)
	if err != nil {
		return fmt.Errorf("remove track from playlist: %w", err)
	}

	return nil
}

func (r *PlaylistRepository) IsTrackFavorited(ctx context.Context, userID, trackID int64) (bool, error) {
	if r == nil || r.db == nil {
		return false, errors.New("playlist repository database is nil")
	}

	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM playlist_tracks pt
		INNER JOIN playlists p ON p.id = pt.playlist_id
		WHERE p.user_id = $1 AND p.name = 'Мои любимые треки' AND pt.track_id = $2
	`, userID, trackID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check track favorited: %w", err)
	}

	return count > 0, nil
}

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
