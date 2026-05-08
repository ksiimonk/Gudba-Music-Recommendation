package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"music-recommender-backend/internal/entity"
)

var (
	ErrInvalidEntityType = errors.New("entity_type must be 'track' or 'playlist'")
	ErrInvalidEventType  = errors.New("event_type must be 'play', 'like', 'unlike', 'dislike', 'skip', or 'open'")
)

var validEntityTypes = map[string]bool{
	"track":    true,
	"playlist": true,
}

var validEventTypes = map[string]bool{
	"play":    true,
	"like":    true,
	"unlike":  true,
	"dislike": true,
	"skip":    true,
	"open":    true,
}

var ErrPlaylistNotFoundForUser = errors.New("favorites playlist not found for user")

type eventRepository interface {
	CreateEvent(ctx context.Context, event *entity.Event) error
}

type playlistToggleRepository interface {
	GetFavoritesPlaylistByUserID(ctx context.Context, userID int64) (*entity.Playlist, error)
	AddTrackToPlaylist(ctx context.Context, playlistID, trackID int64) error
	RemoveTrackFromPlaylist(ctx context.Context, playlistID, trackID int64) error
	IsTrackFavorited(ctx context.Context, userID, trackID int64) (bool, error)
}

type EventUseCase struct {
	eventRepository          eventRepository
	playlistToggleRepository playlistToggleRepository
}

func NewEventUseCase(eventRepository eventRepository, playlistToggleRepository playlistToggleRepository) *EventUseCase {
	return &EventUseCase{eventRepository: eventRepository, playlistToggleRepository: playlistToggleRepository}
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

func (u *EventUseCase) ToggleLike(ctx context.Context, userID, trackID int64) (liked bool, err error) {
	if u == nil || u.eventRepository == nil {
		return false, errors.New("event use case repository is nil")
	}

	if trackID <= 0 {
		return false, ErrInvalidID
	}

	alreadyFavorited, err := u.playlistToggleRepository.IsTrackFavorited(ctx, userID, trackID)
	if err != nil {
		return false, fmt.Errorf("check favorited: %w", err)
	}

	favPlaylist, err := u.playlistToggleRepository.GetFavoritesPlaylistByUserID(ctx, userID)
	if err != nil {
		return false, ErrPlaylistNotFoundForUser
	}

	if alreadyFavorited {
		if err := u.playlistToggleRepository.RemoveTrackFromPlaylist(ctx, favPlaylist.ID, trackID); err != nil {
			return false, fmt.Errorf("remove from favorites: %w", err)
		}
		if err := u.RecordEvent(ctx, userID, "track", "unlike", trackID); err != nil {
			log.Printf("record unlike event: %v", err)
		}
		return false, nil
	}

	if err := u.playlistToggleRepository.AddTrackToPlaylist(ctx, favPlaylist.ID, trackID); err != nil {
		return false, fmt.Errorf("add to favorites: %w", err)
	}
	if err := u.RecordEvent(ctx, userID, "track", "like", trackID); err != nil {
		log.Printf("record like event: %v", err)
	}
	return true, nil
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
