package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"music-recommender-backend/internal/entity"
)

var ErrProfileNotFound = errors.New("user profile not found")

type OnboardingRepository struct {
	db *sql.DB
}

func NewOnboardingRepository(db *sql.DB) *OnboardingRepository {
	return &OnboardingRepository{db: db}
}

func (r *OnboardingRepository) SaveOnboarding(ctx context.Context, userID int64, request *entity.OnboardingRequest) error {
	if r == nil || r.db == nil {
		return errors.New("onboarding repository database is nil")
	}

	if request == nil {
		return errors.New("onboarding request is nil")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM user_preferences WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("delete old preferences: %w", err)
	}

	insertQuery := `
		INSERT INTO user_preferences (user_id, preference_type, reference_id, value)
		VALUES ($1, $2, $3, $4)
	`

	for _, id := range request.GenreIDs {
		if _, err := tx.ExecContext(ctx, insertQuery, userID, "genre", id, nil); err != nil {
			return fmt.Errorf("insert genre preference: %w", err)
		}
	}

	for _, id := range request.ArtistIDs {
		if _, err := tx.ExecContext(ctx, insertQuery, userID, "artist", id, nil); err != nil {
			return fmt.Errorf("insert artist preference: %w", err)
		}
	}

	for _, id := range request.TrackIDs {
		if _, err := tx.ExecContext(ctx, insertQuery, userID, "track", id, nil); err != nil {
			return fmt.Errorf("insert track preference: %w", err)
		}
	}

	for _, context := range request.Contexts {
		if _, err := tx.ExecContext(ctx, insertQuery, userID, "context", nil, context); err != nil {
			return fmt.Errorf("insert context preference: %w", err)
		}
	}

	if err := r.upsertProfileTx(ctx, tx, userID, request); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit onboarding: %w", err)
	}

	return nil
}

func (r *OnboardingRepository) upsertProfileTx(ctx context.Context, tx *sql.Tx, userID int64, request *entity.OnboardingRequest) error {
	genreIDs, _ := json.Marshal(request.GenreIDs)
	artistIDs, _ := json.Marshal(request.ArtistIDs)
	trackIDs, _ := json.Marshal(request.TrackIDs)
	contexts, _ := json.Marshal(request.Contexts)

	_, err := tx.ExecContext(ctx, `
		INSERT INTO user_profiles (user_id, favorite_genre_ids, favorite_artist_ids, starter_track_ids, contexts, updated_at)
		VALUES ($1, $2, $3, $4, $5, now())
		ON CONFLICT (user_id) DO UPDATE SET
			favorite_genre_ids = $2,
			favorite_artist_ids = $3,
			starter_track_ids = $4,
			contexts = $5,
			updated_at = now()
	`, userID, string(genreIDs), string(artistIDs), string(trackIDs), string(contexts))

	if err != nil {
		return fmt.Errorf("upsert user profile: %w", err)
	}

	return nil
}

func (r *OnboardingRepository) GetUserProfile(ctx context.Context, userID int64) (*entity.UserProfile, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("onboarding repository database is nil")
	}

	var profile entity.UserProfile
	var genreRaw, artistRaw, trackRaw, contextRaw []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, favorite_genre_ids, favorite_artist_ids, starter_track_ids, contexts, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`, userID).Scan(
		&profile.UserID,
		&genreRaw,
		&artistRaw,
		&trackRaw,
		&contextRaw,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProfileNotFound
		}

		return nil, fmt.Errorf("get user profile: %w", err)
	}

	if err := json.Unmarshal(genreRaw, &profile.FavoriteGenreIDs); err != nil {
		return nil, fmt.Errorf("unmarshal genre ids: %w", err)
	}

	if err := json.Unmarshal(artistRaw, &profile.FavoriteArtistIDs); err != nil {
		return nil, fmt.Errorf("unmarshal artist ids: %w", err)
	}

	if err := json.Unmarshal(trackRaw, &profile.StarterTrackIDs); err != nil {
		return nil, fmt.Errorf("unmarshal track ids: %w", err)
	}

	if err := json.Unmarshal(contextRaw, &profile.Contexts); err != nil {
		return nil, fmt.Errorf("unmarshal contexts: %w", err)
	}

	return &profile, nil
}
