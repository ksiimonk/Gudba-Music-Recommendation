package http

import (
	nethttp "net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	trustedLocalProxy = "127.0.0.1"

	healthRoute = "/healthz"
	apiPrefix   = "/api/v1"
	authPrefix  = "/auth"

	registerRoute       = "/register"
	loginRoute          = "/login"
	meRoute             = "/me"
	tracksRoute         = "/tracks"
	trackDetailRoute    = "/tracks/:id"
	playlistsRoute      = "/playlists"
	playlistDetailRoute = "/playlists/:id"
	genresRoute         = "/genres"
	artistsRoute        = "/artists"

	corsAllowHeaders = "Authorization, Content-Type"
	corsAllowMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	originHeader     = "Origin"
	varyHeader       = "Origin"
)

var allowedFrontendOriginPrefixes = []string{
	"http://127.0.0.1:",
	"http://localhost:",
}

func NewRouter(userHandler *UserHandler, authMiddleware *AuthMiddleware, routeHandlers ...any) (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), corsMiddleware())

	if err := router.SetTrustedProxies([]string{trustedLocalProxy}); err != nil {
		return nil, err
	}

	registerSwaggerRoutes(router)

	router.GET(healthRoute, func(c *gin.Context) {
		c.JSON(nethttp.StatusOK, gin.H{
			"message": "ok",
		})
	})

	api := router.Group(apiPrefix)
	{
		auth := api.Group(authPrefix)
		{
			auth.POST(registerRoute, userHandler.Register)
			auth.POST(loginRoute, userHandler.Login)
		}

		api.GET(meRoute, authMiddleware.RequireAuth, userHandler.Me)

		var trackHandler *TrackHandler
		var playlistHandler *PlaylistHandler
		var genreHandler *GenreHandler
		var artistHandler *ArtistHandler

		for _, handler := range routeHandlers {
			switch typedHandler := handler.(type) {
			case *TrackHandler:
				trackHandler = typedHandler
			case *PlaylistHandler:
				playlistHandler = typedHandler
			case *GenreHandler:
				genreHandler = typedHandler
			case *ArtistHandler:
				artistHandler = typedHandler
			}
		}

		if trackHandler != nil {
			api.GET(tracksRoute, trackHandler.ListTracks)
			api.GET(trackDetailRoute, trackHandler.GetTrack)
		}

		if playlistHandler != nil {
			api.GET(playlistsRoute, playlistHandler.ListPlaylists)
			api.GET(playlistDetailRoute, playlistHandler.GetPlaylist)
		}

		if genreHandler != nil {
			api.GET(genresRoute, genreHandler.ListGenres)
		}

		if artistHandler != nil {
			api.GET(artistsRoute, artistHandler.ListArtists)
		}
	}

	return router, nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader(originHeader)
		if isAllowedFrontendOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Headers", corsAllowHeaders)
			c.Header("Access-Control-Allow-Methods", corsAllowMethods)
			c.Header(varyHeader, originHeader)
		}

		if c.Request.Method == nethttp.MethodOptions {
			c.AbortWithStatus(nethttp.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedFrontendOrigin(origin string) bool {
	for _, prefix := range allowedFrontendOriginPrefixes {
		if strings.HasPrefix(origin, prefix) {
			return true
		}
	}

	return false
}
