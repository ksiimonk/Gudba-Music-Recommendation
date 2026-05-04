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

type fakeArtistUseCase struct {
	listArtistsFn func(ctx context.Context) ([]entity.Artist, error)
}

func (f *fakeArtistUseCase) ListArtists(ctx context.Context) ([]entity.Artist, error) {
	if f.listArtistsFn != nil {
		return f.listArtistsFn(ctx)
	}

	return []entity.Artist{}, nil
}

func newArtistTestRouter(t *testing.T, artistUseCase *fakeArtistUseCase) *gin.Engine {
	t.Helper()

	userHandler := NewUserHandler(&fakeUserRegistrationUseCase{}, &fakeTokenManager{})
	authMiddleware := NewAuthMiddleware(&fakeTokenManager{})
	artistHandler := NewArtistHandler(artistUseCase)

	router, err := NewRouter(userHandler, authMiddleware, artistHandler)
	if err != nil {
		t.Fatalf("NewRouter() error = %v", err)
	}

	return router
}

func TestListArtistsRouteSuccess(t *testing.T) {
	t.Parallel()

	router := newArtistTestRouter(t, &fakeArtistUseCase{
		listArtistsFn: func(_ context.Context) ([]entity.Artist, error) {
			return []entity.Artist{{ID: 1, Name: "Ideal"}}, nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/artists", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var responseBody struct {
		Artists []entity.Artist `json:"artists"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(responseBody.Artists) != 1 {
		t.Fatalf("len(artists) = %d, want %d", len(responseBody.Artists), 1)
	}
}
