package usecase

import (
	"context"
	"errors"
	"testing"

	"music-recommender-backend/internal/entity"
	"music-recommender-backend/internal/repository"
)

type fakeOnboardingRepository struct {
	saveOnboardingFn func(ctx context.Context, userID int64, request *entity.OnboardingRequest) error
	getUserProfileFn func(ctx context.Context, userID int64) (*entity.UserProfile, error)
}

func (f *fakeOnboardingRepository) SaveOnboarding(ctx context.Context, userID int64, request *entity.OnboardingRequest) error {
	if f.saveOnboardingFn != nil {
		return f.saveOnboardingFn(ctx, userID, request)
	}

	return nil
}

func (f *fakeOnboardingRepository) GetUserProfile(ctx context.Context, userID int64) (*entity.UserProfile, error) {
	if f.getUserProfileFn != nil {
		return f.getUserProfileFn(ctx, userID)
	}

	return nil, repository.ErrProfileNotFound
}

func TestSaveOnboardingSuccess(t *testing.T) {
	t.Parallel()

	var savedUserID int64
	var savedRequest *entity.OnboardingRequest

	onboardingUseCase := NewOnboardingUseCase(&fakeOnboardingRepository{
		saveOnboardingFn: func(_ context.Context, userID int64, request *entity.OnboardingRequest) error {
			savedUserID = userID
			savedRequest = request
			return nil
		},
	})

	request := &entity.OnboardingRequest{
		GenreIDs:  []int64{1, 2, 3},
		ArtistIDs: []int64{5},
		TrackIDs:  []int64{},
		Contexts:  []string{"focus", "relax"},
	}

	err := onboardingUseCase.SaveOnboarding(context.Background(), 42, request)
	if err != nil {
		t.Fatalf("SaveOnboarding() error = %v", err)
	}

	if savedUserID != 42 {
		t.Fatalf("userID = %d, want %d", savedUserID, 42)
	}

	if savedRequest == nil {
		t.Fatal("request was not passed to repository")
	}
}

func TestSaveOnboardingEmpty(t *testing.T) {
	t.Parallel()

	onboardingUseCase := NewOnboardingUseCase(&fakeOnboardingRepository{})

	request := &entity.OnboardingRequest{
		GenreIDs:  []int64{},
		ArtistIDs: []int64{},
		TrackIDs:  []int64{},
		Contexts:  []string{},
	}

	err := onboardingUseCase.SaveOnboarding(context.Background(), 42, request)
	if !errors.Is(err, ErrEmptyOnboarding) {
		t.Fatalf("error = %v, want %v", err, ErrEmptyOnboarding)
	}
}

func TestSaveOnboardingInvalidIDs(t *testing.T) {
	t.Parallel()

	onboardingUseCase := NewOnboardingUseCase(&fakeOnboardingRepository{})

	request := &entity.OnboardingRequest{
		GenreIDs:  []int64{0, -1},
		ArtistIDs: []int64{1},
	}

	err := onboardingUseCase.SaveOnboarding(context.Background(), 42, request)
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidID)
	}
}

func TestGetProfileSuccess(t *testing.T) {
	t.Parallel()

	onboardingUseCase := NewOnboardingUseCase(&fakeOnboardingRepository{
		getUserProfileFn: func(_ context.Context, userID int64) (*entity.UserProfile, error) {
			return &entity.UserProfile{
				UserID:            userID,
				FavoriteGenreIDs:  []int64{1, 2},
				FavoriteArtistIDs: []int64{},
				StarterTrackIDs:   []int64{},
				Contexts:          []string{"focus"},
			}, nil
		},
	})

	profile, err := onboardingUseCase.GetProfile(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}

	if profile.UserID != 42 {
		t.Fatalf("UserID = %d, want %d", profile.UserID, 42)
	}

	if len(profile.FavoriteGenreIDs) != 2 {
		t.Fatalf("len(FavoriteGenreIDs) = %d, want %d", len(profile.FavoriteGenreIDs), 2)
	}
}

func TestGetProfileNotFound(t *testing.T) {
	t.Parallel()

	onboardingUseCase := NewOnboardingUseCase(&fakeOnboardingRepository{})

	_, err := onboardingUseCase.GetProfile(context.Background(), 99)
	if !errors.Is(err, repository.ErrProfileNotFound) {
		t.Fatalf("error = %v, want %v", err, repository.ErrProfileNotFound)
	}
}
