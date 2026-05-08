package usecase

import (
	"context"
	"errors"
	"testing"

	"music-recommender-backend/internal/entity"
)

type fakePlaylistRepository struct {
	listPlaylistsFn   func(ctx context.Context) ([]entity.Playlist, error)
	getPlaylistByIDFn func(ctx context.Context, id int64) (*entity.Playlist, error)
}

func (f *fakePlaylistRepository) GetFavoritesPlaylistByUserID(ctx context.Context, userID int64) (*entity.Playlist, error) {
	return &entity.Playlist{ID: 1, UserID: userID, Name: "Мои любимые треки"}, nil
}

func (f *fakePlaylistRepository) ListPlaylists(ctx context.Context) ([]entity.Playlist, error) {
	if f.listPlaylistsFn != nil {
		return f.listPlaylistsFn(ctx)
	}

	return []entity.Playlist{}, nil
}

func (f *fakePlaylistRepository) GetPlaylistByID(ctx context.Context, id int64) (*entity.Playlist, error) {
	if f.getPlaylistByIDFn != nil {
		return f.getPlaylistByIDFn(ctx, id)
	}

	return &entity.Playlist{ID: id}, nil
}

func TestListPlaylists(t *testing.T) {
	t.Parallel()

	playlistUseCase := NewPlaylistUseCase(&fakePlaylistRepository{
		listPlaylistsFn: func(_ context.Context) ([]entity.Playlist, error) {
			return []entity.Playlist{{ID: 1, Name: "Lo-Fi Hip Hop Mix"}}, nil
		},
	}, nil)

	playlists, err := playlistUseCase.ListPlaylists(context.Background())
	if err != nil {
		t.Fatalf("ListPlaylists() error = %v", err)
	}

	if len(playlists) != 1 {
		t.Fatalf("len(playlists) = %d, want %d", len(playlists), 1)
	}

	if playlists[0].Name != "Lo-Fi Hip Hop Mix" {
		t.Fatalf("Name = %q, want %q", playlists[0].Name, "Lo-Fi Hip Hop Mix")
	}
}

func TestGetPlaylistByIDValidatesID(t *testing.T) {
	t.Parallel()

	playlistUseCase := NewPlaylistUseCase(&fakePlaylistRepository{}, nil)

	_, err := playlistUseCase.GetPlaylistByID(context.Background(), 0)
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidID)
	}
}

func TestGetPlaylistByID(t *testing.T) {
	t.Parallel()

	playlistUseCase := NewPlaylistUseCase(&fakePlaylistRepository{
		getPlaylistByIDFn: func(_ context.Context, id int64) (*entity.Playlist, error) {
			return &entity.Playlist{ID: id, Name: "Lo-Fi Hip Hop Mix"}, nil
		},
	}, nil)

	playlist, err := playlistUseCase.GetPlaylistByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetPlaylistByID() error = %v", err)
	}

	if playlist.ID != 1 {
		t.Fatalf("ID = %d, want %d", playlist.ID, 1)
	}
}
