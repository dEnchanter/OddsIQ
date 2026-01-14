package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	APIFootballKey   string
	OddsAPIKey       string
	MLServiceURL     string
	Port             string
	Env              string
	InitialBankroll  float64
	KellyFraction    float64
	MinEVThreshold   float64
	MaxBetPercentage float64
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	initialBankroll, _ := strconv.ParseFloat(getEnv("INITIAL_BANKROLL", "10000.00"), 64)
	kellyFraction, _ := strconv.ParseFloat(getEnv("KELLY_FRACTION", "0.25"), 64)
	minEVThreshold, _ := strconv.ParseFloat(getEnv("MIN_EV_THRESHOLD", "0.03"), 64)
	maxBetPercentage, _ := strconv.ParseFloat(getEnv("MAX_BET_PERCENTAGE", "0.05"), 64)

	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://localhost:5432/oddsiq?sslmode=disable"),
		APIFootballKey:   getEnv("API_FOOTBALL_KEY", ""),
		OddsAPIKey:       getEnv("ODDS_API_KEY", ""),
		MLServiceURL:     getEnv("ML_SERVICE_URL", "http://localhost:8001"),
		Port:             getEnv("PORT", "8000"),
		Env:              getEnv("ENV", "development"),
		InitialBankroll:  initialBankroll,
		KellyFraction:    kellyFraction,
		MinEVThreshold:   minEVThreshold,
		MaxBetPercentage: maxBetPercentage,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
