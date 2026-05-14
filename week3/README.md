# Week 3 — Weather MCP Server

A **Model Context Protocol (MCP)** server written in Go that wraps the [Open-Meteo](https://open-meteo.com/) weather API. It runs locally over **STDIO transport** and integrates with Claude Desktop (or any MCP-compatible client).

## API Endpoints Used

This server wraps three Open-Meteo endpoints (all free, no API key required):

| Endpoint | Base URL | Purpose |
|---|---|---|
| Geocoding | `geocoding-api.open-meteo.com/v1/search` | Resolve city/place names → coordinates |
| Forecast | `api.open-meteo.com/v1/forecast` | Current conditions + daily forecast (up to 16 days) |
| Archive | `archive-api.open-meteo.com/v1/archive` | Historical weather data (back to 1940) |

## Prerequisites

- **Go 1.23+** — [install Go](https://go.dev/dl/)
- **Claude Desktop** (or another MCP client) — [download](https://claude.ai/download)
- No API keys needed — Open-Meteo is free for non-commercial use (≤10,000 requests/day)

## Setup & Build

```bash
# Clone and navigate to the server directory
cd week3/server

# Download dependencies
go mod tidy

# Build the binary
go build -o weather-mcp-server .
```

This produces a `weather-mcp-server` binary in the `server/` directory.

## Configure Claude Desktop

Add the following to your Claude Desktop MCP configuration file:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

### Option A: Run from built binary

```json
{
  "mcpServers": {
    "weather": {
      "command": "/absolute/path/to/week3/server/weather-mcp-server"
    }
  }
}
```

### Option B: Run with `go run` (no build step)

```json
{
  "mcpServers": {
    "weather": {
      "command": "go",
      "args": ["run", "."],
      "cwd": "/absolute/path/to/week3/server"
    }
  }
}
```

After saving, **restart Claude Desktop**. You should see a 🔨 tools icon indicating the server is connected.

## Testing with MCP Inspector

You can interactively test the server tools using the official MCP Inspector. This provides a web UI to view and execute the tools.

Run the compiled binary with the inspector:

```bash
npx @modelcontextprotocol/inspector /absolute/path/to/week3/server/weather-mcp-server
```

Or, run it directly from the source code (make sure to run this command from inside the `week3/server` directory so Go finds the `go.mod` file):

```bash
cd week3/server
npx @modelcontextprotocol/inspector go run .
```

After running the command, open the provided local URL (usually `http://localhost:5173`) in your browser to test the tools.

## Tool Reference

### 1. `search_location`

Search for a city or place by name and get its geographic coordinates.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `name` | string | ✅ | City or place name to search for |
| `count` | int | ❌ | Number of results (1–10, default 5) |

**Example input:**
```json
{ "name": "Oslo", "count": 3 }
```

**Example output:**
```
Found 3 location(s) for "Oslo":

1. Oslo, Norway (Oslo)
   Lat: 59.9133, Lon: 10.7389 | Elevation: 12m | Pop: 580000 | TZ: Europe/Oslo
2. Oslo, United States (Minnesota)
   Lat: 48.1947, Lon: -97.1262 | Elevation: 253m | Pop: 347 | TZ: America/Chicago
3. ...
```

### 2. `get_forecast`

Get current weather conditions and a multi-day daily forecast.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `latitude` | float | ✅ | Latitude (-90 to 90) |
| `longitude` | float | ✅ | Longitude (-180 to 180) |
| `days` | int | ❌ | Forecast days (1–16, default 3) |
| `temperature_unit` | string | ❌ | `celsius` or `fahrenheit` (default `celsius`) |

**Example input:**
```json
{ "latitude": 59.91, "longitude": 10.75, "days": 3 }
```

**Example output:**
```
Weather for (59.9100, 10.7500) — TZ: Europe/Oslo

── Current Conditions (2025-05-13T10:00) ──
  Temperature: 14.2°C
  Humidity:    62%
  Wind:        12.3 km/h (dir: 220°)
  Conditions:  Partly cloudy

── Daily Forecast ──

  2025-05-13 — Partly cloudy
    High: 17.5°C  Low: 8.2°C
    Precipitation: 0.0 mm | Max Wind: 18.4 km/h
    Sunrise: 04:52  Sunset: 21:32
  ...
```

### 3. `get_historical_weather`

Get historical/past weather data for a location and date range.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `latitude` | float | ✅ | Latitude (-90 to 90) |
| `longitude` | float | ✅ | Longitude (-180 to 180) |
| `start_date` | string | ✅ | Start date (YYYY-MM-DD) |
| `end_date` | string | ✅ | End date (YYYY-MM-DD) |
| `temperature_unit` | string | ❌ | `celsius` or `fahrenheit` (default `celsius`) |

**Example input:**
```json
{
  "latitude": 59.91,
  "longitude": 10.75,
  "start_date": "2024-12-24",
  "end_date": "2024-12-26"
}
```

**Example output:**
```
Historical Weather for (59.9100, 10.7500) — 2024-12-24 to 2024-12-26 (TZ: Europe/Oslo)

  2024-12-24 — Slight rain
    High: 5.2°C  Low: 1.8°C
    Precipitation: 4.3 mm | Max Wind: 22.1 km/h

  2024-12-25 — Overcast
    High: 3.1°C  Low: -0.4°C
    Precipitation: 0.2 mm | Max Wind: 15.6 km/h
  ...
```

## Example Invocation Flow

1. Open Claude Desktop with the server configured.
2. Ask: **"What's the weather like in Berlin right now?"**
   - Claude calls `search_location` with `{"name": "Berlin"}` to get coordinates.
   - Claude calls `get_forecast` with the Berlin coordinates.
   - Claude presents the current conditions and forecast in natural language.
3. Ask: **"How was the weather in Tokyo last Christmas?"**
   - Claude calls `search_location` with `{"name": "Tokyo"}`.
   - Claude calls `get_historical_weather` with `{"start_date": "2024-12-24", "end_date": "2024-12-26", ...}`.

## Resilience & Error Handling

- **Rate limiting**: Built-in throttle (150ms between requests) to stay well within Open-Meteo's 10,000 req/day free tier. HTTP 429 responses are caught and reported as user-friendly errors.
- **Timeouts**: All HTTP requests have a 10-second timeout.
- **Input validation**: Coordinates, date formats, and date ranges are validated before making API calls.
- **Graceful errors**: All API errors (network failures, bad responses, empty results) are returned as MCP error results — the server never crashes.
- **Logging**: All logs go to **stderr** (never stdout) to avoid interfering with the STDIO JSON-RPC transport.

## Project Structure

```
week3/
├── README.md                  # This file
├── assignment.md              # Assignment description
└── server/
    ├── main.go                # MCP server entry point + tool handlers
    ├── go.mod                 # Go module definition
    ├── go.sum                 # Dependency checksums
    ├── .gitignore             # Ignores compiled binary
    └── openmeteo/
        └── client.go          # Open-Meteo API client (geocoding, forecast, historical)
```

## Technology Choices

- **Language**: Go — strong typing, excellent HTTP stdlib, compiles to a single static binary.
- **MCP SDK**: [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) v1.6.0 — the official Go SDK maintained in collaboration with Google.
- **Transport**: STDIO — simplest local deployment, natively supported by Claude Desktop.
- **API**: [Open-Meteo](https://open-meteo.com/) — free, no API key required, global coverage, rich data.
