package usecase

import (
	"context"
	"testing"

	"music-recommender-backend/internal/entity"
)

type fakeArtistRepository struct {
	listArtistsFn func(ctx context.Context) ([]entity.Artist, error)
}

func (f *fakeArtistRepository) ListArtists(ctx context.Context) ([]entity.Artist, error) {
	if f.listArtistsFn != nil {
		return f.listArtistsFn(ctx)
	}

	return []entity.Artist{}, nil
}

func TestListArtists(t *testing.T) {
	t.Parallel()

	artistUseCase := NewArtistUseCase(&fakeArtistRepository{
		listArtistsFn: func(_ context.Context) ([]entity.Artist, error) {
			return []entity.Artist{{ID: 1, Name: "Ideal"}}, nil
		},
	})

	artists, err := artistUseCase.ListArtists(context.Background())
	if err != nil {
		t.Fatalf("ListArtists() error = %v", err)
	}

	if len(artists) != 1 {
		t.Fatalf("len(artists) = %d, want %d", len(artists), 1)
	}

	if artists[0].Name != "Ideal" {
		t.Fatalf("Name = %q, want %q", artists[0].Name, "Ideal")
	}
}
