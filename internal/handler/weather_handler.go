package handler

import (
	"context"
	"net/http"
	"weather-api/internal/model"
	"weather-api/internal/response"

	"github.com/go-chi/chi/v5"
)

type weatherService interface {
	GetWeather(ctx context.Context, location string) (*model.WeatherData, error)
}

type WeatherHandler struct {
	svc weatherService
}

func NewWeatherHandler(svc weatherService) *WeatherHandler {
	return &WeatherHandler{svc: svc}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	location := chi.URLParam(r, "location")
	if location == "" {
		response.Error(w, http.StatusBadRequest, "location is required")
		return
	}

	data, err := h.svc.GetWeather(r.Context(), location)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, data)
}
