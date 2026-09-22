package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GenZM0nks/AscendingMonk/internal/templates"
)

const (
	weatherLocation  = "Copenhagen, Denmark"
	weatherLatitude  = 55.6761
	weatherLongitude = 12.5683
)

var weatherHTTPClient = &http.Client{
	Timeout: 6 * time.Second,
}

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

	response, err := weatherHTTPClient.Get(requestURL)
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

func formatWeatherDate(date string) string {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}

	return parsedDate.Format("Monday, 2 January")
}

func weatherDescription(code int) string {
	switch code {
	case 0:
		return "Clear sky"
	case 1:
		return "Mainly clear"
	case 2:
		return "Partly cloudy"
	case 3:
		return "Overcast"
	case 45:
		return "Fog"
	case 48:
		return "Depositing rime fog"
	case 51:
		return "Light drizzle"
	case 53:
		return "Moderate drizzle"
	case 55:
		return "Dense drizzle"
	case 56:
		return "Light freezing drizzle"
	case 57:
		return "Dense freezing drizzle"
	case 61:
		return "Slight rain"
	case 63:
		return "Moderate rain"
	case 65:
		return "Heavy rain"
	case 66:
		return "Light freezing rain"
	case 67:
		return "Heavy freezing rain"
	case 71:
		return "Slight snowfall"
	case 73:
		return "Moderate snowfall"
	case 75:
		return "Heavy snowfall"
	case 77:
		return "Snow grains"
	case 80:
		return "Slight rain showers"
	case 81:
		return "Moderate rain showers"
	case 82:
		return "Violent rain showers"
	case 85:
		return "Slight snow showers"
	case 86:
		return "Heavy snow showers"
	case 95:
		return "Thunderstorm"
	case 96:
		return "Thunderstorm with slight hail"
	case 99:
		return "Thunderstorm with heavy hail"
	default:
		return "Unknown"
	}
}

func windDirection(degrees int) string {
	directions := []string{
		"N", "NNE", "NE", "ENE",
		"E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW",
		"W", "WNW", "NW", "NNW",
	}

	index := ((degrees + 11) / 22) % len(directions)

	return directions[index]
}

// Weather renders the weather forecast page.
func Weather(responseWriter http.ResponseWriter, _ *http.Request) {
	weatherData, err := fetchWeatherData()
	if err != nil {
		http.Error(
			responseWriter,
			"Failed to fetch weather forecast",
			http.StatusInternalServerError,
		)
		return
	}

	days := make([]WeatherDayView, 0, len(weatherData.Days))

	for _, day := range weatherData.Days {
		days = append(days, WeatherDayView{
			Date:                     formatWeatherDate(day.Date),
			TemperatureMaximum:       day.TemperatureMaximum,
			TemperatureMinimum:       day.TemperatureMinimum,
			PrecipitationProbability: day.PrecipitationProbability,
			WindSpeedMaximum:         day.WindSpeedMaximum,
			WindDirection:            windDirection(day.WindDirectionDominant),
			WeatherDescription:       weatherDescription(day.WeatherCode),
		})
	}

	pageData := WeatherPageData{
		Data: PageData{
			PageTitle: "Weather",
			Flashes:   nil,
		},
		Location: weatherData.Location,
		Days:     days,
	}

	templates.LoadAndExecuteTemplate(
		"web/templates/weather.html",
		pageData,
		responseWriter,
	)
}
