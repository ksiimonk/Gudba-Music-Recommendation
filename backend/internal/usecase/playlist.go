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

type favoritesRepository interface {
	GetFavoritesPlaylistByUserID(ctx context.Context, userID int64) (*entity.Playlist, error)
	CreateFavoritesPlaylist(ctx context.Context, userID int64) (*entity.Playlist, error)
}

type PlaylistUseCase struct {
	playlistRepository  playlistRepository
	favoritesRepository favoritesRepository
}

func NewPlaylistUseCase(playlistRepository playlistRepository, favoritesRepository favoritesRepository) *PlaylistUseCase {
	return &PlaylistUseCase{playlistRepository: playlistRepository, favoritesRepository: favoritesRepository}
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

func (u *PlaylistUseCase) GetFavorites(ctx context.Context, userID int64) (*entity.Playlist, error) {
	if u == nil || u.favoritesRepository == nil {
		return nil, errors.New("playlist use case favorites repository is nil")
	}

	playlist, err := u.favoritesRepository.GetFavoritesPlaylistByUserID(ctx, userID)
	if err != nil {
		playlist, err = u.favoritesRepository.CreateFavoritesPlaylist(ctx, userID)
		if err != nil {
			return nil, err
		}
		return playlist, nil
	}

	return playlist, nil
}
