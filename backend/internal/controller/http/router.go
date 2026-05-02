package http

import (
	nethttp "net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewRouter(userHandler *UserHandler, authMiddleware *AuthMiddleware, routeHandlers ...any) (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), corsMiddleware())

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil, err
	}

	registerSwaggerRoutes(router)

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(nethttp.StatusOK, gin.H{
			"message": "ok",
		})
	})

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		api.GET("/me", authMiddleware.RequireAuth, userHandler.Me)

		var trackHandler *TrackHandler
		var playlistHandler *PlaylistHandler

		for _, handler := range routeHandlers {
			switch typedHandler := handler.(type) {
			case *TrackHandler:
				trackHandler = typedHandler
			case *PlaylistHandler:
				playlistHandler = typedHandler
			}
		}

		if trackHandler != nil {
			api.GET("/tracks", trackHandler.ListTracks)
			api.GET("/tracks/:id", trackHandler.GetTrack)
		}

		if playlistHandler != nil {
			api.GET("/playlists", playlistHandler.ListPlaylists)
			api.GET("/playlists/:id", playlistHandler.GetPlaylist)
		}
	}

	return router, nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if isAllowedFrontendOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Vary", "Origin")
		}

		if c.Request.Method == nethttp.MethodOptions {
			c.AbortWithStatus(nethttp.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedFrontendOrigin(origin string) bool {
	return strings.HasPrefix(origin, "http://127.0.0.1:") ||
		strings.HasPrefix(origin, "http://localhost:")
}
