package usecase

import (
	"context"
	"errors"

	"music-recommender-backend/internal/entity"
)

var ErrInvalidID = errors.New("id must be greater than zero")

type trackRepository interface {
	ListTracks(ctx context.Context) ([]entity.Track, error)
	GetTrackByID(ctx context.Context, id int64) (*entity.Track, error)
}

type TrackUseCase struct {
	trackRepository trackRepository
}

func NewTrackUseCase(trackRepository trackRepository) *TrackUseCase {
	return &TrackUseCase{trackRepository: trackRepository}
}

func (u *TrackUseCase) ListTracks(ctx context.Context) ([]entity.Track, error) {
	if u == nil || u.trackRepository == nil {
		return nil, errors.New("track use case repository is nil")
	}

	return u.trackRepository.ListTracks(ctx)
}

func (u *TrackUseCase) GetTrackByID(ctx context.Context, id int64) (*entity.Track, error) {
	if u == nil || u.trackRepository == nil {
		return nil, errors.New("track use case repository is nil")
	}

	if id <= 0 {
		return nil, ErrInvalidID
	}

	return u.trackRepository.GetTrackByID(ctx, id)
}
