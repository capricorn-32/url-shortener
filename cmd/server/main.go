package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	shutdownSignalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("starting server on port %s", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start web server: %v", err)
		}
	}()

	<-shutdownSignalCtx.Done()
	log.Printf("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	if err := storageService.Close(); err != nil {
		log.Printf("failed closing redis client: %v", err)
	}
}
