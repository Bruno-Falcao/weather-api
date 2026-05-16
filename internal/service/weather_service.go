package service

import (
	"context"
	"log/slog"
	"weather-api/internal/model"
)

type WeatherProvider interface {
	GetCurrentWeather(ctx context.Context, location string) (*model.WeatherData, error)
}

type WeatherCacher interface {
	GetWeather(ctx context.Context, location string) (*model.WeatherData, error)
	SetWeather(ctx context.Context, location string, data *model.WeatherData) error
}

type WeatherService struct {
	provider WeatherProvider
	cache    WeatherCacher
}

func NewWeatherService(provider WeatherProvider, cache WeatherCacher) *WeatherService {
	return &WeatherService{provider: provider, cache: cache}
}

func (s *WeatherService) GetWeather(ctx context.Context, location string) (*model.WeatherData, error) {
	cached, err := s.cache.GetWeather(ctx, location)
	if err != nil {
		slog.Warn("weather cache get failed", "err", err)
	}
	if cached != nil {
		slog.Info("weather cache hit", "location", location)
		return cached, nil
	}

	slog.Info("weather cache miss", "location", location)
	data, err := s.provider.GetCurrentWeather(ctx, location)
	if err != nil {
		return nil, err
	}
	if err := s.cache.SetWeather(ctx, location, data); err != nil {
		slog.Warn("weather cache set failed", "err", err)
	}
	return data, nil
}
