package model

type WeatherData struct {
	Location string        `json:"location"`
	Days     []ForecastDay `json:"days"`
}
type ForecastDay struct {
	Date          string  `json:"datetime"`
	Temperature   float64 `json:"temp"`
	Precipitation float64 `json:"precip"`
}
