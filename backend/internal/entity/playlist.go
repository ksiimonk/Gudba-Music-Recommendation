package entity

import "time"

type Playlist struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	TrackCount  int       `json:"track_count"`
	Tracks      []Track   `json:"tracks,omitempty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
