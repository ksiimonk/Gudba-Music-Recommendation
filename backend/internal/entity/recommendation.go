package entity

import "time"

type RecommendationRequest struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	RequestType string    `json:"request_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type RecommendationImpression struct {
	ID          int64   `json:"id"`
	RequestID   int64   `json:"request_id"`
	UserID      int64   `json:"user_id"`
	EntityType  string  `json:"entity_type"`
	EntityID    int64   `json:"entity_id"`
	Rank        int     `json:"rank"`
	Score       float64 `json:"score"`
	Explanation string  `json:"explanation"`
}

type RecommendationFactor struct {
	ID           int64   `json:"id"`
	ImpressionID int64   `json:"impression_id"`
	FactorName   string  `json:"factor_name"`
	FactorValue  float64 `json:"factor_value"`
	Weight       float64 `json:"weight"`
	Contribution float64 `json:"contribution"`
}

type TrackRecommendation struct {
	Track       Track   `json:"track"`
	Score       float64 `json:"score"`
	Explanation string  `json:"explanation"`
	IsFavorited bool    `json:"is_favorited"`
}

type PlaylistRecommendation struct {
	Playlist    Playlist `json:"playlist"`
	Score       float64  `json:"score"`
	Explanation string   `json:"explanation"`
}
