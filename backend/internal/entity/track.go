package entity

import "time"

type Artist struct {
	ID         int64     `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	ImageURL   string    `json:"image_url,omitempty" db:"image_url"`
	Bio        string    `json:"bio,omitempty" db:"bio"`
	SpotifyURL string    `json:"spotify_url,omitempty" db:"spotify_url"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Genre struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

type Track struct {
	ID              int64     `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	ArtistID        int64     `json:"artist_id" db:"artist_id"`
	Artist          Artist    `json:"artist"`
	Genres          []Genre   `json:"genres"`
	DurationMS      int       `json:"duration_ms" db:"duration_ms"`
	PreviewURL      string    `json:"preview_url,omitempty" db:"preview_url"`
	SpotifyURL      string    `json:"spotify_url,omitempty" db:"spotify_url"`
	CoverURL        string    `json:"cover_url,omitempty" db:"cover_url"`
	PopularityScore int       `json:"popularity_score" db:"popularity_score"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
