package http

import (
	_ "embed"
	nethttp "net/http"

	"github.com/gin-gonic/gin"
)

//go:embed swagger/openapi.json
var openAPISpec []byte

//go:embed swagger/index.html
var swaggerIndexHTML []byte

func registerSwaggerRoutes(router *gin.Engine) {
	serveSwaggerIndex := func(c *gin.Context) {
		c.Data(nethttp.StatusOK, "text/html; charset=utf-8", swaggerIndexHTML)
	}

	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(nethttp.StatusMovedPermanently, "/swagger/")
	})

	router.GET("/swagger/", serveSwaggerIndex)
	router.GET("/swagger/index.html", serveSwaggerIndex)

	router.GET("/swagger/doc.json", func(c *gin.Context) {
		c.Data(nethttp.StatusOK, "application/json; charset=utf-8", openAPISpec)
	})
}
