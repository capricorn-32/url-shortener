package handler

import (
	"errors"
	"net/http"
	"net/url"
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

	if err := validateLongURL(creationRequest.LongURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const maxGenerateAttempts = 5
	for attempt := 0; attempt < maxGenerateAttempts; attempt++ {
		shortURL := shortener.GenerateShortLinkWithSalt(creationRequest.LongURL, creationRequest.UserID, attempt)
		if shortURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate short URL"})
			return
		}

		err := h.store.SaveURLMapping(shortURL, creationRequest.LongURL)
		if err == nil {
			h.respondCreateSuccess(c, shortURL)
			return
		}

		if errors.Is(err, store.ErrShortURLExists) {
			existingURL, getErr := h.store.RetrieveInitialURL(shortURL)
			if getErr != nil && !errors.Is(getErr, redis.Nil) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify existing short URL"})
				return
			}

			if getErr == nil && existingURL == creationRequest.LongURL {
				h.respondCreateSuccess(c, shortURL)
				return
			}

			continue
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save URL mapping"})
		return
	}

	c.JSON(http.StatusConflict, gin.H{"error": "failed to generate unique short URL"})
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

func validateLongURL(rawURL string) error {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return errors.New("long_url must be a valid URL")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("long_url must use http or https scheme")
	}

	if parsedURL.Host == "" {
		return errors.New("long_url must include a valid host")
	}

	return nil
}

func (h *Handler) respondCreateSuccess(c *gin.Context, shortURL string) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "short url created successfully",
		"short_url": h.baseURL + "/" + shortURL,
	})
}
