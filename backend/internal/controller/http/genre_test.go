package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/entity"
)

type fakeGenreUseCase struct {
	listGenresFn func(ctx context.Context) ([]entity.Genre, error)
}

func (f *fakeGenreUseCase) ListGenres(ctx context.Context) ([]entity.Genre, error) {
	if f.listGenresFn != nil {
		return f.listGenresFn(ctx)
	}

	return []entity.Genre{}, nil
}

func newGenreTestRouter(t *testing.T, genreUseCase *fakeGenreUseCase) *gin.Engine {
	t.Helper()

	userHandler := NewUserHandler(&fakeUserRegistrationUseCase{}, &fakeTokenManager{})
	authMiddleware := NewAuthMiddleware(&fakeTokenManager{})
	genreHandler := NewGenreHandler(genreUseCase)

	router, err := NewRouter(userHandler, authMiddleware, genreHandler)
	if err != nil {
		t.Fatalf("NewRouter() error = %v", err)
	}

	return router
}

func TestListGenresRouteSuccess(t *testing.T) {
	t.Parallel()

	router := newGenreTestRouter(t, &fakeGenreUseCase{
		listGenresFn: func(_ context.Context) ([]entity.Genre, error) {
			return []entity.Genre{{ID: 1, Name: "Lo-Fi Hip Hop"}}, nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/genres", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var responseBody struct {
		Genres []entity.Genre `json:"genres"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(responseBody.Genres) != 1 {
		t.Fatalf("len(genres) = %d, want %d", len(responseBody.Genres), 1)
	}
}
