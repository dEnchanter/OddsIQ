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
)

// AccumulatorLeg represents a single leg in an accumulator
type AccumulatorLeg struct {
	FixtureID   int        `json:"fixture_id"`
	Fixture     models.Fixture `json:"fixture"`
	Market      MarketType `json:"market"`
	Outcome     string     `json:"outcome"`
	Description string     `json:"description"`
	Probability float64    `json:"probability"`
	Odds        float64    `json:"odds"`
	Bookmaker   string     `json:"bookmaker"`
	SingleEV    float64    `json:"single_ev"`
}

// Accumulator represents a multi-leg parlay bet
type Accumulator struct {
	ID                string           `json:"id"`
	Legs              []AccumulatorLeg `json:"legs"`
	NumLegs           int              `json:"num_legs"`
	CombinedProbability float64        `json:"combined_probability"`
	CombinedOdds      float64          `json:"combined_odds"`
	ExpectedValue     float64          `json:"expected_value"`
	EVPercent         float64          `json:"ev_percent"`
	SuggestedStake    float64          `json:"suggested_stake"`
	PotentialReturn   float64          `json:"potential_return"`
	Confidence        string           `json:"confidence"`
	GeneratedAt       time.Time        `json:"generated_at"`
}

// AccumulatorConfig holds configuration for accumulator generation
type AccumulatorConfig struct {
	MinLegs              int     // Minimum legs (default 2)
	MaxLegs              int     // Maximum legs (default 3)
	MinEVThreshold       float64 // Minimum EV for accumulators (default 5%)
	MinLegEV             float64 // Minimum EV per leg (default 0%)
	MinLegProbability    float64 // Minimum probability per leg (default 40%)
	KellyFraction        float64 // Kelly fraction for accumulators (default 1/8)
	MaxStakePercent      float64 // Max % of bankroll on accumulators (default 20%)
	AllowSameTeam        bool    // Allow same team in different fixtures
	AllowSameFixture     bool    // Allow multiple markets from same fixture (default false)
}

// DefaultAccumulatorConfig returns default configuration
func DefaultAccumulatorConfig() AccumulatorConfig {
	return AccumulatorConfig{
		MinLegs:           2,
		MaxLegs:           3,
		MinEVThreshold:    0.05,  // 5% minimum EV for accumulators
		MinLegEV:          0.0,   // Individual legs can have 0% EV
		MinLegProbability: 0.40,  // Minimum 40% probability per leg
		KellyFraction:     0.125, // 1/8 Kelly for accumulators
		MaxStakePercent:   0.20,  // Max 20% of bankroll on accumulators
		AllowSameTeam:     false, // Don't allow same team
		AllowSameFixture:  false, // Don't allow same fixture
	}
}

// AccumulatorService handles accumulator generation and calculations
type AccumulatorService struct {
	bettingService *BettingService
	config         *config.Config
	accConfig      AccumulatorConfig
}

// NewAccumulatorService creates a new accumulator service
func NewAccumulatorService(
	bettingService *BettingService,
	cfg *config.Config,
) *AccumulatorService {
	return &AccumulatorService{
		bettingService: bettingService,
		config:         cfg,
		accConfig:      DefaultAccumulatorConfig(),
	}
}

// SetConfig updates the accumulator configuration
func (s *AccumulatorService) SetConfig(cfg AccumulatorConfig) {
	s.accConfig = cfg
}

// IsCorrelated checks if two legs are correlated and should not be combined
func (s *AccumulatorService) IsCorrelated(leg1, leg2 AccumulatorLeg) bool {
	// Same fixture - always correlated
	if leg1.FixtureID == leg2.FixtureID {
		if !s.accConfig.AllowSameFixture {
			return true
		}
	}

	// Same team involved - potentially correlated
	if !s.accConfig.AllowSameTeam {
		// Check if same team appears in both fixtures
		teams1 := []int{leg1.Fixture.HomeTeamID, leg1.Fixture.AwayTeamID}
		teams2 := []int{leg2.Fixture.HomeTeamID, leg2.Fixture.AwayTeamID}

		for _, t1 := range teams1 {
			for _, t2 := range teams2 {
				if t1 == t2 {
					return true
				}
			}
		}
	}

	return false
}

// CalculateAccumulatorEV calculates EV for an accumulator
// EV = (combined_probability × combined_odds) - 1
func (s *AccumulatorService) CalculateAccumulatorEV(legs []AccumulatorLeg) (combinedProb, combinedOdds, ev float64) {
	if len(legs) == 0 {
		return 0, 0, -1
	}

	combinedProb = 1.0
	combinedOdds = 1.0

	for _, leg := range legs {
		combinedProb *= leg.Probability
		combinedOdds *= leg.Odds
	}

	ev = (combinedProb * combinedOdds) - 1
	return combinedProb, combinedOdds, ev
}

// CalculateAccumulatorKelly calculates Kelly stake for accumulator
func (s *AccumulatorService) CalculateAccumulatorKelly(combinedProb, combinedOdds, bankroll float64) float64 {
	b := combinedOdds - 1
	p := combinedProb
	q := 1 - p

	if b <= 0 {
		return 0
	}

	kellyFraction := (b*p - q) / b

	// Apply conservative Kelly fraction for accumulators
	adjustedKelly := kellyFraction * s.accConfig.KellyFraction

	// Cap at max accumulator stake percentage
	maxStake := bankroll * s.accConfig.MaxStakePercent
	stake := adjustedKelly * bankroll

	if stake > maxStake {
		stake = maxStake
	}

	// No negative stakes
	if stake < 0 {
		return 0
	}

	return math.Round(stake*100) / 100
}

// GetConfidenceLevel returns confidence level based on EV
func (s *AccumulatorService) GetConfidenceLevel(ev float64) string {
	if ev >= 0.20 {
		return "high"
	} else if ev >= 0.10 {
		return "medium"
	}
	return "low"
}

// ConvertToLeg converts a BetOutcome to an AccumulatorLeg
func (s *AccumulatorService) ConvertToLeg(outcome BetOutcome, fixture models.Fixture) AccumulatorLeg {
	return AccumulatorLeg{
		FixtureID:   fixture.ID,
		Fixture:     fixture,
		Market:      outcome.Market,
		Outcome:     outcome.Outcome,
		Description: outcome.Description,
		Probability: outcome.Probability,
		Odds:        outcome.BestOdds,
		Bookmaker:   outcome.Bookmaker,
		SingleEV:    outcome.EV,
	}
}

// FilterLegsForAccumulator filters legs suitable for accumulator
func (s *AccumulatorService) FilterLegsForAccumulator(picks []*MultiMarketPick) []AccumulatorLeg {
	var legs []AccumulatorLeg

	for _, pick := range picks {
		// Get all value outcomes from each pick
		for _, outcome := range pick.ValueOutcomes {
			// Filter by minimum probability
			if outcome.Probability < s.accConfig.MinLegProbability {
				continue
			}

			// Filter by minimum leg EV
			if outcome.EV < s.accConfig.MinLegEV {
				continue
			}

			leg := s.ConvertToLeg(outcome, pick.Fixture)
			legs = append(legs, leg)
		}
	}

	// Sort by EV (highest first)
	sort.Slice(legs, func(i, j int) bool {
		return legs[i].SingleEV > legs[j].SingleEV
	})

	return legs
}

// GenerateAccumulators generates optimal accumulators from available picks
func (s *AccumulatorService) GenerateAccumulators(
	ctx context.Context,
	bankroll float64,
	maxAccumulators int,
) ([]*Accumulator, error) {
	// Get multi-market picks
	picks, err := s.bettingService.GetMultiMarketWeeklyPicks(ctx, bankroll)
	if err != nil {
		return nil, fmt.Errorf("failed to get picks: %w", err)
	}

	if len(picks) < s.accConfig.MinLegs {
		log.Printf("Not enough picks for accumulators: %d", len(picks))
		return []*Accumulator{}, nil
	}

	// Filter legs suitable for accumulator
	allLegs := s.FilterLegsForAccumulator(picks)

	if len(allLegs) < s.accConfig.MinLegs {
		log.Printf("Not enough qualifying legs for accumulators: %d", len(allLegs))
		return []*Accumulator{}, nil
	}

	// Generate accumulators of different sizes
	var accumulators []*Accumulator

	// Generate 2-leg accumulators
	if s.accConfig.MinLegs <= 2 && s.accConfig.MaxLegs >= 2 {
		doubles := s.generateNLegAccumulators(allLegs, 2, bankroll)
		accumulators = append(accumulators, doubles...)
	}

	// Generate 3-leg accumulators
	if s.accConfig.MinLegs <= 3 && s.accConfig.MaxLegs >= 3 && len(allLegs) >= 3 {
		trebles := s.generateNLegAccumulators(allLegs, 3, bankroll)
		accumulators = append(accumulators, trebles...)
	}

	// Sort by EV
	sort.Slice(accumulators, func(i, j int) bool {
		return accumulators[i].ExpectedValue > accumulators[j].ExpectedValue
	})

	// Filter by minimum EV threshold
	var filtered []*Accumulator
	for _, acc := range accumulators {
		if acc.ExpectedValue >= s.accConfig.MinEVThreshold {
			filtered = append(filtered, acc)
		}
	}

	// Limit to max accumulators
	if len(filtered) > maxAccumulators {
		filtered = filtered[:maxAccumulators]
	}

	return filtered, nil
}

// generateNLegAccumulators generates all valid N-leg accumulators
func (s *AccumulatorService) generateNLegAccumulators(legs []AccumulatorLeg, n int, bankroll float64) []*Accumulator {
	if len(legs) < n {
		return nil
	}

	var accumulators []*Accumulator

	// Generate combinations
	combinations := s.generateCombinations(len(legs), n)

	for _, combo := range combinations {
		selectedLegs := make([]AccumulatorLeg, n)
		for i, idx := range combo {
			selectedLegs[i] = legs[idx]
		}

		// Check for correlations
		if s.hasCorrelation(selectedLegs) {
			continue
		}

		// Calculate accumulator metrics
		combinedProb, combinedOdds, ev := s.CalculateAccumulatorEV(selectedLegs)
		stake := s.CalculateAccumulatorKelly(combinedProb, combinedOdds, bankroll)

		if stake <= 0 {
			continue
		}

		acc := &Accumulator{
			ID:                  fmt.Sprintf("acc_%d_%d", n, len(accumulators)+1),
			Legs:                selectedLegs,
			NumLegs:             n,
			CombinedProbability: combinedProb,
			CombinedOdds:        math.Round(combinedOdds*100) / 100,
			ExpectedValue:       ev,
			EVPercent:           ev * 100,
			SuggestedStake:      stake,
			PotentialReturn:     math.Round(stake*combinedOdds*100) / 100,
			Confidence:          s.GetConfidenceLevel(ev),
			GeneratedAt:         time.Now(),
		}

		accumulators = append(accumulators, acc)
	}

	return accumulators
}

// hasCorrelation checks if any pair of legs in the selection are correlated
func (s *AccumulatorService) hasCorrelation(legs []AccumulatorLeg) bool {
	for i := 0; i < len(legs); i++ {
		for j := i + 1; j < len(legs); j++ {
			if s.IsCorrelated(legs[i], legs[j]) {
				return true
			}
		}
	}
	return false
}

// generateCombinations generates all combinations of n items from total
func (s *AccumulatorService) generateCombinations(total, n int) [][]int {
	var result [][]int
	combo := make([]int, n)

	var generate func(start, idx int)
	generate = func(start, idx int) {
		if idx == n {
			// Make a copy
			c := make([]int, n)
			copy(c, combo)
			result = append(result, c)
			return
		}
		for i := start; i <= total-n+idx; i++ {
			combo[idx] = i
			generate(i+1, idx+1)
		}
	}

	generate(0, 0)
	return result
}

// AccumulatorSummary represents a summary of generated accumulators
type AccumulatorSummary struct {
	TotalAccumulators   int     `json:"total_accumulators"`
	TotalDoubles        int     `json:"total_doubles"`
	TotalTrebles        int     `json:"total_trebles"`
	TotalSuggestedStake float64 `json:"total_suggested_stake"`
	TotalPotentialReturn float64 `json:"total_potential_return"`
	AverageEV           float64 `json:"average_ev"`
	BestEV              float64 `json:"best_ev"`
	Bankroll            float64 `json:"bankroll"`
	MaxStakeAllocation  float64 `json:"max_stake_allocation"`
}

// GetAccumulatorSummary calculates summary statistics for accumulators
func (s *AccumulatorService) GetAccumulatorSummary(accumulators []*Accumulator, bankroll float64) *AccumulatorSummary {
	summary := &AccumulatorSummary{
		TotalAccumulators:  len(accumulators),
		Bankroll:           bankroll,
		MaxStakeAllocation: bankroll * s.accConfig.MaxStakePercent,
	}

	if len(accumulators) == 0 {
		return summary
	}

	totalEV := 0.0
	for _, acc := range accumulators {
		summary.TotalSuggestedStake += acc.SuggestedStake
		summary.TotalPotentialReturn += acc.PotentialReturn
		totalEV += acc.ExpectedValue

		if acc.NumLegs == 2 {
			summary.TotalDoubles++
		} else if acc.NumLegs == 3 {
			summary.TotalTrebles++
		}

		if acc.ExpectedValue > summary.BestEV {
			summary.BestEV = acc.ExpectedValue
		}
	}

	summary.AverageEV = totalEV / float64(len(accumulators))

	// Cap total stake at max allocation
	if summary.TotalSuggestedStake > summary.MaxStakeAllocation {
		// Scale down stakes proportionally
		scale := summary.MaxStakeAllocation / summary.TotalSuggestedStake
		summary.TotalSuggestedStake = summary.MaxStakeAllocation
		summary.TotalPotentialReturn *= scale
	}

	return summary
}

// WeeklyAccumulatorPicks represents weekly accumulator recommendations
type WeeklyAccumulatorPicks struct {
	Accumulators []*Accumulator     `json:"accumulators"`
	Summary      *AccumulatorSummary `json:"summary"`
	Config       AccumulatorConfig   `json:"config"`
	GeneratedAt  time.Time          `json:"generated_at"`
}

// GetWeeklyAccumulators generates weekly accumulator recommendations
func (s *AccumulatorService) GetWeeklyAccumulators(ctx context.Context, bankroll float64) (*WeeklyAccumulatorPicks, error) {
	// Generate up to 3 accumulators (2 doubles + 1 treble recommended)
	accumulators, err := s.GenerateAccumulators(ctx, bankroll, 5)
	if err != nil {
		return nil, err
	}

	summary := s.GetAccumulatorSummary(accumulators, bankroll)

	return &WeeklyAccumulatorPicks{
		Accumulators: accumulators,
		Summary:      summary,
		Config:       s.accConfig,
		GeneratedAt:  time.Now(),
	}, nil
}
