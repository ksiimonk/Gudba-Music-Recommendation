package http

import (
	"context"
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/entity"
)

type analyticsUseCase interface {
	GetStats(ctx context.Context) (*entity.Stats, error)
	GetRecommendationMetrics(ctx context.Context) (*entity.RecommendationMetrics, error)
}

type AnalyticsHandler struct {
	analyticsUseCase analyticsUseCase
}

func NewAnalyticsHandler(analyticsUseCase analyticsUseCase) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsUseCase: analyticsUseCase}
}

func (h *AnalyticsHandler) GetStats(c *gin.Context) {
	stats, err := h.analyticsUseCase.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"message": "failed to load stats"})
		return
	}

	c.JSON(nethttp.StatusOK, stats)
}

func (h *AnalyticsHandler) GetRecommendationMetrics(c *gin.Context) {
	metrics, err := h.analyticsUseCase.GetRecommendationMetrics(c.Request.Context())
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"message": "failed to load metrics"})
		return
	}

	c.JSON(nethttp.StatusOK, metrics)
}
