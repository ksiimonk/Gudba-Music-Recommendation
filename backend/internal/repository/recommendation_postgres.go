package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"music-recommender-backend/internal/entity"
)

type RecommendationRepository struct {
	db *sql.DB
}

func NewRecommendationRepository(db *sql.DB) *RecommendationRepository {
	return &RecommendationRepository{db: db}
}

func (r *RecommendationRepository) GetUserProfile(ctx context.Context, userID int64) (*entity.UserProfile, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("recommendation repository database is nil")
	}

	var profile entity.UserProfile
	var genreRaw, artistRaw, trackRaw, contextRaw []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, favorite_genre_ids, favorite_artist_ids, starter_track_ids, contexts, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`, userID).Scan(&profile.UserID, &genreRaw, &artistRaw, &trackRaw, &contextRaw, &profile.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get user profile: %w", err)
	}

	if err := json.Unmarshal(genreRaw, &profile.FavoriteGenreIDs); err != nil {
		return nil, fmt.Errorf("unmarshal favorite_genre_ids: %w", err)
	}
	if err := json.Unmarshal(artistRaw, &profile.FavoriteArtistIDs); err != nil {
		return nil, fmt.Errorf("unmarshal favorite_artist_ids: %w", err)
	}
	if err := json.Unmarshal(trackRaw, &profile.StarterTrackIDs); err != nil {
		return nil, fmt.Errorf("unmarshal starter_track_ids: %w", err)
	}
	if err := json.Unmarshal(contextRaw, &profile.Contexts); err != nil {
		return nil, fmt.Errorf("unmarshal contexts: %w", err)
	}

	return &profile, nil
}

func (r *RecommendationRepository) GetRecentEvents(ctx context.Context, userID int64, limit int) ([]entity.Event, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("recommendation repository database is nil")
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, entity_type, entity_id, event_type, created_at
		FROM events
		WHERE user_id = $1 AND created_at > $2
		ORDER BY created_at DESC
		LIMIT $3
	`, userID, time.Now().Add(-24*time.Hour), limit)

	if err != nil {
		return nil, fmt.Errorf("list recent events: %w", err)
	}

	defer rows.Close()

	var events []entity.Event
	for rows.Next() {
		var event entity.Event
		if err := rows.Scan(&event.ID, &event.UserID, &event.EntityType, &event.EntityID, &event.EventType, &event.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate events: %w", err)
	}

	return events, nil
}

func (r *RecommendationRepository) ListCandidateTracks(ctx context.Context) ([]entity.Track, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("recommendation repository database is nil")
	}

	rows, err := r.db.QueryContext(ctx, tracksSelectQuery + `
		GROUP BY t.id, a.id
		ORDER BY t.popularity_score DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list candidate tracks: %w", err)
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
		return nil, fmt.Errorf("iterate candidate tracks: %w", err)
	}

	return tracks, nil
}

func (r *RecommendationRepository) SaveRecommendationRequest(ctx context.Context, userID int64, requestType string) (int64, error) {
	if r == nil || r.db == nil {
		return 0, errors.New("recommendation repository database is nil")
	}

	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO recommendation_requests (user_id, request_type)
		VALUES ($1, $2)
		RETURNING id
	`, userID, requestType).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("save recommendation request: %w", err)
	}

	return id, nil
}

func (r *RecommendationRepository) SaveImpressions(ctx context.Context, impressions []entity.RecommendationImpression) ([]int64, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("recommendation repository database is nil")
	}

	var ids []int64

	for _, imp := range impressions {
		var id int64
		err := r.db.QueryRowContext(ctx, `
			INSERT INTO recommendation_impressions (request_id, user_id, entity_type, entity_id, rank, score, explanation)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`, imp.RequestID, imp.UserID, imp.EntityType, imp.EntityID, imp.Rank, imp.Score, imp.Explanation).Scan(&id)

		if err != nil {
			return ids, fmt.Errorf("save impression: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (r *RecommendationRepository) SaveFactors(ctx context.Context, factors []entity.RecommendationFactor) error {
	if r == nil || r.db == nil {
		return errors.New("recommendation repository database is nil")
	}

	for _, f := range factors {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO recommendation_factors (impression_id, factor_name, factor_value, weight, contribution)
			VALUES ($1, $2, $3, $4, $5)
		`, f.ImpressionID, f.FactorName, f.FactorValue, f.Weight, f.Contribution)

		if err != nil {
			return fmt.Errorf("save factor: %w", err)
		}
	}

	return nil
}

func (r *RecommendationRepository) GetFavoriteTrackIDs(ctx context.Context, userID int64) (map[int64]bool, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("recommendation repository database is nil")
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT pt.track_id
		FROM playlist_tracks pt
		INNER JOIN playlists p ON p.id = pt.playlist_id
		WHERE p.user_id = $1 AND p.name = 'Мои любимые треки' AND p.is_public = false
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("get favorite track ids: %w", err)
	}
	defer rows.Close()

	result := make(map[int64]bool)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan favorite track id: %w", err)
		}
		result[id] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate favorite track ids: %w", err)
	}

	return result, nil
}

func (r *RecommendationRepository) ListPlaylistsWithTracks(ctx context.Context) ([]entity.Playlist, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("recommendation repository database is nil")
	}

	playlistRows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, name, description, is_public, created_at, updated_at
		FROM playlists
		WHERE is_public = true
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list playlists: %w", err)
	}

	var playlists []entity.Playlist
	for playlistRows.Next() {
		var pl entity.Playlist
		var desc sql.NullString

		if err := playlistRows.Scan(&pl.ID, &pl.UserID, &pl.Name, &desc, &pl.IsPublic, &pl.CreatedAt, &pl.UpdatedAt); err != nil {
			playlistRows.Close()
			return nil, fmt.Errorf("scan playlist: %w", err)
		}

		pl.Description = nullableString(desc)
		playlists = append(playlists, pl)
	}

	playlistRows.Close()

	if err := playlistRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate playlists: %w", err)
	}

	for i, pl := range playlists {
		trackRows, err := r.db.QueryContext(ctx, playlistTracksSelectQuery, pl.ID)
		if err != nil {
			return nil, fmt.Errorf("list tracks for playlist %d: %w", pl.ID, err)
		}

		for trackRows.Next() {
			track, err := scanTrack(trackRows)
			if err != nil {
				trackRows.Close()
				return nil, err
			}

			playlists[i].Tracks = append(playlists[i].Tracks, track)
		}

		if err := trackRows.Err(); err != nil {
			trackRows.Close()
			return nil, fmt.Errorf("iterate tracks for playlist %d: %w", pl.ID, err)
		}

		trackRows.Close()

		playlists[i].TrackCount = len(playlists[i].Tracks)
	}

	return playlists, nil
}
