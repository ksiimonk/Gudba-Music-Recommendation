package usecase

import (
	"context"
	"errors"

	"music-recommender-backend/internal/entity"
)

type analyticsRepository interface {
	GetStats(ctx context.Context) (*entity.Stats, error)
	GetRecommendationMetrics(ctx context.Context) (*entity.RecommendationMetrics, error)
}

type AnalyticsUseCase struct {
	analyticsRepository analyticsRepository
}

func NewAnalyticsUseCase(analyticsRepository analyticsRepository) *AnalyticsUseCase {
	return &AnalyticsUseCase{analyticsRepository: analyticsRepository}
}

func (u *AnalyticsUseCase) GetStats(ctx context.Context) (*entity.Stats, error) {
	if u == nil || u.analyticsRepository == nil {
		return nil, errors.New("analytics use case repository is nil")
	}

	return u.analyticsRepository.GetStats(ctx)
}

func (u *AnalyticsUseCase) GetRecommendationMetrics(ctx context.Context) (*entity.RecommendationMetrics, error) {
	if u == nil || u.analyticsRepository == nil {
		return nil, errors.New("analytics use case repository is nil")
	}

	return u.analyticsRepository.GetRecommendationMetrics(ctx)
}
