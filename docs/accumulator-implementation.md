# Accumulator (Parlay) Implementation Guide

## Overview

This guide details how to implement the accumulator (parlay) feature in Week 7 of the MVP.

**What are Accumulators?**
- Multiple bets combined into one
- All legs must win for payout
- Odds multiply together
- Higher potential returns, higher risk

**Example:**
```
3-leg accumulator:
✓ Arsenal Over 2.5 @ 1.80
✓ Brighton BTTS Yes @ 1.70
✓ Liverpool Away @ 2.20

Combined Odds: 1.80 × 1.70 × 2.20 = 6.73
Stake: $75
If all win: Payout = $504.75
Profit: $429.75
```

## Implementation Steps

### Step 1: Database Migration (15 min)

Run the migration to add accumulator tables:

```bash
psql -U oddsiq_user -d oddsiq -f database/migrations/003_add_accumulators.up.sql
```

This creates:
- `accumulators` table - Tracks accumulator bets
- `accumulator_legs` table - Links bets to accumulators
- `accumulator_details` view - Easy querying

### Step 2: Go Models (30 min)

Add accumulator models to `backend/internal/models/models.go`:

```go
// Accumulator represents a parlay bet
type Accumulator struct {
    ID                 int       `json:"id"`
    Name               string    `json:"name"`
    NumLegs            int       `json:"num_legs"`
    Stake              float64   `json:"stake"`
    CombinedOdds       float64   `json:"combined_odds"`
    CombinedProbability float64  `json:"combined_probability"`
    ExpectedValue      float64   `json:"expected_value"`
    PotentialPayout    float64   `json:"potential_payout"`
    PlacedAt           time.Time `json:"placed_at"`
    Status             string    `json:"status"` // pending, won, lost, void
    ActualPayout       *float64  `json:"actual_payout"`
    ProfitLoss         *float64  `json:"profit_loss"`
    SettledAt          *time.Time `json:"settled_at"`
    Notes              string    `json:"notes"`
    Legs               []Bet     `json:"legs"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
}

// AccumulatorRecommendation for weekly picks
type AccumulatorRecommendation struct {
    Name               string    `json:"name"`
    Legs               []WeeklyPick `json:"legs"`
    NumLegs            int       `json:"num_legs"`
    CombinedOdds       float64   `json:"combined_odds"`
    CombinedProbability float64  `json:"combined_probability"`
    ExpectedValue      float64   `json:"expected_value"`
    EVPercentage       float64   `json:"ev_percentage"`
    SuggestedStake     float64   `json:"suggested_stake"`
    PotentialPayout    float64   `json:"potential_payout"`
    Confidence         string    `json:"confidence"`
}
```

### Step 3: Accumulator Builder Service (3 hours)

Create `backend/internal/services/accumulator.go`:

```go
package services

import (
    "math"
    "sort"
    "github.com/oddsiq/backend/internal/models"
    "github.com/oddsiq/backend/config"
)

type AccumulatorService struct {
    cfg *config.Config
}

func NewAccumulatorService(cfg *config.Config) *AccumulatorService {
    return &AccumulatorService{cfg: cfg}
}

// BuildAccumulators generates 2-3 accumulator recommendations
func (s *AccumulatorService) BuildAccumulators(
    picks []models.WeeklyPick,
    bankroll float64,
) []models.AccumulatorRecommendation {

    // 1. Filter high-quality picks (EV > 3%)
    highQuality := s.filterHighQualityPicks(picks, 0.03)

    // 2. Remove correlated picks
    uncorrelated := s.removeCorrelatedPicks(highQuality)

    // 3. Generate 2-3 leg combinations
    combinations := s.generateCombinations(uncorrelated, 2, 3)

    // 4. Calculate accumulator metrics
    accumulators := []models.AccumulatorRecommendation{}
    for i, combo := range combinations {
        acc := s.evaluateAccumulator(combo, i+1, bankroll)

        // Filter by minimum EV (5%)
        if acc.ExpectedValue >= 0.05 {
            accumulators = append(accumulators, acc)
        }
    }

    // 5. Sort by EV, take top 2-3
    sort.Slice(accumulators, func(i, j int) bool {
        return accumulators[i].ExpectedValue > accumulators[j].ExpectedValue
    })

    maxAccumulators := 3
    if len(accumulators) > maxAccumulators {
        accumulators = accumulators[:maxAccumulators]
    }

    // 6. Ensure total stake <= 20% of weekly budget
    accumulators = s.adjustStakesForBudget(accumulators, bankroll)

    return accumulators
}

func (s *AccumulatorService) filterHighQualityPicks(
    picks []models.WeeklyPick,
    minEV float64,
) []models.WeeklyPick {
    filtered := []models.WeeklyPick{}
    for _, pick := range picks {
        if pick.ExpectedValue >= minEV {
            filtered = append(filtered, pick)
        }
    }
    return filtered
}

func (s *AccumulatorService) removeCorrelatedPicks(
    picks []models.WeeklyPick,
) []models.WeeklyPick {
    usedFixtures := make(map[int]bool)
    uncorrelated := []models.WeeklyPick{}

    for _, pick := range picks {
        fixtureID := pick.Fixture.ID

        // Skip if fixture already used
        if usedFixtures[fixtureID] {
            continue
        }

        uncorrelated = append(uncorrelated, pick)
        usedFixtures[fixtureID] = true
    }

    return uncorrelated
}

func (s *AccumulatorService) generateCombinations(
    picks []models.WeeklyPick,
    minLegs, maxLegs int,
) [][]models.WeeklyPick {
    combinations := [][]models.WeeklyPick{}

    // Generate 2-leg combinations
    for i := 0; i < len(picks); i++ {
        for j := i + 1; j < len(picks); j++ {
            combo := []models.WeeklyPick{picks[i], picks[j]}
            combinations = append(combinations, combo)
        }
    }

    // Generate 3-leg combinations
    for i := 0; i < len(picks); i++ {
        for j := i + 1; j < len(picks); j++ {
            for k := j + 1; k < len(picks); k++ {
                combo := []models.WeeklyPick{picks[i], picks[j], picks[k]}
                combinations = append(combinations, combo)
            }
        }
    }

    return combinations
}

func (s *AccumulatorService) evaluateAccumulator(
    legs []models.WeeklyPick,
    number int,
    bankroll float64,
) models.AccumulatorRecommendation {

    // Calculate combined probability and odds
    combinedProb := 1.0
    combinedOdds := 1.0

    for _, leg := range legs {
        // Get probability from prediction
        prob := s.getOutcomeProbability(leg)
        combinedProb *= prob
        combinedOdds *= leg.BestOdds
    }

    // Calculate EV
    ev := (combinedProb * combinedOdds) - 1

    // Calculate stake using 1/8 Kelly
    stake := s.calculateKellyStake(combinedProb, combinedOdds, bankroll, 0.125)

    // Cap stake at reasonable maximum
    maxStake := bankroll * 0.02 // 2% max per accumulator
    if stake > maxStake {
        stake = maxStake
    }

    potentialPayout := stake * combinedOdds

    return models.AccumulatorRecommendation{
        Name:                fmt.Sprintf("Accumulator #%d", number),
        Legs:                legs,
        NumLegs:             len(legs),
        CombinedOdds:        combinedOdds,
        CombinedProbability: combinedProb,
        ExpectedValue:       ev,
        EVPercentage:        ev * 100,
        SuggestedStake:      stake,
        PotentialPayout:     potentialPayout,
        Confidence:          s.getConfidenceLevel(ev),
    }
}

func (s *AccumulatorService) getOutcomeProbability(pick models.WeeklyPick) float64 {
    // Extract probability from prediction based on bet type
    switch pick.BetType {
    case "h2h_home":
        return pick.Prediction.HomeWinProb
    case "h2h_draw":
        return pick.Prediction.DrawProb
    case "h2h_away":
        return pick.Prediction.AwayWinProb
    case "totals_over_2_5":
        return pick.Prediction.OverProb // Assuming added to prediction
    case "totals_under_2_5":
        return pick.Prediction.UnderProb
    case "btts_yes":
        return pick.Prediction.BTTSYesProb
    case "btts_no":
        return pick.Prediction.BTTSNoProb
    default:
        return 0.5 // Fallback
    }
}

func (s *AccumulatorService) calculateKellyStake(
    probability, odds, bankroll, fraction float64,
) float64 {
    // Kelly formula: (bp - q) / b
    // b = odds - 1
    // p = probability
    // q = 1 - p

    b := odds - 1
    p := probability
    q := 1 - p

    kelly := (b*p - q) / b

    if kelly <= 0 {
        return 0
    }

    // Apply fraction (1/8 for accumulators)
    stake := bankroll * kelly * fraction

    // Minimum stake
    if stake < 10 {
        return 0
    }

    return math.Round(stake*100) / 100
}

func (s *AccumulatorService) adjustStakesForBudget(
    accumulators []models.AccumulatorRecommendation,
    bankroll float64,
) []models.AccumulatorRecommendation {

    // Calculate total stake
    totalStake := 0.0
    for _, acc := range accumulators {
        totalStake += acc.SuggestedStake
    }

    // Max 20% of bankroll on accumulators
    maxBudget := bankroll * 0.20

    // If over budget, scale down proportionally
    if totalStake > maxBudget {
        scale := maxBudget / totalStake
        for i := range accumulators {
            accumulators[i].SuggestedStake *= scale
            accumulators[i].SuggestedStake = math.Round(accumulators[i].SuggestedStake*100) / 100
            accumulators[i].PotentialPayout = accumulators[i].SuggestedStake * accumulators[i].CombinedOdds
        }
    }

    return accumulators
}

func (s *AccumulatorService) getConfidenceLevel(ev float64) string {
    if ev >= 0.15 {
        return "high"
    } else if ev >= 0.08 {
        return "medium"
    }
    return "low"
}
```

### Step 4: API Endpoints (1 hour)

Add to `backend/internal/api/handlers.go`:

```go
// getWeeklyAccumulators returns accumulator recommendations
func getWeeklyAccumulators(db *database.DB, cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get weekly singles
        singles := getWeeklyPicksInternal(db, cfg)

        // Build accumulators
        accService := services.NewAccumulatorService(cfg)
        bankroll := cfg.InitialBankroll // Or get current bankroll
        accumulators := accService.BuildAccumulators(singles, bankroll)

        c.JSON(http.StatusOK, gin.H{
            "accumulators": accumulators,
            "summary": gin.H{
                "total_accumulators": len(accumulators),
                "total_stake":        calculateTotalStake(accumulators),
                "avg_ev":             calculateAvgEV(accumulators),
                "avg_legs":           calculateAvgLegs(accumulators),
            },
        })
    }
}

// createAccumulator records an accumulator placement
func createAccumulator(db *database.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO: Implement create accumulator
        c.JSON(http.StatusCreated, gin.H{
            "id":     1,
            "status": "pending",
        })
    }
}

// settleAccumulator settles an accumulator (all legs must win)
func settleAccumulator(db *database.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO: Implement settle accumulator
        accumulatorID := c.Param("id")
        c.JSON(http.StatusOK, gin.H{
            "id":     accumulatorID,
            "status": "won",
        })
    }
}
```

Add routes in `backend/internal/api/routes.go`:

```go
// Accumulator endpoints
accumulators := v1.Group("/accumulators")
{
    accumulators.GET("/weekly", getWeeklyAccumulators(db, cfg))
    accumulators.POST("", createAccumulator(db))
    accumulators.PUT("/:id/settle", settleAccumulator(db))
}
```

### Step 5: Frontend Components (2 hours)

Create `frontend/components/AccumulatorCard.tsx`:

```typescript
interface AccumulatorCardProps {
  accumulator: {
    name: string;
    legs: any[];
    num_legs: number;
    combined_odds: number;
    ev_percentage: number;
    suggested_stake: number;
    potential_payout: number;
    confidence: string;
  };
}

export function AccumulatorCard({ accumulator }: AccumulatorCardProps) {
  return (
    <div className="border rounded-lg p-6 bg-gradient-to-br from-blue-50 to-purple-50">
      <div className="flex justify-between items-start mb-4">
        <div>
          <h3 className="text-xl font-bold">{accumulator.name}</h3>
          <p className="text-sm text-gray-600">{accumulator.num_legs} legs</p>
        </div>
        <div className="text-right">
          <div className="text-2xl font-bold text-blue-600">
            {accumulator.combined_odds.toFixed(2)}
          </div>
          <div className="text-xs text-gray-500">Combined Odds</div>
        </div>
      </div>

      {/* Legs */}
      <div className="space-y-2 mb-4">
        {accumulator.legs.map((leg, idx) => (
          <div key={idx} className="flex items-center justify-between bg-white rounded p-2">
            <div className="flex-1">
              <div className="font-semibold text-sm">
                {leg.fixture.home_team} vs {leg.fixture.away_team}
              </div>
              <div className="text-xs text-gray-600">
                {leg.recommendation.bet_type} @ {leg.recommendation.best_odds.toFixed(2)}
              </div>
            </div>
            <div className="text-sm font-semibold text-green-600">
              {(leg.recommendation.model_probability * 100).toFixed(1)}%
            </div>
          </div>
        ))}
      </div>

      {/* Stake and Payout */}
      <div className="border-t pt-4 grid grid-cols-2 gap-4">
        <div>
          <div className="text-xs text-gray-600">Suggested Stake</div>
          <div className="text-lg font-bold">${accumulator.suggested_stake.toFixed(2)}</div>
        </div>
        <div>
          <div className="text-xs text-gray-600">Potential Payout</div>
          <div className="text-lg font-bold text-green-600">
            ${accumulator.potential_payout.toFixed(2)}
          </div>
        </div>
      </div>

      {/* EV Badge */}
      <div className="mt-4">
        <span className={`inline-block px-3 py-1 rounded-full text-sm font-semibold ${
          accumulator.ev_percentage > 10
            ? 'bg-green-100 text-green-800'
            : 'bg-yellow-100 text-yellow-800'
        }`}>
          EV: +{accumulator.ev_percentage.toFixed(1)}%
        </span>
        <span className="ml-2 text-xs text-gray-600">
          Confidence: {accumulator.confidence}
        </span>
      </div>
    </div>
  );
}
```

## Testing Checklist

- [ ] Accumulators generate from uncorrelated picks
- [ ] No accumulator contains picks from same fixture
- [ ] All accumulators have EV ≥ 5%
- [ ] Stake uses 1/8 Kelly (more conservative)
- [ ] Total accumulator stake ≤ 20% of weekly budget
- [ ] Dashboard displays accumulators correctly
- [ ] Legs are shown in order
- [ ] Settlement logic works (all legs must win)
- [ ] Performance tracking separate from singles

## Example Output

**Week 7 Picks:**

**Singles (15 picks):**
- Arsenal Over 2.5 @ 1.80 (EV: +8.3%) - Stake: $125
- Brighton BTTS @ 1.70 (EV: +6.1%) - Stake: $110
- Liverpool Away @ 2.20 (EV: +5.2%) - Stake: $95
- ... (12 more)

**Accumulators (2 picks):**

**Accumulator #1** (3 legs)
- Arsenal Over 2.5 @ 1.80
- Brighton BTTS @ 1.70
- Liverpool Away @ 2.20
- **Combined Odds**: 6.73
- **EV**: +63.5%
- **Stake**: $75
- **Potential**: $504.75

**Accumulator #2** (2 legs)
- Man City -1 @ 2.10
- Tottenham BTTS @ 1.65
- **Combined Odds**: 3.47
- **EV**: +28.2%
- **Stake**: $60
- **Potential**: $208.00

**Total Weekly Stake**: $1,850 singles + $135 accumulators = $1,985

## Performance Monitoring

Track these metrics separately:

**Singles:**
- Count, total stake, wins, ROI, win rate

**Accumulators:**
- Count, total stake, wins, ROI, win rate, avg legs, avg odds

**Overall:**
- Weighted ROI, total profit, Sharpe ratio

This helps identify if accumulators outperform or underperform singles.

## Next Steps

After Week 7 implementation:
1. Paper trade for 2-3 weeks
2. Monitor accumulator win rate (expect 10-30%)
3. Compare ROI vs singles
4. Adjust parameters if needed
5. Document results for investors
