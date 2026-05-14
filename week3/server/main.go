// Weather MCP Server — wraps the Open-Meteo API and exposes it via MCP over STDIO.
//
// Tools provided:
//   - search_location: Find coordinates for a city/place name.
//   - get_forecast:    Get current conditions + multi-day forecast.
//   - get_historical_weather: Get past weather data for a date range.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/holm/weather-mcp-server/openmeteo"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Shared Open-Meteo client (rate-limited, reused across tool calls).
var client = openmeteo.NewClient()

// --- Tool Input Types ---

// SearchLocationInput is the typed input for the search_location tool.
type SearchLocationInput struct {
	Name  string `json:"name"                    jsonschema:"City or place name to search for"`
	Count int    `json:"count,omitzero"           jsonschema:"Number of results to return (1-10 default 5)"`
}

// GetForecastInput is the typed input for the get_forecast tool.
type GetForecastInput struct {
	Latitude        float64 `json:"latitude"                jsonschema:"Latitude (-90 to 90)"`
	Longitude       float64 `json:"longitude"               jsonschema:"Longitude (-180 to 180)"`
	Days            int     `json:"days,omitzero"            jsonschema:"Forecast days (1-16 default 3)"`
	TemperatureUnit string  `json:"temperature_unit,omitempty" jsonschema:"Temperature unit: celsius or fahrenheit (default celsius)"`
}

// GetHistoricalInput is the typed input for the get_historical_weather tool.
type GetHistoricalInput struct {
	Latitude        float64 `json:"latitude"                jsonschema:"Latitude (-90 to 90)"`
	Longitude       float64 `json:"longitude"               jsonschema:"Longitude (-180 to 180)"`
	StartDate       string  `json:"start_date"              jsonschema:"Start date in YYYY-MM-DD format"`
	EndDate         string  `json:"end_date"                jsonschema:"End date in YYYY-MM-DD format"`
	TemperatureUnit string  `json:"temperature_unit,omitempty" jsonschema:"Temperature unit: celsius or fahrenheit (default celsius)"`
}

// --- Tool Handlers ---

// handleSearchLocation resolves a place name to geographic coordinates.
func handleSearchLocation(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SearchLocationInput,
) (*mcp.CallToolResult, any, error) {
	log.Printf("[tool] search_location: name=%q count=%d", input.Name, input.Count)

	results, err := client.SearchLocation(input.Name, input.Count)
	if err != nil {
		return errResult(fmt.Sprintf("Location search failed: %v", err)), nil, nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d location(s) for %q:\n\n", len(results), input.Name))
	for i, loc := range results {
		region := loc.Admin1
		if region == "" {
			region = "N/A"
		}
		sb.WriteString(fmt.Sprintf(
			"%d. %s, %s (%s)\n   Lat: %.4f, Lon: %.4f | Elevation: %.0fm | Pop: %d | TZ: %s\n\n",
			i+1, loc.Name, loc.Country, region,
			loc.Latitude, loc.Longitude, loc.Elevation, loc.Population, loc.Timezone,
		))
	}

	return textResult(sb.String()), nil, nil
}

// handleGetForecast returns current weather + daily forecast.
func handleGetForecast(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetForecastInput,
) (*mcp.CallToolResult, any, error) {
	log.Printf("[tool] get_forecast: lat=%.4f lon=%.4f days=%d unit=%s",
		input.Latitude, input.Longitude, input.Days, input.TemperatureUnit)

	forecast, err := client.GetForecast(openmeteo.ForecastParams{
		Latitude:        input.Latitude,
		Longitude:       input.Longitude,
		Days:            input.Days,
		TemperatureUnit: input.TemperatureUnit,
	})
	if err != nil {
		return errResult(fmt.Sprintf("Forecast request failed: %v", err)), nil, nil
	}

	tempUnit := "°C"
	if input.TemperatureUnit == "fahrenheit" {
		tempUnit = "°F"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Weather for (%.4f, %.4f) — TZ: %s\n",
		forecast.Latitude, forecast.Longitude, forecast.Timezone))

	// Current conditions
	if c := forecast.Current; c != nil {
		sb.WriteString(fmt.Sprintf(
			"\n── Current Conditions (%s) ──\n"+
				"  Temperature: %.1f%s\n"+
				"  Humidity:    %.0f%%\n"+
				"  Wind:        %.1f km/h (dir: %d°)\n"+
				"  Conditions:  %s\n",
			c.Time, c.Temperature2m, tempUnit,
			c.Humidity, c.WindSpeed10m, c.WindDirection,
			openmeteo.WeatherCodeDescription(c.WeatherCode),
		))
	}

	// Daily forecast
	if d := forecast.Daily; d != nil && len(d.Time) > 0 {
		sb.WriteString("\n── Daily Forecast ──\n")
		for i, date := range d.Time {
			sb.WriteString(fmt.Sprintf(
				"\n  %s — %s\n"+
					"    High: %.1f%s  Low: %.1f%s\n"+
					"    Precipitation: %.1f mm | Max Wind: %.1f km/h\n"+
					"    Sunrise: %s  Sunset: %s\n",
				date, openmeteo.WeatherCodeDescription(d.WeatherCode[i]),
				d.Temperature2mMax[i], tempUnit, d.Temperature2mMin[i], tempUnit,
				d.PrecipitationSum[i], d.WindSpeed10mMax[i],
				d.Sunrise[i], d.Sunset[i],
			))
		}
	}

	return textResult(sb.String()), nil, nil
}

// handleGetHistoricalWeather returns past weather data for a date range.
func handleGetHistoricalWeather(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetHistoricalInput,
) (*mcp.CallToolResult, any, error) {
	log.Printf("[tool] get_historical_weather: lat=%.4f lon=%.4f start=%s end=%s",
		input.Latitude, input.Longitude, input.StartDate, input.EndDate)

	hist, err := client.GetHistoricalWeather(openmeteo.HistoricalParams{
		Latitude:        input.Latitude,
		Longitude:       input.Longitude,
		StartDate:       input.StartDate,
		EndDate:         input.EndDate,
		TemperatureUnit: input.TemperatureUnit,
	})
	if err != nil {
		return errResult(fmt.Sprintf("Historical weather request failed: %v", err)), nil, nil
	}

	tempUnit := "°C"
	if input.TemperatureUnit == "fahrenheit" {
		tempUnit = "°F"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"Historical Weather for (%.4f, %.4f) — %s to %s (TZ: %s)\n\n",
		hist.Latitude, hist.Longitude, input.StartDate, input.EndDate, hist.Timezone,
	))

	if d := hist.Daily; d != nil && len(d.Time) > 0 {
		for i, date := range d.Time {
			sb.WriteString(fmt.Sprintf(
				"  %s — %s\n"+
					"    High: %.1f%s  Low: %.1f%s\n"+
					"    Precipitation: %.1f mm | Max Wind: %.1f km/h\n\n",
				date, openmeteo.WeatherCodeDescription(d.WeatherCode[i]),
				d.Temperature2mMax[i], tempUnit, d.Temperature2mMin[i], tempUnit,
				d.PrecipitationSum[i], d.WindSpeed10mMax[i],
			))
		}
	} else {
		sb.WriteString("  No data available for the requested date range.\n")
	}

	return textResult(sb.String()), nil, nil
}

// --- Helpers ---

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

func errResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}
}

// --- Main ---

func main() {
	// Log to stderr so we don't interfere with STDIO JSON-RPC on stdout.
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Weather MCP Server v1.0.0")

	server := mcp.NewServer(
		&mcp.Implementation{Name: "weather-mcp-server", Version: "1.0.0"},
		nil,
	)

	// Register tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_location",
		Description: "Search for a city or place by name and get its geographic coordinates (latitude/longitude). Use this before get_forecast if you only have a city name.",
	}, handleSearchLocation)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_forecast",
		Description: "Get current weather conditions and a multi-day forecast for a location specified by latitude and longitude. Returns temperature, precipitation, wind, humidity, and weather conditions.",
	}, handleGetForecast)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_historical_weather",
		Description: "Get historical weather data for a location and date range (up to 1 year). Returns daily high/low temperatures, precipitation, wind speed, and conditions. Dates must be in YYYY-MM-DD format.",
	}, handleGetHistoricalWeather)

	log.Println("Tools registered: search_location, get_forecast, get_historical_weather")
	log.Println("Waiting for MCP client connection via STDIO...")

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
