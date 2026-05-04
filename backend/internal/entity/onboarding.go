package entity

import "time"

type OnboardingRequest struct {
	GenreIDs  []int64  `json:"genre_ids"`
	ArtistIDs []int64  `json:"artist_ids"`
	TrackIDs  []int64  `json:"track_ids"`
	Contexts  []string `json:"contexts"`
}

type UserProfile struct {
	UserID            int64     `json:"user_id"`
	FavoriteGenreIDs  []int64   `json:"favorite_genre_ids"`
	FavoriteArtistIDs []int64   `json:"favorite_artist_ids"`
	StarterTrackIDs   []int64   `json:"starter_track_ids"`
	Contexts          []string  `json:"contexts"`
	UpdatedAt         time.Time `json:"updated_at"`
}
