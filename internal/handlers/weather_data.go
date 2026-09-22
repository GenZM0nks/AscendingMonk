package handlers

// StandardResponse represents the standard API response containing weather data.
type StandardResponse struct {
	Data WeatherData `json:"data"`
}

// WeatherData contains the weather forecast for a location.
type WeatherData struct {
	Location string       `json:"location"`
	Days     []WeatherDay `json:"days"`
}

// WeatherDay contains the forecast for one day.
type WeatherDay struct {
	Date                     string  `json:"date"`
	TemperatureMaximum       float64 `json:"temperatureMaximum"`
	TemperatureMinimum       float64 `json:"temperatureMinimum"`
	PrecipitationProbability int     `json:"precipitationProbability"`
	WindSpeedMaximum         float64 `json:"windSpeedMaximum"`
	WindDirectionDominant    int     `json:"windDirectionDominant"`
	WeatherCode              int     `json:"weatherCode"`
}

// WeatherPageData contains the data needed to render the weather page.
type WeatherPageData struct {
	Data     PageData
	Location string
	Days     []WeatherDayView
}

// WeatherDayView contains display-friendly weather data for one day.
type WeatherDayView struct {
	Date                     string
	TemperatureMaximum       float64
	TemperatureMinimum       float64
	PrecipitationProbability int
	WindSpeedMaximum         float64
	WindDirection            string
	WeatherDescription       string
}
