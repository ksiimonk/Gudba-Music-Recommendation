package usecase

import (
	"context"
	"testing"

	"music-recommender-backend/internal/entity"
)

type fakeGenreRepository struct {
	listGenresFn func(ctx context.Context) ([]entity.Genre, error)
}

func (f *fakeGenreRepository) ListGenres(ctx context.Context) ([]entity.Genre, error) {
	if f.listGenresFn != nil {
		return f.listGenresFn(ctx)
	}

	return []entity.Genre{}, nil
}

func TestListGenres(t *testing.T) {
	t.Parallel()

	genreUseCase := NewGenreUseCase(&fakeGenreRepository{
		listGenresFn: func(_ context.Context) ([]entity.Genre, error) {
			return []entity.Genre{{ID: 1, Name: "Lo-Fi Hip Hop"}}, nil
		},
	})

	genres, err := genreUseCase.ListGenres(context.Background())
	if err != nil {
		t.Fatalf("ListGenres() error = %v", err)
	}

	if len(genres) != 1 {
		t.Fatalf("len(genres) = %d, want %d", len(genres), 1)
	}

	if genres[0].Name != "Lo-Fi Hip Hop" {
		t.Fatalf("Name = %q, want %q", genres[0].Name, "Lo-Fi Hip Hop")
	}
}
