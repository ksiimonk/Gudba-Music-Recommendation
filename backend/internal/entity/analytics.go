package entity

type Stats struct {
	UsersCount                 int64 `json:"users_count"`
	TracksCount                int64 `json:"tracks_count"`
	PlaylistsCount             int64 `json:"playlists_count"`
	EventsCount                int64 `json:"events_count"`
	LikesCount                 int64 `json:"likes_count"`
	DislikesCount              int64 `json:"dislikes_count"`
	SkipsCount                 int64 `json:"skips_count"`
	PlaysCount                 int64 `json:"plays_count"`
	RecommendationRequestsCount int64 `json:"recommendation_requests_count"`
	RecommendationImpressionsCount int64 `json:"recommendation_impressions_count"`
}

type RecommendationMetrics struct {
	LikeRate            float64 `json:"like_rate"`
	SkipRate            float64 `json:"skip_rate"`
	DislikeRate         float64 `json:"dislike_rate"`
	CoverageTracksCount int64   `json:"coverage_tracks_count"`
	CoverageTracksPercent float64 `json:"coverage_tracks_percent"`
}
