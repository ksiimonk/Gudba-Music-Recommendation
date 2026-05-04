package usecase

import (
	"context"
	"errors"

	"music-recommender-backend/internal/entity"
)

var ErrEmptyOnboarding = errors.New("at least one genre, artist, or track must be selected")

type onboardingRepository interface {
	SaveOnboarding(ctx context.Context, userID int64, request *entity.OnboardingRequest) error
	GetUserProfile(ctx context.Context, userID int64) (*entity.UserProfile, error)
}

type OnboardingUseCase struct {
	onboardingRepository onboardingRepository
}

func NewOnboardingUseCase(onboardingRepository onboardingRepository) *OnboardingUseCase {
	return &OnboardingUseCase{onboardingRepository: onboardingRepository}
}

func (u *OnboardingUseCase) SaveOnboarding(ctx context.Context, userID int64, request *entity.OnboardingRequest) error {
	if u == nil || u.onboardingRepository == nil {
		return errors.New("onboarding use case repository is nil")
	}

	if err := validateOnboardingRequest(request); err != nil {
		return err
	}

	return u.onboardingRepository.SaveOnboarding(ctx, userID, request)
}

func (u *OnboardingUseCase) GetProfile(ctx context.Context, userID int64) (*entity.UserProfile, error) {
	if u == nil || u.onboardingRepository == nil {
		return nil, errors.New("onboarding use case repository is nil")
	}

	profile, err := u.onboardingRepository.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func validateOnboardingRequest(request *entity.OnboardingRequest) error {
	if request == nil {
		return errors.New("onboarding request is nil")
	}

	for _, id := range request.GenreIDs {
		if id <= 0 {
			return ErrInvalidID
		}
	}

	for _, id := range request.ArtistIDs {
		if id <= 0 {
			return ErrInvalidID
		}
	}

	for _, id := range request.TrackIDs {
		if id <= 0 {
			return ErrInvalidID
		}
	}

	if len(request.GenreIDs) == 0 && len(request.ArtistIDs) == 0 && len(request.TrackIDs) == 0 {
		return ErrEmptyOnboarding
	}

	return nil
}
