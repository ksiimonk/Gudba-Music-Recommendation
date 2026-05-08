package http

import (
	"context"
	"errors"
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/usecase"
)

type eventUseCase interface {
	RecordEvent(ctx context.Context, userID int64, entityType, eventType string, entityID int64) error
	RecordTrackPlay(ctx context.Context, userID int64, trackID int64) error
	RecordTrackLike(ctx context.Context, userID int64, trackID int64) error
	RecordTrackDislike(ctx context.Context, userID int64, trackID int64) error
	RecordTrackSkip(ctx context.Context, userID int64, trackID int64) error
	RecordPlaylistOpen(ctx context.Context, userID int64, playlistID int64) error
	ToggleLike(ctx context.Context, userID, trackID int64) (liked bool, err error)
}

type EventHandler struct {
	eventUseCase eventUseCase
}

func NewEventHandler(eventUseCase eventUseCase) *EventHandler {
	return &EventHandler{eventUseCase: eventUseCase}
}

type genericEventRequest struct {
	EntityType string `json:"entity_type" binding:"required"`
	EntityID   int64  `json:"entity_id" binding:"required"`
	EventType  string `json:"event_type" binding:"required"`
}

func (h *EventHandler) RecordEvent(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		c.JSON(nethttp.StatusUnauthorized, gin.H{"message": "authorization required"})
		return
	}

	var req genericEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	if err := h.eventUseCase.RecordEvent(c.Request.Context(), claims.UserID, req.EntityType, req.EventType, req.EntityID); err != nil {
		writeEventError(c, err)
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{"message": "event saved"})
}

func (h *EventHandler) RecordTrackPlay(c *gin.Context) {
	h.recordTrackEvent(c, "play")
}

func (h *EventHandler) ToggleTrackLike(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		c.JSON(nethttp.StatusUnauthorized, gin.H{"message": "authorization required"})
		return
	}

	id, ok := parseTrackID(c.Param("id"))
	if !ok {
		c.JSON(nethttp.StatusBadRequest, gin.H{"message": "invalid track id"})
		return
	}

	liked, err := h.eventUseCase.ToggleLike(c.Request.Context(), claims.UserID, id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidID):
			c.JSON(nethttp.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(nethttp.StatusInternalServerError, gin.H{"message": "failed to toggle like"})
		}
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{"liked": liked})
}

func (h *EventHandler) RecordTrackLike(c *gin.Context) {
	h.recordTrackEvent(c, "like")
}

func (h *EventHandler) RecordTrackDislike(c *gin.Context) {
	h.recordTrackEvent(c, "dislike")
}

func (h *EventHandler) RecordTrackSkip(c *gin.Context) {
	h.recordTrackEvent(c, "skip")
}

func (h *EventHandler) RecordPlaylistOpen(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		c.JSON(nethttp.StatusUnauthorized, gin.H{"message": "authorization required"})
		return
	}

	id, ok := parsePlaylistID(c.Param("id"))
	if !ok {
		c.JSON(nethttp.StatusBadRequest, gin.H{"message": "invalid playlist id"})
		return
	}

	if err := h.eventUseCase.RecordPlaylistOpen(c.Request.Context(), claims.UserID, id); err != nil {
		writeEventError(c, err)
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{"message": "event saved"})
}

func (h *EventHandler) recordTrackEvent(c *gin.Context, eventType string) {
	claims, ok := getAuthClaims(c)
	if !ok {
		c.JSON(nethttp.StatusUnauthorized, gin.H{"message": "authorization required"})
		return
	}

	id, ok := parseTrackID(c.Param("id"))
	if !ok {
		c.JSON(nethttp.StatusBadRequest, gin.H{"message": "invalid track id"})
		return
	}

	var err error
	switch eventType {
	case "play":
		err = h.eventUseCase.RecordTrackPlay(c.Request.Context(), claims.UserID, id)
	case "like":
		err = h.eventUseCase.RecordTrackLike(c.Request.Context(), claims.UserID, id)
	case "dislike":
		err = h.eventUseCase.RecordTrackDislike(c.Request.Context(), claims.UserID, id)
	case "skip":
		err = h.eventUseCase.RecordTrackSkip(c.Request.Context(), claims.UserID, id)
	}

	if err != nil {
		writeEventError(c, err)
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{"message": "event saved"})
}

func writeEventError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidEntityType), errors.Is(err, usecase.ErrInvalidEventType), errors.Is(err, usecase.ErrInvalidID):
		c.JSON(nethttp.StatusBadRequest, gin.H{"message": err.Error()})
	default:
		c.JSON(nethttp.StatusInternalServerError, gin.H{"message": "failed to save event"})
	}
}
