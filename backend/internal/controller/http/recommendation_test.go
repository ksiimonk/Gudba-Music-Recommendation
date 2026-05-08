package http

import (
	"context"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/auth"
	"music-recommender-backend/internal/entity"
)

type fakeRecommendationUseCase struct {
	getTrackRecommendationsFn     func(ctx context.Context, userID int64, limit int) ([]entity.TrackRecommendation, error)
	getPlaylistRecommendationsFn  func(ctx context.Context, userID int64, limit int) ([]entity.PlaylistRecommendation, error)
}

func (f *fakeRecommendationUseCase) GetTrackRecommendations(ctx context.Context, userID int64, limit int) ([]entity.TrackRecommendation, error) {
	if f.getTrackRecommendationsFn != nil {
		return f.getTrackRecommendationsFn(ctx, userID, limit)
	}

	return []entity.TrackRecommendation{}, nil
}

func (f *fakeRecommendationUseCase) GetPlaylistRecommendations(ctx context.Context, userID int64, limit int) ([]entity.PlaylistRecommendation, error) {
	if f.getPlaylistRecommendationsFn != nil {
		return f.getPlaylistRecommendationsFn(ctx, userID, limit)
	}

	return []entity.PlaylistRecommendation{}, nil
}

func newRecommendationTestRouter(t *testing.T, recUseCase *fakeRecommendationUseCase) *gin.Engine {
	t.Helper()

	tokenManager := &fakeTokenManager{
		parseTokenFn: func(token string) (*auth.Claims, error) {
			return &auth.Claims{UserID: 1, Email: "test@example.com"}, nil
		},
	}

	userHandler := NewUserHandler(&fakeUserRegistrationUseCase{}, tokenManager)
	authMiddleware := NewAuthMiddleware(tokenManager)
	recHandler := NewRecommendationHandler(recUseCase)

	router, err := NewRouter(userHandler, authMiddleware, recHandler)
	if err != nil {
		t.Fatalf("NewRouter() error = %v", err)
	}

	return router
}

func TestRecommendationRouteSuccess(t *testing.T) {
	t.Parallel()

	router := newRecommendationTestRouter(t, &fakeRecommendationUseCase{
		getTrackRecommendationsFn: func(_ context.Context, userID int64, limit int) ([]entity.TrackRecommendation, error) {
			return []entity.TrackRecommendation{
				{
					Track:       entity.Track{ID: 1, Title: "Night Drive", Artist: entity.Artist{ID: 1, Name: "Ideal"}, Genres: []entity.Genre{{ID: 1, Name: "Lo-Fi"}}},
					Score:       85.5,
					Explanation: "Совпадает с твоими любимыми жанрами",
				},
			}, nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/recommendations/tracks?limit=5", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}
}

func TestRecommendationRouteUnauthorized(t *testing.T) {
	t.Parallel()

	router := newRecommendationTestRouter(t, &fakeRecommendationUseCase{})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/recommendations/tracks", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusUnauthorized)
	}
}

func TestRecommendationRouteEmptyResult(t *testing.T) {
	t.Parallel()

	router := newRecommendationTestRouter(t, &fakeRecommendationUseCase{})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/recommendations/tracks", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}
}
