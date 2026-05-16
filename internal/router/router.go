package router

import (
	"log/slog"
	"weather-api/internal/config"
	"weather-api/internal/handler"
	"weather-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Dependencies struct {
	WeatherHandler *handler.WeatherHandler
	Logger         *slog.Logger
	Config         *config.Config
}

func New(deps Dependencies) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger(deps.Logger))
	r.Use(middleware.RateLimit(deps.Config.RateLimitRequests, deps.Config.RateLimitWindow))

	r.Route("/api/v1", func(r chi.Router) {
		mountWeather(r, deps.WeatherHandler)
	})

	return r
}
