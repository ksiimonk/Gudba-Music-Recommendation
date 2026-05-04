package http

import (
	"context"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/auth"
)

type fakeEventUseCase struct {
	recordEventFn         func(ctx context.Context, userID int64, entityType, eventType string, entityID int64) error
	recordTrackPlayFn     func(ctx context.Context, userID int64, trackID int64) error
	recordTrackLikeFn     func(ctx context.Context, userID int64, trackID int64) error
	recordTrackDislikeFn  func(ctx context.Context, userID int64, trackID int64) error
	recordTrackSkipFn     func(ctx context.Context, userID int64, trackID int64) error
	recordPlaylistOpenFn  func(ctx context.Context, userID int64, playlistID int64) error
}

func (f *fakeEventUseCase) RecordEvent(ctx context.Context, userID int64, entityType, eventType string, entityID int64) error {
	if f.recordEventFn != nil {
		return f.recordEventFn(ctx, userID, entityType, eventType, entityID)
	}
	return nil
}

func (f *fakeEventUseCase) RecordTrackPlay(ctx context.Context, userID int64, trackID int64) error {
	if f.recordTrackPlayFn != nil {
		return f.recordTrackPlayFn(ctx, userID, trackID)
	}
	return nil
}

func (f *fakeEventUseCase) RecordTrackLike(ctx context.Context, userID int64, trackID int64) error {
	if f.recordTrackLikeFn != nil {
		return f.recordTrackLikeFn(ctx, userID, trackID)
	}
	return nil
}

func (f *fakeEventUseCase) RecordTrackDislike(ctx context.Context, userID int64, trackID int64) error {
	if f.recordTrackDislikeFn != nil {
		return f.recordTrackDislikeFn(ctx, userID, trackID)
	}
	return nil
}

func (f *fakeEventUseCase) RecordTrackSkip(ctx context.Context, userID int64, trackID int64) error {
	if f.recordTrackSkipFn != nil {
		return f.recordTrackSkipFn(ctx, userID, trackID)
	}
	return nil
}

func (f *fakeEventUseCase) RecordPlaylistOpen(ctx context.Context, userID int64, playlistID int64) error {
	if f.recordPlaylistOpenFn != nil {
		return f.recordPlaylistOpenFn(ctx, userID, playlistID)
	}
	return nil
}

func newEventTestRouter(t *testing.T, eventUseCase *fakeEventUseCase) *gin.Engine {
	t.Helper()

	tokenManager := &fakeTokenManager{
		parseTokenFn: func(token string) (*auth.Claims, error) {
			return &auth.Claims{UserID: 1, Email: "test@example.com"}, nil
		},
	}

	userHandler := NewUserHandler(&fakeUserRegistrationUseCase{}, tokenManager)
	authMiddleware := NewAuthMiddleware(tokenManager)
	eventHandler := NewEventHandler(eventUseCase)

	router, err := NewRouter(userHandler, authMiddleware, eventHandler)
	if err != nil {
		t.Fatalf("NewRouter() error = %v", err)
	}

	return router
}

func TestTrackPlayRoute(t *testing.T) {
	t.Parallel()

	var savedUserID, savedTrackID int64

	router := newEventTestRouter(t, &fakeEventUseCase{
		recordTrackPlayFn: func(_ context.Context, userID, trackID int64) error {
			savedUserID = userID
			savedTrackID = trackID
			return nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/tracks/1/play", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}

	if savedUserID != 1 || savedTrackID != 1 {
		t.Fatalf("saved user_id=%d track_id=%d, want 1, 1", savedUserID, savedTrackID)
	}
}

func TestTrackLikeRoute(t *testing.T) {
	t.Parallel()

	router := newEventTestRouter(t, &fakeEventUseCase{
		recordTrackLikeFn: func(_ context.Context, userID, trackID int64) error {
			return nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/tracks/1/like", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}
}

func TestTrackDislikeRoute(t *testing.T) {
	t.Parallel()

	router := newEventTestRouter(t, &fakeEventUseCase{
		recordTrackDislikeFn: func(_ context.Context, userID, trackID int64) error {
			return nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/tracks/1/dislike", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}
}

func TestTrackSkipRoute(t *testing.T) {
	t.Parallel()

	router := newEventTestRouter(t, &fakeEventUseCase{
		recordTrackSkipFn: func(_ context.Context, userID, trackID int64) error {
			return nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/tracks/1/skip", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}
}

func TestPlaylistOpenRoute(t *testing.T) {
	t.Parallel()

	router := newEventTestRouter(t, &fakeEventUseCase{
		recordPlaylistOpenFn: func(_ context.Context, userID, playlistID int64) error {
			return nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/playlists/1/open", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}
}

func TestGenericEventRoute(t *testing.T) {
	t.Parallel()

	router := newEventTestRouter(t, &fakeEventUseCase{
		recordEventFn: func(_ context.Context, userID int64, entityType, eventType string, entityID int64) error {
			return nil
		},
	})

	body := `{"entity_type":"track","entity_id":1,"event_type":"play"}`
	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/events", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer test-token")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}
}

func TestEventRouteUnauthorized(t *testing.T) {
	t.Parallel()

	router := newEventTestRouter(t, &fakeEventUseCase{})

	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/tracks/1/like", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusUnauthorized)
	}
}

func TestEventRouteInvalidTrackID(t *testing.T) {
	t.Parallel()

	router := newEventTestRouter(t, &fakeEventUseCase{})

	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/tracks/abc/like", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusBadRequest)
	}
}
