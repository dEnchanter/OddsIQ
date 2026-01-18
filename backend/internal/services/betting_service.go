package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/dEnchanter/OddsIQ/backend/config"
	"github.com/dEnchanter/OddsIQ/backend/internal/models"
	"github.com/dEnchanter/OddsIQ/backend/internal/repository"
)

// MarketType represents different betting markets
type MarketType string

const (
	MarketType1X2       MarketType = "1x2"
	MarketTypeOverUnder MarketType = "over_under"
	MarketTypeBTTS      MarketType = "btts"
)

// BetOutcome represents a specific betting outcome within a market
type BetOutcome struct {
	Market      MarketType `json:"market"`
	Outcome     string     `json:"outcome"`      // e.g., "home_win", "over_2_5", "yes"
	Description string     `json:"description"`  // Human-readable description
	Probability float64    `json:"probability"`  // Model probability
	BestOdds    float64    `json:"best_odds"`    // Best available odds
	Bookmaker   string     `json:"bookmaker"`    // Source of odds
	EV          float64    `json:"ev"`           // Expected Value
	EVPercent   float64    `json:"ev_percent"`   // EV as percentage
	KellyStake  float64    `json:"kelly_stake"`  // Recommended stake (Kelly)
	Confidence  float64    `json:"confidence"`   // Model confidence
}

// MultiMarketPick represents a recommended bet with all market options evaluated
type MultiMarketPick struct {
	Fixture          models.Fixture   `json:"fixture"`
	AllOutcomes      []BetOutcome     `json:"all_outcomes"`      // All evaluated outcomes
	BestOutcome      *BetOutcome      `json:"best_outcome"`      // Highest EV outcome
	ValueOutcomes    []BetOutcome     `json:"value_outcomes"`    // All outcomes with +EV
	SuggestedStake   float64          `json:"suggested_stake"`   // Stake for best outcome
	TotalEV          float64          `json:"total_ev"`          // Sum of positive EVs
	EvaluatedAt      time.Time        `json:"evaluated_at"`
}

// BettingService handles betting calculations and recommendations
type BettingService struct {
	mlClient     *MLClient
	fixturesRepo *repository.FixturesRepository
	oddsRepo     *repository.OddsRepository
	config       *config.Config
}

// NewBettingService creates a new betting service
func NewBettingService(
	cfg *config.Config,
	mlClient *MLClient,
	fixturesRepo *repository.FixturesRepository,
	oddsRepo *repository.OddsRepository,
) *BettingService {
	return &BettingService{
		mlClient:     mlClient,
		fixturesRepo: fixturesRepo,
		oddsRepo:     oddsRepo,
		config:       cfg,
	}
}

// CalculateEV calculates Expected Value for a bet
// EV = (probability * odds) - 1
func (s *BettingService) CalculateEV(probability, odds float64) float64 {
	return (probability * odds) - 1
}

// CalculateKellyStake calculates optimal stake using Kelly Criterion
// Kelly formula: f* = (bp - q) / b
// where b = odds - 1, p = probability of winning, q = 1 - p
func (s *BettingService) CalculateKellyStake(probability, odds, bankroll float64, market MarketType) float64 {
	b := odds - 1
	p := probability
	q := 1 - p

	if b <= 0 {
		return 0
	}

	kellyFraction := (b*p - q) / b

	// Apply fractional Kelly based on market type
	fraction := s.config.KellyFraction // Default 1/4 Kelly

	// Different Kelly fractions per market (O/U and BTTS are riskier)
	switch market {
	case MarketType1X2:
		fraction = s.config.KellyFraction // 0.25
	case MarketTypeOverUnder:
		fraction = s.config.KellyFraction * 0.8 // Slightly more conservative
	case MarketTypeBTTS:
		fraction = s.config.KellyFraction * 0.8 // Slightly more conservative
	}

	adjustedKelly := kellyFraction * fraction

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

// GetOutcomeDescription returns a human-readable description for an outcome
func GetOutcomeDescription(market MarketType, outcome string) string {
	descriptions := map[MarketType]map[string]string{
		MarketType1X2: {
			"home_win": "Home Win",
			"draw":     "Draw",
			"away_win": "Away Win",
		},
		MarketTypeOverUnder: {
			"under_2_5": "Under 2.5 Goals",
			"over_2_5":  "Over 2.5 Goals",
		},
		MarketTypeBTTS: {
			"no":  "BTTS No",
			"yes": "BTTS Yes",
		},
	}

	if marketDescs, ok := descriptions[market]; ok {
		if desc, ok := marketDescs[outcome]; ok {
			return desc
		}
	}
	return outcome
}

// EvaluateFixture evaluates all markets for a single fixture
func (s *BettingService) EvaluateFixture(
	ctx context.Context,
	fixture *models.Fixture,
	bankroll float64,
) (*MultiMarketPick, error) {
	// Get multi-market predictions from ML service
	predictions, err := s.mlClient.PredictMultiMarket(ctx, fixture)
	if err != nil {
		return nil, fmt.Errorf("failed to get predictions: %w", err)
	}

	// Get odds for all markets
	odds, err := s.oddsRepo.GetLatestByFixture(ctx, fixture.ID)
	if err != nil {
		log.Printf("Warning: Could not get odds for fixture %d: %v", fixture.ID, err)
		// Continue with synthetic odds
	}

	// Build odds map by market/outcome
	oddsMap := s.buildOddsMap(odds, predictions)

	// Evaluate all outcomes
	var allOutcomes []BetOutcome
	var valueOutcomes []BetOutcome

	for marketStr, marketPred := range predictions.Predictions {
		market := MarketType(marketStr)

		for outcome, prob := range marketPred.Probabilities {
			oddsKey := fmt.Sprintf("%s_%s", marketStr, outcome)
			bestOdds, bookmaker := oddsMap[oddsKey], "synthetic"

			// If no real odds, use synthetic odds (fair odds with 5% margin)
			if bestOdds == 0 && prob > 0 {
				bestOdds = (1.0 / prob) * 0.95
			}

			if bestOdds <= 1 {
				continue // Invalid odds
			}

			ev := s.CalculateEV(prob, bestOdds)
			stake := s.CalculateKellyStake(prob, bestOdds, bankroll, market)

			betOutcome := BetOutcome{
				Market:      market,
				Outcome:     outcome,
				Description: GetOutcomeDescription(market, outcome),
				Probability: prob,
				BestOdds:    bestOdds,
				Bookmaker:   bookmaker,
				EV:          ev,
				EVPercent:   ev * 100,
				KellyStake:  math.Round(stake*100) / 100,
				Confidence:  marketPred.Confidence,
			}

			allOutcomes = append(allOutcomes, betOutcome)

			// Check if this is a value bet (meets minimum EV threshold)
			if ev >= s.config.MinEVThreshold {
				valueOutcomes = append(valueOutcomes, betOutcome)
			}
		}
	}

	// Sort all outcomes by EV (highest first)
	sort.Slice(allOutcomes, func(i, j int) bool {
		return allOutcomes[i].EV > allOutcomes[j].EV
	})

	sort.Slice(valueOutcomes, func(i, j int) bool {
		return valueOutcomes[i].EV > valueOutcomes[j].EV
	})

	// Find best outcome
	var bestOutcome *BetOutcome
	var suggestedStake float64
	if len(valueOutcomes) > 0 {
		bestOutcome = &valueOutcomes[0]
		suggestedStake = bestOutcome.KellyStake
	}

	// Calculate total EV from all value bets
	totalEV := 0.0
	for _, vo := range valueOutcomes {
		totalEV += vo.EV
	}

	return &MultiMarketPick{
		Fixture:        *fixture,
		AllOutcomes:    allOutcomes,
		BestOutcome:    bestOutcome,
		ValueOutcomes:  valueOutcomes,
		SuggestedStake: suggestedStake,
		TotalEV:        totalEV,
		EvaluatedAt:    time.Now(),
	}, nil
}

// buildOddsMap creates a map of odds by market_outcome key
func (s *BettingService) buildOddsMap(odds []models.Odds, predictions *MultiMarketPredictionResponse) map[string]float64 {
	oddsMap := make(map[string]float64)

	for _, odd := range odds {
		// Map market types from odds to our keys
		switch odd.MarketType {
		case "h2h", "1x2":
			// Home/Draw/Away odds
			if odd.Outcome == "Home" || odd.Outcome == "home" {
				oddsMap["1x2_home_win"] = odd.OddsValue
			} else if odd.Outcome == "Draw" || odd.Outcome == "draw" {
				oddsMap["1x2_draw"] = odd.OddsValue
			} else if odd.Outcome == "Away" || odd.Outcome == "away" {
				oddsMap["1x2_away_win"] = odd.OddsValue
			}
		case "totals", "over_under":
			// Over/Under odds
			if odd.Outcome == "Over" || odd.Outcome == "over" {
				oddsMap["over_under_over_2_5"] = odd.OddsValue
			} else if odd.Outcome == "Under" || odd.Outcome == "under" {
				oddsMap["over_under_under_2_5"] = odd.OddsValue
			}
		case "btts":
			// Both Teams To Score odds
			if odd.Outcome == "Yes" || odd.Outcome == "yes" {
				oddsMap["btts_yes"] = odd.OddsValue
			} else if odd.Outcome == "No" || odd.Outcome == "no" {
				oddsMap["btts_no"] = odd.OddsValue
			}
		}
	}

	return oddsMap
}

// GetMultiMarketWeeklyPicks generates weekly picks across all markets
func (s *BettingService) GetMultiMarketWeeklyPicks(ctx context.Context, bankroll float64) ([]*MultiMarketPick, error) {
	// Get upcoming fixtures
	fixtures, err := s.fixturesRepo.GetUpcoming(ctx, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming fixtures: %w", err)
	}

	if len(fixtures) == 0 {
		log.Println("No upcoming fixtures found")
		return []*MultiMarketPick{}, nil
	}

	var picks []*MultiMarketPick

	for i := range fixtures {
		fixture := &fixtures[i]
		pick, err := s.EvaluateFixture(ctx, fixture, bankroll)
		if err != nil {
			log.Printf("Warning: Failed to evaluate fixture %d: %v", fixture.ID, err)
			continue
		}

		// Only include fixtures with at least one value bet
		if pick.BestOutcome != nil {
			picks = append(picks, pick)
		}
	}

	// Sort picks by best outcome EV (highest first)
	sort.Slice(picks, func(i, j int) bool {
		if picks[i].BestOutcome == nil {
			return false
		}
		if picks[j].BestOutcome == nil {
			return true
		}
		return picks[i].BestOutcome.EV > picks[j].BestOutcome.EV
	})

	return picks, nil
}

// GetTopPicks returns the top N picks by EV
func (s *BettingService) GetTopPicks(ctx context.Context, bankroll float64, limit int) ([]*MultiMarketPick, error) {
	allPicks, err := s.GetMultiMarketWeeklyPicks(ctx, bankroll)
	if err != nil {
		return nil, err
	}

	if len(allPicks) <= limit {
		return allPicks, nil
	}

	return allPicks[:limit], nil
}

// PicksSummary represents a summary of weekly picks
type PicksSummary struct {
	TotalPicks         int                    `json:"total_picks"`
	TotalValueBets     int                    `json:"total_value_bets"`
	TotalSuggestedStake float64               `json:"total_suggested_stake"`
	TotalExpectedValue float64               `json:"total_expected_value"`
	PicksByMarket      map[string]int         `json:"picks_by_market"`
	AverageEV          float64               `json:"average_ev"`
	Bankroll           float64               `json:"bankroll"`
}

// GetPicksSummary calculates summary statistics for picks
func (s *BettingService) GetPicksSummary(picks []*MultiMarketPick, bankroll float64) *PicksSummary {
	summary := &PicksSummary{
		TotalPicks:    len(picks),
		PicksByMarket: make(map[string]int),
		Bankroll:      bankroll,
	}

	for _, pick := range picks {
		if pick.BestOutcome != nil {
			summary.TotalSuggestedStake += pick.SuggestedStake
			summary.TotalExpectedValue += pick.BestOutcome.EV * pick.SuggestedStake
			summary.PicksByMarket[string(pick.BestOutcome.Market)]++
		}
		summary.TotalValueBets += len(pick.ValueOutcomes)
	}

	if summary.TotalPicks > 0 {
		totalEV := 0.0
		for _, pick := range picks {
			if pick.BestOutcome != nil {
				totalEV += pick.BestOutcome.EV
			}
		}
		summary.AverageEV = totalEV / float64(summary.TotalPicks)
	}

	return summary
}
