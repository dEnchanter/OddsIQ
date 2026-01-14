package apifootball

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL        = "https://v3.football.api-sports.io"
	PremierLeagueID = 39 // England Premier League
)

// Client represents API-Football client
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new API-Football client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: BaseURL,
	}
}

// doRequest performs HTTP request with API key header
func (c *Client) doRequest(endpoint string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key header
	req.Header.Add("x-apisports-key", c.apiKey)

	// Add query parameters
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

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

// APIResponse is the generic response wrapper
type APIResponse struct {
	Get        string                   `json:"get"`
	Parameters map[string]interface{}   `json:"parameters"`
	Errors     map[string]interface{}   `json:"errors"`
	Results    int                      `json:"results"`
	Paging     map[string]interface{}   `json:"paging"`
	Response   json.RawMessage          `json:"response"`
}

// Team represents a team
type Team struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Country  string `json:"country"`
	Founded  int    `json:"founded"`
	National bool   `json:"national"`
	Logo     string `json:"logo"`
}

// Venue represents a stadium/venue
type Venue struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Capacity int    `json:"capacity"`
	Surface  string `json:"surface"`
	Image    string `json:"image"`
}

// Fixture represents a match fixture
type Fixture struct {
	ID        int       `json:"id"`
	Referee   string    `json:"referee"`
	Timezone  string    `json:"timezone"`
	Date      time.Time `json:"date"`
	Timestamp int64     `json:"timestamp"`
	Venue     Venue     `json:"venue"`
	Status    struct {
		Long    string `json:"long"`
		Short   string `json:"short"`
		Elapsed int    `json:"elapsed"`
	} `json:"status"`
}

// FixtureResponse represents fixture data from API
type FixtureResponse struct {
	Fixture Fixture `json:"fixture"`
	League  struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
		Logo    string `json:"logo"`
		Flag    string `json:"flag"`
		Season  int    `json:"season"`
		Round   string `json:"round"`
	} `json:"league"`
	Teams struct {
		Home Team `json:"home"`
		Away Team `json:"away"`
	} `json:"teams"`
	Goals struct {
		Home int `json:"home"`
		Away int `json:"away"`
	} `json:"goals"`
	Score struct {
		Halftime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"halftime"`
		Fulltime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"fulltime"`
	} `json:"score"`
}

// Standings represents league standings
type StandingsResponse struct {
	League struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Country   string `json:"country"`
		Season    int    `json:"season"`
		Standings [][]struct {
			Rank   int  `json:"rank"`
			Team   Team `json:"team"`
			Points int  `json:"points"`
			All    struct {
				Played int `json:"played"`
				Win    int `json:"win"`
				Draw   int `json:"draw"`
				Lose   int `json:"lose"`
				Goals  struct {
					For     int `json:"for"`
					Against int `json:"against"`
				} `json:"goals"`
			} `json:"all"`
			Home struct {
				Played int `json:"played"`
				Win    int `json:"win"`
				Draw   int `json:"draw"`
				Lose   int `json:"lose"`
			} `json:"home"`
			Away struct {
				Played int `json:"played"`
				Win    int `json:"win"`
				Draw   int `json:"draw"`
				Lose   int `json:"lose"`
			} `json:"away"`
			Form        string `json:"form"`
			Description string `json:"description"`
			GoalsDiff   int    `json:"goalsDiff"`
		} `json:"standings"`
	} `json:"league"`
}
