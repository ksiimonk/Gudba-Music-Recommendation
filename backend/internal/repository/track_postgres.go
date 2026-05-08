package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"music-recommender-backend/internal/entity"
)

var ErrTrackNotFound = errors.New("track not found")

type TrackRepository struct {
	db *sql.DB
}

func NewTrackRepository(db *sql.DB) *TrackRepository {
	return &TrackRepository{db: db}
}

func (r *TrackRepository) ListTracks(ctx context.Context) ([]entity.Track, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("track repository database is nil")
	}

	rows, err := r.db.QueryContext(ctx, tracksSelectQuery+`
		GROUP BY t.id, a.id
		ORDER BY t.popularity_score DESC, t.id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list tracks: %w", err)
	}
	defer rows.Close()

	var tracks []entity.Track
	for rows.Next() {
		track, err := scanTrack(rows)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tracks: %w", err)
	}

	return tracks, nil
}

func (r *TrackRepository) GetTrackByID(ctx context.Context, id int64) (*entity.Track, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("track repository database is nil")
	}

	row := r.db.QueryRowContext(ctx, tracksSelectQuery+`
		WHERE t.id = $1
		GROUP BY t.id, a.id
	`, id)

	track, err := scanTrack(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTrackNotFound
		}

		return nil, err
	}

	return &track, nil
}

const trackFields = `
	t.id,
	t.title,
	t.duration_ms,
	t.preview_url,
	t.spotify_url,
	t.cover_url,
	t.popularity_score,
	t.created_at,
	t.updated_at,
	a.id,
	a.name,
	a.image_url,
	a.bio,
	a.spotify_url,
	a.created_at,
	a.updated_at,
	COALESCE(
		json_agg(
			json_build_object('id', g.id, 'name', g.name)
			ORDER BY g.name
		) FILTER (WHERE g.id IS NOT NULL),
		'[]'
	) AS genres
`

const tracksSelectQuery = `
	SELECT ` + trackFields + `
	FROM tracks t
	INNER JOIN artists a ON a.id = t.artist_id
	LEFT JOIN track_genres tg ON tg.track_id = t.id
	LEFT JOIN genres g ON g.id = tg.genre_id
`

type trackRowScanner interface {
	Scan(dest ...any) error
}

func scanTrack(scanner trackRowScanner) (entity.Track, error) {
	var track entity.Track
	var artist entity.Artist
	var trackPreviewURL sql.NullString
	var trackSpotifyURL sql.NullString
	var coverURL sql.NullString
	var artistImageURL sql.NullString
	var artistBio sql.NullString
	var artistSpotifyURL sql.NullString
	var genresPayload []byte

	if err := scanner.Scan(
		&track.ID,
		&track.Title,
		&track.DurationMS,
		&trackPreviewURL,
		&trackSpotifyURL,
		&coverURL,
		&track.PopularityScore,
		&track.CreatedAt,
		&track.UpdatedAt,
		&artist.ID,
		&artist.Name,
		&artistImageURL,
		&artistBio,
		&artistSpotifyURL,
		&artist.CreatedAt,
		&artist.UpdatedAt,
		&genresPayload,
	); err != nil {
		return entity.Track{}, err
	}

	genres, err := decodeTrackGenres(genresPayload)
	if err != nil {
		return entity.Track{}, err
	}

	track.PreviewURL = nullableString(trackPreviewURL)
	track.SpotifyURL = nullableString(trackSpotifyURL)
	track.CoverURL = nullableString(coverURL)
	artist.ImageURL = nullableString(artistImageURL)
	artist.Bio = nullableString(artistBio)
	artist.SpotifyURL = nullableString(artistSpotifyURL)
	track.Artist = artist
	track.Genres = genres

	return track, nil
}

func decodeTrackGenres(payload []byte) ([]entity.Genre, error) {
	var genres []entity.Genre
	if err := json.Unmarshal(payload, &genres); err != nil {
		return nil, fmt.Errorf("decode track genres: %w", err)
	}

	if genres == nil {
		return []entity.Genre{}, nil
	}

	return genres, nil
}

func nullableString(value sql.NullString) string {
	if !value.Valid {
		return ""
	}

	return value.String
}
