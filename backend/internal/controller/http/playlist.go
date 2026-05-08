package http

import (
	"context"
	"errors"
	nethttp "net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/entity"
	"music-recommender-backend/internal/repository"
	"music-recommender-backend/internal/usecase"
)

type playlistUseCase interface {
	ListPlaylists(ctx context.Context) ([]entity.Playlist, error)
	GetPlaylistByID(ctx context.Context, id int64) (*entity.Playlist, error)
	GetFavorites(ctx context.Context, userID int64) (*entity.Playlist, error)
}

type PlaylistHandler struct {
	playlistUseCase playlistUseCase
}

func NewPlaylistHandler(playlistUseCase playlistUseCase) *PlaylistHandler {
	return &PlaylistHandler{playlistUseCase: playlistUseCase}
}

func (h *PlaylistHandler) ListPlaylists(c *gin.Context) {
	playlists, err := h.playlistUseCase.ListPlaylists(c.Request.Context())
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{
			"message": "failed to load playlists",
		})
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"playlists": playlists,
	})
}

func (h *PlaylistHandler) GetPlaylist(c *gin.Context) {
	id, ok := parsePlaylistID(c.Param("id"))
	if !ok {
		c.JSON(nethttp.StatusBadRequest, gin.H{
			"message": "invalid playlist id",
		})
		return
	}

	playlist, err := h.playlistUseCase.GetPlaylistByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidID):
			c.JSON(nethttp.StatusBadRequest, gin.H{
				"message": "invalid playlist id",
			})
		case errors.Is(err, repository.ErrPlaylistNotFound):
			c.JSON(nethttp.StatusNotFound, gin.H{
				"message": "playlist not found",
			})
		default:
			c.JSON(nethttp.StatusInternalServerError, gin.H{
				"message": "failed to load playlist",
			})
		}
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"playlist": playlist,
	})
}

func (h *PlaylistHandler) GetFavorites(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		c.JSON(nethttp.StatusUnauthorized, gin.H{"message": "authorization required"})
		return
	}

	playlist, err := h.playlistUseCase.GetFavorites(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"message": "failed to load favorites"})
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{"playlist": playlist})
}

func parsePlaylistID(value string) (int64, bool) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
