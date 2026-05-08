package http

import (
	"bytes"
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/auth"
	"music-recommender-backend/internal/entity"
	"music-recommender-backend/internal/repository"
	"music-recommender-backend/internal/usecase"
)

type fakeOnboardingUseCase struct {
	saveOnboardingFn func(ctx context.Context, userID int64, request *entity.OnboardingRequest) error
	getProfileFn     func(ctx context.Context, userID int64) (*entity.UserProfile, error)
}

func (f *fakeOnboardingUseCase) SaveOnboarding(ctx context.Context, userID int64, request *entity.OnboardingRequest) error {
	if f.saveOnboardingFn != nil {
		return f.saveOnboardingFn(ctx, userID, request)
	}

	return nil
}

func (f *fakeOnboardingUseCase) GetProfile(ctx context.Context, userID int64) (*entity.UserProfile, error) {
	if f.getProfileFn != nil {
		return f.getProfileFn(ctx, userID)
	}

	return nil, repository.ErrProfileNotFound
}

func newOnboardingTestRouter(t *testing.T, onboardingUseCase *fakeOnboardingUseCase) *gin.Engine {
	t.Helper()

	tokenManager := &fakeTokenManager{
		parseTokenFn: func(token string) (*auth.Claims, error) {
			return &auth.Claims{UserID: 1, Email: "test@example.com"}, nil
		},
	}

	userHandler := NewUserHandler(&fakeUserRegistrationUseCase{}, tokenManager)
	authMiddleware := NewAuthMiddleware(tokenManager)
	onboardingHandler := NewOnboardingHandler(onboardingUseCase)

	router, err := NewRouter(userHandler, authMiddleware, onboardingHandler)
	if err != nil {
		t.Fatalf("NewRouter() error = %v", err)
	}

	return router
}

func TestSaveOnboardingRouteSuccess(t *testing.T) {
	t.Parallel()

	var savedUserID int64
	var savedRequest *entity.OnboardingRequest

	router := newOnboardingTestRouter(t, &fakeOnboardingUseCase{
		saveOnboardingFn: func(_ context.Context, userID int64, request *entity.OnboardingRequest) error {
			savedUserID = userID
			savedRequest = request
			return nil
		},
	})

	body := []byte(`{"genre_ids":[1,2],"artist_ids":[3],"track_ids":[],"contexts":["focus"]}`)
	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/me/onboarding", bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer test-token")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}

	if savedUserID != 1 {
		t.Fatalf("userID = %d, want %d", savedUserID, 1)
	}

	if savedRequest == nil || len(savedRequest.GenreIDs) != 2 {
		t.Fatal("request was not parsed correctly")
	}
}

func TestSaveOnboardingRouteEmptyBodyError(t *testing.T) {
	t.Parallel()

	router := newOnboardingTestRouter(t, &fakeOnboardingUseCase{
		saveOnboardingFn: func(_ context.Context, _ int64, _ *entity.OnboardingRequest) error {
			return usecase.ErrEmptyOnboarding
		},
	})

	body := []byte(`{}`)
	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/me/onboarding", bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer test-token")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusBadRequest)
	}
}

func TestSaveOnboardingRouteUnauthorized(t *testing.T) {
	t.Parallel()

	router := newOnboardingTestRouter(t, &fakeOnboardingUseCase{})

	body := []byte(`{"genre_ids":[1]}`)
	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/me/onboarding", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusUnauthorized)
	}
}

func TestGetProfileRouteSuccess(t *testing.T) {
	t.Parallel()

	router := newOnboardingTestRouter(t, &fakeOnboardingUseCase{
		getProfileFn: func(_ context.Context, userID int64) (*entity.UserProfile, error) {
			return &entity.UserProfile{
				UserID:            userID,
				FavoriteGenreIDs:  []int64{1, 2},
				FavoriteArtistIDs: []int64{},
				StarterTrackIDs:   []int64{},
				Contexts:          []string{"focus"},
			}, nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/me/profile", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var responseBody struct {
		Profile entity.UserProfile `json:"profile"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(responseBody.Profile.FavoriteGenreIDs) != 2 {
		t.Fatalf("len(genres) = %d, want %d", len(responseBody.Profile.FavoriteGenreIDs), 2)
	}
}

func TestGetProfileRouteEmptyProfile(t *testing.T) {
	t.Parallel()

	router := newOnboardingTestRouter(t, &fakeOnboardingUseCase{})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/me/profile", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var responseBody struct {
		Profile entity.UserProfile `json:"profile"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if responseBody.Profile.UserID != 1 {
		t.Fatalf("UserID = %d, want %d", responseBody.Profile.UserID, 1)
	}

	if len(responseBody.Profile.FavoriteGenreIDs) != 0 {
		t.Fatalf("expected empty favorite_genre_ids, got %v", responseBody.Profile.FavoriteGenreIDs)
	}
}

func TestGetProfileRouteUnauthorized(t *testing.T) {
	t.Parallel()

	router := newOnboardingTestRouter(t, &fakeOnboardingUseCase{})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/me/profile", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusUnauthorized)
	}
}
