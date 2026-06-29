package server

import (
	"github.com/gin-gonic/gin"
	"github.com/seniorLikeToCode/url-shortener/internal/handler"
)

func NewRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	// Serving the static file
	r.LoadHTMLGlob("web/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/health", h.Health)
	r.POST("/create-short-url", h.CreateShortURL)
	r.GET("/:shortUrl", h.HandleShortURLRedirect)
	return r
}
