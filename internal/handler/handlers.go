package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/seniorLikeToCode/url-shortener/internal/shortener"
	"github.com/seniorLikeToCode/url-shortener/internal/store"
)

type Handler struct {
	store   *store.StorageService
	baseURL string
}

type URLCreationRequest struct {
	LongURL string `json:"long_url" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

func New(storageService *store.StorageService, baseURL string) *Handler {
	trimmedBaseURL := strings.TrimSuffix(baseURL, "/")
	return &Handler{
		store:   storageService,
		baseURL: trimmedBaseURL,
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hey Go URL Shortener!",
	})
}

func (h *Handler) CreateShortURL(c *gin.Context) {
	var creationRequest URLCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortURL := shortener.GenerateShortLink(creationRequest.LongURL, creationRequest.UserID)
	if err := h.store.SaveURLMapping(shortURL, creationRequest.LongURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save URL mapping"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "short url created successfully",
		"short_url": h.baseURL + "/" + shortURL,
	})
}

func (h *Handler) HandleShortURLRedirect(c *gin.Context) {
	shortURL := c.Param("shortUrl")
	initialURL, err := h.store.RetrieveInitialURL(shortURL)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusNotFound, gin.H{"error": "short URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve URL"})
		return
	}

	c.Redirect(http.StatusFound, initialURL)
}
