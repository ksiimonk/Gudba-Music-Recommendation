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

type fakePlaylistUseCase struct {
	listPlaylistsFn   func(ctx context.Context) ([]entity.Playlist, error)
	getPlaylistByIDFn func(ctx context.Context, id int64) (*entity.Playlist, error)
}

func (f *fakePlaylistUseCase) ListPlaylists(ctx context.Context) ([]entity.Playlist, error) {
	if f.listPlaylistsFn != nil {
		return f.listPlaylistsFn(ctx)
	}

	return []entity.Playlist{}, nil
}

func (f *fakePlaylistUseCase) GetPlaylistByID(ctx context.Context, id int64) (*entity.Playlist, error) {
	if f.getPlaylistByIDFn != nil {
		return f.getPlaylistByIDFn(ctx, id)
	}

	return nil, errors.New("get playlist by id function is not set")
}

func newPlaylistTestRouter(t *testing.T, playlistUseCase *fakePlaylistUseCase) *gin.Engine {
	t.Helper()

	userHandler := NewUserHandler(&fakeUserRegistrationUseCase{}, &fakeTokenManager{})
	authMiddleware := NewAuthMiddleware(&fakeTokenManager{})
	playlistHandler := NewPlaylistHandler(playlistUseCase)

	router, err := NewRouter(userHandler, authMiddleware, playlistHandler)
	if err != nil {
		t.Fatalf("NewRouter() error = %v", err)
	}

	return router
}

func TestListPlaylistsRouteSuccess(t *testing.T) {
	t.Parallel()

	router := newPlaylistTestRouter(t, &fakePlaylistUseCase{
		listPlaylistsFn: func(_ context.Context) ([]entity.Playlist, error) {
			return []entity.Playlist{{
				ID:          1,
				Name:        "Lo-Fi Hip Hop Mix",
				Description: "Chill beats to study, relax, and focus.",
				IsPublic:    true,
				TrackCount:  8,
			}}, nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/playlists", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var responseBody struct {
		Playlists []entity.Playlist `json:"playlists"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(responseBody.Playlists) != 1 {
		t.Fatalf("len(playlists) = %d, want %d", len(responseBody.Playlists), 1)
	}
}

func TestGetPlaylistRouteSuccess(t *testing.T) {
	t.Parallel()

	router := newPlaylistTestRouter(t, &fakePlaylistUseCase{
		getPlaylistByIDFn: func(_ context.Context, id int64) (*entity.Playlist, error) {
			return &entity.Playlist{
				ID:   id,
				Name: "Lo-Fi Hip Hop Mix",
				Tracks: []entity.Track{
					{ID: 1, Title: "Night Drive"},
				},
			}, nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/playlists/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var responseBody struct {
		Playlist entity.Playlist `json:"playlist"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if responseBody.Playlist.ID != 1 {
		t.Fatalf("ID = %d, want %d", responseBody.Playlist.ID, 1)
	}

	if len(responseBody.Playlist.Tracks) != 1 {
		t.Fatalf("len(tracks) = %d, want %d", len(responseBody.Playlist.Tracks), 1)
	}
}

func TestGetPlaylistRouteInvalidID(t *testing.T) {
	t.Parallel()

	router := newPlaylistTestRouter(t, &fakePlaylistUseCase{})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/playlists/abc", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusBadRequest)
	}
}

func TestGetPlaylistRouteNotFound(t *testing.T) {
	t.Parallel()

	router := newPlaylistTestRouter(t, &fakePlaylistUseCase{
		getPlaylistByIDFn: func(_ context.Context, id int64) (*entity.Playlist, error) {
			return nil, repository.ErrPlaylistNotFound
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/playlists/99", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusNotFound)
	}
}
