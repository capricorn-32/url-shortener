package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type StorageService struct {
	redisClient *redis.Client
	ctx         context.Context
}

const CacheDuration = 6 * time.Hour

var ErrShortURLExists = errors.New("short URL already exists")

func NewStorageService(addr, password string, db int) (*StorageService, error) {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &StorageService{
		redisClient: redisClient,
		ctx:         ctx,
	}, nil
}

func (s *StorageService) SaveURLMapping(shortURL, originalURL string) error {
	ok, err := s.redisClient.SetNX(s.ctx, shortURL, originalURL, CacheDuration).Result()
	if err != nil {
		return fmt.Errorf("save URL mapping: %w", err)
	}
	if !ok {
		return ErrShortURLExists
	}
	return nil
}

func (s *StorageService) RetrieveInitialURL(shortURL string) (string, error) {
	result, err := s.redisClient.Get(s.ctx, shortURL).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *StorageService) Close() error {
	if err := s.redisClient.Close(); err != nil {
		return fmt.Errorf("close redis client: %w", err)
	}
	return nil
}
