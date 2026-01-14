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
		// Fixtures endpoints
		fixtures := v1.Group("/fixtures")
		{
			fixtures.GET("", api.getFixtures())
			fixtures.GET("/:id", api.getFixture())
			fixtures.GET("/:id/odds", api.getFixtureOdds())
		}

		// Picks endpoints
		picks := v1.Group("/picks")
		{
			picks.GET("/weekly", api.getWeeklyPicks())
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
