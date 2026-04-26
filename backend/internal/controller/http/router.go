package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(userHandler *UserHandler, authMiddleware *AuthMiddleware) (*gin.Engine, error) {
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
	}

	return router, nil
}
