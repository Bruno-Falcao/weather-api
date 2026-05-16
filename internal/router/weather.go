package router

import (
	"weather-api/internal/handler"

	"github.com/go-chi/chi/v5"
)

func mountWeather(r chi.Router, h *handler.WeatherHandler) {
	r.Route("/weather", func(r chi.Router) {
		r.Get("/{location}", h.GetWeather)
	})
}
