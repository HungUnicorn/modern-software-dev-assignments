// Package openmeteo provides a client for the Open-Meteo weather API.
package openmeteo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	forecastBaseURL  = "https://api.open-meteo.com/v1/forecast"
	geocodingBaseURL = "https://geocoding-api.open-meteo.com/v1/search"
	archiveBaseURL   = "https://archive-api.open-meteo.com/v1/archive"

	minRequestInterval = 150 * time.Millisecond
	httpTimeout        = 10 * time.Second
)

// Client is a rate-limited HTTP client for the Open-Meteo API.
type Client struct {
	httpClient  *http.Client
	mu          sync.Mutex
	lastRequest time.Time
}

// NewClient creates a new Open-Meteo API client.
func NewClient() *Client {
	return &Client{httpClient: &http.Client{Timeout: httpTimeout}}
}

func (c *Client) throttle() {
	c.mu.Lock()
	defer c.mu.Unlock()
	elapsed := time.Since(c.lastRequest)
	if elapsed < minRequestInterval {
		wait := minRequestInterval - elapsed
		log.Printf("[rate-limit] Waiting %v before next request", wait)
		time.Sleep(wait)
	}
	c.lastRequest = time.Now()
}

func (c *Client) doGet(requestURL string) ([]byte, error) {
	c.throttle()
	log.Printf("[http] GET %s", requestURL)

	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limit exceeded (HTTP 429) — Open-Meteo allows ~10,000 requests/day on the free tier")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned HTTP %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func validateCoordinates(lat, lon float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90, got %.4f", lat)
	}
	if lon < -180 || lon > 180 {
		return fmt.Errorf("longitude must be between -180 and 180, got %.4f", lon)
	}
	return nil
}

// WeatherCodeDescription returns a human-readable WMO weather code description.
func WeatherCodeDescription(code int) string {
	m := map[int]string{
		0: "Clear sky", 1: "Mainly clear", 2: "Partly cloudy", 3: "Overcast",
		45: "Fog", 48: "Depositing rime fog",
		51: "Light drizzle", 53: "Moderate drizzle", 55: "Dense drizzle",
		56: "Light freezing drizzle", 57: "Dense freezing drizzle",
		61: "Slight rain", 63: "Moderate rain", 65: "Heavy rain",
		66: "Light freezing rain", 67: "Heavy freezing rain",
		71: "Slight snowfall", 73: "Moderate snowfall", 75: "Heavy snowfall", 77: "Snow grains",
		80: "Slight rain showers", 81: "Moderate rain showers", 82: "Violent rain showers",
		85: "Slight snow showers", 86: "Heavy snow showers",
		95: "Thunderstorm", 96: "Thunderstorm with slight hail", 99: "Thunderstorm with heavy hail",
	}
	if desc, ok := m[code]; ok {
		return desc
	}
	return fmt.Sprintf("Unknown (code %d)", code)
}

// --- Geocoding ---

// GeocodingResult represents a single location result.
type GeocodingResult struct {
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Admin1      string  `json:"admin1"`
	Population  int     `json:"population"`
	Elevation   float64 `json:"elevation"`
	Timezone    string  `json:"timezone"`
}

type geocodingResponse struct {
	Results []GeocodingResult `json:"results"`
}

// SearchLocation searches for locations by name.
func (c *Client) SearchLocation(name string, count int) ([]GeocodingResult, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("location name must not be empty")
	}
	if count < 1 {
		count = 5
	}
	if count > 10 {
		count = 10
	}

	params := url.Values{}
	params.Set("name", name)
	params.Set("count", fmt.Sprintf("%d", count))
	params.Set("language", "en")
	params.Set("format", "json")

	body, err := c.doGet(fmt.Sprintf("%s?%s", geocodingBaseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("geocoding request failed: %w", err)
	}

	var result geocodingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse geocoding response: %w", err)
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("no locations found for %q", name)
	}
	return result.Results, nil
}

// --- Forecast ---

// ForecastParams holds forecast request parameters.
type ForecastParams struct {
	Latitude        float64
	Longitude       float64
	Days            int
	TemperatureUnit string
	WindspeedUnit   string
	Timezone        string
}

// ForecastResponse holds the parsed forecast data.
type ForecastResponse struct {
	Latitude  float64         `json:"latitude"`
	Longitude float64         `json:"longitude"`
	Timezone  string          `json:"timezone"`
	Daily     *DailyForecast  `json:"daily"`
	Current   *CurrentWeather `json:"current"`
}

// DailyForecast holds daily forecast arrays.
type DailyForecast struct {
	Time             []string  `json:"time"`
	WeatherCode      []int     `json:"weather_code"`
	Temperature2mMax []float64 `json:"temperature_2m_max"`
	Temperature2mMin []float64 `json:"temperature_2m_min"`
	PrecipitationSum []float64 `json:"precipitation_sum"`
	WindSpeed10mMax  []float64 `json:"wind_speed_10m_max"`
	Sunrise          []string  `json:"sunrise"`
	Sunset           []string  `json:"sunset"`
}

// CurrentWeather holds current conditions.
type CurrentWeather struct {
	Time          string  `json:"time"`
	Temperature2m float64 `json:"temperature_2m"`
	WeatherCode   int     `json:"weather_code"`
	WindSpeed10m  float64 `json:"wind_speed_10m"`
	WindDirection int     `json:"wind_direction_10m"`
	Humidity      float64 `json:"relative_humidity_2m"`
}

// GetForecast fetches a weather forecast.
func (c *Client) GetForecast(p ForecastParams) (*ForecastResponse, error) {
	if err := validateCoordinates(p.Latitude, p.Longitude); err != nil {
		return nil, err
	}
	if p.Days < 1 {
		p.Days = 3
	}
	if p.Days > 16 {
		p.Days = 16
	}
	if p.TemperatureUnit == "" {
		p.TemperatureUnit = "celsius"
	}
	if p.WindspeedUnit == "" {
		p.WindspeedUnit = "kmh"
	}
	if p.Timezone == "" {
		p.Timezone = "auto"
	}

	q := url.Values{}
	q.Set("latitude", fmt.Sprintf("%.4f", p.Latitude))
	q.Set("longitude", fmt.Sprintf("%.4f", p.Longitude))
	q.Set("forecast_days", fmt.Sprintf("%d", p.Days))
	q.Set("temperature_unit", p.TemperatureUnit)
	q.Set("wind_speed_unit", p.WindspeedUnit)
	q.Set("timezone", p.Timezone)
	q.Set("current", "temperature_2m,weather_code,wind_speed_10m,wind_direction_10m,relative_humidity_2m")
	q.Set("daily", "weather_code,temperature_2m_max,temperature_2m_min,precipitation_sum,wind_speed_10m_max,sunrise,sunset")

	body, err := c.doGet(fmt.Sprintf("%s?%s", forecastBaseURL, q.Encode()))
	if err != nil {
		return nil, fmt.Errorf("forecast request failed: %w", err)
	}

	var result ForecastResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse forecast response: %w", err)
	}
	return &result, nil
}

// --- Historical Weather ---

// HistoricalParams holds historical weather request parameters.
type HistoricalParams struct {
	Latitude        float64
	Longitude       float64
	StartDate       string
	EndDate         string
	TemperatureUnit string
	Timezone        string
}

// HistoricalResponse holds parsed historical weather data.
type HistoricalResponse struct {
	Latitude  float64          `json:"latitude"`
	Longitude float64          `json:"longitude"`
	Timezone  string           `json:"timezone"`
	Daily     *HistoricalDaily `json:"daily"`
}

// HistoricalDaily holds daily historical data arrays.
type HistoricalDaily struct {
	Time             []string  `json:"time"`
	WeatherCode      []int     `json:"weather_code"`
	Temperature2mMax []float64 `json:"temperature_2m_max"`
	Temperature2mMin []float64 `json:"temperature_2m_min"`
	PrecipitationSum []float64 `json:"precipitation_sum"`
	WindSpeed10mMax  []float64 `json:"wind_speed_10m_max"`
}

// GetHistoricalWeather fetches historical weather data.
func (c *Client) GetHistoricalWeather(p HistoricalParams) (*HistoricalResponse, error) {
	if err := validateCoordinates(p.Latitude, p.Longitude); err != nil {
		return nil, err
	}
	if p.StartDate == "" || p.EndDate == "" {
		return nil, fmt.Errorf("start_date and end_date are required (YYYY-MM-DD)")
	}
	start, err := time.Parse("2006-01-02", p.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}
	end, err := time.Parse("2006-01-02", p.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format: %w", err)
	}
	if end.Before(start) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}
	if end.Sub(start) > 366*24*time.Hour {
		return nil, fmt.Errorf("date range must not exceed 1 year")
	}
	if p.TemperatureUnit == "" {
		p.TemperatureUnit = "celsius"
	}
	if p.Timezone == "" {
		p.Timezone = "auto"
	}

	q := url.Values{}
	q.Set("latitude", fmt.Sprintf("%.4f", p.Latitude))
	q.Set("longitude", fmt.Sprintf("%.4f", p.Longitude))
	q.Set("start_date", p.StartDate)
	q.Set("end_date", p.EndDate)
	q.Set("temperature_unit", p.TemperatureUnit)
	q.Set("timezone", p.Timezone)
	q.Set("daily", "weather_code,temperature_2m_max,temperature_2m_min,precipitation_sum,wind_speed_10m_max")

	body, err := c.doGet(fmt.Sprintf("%s?%s", archiveBaseURL, q.Encode()))
	if err != nil {
		return nil, fmt.Errorf("historical weather request failed: %w", err)
	}

	var result HistoricalResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse historical weather response: %w", err)
	}
	return &result, nil
}
