package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"music-recommender-backend/internal/entity"
)

const defaultTokenTTL = 24 * time.Hour

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenManager(secret string) (*TokenManager, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}

	return &TokenManager{
		secret: []byte(secret),
		ttl:    defaultTokenTTL,
	}, nil
}

func (m *TokenManager) GenerateToken(user *entity.User) (string, error) {
	if m == nil || len(m.secret) == 0 {
		return "", errors.New("jwt token manager secret is empty")
	}

	if user == nil {
		return "", errors.New("user is nil")
	}

	header := struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}{
		Alg: "HS256",
		Typ: "JWT",
	}

	now := time.Now()
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Iat:    now.Unix(),
		Exp:    now.Add(m.ttl).Unix(),
	}

	headerPart, err := encodeJWTPart(header)
	if err != nil {
		return "", fmt.Errorf("encode jwt header: %w", err)
	}

	claimsPart, err := encodeJWTPart(claims)
	if err != nil {
		return "", fmt.Errorf("encode jwt claims: %w", err)
	}

	signingInput := headerPart + "." + claimsPart
	signature := m.sign(signingInput)

	return signingInput + "." + signature, nil
}

func (m *TokenManager) ParseToken(token string) (*Claims, error) {
	if m == nil || len(m.secret) == 0 {
		return nil, errors.New("jwt token manager secret is empty")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}

	if err := decodeJWTPart(parts[0], &header); err != nil {
		return nil, ErrInvalidToken
	}

	if header.Alg != "HS256" || header.Typ != "JWT" {
		return nil, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	if !m.verifySignature(signingInput, parts[2]) {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := decodeJWTPart(parts[1], &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if time.Now().Unix() >= claims.Exp {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}

func (m *TokenManager) sign(signingInput string) string {
	hash := hmac.New(sha256.New, m.secret)
	hash.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func (m *TokenManager) verifySignature(signingInput, signature string) bool {
	expectedSignature, err := base64.RawURLEncoding.DecodeString(m.sign(signingInput))
	if err != nil {
		return false
	}

	providedSignature, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	return hmac.Equal(expectedSignature, providedSignature)
}

func encodeJWTPart(value any) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func decodeJWTPart(value string, target any) error {
	payload, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(payload, target)
}
