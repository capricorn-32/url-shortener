package server

import (
	"github.com/gin-gonic/gin"
	"github.com/seniorLikeToCode/url-shortener/internal/handler"
)

func NewRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/", h.Health)
	r.POST("/create-short-url", h.CreateShortURL)
	r.GET("/:shortUrl", h.HandleShortURLRedirect)
	return r
}
