package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/dEnchanter/OddsIQ/backend/config"
	"github.com/dEnchanter/OddsIQ/backend/internal/models"
	"github.com/dEnchanter/OddsIQ/backend/internal/repository"
)

// PredictionService handles predictions and betting recommendations
type PredictionService struct {
	mlClient     *MLClient
	fixturesRepo *repository.FixturesRepository
	oddsRepo     *repository.OddsRepository
	config       *config.Config

	// Cache for predictions (fixture_id -> prediction)
	cache      map[int]*models.Prediction
	cacheMutex sync.RWMutex
	cacheTime  map[int]time.Time
	cacheTTL   time.Duration
}

// NewPredictionService creates a new prediction service
func NewPredictionService(
	cfg *config.Config,
	fixturesRepo *repository.FixturesRepository,
	oddsRepo *repository.OddsRepository,
) *PredictionService {
	return &PredictionService{
		mlClient:     NewMLClient(cfg.MLServiceURL),
		fixturesRepo: fixturesRepo,
		oddsRepo:     oddsRepo,
		config:       cfg,
		cache:        make(map[int]*models.Prediction),
		cacheTime:    make(map[int]time.Time),
		cacheTTL:     1 * time.Hour, // Cache predictions for 1 hour
	}
}

// GetPrediction gets or creates a prediction for a fixture
func (s *PredictionService) GetPrediction(ctx context.Context, fixture *models.Fixture) (*models.Prediction, error) {
	// Check cache first
	s.cacheMutex.RLock()
	if pred, ok := s.cache[fixture.ID]; ok {
		if time.Since(s.cacheTime[fixture.ID]) < s.cacheTTL {
			s.cacheMutex.RUnlock()
			return pred, nil
		}
	}
	s.cacheMutex.RUnlock()

	// Call ML service
	pred, err := s.mlClient.Predict(ctx, fixture)
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction: %w", err)
	}

	// Update cache
	s.cacheMutex.Lock()
	s.cache[fixture.ID] = pred
	s.cacheTime[fixture.ID] = time.Now()
	s.cacheMutex.Unlock()

	return pred, nil
}

// GetPredictions gets predictions for multiple fixtures
func (s *PredictionService) GetPredictions(ctx context.Context, fixtures []*models.Fixture) ([]*models.Prediction, error) {
	// Check which fixtures need predictions
	var needPrediction []*models.Fixture
	predictions := make([]*models.Prediction, len(fixtures))

	s.cacheMutex.RLock()
	for i, f := range fixtures {
		if pred, ok := s.cache[f.ID]; ok {
			if time.Since(s.cacheTime[f.ID]) < s.cacheTTL {
				predictions[i] = pred
				continue
			}
		}
		needPrediction = append(needPrediction, f)
	}
	s.cacheMutex.RUnlock()

	// Get missing predictions from ML service
	if len(needPrediction) > 0 {
		newPreds, err := s.mlClient.PredictBatch(ctx, needPrediction)
		if err != nil {
			return nil, fmt.Errorf("failed to get batch predictions: %w", err)
		}

		// Update cache and fill in predictions array
		s.cacheMutex.Lock()
		for _, pred := range newPreds {
			s.cache[pred.FixtureID] = pred
			s.cacheTime[pred.FixtureID] = time.Now()

			// Find and fill in the predictions array
			for i, f := range fixtures {
				if f.ID == pred.FixtureID {
					predictions[i] = pred
					break
				}
			}
		}
		s.cacheMutex.Unlock()
	}

	return predictions, nil
}

// CalculateExpectedValue calculates EV for a bet
func (s *PredictionService) CalculateExpectedValue(modelProb, odds float64) float64 {
	// EV = (probability * odds) - 1
	return (modelProb * odds) - 1
}

// CalculateKellyStake calculates optimal stake using Kelly Criterion
func (s *PredictionService) CalculateKellyStake(modelProb, odds, bankroll float64) float64 {
	// Kelly formula: f* = (bp - q) / b
	// where b = odds - 1, p = probability of winning, q = 1 - p
	b := odds - 1
	p := modelProb
	q := 1 - p

	kellyFraction := (b*p - q) / b

	// Apply fractional Kelly (configured fraction, default 1/4)
	adjustedKelly := kellyFraction * s.config.KellyFraction

	// Cap at max bet percentage
	if adjustedKelly > s.config.MaxBetPercentage {
		adjustedKelly = s.config.MaxBetPercentage
	}

	// No negative stakes
	if adjustedKelly < 0 {
		return 0
	}

	return adjustedKelly * bankroll
}

// GetWeeklyPicks generates betting recommendations for upcoming fixtures
func (s *PredictionService) GetWeeklyPicks(ctx context.Context, bankroll float64) ([]*models.WeeklyPick, error) {
	// Get upcoming fixtures (limit 20)
	fixtureSlice, err := s.fixturesRepo.GetUpcoming(ctx, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming fixtures: %w", err)
	}

	if len(fixtureSlice) == 0 {
		log.Println("No upcoming fixtures found")
		return []*models.WeeklyPick{}, nil
	}

	// Convert to pointer slice for predictions
	fixtures := make([]*models.Fixture, len(fixtureSlice))
	for i := range fixtureSlice {
		fixtures[i] = &fixtureSlice[i]
	}

	// Get predictions for all fixtures
	predictions, err := s.GetPredictions(ctx, fixtures)
	if err != nil {
		return nil, fmt.Errorf("failed to get predictions: %w", err)
	}

	var picks []*models.WeeklyPick

	for i, fixture := range fixtures {
		pred := predictions[i]
		if pred == nil {
			continue
		}

		// Get best odds for each outcome
		// For now, use synthetic odds based on predictions
		// TODO: Get real odds from odds repository
		homeOdds := 1.0 / pred.HomeWinProb * 0.95  // Add 5% margin
		drawOdds := 1.0 / pred.DrawProb * 0.95
		awayOdds := 1.0 / pred.AwayWinProb * 0.95

		// Check each outcome for value
		outcomes := []struct {
			betType  string
			prob     float64
			odds     float64
			outcome  string
		}{
			{"home_win", pred.HomeWinProb, homeOdds, "home_win"},
			{"draw", pred.DrawProb, drawOdds, "draw"},
			{"away_win", pred.AwayWinProb, awayOdds, "away_win"},
		}

		for _, o := range outcomes {
			ev := s.CalculateExpectedValue(o.prob, o.odds)

			// Only include picks with positive EV above threshold
			if ev >= s.config.MinEVThreshold {
				stake := s.CalculateKellyStake(o.prob, o.odds, bankroll)

				confidence := "low"
				if pred.ConfidenceScore > 0.5 {
					confidence = "medium"
				}
				if pred.ConfidenceScore > 0.6 {
					confidence = "high"
				}

				pick := &models.WeeklyPick{
					Fixture:        *fixture,
					Prediction:     *pred,
					BestOdds:       o.odds,
					Bookmaker:      "synthetic", // TODO: Get from odds repo
					ExpectedValue:  ev,
					EVPercentage:   ev * 100,
					SuggestedStake: math.Round(stake*100) / 100,
					KellyFraction:  s.config.KellyFraction,
					BetType:        o.betType,
					Confidence:     confidence,
				}

				picks = append(picks, pick)
			}
		}
	}

	// Sort by EV (highest first)
	for i := 0; i < len(picks)-1; i++ {
		for j := i + 1; j < len(picks); j++ {
			if picks[j].ExpectedValue > picks[i].ExpectedValue {
				picks[i], picks[j] = picks[j], picks[i]
			}
		}
	}

	return picks, nil
}

// GetModelMetrics returns current model performance metrics
func (s *PredictionService) GetModelMetrics(ctx context.Context) (*ModelMetricsResponse, error) {
	return s.mlClient.GetModelMetrics(ctx)
}

// CheckMLServiceHealth checks if ML service is available
func (s *PredictionService) CheckMLServiceHealth(ctx context.Context) (bool, error) {
	health, err := s.mlClient.HealthCheck(ctx)
	if err != nil {
		return false, err
	}
	return health.Status == "healthy", nil
}

// ClearCache clears the prediction cache
func (s *PredictionService) ClearCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	s.cache = make(map[int]*models.Prediction)
	s.cacheTime = make(map[int]time.Time)
}
