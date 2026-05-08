package usecase

import (
	"context"
	"errors"

	"music-recommender-backend/internal/entity"
)

type artistRepository interface {
	ListArtists(ctx context.Context) ([]entity.Artist, error)
}

type ArtistUseCase struct {
	artistRepository artistRepository
}

func NewArtistUseCase(artistRepository artistRepository) *ArtistUseCase {
	return &ArtistUseCase{artistRepository: artistRepository}
}

func (u *ArtistUseCase) ListArtists(ctx context.Context) ([]entity.Artist, error) {
	if u == nil || u.artistRepository == nil {
		return nil, errors.New("artist use case repository is nil")
	}

	return u.artistRepository.ListArtists(ctx)
}
