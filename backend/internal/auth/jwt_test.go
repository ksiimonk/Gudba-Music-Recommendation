package auth

import (
	"errors"
	"testing"
	"time"

	"music-recommender-backend/internal/entity"
)

func TestGenerateAndParseToken(t *testing.T) {
	t.Parallel()

	tokenManager, err := NewTokenManager("super-secret")
	if err != nil {
		t.Fatalf("NewTokenManager() error = %v", err)
	}

	token, err := tokenManager.GenerateToken(&entity.User{
		ID:    1,
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := tokenManager.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.UserID != 1 {
		t.Fatalf("UserID = %d, want %d", claims.UserID, 1)
	}

	if claims.Email != "test@example.com" {
		t.Fatalf("Email = %q, want %q", claims.Email, "test@example.com")
	}
}

func TestParseTokenExpired(t *testing.T) {
	t.Parallel()

	tokenManager, err := NewTokenManager("super-secret")
	if err != nil {
		t.Fatalf("NewTokenManager() error = %v", err)
	}

	tokenManager.ttl = -time.Minute

	token, err := tokenManager.GenerateToken(&entity.User{
		ID:    1,
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, err = tokenManager.ParseToken(token)
	if !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("error = %v, want %v", err, ErrTokenExpired)
	}
}
