package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/dEnchanter/OddsIQ/backend/config"
	"github.com/dEnchanter/OddsIQ/backend/internal/repository"
)

// API holds all the dependencies for handlers
type API struct {
	db           *pgxpool.Pool
	cfg          *config.Config
	teamsRepo    *repository.TeamsRepository
	fixturesRepo *repository.FixturesRepository
	oddsRepo     *repository.OddsRepository
	statsRepo    *repository.TeamStatsRepository
}

// NewAPI creates a new API instance
func NewAPI(db *pgxpool.Pool, cfg *config.Config) *API {
	return &API{
		db:           db,
		cfg:          cfg,
		teamsRepo:    repository.NewTeamsRepository(db),
		fixturesRepo: repository.NewFixturesRepository(db),
		oddsRepo:     repository.NewOddsRepository(db),
		statsRepo:    repository.NewTeamStatsRepository(db),
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
		var err error

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
func getWeeklyPicks(db *database.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement weekly picks
		c.JSON(http.StatusOK, gin.H{
			"picks": []interface{}{},
			"summary": gin.H{
				"total_picks": 0,
			},
		})
	}
}

// getBets returns bets list handler
func getBets(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get bets
		c.JSON(http.StatusOK, gin.H{
			"bets":  []interface{}{},
			"total": 0,
		})
	}
}

// createBet returns create bet handler
func createBet(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement create bet
		c.JSON(http.StatusCreated, gin.H{
			"id":     1,
			"status": "pending",
		})
	}
}

// settleBet returns settle bet handler
func settleBet(db *database.DB) gin.HandlerFunc {
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
func getPerformanceSummary(db *database.DB) gin.HandlerFunc {
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
func getDailyPerformance(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement daily performance
		c.JSON(http.StatusOK, gin.H{
			"daily_performance": []interface{}{},
		})
	}
}

// getBankrollHistory returns bankroll history handler
func getBankrollHistory(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement bankroll history
		c.JSON(http.StatusOK, gin.H{
			"history": []interface{}{},
		})
	}
}
