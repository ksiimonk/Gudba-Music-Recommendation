package usecase

import (
	"context"
	"errors"

	"music-recommender-backend/internal/entity"
)

type playlistRepository interface {
	ListPlaylists(ctx context.Context) ([]entity.Playlist, error)
	GetPlaylistByID(ctx context.Context, id int64) (*entity.Playlist, error)
}

type PlaylistUseCase struct {
	playlistRepository playlistRepository
}

func NewPlaylistUseCase(playlistRepository playlistRepository) *PlaylistUseCase {
	return &PlaylistUseCase{playlistRepository: playlistRepository}
}

func (u *PlaylistUseCase) ListPlaylists(ctx context.Context) ([]entity.Playlist, error) {
	if u == nil || u.playlistRepository == nil {
		return nil, errors.New("playlist use case repository is nil")
	}

	return u.playlistRepository.ListPlaylists(ctx)
}

func (u *PlaylistUseCase) GetPlaylistByID(ctx context.Context, id int64) (*entity.Playlist, error) {
	if u == nil || u.playlistRepository == nil {
		return nil, errors.New("playlist use case repository is nil")
	}

	if id <= 0 {
		return nil, ErrInvalidID
	}

	return u.playlistRepository.GetPlaylistByID(ctx, id)
}
