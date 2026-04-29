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

type trackUseCase interface {
	ListTracks(ctx context.Context) ([]entity.Track, error)
	GetTrackByID(ctx context.Context, id int64) (*entity.Track, error)
}

type TrackHandler struct {
	trackUseCase trackUseCase
}

func NewTrackHandler(trackUseCase trackUseCase) *TrackHandler {
	return &TrackHandler{trackUseCase: trackUseCase}
}

func (h *TrackHandler) ListTracks(c *gin.Context) {
	tracks, err := h.trackUseCase.ListTracks(c.Request.Context())
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{
			"message": "failed to load tracks",
		})
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"tracks": tracks,
	})
}

func (h *TrackHandler) GetTrack(c *gin.Context) {
	id, ok := parseTrackID(c.Param("id"))
	if !ok {
		c.JSON(nethttp.StatusBadRequest, gin.H{
			"message": "invalid track id",
		})
		return
	}

	track, err := h.trackUseCase.GetTrackByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidID):
			c.JSON(nethttp.StatusBadRequest, gin.H{
				"message": "invalid track id",
			})
		case errors.Is(err, repository.ErrTrackNotFound):
			c.JSON(nethttp.StatusNotFound, gin.H{
				"message": "track not found",
			})
		default:
			c.JSON(nethttp.StatusInternalServerError, gin.H{
				"message": "failed to load track",
			})
		}
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"track": track,
	})
}

func parseTrackID(value string) (int64, bool) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
