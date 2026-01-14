# Market Expansion Roadmap

## Overview

OddsIQ will use an incremental multi-market approach, starting with the most popular markets and gradually expanding to cover more exotic betting options.

## Market Priority Framework

Markets prioritized by:
1. **Volume** - How often bookmakers offer these markets
2. **Liquidity** - Bet volume and odds availability
3. **Predictability** - How well ML models can predict these outcomes
4. **EV Potential** - Historical edge opportunities

## Phase-by-Phase Market Rollout

### 🎯 Phase 1: Core Market (Week 3)
**Market: 1X2 (Home/Draw/Away)**

- Most fundamental market
- Proves ML pipeline works
- 3-way classification (Home/Draw/Away)

**Model:** `h2h_model.py`
**Target Accuracy:** 55-60%
**Volume:** ~100% of fixtures
**Odds Range:** 1.50 - 15.00

---

### 🎯 Phase 2: Goals Market (Week 4)
**Market: Over/Under 2.5 Goals**

- Most popular totals market
- High liquidity, tight odds
- Binary classification (Over/Under)
- Strong correlation with team stats

**Model:** `totals_model.py`
**Target Accuracy:** 60-65%
**Volume:** ~95% of fixtures
**Odds Range:** 1.40 - 3.00

**Features to add:**
- Average goals scored (home/away)
- Average goals conceded (home/away)
- Recent scoring form (last 5 games)
- Head-to-head goal averages

---

### 🎯 Phase 3: BTTS Market (Week 5)
**Market: Both Teams to Score (Yes/No)**

- Very popular with bettors
- Independent of match outcome
- Binary classification

**Model:** `btts_model.py`
**Target Accuracy:** 58-62%
**Volume:** ~90% of fixtures
**Odds Range:** 1.50 - 2.50

**Features to add:**
- Clean sheet percentage
- Teams scoring percentage
- Defensive strength metrics
- Offensive strength metrics

---

### 🎯 Phase 4: Market Selector (Week 6)
**Feature: Intelligent Market Selection**

- Evaluate EV across all 3 markets
- Select highest EV opportunity
- Can recommend multiple markets if EV > threshold

**Example Output:**
```
Fixture: Arsenal vs Liverpool

Market Analysis:
1. Home Win (1X2): EV = -2.5% ❌
2. Draw (1X2): EV = +1.8% ❌ (below 3% threshold)
3. Over 2.5 Goals: EV = +8.3% ✅ RECOMMENDED
4. BTTS Yes: EV = +6.1% ✅ RECOMMENDED

Recommendations:
- PRIMARY: Over 2.5 @ 1.85 (Stake: $127)
- SECONDARY: BTTS Yes @ 1.70 (Stake: $95)
```

---

### 🎯 Phase 5: Double Chance (Week 7-8)
**Markets: Home or Draw / Home or Away / Draw or Away**

- Lower risk, lower odds
- Good for high-confidence scenarios
- 3 binary classifications

**Model:** `double_chance_model.py` or reuse h2h_model
**Target Accuracy:** 65-70% (easier to predict)
**Volume:** ~80% of fixtures
**Odds Range:** 1.10 - 1.80

---

## Future Market Expansion (Post-MVP)

### Phase 6: Additional Totals (Month 3)
- Over/Under 1.5 Goals
- Over/Under 3.5 Goals
- Over/Under 4.5 Goals
- First Half Over/Under 0.5, 1.5

### Phase 7: Half-Time Markets (Month 3-4)
- Half-Time Result (1X2)
- Half-Time/Full-Time (9 outcomes)
- First Half Goals Over/Under

### Phase 8: Team-Specific Totals (Month 4)
- Home Team Over/Under 1.5
- Away Team Over/Under 0.5
- Exact Number of Goals (0, 1, 2, 3, 4+)

### Phase 9: Handicap Markets (Month 4-5)
- Asian Handicap (-1.5, -1, -0.5, 0, +0.5, +1, +1.5)
- European Handicap (0:1, 0:2, 1:0, 2:0)

### Phase 10: Goal Timing (Month 5)
- First Goal (Home/None/Away)
- Last Goal
- Early Goals (Over/Under 1.5 in first 25 mins)
- Late Goals

### Phase 11: Advanced Markets (Month 6+)
- Correct Score (top 6-10 likely scores)
- Winning Margin
- Clean Sheet (Yes/No for each team)
- To Win to Nil
- Corners Over/Under
- Cards Over/Under

### Phase 12: Player Props (Month 7+)
Requires additional data sources:
- Anytime Goalscorer
- First Goalscorer
- Player Shots/Cards/Assists

---

## Market Coverage Goals

### MVP (Week 8)
- **3 markets** (1X2, O/U 2.5, BTTS)
- **~85% coverage** of all available bets
- **Multiple picks per fixture** if multiple EV opportunities

### Phase 2 (Month 3-4)
- **8-10 markets**
- **~95% coverage**
- Including half-time, handicaps, team totals

### Phase 3 (Month 6+)
- **15-20 markets**
- **~99% coverage**
- All major markets covered

---

## Market-Specific Model Strategy

### Simple Binary Markets
**Markets:** Over/Under, BTTS, Double Chance
**Approach:** Single XGBoost binary classifier per market
**Training:** Faster, simpler feature engineering

### Multi-Class Markets
**Markets:** 1X2, Correct Score, Half-Time/Full-Time
**Approach:** XGBoost multi-class classifier
**Training:** More complex, requires probability calibration

### Derivative Markets
**Markets:** Double Chance, Team Totals
**Approach:** Can derive from existing models
**Example:** P(Home or Draw) = P(Home) + P(Draw) from 1X2 model

---

## Technical Implementation

### Database Schema
```sql
-- Market types supported
CREATE TYPE market_type AS ENUM (
    -- Core markets
    'h2h',              -- Home/Draw/Away
    'totals_2_5',       -- Over/Under 2.5 goals
    'btts',             -- Both teams to score
    'double_chance',    -- Home or Draw, etc.

    -- Extended totals
    'totals_1_5',
    'totals_3_5',
    'totals_4_5',

    -- Half markets
    'ht_result',        -- Half-time result
    'ht_ft',            -- Half-time/Full-time
    'ht_totals_0_5',
    'ht_totals_1_5',

    -- Team totals
    'home_totals_0_5',
    'home_totals_1_5',
    'away_totals_0_5',
    'away_totals_1_5',

    -- Handicaps
    'asian_handicap',
    'european_handicap',

    -- Goal markets
    'first_goal',
    'last_goal',
    'early_goals',

    -- Advanced
    'correct_score',
    'winning_margin',
    'clean_sheet',

    -- Future
    'corners',
    'cards',
    'player_props'
);
```

### ML Service Structure
```
ml-service/app/models/
├── base_model.py           # Abstract base class
├── h2h_model.py            # 1X2 predictions
├── totals_model.py         # Over/Under predictions
├── btts_model.py           # BTTS predictions
├── double_chance_model.py  # Double Chance (derived)
├── half_time_model.py      # Half-time markets
├── handicap_model.py       # Handicap markets
├── correct_score_model.py  # Correct score
└── model_manager.py        # Loads and manages all models

ml-service/app/features/
├── base_features.py        # Shared features
├── goals_features.py       # Goal-specific features
├── defensive_features.py   # Defensive metrics
├── timing_features.py      # Goal timing patterns
└── advanced_features.py    # xG, possession, shots
```

### API Endpoint
```python
POST /api/predict/all-markets
{
    "fixture_id": 123,
    "markets": ["h2h", "totals_2_5", "btts", "double_chance"]
}

Response:
{
    "fixture_id": 123,
    "predictions": {
        "h2h": {
            "home": 0.45,
            "draw": 0.28,
            "away": 0.27
        },
        "totals_2_5": {
            "over": 0.65,
            "under": 0.35
        },
        "btts": {
            "yes": 0.68,
            "no": 0.32
        },
        "double_chance": {
            "home_or_draw": 0.73,
            "home_or_away": 0.72,
            "draw_or_away": 0.55
        }
    },
    "model_versions": {
        "h2h": "v1.0",
        "totals_2_5": "v1.0",
        "btts": "v1.0",
        "double_chance": "derived"
    }
}
```

---

## Market Selection Algorithm

```python
def select_best_markets(fixture_id, min_ev=0.03):
    """
    Evaluate all available markets and return ranked recommendations
    """
    markets = get_available_markets(fixture_id)
    recommendations = []

    for market in markets:
        # Get predictions from appropriate model
        predictions = predict_market(fixture_id, market)

        # Get latest odds for this market
        odds = get_latest_odds(fixture_id, market)

        # Calculate EV for each outcome
        for outcome, prob in predictions.items():
            ev = calculate_ev(prob, odds[outcome])

            if ev >= min_ev:
                recommendations.append({
                    'market': market,
                    'outcome': outcome,
                    'probability': prob,
                    'odds': odds[outcome],
                    'ev': ev,
                    'ev_percentage': ev * 100,
                    'stake': calculate_kelly_stake(prob, odds[outcome])
                })

    # Sort by EV (highest first)
    recommendations.sort(key=lambda x: x['ev'], reverse=True)

    return recommendations
```

---

## Success Metrics by Phase

### Week 3 (1X2 Only)
- Model accuracy: 55-60%
- Backtest ROI: +3-5%
- ~10 picks per week

### Week 6 (3 Markets)
- Combined accuracy: 58-63%
- Backtest ROI: +5-8%
- ~20-25 picks per week
- Market distribution: 40% 1X2, 35% Totals, 25% BTTS

### Week 8 (MVP Complete)
- 3-4 markets active
- Backtest ROI: +6-10%
- ~25-30 picks per week
- Smart market selector operational

### Month 3-4 (Extended)
- 8-10 markets
- Backtest ROI: +8-12%
- ~40-50 picks per week

---

## Market Risk Management

### Diversification Rules
- Max 40% of weekly picks from single market type
- If 1X2 has poor week, other markets provide stability
- Uncorrelated markets (e.g., BTTS is independent of 1X2)

### Market-Specific Limits
- Higher stakes for high-accuracy markets (Totals, BTTS)
- Lower stakes for lower-accuracy markets (1X2)
- Different Kelly fractions per market type

### Closing Line Value (CLV) Tracking
- Track CLV separately per market
- Some markets may have better CLV than others
- Adjust strategy based on market-specific CLV performance

---

## Notes

- The Odds API provides data for most major markets
- API-Football provides match stats for feature engineering
- Some exotic markets may have inconsistent odds availability
- Player props require additional data sources (not in MVP scope)
- Market expansion driven by data: add markets with consistent odds availability

**This roadmap is flexible** - we can accelerate certain markets or skip others based on:
- Data availability
- Model performance
- Odds availability
- Investor feedback
- Profitability per market
