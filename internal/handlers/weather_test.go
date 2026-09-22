package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type weatherRoundTripper func(request *http.Request) (*http.Response, error)

func (roundTripper weatherRoundTripper) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	return roundTripper(request)
}

func resetWeatherCache() {
	weatherCache.Lock()
	defer weatherCache.Unlock()

	weatherCache.data = WeatherData{}
	weatherCache.expiresAt = time.Time{}
	weatherCache.hasData = false
}

func TestFetchWeatherData(t *testing.T) {
	resetWeatherCache()
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
		resetWeatherCache()
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

func TestFetchWeatherDataUsesCache(t *testing.T) {
	resetWeatherCache()

	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
		resetWeatherCache()
	})

	requestCount := 0

	weatherHTTPClient = &http.Client{
		Transport: weatherRoundTripper(func(
			_ *http.Request,
		) (*http.Response, error) {
			requestCount++

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

	_, err := fetchWeatherData()
	if err != nil {
		t.Fatalf("first fetchWeatherData() returned an unexpected error: %v", err)
	}

	_, err = fetchWeatherData()
	if err != nil {
		t.Fatalf("second fetchWeatherData() returned an unexpected error: %v", err)
	}

	if requestCount != 1 {
		t.Errorf(
			"expected 1 request to Open-Meteo, got %d",
			requestCount,
		)
	}
}

func TestFetchWeatherDataRefreshesExpiredCache(t *testing.T) {
	resetWeatherCache()

	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
		resetWeatherCache()
	})

	requestCount := 0

	weatherHTTPClient = &http.Client{
		Transport: weatherRoundTripper(func(
			_ *http.Request,
		) (*http.Response, error) {
			requestCount++

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

	_, err := fetchWeatherData()
	if err != nil {
		t.Fatalf("first fetchWeatherData() returned an unexpected error: %v", err)
	}

	weatherCache.Lock()
	weatherCache.expiresAt = time.Now().Add(-time.Minute)
	weatherCache.Unlock()

	_, err = fetchWeatherData()
	if err != nil {
		t.Fatalf("second fetchWeatherData() returned an unexpected error: %v", err)
	}

	if requestCount != 2 {
		t.Errorf(
			"expected 2 requests to Open-Meteo, got %d",
			requestCount,
		)
	}
}

func TestFetchWeatherDataNonOKResponse(t *testing.T) {
	resetWeatherCache()
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
		resetWeatherCache()
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
	resetWeatherCache()
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
		resetWeatherCache()
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
	resetWeatherCache()
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
		resetWeatherCache()
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
	resetWeatherCache()
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
		resetWeatherCache()
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
	resetWeatherCache()
	originalClient := weatherHTTPClient
	t.Cleanup(func() {
		weatherHTTPClient = originalClient
		resetWeatherCache()
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

func TestWeatherHTTPClientTimeout(t *testing.T) {
	expectedTimeout := 6 * time.Second

	if weatherHTTPClient.Timeout != expectedTimeout {
		t.Errorf(
			"expected weather HTTP client timeout %v, got %v",
			expectedTimeout,
			weatherHTTPClient.Timeout,
		)
	}
}
