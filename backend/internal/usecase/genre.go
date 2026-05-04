package usecase

import (
	"context"
	"errors"

	"music-recommender-backend/internal/entity"
)

type genreRepository interface {
	ListGenres(ctx context.Context) ([]entity.Genre, error)
}

type GenreUseCase struct {
	genreRepository genreRepository
}

func NewGenreUseCase(genreRepository genreRepository) *GenreUseCase {
	return &GenreUseCase{genreRepository: genreRepository}
}

func (u *GenreUseCase) ListGenres(ctx context.Context) ([]entity.Genre, error) {
	if u == nil || u.genreRepository == nil {
		return nil, errors.New("genre use case repository is nil")
	}

	return u.genreRepository.ListGenres(ctx)
}
