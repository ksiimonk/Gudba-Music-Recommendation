package http

import (
	"context"
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/entity"
)

type genreUseCase interface {
	ListGenres(ctx context.Context) ([]entity.Genre, error)
}

type GenreHandler struct {
	genreUseCase genreUseCase
}

func NewGenreHandler(genreUseCase genreUseCase) *GenreHandler {
	return &GenreHandler{genreUseCase: genreUseCase}
}

func (h *GenreHandler) ListGenres(c *gin.Context) {
	genres, err := h.genreUseCase.ListGenres(c.Request.Context())
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{
			"message": "failed to load genres",
		})
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"genres": genres,
	})
}
