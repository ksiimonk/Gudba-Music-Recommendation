package usecase

import (
	"context"
	"errors"

	"music-recommender-backend/internal/entity"
)

var (
	ErrInvalidEntityType = errors.New("entity_type must be 'track' or 'playlist'")
	ErrInvalidEventType  = errors.New("event_type must be 'play', 'like', 'dislike', 'skip', or 'open'")
)

var validEntityTypes = map[string]bool{
	"track":    true,
	"playlist": true,
}

var validEventTypes = map[string]bool{
	"play":    true,
	"like":    true,
	"dislike": true,
	"skip":    true,
	"open":    true,
}

type eventRepository interface {
	CreateEvent(ctx context.Context, event *entity.Event) error
}

type EventUseCase struct {
	eventRepository eventRepository
}

func NewEventUseCase(eventRepository eventRepository) *EventUseCase {
	return &EventUseCase{eventRepository: eventRepository}
}

func (u *EventUseCase) RecordEvent(ctx context.Context, userID int64, entityType, entityEvent string, entityID int64) error {
	if u == nil || u.eventRepository == nil {
		return errors.New("event use case repository is nil")
	}

	if !validEntityTypes[entityType] {
		return ErrInvalidEntityType
	}

	if !validEventTypes[entityEvent] {
		return ErrInvalidEventType
	}

	if entityID <= 0 {
		return ErrInvalidID
	}

	event := &entity.Event{
		UserID:     userID,
		EntityType: entityType,
		EntityID:   entityID,
		EventType:  entityEvent,
	}

	return u.eventRepository.CreateEvent(ctx, event)
}

func (u *EventUseCase) RecordTrackPlay(ctx context.Context, userID int64, trackID int64) error {
	return u.RecordEvent(ctx, userID, "track", "play", trackID)
}

func (u *EventUseCase) RecordTrackLike(ctx context.Context, userID int64, trackID int64) error {
	return u.RecordEvent(ctx, userID, "track", "like", trackID)
}

func (u *EventUseCase) RecordTrackDislike(ctx context.Context, userID int64, trackID int64) error {
	return u.RecordEvent(ctx, userID, "track", "dislike", trackID)
}

func (u *EventUseCase) RecordTrackSkip(ctx context.Context, userID int64, trackID int64) error {
	return u.RecordEvent(ctx, userID, "track", "skip", trackID)
}

func (u *EventUseCase) RecordPlaylistOpen(ctx context.Context, userID int64, playlistID int64) error {
	return u.RecordEvent(ctx, userID, "playlist", "open", playlistID)
}
