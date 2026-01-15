package apifootball

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// GetOddsByFixture fetches odds for a specific fixture
func (c *Client) GetOddsByFixture(fixtureID int) ([]OddsResponse, error) {
	params := map[string]string{
		"fixture": strconv.Itoa(fixtureID),
	}

	body, err := c.doRequest("/odds", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var odds []OddsResponse
	if err := json.Unmarshal(apiResp.Response, &odds); err != nil {
		return nil, fmt.Errorf("failed to parse odds: %w", err)
	}

	return odds, nil
}

// GetOddsByLeague fetches odds for all fixtures in a league and season
func (c *Client) GetOddsByLeague(leagueID, season int) ([]OddsResponse, error) {
	params := map[string]string{
		"league": strconv.Itoa(leagueID),
		"season": strconv.Itoa(season),
	}

	body, err := c.doRequest("/odds", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var odds []OddsResponse
	if err := json.Unmarshal(apiResp.Response, &odds); err != nil {
		return nil, fmt.Errorf("failed to parse odds: %w", err)
	}

	return odds, nil
}

// GetLiveOdds fetches live odds for a specific fixture
func (c *Client) GetLiveOdds(fixtureID int) ([]OddsResponse, error) {
	params := map[string]string{
		"fixture": strconv.Itoa(fixtureID),
	}

	body, err := c.doRequest("/odds/live", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var odds []OddsResponse
	if err := json.Unmarshal(apiResp.Response, &odds); err != nil {
		return nil, fmt.Errorf("failed to parse odds: %w", err)
	}

	return odds, nil
}

// GetBookmakers fetches the list of available bookmakers
func (c *Client) GetBookmakers() ([]BookmakerInfo, error) {
	params := map[string]string{}

	body, err := c.doRequest("/odds/bookmakers", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var bookmakers []BookmakerInfo
	if err := json.Unmarshal(apiResp.Response, &bookmakers); err != nil {
		return nil, fmt.Errorf("failed to parse bookmakers: %w", err)
	}

	return bookmakers, nil
}

// GetBetTypes fetches the list of available bet types
func (c *Client) GetBetTypes() ([]BetTypeInfo, error) {
	params := map[string]string{}

	body, err := c.doRequest("/odds/bets", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var betTypes []BetTypeInfo
	if err := json.Unmarshal(apiResp.Response, &betTypes); err != nil {
		return nil, fmt.Errorf("failed to parse bet types: %w", err)
	}

	return betTypes, nil
}

// OddsResponse represents the response structure for odds
type OddsResponse struct {
	League struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
		Logo    string `json:"logo"`
		Flag    string `json:"flag"`
		Season  int    `json:"season"`
	} `json:"league"`
	Fixture struct {
		ID       int    `json:"id"`
		Timezone string `json:"timezone"`
		Date     string `json:"date"`
		Timestamp int64 `json:"timestamp"`
	} `json:"fixture"`
	Update    string `json:"update"` // Last update timestamp
	Bookmakers []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Bets []struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Values []struct {
				Value string `json:"value"` // e.g., "Home", "Draw", "Away"
				Odd   string `json:"odd"`   // e.g., "1.85"
			} `json:"values"`
		} `json:"bets"`
	} `json:"bookmakers"`
}

// BookmakerInfo represents bookmaker information
type BookmakerInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// BetTypeInfo represents bet type information
type BetTypeInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
