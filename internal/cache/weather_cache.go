package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"weather-api/internal/model"

	"github.com/redis/go-redis/v9"
)

const weatherTTL = 10 * time.Minute

type WeatherCache struct {
	client *redis.Client
}

func NewWeatherCache(client *redis.Client) *WeatherCache {
	return &WeatherCache{client: client}
}

func (c *WeatherCache) GetWeather(ctx context.Context, location string) (*model.WeatherData, error) {
	val, err := c.client.Get(ctx, location).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redus get: %w", err)
	}

	var data model.WeatherData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("unmarshaling weather data: %w", err)
	}
	return &data, nil
}

func (c *WeatherCache) SetWeather(ctx context.Context, location string, data *model.WeatherData) error {
	val, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling weather data: %w", err)
	}
	return c.client.Set(ctx, location, val, weatherTTL).Err()
}
