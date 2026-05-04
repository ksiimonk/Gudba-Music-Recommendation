package usecase

import (
	"context"
	"errors"
	"math"
	"sort"

	"music-recommender-backend/internal/entity"
)

type recommendationRepository interface {
	GetUserProfile(ctx context.Context, userID int64) (*entity.UserProfile, error)
	GetRecentEvents(ctx context.Context, userID int64, limit int) ([]entity.Event, error)
	ListCandidateTracks(ctx context.Context) ([]entity.Track, error)
	ListPlaylistsWithTracks(ctx context.Context) ([]entity.Playlist, error)
	SaveRecommendationRequest(ctx context.Context, userID int64, requestType string) (int64, error)
	SaveImpressions(ctx context.Context, impressions []entity.RecommendationImpression) ([]int64, error)
	SaveFactors(ctx context.Context, factors []entity.RecommendationFactor) error
}

type RecommendationUseCase struct {
	recommendationRepository recommendationRepository
}

func NewRecommendationUseCase(recommendationRepository recommendationRepository) *RecommendationUseCase {
	return &RecommendationUseCase{recommendationRepository: recommendationRepository}
}

func (u *RecommendationUseCase) GetTrackRecommendations(ctx context.Context, userID int64, limit int) ([]entity.TrackRecommendation, error) {
	if u == nil || u.recommendationRepository == nil {
		return nil, errors.New("recommendation use case repository is nil")
	}

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	profile, _ := u.recommendationRepository.GetUserProfile(ctx, userID)
	events, _ := u.recommendationRepository.GetRecentEvents(ctx, userID, 200)
	allTracks, err := u.recommendationRepository.ListCandidateTracks(ctx)
	if err != nil {
		return nil, err
	}

	likedGenreIDs := buildLikedGenreIDs(events, allTracks)
	playedIDs, skippedIDs := buildSkippedIDs(events)

	scoredTracks := scoreCandidates(allTracks, profile, likedGenreIDs)

	sort.Slice(scoredTracks, func(i, j int) bool {
		return scoredTracks[i].Score > scoredTracks[j].Score
	})

	scoredTracks = filterSkipped(scoredTracks, skippedIDs)

	scoredTracks = downrankPlayed(scoredTracks, playedIDs)

	reranked := rerankForDiversity(scoredTracks, limit)

	topN := reranked
	if len(topN) > limit {
		topN = topN[:limit]
	}

	requestID, err := u.recommendationRepository.SaveRecommendationRequest(ctx, userID, "tracks")
	if err != nil {
		return nil, err
	}

	impressions := make([]entity.RecommendationImpression, 0, len(topN))
	for i, rec := range topN {
		impressions = append(impressions, entity.RecommendationImpression{
			RequestID:   requestID,
			UserID:      userID,
			EntityType:  "track",
			EntityID:    rec.Track.ID,
			Rank:        i + 1,
			Score:       rec.Score,
			Explanation: rec.Explanation,
		})
	}

	impressionIDs, err := u.recommendationRepository.SaveImpressions(ctx, impressions)
	if err != nil {
		return toTrackRecommendations(topN), nil
	}

	for i, impID := range impressionIDs {
		if i >= len(topN) {
			break
		}

		var factors []entity.RecommendationFactor
		if topN[i].Score > 0 {
			genreContrib := math.Min(topN[i].Score, 40)
			remaining := topN[i].Score - genreContrib
			artistContrib := math.Min(remaining, 25)
			remaining -= artistContrib
			popularityContrib := math.Min(remaining, 20)
			remaining -= popularityContrib
			interestContrib := remaining

			addFactor := func(name string, value float64) {
				weight := 0.0
				switch name {
				case "genre_match":
					weight = 0.40
				case "artist_match":
					weight = 0.25
				case "popularity":
					weight = 0.20
				case "recent_interest":
					weight = 0.15
				}

				factors = append(factors, entity.RecommendationFactor{
					ImpressionID: impID,
					FactorName:   name,
					FactorValue:  value,
					Weight:       weight,
					Contribution: contribution(name, value, genreContrib, artistContrib, popularityContrib, interestContrib),
				})
			}

			addFactor("genre_match", genreContrib)
			addFactor("artist_match", artistContrib)
			addFactor("popularity", popularityContrib)
			addFactor("recent_interest", interestContrib)
		}

		if len(factors) > 0 {
			u.recommendationRepository.SaveFactors(ctx, factors)
		}
	}

	return toTrackRecommendations(topN), nil
}

func (u *RecommendationUseCase) GetPlaylistRecommendations(ctx context.Context, userID int64, limit int) ([]entity.PlaylistRecommendation, error) {
	if u == nil || u.recommendationRepository == nil {
		return nil, errors.New("recommendation use case repository is nil")
	}

	if limit <= 0 || limit > 20 {
		limit = 6
	}

	profile, _ := u.recommendationRepository.GetUserProfile(ctx, userID)
	events, _ := u.recommendationRepository.GetRecentEvents(ctx, userID, 200)
	allTracks, _ := u.recommendationRepository.ListCandidateTracks(ctx)
	allPlaylists, err := u.recommendationRepository.ListPlaylistsWithTracks(ctx)
	if err != nil {
		return nil, err
	}

	likedGenreIDs := buildLikedGenreIDs(events, allTracks)

	var scored []entity.PlaylistRecommendation

	for _, pl := range allPlaylists {
		pProfile := buildPlaylistProfile(pl)
		score, explanation := scorePlaylist(pProfile, profile, likedGenreIDs)

		if score <= 0 {
			continue
		}

		scored = append(scored, entity.PlaylistRecommendation{
			Playlist:    pl,
			Score:       math.Round(score*100) / 100,
			Explanation: explanation,
		})
	}

	if len(scored) == 0 {
		return []entity.PlaylistRecommendation{}, nil
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	if len(scored) > limit {
		scored = scored[:limit]
	}

	requestID, err := u.recommendationRepository.SaveRecommendationRequest(ctx, userID, "playlists")
	if err != nil {
		return scored, nil
	}

	impressions := make([]entity.RecommendationImpression, 0, len(scored))
	for i, rec := range scored {
		impressions = append(impressions, entity.RecommendationImpression{
			RequestID:   requestID,
			UserID:      userID,
			EntityType:  "playlist",
			EntityID:    rec.Playlist.ID,
			Rank:        i + 1,
			Score:       rec.Score,
			Explanation: rec.Explanation,
		})
	}

	u.recommendationRepository.SaveImpressions(ctx, impressions)

	return scored, nil
}

func toTrackRecommendations(scored []scoredTrack) []entity.TrackRecommendation {
	result := make([]entity.TrackRecommendation, len(scored))
	for i, s := range scored {
		result[i] = entity.TrackRecommendation{
			Track:       s.Track,
			Score:       s.Score,
			Explanation: s.Explanation,
		}
	}
	return result
}

func contribution(name string, value, genreContrib, artistContrib, popularityContrib, interestContrib float64) float64 {
	total := genreContrib + artistContrib + popularityContrib + interestContrib
	if total == 0 {
		return 0
	}

	switch name {
	case "genre_match":
		return value / total
	case "artist_match":
		return value / total
	case "popularity":
		return value / total
	case "recent_interest":
		return value / total
	}

	return 0
}

func buildLikedGenreIDs(events []entity.Event, tracks []entity.Track) map[int64]bool {
	trackLikes := make(map[int64]bool)
	for _, e := range events {
		if e.EventType == "like" && e.EntityType == "track" {
			trackLikes[e.EntityID] = true
		}
	}

	genreIDs := make(map[int64]bool)
	for _, t := range tracks {
		if trackLikes[t.ID] {
			for _, g := range t.Genres {
				genreIDs[g.ID] = true
			}
		}
	}

	return genreIDs
}

func buildSkippedIDs(events []entity.Event) (map[int64]bool, map[int64]bool) {
	skipped := make(map[int64]bool)
	played := make(map[int64]bool)

	for _, e := range events {
		if e.EntityType != "track" {
			continue
		}

		switch e.EventType {
		case "skip":
			skipped[e.EntityID] = true
		case "play":
			played[e.EntityID] = true
		}
	}

	return played, skipped
}

func filterSkipped(tracks []scoredTrack, skippedIDs map[int64]bool) []scoredTrack {
	var filtered []scoredTrack

	for _, t := range tracks {
		if !skippedIDs[t.Track.ID] {
			filtered = append(filtered, t)
		}
	}

	return filtered
}

func downrankPlayed(tracks []scoredTrack, playedIDs map[int64]bool) []scoredTrack {
	for i := range tracks {
		if playedIDs[tracks[i].Track.ID] {
			tracks[i].Score *= 0.7
		}
	}

	return tracks
}

func rerankForDiversity(tracks []scoredTrack, limit int) []scoredTrack {
	if len(tracks) <= 1 {
		return tracks
	}

	artistCount := make(map[int64]int)
	genreCount := make(map[int64]int)
	var result []scoredTrack

	for _, t := range tracks {
		artistCount[t.Track.Artist.ID]++

		for _, g := range t.Track.Genres {
			if artistCount[t.Track.Artist.ID] > 2 {
				continue
			}

			if genreCount[g.ID] >= 3 {
				continue
			}
		}

		if artistCount[t.Track.Artist.ID] > 2 {
			continue
		}

		maxGenreCount := 0
		for _, g := range t.Track.Genres {
			if genreCount[g.ID] > maxGenreCount {
				maxGenreCount = genreCount[g.ID]
			}
		}

		if maxGenreCount >= 3 {
			continue
		}

		result = append(result, t)
		artistCount[t.Track.Artist.ID]--
		artistCount[t.Track.Artist.ID]++

		for _, g := range t.Track.Genres {
			genreCount[g.ID]++
		}

		if len(result) >= limit*2 {
			break
		}
	}

	if len(result) == 0 {
		result = tracks
		if len(result) > limit {
			result = result[:limit]
		}
	}

	return result
}
