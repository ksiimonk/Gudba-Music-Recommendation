package usecase

import (
	"context"
	"errors"
	"testing"

	"music-recommender-backend/internal/entity"
)

type fakeEventRepository struct {
	createEventFn func(ctx context.Context, event *entity.Event) error
}

func (f *fakeEventRepository) CreateEvent(ctx context.Context, event *entity.Event) error {
	if f.createEventFn != nil {
		return f.createEventFn(ctx, event)
	}

	return nil
}

func TestRecordTrackPlaySuccess(t *testing.T) {
	t.Parallel()

	var savedEvent *entity.Event

	eventUseCase := NewEventUseCase(&fakeEventRepository{
		createEventFn: func(_ context.Context, event *entity.Event) error {
			saved := *event
			savedEvent = &saved
			return nil
		},
	}, nil)

	err := eventUseCase.RecordTrackPlay(context.Background(), 42, 1)
	if err != nil {
		t.Fatalf("RecordTrackPlay() error = %v", err)
	}

	if savedEvent == nil {
		t.Fatal("event was not saved")
	}

	if savedEvent.UserID != 42 || savedEvent.EntityType != "track" || savedEvent.EntityID != 1 || savedEvent.EventType != "play" {
		t.Fatalf("event = %+v, want user_id=42 track/1 play", savedEvent)
	}
}

func TestRecordTrackLikeSuccess(t *testing.T) {
	t.Parallel()

	eventUseCase := NewEventUseCase(&fakeEventRepository{
		createEventFn: func(_ context.Context, event *entity.Event) error {
			return nil
		},
	}, nil)
	err := eventUseCase.RecordTrackLike(context.Background(), 42, 1)
	if err != nil {
		t.Fatalf("RecordTrackLike() error = %v", err)
	}
}

func TestRecordEventInvalidEntityType(t *testing.T) {
	t.Parallel()

	eventUseCase := NewEventUseCase(&fakeEventRepository{}, nil)

	err := eventUseCase.RecordEvent(context.Background(), 1, "album", "play", 1)
	if !errors.Is(err, ErrInvalidEntityType) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidEntityType)
	}
}

func TestRecordEventInvalidEventType(t *testing.T) {
	t.Parallel()

	eventUseCase := NewEventUseCase(&fakeEventRepository{}, nil)

	err := eventUseCase.RecordEvent(context.Background(), 1, "track", "share", 1)
	if !errors.Is(err, ErrInvalidEventType) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidEventType)
	}
}

func TestRecordEventInvalidID(t *testing.T) {
	t.Parallel()

	eventUseCase := NewEventUseCase(&fakeEventRepository{}, nil)

	err := eventUseCase.RecordEvent(context.Background(), 1, "track", "play", 0)
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidID)
	}
}
