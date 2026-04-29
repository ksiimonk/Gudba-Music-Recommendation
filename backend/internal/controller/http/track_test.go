package http

import (
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/entity"
	"music-recommender-backend/internal/repository"
)

type fakeTrackUseCase struct {
	listTracksFn   func(ctx context.Context) ([]entity.Track, error)
	getTrackByIDFn func(ctx context.Context, id int64) (*entity.Track, error)
}

func (f *fakeTrackUseCase) ListTracks(ctx context.Context) ([]entity.Track, error) {
	if f.listTracksFn != nil {
		return f.listTracksFn(ctx)
	}

	return []entity.Track{}, nil
}

func (f *fakeTrackUseCase) GetTrackByID(ctx context.Context, id int64) (*entity.Track, error) {
	if f.getTrackByIDFn != nil {
		return f.getTrackByIDFn(ctx, id)
	}

	return nil, errors.New("get track by id function is not set")
}

func newTrackTestRouter(t *testing.T, trackUseCase *fakeTrackUseCase) *gin.Engine {
	t.Helper()

	userHandler := NewUserHandler(&fakeUserRegistrationUseCase{}, &fakeTokenManager{})
	authMiddleware := NewAuthMiddleware(&fakeTokenManager{})
	trackHandler := NewTrackHandler(trackUseCase)

	router, err := NewRouter(userHandler, authMiddleware, trackHandler)
	if err != nil {
		t.Fatalf("NewRouter() error = %v", err)
	}

	return router
}

func TestListTracksRouteSuccess(t *testing.T) {
	t.Parallel()

	router := newTrackTestRouter(t, &fakeTrackUseCase{
		listTracksFn: func(_ context.Context) ([]entity.Track, error) {
			return []entity.Track{{
				ID:    1,
				Title: "Night Drive",
				Artist: entity.Artist{
					ID:   1,
					Name: "Ideal",
				},
				Genres: []entity.Genre{{ID: 1, Name: "Lo-Fi Hip Hop"}},
			}}, nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/tracks", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var responseBody struct {
		Tracks []entity.Track `json:"tracks"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(responseBody.Tracks) != 1 {
		t.Fatalf("len(tracks) = %d, want %d", len(responseBody.Tracks), 1)
	}
}

func TestGetTrackRouteSuccess(t *testing.T) {
	t.Parallel()

	router := newTrackTestRouter(t, &fakeTrackUseCase{
		getTrackByIDFn: func(_ context.Context, id int64) (*entity.Track, error) {
			return &entity.Track{ID: id, Title: "Night Drive"}, nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/tracks/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var responseBody struct {
		Track entity.Track `json:"track"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if responseBody.Track.ID != 1 {
		t.Fatalf("ID = %d, want %d", responseBody.Track.ID, 1)
	}
}

func TestGetTrackRouteInvalidID(t *testing.T) {
	t.Parallel()

	router := newTrackTestRouter(t, &fakeTrackUseCase{})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/tracks/abc", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusBadRequest)
	}
}

func TestGetTrackRouteNotFound(t *testing.T) {
	t.Parallel()

	router := newTrackTestRouter(t, &fakeTrackUseCase{
		getTrackByIDFn: func(_ context.Context, id int64) (*entity.Track, error) {
			return nil, repository.ErrTrackNotFound
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/tracks/99", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusNotFound)
	}
}
