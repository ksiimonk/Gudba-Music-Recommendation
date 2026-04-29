package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(userHandler *UserHandler, authMiddleware *AuthMiddleware, routeHandlers ...any) (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

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
