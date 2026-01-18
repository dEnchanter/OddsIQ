package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/dEnchanter/OddsIQ/backend/config"
)

func SetupRoutes(router *gin.Engine, db *pgxpool.Pool, cfg *config.Config) {
	// Create API instance with repositories
	api := NewAPI(db, cfg)

	// Health check endpoint
	router.GET("/health", api.healthCheck())

	// API v1 group
	v1 := router.Group("/api")
	{
		// Teams endpoint (for manual entry dropdowns)
		v1.GET("/teams", api.getTeams())

		// Fixtures endpoints
		fixtures := v1.Group("/fixtures")
		{
			fixtures.GET("", api.getFixtures())
			fixtures.GET("/upcoming", api.getManualFixtures()) // List upcoming fixtures with odds status
			fixtures.GET("/:id", api.getFixture())
			fixtures.GET("/:id/odds", api.getFixtureOdds())
			fixtures.POST("/manual", api.createManualFixture())     // Manual fixture entry
			fixtures.DELETE("/:id", api.deleteManualFixture())      // Delete fixture
		}

		// Odds endpoints (manual entry)
		odds := v1.Group("/odds")
		{
			odds.POST("/manual", api.createManualOdds())        // Add single odds entry
			odds.POST("/manual/batch", api.createManualOddsBatch()) // Add multiple odds at once
		}

		// Picks endpoints
		picks := v1.Group("/picks")
		{
			picks.GET("/weekly", api.getWeeklyPicks())             // Legacy 1X2 only
			picks.GET("/multi", api.getMultiMarketPicks())         // Smart Market Selector (all markets)
		}

		// Accumulators endpoints
		accumulators := v1.Group("/accumulators")
		{
			accumulators.GET("/weekly", api.getWeeklyAccumulators())   // Weekly accumulator recommendations
			accumulators.GET("/config", api.getAccumulatorConfig())    // Get accumulator configuration
		}

		// Predictions endpoints
		predictions := v1.Group("/predictions")
		{
			predictions.GET("/fixture/:id", api.getPrediction())
			predictions.GET("/fixture/:id/evaluate", api.evaluateFixture())  // Evaluate all markets
		}

		// Model endpoints
		model := v1.Group("/model")
		{
			model.GET("/metrics", api.getModelMetrics())
			model.GET("/metrics/all", api.getAllMarketsMetrics())  // All market models
			model.GET("/health", api.getMLHealth())
		}

		// Bets endpoints
		bets := v1.Group("/bets")
		{
			bets.GET("", api.getBets())
			bets.POST("", api.createBet())
			bets.PUT("/:id/settle", api.settleBet())
		}

		// Performance endpoints
		performance := v1.Group("/performance")
		{
			performance.GET("/summary", api.getPerformanceSummary())
			performance.GET("/daily", api.getDailyPerformance())
		}

		// Bankroll endpoints
		bankroll := v1.Group("/bankroll")
		{
			bankroll.GET("/history", api.getBankrollHistory())
		}
	}
}
