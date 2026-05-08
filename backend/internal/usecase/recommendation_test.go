package usecase

import (
	"context"
	"errors"
	"testing"

	"music-recommender-backend/internal/entity"
)

type fakeRecommendationRepository struct {
	getUserProfileFn          func(ctx context.Context, userID int64) (*entity.UserProfile, error)
	getRecentEventsFn         func(ctx context.Context, userID int64, limit int) ([]entity.Event, error)
	listCandidateTracksFn     func(ctx context.Context) ([]entity.Track, error)
	listPlaylistsWithTracksFn func(ctx context.Context) ([]entity.Playlist, error)
	saveRequestFn             func(ctx context.Context, userID int64, requestType string) (int64, error)
	saveImpressionsFn         func(ctx context.Context, impressions []entity.RecommendationImpression) ([]int64, error)
	saveFactorsFn             func(ctx context.Context, factors []entity.RecommendationFactor) error
}

func (f *fakeRecommendationRepository) GetUserProfile(ctx context.Context, userID int64) (*entity.UserProfile, error) {
	if f.getUserProfileFn != nil {
		return f.getUserProfileFn(ctx, userID)
	}
	return nil, nil
}

func (f *fakeRecommendationRepository) GetRecentEvents(ctx context.Context, userID int64, limit int) ([]entity.Event, error) {
	if f.getRecentEventsFn != nil {
		return f.getRecentEventsFn(ctx, userID, limit)
	}
	return nil, nil
}

func (f *fakeRecommendationRepository) ListCandidateTracks(ctx context.Context) ([]entity.Track, error) {
	if f.listCandidateTracksFn != nil {
		return f.listCandidateTracksFn(ctx)
	}
	return nil, errors.New("no tracks")
}

func (f *fakeRecommendationRepository) GetFavoriteTrackIDs(ctx context.Context, userID int64) (map[int64]bool, error) {
	return nil, nil
}

func (f *fakeRecommendationRepository) ListPlaylistsWithTracks(ctx context.Context) ([]entity.Playlist, error) {
	if f.listPlaylistsWithTracksFn != nil {
		return f.listPlaylistsWithTracksFn(ctx)
	}
	return []entity.Playlist{}, nil
}

func (f *fakeRecommendationRepository) SaveRecommendationRequest(ctx context.Context, userID int64, requestType string) (int64, error) {
	if f.saveRequestFn != nil {
		return f.saveRequestFn(ctx, userID, requestType)
	}
	return 1, nil
}

func (f *fakeRecommendationRepository) SaveImpressions(ctx context.Context, impressions []entity.RecommendationImpression) ([]int64, error) {
	if f.saveImpressionsFn != nil {
		return f.saveImpressionsFn(ctx, impressions)
	}
	ids := make([]int64, len(impressions))
	for i := range impressions {
		ids[i] = int64(i + 1)
	}
	return ids, nil
}

func (f *fakeRecommendationRepository) SaveFactors(ctx context.Context, factors []entity.RecommendationFactor) error {
	if f.saveFactorsFn != nil {
		return f.saveFactorsFn(ctx, factors)
	}
	return nil
}

func TestGetTrackRecommendationsWithProfile(t *testing.T) {
	t.Parallel()

	recUseCase := NewRecommendationUseCase(&fakeRecommendationRepository{
		getUserProfileFn: func(_ context.Context, userID int64) (*entity.UserProfile, error) {
			return &entity.UserProfile{
				UserID:            userID,
				FavoriteGenreIDs:  []int64{1},
				FavoriteArtistIDs: []int64{1},
			}, nil
		},
		listCandidateTracksFn: func(_ context.Context) ([]entity.Track, error) {
			return []entity.Track{{
				ID:    1,
				Title: "Night Drive",
				Artist: entity.Artist{ID: 1, Name: "Ideal"},
				Genres: []entity.Genre{{ID: 1, Name: "Lo-Fi Hip Hop"}},
				PopularityScore: 89,
			}}, nil
		},
	})

	recs, err := recUseCase.GetTrackRecommendations(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetTrackRecommendations() error = %v", err)
	}

	if len(recs) != 1 {
		t.Fatalf("len(recs) = %d, want %d", len(recs), 1)
	}

	if recs[0].Score <= 0 {
		t.Fatalf("Score = %f, want > 0", recs[0].Score)
	}

	if recs[0].Explanation == "" {
		t.Fatal("Explanation should not be empty")
	}
}

func TestGetTrackRecommendationsColdStart(t *testing.T) {
	t.Parallel()

	recUseCase := NewRecommendationUseCase(&fakeRecommendationRepository{
		listCandidateTracksFn: func(_ context.Context) ([]entity.Track, error) {
			return []entity.Track{{
				ID:    1,
				Title: "Night Drive",
				Artist: entity.Artist{ID: 1, Name: "Ideal"},
				Genres: []entity.Genre{{ID: 1}},
				PopularityScore: 89,
			}}, nil
		},
	})

	recs, err := recUseCase.GetTrackRecommendations(context.Background(), 2, 10)
	if err != nil {
		t.Fatalf("GetTrackRecommendations() error = %v", err)
	}

	if len(recs) != 1 {
		t.Fatalf("len(recs) = %d, want %d", len(recs), 1)
	}

	if recs[0].Score <= 0 {
		t.Fatalf("Score = %f, want > 0", recs[0].Score)
	}
}

func TestScoringGenreMatch(t *testing.T) {
	t.Parallel()

	profile := &entity.UserProfile{
		FavoriteGenreIDs: []int64{1, 2},
	}

	track := entity.Track{
		ID: 1,
		Genres: []entity.Genre{{ID: 1, Name: "Lo-Fi"}, {ID: 2, Name: "Chill"}},
	}

	favGenres := makeSet(profile.FavoriteGenreIDs)
	score := calcGenreScore(track, favGenres)

	if score != 40 {
		t.Fatalf("genre score = %f, want %f", score, 40.0)
	}
}

func TestScoringArtistMatch(t *testing.T) {
	t.Parallel()

	favArtists := makeSet([]int64{1})

	track := entity.Track{
		Artist: entity.Artist{ID: 1, Name: "Ideal"},
	}

	score := calcArtistScore(track, favArtists)
	if score != 25 {
		t.Fatalf("artist score = %f, want %f", score, 25.0)
	}
}

func TestScoringPopularity(t *testing.T) {
	t.Parallel()

	track := entity.Track{PopularityScore: 80}
	score := calcPopularityScore(track)

	if score != 16 {
		t.Fatalf("popularity score = %f, want %f", score, 16.0)
	}
}

func TestColdStartScores(t *testing.T) {
	t.Parallel()

	tracks := []entity.Track{
		{ID: 1, Title: "A", PopularityScore: 90},
		{ID: 2, Title: "B", PopularityScore: 50},
	}

	scored := coldStartScores(tracks)

	if len(scored) != 2 {
		t.Fatalf("len = %d, want %d", len(scored), 2)
	}

	if scored[0].Score <= scored[1].Score {
		t.Fatalf("track 1 (%f) should score higher than track 2 (%f)", scored[0].Score, scored[1].Score)
	}
}

func TestFilterSkipped(t *testing.T) {
	t.Parallel()

	tracks := []scoredTrack{
		{Track: entity.Track{ID: 1}},
		{Track: entity.Track{ID: 2}},
		{Track: entity.Track{ID: 3}},
	}

	skipped := map[int64]bool{2: true}

	filtered := filterSkipped(tracks, skipped)

	if len(filtered) != 2 {
		t.Fatalf("len = %d, want %d", len(filtered), 2)
	}

	if filtered[0].Track.ID == 2 || filtered[1].Track.ID == 2 {
		t.Fatal("skipped track was not filtered")
	}
}

func TestDownrankPlayed(t *testing.T) {
	t.Parallel()

	tracks := []scoredTrack{
		{Track: entity.Track{ID: 1}, Score: 100},
		{Track: entity.Track{ID: 2}, Score: 80},
	}

	played := map[int64]bool{1: true}

	downrankPlayed(tracks, played)

	if tracks[0].Score != 70 {
		t.Fatalf("played track score = %f, want %f", tracks[0].Score, 70.0)
	}

	if tracks[1].Score != 80 {
		t.Fatalf("non-played track score = %f, want %f", tracks[1].Score, 80.0)
	}
}
