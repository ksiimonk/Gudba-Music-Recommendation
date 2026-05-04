package usecase

import (
	"math"

	"music-recommender-backend/internal/entity"
)

type scoredTrack struct {
	Track       entity.Track
	Score       float64
	Explanation string
}

func scoreCandidates(tracks []entity.Track, profile *entity.UserProfile, likedGenreIDs map[int64]bool) []scoredTrack {
	if profile == nil || (len(profile.FavoriteGenreIDs) == 0 && len(profile.FavoriteArtistIDs) == 0) {
		return coldStartScores(tracks)
	}

	favGenres := makeSet(profile.FavoriteGenreIDs)
	favArtists := makeSet(profile.FavoriteArtistIDs)

	result := make([]scoredTrack, 0, len(tracks))

	for _, track := range tracks {
		genreScore := calcGenreScore(track, favGenres)
		artistScore := calcArtistScore(track, favArtists)
		popularityScore := calcPopularityScore(track)
		interestScore := calcInterestScore(track, likedGenreIDs)

		total := genreScore + artistScore + popularityScore + interestScore

		explanation := buildExplanation(genreScore, artistScore, popularityScore, interestScore)

		result = append(result, scoredTrack{
			Track:       track,
			Score:       math.Round(total*100) / 100,
			Explanation: explanation,
		})
	}

	return result
}

func coldStartScores(tracks []entity.Track) []scoredTrack {
	result := make([]scoredTrack, 0, len(tracks))

	for _, track := range tracks {
		score := calcPopularityScore(track)

		result = append(result, scoredTrack{
			Track:       track,
			Score:       math.Round(score*100) / 100,
			Explanation: "Популярен среди слушателей",
		})
	}

	return result
}

func calcGenreScore(track entity.Track, favGenres map[int64]bool) float64 {
	if len(favGenres) == 0 || len(track.Genres) == 0 {
		return 0
	}

	matchCount := 0
	for _, g := range track.Genres {
		if favGenres[g.ID] {
			matchCount++
		}
	}

	if matchCount == 0 {
		return 0
	}

	return (float64(matchCount) / float64(len(track.Genres))) * 40
}

func calcArtistScore(track entity.Track, favArtists map[int64]bool) float64 {
	if favArtists[track.Artist.ID] {
		return 25
	}

	return 0
}

func calcPopularityScore(track entity.Track) float64 {
	return (float64(track.PopularityScore) / 100) * 20
}

func calcInterestScore(track entity.Track, likedGenreIDs map[int64]bool) float64 {
	if len(likedGenreIDs) == 0 {
		return 0
	}

	for _, g := range track.Genres {
		if likedGenreIDs[g.ID] {
			return 15
		}
	}

	return 0
}

func buildExplanation(genreScore, artistScore, popularityScore, interestScore float64) string {
	if genreScore > 0 && artistScore > 0 {
		return "Совпадает с твоими любимыми жанрами и артистами"
	}

	if genreScore > 0 {
		return "Совпадает с твоими любимыми жанрами"
	}

	if artistScore > 0 {
		return "Похож на артистов из твоего профиля"
	}

	if interestScore > 0 {
		return "Учитывает твои недавние лайки"
	}

	return "Популярен в выбранных жанрах"
}

func makeSet(ids []int64) map[int64]bool {
	set := make(map[int64]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}

	return set
}
