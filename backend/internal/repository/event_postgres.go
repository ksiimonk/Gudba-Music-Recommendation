package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"music-recommender-backend/internal/entity"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) CreateEvent(ctx context.Context, event *entity.Event) error {
	if r == nil || r.db == nil {
		return errors.New("event repository database is nil")
	}

	if event == nil {
		return errors.New("event is nil")
	}

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO events (user_id, entity_type, entity_id, event_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`, event.UserID, event.EntityType, event.EntityID, event.EventType).Scan(&event.ID, &event.CreatedAt)

	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}

	return nil
}
