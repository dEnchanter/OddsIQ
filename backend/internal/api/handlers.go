package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/dEnchanter/OddsIQ/backend/config"
	"github.com/dEnchanter/OddsIQ/backend/internal/models"
	"github.com/dEnchanter/OddsIQ/backend/internal/repository"
	"github.com/dEnchanter/OddsIQ/backend/internal/services"
)

// ManualFixtureRequest represents a request to create a fixture manually
type ManualFixtureRequest struct {
	HomeTeamID int    `json:"home_team_id" binding:"required"`
	AwayTeamID int    `json:"away_team_id" binding:"required"`
	MatchDate  string `json:"match_date" binding:"required"` // Format: "2025-01-20T15:00:00Z"
	Season     int    `json:"season" binding:"required"`
	Round      string `json:"round"`
	VenueName  string `json:"venue_name"`
}

// ManualOddsRequest represents a request to add odds manually
type ManualOddsRequest struct {
	FixtureID  int     `json:"fixture_id" binding:"required"`
	Bookmaker  string  `json:"bookmaker" binding:"required"`
	MarketType string  `json:"market_type" binding:"required"` // h2h, totals, btts
	Outcome    string  `json:"outcome" binding:"required"`     // Home, Draw, Away, Over, Under, Yes, No
	OddsValue  float64 `json:"odds_value" binding:"required"`
}

// ManualOddsBatchRequest represents a request to add multiple odds at once
type ManualOddsBatchRequest struct {
	FixtureID int              `json:"fixture_id" binding:"required"`
	Bookmaker string           `json:"bookmaker" binding:"required"`
	Odds      []OddsEntryInput `json:"odds" binding:"required"`
}

// OddsEntryInput represents a single odds entry
type OddsEntryInput struct {
	MarketType string  `json:"market_type" binding:"required"`
	Outcome    string  `json:"outcome" binding:"required"`
	OddsValue  float64 `json:"odds_value" binding:"required"`
}

// API holds all the dependencies for handlers
type API struct {
	db                  *pgxpool.Pool
	cfg                 *config.Config
	teamsRepo           *repository.TeamsRepository
	fixturesRepo        *repository.FixturesRepository
	oddsRepo            *repository.OddsRepository
	statsRepo           *repository.TeamStatsRepository
	predictionService   *services.PredictionService
	bettingService      *services.BettingService
	accumulatorService  *services.AccumulatorService
}

// NewAPI creates a new API instance
func NewAPI(db *pgxpool.Pool, cfg *config.Config) *API {
	fixturesRepo := repository.NewFixturesRepository(db)
	oddsRepo := repository.NewOddsRepository(db)
	mlClient := services.NewMLClient(cfg.MLServiceURL)
	bettingService := services.NewBettingService(cfg, mlClient, fixturesRepo, oddsRepo)

	return &API{
		db:                  db,
		cfg:                 cfg,
		teamsRepo:           repository.NewTeamsRepository(db),
		fixturesRepo:        fixturesRepo,
		oddsRepo:            oddsRepo,
		statsRepo:           repository.NewTeamStatsRepository(db),
		predictionService:   services.NewPredictionService(cfg, fixturesRepo, oddsRepo),
		bettingService:      bettingService,
		accumulatorService:  services.NewAccumulatorService(bettingService, cfg),
	}
}

// healthCheck returns a health check handler
func (api *API) healthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Check database connection
		if err := api.db.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "oddsiq-backend",
			"version": "0.1.0",
		})
	}
}

// getFixtures returns fixtures list handler
func (api *API) getFixtures() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get query parameters
		seasonStr := c.Query("season")
		status := c.Query("status")

		var fixtures []interface{}

		if seasonStr != "" {
			season, parseErr := strconv.Atoi(seasonStr)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid season parameter"})
				return
			}
			fixturesList, err := api.fixturesRepo.GetBySeason(ctx, season)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			for _, f := range fixturesList {
				fixtures = append(fixtures, f)
			}
		} else if status != "" {
			fixturesList, err := api.fixturesRepo.GetByStatus(ctx, status)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			for _, f := range fixturesList {
				fixtures = append(fixtures, f)
			}
		} else {
			// Get upcoming fixtures by default
			limit := 20
			fixturesList, err := api.fixturesRepo.GetUpcoming(ctx, limit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			for _, f := range fixturesList {
				fixtures = append(fixtures, f)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"fixtures": fixtures,
			"total":    len(fixtures),
		})
	}
}

// getFixture returns single fixture handler
func (api *API) getFixture() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		fixtureID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fixture ID"})
			return
		}

		fixture, err := api.fixturesRepo.GetByID(ctx, fixtureID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "fixture not found"})
			return
		}

		// Get teams
		homeTeam, _ := api.teamsRepo.GetByID(ctx, fixture.HomeTeamID)
		awayTeam, _ := api.teamsRepo.GetByID(ctx, fixture.AwayTeamID)

		c.JSON(http.StatusOK, gin.H{
			"fixture":   fixture,
			"home_team": homeTeam,
			"away_team": awayTeam,
		})
	}
}

// getFixtureOdds returns fixture odds handler
func (api *API) getFixtureOdds() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		fixtureID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fixture ID"})
			return
		}

		// Get latest odds for the fixture
		odds, err := api.oddsRepo.GetLatestByFixture(ctx, fixtureID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get market types
		marketTypes, _ := api.oddsRepo.GetMarketTypes(ctx)

		c.JSON(http.StatusOK, gin.H{
			"fixture_id":   fixtureID,
			"odds":         odds,
			"market_types": marketTypes,
			"total":        len(odds),
		})
	}
}

// getWeeklyPicks returns weekly picks handler
func (api *API) getWeeklyPicks() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get bankroll from query or use default
		bankroll := api.cfg.InitialBankroll
		if bankrollStr := c.Query("bankroll"); bankrollStr != "" {
			if b, err := strconv.ParseFloat(bankrollStr, 64); err == nil {
				bankroll = b
			}
		}

		picks, err := api.predictionService.GetWeeklyPicks(ctx, bankroll)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Calculate summary
		var totalEV float64
		var totalStake float64
		for _, p := range picks {
			totalEV += p.ExpectedValue
			totalStake += p.SuggestedStake
		}

		c.JSON(http.StatusOK, gin.H{
			"picks": picks,
			"summary": gin.H{
				"total_picks":        len(picks),
				"total_suggested_stake": totalStake,
				"total_expected_value":  totalEV,
				"bankroll":           bankroll,
			},
		})
	}
}

// getPrediction returns prediction for a single fixture
func (api *API) getPrediction() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		fixtureID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fixture ID"})
			return
		}

		fixture, err := api.fixturesRepo.GetByID(ctx, fixtureID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "fixture not found"})
			return
		}

		prediction, err := api.predictionService.GetPrediction(ctx, fixture)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"fixture":    fixture,
			"prediction": prediction,
		})
	}
}

// getMultiMarketPicks returns weekly picks across all markets (Smart Market Selector)
func (api *API) getMultiMarketPicks() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get bankroll from query or use default
		bankroll := api.cfg.InitialBankroll
		if bankrollStr := c.Query("bankroll"); bankrollStr != "" {
			if b, err := strconv.ParseFloat(bankrollStr, 64); err == nil {
				bankroll = b
			}
		}

		// Get limit from query (default 15)
		limit := 15
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}

		picks, err := api.bettingService.GetTopPicks(ctx, bankroll, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get summary
		summary := api.bettingService.GetPicksSummary(picks, bankroll)

		c.JSON(http.StatusOK, gin.H{
			"picks":   picks,
			"summary": summary,
		})
	}
}

// evaluateFixture evaluates all markets for a single fixture
func (api *API) evaluateFixture() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		fixtureID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fixture ID"})
			return
		}

		fixture, err := api.fixturesRepo.GetByID(ctx, fixtureID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "fixture not found"})
			return
		}

		// Get bankroll from query or use default
		bankroll := api.cfg.InitialBankroll
		if bankrollStr := c.Query("bankroll"); bankrollStr != "" {
			if b, err := strconv.ParseFloat(bankrollStr, 64); err == nil {
				bankroll = b
			}
		}

		evaluation, err := api.bettingService.EvaluateFixture(ctx, fixture, bankroll)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get teams for response
		homeTeam, _ := api.teamsRepo.GetByID(ctx, fixture.HomeTeamID)
		awayTeam, _ := api.teamsRepo.GetByID(ctx, fixture.AwayTeamID)

		c.JSON(http.StatusOK, gin.H{
			"fixture":    fixture,
			"home_team":  homeTeam,
			"away_team":  awayTeam,
			"evaluation": evaluation,
		})
	}
}

// getAllMarketsMetrics returns metrics for all market models
func (api *API) getAllMarketsMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		metrics, err := api.predictionService.GetAllMarketsMetrics(ctx)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "ML service unavailable",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, metrics)
	}
}

// getWeeklyAccumulators returns weekly accumulator recommendations
func (api *API) getWeeklyAccumulators() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get bankroll from query or use default
		bankroll := api.cfg.InitialBankroll
		if bankrollStr := c.Query("bankroll"); bankrollStr != "" {
			if b, err := strconv.ParseFloat(bankrollStr, 64); err == nil {
				bankroll = b
			}
		}

		result, err := api.accumulatorService.GetWeeklyAccumulators(ctx, bankroll)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// getAccumulatorConfig returns current accumulator configuration
func (api *API) getAccumulatorConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := services.DefaultAccumulatorConfig()
		c.JSON(http.StatusOK, gin.H{
			"config": config,
			"description": gin.H{
				"min_legs":           "Minimum number of legs in accumulator",
				"max_legs":           "Maximum number of legs in accumulator",
				"min_ev_threshold":   "Minimum EV required for accumulator (5% = 0.05)",
				"min_leg_probability": "Minimum probability per leg (40% = 0.40)",
				"kelly_fraction":     "Kelly fraction for stake sizing (1/8 = 0.125)",
				"max_stake_percent":  "Maximum % of bankroll on accumulators (20% = 0.20)",
				"allow_same_team":    "Allow same team in different fixtures",
				"allow_same_fixture": "Allow multiple markets from same fixture",
			},
		})
	}
}

// getModelMetrics returns ML model performance metrics
func (api *API) getModelMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		metrics, err := api.predictionService.GetModelMetrics(ctx)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "ML service unavailable",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"model": metrics,
		})
	}
}

// getMLHealth returns ML service health status
func (api *API) getMLHealth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		healthy, err := api.predictionService.CheckMLServiceHealth(ctx)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"service": "ml-service",
				"error":   err.Error(),
			})
			return
		}

		status := "unhealthy"
		if healthy {
			status = "healthy"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"service": "ml-service",
		})
	}
}

// getBets returns bets list handler
func (api *API) getBets() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get bets
		c.JSON(http.StatusOK, gin.H{
			"bets":  []interface{}{},
			"total": 0,
		})
	}
}

// createBet returns create bet handler
func (api *API) createBet() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement create bet
		c.JSON(http.StatusCreated, gin.H{
			"id":     1,
			"status": "pending",
		})
	}
}

// settleBet returns settle bet handler
func (api *API) settleBet() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement settle bet
		betID := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"id":     betID,
			"status": "won",
		})
	}
}

// getPerformanceSummary returns performance summary handler
func (api *API) getPerformanceSummary() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement performance summary
		c.JSON(http.StatusOK, gin.H{
			"metrics": gin.H{
				"total_bets": 0,
				"roi":        0.0,
			},
		})
	}
}

// getDailyPerformance returns daily performance handler
func (api *API) getDailyPerformance() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement daily performance
		c.JSON(http.StatusOK, gin.H{
			"daily_performance": []interface{}{},
		})
	}
}

// getBankrollHistory returns bankroll history handler
func (api *API) getBankrollHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement bankroll history
		c.JSON(http.StatusOK, gin.H{
			"history": []interface{}{},
		})
	}
}

// ===============================================================
// Manual Entry Handlers - For entering fixtures and odds manually
// ===============================================================

// getTeams returns all teams for selection dropdowns
func (api *API) getTeams() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		teams, err := api.teamsRepo.GetAll(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"teams": teams,
			"total": len(teams),
		})
	}
}

// createManualFixture creates a fixture manually
func (api *API) createManualFixture() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req ManualFixtureRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse match date
		matchDate, err := time.Parse(time.RFC3339, req.MatchDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match_date format, use RFC3339 (e.g., 2025-01-20T15:00:00Z)"})
			return
		}

		// Validate teams exist
		homeTeam, err := api.teamsRepo.GetByID(ctx, req.HomeTeamID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "home team not found"})
			return
		}

		awayTeam, err := api.teamsRepo.GetByID(ctx, req.AwayTeamID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "away team not found"})
			return
		}

		// Validate teams are different
		if req.HomeTeamID == req.AwayTeamID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "home team and away team must be different"})
			return
		}

		// Generate a manual API Football ID (negative to distinguish from real API IDs)
		manualAPIID := -int(time.Now().UnixNano() % 1000000000)

		// Set default round if not provided
		round := req.Round
		if round == "" {
			round = "Manual Entry"
		}

		fixture := &models.Fixture{
			APIFootballID: manualAPIID,
			Season:        req.Season,
			Round:         round,
			MatchDate:     matchDate,
			HomeTeamID:    req.HomeTeamID,
			AwayTeamID:    req.AwayTeamID,
			Status:        "NS", // Not Started
			VenueName:     req.VenueName,
		}

		if err := api.fixturesRepo.Create(ctx, fixture); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create fixture: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"fixture": fixture,
			"home_team": homeTeam,
			"away_team": awayTeam,
			"message": "Fixture created successfully. Now add odds using POST /api/odds/manual",
		})
	}
}

// createManualOdds adds odds for a fixture manually
func (api *API) createManualOdds() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req ManualOddsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate fixture exists
		fixture, err := api.fixturesRepo.GetByID(ctx, req.FixtureID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fixture not found"})
			return
		}

		// Validate odds value
		if req.OddsValue <= 1.0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "odds_value must be greater than 1.0"})
			return
		}

		// Validate market type and outcome
		if !isValidMarketOutcome(req.MarketType, req.Outcome) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid market_type/outcome combination",
				"valid_combinations": gin.H{
					"h2h":    []string{"Home", "Draw", "Away"},
					"totals": []string{"Over", "Under"},
					"btts":   []string{"Yes", "No"},
				},
			})
			return
		}

		odds := &models.Odds{
			FixtureID:  req.FixtureID,
			Bookmaker:  req.Bookmaker,
			MarketType: req.MarketType,
			Outcome:    req.Outcome,
			OddsValue:  req.OddsValue,
			Timestamp:  time.Now(),
		}

		if err := api.oddsRepo.Create(ctx, odds); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create odds: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"odds":    odds,
			"fixture": fixture,
			"message": "Odds added successfully",
		})
	}
}

// createManualOddsBatch adds multiple odds for a fixture at once
func (api *API) createManualOddsBatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req ManualOddsBatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate fixture exists
		fixture, err := api.fixturesRepo.GetByID(ctx, req.FixtureID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fixture not found"})
			return
		}

		// Validate and prepare odds
		var oddsList []models.Odds
		now := time.Now()

		for i, entry := range req.Odds {
			if entry.OddsValue <= 1.0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "odds_value must be greater than 1.0",
					"index": i,
				})
				return
			}

			if !isValidMarketOutcome(entry.MarketType, entry.Outcome) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid market_type/outcome combination",
					"index": i,
					"valid_combinations": gin.H{
						"h2h":    []string{"Home", "Draw", "Away"},
						"totals": []string{"Over", "Under"},
						"btts":   []string{"Yes", "No"},
					},
				})
				return
			}

			oddsList = append(oddsList, models.Odds{
				FixtureID:  req.FixtureID,
				Bookmaker:  req.Bookmaker,
				MarketType: entry.MarketType,
				Outcome:    entry.Outcome,
				OddsValue:  entry.OddsValue,
				Timestamp:  now,
			})
		}

		// Insert all odds
		if err := api.oddsRepo.CreateBatch(ctx, oddsList); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create odds: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"odds_count": len(oddsList),
			"fixture":    fixture,
			"message":    "Odds added successfully. Fixture is now ready for predictions.",
		})
	}
}

// getManualFixtures returns manually entered upcoming fixtures
func (api *API) getManualFixtures() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get upcoming fixtures (includes manual entries)
		fixtures, err := api.fixturesRepo.GetUpcoming(ctx, 50)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Enrich with team names and odds status
		type EnrichedFixture struct {
			models.Fixture
			HomeTeamName string `json:"home_team_name"`
			AwayTeamName string `json:"away_team_name"`
			HasOdds      bool   `json:"has_odds"`
			OddsCount    int    `json:"odds_count"`
		}

		var enriched []EnrichedFixture
		for _, f := range fixtures {
			ef := EnrichedFixture{Fixture: f}

			// Get team names
			if homeTeam, err := api.teamsRepo.GetByID(ctx, f.HomeTeamID); err == nil {
				ef.HomeTeamName = homeTeam.Name
			}
			if awayTeam, err := api.teamsRepo.GetByID(ctx, f.AwayTeamID); err == nil {
				ef.AwayTeamName = awayTeam.Name
			}

			// Check if odds exist
			odds, err := api.oddsRepo.GetLatestByFixture(ctx, f.ID)
			if err == nil {
				ef.HasOdds = len(odds) > 0
				ef.OddsCount = len(odds)
			}

			enriched = append(enriched, ef)
		}

		c.JSON(http.StatusOK, gin.H{
			"fixtures": enriched,
			"total":    len(enriched),
		})
	}
}

// deleteManualFixture deletes a manually entered fixture
func (api *API) deleteManualFixture() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		fixtureID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fixture ID"})
			return
		}

		// Verify fixture exists
		fixture, err := api.fixturesRepo.GetByID(ctx, fixtureID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "fixture not found"})
			return
		}

		// Only allow deleting upcoming fixtures (NS status)
		if fixture.Status != "NS" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can only delete upcoming fixtures (status: NS)"})
			return
		}

		if err := api.fixturesRepo.Delete(ctx, fixtureID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete fixture: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Fixture deleted successfully",
			"fixture_id": fixtureID,
		})
	}
}

// isValidMarketOutcome validates market type and outcome combinations
func isValidMarketOutcome(marketType, outcome string) bool {
	validCombinations := map[string][]string{
		"h2h":    {"Home", "Draw", "Away"},
		"totals": {"Over", "Under"},
		"btts":   {"Yes", "No"},
	}

	validOutcomes, exists := validCombinations[marketType]
	if !exists {
		return false
	}

	for _, valid := range validOutcomes {
		if outcome == valid {
			return true
		}
	}
	return false
}
