package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/dEnchanter/OddsIQ/backend/config"
	"github.com/dEnchanter/OddsIQ/backend/internal/repository"
	"github.com/dEnchanter/OddsIQ/backend/internal/services"
)

// API holds all the dependencies for handlers
type API struct {
	db                *pgxpool.Pool
	cfg               *config.Config
	teamsRepo         *repository.TeamsRepository
	fixturesRepo      *repository.FixturesRepository
	oddsRepo          *repository.OddsRepository
	statsRepo         *repository.TeamStatsRepository
	predictionService *services.PredictionService
}

// NewAPI creates a new API instance
func NewAPI(db *pgxpool.Pool, cfg *config.Config) *API {
	fixturesRepo := repository.NewFixturesRepository(db)
	oddsRepo := repository.NewOddsRepository(db)

	return &API{
		db:                db,
		cfg:               cfg,
		teamsRepo:         repository.NewTeamsRepository(db),
		fixturesRepo:      fixturesRepo,
		oddsRepo:          oddsRepo,
		statsRepo:         repository.NewTeamStatsRepository(db),
		predictionService: services.NewPredictionService(cfg, fixturesRepo, oddsRepo),
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
