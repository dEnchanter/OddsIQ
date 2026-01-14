package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dEnchanter/OddsIQ/backend/config"
	"github.com/dEnchanter/OddsIQ/backend/internal/repository"
	"github.com/dEnchanter/OddsIQ/backend/internal/services"
	"github.com/dEnchanter/OddsIQ/backend/pkg/apifootball"
	"github.com/dEnchanter/OddsIQ/backend/pkg/database"
)

func main() {
	// Command-line flags
	seasonsFlag := flag.String("seasons", "2022,2023,2024", "Comma-separated list of seasons to backfill")
	teamsOnly := flag.Bool("teams-only", false, "Only sync teams, skip fixtures")
	fixturesOnly := flag.Bool("fixtures-only", false, "Only sync fixtures, skip teams")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *help {
		printHelp()
		return
	}

	// Parse seasons
	seasons, err := parseSeasons(*seasonsFlag)
	if err != nil {
		log.Fatalf("Invalid seasons format: %v", err)
	}

	log.Printf("Starting backfill for seasons: %v", seasons)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	// Initialize API clients
	apiFootballClient := apifootball.NewClient(cfg.APIFootballKey)

	// Initialize repositories
	teamsRepo := repository.NewTeamsRepository(db)
	fixturesRepo := repository.NewFixturesRepository(db)

	// Initialize sync service
	fixtureSyncService := services.NewFixtureSyncService(
		apiFootballClient,
		teamsRepo,
		fixturesRepo,
	)

	// Create context
	ctx := context.Background()

	// Execute backfill
	for _, season := range seasons {
		log.Printf("\n=== Processing Season %d ===\n", season)

		// Sync teams (unless fixtures-only)
		if !*fixturesOnly {
			log.Printf("Syncing teams for season %d...", season)
			if err := fixtureSyncService.SyncTeams(ctx, season); err != nil {
				log.Printf("ERROR: Failed to sync teams: %v", err)
				continue
			}
			log.Println("✓ Teams synced successfully")
		}

		// Sync fixtures (unless teams-only)
		if !*teamsOnly {
			log.Printf("Syncing fixtures for season %d...", season)
			if err := fixtureSyncService.SyncFixturesBySeason(ctx, season); err != nil {
				log.Printf("ERROR: Failed to sync fixtures: %v", err)
				continue
			}
			log.Println("✓ Fixtures synced successfully")
		}

		log.Printf("=== Completed Season %d ===\n", season)
	}

	// Print summary
	printSummary(ctx, teamsRepo, fixturesRepo)

	log.Println("\n✓ Backfill completed successfully")
}

func parseSeasons(seasonsStr string) ([]int, error) {
	parts := strings.Split(seasonsStr, ",")
	seasons := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		season, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid season: %s", part)
		}
		seasons = append(seasons, season)
	}

	return seasons, nil
}

func printSummary(ctx context.Context, teamsRepo *repository.TeamsRepository, fixturesRepo *repository.FixturesRepository) {
	log.Println("\n=== Backfill Summary ===")

	// Count teams
	teams, err := teamsRepo.GetAll(ctx)
	if err != nil {
		log.Printf("Failed to count teams: %v", err)
	} else {
		log.Printf("Total teams in database: %d", len(teams))
	}

	// Count fixtures by season
	seasons := []int{2022, 2023, 2024}
	for _, season := range seasons {
		fixtures, err := fixturesRepo.GetBySeason(ctx, season)
		if err != nil {
			log.Printf("Failed to count fixtures for season %d: %v", season, err)
		} else {
			log.Printf("Season %d fixtures: %d", season, len(fixtures))
		}
	}

	log.Println("======================")
}

func printHelp() {
	fmt.Println("OddsIQ Historical Data Backfill Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/backfill/main.go [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -seasons string")
	fmt.Println("        Comma-separated list of seasons to backfill (default \"2022,2023,2024\")")
	fmt.Println("  -teams-only")
	fmt.Println("        Only sync teams, skip fixtures")
	fmt.Println("  -fixtures-only")
	fmt.Println("        Only sync fixtures, skip teams")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Backfill all data for 2022-2024")
	fmt.Println("  go run cmd/backfill/main.go")
	fmt.Println()
	fmt.Println("  # Backfill only 2024 season")
	fmt.Println("  go run cmd/backfill/main.go -seasons 2024")
	fmt.Println()
	fmt.Println("  # Backfill only teams")
	fmt.Println("  go run cmd/backfill/main.go -teams-only")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  DATABASE_URL         PostgreSQL connection string")
	fmt.Println("  API_FOOTBALL_KEY     API-Football API key")
	fmt.Println()
}
