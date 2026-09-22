package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	weatherLocation  = "Copenhagen, Denmark"
	weatherLatitude  = 55.6761
	weatherLongitude = 12.5683
)

type openMeteoResponse struct {
	Daily openMeteoDailyData `json:"daily"`
}

type openMeteoDailyData struct {
	Time                     []string  `json:"time"`
	TemperatureMaximum       []float64 `json:"temperature_2m_max"`
	TemperatureMinimum       []float64 `json:"temperature_2m_min"`
	PrecipitationProbability []int     `json:"precipitation_probability_max"`
	WindSpeedMaximum         []float64 `json:"wind_speed_10m_max"`
	WindDirectionDominant    []int     `json:"wind_direction_10m_dominant"`
	WeatherCode              []int     `json:"weather_code"`
}

func fetchWeatherData() (WeatherData, error) {
	requestURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast"+
			"?latitude=%f"+
			"&longitude=%f"+
			"&daily=weather_code,temperature_2m_max,temperature_2m_min,"+
			"precipitation_probability_max,wind_speed_10m_max,"+
			"wind_direction_10m_dominant"+
			"&temperature_unit=celsius"+
			"&wind_speed_unit=ms"+
			"&forecast_days=7",
		weatherLatitude,
		weatherLongitude,
	)

	response, err := http.Get(requestURL)
	if err != nil {
		return WeatherData{}, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return WeatherData{}, fmt.Errorf(
			"weather service returned status %d",
			response.StatusCode,
		)
	}

	var openMeteoData openMeteoResponse

	err = json.NewDecoder(response.Body).Decode(&openMeteoData)
	if err != nil {
		return WeatherData{}, fmt.Errorf("failed to decode weather data: %w", err)
	}

	dayCount := len(openMeteoData.Daily.Time)

	if len(openMeteoData.Daily.TemperatureMaximum) != dayCount ||
		len(openMeteoData.Daily.TemperatureMinimum) != dayCount ||
		len(openMeteoData.Daily.PrecipitationProbability) != dayCount ||
		len(openMeteoData.Daily.WindSpeedMaximum) != dayCount ||
		len(openMeteoData.Daily.WindDirectionDominant) != dayCount ||
		len(openMeteoData.Daily.WeatherCode) != dayCount {
		return WeatherData{}, fmt.Errorf(
			"weather service returned inconsistent daily data",
		)
	}

	days := make([]WeatherDay, 0, dayCount)

	for index, date := range openMeteoData.Daily.Time {
		day := WeatherDay{
			Date:                     date,
			TemperatureMaximum:       openMeteoData.Daily.TemperatureMaximum[index],
			TemperatureMinimum:       openMeteoData.Daily.TemperatureMinimum[index],
			PrecipitationProbability: openMeteoData.Daily.PrecipitationProbability[index],
			WindSpeedMaximum:         openMeteoData.Daily.WindSpeedMaximum[index],
			WindDirectionDominant:    openMeteoData.Daily.WindDirectionDominant[index],
			WeatherCode:              openMeteoData.Daily.WeatherCode[index],
		}

		days = append(days, day)
	}

	return WeatherData{
		Location: weatherLocation,
		Days:     days,
	}, nil
}

// APIWeather returns the weather forecast as JSON.
//
// @Summary Fetch weather forecast.
// @Description Fetch a seven-day weather forecast for Copenhagen.
// @Produce json
// @Success 200 {object} StandardResponse
// @Failure 500 {string} string "error"
// @Tags API
// @Router /weather [get]
func APIWeather(responseWriter http.ResponseWriter, _ *http.Request) {
	weatherData, err := fetchWeatherData()
	if err != nil {
		http.Error(
			responseWriter,
			"Failed to fetch weather forecast",
			http.StatusInternalServerError,
		)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(responseWriter).Encode(StandardResponse{
		Data: weatherData,
	})
	if err != nil {
		http.Error(
			responseWriter,
			"Failed to encode weather forecast",
			http.StatusInternalServerError,
		)
	}
}
