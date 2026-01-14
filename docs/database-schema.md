# Database Schema Design

## Overview
PostgreSQL database schema for OddsIQ MVP, designed to store 3 seasons of Premier League data, odds, predictions, and betting results.

## Tables

### 1. teams
Stores Premier League team information.

```sql
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    api_football_id INTEGER UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10),
    logo_url TEXT,
    founded INTEGER,
    venue_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_teams_api_football_id ON teams(api_football_id);
```

### 2. fixtures
Match fixtures from Premier League.

```sql
CREATE TABLE fixtures (
    id SERIAL PRIMARY KEY,
    api_football_id INTEGER UNIQUE NOT NULL,
    season INTEGER NOT NULL,
    round VARCHAR(50),
    match_date TIMESTAMP NOT NULL,
    home_team_id INTEGER REFERENCES teams(id),
    away_team_id INTEGER REFERENCES teams(id),
    home_score INTEGER,
    away_score INTEGER,
    status VARCHAR(20), -- 'scheduled', 'live', 'finished', 'postponed'
    venue VARCHAR(100),
    referee VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_fixtures_match_date ON fixtures(match_date);
CREATE INDEX idx_fixtures_season ON fixtures(season);
CREATE INDEX idx_fixtures_status ON fixtures(status);
CREATE INDEX idx_fixtures_home_team ON fixtures(home_team_id);
CREATE INDEX idx_fixtures_away_team ON fixtures(away_team_id);
```

### 3. odds
Bookmaker odds for fixtures.

```sql
CREATE TABLE odds (
    id SERIAL PRIMARY KEY,
    fixture_id INTEGER REFERENCES fixtures(id) ON DELETE CASCADE,
    bookmaker VARCHAR(50) NOT NULL,
    market_type VARCHAR(50) NOT NULL, -- 'h2h', 'totals', 'spreads'
    outcome VARCHAR(20), -- 'home', 'draw', 'away', 'over', 'under'
    odds_value DECIMAL(10, 2) NOT NULL,
    recorded_at TIMESTAMP NOT NULL,
    is_closing_line BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_odds_fixture ON odds(fixture_id);
CREATE INDEX idx_odds_bookmaker ON odds(bookmaker);
CREATE INDEX idx_odds_recorded_at ON odds(recorded_at);
CREATE INDEX idx_odds_closing_line ON odds(is_closing_line) WHERE is_closing_line = TRUE;
```

### 4. team_stats
Aggregated team statistics for feature engineering.

```sql
CREATE TABLE team_stats (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id),
    season INTEGER NOT NULL,
    match_date DATE NOT NULL,
    games_played INTEGER DEFAULT 0,
    wins INTEGER DEFAULT 0,
    draws INTEGER DEFAULT 0,
    losses INTEGER DEFAULT 0,
    goals_for INTEGER DEFAULT 0,
    goals_against INTEGER DEFAULT 0,
    points INTEGER DEFAULT 0,
    position INTEGER,
    form_last_5 VARCHAR(5), -- e.g., 'WWDLW'
    home_wins INTEGER DEFAULT 0,
    home_draws INTEGER DEFAULT 0,
    home_losses INTEGER DEFAULT 0,
    away_wins INTEGER DEFAULT 0,
    away_draws INTEGER DEFAULT 0,
    away_losses INTEGER DEFAULT 0,
    xg_for DECIMAL(5, 2),
    xg_against DECIMAL(5, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(team_id, season, match_date)
);

CREATE INDEX idx_team_stats_team_season ON team_stats(team_id, season);
CREATE INDEX idx_team_stats_match_date ON team_stats(match_date);
```

### 5. predictions
Model predictions for fixtures.

```sql
CREATE TABLE predictions (
    id SERIAL PRIMARY KEY,
    fixture_id INTEGER REFERENCES fixtures(id) ON DELETE CASCADE,
    model_version VARCHAR(20) NOT NULL,
    home_win_prob DECIMAL(5, 4) NOT NULL,
    draw_prob DECIMAL(5, 4) NOT NULL,
    away_win_prob DECIMAL(5, 4) NOT NULL,
    predicted_outcome VARCHAR(10), -- 'home', 'draw', 'away'
    confidence_score DECIMAL(5, 4),
    features JSONB, -- Store feature values used
    predicted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_predictions_fixture ON predictions(fixture_id);
CREATE INDEX idx_predictions_model_version ON predictions(model_version);
CREATE INDEX idx_predictions_predicted_at ON predictions(predicted_at);
```

### 6. bets
Tracking of actual bets placed.

```sql
CREATE TABLE bets (
    id SERIAL PRIMARY KEY,
    fixture_id INTEGER REFERENCES fixtures(id),
    prediction_id INTEGER REFERENCES predictions(id),
    bet_type VARCHAR(50) NOT NULL, -- 'h2h_home', 'h2h_draw', 'h2h_away'
    stake DECIMAL(10, 2) NOT NULL,
    odds DECIMAL(10, 2) NOT NULL,
    expected_value DECIMAL(10, 4) NOT NULL, -- EV percentage
    bookmaker VARCHAR(50),
    placed_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'won', 'lost', 'void'
    payout DECIMAL(10, 2),
    profit_loss DECIMAL(10, 2),
    settled_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bets_fixture ON bets(fixture_id);
CREATE INDEX idx_bets_status ON bets(status);
CREATE INDEX idx_bets_placed_at ON bets(placed_at);
```

### 7. bankroll
Bankroll tracking over time.

```sql
CREATE TABLE bankroll (
    id SERIAL PRIMARY KEY,
    balance DECIMAL(10, 2) NOT NULL,
    total_staked DECIMAL(10, 2) DEFAULT 0,
    total_returned DECIMAL(10, 2) DEFAULT 0,
    total_profit_loss DECIMAL(10, 2) DEFAULT 0,
    roi_percentage DECIMAL(10, 4),
    num_bets INTEGER DEFAULT 0,
    num_wins INTEGER DEFAULT 0,
    num_losses INTEGER DEFAULT 0,
    win_rate DECIMAL(5, 4),
    recorded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bankroll_recorded_at ON bankroll(recorded_at);
```

### 8. model_performance
Track model performance metrics over time.

```sql
CREATE TABLE model_performance (
    id SERIAL PRIMARY KEY,
    model_version VARCHAR(20) NOT NULL,
    evaluation_date DATE NOT NULL,
    accuracy DECIMAL(5, 4),
    precision DECIMAL(5, 4),
    recall DECIMAL(5, 4),
    f1_score DECIMAL(5, 4),
    brier_score DECIMAL(5, 4),
    log_loss DECIMAL(10, 6),
    roc_auc DECIMAL(5, 4),
    num_predictions INTEGER,
    backtest_roi DECIMAL(10, 4),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_model_performance_version ON model_performance(model_version);
CREATE INDEX idx_model_performance_date ON model_performance(evaluation_date);
```

### 9. accumulators
Accumulator (parlay) bet tracking.

```sql
CREATE TABLE accumulators (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    num_legs INTEGER NOT NULL,
    stake DECIMAL(10, 2) NOT NULL,
    combined_odds DECIMAL(10, 2) NOT NULL,
    combined_probability DECIMAL(5, 4) NOT NULL,
    expected_value DECIMAL(10, 4) NOT NULL,
    potential_payout DECIMAL(10, 2),
    placed_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    actual_payout DECIMAL(10, 2),
    profit_loss DECIMAL(10, 2),
    settled_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_accumulators_status ON accumulators(status);
CREATE INDEX idx_accumulators_placed_at ON accumulators(placed_at);
```

### 10. accumulator_legs
Individual legs within an accumulator.

```sql
CREATE TABLE accumulator_legs (
    id SERIAL PRIMARY KEY,
    accumulator_id INTEGER REFERENCES accumulators(id) ON DELETE CASCADE,
    bet_id INTEGER REFERENCES bets(id) ON DELETE CASCADE,
    leg_order INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(accumulator_id, bet_id)
);

CREATE INDEX idx_accumulator_legs_accumulator ON accumulator_legs(accumulator_id);
CREATE INDEX idx_accumulator_legs_bet ON accumulator_legs(bet_id);
```

## Data Volume Estimates

### 3 Seasons of Premier League Data
- **Teams**: 20 teams
- **Fixtures**: ~380 matches per season × 3 = ~1,140 fixtures
- **Odds**: ~1,140 fixtures × 5 bookmakers × 3 outcomes × 10 updates = ~171,000 odds records
- **Team Stats**: 20 teams × 3 seasons × 38 weeks = ~2,280 stat snapshots
- **Predictions**: ~1,140 predictions (one per fixture)
- **Bets**: 8 weeks × ~10 bets/week = ~80 bets (initial MVP period)

### Storage Estimate
- Total initial data: ~200 MB
- With indexes: ~350 MB
- Growth rate: ~5-10 MB/week during active betting

## Indexes Strategy

**Query patterns prioritized:**
1. Get fixtures by date range (for upcoming matches)
2. Get latest odds for a fixture
3. Get team stats at a specific point in time
4. Get predictions for upcoming fixtures
5. Calculate performance metrics over time periods

## Sample Queries

### Get upcoming fixtures with latest odds
```sql
SELECT
    f.id,
    f.match_date,
    ht.name as home_team,
    at.name as away_team,
    o.odds_value as home_odds,
    o2.odds_value as away_odds
FROM fixtures f
JOIN teams ht ON f.home_team_id = ht.id
JOIN teams at ON f.away_team_id = at.id
LEFT JOIN LATERAL (
    SELECT odds_value
    FROM odds
    WHERE fixture_id = f.id
      AND outcome = 'home'
      AND bookmaker = 'bet365'
    ORDER BY recorded_at DESC
    LIMIT 1
) o ON true
LEFT JOIN LATERAL (
    SELECT odds_value
    FROM odds
    WHERE fixture_id = f.id
      AND outcome = 'away'
      AND bookmaker = 'bet365'
    ORDER BY recorded_at DESC
    LIMIT 1
) o2 ON true
WHERE f.status = 'scheduled'
  AND f.match_date > NOW()
ORDER BY f.match_date;
```

### Calculate current performance metrics
```sql
SELECT
    COUNT(*) as total_bets,
    SUM(CASE WHEN status = 'won' THEN 1 ELSE 0 END) as wins,
    SUM(CASE WHEN status = 'lost' THEN 1 ELSE 0 END) as losses,
    ROUND(SUM(CASE WHEN status = 'won' THEN 1 ELSE 0 END)::DECIMAL /
          NULLIF(COUNT(*), 0) * 100, 2) as win_rate,
    SUM(stake) as total_staked,
    SUM(COALESCE(profit_loss, 0)) as total_pl,
    ROUND(SUM(COALESCE(profit_loss, 0)) / NULLIF(SUM(stake), 0) * 100, 2) as roi
FROM bets
WHERE status IN ('won', 'lost');
```
