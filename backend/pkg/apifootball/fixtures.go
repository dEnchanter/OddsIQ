package apifootball

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// GetFixtures fetches fixtures for a specific league and season
func (c *Client) GetFixtures(leagueID, season int) ([]FixtureResponse, error) {
	params := map[string]string{
		"league": strconv.Itoa(leagueID),
		"season": strconv.Itoa(season),
	}

	body, err := c.doRequest("/fixtures", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var fixtures []FixtureResponse
	if err := json.Unmarshal(apiResp.Response, &fixtures); err != nil {
		return nil, fmt.Errorf("failed to parse fixtures: %w", err)
	}

	return fixtures, nil
}

// GetFixturesByDate fetches fixtures for a specific date
func (c *Client) GetFixturesByDate(date string) ([]FixtureResponse, error) {
	params := map[string]string{
		"date":   date, // Format: YYYY-MM-DD
		"league": strconv.Itoa(PremierLeagueID),
	}

	body, err := c.doRequest("/fixtures", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var fixtures []FixtureResponse
	if err := json.Unmarshal(apiResp.Response, &fixtures); err != nil {
		return nil, fmt.Errorf("failed to parse fixtures: %w", err)
	}

	return fixtures, nil
}

// GetFixturesByDateRange fetches fixtures between two dates
func (c *Client) GetFixturesByDateRange(from, to string) ([]FixtureResponse, error) {
	params := map[string]string{
		"from":   from, // Format: YYYY-MM-DD
		"to":     to,   // Format: YYYY-MM-DD
		"league": strconv.Itoa(PremierLeagueID),
	}

	body, err := c.doRequest("/fixtures", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var fixtures []FixtureResponse
	if err := json.Unmarshal(apiResp.Response, &fixtures); err != nil {
		return nil, fmt.Errorf("failed to parse fixtures: %w", err)
	}

	return fixtures, nil
}

// GetFixture fetches a single fixture by ID
func (c *Client) GetFixture(fixtureID int) (*FixtureResponse, error) {
	params := map[string]string{
		"id": strconv.Itoa(fixtureID),
	}

	body, err := c.doRequest("/fixtures", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var fixtures []FixtureResponse
	if err := json.Unmarshal(apiResp.Response, &fixtures); err != nil {
		return nil, fmt.Errorf("failed to parse fixtures: %w", err)
	}

	if len(fixtures) == 0 {
		return nil, fmt.Errorf("fixture not found")
	}

	return &fixtures[0], nil
}
