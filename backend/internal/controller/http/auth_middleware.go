package http

import (
	"errors"
	nethttp "net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/auth"
)

const authClaimsContextKey = "authClaims"

type tokenParser interface {
	ParseToken(token string) (*auth.Claims, error)
}

type AuthMiddleware struct {
	tokenParser tokenParser
}

func NewAuthMiddleware(tokenParser tokenParser) *AuthMiddleware {
	return &AuthMiddleware{tokenParser: tokenParser}
}

func (m *AuthMiddleware) RequireAuth(c *gin.Context) {
	if m == nil || m.tokenParser == nil {
		c.AbortWithStatusJSON(nethttp.StatusInternalServerError, gin.H{
			"message": "auth middleware is not configured",
		})
		return
	}

	token := extractBearerToken(c.GetHeader("Authorization"))
	if token == "" {
		c.AbortWithStatusJSON(nethttp.StatusUnauthorized, gin.H{
			"message": "authorization token is required",
		})
		return
	}

	claims, err := m.tokenParser.ParseToken(token)
	if err != nil {
		message := "invalid token"
		if errors.Is(err, auth.ErrTokenExpired) {
			message = "token expired"
		}

		c.AbortWithStatusJSON(nethttp.StatusUnauthorized, gin.H{
			"message": message,
		})
		return
	}

	c.Set(authClaimsContextKey, claims)
	c.Next()
}

func getAuthClaims(c *gin.Context) (*auth.Claims, bool) {
	value, ok := c.Get(authClaimsContextKey)
	if !ok {
		return nil, false
	}

	claims, ok := value.(*auth.Claims)
	return claims, ok
}

func extractBearerToken(header string) string {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(header, bearerPrefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
}
