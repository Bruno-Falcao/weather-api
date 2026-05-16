package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
	"weather-api/internal/model"
)

const baseURL = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline"

type vcResponse struct {
	ResolvedAddress string `json:"resolvedAddress"`
	Days            []struct {
		Date          string  `json:"datetime"`
		Temperature   float64 `json:"temp"`
		Precipitation float64 `json:"precip"`
	} `json:"days"`
}

type VisualCrossingClient struct {
	http   *http.Client
	apiKey string
}

func NewVisualCrossing(apiKey string) *VisualCrossingClient {
	return &VisualCrossingClient{
		apiKey: apiKey,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *VisualCrossingClient) GetCurrentWeather(ctx context.Context, location string) (*model.WeatherData, error) {
	resp, err := c.fetchTimeline(ctx, location)
	if err != nil {
		return nil, err
	}

	weatherData := &model.WeatherData{Location: resp.ResolvedAddress}
	for _, day := range resp.Days {
		weatherData.Days = append(weatherData.Days, model.ForecastDay{
			Temperature:   day.Temperature,
			Date:          day.Date,
			Precipitation: day.Precipitation,
		},
		)
	}

	return weatherData, nil
}

// fetchTimeline fetches the weather data from the Visual Crossing API.
func (c *VisualCrossingClient) fetchTimeline(ctx context.Context, location string) (*vcResponse, error) {
	endpoint := fmt.Sprintf("%s/%s/",
		baseURL,
		url.PathEscape(location),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	q := req.URL.Query()
	q.Set("key", c.apiKey)
	q.Set("unitGroup", "metric")
	q.Set("include", "current")
	q.Set("contentType", "json")

	req.URL.RawQuery = q.Encode()

	slog.Info("visual crossing request", "url", req.URL.String(), "key_length", len(c.apiKey))
	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("visual crossing API error: status %d", res.StatusCode)
	}

	var payload vcResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &payload, nil
}
