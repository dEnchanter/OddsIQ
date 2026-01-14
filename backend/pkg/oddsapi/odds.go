package oddsapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetOdds fetches odds for a specific sport and markets
// markets: h2h, totals, btts, spreads (comma-separated)
// regions: uk, eu, us (comma-separated)
func (c *Client) GetOdds(sport string, markets []string, regions []string) ([]Event, error) {
	params := map[string]string{
		"markets": strings.Join(markets, ","),
		"regions": strings.Join(regions, ","),
	}

	endpoint := fmt.Sprintf("/sports/%s/odds", sport)
	body, err := c.doRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("failed to parse odds response: %w", err)
	}

	return events, nil
}

// GetEventOdds fetches odds for a specific event by ID
func (c *Client) GetEventOdds(sport, eventID string, markets []string, regions []string) (*Event, error) {
	params := map[string]string{
		"markets": strings.Join(markets, ","),
		"regions": strings.Join(regions, ","),
	}

	endpoint := fmt.Sprintf("/sports/%s/events/%s/odds", sport, eventID)
	body, err := c.doRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("failed to parse event odds: %w", err)
	}

	return &event, nil
}

// GetEPLOdds fetches odds for English Premier League matches
// This is a convenience method for the most common use case
func (c *Client) GetEPLOdds(markets []string) ([]Event, error) {
	regions := []string{RegionUK, RegionEU}
	return c.GetOdds(SportEPL, markets, regions)
}

// GetAllMarketsEPL fetches all available markets for EPL
func (c *Client) GetAllMarketsEPL() ([]Event, error) {
	markets := []string{MarketH2H, MarketTotals, MarketBTTS}
	return c.GetEPLOdds(markets)
}

// GetH2HOdds fetches 1X2 (Home/Draw/Away) odds for EPL
func (c *Client) GetH2HOdds() ([]Event, error) {
	return c.GetEPLOdds([]string{MarketH2H})
}

// GetTotalsOdds fetches Over/Under odds for EPL
func (c *Client) GetTotalsOdds() ([]Event, error) {
	return c.GetEPLOdds([]string{MarketTotals})
}

// GetBTTSOdds fetches Both Teams to Score odds for EPL
func (c *Client) GetBTTSOdds() ([]Event, error) {
	return c.GetEPLOdds([]string{MarketBTTS})
}

// GetSports fetches list of available sports
func (c *Client) GetSports() ([]Sport, error) {
	body, err := c.doRequest("/sports", nil)
	if err != nil {
		return nil, err
	}

	var sports []Sport
	if err := json.Unmarshal(body, &sports); err != nil {
		return nil, fmt.Errorf("failed to parse sports response: %w", err)
	}

	return sports, nil
}

// OddsHelper provides utility functions for working with odds data

// GetBestOdds finds the best odds for a specific outcome across all bookmakers
func GetBestOdds(event Event, marketKey, outcomeName string) *Outcome {
	var bestOdds *Outcome

	for _, bookmaker := range event.Bookmakers {
		for _, market := range bookmaker.Markets {
			if market.Key != marketKey {
				continue
			}

			for _, outcome := range market.Outcomes {
				if outcome.Name != outcomeName {
					continue
				}

				if bestOdds == nil || outcome.Price > bestOdds.Price {
					bestOdds = &outcome
				}
			}
		}
	}

	return bestOdds
}

// GetAverageOdds calculates average odds for a specific outcome across all bookmakers
func GetAverageOdds(event Event, marketKey, outcomeName string) float64 {
	var sum float64
	var count int

	for _, bookmaker := range event.Bookmakers {
		for _, market := range bookmaker.Markets {
			if market.Key != marketKey {
				continue
			}

			for _, outcome := range market.Outcomes {
				if outcome.Name != outcomeName {
					continue
				}

				sum += outcome.Price
				count++
			}
		}
	}

	if count == 0 {
		return 0
	}

	return sum / float64(count)
}

// GetBookmakerOdds gets odds from a specific bookmaker
func GetBookmakerOdds(event Event, bookmakerKey, marketKey, outcomeName string) *Outcome {
	for _, bookmaker := range event.Bookmakers {
		if bookmaker.Key != bookmakerKey {
			continue
		}

		for _, market := range bookmaker.Markets {
			if market.Key != marketKey {
				continue
			}

			for _, outcome := range market.Outcomes {
				if outcome.Name == outcomeName {
					return &outcome
				}
			}
		}
	}

	return nil
}

// ExtractH2HOdds extracts 1X2 odds from an event
func ExtractH2HOdds(event Event) (home, draw, away float64, found bool) {
	bestHome := GetBestOdds(event, MarketH2H, event.HomeTeam)
	bestDraw := GetBestOdds(event, MarketH2H, "Draw")
	bestAway := GetBestOdds(event, MarketH2H, event.AwayTeam)

	if bestHome != nil && bestDraw != nil && bestAway != nil {
		return bestHome.Price, bestDraw.Price, bestAway.Price, true
	}

	return 0, 0, 0, false
}

// ExtractOverUnderOdds extracts Over/Under 2.5 odds from an event
func ExtractOverUnderOdds(event Event, point float64) (over, under float64, found bool) {
	for _, bookmaker := range event.Bookmakers {
		for _, market := range bookmaker.Markets {
			if market.Key != MarketTotals {
				continue
			}

			var overOdds, underOdds *Outcome
			for _, outcome := range market.Outcomes {
				if outcome.Point != point {
					continue
				}

				if outcome.Name == "Over" {
					overOdds = &outcome
				} else if outcome.Name == "Under" {
					underOdds = &outcome
				}
			}

			if overOdds != nil && underOdds != nil {
				if over == 0 || overOdds.Price > over {
					over = overOdds.Price
				}
				if under == 0 || underOdds.Price > under {
					under = underOdds.Price
				}
			}
		}
	}

	if over > 0 && under > 0 {
		return over, under, true
	}

	return 0, 0, false
}
