package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	api := router.Group("/api/v1")
	{
		api.GET("/tracks", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"tracks": []gin.H{
					{
						"id":    1,
						"title": "Night Walk",
						"genre": "lo-fi",
					},
					{
						"id":    2,
						"title": "Blue Jazz",
						"genre": "jazz",
					},
				},
			})
		})
	}

	return router
}
