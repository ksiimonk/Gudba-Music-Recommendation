package http

import (
	"context"
	nethttp "net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/entity"
)

type recommendationUseCase interface {
	GetTrackRecommendations(ctx context.Context, userID int64, limit int) ([]entity.TrackRecommendation, error)
	GetPlaylistRecommendations(ctx context.Context, userID int64, limit int) ([]entity.PlaylistRecommendation, error)
}

type RecommendationHandler struct {
	recommendationUseCase recommendationUseCase
}

func NewRecommendationHandler(recommendationUseCase recommendationUseCase) *RecommendationHandler {
	return &RecommendationHandler{recommendationUseCase: recommendationUseCase}
}

func (h *RecommendationHandler) ListTrackRecommendations(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		c.JSON(nethttp.StatusUnauthorized, gin.H{"message": "authorization required"})
		return
	}

	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsed, err := strconv.Atoi(limitParam); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	recommendations, err := h.recommendationUseCase.GetTrackRecommendations(c.Request.Context(), claims.UserID, limit)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"message": "failed to load recommendations"})
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"recommendations": recommendations,
	})
}

func (h *RecommendationHandler) ListPlaylistRecommendations(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		c.JSON(nethttp.StatusUnauthorized, gin.H{"message": "authorization required"})
		return
	}

	limit := 6
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsed, err := strconv.Atoi(limitParam); err == nil && parsed > 0 && parsed <= 20 {
			limit = parsed
		}
	}

	recommendations, err := h.recommendationUseCase.GetPlaylistRecommendations(c.Request.Context(), claims.UserID, limit)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"message": "failed to load playlist recommendations"})
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"recommendations": recommendations,
	})
}
