package apifootball

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// TeamResponse represents team data from API
type TeamResponse struct {
	Team  Team  `json:"team"`
	Venue Venue `json:"venue"`
}

// GetTeams fetches all teams for a specific league and season
func (c *Client) GetTeams(leagueID, season int) ([]TeamResponse, error) {
	params := map[string]string{
		"league": strconv.Itoa(leagueID),
		"season": strconv.Itoa(season),
	}

	body, err := c.doRequest("/teams", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var teams []TeamResponse
	if err := json.Unmarshal(apiResp.Response, &teams); err != nil {
		return nil, fmt.Errorf("failed to parse teams: %w", err)
	}

	return teams, nil
}

// GetTeam fetches a single team by ID
func (c *Client) GetTeam(teamID int) (*TeamResponse, error) {
	params := map[string]string{
		"id": strconv.Itoa(teamID),
	}

	body, err := c.doRequest("/teams", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var teams []TeamResponse
	if err := json.Unmarshal(apiResp.Response, &teams); err != nil {
		return nil, fmt.Errorf("failed to parse teams: %w", err)
	}

	if len(teams) == 0 {
		return nil, fmt.Errorf("team not found")
	}

	return &teams[0], nil
}
