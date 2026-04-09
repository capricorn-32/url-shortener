package main

import (
	"log"

	"github.com/seniorLikeToCode/url-shortener/internal/config"
	"github.com/seniorLikeToCode/url-shortener/internal/handler"
	"github.com/seniorLikeToCode/url-shortener/internal/server"
	"github.com/seniorLikeToCode/url-shortener/internal/store"
)

func main() {
	cfg := config.Load()

	storageService, err := store.NewStorageService(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	h := handler.New(storageService, cfg.BaseURL)
	r := server.NewRouter(h)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start web server: %v", err)
	}
}
