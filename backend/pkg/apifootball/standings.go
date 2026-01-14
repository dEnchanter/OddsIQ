package apifootball

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// GetStandings fetches league standings for a specific season
func (c *Client) GetStandings(leagueID, season int) (*StandingsResponse, error) {
	params := map[string]string{
		"league": strconv.Itoa(leagueID),
		"season": strconv.Itoa(season),
	}

	body, err := c.doRequest("/standings", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var standings []StandingsResponse
	if err := json.Unmarshal(apiResp.Response, &standings); err != nil {
		return nil, fmt.Errorf("failed to parse standings: %w", err)
	}

	if len(standings) == 0 {
		return nil, fmt.Errorf("standings not found")
	}

	return &standings[0], nil
}
