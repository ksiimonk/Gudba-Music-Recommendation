package usecase

import (
	"context"
	"testing"

	"music-recommender-backend/internal/entity"
)

type fakeAnalyticsRepository struct {
	getStatsFn                  func(ctx context.Context) (*entity.Stats, error)
	getRecommendationMetricsFn  func(ctx context.Context) (*entity.RecommendationMetrics, error)
}

func (f *fakeAnalyticsRepository) GetStats(ctx context.Context) (*entity.Stats, error) {
	if f.getStatsFn != nil {
		return f.getStatsFn(ctx)
	}
	return &entity.Stats{}, nil
}

func (f *fakeAnalyticsRepository) GetRecommendationMetrics(ctx context.Context) (*entity.RecommendationMetrics, error) {
	if f.getRecommendationMetricsFn != nil {
		return f.getRecommendationMetricsFn(ctx)
	}
	return &entity.RecommendationMetrics{}, nil
}

func TestGetStats(t *testing.T) {
	t.Parallel()

	analyticsUseCase := NewAnalyticsUseCase(&fakeAnalyticsRepository{
		getStatsFn: func(_ context.Context) (*entity.Stats, error) {
			return &entity.Stats{UsersCount: 5, TracksCount: 41, EventsCount: 10}, nil
		},
	})

	stats, err := analyticsUseCase.GetStats(context.Background())
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}

	if stats.UsersCount != 5 {
		t.Fatalf("UsersCount = %d, want %d", stats.UsersCount, 5)
	}

	if stats.TracksCount != 41 {
		t.Fatalf("TracksCount = %d, want %d", stats.TracksCount, 41)
	}
}

func TestGetRecommendationMetrics(t *testing.T) {
	t.Parallel()

	analyticsUseCase := NewAnalyticsUseCase(&fakeAnalyticsRepository{
		getRecommendationMetricsFn: func(_ context.Context) (*entity.RecommendationMetrics, error) {
			return &entity.RecommendationMetrics{
				LikeRate:            0.5,
				SkipRate:            0.3,
				DislikeRate:         0.2,
				CoverageTracksCount: 20,
				CoverageTracksPercent: 48.78,
			}, nil
		},
	})

	metrics, err := analyticsUseCase.GetRecommendationMetrics(context.Background())
	if err != nil {
		t.Fatalf("GetRecommendationMetrics() error = %v", err)
	}

	if metrics.LikeRate != 0.5 {
		t.Fatalf("LikeRate = %f, want %f", metrics.LikeRate, 0.5)
	}
}
