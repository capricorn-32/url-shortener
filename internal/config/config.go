package config

import "os"

type Config struct {
	Port          string
	BaseURL       string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() Config {
	return Config{
		Port:          getEnv("APP_PORT", "5000"),
		BaseURL:       getEnv("APP_BASE_URL", "http://localhost:5000"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
