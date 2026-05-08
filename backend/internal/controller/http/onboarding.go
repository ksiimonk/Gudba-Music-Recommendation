package http

import (
	"context"
	"errors"
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/entity"
	"music-recommender-backend/internal/repository"
	"music-recommender-backend/internal/usecase"
)

type onboardingUseCase interface {
	SaveOnboarding(ctx context.Context, userID int64, request *entity.OnboardingRequest) error
	GetProfile(ctx context.Context, userID int64) (*entity.UserProfile, error)
}

type OnboardingHandler struct {
	onboardingUseCase onboardingUseCase
}

func NewOnboardingHandler(onboardingUseCase onboardingUseCase) *OnboardingHandler {
	return &OnboardingHandler{onboardingUseCase: onboardingUseCase}
}

func (h *OnboardingHandler) SaveOnboarding(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		c.JSON(nethttp.StatusUnauthorized, gin.H{
			"message": "authorization required",
		})
		return
	}

	var request entity.OnboardingRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	if err := h.onboardingUseCase.SaveOnboarding(c.Request.Context(), claims.UserID, &request); err != nil {
		switch {
		case errors.Is(err, usecase.ErrEmptyOnboarding):
			c.JSON(nethttp.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		case errors.Is(err, usecase.ErrInvalidID):
			c.JSON(nethttp.StatusBadRequest, gin.H{
				"message": "all ids must be greater than zero",
			})
		default:
			c.JSON(nethttp.StatusInternalServerError, gin.H{
				"message": "failed to save preferences",
			})
		}
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"message": "preferences saved",
	})
}

func (h *OnboardingHandler) GetProfile(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		c.JSON(nethttp.StatusUnauthorized, gin.H{
			"message": "authorization required",
		})
		return
	}

	profile, err := h.onboardingUseCase.GetProfile(c.Request.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrProfileNotFound) {
			c.JSON(nethttp.StatusOK, gin.H{
				"profile": entity.UserProfile{
					UserID:            claims.UserID,
					FavoriteGenreIDs:  []int64{},
					FavoriteArtistIDs: []int64{},
					StarterTrackIDs:   []int64{},
					Contexts:          []string{},
				},
			})
			return
		}

		c.JSON(nethttp.StatusInternalServerError, gin.H{
			"message": "failed to load profile",
		})
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"profile": profile,
	})
}
