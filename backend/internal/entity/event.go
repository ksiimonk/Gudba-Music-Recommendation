package entity

import "time"

type Event struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	EntityType string    `json:"entity_type"`
	EntityID   int64     `json:"entity_id"`
	EventType  string    `json:"event_type"`
	CreatedAt  time.Time `json:"created_at"`
}
