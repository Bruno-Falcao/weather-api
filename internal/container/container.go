package container

import (
	"fmt"
	"log/slog"
	"weather-api/internal/cache"
	"weather-api/internal/client"
	"weather-api/internal/config"
	"weather-api/internal/handler"
	"weather-api/internal/router"
	"weather-api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	Router      *chi.Mux
	RedisClient *redis.Client
}

func New(cfg *config.Config, logger *slog.Logger) (*Container, error) {
	// infra
	redisClient, err := cache.NewRedisClient(cfg.RedisAddr)
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	// external clients
	vcClient := client.NewVisualCrossing(cfg.VisualCrossingAPIKey)

	// caches
	weatherCache := cache.NewWeatherCache(redisClient)

	// services
	weatherSvc := service.NewWeatherService(vcClient, weatherCache)

	// handlers
	deps := router.Dependencies{
		WeatherHandler: handler.NewWeatherHandler(weatherSvc),
		Logger:         logger,
		Config:         cfg,
	}

	return &Container{
		Router:      router.New(deps),
		RedisClient: redisClient,
	}, nil
}

func (c *Container) Close() {
	if err := c.RedisClient.Close(); err != nil {
		slog.Error("closing redis", "err", err)
	}
}
