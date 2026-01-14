package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dEnchanter/OddsIQ/backend/internal/models"
	"github.com/dEnchanter/OddsIQ/backend/internal/repository"
	"github.com/dEnchanter/OddsIQ/backend/pkg/apifootball"
)

// FixtureSyncService handles syncing fixtures from API-Football
type FixtureSyncService struct {
	apiClient   *apifootball.Client
	teamsRepo   *repository.TeamsRepository
	fixturesRepo *repository.FixturesRepository
}

// NewFixtureSyncService creates a new fixture sync service
func NewFixtureSyncService(
	apiClient *apifootball.Client,
	teamsRepo *repository.TeamsRepository,
	fixturesRepo *repository.FixturesRepository,
) *FixtureSyncService {
	return &FixtureSyncService{
		apiClient:   apiClient,
		teamsRepo:   teamsRepo,
		fixturesRepo: fixturesRepo,
	}
}

// SyncTeams fetches and stores Premier League teams
func (s *FixtureSyncService) SyncTeams(ctx context.Context, season int) error {
	log.Printf("Syncing teams for season %d...", season)

	// Fetch teams from API
	teamsResp, err := s.apiClient.GetTeams(apifootball.PremierLeagueID, season)
	if err != nil {
		return fmt.Errorf("failed to fetch teams: %w", err)
	}

	log.Printf("Fetched %d teams from API", len(teamsResp))

	// Upsert each team
	for _, teamResp := range teamsResp {
		team := &models.Team{
			APIFootballID: teamResp.Team.ID,
			Name:          teamResp.Team.Name,
			Code:          teamResp.Team.Code,
			LogoURL:       teamResp.Team.Logo,
			Founded:       teamResp.Team.Founded,
			VenueName:     teamResp.Venue.Name,
			VenueCity:     teamResp.Venue.City,
			VenueCapacity: teamResp.Venue.Capacity,
		}

		if err := s.teamsRepo.Upsert(ctx, team); err != nil {
			log.Printf("Failed to upsert team %s: %v", team.Name, err)
			continue
		}

		log.Printf("Upserted team: %s (ID: %d)", team.Name, team.ID)
	}

	log.Printf("Successfully synced %d teams", len(teamsResp))
	return nil
}

// SyncFixturesBySeason fetches and stores all fixtures for a season
func (s *FixtureSyncService) SyncFixturesBySeason(ctx context.Context, season int) error {
	log.Printf("Syncing fixtures for season %d...", season)

	// Fetch fixtures from API
	fixturesResp, err := s.apiClient.GetFixtures(apifootball.PremierLeagueID, season)
	if err != nil {
		return fmt.Errorf("failed to fetch fixtures: %w", err)
	}

	log.Printf("Fetched %d fixtures from API", len(fixturesResp))

	// Process each fixture
	successCount := 0
	for _, fixtureResp := range fixturesResp {
		if err := s.processFixture(ctx, fixtureResp, season); err != nil {
			log.Printf("Failed to process fixture %d: %v", fixtureResp.Fixture.ID, err)
			continue
		}
		successCount++
	}

	log.Printf("Successfully synced %d/%d fixtures", successCount, len(fixturesResp))
	return nil
}

// SyncFixturesByDateRange fetches and stores fixtures within a date range
func (s *FixtureSyncService) SyncFixturesByDateRange(ctx context.Context, from, to time.Time) error {
	fromStr := from.Format("2006-01-02")
	toStr := to.Format("2006-01-02")

	log.Printf("Syncing fixtures from %s to %s...", fromStr, toStr)

	// Fetch fixtures from API
	fixturesResp, err := s.apiClient.GetFixturesByDateRange(fromStr, toStr)
	if err != nil {
		return fmt.Errorf("failed to fetch fixtures: %w", err)
	}

	log.Printf("Fetched %d fixtures from API", len(fixturesResp))

	// Process each fixture
	successCount := 0
	for _, fixtureResp := range fixturesResp {
		// Extract season from fixture date
		season := fixtureResp.League.Season

		if err := s.processFixture(ctx, fixtureResp, season); err != nil {
			log.Printf("Failed to process fixture %d: %v", fixtureResp.Fixture.ID, err)
			continue
		}
		successCount++
	}

	log.Printf("Successfully synced %d/%d fixtures", successCount, len(fixturesResp))
	return nil
}

// SyncUpcomingFixtures syncs upcoming fixtures (next 7 days)
func (s *FixtureSyncService) SyncUpcomingFixtures(ctx context.Context) error {
	now := time.Now()
	to := now.AddDate(0, 0, 7) // Next 7 days

	return s.SyncFixturesByDateRange(ctx, now, to)
}

// UpdateFixtureResults updates scores and status for recently completed fixtures
func (s *FixtureSyncService) UpdateFixtureResults(ctx context.Context) error {
	log.Println("Updating fixture results...")

	// Get fixtures from last 2 days that might have been completed
	from := time.Now().AddDate(0, 0, -2)
	to := time.Now()

	// Fetch latest fixture data
	fromStr := from.Format("2006-01-02")
	toStr := to.Format("2006-01-02")

	fixturesResp, err := s.apiClient.GetFixturesByDateRange(fromStr, toStr)
	if err != nil {
		return fmt.Errorf("failed to fetch fixtures: %w", err)
	}

	log.Printf("Checking %d fixtures for result updates", len(fixturesResp))

	// Update each fixture
	successCount := 0
	for _, fixtureResp := range fixturesResp {
		season := fixtureResp.League.Season

		if err := s.processFixture(ctx, fixtureResp, season); err != nil {
			log.Printf("Failed to update fixture %d: %v", fixtureResp.Fixture.ID, err)
			continue
		}
		successCount++
	}

	log.Printf("Successfully updated %d/%d fixtures", successCount, len(fixturesResp))
	return nil
}

// processFixture converts API fixture to model and upserts to database
func (s *FixtureSyncService) processFixture(ctx context.Context, fixtureResp apifootball.FixtureResponse, season int) error {
	// Get team IDs from database using API-Football IDs
	homeTeam, err := s.teamsRepo.GetByAPIFootballID(ctx, fixtureResp.Teams.Home.ID)
	if err != nil {
		return fmt.Errorf("home team not found: %w", err)
	}

	awayTeam, err := s.teamsRepo.GetByAPIFootballID(ctx, fixtureResp.Teams.Away.ID)
	if err != nil {
		return fmt.Errorf("away team not found: %w", err)
	}

	// Extract scores (may be nil if match hasn't started)
	var homeScore, awayScore *int
	if fixtureResp.Goals.Home >= 0 {
		homeScore = &fixtureResp.Goals.Home
	}
	if fixtureResp.Goals.Away >= 0 {
		awayScore = &fixtureResp.Goals.Away
	}

	// Create fixture model
	fixture := &models.Fixture{
		APIFootballID: fixtureResp.Fixture.ID,
		Season:        season,
		MatchDate:     fixtureResp.Fixture.Date,
		Round:         fixtureResp.League.Round,
		HomeTeamID:    homeTeam.ID,
		AwayTeamID:    awayTeam.ID,
		Status:        fixtureResp.Fixture.Status.Short,
		HomeScore:     homeScore,
		AwayScore:     awayScore,
		VenueName:     fixtureResp.Fixture.Venue.Name,
		Referee:       fixtureResp.Fixture.Referee,
	}

	// Upsert fixture
	if err := s.fixturesRepo.Upsert(ctx, fixture); err != nil {
		return fmt.Errorf("failed to upsert fixture: %w", err)
	}

	return nil
}

// SyncAllSeasons syncs teams and fixtures for multiple seasons
func (s *FixtureSyncService) SyncAllSeasons(ctx context.Context, seasons []int) error {
	for _, season := range seasons {
		log.Printf("=== Syncing season %d ===", season)

		// First sync teams
		if err := s.SyncTeams(ctx, season); err != nil {
			log.Printf("Failed to sync teams for season %d: %v", season, err)
			continue
		}

		// Then sync fixtures
		if err := s.SyncFixturesBySeason(ctx, season); err != nil {
			log.Printf("Failed to sync fixtures for season %d: %v", season, err)
			continue
		}

		log.Printf("=== Completed season %d ===", season)

		// Small delay to respect API rate limits
		time.Sleep(2 * time.Second)
	}

	return nil
}
