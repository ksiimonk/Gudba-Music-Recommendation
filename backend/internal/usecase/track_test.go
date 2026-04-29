package usecase

import (
	"context"
	"errors"
	"testing"

	"music-recommender-backend/internal/entity"
)

type fakeTrackRepository struct {
	listTracksFn   func(ctx context.Context) ([]entity.Track, error)
	getTrackByIDFn func(ctx context.Context, id int64) (*entity.Track, error)
}

func (f *fakeTrackRepository) ListTracks(ctx context.Context) ([]entity.Track, error) {
	if f.listTracksFn != nil {
		return f.listTracksFn(ctx)
	}

	return []entity.Track{}, nil
}

func (f *fakeTrackRepository) GetTrackByID(ctx context.Context, id int64) (*entity.Track, error) {
	if f.getTrackByIDFn != nil {
		return f.getTrackByIDFn(ctx, id)
	}

	return &entity.Track{ID: id}, nil
}

func TestListTracks(t *testing.T) {
	t.Parallel()

	trackUseCase := NewTrackUseCase(&fakeTrackRepository{
		listTracksFn: func(_ context.Context) ([]entity.Track, error) {
			return []entity.Track{{ID: 1, Title: "Night Drive"}}, nil
		},
	})

	tracks, err := trackUseCase.ListTracks(context.Background())
	if err != nil {
		t.Fatalf("ListTracks() error = %v", err)
	}

	if len(tracks) != 1 {
		t.Fatalf("len(tracks) = %d, want %d", len(tracks), 1)
	}

	if tracks[0].Title != "Night Drive" {
		t.Fatalf("Title = %q, want %q", tracks[0].Title, "Night Drive")
	}
}

func TestGetTrackByIDValidatesID(t *testing.T) {
	t.Parallel()

	trackUseCase := NewTrackUseCase(&fakeTrackRepository{})

	_, err := trackUseCase.GetTrackByID(context.Background(), 0)
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidID)
	}
}

func TestGetTrackByID(t *testing.T) {
	t.Parallel()

	trackUseCase := NewTrackUseCase(&fakeTrackRepository{
		getTrackByIDFn: func(_ context.Context, id int64) (*entity.Track, error) {
			return &entity.Track{ID: id, Title: "Night Drive"}, nil
		},
	})

	track, err := trackUseCase.GetTrackByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTrackByID() error = %v", err)
	}

	if track.ID != 1 {
		t.Fatalf("ID = %d, want %d", track.ID, 1)
	}
}
