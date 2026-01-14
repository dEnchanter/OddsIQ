# API Specification

## Overview
RESTful API specifications for OddsIQ system services.

## Go Backend API (Port 8000)

Base URL: `http://localhost:8000/api`

### Fixtures

#### GET /fixtures
Get list of fixtures with optional filters.

**Query Parameters:**
- `season` (integer): Filter by season (e.g., 2023, 2024)
- `status` (string): Filter by status (scheduled, live, finished)
- `from_date` (ISO 8601): Start date for date range
- `to_date` (ISO 8601): End date for date range
- `team_id` (integer): Filter by team ID
- `limit` (integer): Number of results (default: 50, max: 100)
- `offset` (integer): Pagination offset (default: 0)

**Response:**
```json
{
  "fixtures": [
    {
      "id": 1,
      "api_football_id": 867946,
      "season": 2024,
      "round": "Regular Season - 20",
      "match_date": "2024-01-20T15:00:00Z",
      "home_team": {
        "id": 1,
        "name": "Arsenal",
        "code": "ARS"
      },
      "away_team": {
        "id": 2,
        "name": "Liverpool",
        "code": "LIV"
      },
      "status": "scheduled",
      "venue": "Emirates Stadium"
    }
  ],
  "total": 380,
  "limit": 50,
  "offset": 0
}
```

#### GET /fixtures/:id
Get single fixture with detailed information.

**Response:**
```json
{
  "id": 1,
  "api_football_id": 867946,
  "season": 2024,
  "round": "Regular Season - 20",
  "match_date": "2024-01-20T15:00:00Z",
  "home_team": {
    "id": 1,
    "name": "Arsenal",
    "code": "ARS",
    "logo_url": "https://..."
  },
  "away_team": {
    "id": 2,
    "name": "Liverpool",
    "code": "LIV",
    "logo_url": "https://..."
  },
  "home_score": null,
  "away_score": null,
  "status": "scheduled",
  "venue": "Emirates Stadium",
  "referee": "Michael Oliver"
}
```

### Odds

#### GET /fixtures/:id/odds
Get odds for a specific fixture.

**Query Parameters:**
- `bookmaker` (string): Filter by bookmaker name
- `market_type` (string): Filter by market type (h2h, totals, spreads)
- `latest` (boolean): Only get latest odds per bookmaker/outcome

**Response:**
```json
{
  "fixture_id": 1,
  "odds": [
    {
      "id": 1,
      "bookmaker": "bet365",
      "market_type": "h2h",
      "outcome": "home",
      "odds_value": 2.10,
      "recorded_at": "2024-01-19T10:00:00Z",
      "is_closing_line": false
    },
    {
      "id": 2,
      "bookmaker": "bet365",
      "market_type": "h2h",
      "outcome": "draw",
      "odds_value": 3.50,
      "recorded_at": "2024-01-19T10:00:00Z",
      "is_closing_line": false
    },
    {
      "id": 3,
      "bookmaker": "bet365",
      "market_type": "h2h",
      "outcome": "away",
      "odds_value": 3.20,
      "recorded_at": "2024-01-19T10:00:00Z",
      "is_closing_line": false
    }
  ]
}
```

### Weekly Picks

#### GET /picks/weekly
Get betting recommendations for the current week.

**Query Parameters:**
- `week` (ISO 8601 date): Specific week to get picks for (defaults to current week)
- `min_ev` (float): Minimum EV threshold (default: 0.03 = 3%)

**Response:**
```json
{
  "week_start": "2024-01-15",
  "week_end": "2024-01-21",
  "picks": [
    {
      "fixture": {
        "id": 1,
        "match_date": "2024-01-20T15:00:00Z",
        "home_team": "Arsenal",
        "away_team": "Liverpool"
      },
      "recommendation": {
        "bet_type": "h2h_home",
        "outcome": "home",
        "model_probability": 0.52,
        "best_odds": 2.10,
        "bookmaker": "bet365",
        "expected_value": 0.092,
        "ev_percentage": 9.2,
        "suggested_stake": 125.50,
        "kelly_fraction": 0.25,
        "confidence": "medium"
      },
      "model_prediction": {
        "home_win_prob": 0.52,
        "draw_prob": 0.25,
        "away_win_prob": 0.23,
        "model_version": "v1.0"
      }
    }
  ],
  "summary": {
    "total_picks": 5,
    "total_stake_recommended": 625.00,
    "avg_ev": 0.067,
    "avg_odds": 2.35
  }
}
```

### Bets

#### POST /bets
Record a placed bet.

**Request Body:**
```json
{
  "fixture_id": 1,
  "prediction_id": 1,
  "bet_type": "h2h_home",
  "stake": 125.50,
  "odds": 2.10,
  "expected_value": 0.092,
  "bookmaker": "bet365",
  "placed_at": "2024-01-19T14:30:00Z",
  "notes": "Strong home form, value at 2.10"
}
```

**Response:**
```json
{
  "id": 1,
  "fixture_id": 1,
  "status": "pending",
  "stake": 125.50,
  "odds": 2.10,
  "potential_payout": 263.55,
  "created_at": "2024-01-19T14:30:00Z"
}
```

#### PUT /bets/:id/settle
Settle a bet with result.

**Request Body:**
```json
{
  "status": "won",
  "payout": 263.55,
  "settled_at": "2024-01-20T17:00:00Z"
}
```

**Response:**
```json
{
  "id": 1,
  "status": "won",
  "stake": 125.50,
  "payout": 263.55,
  "profit_loss": 138.05,
  "roi_percentage": 110.0
}
```

#### GET /bets
Get list of bets.

**Query Parameters:**
- `status` (string): Filter by status (pending, won, lost, void)
- `from_date` (ISO 8601): Start date
- `to_date` (ISO 8601): End date
- `limit` (integer): Results per page
- `offset` (integer): Pagination offset

**Response:**
```json
{
  "bets": [
    {
      "id": 1,
      "fixture": {
        "match_date": "2024-01-20T15:00:00Z",
        "home_team": "Arsenal",
        "away_team": "Liverpool"
      },
      "bet_type": "h2h_home",
      "stake": 125.50,
      "odds": 2.10,
      "status": "won",
      "payout": 263.55,
      "profit_loss": 138.05,
      "placed_at": "2024-01-19T14:30:00Z",
      "settled_at": "2024-01-20T17:00:00Z"
    }
  ],
  "total": 45
}
```

### Performance Metrics

#### GET /performance/summary
Get overall performance summary.

**Query Parameters:**
- `from_date` (ISO 8601): Start date for calculations
- `to_date` (ISO 8601): End date for calculations

**Response:**
```json
{
  "period": {
    "from": "2024-01-01T00:00:00Z",
    "to": "2024-01-31T23:59:59Z",
    "days": 31
  },
  "metrics": {
    "total_bets": 45,
    "total_staked": 5625.00,
    "total_returned": 6187.50,
    "total_profit": 562.50,
    "roi_percentage": 10.0,
    "win_rate": 0.533,
    "avg_odds": 2.25,
    "avg_stake": 125.00,
    "num_wins": 24,
    "num_losses": 21,
    "biggest_win": 312.50,
    "biggest_loss": -125.00,
    "max_drawdown": -375.00,
    "sharpe_ratio": 1.85,
    "clv_average": 0.025
  },
  "bankroll": {
    "starting_balance": 10000.00,
    "current_balance": 10562.50,
    "peak_balance": 10875.00,
    "growth_percentage": 5.625
  }
}
```

#### GET /performance/daily
Get daily performance breakdown.

**Response:**
```json
{
  "daily_performance": [
    {
      "date": "2024-01-20",
      "num_bets": 3,
      "total_staked": 375.00,
      "profit_loss": 138.05,
      "roi": 36.8,
      "balance_eod": 10138.05
    }
  ]
}
```

### Bankroll

#### GET /bankroll/history
Get bankroll history over time.

**Query Parameters:**
- `from_date` (ISO 8601)
- `to_date` (ISO 8601)
- `granularity` (string): hourly, daily, weekly (default: daily)

**Response:**
```json
{
  "history": [
    {
      "id": 1,
      "balance": 10562.50,
      "total_staked": 5625.00,
      "total_profit_loss": 562.50,
      "roi_percentage": 10.0,
      "num_bets": 45,
      "win_rate": 0.533,
      "recorded_at": "2024-01-31T23:59:59Z"
    }
  ]
}
```

## Python ML Service API (Port 8001)

Base URL: `http://localhost:8001`

### Predictions

#### POST /predict
Get prediction for a single fixture.

**Request Body:**
```json
{
  "fixture_id": 1,
  "home_team_id": 1,
  "away_team_id": 2,
  "match_date": "2024-01-20T15:00:00Z"
}
```

**Response:**
```json
{
  "fixture_id": 1,
  "model_version": "v1.0",
  "predictions": {
    "home_win_prob": 0.52,
    "draw_prob": 0.25,
    "away_win_prob": 0.23
  },
  "predicted_outcome": "home",
  "confidence_score": 0.52,
  "features": {
    "home_form_last_5": 13,
    "away_form_last_5": 11,
    "home_goals_avg": 2.1,
    "away_goals_avg": 1.8,
    "h2h_home_wins_pct": 0.45,
    "position_differential": 2
  },
  "predicted_at": "2024-01-19T12:00:00Z"
}
```

#### POST /predict/batch
Get predictions for multiple fixtures.

**Request Body:**
```json
{
  "fixtures": [
    {
      "fixture_id": 1,
      "home_team_id": 1,
      "away_team_id": 2,
      "match_date": "2024-01-20T15:00:00Z"
    },
    {
      "fixture_id": 2,
      "home_team_id": 3,
      "away_team_id": 4,
      "match_date": "2024-01-20T17:30:00Z"
    }
  ]
}
```

**Response:**
```json
{
  "predictions": [
    {
      "fixture_id": 1,
      "predictions": {
        "home_win_prob": 0.52,
        "draw_prob": 0.25,
        "away_win_prob": 0.23
      }
    },
    {
      "fixture_id": 2,
      "predictions": {
        "home_win_prob": 0.38,
        "draw_prob": 0.28,
        "away_win_prob": 0.34
      }
    }
  ],
  "model_version": "v1.0",
  "batch_predicted_at": "2024-01-19T12:00:00Z"
}
```

### Model Management

#### GET /model/metrics
Get model performance metrics.

**Response:**
```json
{
  "model_version": "v1.0",
  "training_date": "2024-01-15T00:00:00Z",
  "metrics": {
    "accuracy": 0.58,
    "precision": 0.61,
    "recall": 0.58,
    "f1_score": 0.59,
    "brier_score": 0.21,
    "log_loss": 1.02,
    "roc_auc": 0.65
  },
  "backtest": {
    "num_matches": 380,
    "theoretical_roi": 0.075,
    "win_rate": 0.58
  },
  "training_data": {
    "num_samples": 1140,
    "seasons": [2021, 2022, 2023],
    "features_count": 18
  }
}
```

#### POST /model/train
Trigger model retraining.

**Request Body:**
```json
{
  "seasons": [2021, 2022, 2023, 2024],
  "hyperparameters": {
    "n_estimators": 200,
    "max_depth": 6,
    "learning_rate": 0.1
  }
}
```

**Response:**
```json
{
  "status": "training_started",
  "job_id": "train_20240119_120000",
  "estimated_duration_minutes": 15
}
```

#### GET /model/train/:job_id
Check training job status.

**Response:**
```json
{
  "job_id": "train_20240119_120000",
  "status": "completed",
  "started_at": "2024-01-19T12:00:00Z",
  "completed_at": "2024-01-19T12:14:23Z",
  "new_model_version": "v1.1",
  "metrics": {
    "accuracy": 0.59,
    "improvement_over_previous": 0.01
  }
}
```

## Error Responses

All APIs use consistent error response format:

```json
{
  "error": {
    "code": "INVALID_FIXTURE",
    "message": "Fixture not found",
    "details": {
      "fixture_id": 999
    }
  }
}
```

**Common HTTP Status Codes:**
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `404` - Not Found
- `422` - Validation Error
- `500` - Internal Server Error
- `503` - Service Unavailable
