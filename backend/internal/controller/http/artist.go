package http

import (
	"context"
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/entity"
)

type artistUseCase interface {
	ListArtists(ctx context.Context) ([]entity.Artist, error)
}

type ArtistHandler struct {
	artistUseCase artistUseCase
}

func NewArtistHandler(artistUseCase artistUseCase) *ArtistHandler {
	return &ArtistHandler{artistUseCase: artistUseCase}
}

func (h *ArtistHandler) ListArtists(c *gin.Context) {
	artists, err := h.artistUseCase.ListArtists(c.Request.Context())
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{
			"message": "failed to load artists",
		})
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"artists": artists,
	})
}
