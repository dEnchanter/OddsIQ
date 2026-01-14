# Architecture Decision Records (ADR)

## ADR-001: Multi-Service Architecture with Go and Python

**Date:** 2026-01-12
**Status:** Accepted

### Context
OddsIQ MVP requires both high-performance API/data services and ML capabilities. The original roadmap specified Python + FastAPI, but user preference is for Go backend.

### Decision
Implement a multi-service architecture:
- **Go service** for API, data ingestion, business logic
- **Python service** for ML predictions and feature engineering
- Services communicate via HTTP/REST

### Rationale
**Go Backend:**
- Superior performance for API and concurrent data fetching
- Better resource efficiency
- Easier deployment and containerization
- Strong typing and robust error handling
- Native support for concurrent API calls to multiple bookmakers

**Python ML Service:**
- Required for ML libraries (XGBoost, scikit-learn, pandas)
- Excellent ecosystem for data science
- FastAPI provides fast async endpoints
- Easy integration with Jupyter notebooks for experimentation

**Service Separation:**
- Clear separation of concerns
- Independent scaling (ML predictions are compute-heavy)
- Different update cycles (ML models retrain less frequently)
- Team can work on services independently

### Consequences
**Positive:**
- Best tool for each job
- Services can scale independently
- Clear boundaries and responsibilities
- Go service handles high-frequency operations efficiently

**Negative:**
- Network latency between services (mitigated by local deployment initially)
- More complex deployment (two services instead of one)
- Need to maintain two codebases with different languages

### Alternatives Considered
1. **Pure Python** - Simpler, but slower API performance and higher resource usage
2. **Pure Go with ML libraries** - Limited ML ecosystem, would need to reimplement many algorithms
3. **Monolith** - Simpler deployment, but loses benefits of service separation

---

## ADR-002: PostgreSQL for Data Storage

**Date:** 2026-01-12
**Status:** Accepted

### Context
Need to store 3 seasons of fixtures, odds, predictions, and betting results. MVP roadmap suggests PostgreSQL.

### Decision
Use PostgreSQL 15+ for all persistent data storage, starting with local instance.

### Rationale
- **Relational data model fits domain well** - fixtures, teams, odds, bets have clear relationships
- **Time-series capabilities** - excellent for historical odds tracking and performance metrics
- **Mature and robust** - battle-tested, excellent tooling
- **JSON support** - can store flexible data (features, metadata) when needed
- **Free and open source** - no licensing costs
- **Easy cloud migration** - AWS RDS, Google Cloud SQL support when scaling

### Consequences
**Positive:**
- No need to migrate later (SQLite → PostgreSQL is painful)
- Powerful query capabilities for analytics
- ACID guarantees for bet tracking
- Excellent indexing for time-based queries

**Negative:**
- Requires PostgreSQL installation for local development
- Slightly more setup than SQLite
- Need to manage migrations carefully

### Alternatives Considered
1. **SQLite** - Simpler initial setup, but would require migration for production
2. **MongoDB** - Flexible schema, but relational queries are more natural for betting domain
3. **TimescaleDB** - Specialized for time-series, but overkill for MVP scale

---

## ADR-003: Next.js for Frontend

**Date:** 2026-01-12
**Status:** Accepted

### Context
Need investor-facing dashboard to display picks, performance metrics, and betting history. Original roadmap suggests React.

### Decision
Use Next.js 14+ with App Router for the frontend dashboard.

### Rationale
- **Built on React** - aligns with roadmap vision
- **Server-side rendering** - faster initial load, better for dashboards with real data
- **API routes** - can add BFF (Backend for Frontend) layer if needed
- **TypeScript support** - type safety across frontend
- **Modern DX** - excellent developer experience, fast refresh
- **Production-ready** - built-in optimization, image optimization
- **Active ecosystem** - shadcn/ui, Tailwind CSS integration

### Consequences
**Positive:**
- Single framework for routing, SSR, and client-side interactivity
- Better performance than pure SPA
- TypeScript ensures type safety with backend contracts
- Easy to add authentication later (NextAuth)

**Negative:**
- Larger framework than vanilla React
- Learning curve for App Router paradigm
- May be overpowered for simple dashboard (but scales well)

### Alternatives Considered
1. **Vite + React** - Lighter, faster dev server, but need separate routing solution
2. **SvelteKit** - Modern and fast, but smaller ecosystem
3. **Plain React** - Maximum simplicity, but loses SSR benefits

---

## ADR-004: Fractional Kelly Criterion for Stake Sizing

**Date:** 2026-01-12
**Status:** Accepted

### Context
MVP needs a bankroll management strategy. Full Kelly Criterion maximizes growth but has high variance.

### Decision
Implement 1/4 Kelly Criterion for stake sizing with configurable fraction.

### Formula
```
Kelly % = (bp - q) / b
Stake = Bankroll × Kelly % × Fraction

Where:
- b = decimal odds - 1
- p = model probability of winning
- q = 1 - p
- Fraction = 0.25 (1/4 Kelly)
```

### Rationale
- **Reduces variance** - fractional Kelly trades growth for stability
- **Safer for MVP** - conservative approach during validation period
- **Mathematically sound** - still based on Kelly's proven formula
- **Investor-friendly** - demonstrates risk management awareness
- **Configurable** - can adjust fraction based on confidence

### Consequences
**Positive:**
- Lower risk of significant drawdowns
- More palatable for conservative investors
- Still captures value from +EV opportunities
- Easy to explain and justify

**Negative:**
- Lower growth rate than full Kelly
- May be too conservative if model is highly accurate
- Requires accurate probability estimates

### Implementation Notes
- Minimum bet size: $10
- Maximum bet size: 5% of bankroll (safety cap)
- Only bet when EV > 3%
- Track Kelly fraction performance for future optimization

---

## ADR-005: API-Football and The Odds API for Data

**Date:** 2026-01-12
**Status:** Accepted

### Context
Need reliable sources for fixture data and odds. MVP budget is $25k-$40k.

### Decision
- **API-Football** for fixtures, teams, standings, statistics
- **The Odds API** for bookmaker odds across multiple bookmakers

### Rationale
**API-Football:**
- Comprehensive coverage of Premier League
- Historical data available (3+ seasons)
- Team stats, lineups, injuries
- Reasonable pricing (~$50-100/month for needed tier)
- RESTful API, well-documented

**The Odds API:**
- Multiple bookmakers in single API
- Historical odds available
- Real-time updates
- ~$200/month for needed tier
- Specifically designed for betting applications

**Total API costs:** ~$250-300/month (within budget)

### Consequences
**Positive:**
- Professional, reliable data sources
- Historical data for model training
- Real-time updates for live system
- Good documentation and support

**Negative:**
- Recurring monthly costs
- Dependent on third-party services
- API rate limits to manage
- Need fallback if service is down

### Alternatives Considered
1. **Web scraping** - Free but fragile, legally questionable, time-consuming
2. **Sportmonks** - More expensive, overkill for MVP
3. **Multiple free APIs** - Inconsistent data quality, no guaranteed uptime

---

## ADR-006: XGBoost for Initial ML Model

**Date:** 2026-01-12
**Status:** Accepted

### Context
Need to build prediction model that achieves 55-60% accuracy. MVP timeline is 8 weeks.

### Decision
Start with XGBoost gradient boosting model for match outcome predictions.

### Rationale
- **Industry standard** - proven in sports betting and competitions
- **Fast training** - can iterate quickly during MVP phase
- **Handles tabular data well** - perfect for engineered features
- **Built-in regularization** - reduces overfitting
- **Feature importance** - can explain which factors drive predictions
- **Good probability calibration** - important for expected value calculations

### Model Architecture
```python
XGBClassifier(
    n_estimators=200,
    max_depth=6,
    learning_rate=0.1,
    objective='multi:softprob',  # Home/Draw/Away
    eval_metric='logloss'
)
```

### Features (15-20 core features)
- Form metrics (last 5 games points)
- Goals for/against averages
- Home/away performance splits
- Head-to-head history
- League position differential
- xG metrics (if available)

### Consequences
**Positive:**
- Quick to implement and train
- Interpretable results
- Can achieve target 55-60% accuracy
- Easy to add features incrementally

**Negative:**
- May plateau at ~60% accuracy (okay for MVP)
- Not ideal for sequence data (no match context)
- Requires good feature engineering

### Future Enhancements (Post-MVP)
- Ensemble with CatBoost, LightGBM
- Neural network for comparison
- Separate models for different markets (O/U, BTTS)

---

## ADR-007: Minimum 3% EV Threshold for Bet Selection

**Date:** 2026-01-12
**Status:** Accepted

### Context
Need criteria for filtering predictions into actionable bets. MVP should be conservative.

### Decision
Only recommend bets with Expected Value (EV) ≥ 3%.

### Formula
```
EV = (Model Probability × Odds) - 1
Bet if EV ≥ 0.03 (3%)
```

### Example
```
Model: Arsenal 52% to win
Odds: 2.10
EV = (0.52 × 2.10) - 1 = 0.092 = 9.2% ✅ Recommend

Model: Chelsea 48% to win
Odds: 2.00
EV = (0.48 × 2.00) - 1 = -0.04 = -4% ❌ No bet
```

### Rationale
- **Conservative threshold** - filters noise, higher quality picks
- **Covers uncertainty** - model probability estimates aren't perfect
- **Reduces bet volume** - focus on best opportunities
- **Better for investor demo** - shows discipline and strategy
- **Buffer for vig** - typical bookmaker margin is 2-5%

### Consequences
**Positive:**
- Higher quality picks (fewer but better)
- Less capital at risk
- Easier to track and manage
- Demonstrates value-betting discipline

**Negative:**
- May miss some profitable bets (2-3% EV)
- Lower volume = slower bankroll growth
- Need accurate probability calibration

### Configurable Thresholds
Allow adjusting threshold based on:
- Model confidence level
- Closing line value (CLV) tracking
- Historical performance
- Risk appetite

---

## ADR-008: Manual Bet Placement for MVP

**Date:** 2026-01-12
**Status:** Accepted

### Context
MVP timeline is 8 weeks. Automated bet placement APIs are complex and time-consuming.

### Decision
MVP will generate recommendations but require manual bet placement and result entry.

### Rationale
- **Faster MVP development** - saves 2-3 weeks of API integration work
- **Lower risk** - human verification before real money bets
- **Regulatory complexity** - automated betting may have legal restrictions
- **Cost savings** - no betting API subscriptions needed for MVP
- **Better for demo** - can show picks without needing real betting accounts

### MVP Workflow
1. System generates weekly picks
2. Dashboard displays recommendations with stake sizes
3. User manually places bets at bookmaker
4. User records bet placement in system
5. User settles bets after match completion

### Consequences
**Positive:**
- MVP can launch without betting API integrations
- Manual review catches potential model errors
- No regulatory concerns for demo period
- Can test across multiple bookmakers easily

**Negative:**
- Manual work required (not scalable)
- Risk of execution errors
- Can't capture exact odds (line movement)
- Time-sensitive (odds change)

### Future Automation (Phase 2+)
- Integrate with Betfair API
- Automated bet placement
- Real-time odds monitoring
- Auto-settlement via API

---

## ADR-009: Weekly Betting Cycle for MVP

**Date:** 2026-01-12
**Status:** Accepted

### Context
Premier League has weekend fixtures primarily. Need to establish betting rhythm for MVP.

### Decision
Generate picks on **weekly cycle** (Monday-Sunday), published Thursday/Friday.

### Weekly Cycle
```
Monday-Wednesday: Data updates, model retraining if needed
Thursday: Generate predictions for weekend fixtures
Thursday PM: Calculate EV, filter picks, publish recommendations
Friday-Saturday: Users place bets manually
Weekend: Matches played
Sunday-Monday: Settle bets, update performance metrics
```

### Rationale
- **Aligns with fixture schedule** - most matches on weekends
- **Time for analysis** - 1-2 days to review picks before placement
- **Captures best odds** - Thursday odds typically better than day-of
- **Manageable cadence** - weekly reviews for investors
- **Allows iteration** - can adjust strategy week-to-week

### Consequences
**Positive:**
- Clear rhythm for users and investors
- Time to validate picks before betting
- Weekly performance reporting natural
- Reduces decision fatigue

**Negative:**
- May miss mid-week fixtures (smaller volume)
- Line movement between Thursday and Saturday
- Can't adjust for late team news

### MVP Target
- 8-12 weeks of weekly picks
- ~20-30 picks per week on average (multi-market)
- ~160-250 total bets for statistical significance

---

## ADR-010: Multi-Market Strategy with Incremental Rollout

**Date:** 2026-01-12
**Status:** Accepted

### Context

Bookmakers offer numerous betting markets beyond just 1X2 (Home/Draw/Away):
- Over/Under goals (various lines: 0.5, 1.5, 2.5, 3.5, 4.5, etc.)
- Both Teams to Score (BTTS)
- Double Chance (Home or Draw, etc.)
- Handicaps
- First Goal
- Correct Score
- Half-time markets
- And many more exotic markets

**The original MVP plan only focused on 1X2**, which would:
- Miss value opportunities in other markets
- Limit weekly picks to ~10 bets
- Not leverage the full odds data available
- Compete only in the most efficient (lowest edge) market

A smarter system should evaluate ALL available markets and recommend the highest EV opportunity, regardless of market type.

### Decision

Implement a **multi-market strategy with incremental rollout**:

**Week 3:** Launch with 1X2 model (proves pipeline works)
**Week 4:** Add Over/Under 2.5 Goals model
**Week 5:** Add BTTS model
**Week 6:** Implement Smart Market Selector

**Smart Market Selector Logic:**
```
For each fixture:
  1. Get predictions from all active models (1X2, O/U, BTTS, etc.)
  2. Fetch odds for all markets from bookmakers
  3. Calculate EV for EVERY market/outcome combination
  4. Rank all opportunities by EV
  5. Filter by minimum threshold (3% EV)
  6. Return top recommendations (can be from different markets)
```

**Example Output:**
```
Fixture: Arsenal vs Liverpool

Market Analysis:
- Home Win (1X2): 45% prob, 2.10 odds → EV = -5.5% ❌
- Draw (1X2): 28% prob, 3.50 odds → EV = -2% ❌
- Away Win (1X2): 27% prob, 3.20 odds → EV = -13.6% ❌
- Over 2.5: 65% prob, 1.80 odds → EV = +17% ✅ BEST
- BTTS Yes: 70% prob, 1.65 odds → EV = +15.5% ✅

Recommendation: Over 2.5 Goals @ 1.80 (stake: $142)
Alternative: BTTS Yes @ 1.65 (stake: $128)
```

### Rationale

**Maximizes Value:**
- Not constrained to single market type
- Finds best EV across entire opportunity set
- More betting opportunities (20-30 picks/week vs 10)

**Better Risk-Adjusted Returns:**
- Diversification across uncorrelated markets
- BTTS is independent of match outcome
- Totals markets often less efficient than 1X2

**Competitive Advantage:**
- Most casual systems only predict 1X2
- Professional bettors already use multi-market approach
- Demonstrates sophistication to investors

**Incremental Approach Reduces Risk:**
- Week 3: Prove ML pipeline works with single model
- Week 4-5: Add models one at a time, validate each
- Week 6: Integrate when confident all models work
- Can rollback if any model underperforms

**Market Selection by Predictability:**
- O/U 2.5: Easier to predict than 1X2 (60-65% accuracy)
- BTTS: Strong historical patterns, good features
- Starting with most predictable markets

### Consequences

**Positive:**
- Higher overall ROI (target: 6-10% vs 3-5% single market)
- More weekly picks (20-30 vs 10)
- Better diversification
- Market-specific performance tracking
- Can disable underperforming markets
- More attractive to investors (comprehensive coverage)

**Negative:**
- More complex to build (3 models instead of 1)
- More data to fetch (odds for multiple markets)
- More complex testing and validation
- Need separate backtests per market
- MVP timeline extends from 6 weeks to 8 weeks
- Higher computational requirements

**Risk Mitigation:**
- Incremental rollout limits downside
- Can launch with just 1X2 if needed
- Market-specific Kelly fractions (conservative on new markets)
- Max 40% of picks from any single market (diversification rule)
- Independent backtesting per market before enabling

### Implementation Notes

**Database:**
- `odds.market_type` already supports multiple markets
- No schema changes needed

**ML Service:**
```python
ml-service/app/models/
├── h2h_model.py          # 1X2 (Week 3)
├── totals_model.py       # O/U 2.5 (Week 4)
├── btts_model.py         # BTTS (Week 5)
└── model_manager.py      # Orchestrates all models
```

**Betting Engine:**
```go
// backend/internal/services/market_selector.go
func SelectBestMarkets(fixture Fixture, minEV float64) []Recommendation {
    markets := []string{"h2h", "totals_2_5", "btts"}
    recommendations := []Recommendation{}

    for _, market := range markets {
        predictions := mlService.Predict(fixture, market)
        odds := getOdds(fixture, market)

        for outcome, prob := range predictions {
            ev := calculateEV(prob, odds[outcome])
            if ev >= minEV {
                recommendations = append(recommendations, Recommendation{
                    Market: market,
                    Outcome: outcome,
                    EV: ev,
                    // ... other fields
                })
            }
        }
    }

    sort.Slice(recommendations, func(i, j int) bool {
        return recommendations[i].EV > recommendations[j].EV
    })

    return recommendations
}
```

**Performance Tracking:**
- Track ROI, win rate, CLV separately per market
- Overall performance = weighted average
- Can disable markets with poor CLV

### Future Expansion

Post-MVP markets (see `docs/market-expansion-roadmap.md`):
- Double Chance (Month 3)
- Other totals: O/U 1.5, 3.5, 4.5 (Month 3)
- Half-time markets (Month 3-4)
- Handicaps (Month 4-5)
- Correct Score (Month 5-6)
- Player props (Month 7+)

Target: 15-20 markets by Month 6

### Alternatives Considered

1. **Single Market (1X2 only)** - Simpler but leaves money on table
2. **All Markets at Once** - Higher risk, harder to validate
3. **Market-Agnostic Model** - Single model predicts all markets (technically difficult)

---

## ADR-011: Simple Accumulators (Parlays) in MVP

**Date:** 2026-01-12
**Status:** Accepted

### Context

Accumulators (parlays) are a popular betting format where multiple selections are combined into a single bet. All legs must win for the accumulator to pay out.

**Benefits:**
- Higher potential returns (odds multiply)
- Popular with casual bettors (investor appeal)
- Can create value from multiple +EV opportunities
- Reduces number of individual transactions

**Risks:**
- Higher variance (all legs must win)
- Correlation risk (events may not be independent)
- Harder to achieve positive EV
- One loss kills entire accumulator

**Example:**
```
3-leg accumulator:
Leg 1: Arsenal Over 2.5 @ 1.80 (65% prob)
Leg 2: Brighton BTTS @ 1.70 (68% prob)
Leg 3: Liverpool Away @ 2.20 (55% prob)

Combined odds: 6.73
Combined probability: 24.3%
EV = (0.243 × 6.73) - 1 = 63.5%
```

### Decision

Add **simple accumulator building** to MVP (Week 7):

**Scope:**
- 2-3 leg accumulators only (no 4+ leg parlays)
- 2-3 accumulators recommended per week
- Maximum 20% of weekly stake allocation
- Conservative Kelly sizing (1/8 instead of 1/4)
- Basic correlation detection

**Correlation Rules:**
- ❌ Can't combine picks from same fixture
- ❌ Can't combine picks involving same team (optional, configurable)
- ✅ Can combine uncorrelated markets (Arsenal Over 2.5 + Brighton BTTS + Liverpool Away)

**Accumulator EV Calculation:**
```
Combined Probability = P(Leg1) × P(Leg2) × P(Leg3)
Combined Odds = Odds1 × Odds2 × Odds3
EV = (Combined Probability × Combined Odds) - 1

Minimum EV for accumulators: 5% (vs 3% for singles)
```

**Stake Sizing:**
- Use 1/8 Kelly instead of 1/4 (more conservative)
- Cap total accumulator stake at 20% of weekly budget
- Smaller stakes reflect higher variance

### Rationale

**Why Include in MVP:**
- **Investor appeal** - Shows sophistication, demonstrates multiple betting formats
- **Higher returns** - Can showcase impressive wins (6-10x returns)
- **Popular format** - Most bettors use accumulators, good for demos
- **Natural fit** - We already have multiple high-EV picks, makes sense to combine
- **Not too complex** - 1 week of dev time, manageable risk

**Why Keep Simple:**
- **Lower risk** - 2-3 legs max limits variance
- **Conservative stakes** - 20% max allocation protects bankroll
- **Basic correlation** - Simple rules are easier to explain and validate
- **Higher EV threshold** - 5% minimum filters out marginal combinations
- **Incremental approach** - Can enhance in Phase 2 if successful

**Why Not Full Portfolio System:**
- **Too complex** - Multi-ticket optimization adds 2-3 weeks
- **Need more data** - Portfolio theory requires statistical validation
- **Harder to explain** - Investors may not understand portfolio optimization
- **Can add later** - Phase 2 (Months 3-4) is appropriate timeline

### Consequences

**Positive:**
- More attractive to investors (demonstrates comprehensive strategy)
- Potential for higher returns (6-10x payouts)
- Shows smart combination logic (correlation detection)
- Popular betting format (good for marketing)
- Only adds 1 week to timeline (Week 7)
- Clear differentiation (singles + accumulators)

**Negative:**
- Higher variance (accumulators are riskier)
- More complex to track (separate tables needed)
- Correlation detection may miss edge cases
- Requires separate Kelly sizing logic
- Testing is more complex (need multiple wins)
- Could lose all accumulator stake on single fixture upset

**Risk Mitigation:**
- **Conservative sizing**: 1/8 Kelly + 20% max allocation
- **Higher EV threshold**: 5% minimum (vs 3% singles)
- **Leg limits**: 2-3 legs max (no crazy 10-leg parlays)
- **Correlation detection**: Automated filtering
- **Independent tracking**: Separate performance metrics
- **Can disable**: If accumulators underperform, remove from recommendations

### Implementation Notes

**Database:**
```sql
CREATE TABLE accumulators (
    id SERIAL PRIMARY KEY,
    num_legs INTEGER NOT NULL,
    stake DECIMAL(10, 2) NOT NULL,
    combined_odds DECIMAL(10, 2) NOT NULL,
    combined_probability DECIMAL(5, 4) NOT NULL,
    expected_value DECIMAL(10, 4) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    ...
);

CREATE TABLE accumulator_legs (
    id SERIAL PRIMARY KEY,
    accumulator_id INTEGER REFERENCES accumulators(id),
    bet_id INTEGER REFERENCES bets(id),
    leg_order INTEGER NOT NULL,
    ...
);
```

**Accumulator Builder Algorithm:**
```go
func BuildAccumulators(picks []Pick, maxLegs int, minEV float64) []Accumulator {
    // 1. Filter to high-quality picks (singles with EV > 3%)
    highQuality := filterByEV(picks, 0.03)

    // 2. Remove correlated picks
    uncorrelated := removeCorrelatedPicks(highQuality)

    // 3. Generate combinations (2-3 legs)
    combinations := generateCombinations(uncorrelated, maxLegs)

    // 4. Calculate accumulator EV for each combo
    accumulators := []Accumulator{}
    for _, combo := range combinations {
        combinedProb := multiplyProbabilities(combo)
        combinedOdds := multiplyOdds(combo)
        ev := (combinedProb * combinedOdds) - 1

        if ev >= minEV { // 5% minimum
            stake := calculateKellyStake(combinedProb, combinedOdds, 0.125) // 1/8 Kelly
            accumulators = append(accumulators, createAccumulator(combo, stake, ev))
        }
    }

    // 5. Sort by EV, take top 2-3
    sort.Slice(accumulators, func(i, j int) bool {
        return accumulators[i].EV > accumulators[j].EV
    })

    return accumulators[:min(3, len(accumulators))]
}
```

**Correlation Detection:**
```go
func removeCorrelatedPicks(picks []Pick) []Pick {
    used := make(map[int]bool) // Track used fixture IDs
    uncorrelated := []Pick{}

    for _, pick := range picks {
        if used[pick.FixtureID] {
            continue // Skip if fixture already used
        }

        // Optional: Check team correlation
        // if hasTeamOverlap(pick, uncorrelated) { continue }

        uncorrelated = append(uncorrelated, pick)
        used[pick.FixtureID] = true
    }

    return uncorrelated
}
```

**Settlement Logic:**
```go
func SettleAccumulator(accumulator Accumulator) {
    // Check all legs
    allWon := true
    for _, leg := range accumulator.Legs {
        if leg.Status != "won" {
            allWon = false
            break
        }
    }

    if allWon {
        accumulator.Status = "won"
        accumulator.Payout = accumulator.Stake * accumulator.CombinedOdds
        accumulator.ProfitLoss = accumulator.Payout - accumulator.Stake
    } else {
        accumulator.Status = "lost"
        accumulator.Payout = 0
        accumulator.ProfitLoss = -accumulator.Stake
    }
}
```

### Performance Tracking

Track separately from singles:
- **Singles ROI**: Individual bet performance
- **Accumulator ROI**: Parlay performance
- **Overall ROI**: Weighted average
- **Win rate**: Accumulators typically 10-30% (vs 50-60% singles)
- **Average return**: Higher when won (6-10x vs 2-3x)

**Example Metrics Dashboard:**
```
Singles (15-20 per week):
- ROI: 8.2%
- Win Rate: 56%
- Avg Odds: 2.15

Accumulators (2-3 per week):
- ROI: 12.5%
- Win Rate: 22%
- Avg Odds: 7.2
- Avg Legs: 2.7

Overall:
- ROI: 9.1%
- Total Profit: $1,450
```

### Testing Requirements

Before enabling accumulators:
1. Backtest on historical data
2. Verify correlation detection works
3. Validate EV calculations are accurate
4. Test settlement logic (all legs must win)
5. Ensure Kelly sizing is conservative
6. Paper trading for 2-3 weeks

### Future Enhancements (Phase 2)

**Month 3-4:**
- 4-5 leg accumulators (if 2-3 leg performs well)
- Advanced correlation detection (statistical correlation matrix)
- Same-game parlays (combining correlated markets intelligently)
- Accumulator insurance (void if 1 leg loses)
- System bets (multiple accumulators from same picks)

**Not in Scope for MVP:**
- Multi-ticket portfolio optimization
- Dynamic leg count based on EV
- Live/in-play accumulator building
- Accumulator cash-out logic

### Alternatives Considered

1. **No Accumulators** - Simpler, but less attractive to investors
2. **Full Multi-Ticket Portfolio** - Better diversification, but adds 2-3 weeks
3. **Same-Game Parlays** - Higher risk due to correlation, too complex for MVP
4. **Accumulator Only (No Singles)** - Too risky, variance too high

### Acceptance Criteria

- [ ] Can generate 2-3 accumulators per week
- [ ] All legs are from different fixtures
- [ ] Combined EV ≥ 5%
- [ ] Uses 1/8 Kelly sizing
- [ ] Max 20% of weekly stake allocated
- [ ] Correlation detection prevents same-fixture combinations
- [ ] Settlement logic works correctly (all legs must win)
- [ ] Dashboard displays accumulators with legs breakdown
- [ ] Performance tracking separate from singles

