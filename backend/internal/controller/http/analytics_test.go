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

type fakeAnalyticsUseCase struct {
	getStatsFn                 func(ctx context.Context) (*entity.Stats, error)
	getRecommendationMetricsFn func(ctx context.Context) (*entity.RecommendationMetrics, error)
}

func (f *fakeAnalyticsUseCase) GetStats(ctx context.Context) (*entity.Stats, error) {
	if f.getStatsFn != nil {
		return f.getStatsFn(ctx)
	}
	return &entity.Stats{}, nil
}

func (f *fakeAnalyticsUseCase) GetRecommendationMetrics(ctx context.Context) (*entity.RecommendationMetrics, error) {
	if f.getRecommendationMetricsFn != nil {
		return f.getRecommendationMetricsFn(ctx)
	}
	return &entity.RecommendationMetrics{}, nil
}

func newAnalyticsTestRouter(t *testing.T, analyticsUseCase *fakeAnalyticsUseCase) *gin.Engine {
	t.Helper()

	tokenManager := &fakeTokenManager{
		parseTokenFn: func(token string) (*auth.Claims, error) {
			return &auth.Claims{UserID: 1, Email: "admin@example.com"}, nil
		},
	}

	userHandler := NewUserHandler(&fakeUserRegistrationUseCase{}, tokenManager)
	authMiddleware := NewAuthMiddleware(tokenManager)
	analyticsHandler := NewAnalyticsHandler(analyticsUseCase)

	router, err := NewRouter(userHandler, authMiddleware, analyticsHandler)
	if err != nil {
		t.Fatalf("NewRouter() error = %v", err)
	}

	return router
}

func TestAdminStatsRoute(t *testing.T) {
	t.Parallel()

	router := newAnalyticsTestRouter(t, &fakeAnalyticsUseCase{
		getStatsFn: func(_ context.Context) (*entity.Stats, error) {
			return &entity.Stats{UsersCount: 5, TracksCount: 41}, nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/admin/stats", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}
}

func TestAdminMetricsRoute(t *testing.T) {
	t.Parallel()

	router := newAnalyticsTestRouter(t, &fakeAnalyticsUseCase{})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/admin/recommendation-metrics", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}
}

func TestAdminRoutesUnauthorized(t *testing.T) {
	t.Parallel()

	router := newAnalyticsTestRouter(t, &fakeAnalyticsUseCase{})

	for _, path := range []string{"/api/v1/admin/stats", "/api/v1/admin/recommendation-metrics"} {
		request := httptest.NewRequest(nethttp.MethodGet, path, nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		if response.Code != nethttp.StatusUnauthorized {
			t.Fatalf("%s: status = %d, want %d", path, response.Code, nethttp.StatusUnauthorized)
		}
	}
}
