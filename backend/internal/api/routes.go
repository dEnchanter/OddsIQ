package api

import (
	"github.com/gin-gonic/gin"
	"github.com/dEnchanter/OddsIQ/backend/config"
	"github.com/dEnchanter/OddsIQ/backend/pkg/database"
)

func SetupRoutes(router *gin.Engine, db *database.DB, cfg *config.Config) {
	// Health check endpoint
	router.GET("/health", healthCheck(db))

	// API v1 group
	v1 := router.Group("/api")
	{
		// Fixtures endpoints
		fixtures := v1.Group("/fixtures")
		{
			fixtures.GET("", getFixtures(db))
			fixtures.GET("/:id", getFixture(db))
			fixtures.GET("/:id/odds", getFixtureOdds(db))
		}

		// Picks endpoints
		picks := v1.Group("/picks")
		{
			picks.GET("/weekly", getWeeklyPicks(db, cfg))
		}

		// Bets endpoints
		bets := v1.Group("/bets")
		{
			bets.GET("", getBets(db))
			bets.POST("", createBet(db))
			bets.PUT("/:id/settle", settleBet(db))
		}

		// Performance endpoints
		performance := v1.Group("/performance")
		{
			performance.GET("/summary", getPerformanceSummary(db))
			performance.GET("/daily", getDailyPerformance(db))
		}

		// Bankroll endpoints
		bankroll := v1.Group("/bankroll")
		{
			bankroll.GET("/history", getBankrollHistory(db))
		}
	}
}
