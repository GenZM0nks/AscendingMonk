package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type weatherRoundTripper func(request *http.Request) (*http.Response, error)

func (roundTripper weatherRoundTripper) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	return roundTripper(request)
}

func TestFetchWeatherData(t *testing.T) {
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
	})

	weatherHTTPClient = &http.Client{
		Transport: weatherRoundTripper(func(
			_ *http.Request,
		) (*http.Response, error) {
			responseBody := `{
			"daily": {
				"time": ["2026-09-22"],
				"weather_code": [2],
				"temperature_2m_max": [18.5],
				"temperature_2m_min": [11.2],
				"precipitation_probability_max": [25],
				"wind_speed_10m_max": [7.4],
				"wind_direction_10m_dominant": [225]
			}
		}`

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		}),
	}

	weatherData, err := fetchWeatherData()
	if err != nil {
		t.Fatalf("fetchWeatherData() returned an unexpected error: %v", err)
	}

	if weatherData.Location != weatherLocation {
		t.Errorf(
			"expected location %q, got %q",
			weatherLocation,
			weatherData.Location,
		)
	}

	if len(weatherData.Days) != 1 {
		t.Fatalf("expected 1 weather day, got %d", len(weatherData.Days))
	}

	day := weatherData.Days[0]

	if day.Date != "2026-09-22" {
		t.Errorf("expected date %q, got %q", "2026-09-22", day.Date)
	}

	if day.WeatherCode != 2 {
		t.Errorf("expected weather code 2, got %d", day.WeatherCode)
	}

	if day.TemperatureMaximum != 18.5 {
		t.Errorf(
			"expected maximum temperature 18.5, got %v",
			day.TemperatureMaximum,
		)
	}

	if day.TemperatureMinimum != 11.2 {
		t.Errorf(
			"expected minimum temperature 11.2, got %v",
			day.TemperatureMinimum,
		)
	}

	if day.PrecipitationProbability != 25 {
		t.Errorf(
			"expected precipitation probability 25, got %d",
			day.PrecipitationProbability,
		)
	}

	if day.WindSpeedMaximum != 7.4 {
		t.Errorf(
			"expected maximum wind speed 7.4, got %v",
			day.WindSpeedMaximum,
		)
	}

	if day.WindDirectionDominant != 225 {
		t.Errorf(
			"expected dominant wind direction 225, got %d",
			day.WindDirectionDominant,
		)
	}
}

func TestFetchWeatherDataNonOKResponse(t *testing.T) {
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
	})

	weatherHTTPClient = &http.Client{
		Transport: weatherRoundTripper(func(
			_ *http.Request,
		) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		}),
	}

	_, err := fetchWeatherData()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestFetchWeatherDataMalformedJSON(t *testing.T) {
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
	})

	weatherHTTPClient = &http.Client{
		Transport: weatherRoundTripper(func(
			_ *http.Request,
		) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(
					strings.NewReader(`{"daily": invalid}`),
				),
			}, nil
		}),
	}

	_, err := fetchWeatherData()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestFetchWeatherDataInconsistentDailyData(t *testing.T) {
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
	})

	weatherHTTPClient = &http.Client{
		Transport: weatherRoundTripper(func(
			_ *http.Request,
		) (*http.Response, error) {
			responseBody := `{
				"daily": {
					"time": ["2026-09-22", "2026-09-23"],
					"weather_code": [2],
					"temperature_2m_max": [18.5],
					"temperature_2m_min": [11.2],
					"precipitation_probability_max": [25],
					"wind_speed_10m_max": [7.4],
					"wind_direction_10m_dominant": [225]
				}
			}`

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		}),
	}

	_, err := fetchWeatherData()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestFetchWeatherDataRequestError(t *testing.T) {
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
	})

	weatherHTTPClient = &http.Client{
		Transport: weatherRoundTripper(func(
			_ *http.Request,
		) (*http.Response, error) {
			return nil, errors.New("test request failure")
		}),
	}

	_, err := fetchWeatherData()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestAPIWeather(t *testing.T) {
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
	})

	weatherHTTPClient = &http.Client{
		Transport: weatherRoundTripper(func(
			_ *http.Request,
		) (*http.Response, error) {
			responseBody := `{
				"daily": {
					"time": ["2026-09-22"],
					"weather_code": [2],
					"temperature_2m_max": [18.5],
					"temperature_2m_min": [11.2],
					"precipitation_probability_max": [25],
					"wind_speed_10m_max": [7.4],
					"wind_direction_10m_dominant": [225]
				}
			}`

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		}),
	}

	request := httptest.NewRequest(http.MethodGet, "/api/weather", nil)
	responseRecorder := httptest.NewRecorder()

	APIWeather(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf(
			"expected status code %d, got %d",
			http.StatusOK,
			responseRecorder.Code,
		)
	}

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf(
			"expected Content-Type %q, got %q",
			"application/json",
			contentType,
		)
	}

	responseBody := responseRecorder.Body.String()

	if !strings.Contains(responseBody, `"data"`) {
		t.Errorf("expected response body to contain data field, got %s", responseBody)
	}

	if !strings.Contains(responseBody, `"location":"Copenhagen, Denmark"`) {
		t.Errorf(
			"expected response body to contain Copenhagen location, got %s",
			responseBody,
		)
	}
}
