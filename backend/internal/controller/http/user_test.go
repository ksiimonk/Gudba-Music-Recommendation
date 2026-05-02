package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/auth"
	"music-recommender-backend/internal/entity"
	"music-recommender-backend/internal/repository"
	"music-recommender-backend/internal/usecase"
)

type fakeUserRegistrationUseCase struct {
	registerUserFn   func(ctx context.Context, email, password string) (*entity.User, error)
	loginUserFn      func(ctx context.Context, email, password string) (*entity.User, error)
	getUserByEmailFn func(ctx context.Context, email string) (*entity.User, error)
}

func (f *fakeUserRegistrationUseCase) RegisterUser(ctx context.Context, email, password string) (*entity.User, error) {
	if f.registerUserFn != nil {
		return f.registerUserFn(ctx, email, password)
	}

	return nil, errors.New("register user function is not set")
}

func (f *fakeUserRegistrationUseCase) LoginUser(ctx context.Context, email, password string) (*entity.User, error) {
	if f.loginUserFn != nil {
		return f.loginUserFn(ctx, email, password)
	}

	return nil, errors.New("login user function is not set")
}

func (f *fakeUserRegistrationUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	if f.getUserByEmailFn != nil {
		return f.getUserByEmailFn(ctx, email)
	}

	return nil, errors.New("get user by email function is not set")
}

type fakeTokenManager struct {
	generateTokenFn func(user *entity.User) (string, error)
	parseTokenFn    func(token string) (*auth.Claims, error)
}

func (f *fakeTokenManager) GenerateToken(user *entity.User) (string, error) {
	if f.generateTokenFn != nil {
		return f.generateTokenFn(user)
	}

	return "", errors.New("generate token function is not set")
}

func (f *fakeTokenManager) ParseToken(token string) (*auth.Claims, error) {
	if f.parseTokenFn != nil {
		return f.parseTokenFn(token)
	}

	return nil, errors.New("parse token function is not set")
}

func newTestRouter(t *testing.T, userUseCase *fakeUserRegistrationUseCase, tokenManager *fakeTokenManager) *gin.Engine {
	t.Helper()

	userHandler := NewUserHandler(userUseCase, tokenManager)
	authMiddleware := NewAuthMiddleware(tokenManager)

	router, err := NewRouter(userHandler, authMiddleware)
	if err != nil {
		t.Fatalf("NewRouter() error = %v", err)
	}

	return router
}

func TestRegisterRouteSuccess(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t, &fakeUserRegistrationUseCase{
		registerUserFn: func(_ context.Context, email, password string) (*entity.User, error) {
			return &entity.User{
				ID:        1,
				Email:     "test@example.com",
				CreatedAt: time.Date(2026, 4, 23, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 4, 23, 12, 0, 0, 0, time.UTC),
			}, nil
		},
	}, &fakeTokenManager{})

	body := []byte(`{"email":"test@example.com","password":"plain-password"}`)
	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusCreated {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusCreated)
	}

	var responseBody struct {
		Message string `json:"message"`
		User    struct {
			ID    int64  `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if responseBody.Message != "user registered successfully" {
		t.Fatalf("message = %q, want %q", responseBody.Message, "user registered successfully")
	}

	if responseBody.User.Email != "test@example.com" {
		t.Fatalf("email = %q, want %q", responseBody.User.Email, "test@example.com")
	}
}

func TestRegisterRouteValidationError(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t, &fakeUserRegistrationUseCase{}, &fakeTokenManager{})

	body := []byte(`{"email":" ","password":" "}`)
	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusBadRequest)
	}

	var responseBody struct {
		Message string            `json:"message"`
		Errors  map[string]string `json:"errors"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if responseBody.Message != "validation failed" {
		t.Fatalf("message = %q, want %q", responseBody.Message, "validation failed")
	}

	if responseBody.Errors["email"] != "email is required" {
		t.Fatalf("email error = %q, want %q", responseBody.Errors["email"], "email is required")
	}

	if responseBody.Errors["password"] != "password is required" {
		t.Fatalf("password error = %q, want %q", responseBody.Errors["password"], "password is required")
	}
}

func TestRegisterRouteDuplicateEmail(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t, &fakeUserRegistrationUseCase{
		registerUserFn: func(_ context.Context, email, password string) (*entity.User, error) {
			return nil, repository.ErrUserAlreadyExists
		},
	}, &fakeTokenManager{})

	body := []byte(`{"email":"test@example.com","password":"plain-password"}`)
	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusConflict {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusConflict)
	}
}

func TestLoginRouteSuccess(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t, &fakeUserRegistrationUseCase{
		loginUserFn: func(_ context.Context, email, password string) (*entity.User, error) {
			return &entity.User{ID: 1, Email: email}, nil
		},
	}, &fakeTokenManager{
		generateTokenFn: func(user *entity.User) (string, error) {
			return "jwt-token", nil
		},
	})

	body := []byte(`{"email":"test@example.com","password":"plain-password"}`)
	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var responseBody struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if responseBody.Message != "login successful" {
		t.Fatalf("message = %q, want %q", responseBody.Message, "login successful")
	}

	if responseBody.Token != "jwt-token" {
		t.Fatalf("token = %q, want %q", responseBody.Token, "jwt-token")
	}
}

func TestLoginRouteValidationError(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t, &fakeUserRegistrationUseCase{}, &fakeTokenManager{})

	body := []byte(`{"email":" ","password":" "}`)
	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusBadRequest)
	}
}

func TestLoginRouteInvalidCredentials(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t, &fakeUserRegistrationUseCase{
		loginUserFn: func(_ context.Context, email, password string) (*entity.User, error) {
			return nil, usecase.ErrInvalidCredentials
		},
	}, &fakeTokenManager{})

	body := []byte(`{"email":"test@example.com","password":"plain-password"}`)
	request := httptest.NewRequest(nethttp.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusUnauthorized)
	}
}

func TestMeRouteSuccess(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t, &fakeUserRegistrationUseCase{
		getUserByEmailFn: func(_ context.Context, email string) (*entity.User, error) {
			return &entity.User{
				ID:        1,
				Email:     email,
				CreatedAt: time.Date(2026, 4, 23, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 4, 23, 12, 0, 0, 0, time.UTC),
			}, nil
		},
	}, &fakeTokenManager{
		parseTokenFn: func(token string) (*auth.Claims, error) {
			return &auth.Claims{UserID: 1, Email: "test@example.com"}, nil
		},
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/me", nil)
	request.Header.Set("Authorization", "Bearer valid-token")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusOK)
	}
}

func TestMeRouteUnauthorized(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t, &fakeUserRegistrationUseCase{}, &fakeTokenManager{})

	request := httptest.NewRequest(nethttp.MethodGet, "/api/v1/me", nil)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusUnauthorized)
	}
}

func TestCorsPreflightAllowsLocalFrontend(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t, &fakeUserRegistrationUseCase{}, &fakeTokenManager{})

	request := httptest.NewRequest(nethttp.MethodOptions, "/api/v1/auth/login", nil)
	request.Header.Set("Origin", "http://127.0.0.1:5173")
	request.Header.Set("Access-Control-Request-Method", nethttp.MethodPost)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, nethttp.StatusNoContent)
	}

	if response.Header().Get("Access-Control-Allow-Origin") != "http://127.0.0.1:5173" {
		t.Fatalf("Access-Control-Allow-Origin = %q", response.Header().Get("Access-Control-Allow-Origin"))
	}
}
