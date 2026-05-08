package usecase

import "music-recommender-backend/internal/entity"

type playlistProfile struct {
	genreCounts  map[int64]int
	artistCounts map[int64]int
	avgPopularity float64
	totalTracks  int
}

func buildPlaylistProfile(playlist entity.Playlist) playlistProfile {
	genreCounts := make(map[int64]int)
	artistCounts := make(map[int64]int)
	totalPopularity := 0

	for _, track := range playlist.Tracks {
		artistCounts[track.Artist.ID]++

		for _, g := range track.Genres {
			genreCounts[g.ID]++
		}

		totalPopularity += track.PopularityScore
	}

	totalTracks := len(playlist.Tracks)
	avgPop := 0.0
	if totalTracks > 0 {
		avgPop = float64(totalPopularity) / float64(totalTracks)
	}

	return playlistProfile{
		genreCounts:   genreCounts,
		artistCounts:  artistCounts,
		avgPopularity: avgPop,
		totalTracks:   totalTracks,
	}
}

func scorePlaylist(profile playlistProfile, userProfile *entity.UserProfile, likedGenreIDs map[int64]bool) (float64, string) {
	if profile.totalTracks == 0 {
		return 0, ""
	}

	if userProfile == nil || (len(userProfile.FavoriteGenreIDs) == 0 && len(userProfile.FavoriteArtistIDs) == 0) {
		return calcPlaylistPopularity(profile.avgPopularity), "Популярный плейлист"
	}

	favGenres := makeSet(userProfile.FavoriteGenreIDs)
	favArtists := makeSet(userProfile.FavoriteArtistIDs)

	genreOverlap := calcPlaylistGenreOverlap(profile.genreCounts, favGenres)
	artistOverlap := calcPlaylistArtistOverlap(profile.artistCounts, favArtists)
	likedBoost := calcLikedGenreBoost(profile.genreCounts, likedGenreIDs)
	popularityScore := calcPlaylistPopularity(profile.avgPopularity)

	total := genreOverlap + artistOverlap + likedBoost + popularityScore

	explanation := buildPlaylistExplanation(genreOverlap, artistOverlap, likedBoost, popularityScore)

	return total, explanation
}

func calcPlaylistGenreOverlap(genreCounts map[int64]int, favGenres map[int64]bool) float64 {
	totalGenres := 0
	matchGenres := 0

	for gID := range genreCounts {
		totalGenres++
		if favGenres[gID] {
			matchGenres++
		}
	}

	if totalGenres == 0 {
		return 0
	}

	return (float64(matchGenres) / float64(totalGenres)) * 40
}

func calcPlaylistArtistOverlap(artistCounts map[int64]int, favArtists map[int64]bool) float64 {
	totalArtists := 0
	matchArtists := 0

	for aID := range artistCounts {
		totalArtists++
		if favArtists[aID] {
			matchArtists++
		}
	}

	if totalArtists == 0 {
		return 0
	}

	return (float64(matchArtists) / float64(totalArtists)) * 25
}

func calcLikedGenreBoost(genreCounts map[int64]int, likedGenreIDs map[int64]bool) float64 {
	if len(likedGenreIDs) == 0 {
		return 0
	}

	for gID := range genreCounts {
		if likedGenreIDs[gID] {
			return 20
		}
	}

	return 0
}

func calcPlaylistPopularity(avgPopularity float64) float64 {
	return (avgPopularity / 100) * 15
}

func buildPlaylistExplanation(genreOverlap, artistOverlap, likedBoost, popularityScore float64) string {
	if genreOverlap > 0 && artistOverlap > 0 {
		return "Совпадает с твоими любимыми жанрами и артистами"
	}

	if genreOverlap > 0 {
		return "Совпадает с твоими любимыми жанрами"
	}

	if artistOverlap > 0 {
		return "Содержит твоих любимых артистов"
	}

	if likedBoost > 0 {
		return "Содержит треки, похожие на твои лайки"
	}

	return "Популярный плейлист"
}
