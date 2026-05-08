package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"music-recommender-backend/internal/entity"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) GetStats(ctx context.Context) (*entity.Stats, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("analytics repository database is nil")
	}

	var stats entity.Stats

	err := r.db.QueryRowContext(ctx, `
		SELECT
			(SELECT COUNT(*) FROM users) AS users_count,
			(SELECT COUNT(*) FROM tracks) AS tracks_count,
			(SELECT COUNT(*) FROM playlists WHERE is_public = true) AS playlists_count,
			(SELECT COUNT(*) FROM events) AS events_count,
			(SELECT COUNT(*) FROM events WHERE event_type = 'like') AS likes_count,
			(SELECT COUNT(*) FROM events WHERE event_type = 'dislike') AS dislikes_count,
			(SELECT COUNT(*) FROM events WHERE event_type = 'skip') AS skips_count,
			(SELECT COUNT(*) FROM events WHERE event_type = 'play') AS plays_count,
			(SELECT COUNT(*) FROM recommendation_requests) AS req_count,
			(SELECT COUNT(*) FROM recommendation_impressions) AS imp_count
	`).Scan(
		&stats.UsersCount,
		&stats.TracksCount,
		&stats.PlaylistsCount,
		&stats.EventsCount,
		&stats.LikesCount,
		&stats.DislikesCount,
		&stats.SkipsCount,
		&stats.PlaysCount,
		&stats.RecommendationRequestsCount,
		&stats.RecommendationImpressionsCount,
	)

	if err != nil {
		return nil, fmt.Errorf("get stats: %w", err)
	}

	return &stats, nil
}

func (r *AnalyticsRepository) GetRecommendationMetrics(ctx context.Context) (*entity.RecommendationMetrics, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("analytics repository database is nil")
	}

	var metrics entity.RecommendationMetrics
	var likesCount, dislikesCount, skipsCount int64
	var totalTracksCount int64

	err := r.db.QueryRowContext(ctx, `
		SELECT
			COALESCE((SELECT COUNT(*) FROM events WHERE event_type = 'like'), 0),
			COALESCE((SELECT COUNT(*) FROM events WHERE event_type = 'dislike'), 0),
			COALESCE((SELECT COUNT(*) FROM events WHERE event_type = 'skip'), 0),
			COALESCE((SELECT COUNT(DISTINCT entity_id) FROM recommendation_impressions WHERE entity_type = 'track'), 0),
			(SELECT COUNT(*) FROM tracks)
	`).Scan(
		&likesCount,
		&dislikesCount,
		&skipsCount,
		&metrics.CoverageTracksCount,
		&totalTracksCount,
	)

	if err != nil {
		return nil, fmt.Errorf("get recommendation metrics: %w", err)
	}

	totalInteractions := float64(likesCount + dislikesCount + skipsCount)

	if totalInteractions > 0 {
		metrics.LikeRate = float64(likesCount) / totalInteractions
		metrics.DislikeRate = float64(dislikesCount) / totalInteractions
		metrics.SkipRate = float64(skipsCount) / totalInteractions
	}

	if totalTracksCount > 0 {
		metrics.CoverageTracksPercent = (float64(metrics.CoverageTracksCount) / float64(totalTracksCount)) * 100
	}

	return &metrics, nil
}
