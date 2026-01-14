package oddsapi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL      = "https://api.the-odds-api.com/v4"
	SportEPL     = "soccer_epl" // English Premier League
	RegionUK     = "uk"
	RegionEU     = "eu"
	RegionUS     = "us"
)

// Market types
const (
	MarketH2H    = "h2h"         // 1X2 (Home/Draw/Away)
	MarketTotals = "totals"      // Over/Under
	MarketBTTS   = "btts"        // Both Teams to Score
	MarketSpread = "spreads"     // Handicap
)

// Client represents The Odds API client
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Odds API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: BaseURL,
	}
}

// doRequest performs HTTP request with API key parameter
func (c *Client) doRequest(endpoint string, params map[string]string) ([]byte, error) {
	// Build URL
	reqURL, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add query parameters
	q := reqURL.Query()
	q.Add("apiKey", c.apiKey)
	for key, value := range params {
		q.Add(key, value)
	}
	reqURL.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Response structures

// Event represents a single event/match with odds
type Event struct {
	ID           string       `json:"id"`
	SportKey     string       `json:"sport_key"`
	SportTitle   string       `json:"sport_title"`
	CommenceTime time.Time    `json:"commence_time"`
	HomeTeam     string       `json:"home_team"`
	AwayTeam     string       `json:"away_team"`
	Bookmakers   []Bookmaker  `json:"bookmakers"`
}

// Bookmaker represents a bookmaker with their odds
type Bookmaker struct {
	Key        string    `json:"key"`
	Title      string    `json:"title"`
	LastUpdate time.Time `json:"last_update"`
	Markets    []Market  `json:"markets"`
}

// Market represents a specific betting market
type Market struct {
	Key        string    `json:"key"` // h2h, totals, btts, etc.
	LastUpdate time.Time `json:"last_update"`
	Outcomes   []Outcome `json:"outcomes"`
}

// Outcome represents a specific betting outcome
type Outcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"` // Decimal odds
	Point float64 `json:"point,omitempty"` // For totals/spreads (e.g., 2.5 goals)
}

// Sport represents available sport information
type Sport struct {
	Key          string `json:"key"`
	Group        string `json:"group"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Active       bool   `json:"active"`
	HasOutrights bool   `json:"has_outrights"`
}
