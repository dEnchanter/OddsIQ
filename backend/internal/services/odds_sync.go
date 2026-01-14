package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dEnchanter/OddsIQ/backend/internal/models"
	"github.com/dEnchanter/OddsIQ/backend/internal/repository"
	"github.com/dEnchanter/OddsIQ/backend/pkg/oddsapi"
)

// OddsSyncService handles syncing odds from The Odds API
type OddsSyncService struct {
	apiClient    *oddsapi.Client
	fixturesRepo *repository.FixturesRepository
	oddsRepo     *repository.OddsRepository
	teamsRepo    *repository.TeamsRepository
}

// NewOddsSyncService creates a new odds sync service
func NewOddsSyncService(
	apiClient *oddsapi.Client,
	fixturesRepo *repository.FixturesRepository,
	oddsRepo *repository.OddsRepository,
	teamsRepo *repository.TeamsRepository,
) *OddsSyncService {
	return &OddsSyncService{
		apiClient:    apiClient,
		fixturesRepo: fixturesRepo,
		oddsRepo:     oddsRepo,
		teamsRepo:    teamsRepo,
	}
}

// SyncAllMarkets syncs odds for all supported markets (1X2, Over/Under, BTTS)
func (s *OddsSyncService) SyncAllMarkets(ctx context.Context) error {
	log.Println("Syncing odds for all markets...")

	// Fetch events with all markets
	events, err := s.apiClient.GetAllMarketsEPL()
	if err != nil {
		return fmt.Errorf("failed to fetch odds: %w", err)
	}

	log.Printf("Fetched odds for %d events", len(events))

	// Process each event
	successCount := 0
	for _, event := range events {
		if err := s.processEvent(ctx, event); err != nil {
			log.Printf("Failed to process event %s: %v", event.ID, err)
			continue
		}
		successCount++
	}

	log.Printf("Successfully synced odds for %d/%d events", successCount, len(events))
	return nil
}

// SyncMarket syncs odds for a specific market type
func (s *OddsSyncService) SyncMarket(ctx context.Context, marketType string) error {
	log.Printf("Syncing odds for market: %s...", marketType)

	// Map market types
	markets := []string{marketType}

	// Fetch events
	events, err := s.apiClient.GetEPLOdds(markets)
	if err != nil {
		return fmt.Errorf("failed to fetch odds: %w", err)
	}

	log.Printf("Fetched odds for %d events", len(events))

	// Process each event
	successCount := 0
	for _, event := range events {
		if err := s.processEvent(ctx, event); err != nil {
			log.Printf("Failed to process event %s: %v", event.ID, err)
			continue
		}
		successCount++
	}

	log.Printf("Successfully synced odds for %d/%d events", successCount, len(events))
	return nil
}

// SyncH2HOdds syncs 1X2 (Home/Draw/Away) odds
func (s *OddsSyncService) SyncH2HOdds(ctx context.Context) error {
	return s.SyncMarket(ctx, oddsapi.MarketH2H)
}

// SyncTotalsOdds syncs Over/Under odds
func (s *OddsSyncService) SyncTotalsOdds(ctx context.Context) error {
	return s.SyncMarket(ctx, oddsapi.MarketTotals)
}

// SyncBTTSOdds syncs Both Teams to Score odds
func (s *OddsSyncService) SyncBTTSOdds(ctx context.Context) error {
	return s.SyncMarket(ctx, oddsapi.MarketBTTS)
}

// processEvent processes a single event and stores odds in database
func (s *OddsSyncService) processEvent(ctx context.Context, event oddsapi.Event) error {
	// Find matching fixture in database
	fixture, err := s.findMatchingFixture(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to find matching fixture: %w", err)
	}

	if fixture == nil {
		// No matching fixture found, skip
		log.Printf("No matching fixture found for event: %s vs %s", event.HomeTeam, event.AwayTeam)
		return nil
	}

	// Extract and store odds from all bookmakers and markets
	oddsList := s.extractOddsFromEvent(fixture.ID, event)

	// Batch insert odds
	if len(oddsList) > 0 {
		if err := s.oddsRepo.CreateBatch(ctx, oddsList); err != nil {
			return fmt.Errorf("failed to store odds: %w", err)
		}
		log.Printf("Stored %d odds entries for fixture %d", len(oddsList), fixture.ID)
	}

	return nil
}

// findMatchingFixture finds a fixture in database matching the odds API event
func (s *OddsSyncService) findMatchingFixture(ctx context.Context, event oddsapi.Event) (*models.Fixture, error) {
	// Get upcoming fixtures around the event commence time
	from := event.CommenceTime.Add(-12 * time.Hour)
	to := event.CommenceTime.Add(12 * time.Hour)

	fixtures, err := s.fixturesRepo.GetByDateRange(ctx, from, to)
	if err != nil {
		return nil, err
	}

	// Try to match by team names
	for _, fixture := range fixtures {
		// Get team names
		homeTeam, err := s.teamsRepo.GetByID(ctx, fixture.HomeTeamID)
		if err != nil {
			continue
		}

		awayTeam, err := s.teamsRepo.GetByID(ctx, fixture.AwayTeamID)
		if err != nil {
			continue
		}

		// Match by team names (case-insensitive, partial match)
		if s.matchTeamNames(homeTeam.Name, event.HomeTeam) && s.matchTeamNames(awayTeam.Name, event.AwayTeam) {
			return &fixture, nil
		}
	}

	return nil, nil
}

// matchTeamNames checks if two team names match (handles variations)
func (s *OddsSyncService) matchTeamNames(dbName, apiName string) bool {
	// Normalize names (lowercase, remove spaces)
	normalize := func(name string) string {
		return strings.ToLower(strings.ReplaceAll(name, " ", ""))
	}

	dbNorm := normalize(dbName)
	apiNorm := normalize(apiName)

	// Exact match
	if dbNorm == apiNorm {
		return true
	}

	// Partial match (one contains the other)
	if strings.Contains(dbNorm, apiNorm) || strings.Contains(apiNorm, dbNorm) {
		return true
	}

	// Handle common abbreviations
	abbreviations := map[string]string{
		"manchester":  "man",
		"united":      "utd",
		"city":        "",
		"tottenham":   "spurs",
		"newcastle":   "newcastle",
		"brighton":    "brighton",
		"westham":     "westham",
		"wolverhampton": "wolves",
		"nottingham":  "nottingham",
	}

	for full, abbr := range abbreviations {
		dbNorm = strings.ReplaceAll(dbNorm, full, abbr)
		apiNorm = strings.ReplaceAll(apiNorm, full, abbr)
	}

	return dbNorm == apiNorm
}

// extractOddsFromEvent extracts all odds from an event
func (s *OddsSyncService) extractOddsFromEvent(fixtureID int, event oddsapi.Event) []models.Odds {
	var oddsList []models.Odds
	timestamp := time.Now()

	for _, bookmaker := range event.Bookmakers {
		for _, market := range bookmaker.Markets {
			for _, outcome := range market.Outcomes {
				odds := models.Odds{
					FixtureID:  fixtureID,
					Bookmaker:  bookmaker.Key,
					MarketType: market.Key,
					Outcome:    s.normalizeOutcome(outcome.Name, market.Key),
					OddsValue:  outcome.Price,
					Timestamp:  timestamp,
				}
				oddsList = append(oddsList, odds)
			}
		}
	}

	return oddsList
}

// normalizeOutcome normalizes outcome names for consistency
func (s *OddsSyncService) normalizeOutcome(name, marketType string) string {
	// For h2h market, normalize to Home/Draw/Away
	if marketType == oddsapi.MarketH2H {
		// Names from API are team names or "Draw"
		if strings.ToLower(name) == "draw" {
			return "Draw"
		}
		// We'll keep team names as-is and normalize later when needed
		return name
	}

	// For totals market, normalize to Over/Under
	if marketType == oddsapi.MarketTotals {
		return name // Already "Over" or "Under"
	}

	// For BTTS market, normalize to Yes/No
	if marketType == oddsapi.MarketBTTS {
		if strings.ToLower(name) == "yes" {
			return "Yes"
		}
		return "No"
	}

	return name
}

// CleanupOldOdds removes odds older than specified days
func (s *OddsSyncService) CleanupOldOdds(ctx context.Context, daysToKeep int) error {
	log.Printf("Cleaning up odds older than %d days...", daysToKeep)

	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)
	deleted, err := s.oddsRepo.DeleteOldOdds(ctx, cutoffDate)
	if err != nil {
		return fmt.Errorf("failed to cleanup old odds: %w", err)
	}

	log.Printf("Deleted %d old odds entries", deleted)
	return nil
}

// GetOddsSummary returns a summary of stored odds
func (s *OddsSyncService) GetOddsSummary(ctx context.Context) (map[string]interface{}, error) {
	marketTypes, err := s.oddsRepo.GetMarketTypes(ctx)
	if err != nil {
		return nil, err
	}

	bookmakers, err := s.oddsRepo.GetBookmakers(ctx)
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"market_types": marketTypes,
		"bookmakers":   bookmakers,
		"total_markets": len(marketTypes),
		"total_bookmakers": len(bookmakers),
	}

	return summary, nil
}
