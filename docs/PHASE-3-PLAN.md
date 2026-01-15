# Phase 3: ML Model Development - Action Plan

**Timeline:** Weeks 3-5 (3 weeks)
**Goal:** Build trained XGBoost model that predicts match outcomes
**Success:** >55% accuracy, positive backtested ROI

---

## Overview

### What We're Building

A complete machine learning pipeline that:
1. **Extracts features** from 1,140 historical matches
2. **Trains XGBoost model** to predict Home/Draw/Away outcomes
3. **Validates** on held-out test data
4. **Backtests** strategy with synthetic odds
5. **Serves predictions** via FastAPI

### Input
- ✅ 1,140 matches in PostgreSQL (2022-2024)
- ✅ 24 teams with details
- ✅ Match results, scores, dates

### Output
- ✅ Trained ML model (xgboost_v1.pkl)
- ✅ Feature engineering pipeline
- ✅ Prediction API service (FastAPI on port 8001)
- ✅ Backtest results showing theoretical ROI
- ✅ Model performance metrics

---

## Week-by-Week Breakdown

### Week 3: Feature Engineering (Days 1-7)

**Goal:** Extract predictive features from 1,140 matches

**Days 1-2: Setup Python ML Service**
- Create ml-service directory structure
- Set up Python virtual environment
- Install dependencies (pandas, numpy, scikit-learn, xgboost, fastapi)
- Create database connection from Python
- Test reading fixtures from PostgreSQL

**Days 3-5: Build Feature Extractors**
- Form features (last 5 games, points, goals, wins/draws/losses)
- Head-to-head features (past meetings, goal differences)
- League position features (table position, points difference)
- Home/away splits (performance at home vs away)
- Goal metrics (avg goals scored/conceded, clean sheets)

**Days 6-7: Create Training Dataset**
- Extract features for all 1,140 matches
- Create labels (Home Win=1, Draw=X, Away Win=2)
- Handle missing data (early season matches)
- Save to CSV for analysis
- Validate feature distributions

**Deliverables:**
- [ ] Python ML service structure
- [ ] Feature engineering pipeline
- [ ] Training dataset CSV (1,140 rows × 30+ features)
- [ ] Feature analysis notebook

---

### Week 4: Model Training & Validation (Days 8-14)

**Goal:** Train model and validate it works

**Days 8-9: Data Preparation**
- Train/test split (80/20, chronologically)
- Feature scaling/normalization
- Handle class imbalance (more home wins than draws)
- Create validation sets

**Days 10-11: Model Training**
- Train XGBoost classifier
- Tune hyperparameters (max_depth, learning_rate, n_estimators)
- Cross-validation with time-series splits
- Feature importance analysis

**Days 12-13: Model Validation**
- Test on held-out data (20% of matches)
- Calculate accuracy, precision, recall
- Probability calibration
- Confusion matrix analysis
- Compare against baseline (always predict home win)

**Day 14: Backtesting Framework**
- Create synthetic odds generator (realistic market odds)
- Build backtesting engine
- Simulate betting strategy (Kelly Criterion, EV threshold)
- Calculate theoretical ROI, Sharpe ratio, max drawdown

**Deliverables:**
- [ ] Trained XGBoost model (xgboost_v1.pkl)
- [ ] Model achieves >55% accuracy
- [ ] Backtesting shows positive ROI
- [ ] Model training notebook/report
- [ ] Feature importance report

---

### Week 5: API Service & Integration (Days 15-21)

**Goal:** Make model accessible via API

**Days 15-16: Build FastAPI Service**
- Create prediction endpoints
- POST /predict - Single fixture prediction
- POST /batch-predict - Multiple fixtures
- GET /model/metrics - Model performance stats
- Input validation and error handling

**Days 17-18: Go Backend Integration**
- Create ml_client.go (calls Python API)
- Create prediction storage in database
- Test end-to-end prediction flow
- Error handling and retries

**Days 19-20: Testing & Refinement**
- Integration tests
- Performance testing
- Model refinement based on results
- Documentation

**Day 21: Review & Next Phase Prep**
- Review all Phase 3 deliverables
- Analyze model performance
- Document findings
- Plan Phase 4 (web scraping setup)

**Deliverables:**
- [ ] FastAPI service running on port 8001
- [ ] Go backend can call ML service
- [ ] Predictions stored in database
- [ ] End-to-end workflow tested
- [ ] API documentation
- [ ] Phase 3 completion report

---

## Technical Architecture

### ML Service Structure

```
ml-service/
├── app/
│   ├── __init__.py
│   ├── main.py                  # FastAPI app entry point
│   ├── config.py                # Configuration settings
│   │
│   ├── database/
│   │   ├── __init__.py
│   │   └── connection.py        # PostgreSQL connection
│   │
│   ├── models/
│   │   ├── __init__.py
│   │   └── xgboost_model.py     # Model training & prediction
│   │
│   ├── features/
│   │   ├── __init__.py
│   │   ├── form_metrics.py      # Team form features
│   │   ├── h2h_stats.py         # Head-to-head features
│   │   ├── league_position.py   # Position-based features
│   │   ├── team_stats.py        # General team stats
│   │   └── feature_builder.py   # Orchestrates all features
│   │
│   ├── backtesting/
│   │   ├── __init__.py
│   │   ├── synthetic_odds.py    # Generate realistic odds
│   │   └── backtest_engine.py   # Strategy backtesting
│   │
│   └── api/
│       ├── __init__.py
│       ├── predictions.py       # Prediction endpoints
│       └── schemas.py           # Pydantic models
│
├── notebooks/
│   ├── 01_data_exploration.ipynb
│   ├── 02_feature_engineering.ipynb
│   ├── 03_model_training.ipynb
│   └── 04_backtest_analysis.ipynb
│
├── models/
│   └── xgboost_v1.pkl          # Trained model
│
├── data/
│   ├── training_data.csv       # Feature dataset
│   └── backtest_results.csv    # Backtest results
│
├── tests/
│   ├── __init__.py
│   ├── test_features.py
│   └── test_predictions.py
│
├── requirements.txt
├── pyproject.toml
├── .env.example
└── README.md
```

---

## Features to Extract

### 1. Team Form Features (Last 5 Games)

**For each team:**
- Last 5 results (WWDLL → 7 points)
- Last 5 goals scored
- Last 5 goals conceded
- Last 5 points earned
- Win percentage last 5
- Clean sheets last 5

### 2. Overall Season Form

**For each team:**
- Total points this season
- Goals scored this season
- Goals conceded this season
- Goal difference
- Win/Draw/Loss counts
- Average goals per game
- Clean sheet percentage

### 3. Home/Away Splits

**For home team:**
- Home win percentage
- Home goals scored avg
- Home goals conceded avg
- Home points avg

**For away team:**
- Away win percentage
- Away goals scored avg
- Away goals conceded avg
- Away points avg

### 4. League Position Features

- Home team position
- Away team position
- Position difference
- Points difference between teams
- Games behind leader

### 5. Head-to-Head History

- Last 5 H2H results
- H2H goals for home team
- H2H goals for away team
- Home win percentage in H2H
- Average goals in H2H matches

### 6. Recent Momentum

- Points earned last 3 games
- Goal difference last 3 games
- Scoring streak (consecutive games scoring)
- Conceding streak
- Unbeaten run

### 7. Match Context

- Day of week
- Month of season
- Games played so far (fatigue)
- Days since last match (rest)

**Total: ~30-40 features per match**

---

## Model Training Approach

### 1. Data Preparation

```python
# Train/test split (chronologically)
train_data = matches[matches['season'] < 2024]  # 2022-2023: ~760 matches
test_data = matches[matches['season'] == 2024]  # 2024: 380 matches
```

### 2. XGBoost Configuration

```python
params = {
    'objective': 'multi:softprob',  # Multi-class classification
    'num_class': 3,                 # Home/Draw/Away
    'max_depth': 6,
    'learning_rate': 0.1,
    'n_estimators': 100,
    'subsample': 0.8,
    'colsample_bytree': 0.8,
    'eval_metric': 'mlogloss'
}
```

### 3. Training Process

1. Load training data
2. Extract features
3. Train XGBoost model
4. Validate on test set
5. Calibrate probabilities
6. Save model to disk

### 4. Evaluation Metrics

- **Accuracy:** Overall prediction accuracy (target: >55%)
- **Precision/Recall:** Per class (Home/Draw/Away)
- **Log Loss:** Probability calibration quality
- **ROI:** Backtested return on investment
- **Sharpe Ratio:** Risk-adjusted returns
- **Max Drawdown:** Worst losing streak

---

## Backtesting Strategy

### Synthetic Odds Generation

**Approach:** Generate realistic odds based on:
- Historical market averages by league position difference
- Home advantage factor (~1.5 goals)
- Recent form adjustments
- Add market margin (~5-10% overround)

**Example:**
```python
# Home team: 2nd place, good form
# Away team: 15th place, bad form
# Expected probabilities: Home 60%, Draw 25%, Away 15%
# Add overround (5%): 63%, 26%, 16%
# Convert to odds: 1.59, 3.85, 6.25
```

### Betting Strategy

**Rules:**
1. Model predicts probabilities for Home/Draw/Away
2. Convert odds to implied probabilities
3. Calculate Expected Value (EV) = (Model Prob × Odds) - 1
4. Only bet if EV > 3% (value threshold)
5. Use Kelly Criterion for stake sizing (¼ Kelly for safety)
6. Track P&L over entire test set

**Example:**
```
Model: Home Win 65%
Odds: 1.80 (implied 55.6%)
EV = (0.65 × 1.80) - 1 = 0.17 = 17% edge ✅ BET
Kelly Stake = (0.65 × 1.80 - 1) / 1.80 = 9.4% of bankroll
¼ Kelly = 2.35% of bankroll
```

---

## API Endpoints

### 1. POST /predict

**Request:**
```json
{
  "home_team_id": 14,  // Man City
  "away_team_id": 6,   // Liverpool
  "match_date": "2026-02-15"
}
```

**Response:**
```json
{
  "prediction": {
    "home_win_prob": 0.45,
    "draw_prob": 0.28,
    "away_win_prob": 0.27,
    "predicted_outcome": "home_win",
    "confidence": 0.45
  },
  "features": {
    "home_form_last5": 9,
    "away_form_last5": 13,
    "position_diff": 2
  },
  "model_version": "v1.0"
}
```

### 2. POST /batch-predict

**Request:**
```json
{
  "fixtures": [
    {"home_team_id": 14, "away_team_id": 6, "match_date": "2026-02-15"},
    {"home_team_id": 8, "away_team_id": 13, "match_date": "2026-02-15"}
  ]
}
```

**Response:**
```json
{
  "predictions": [
    {...},
    {...}
  ]
}
```

### 3. GET /model/metrics

**Response:**
```json
{
  "model_version": "v1.0",
  "training_date": "2026-01-15",
  "training_samples": 760,
  "test_samples": 380,
  "metrics": {
    "accuracy": 0.567,
    "home_precision": 0.61,
    "draw_precision": 0.32,
    "away_precision": 0.54,
    "log_loss": 1.02
  },
  "backtest": {
    "roi": 0.087,
    "sharpe_ratio": 1.24,
    "max_drawdown": 0.15,
    "total_bets": 82,
    "winning_bets": 48
  }
}
```

---

## Go Backend Integration

### ml_client.go

```go
package services

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type MLClient struct {
    baseURL    string
    httpClient *http.Client
}

type PredictionRequest struct {
    HomeTeamID int    `json:"home_team_id"`
    AwayTeamID int    `json:"away_team_id"`
    MatchDate  string `json:"match_date"`
}

type PredictionResponse struct {
    Prediction struct {
        HomeWinProb   float64 `json:"home_win_prob"`
        DrawProb      float64 `json:"draw_prob"`
        AwayWinProb   float64 `json:"away_win_prob"`
        Predicted     string  `json:"predicted_outcome"`
        Confidence    float64 `json:"confidence"`
    } `json:"prediction"`
}

func NewMLClient(baseURL string) *MLClient {
    return &MLClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *MLClient) Predict(req PredictionRequest) (*PredictionResponse, error) {
    body, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    resp, err := c.httpClient.Post(
        c.baseURL+"/predict",
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var prediction PredictionResponse
    if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
        return nil, err
    }

    return &prediction, nil
}
```

---

## Success Criteria

### Model Performance
- [ ] Accuracy >55% on test set (2024 season)
- [ ] Better than baseline (always predict home win: ~46%)
- [ ] Predictions are calibrated (predicted 60% → actual 60%)
- [ ] Home wins: precision >60%
- [ ] Draws: precision >30% (draws are hard!)
- [ ] Away wins: precision >50%

### Backtesting Results
- [ ] Positive ROI on test set (>5%)
- [ ] Sharpe ratio >1.0 (risk-adjusted returns)
- [ ] Max drawdown <20% (reasonable risk)
- [ ] 50+ value bets identified (EV >3%)
- [ ] Win rate on value bets >55%

### Technical Requirements
- [ ] Model trains in <10 minutes
- [ ] Predictions return in <100ms
- [ ] API handles concurrent requests
- [ ] Features extract correctly
- [ ] No data leakage (future info in training)

### Documentation
- [ ] Model training documented
- [ ] Feature descriptions written
- [ ] API documentation complete
- [ ] Backtest results analyzed
- [ ] Next steps identified

---

## THIS WEEK: Action Plan

### Day 1 (Today/Tomorrow)

**Setup (2-3 hours):**
1. Create ml-service directory structure
2. Set up Python virtual environment
3. Install dependencies
4. Test PostgreSQL connection from Python
5. Read sample fixtures to validate

**Files to create:**
```
ml-service/
├── requirements.txt
├── app/
│   ├── __init__.py
│   ├── config.py
│   └── database/
│       ├── __init__.py
│       └── connection.py
```

**Commands:**
```bash
# Create virtual environment
cd ml-service
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install pandas numpy scikit-learn xgboost fastapi uvicorn psycopg2-binary python-dotenv

# Save requirements
pip freeze > requirements.txt

# Test connection
python -c "import psycopg2; print('OK')"
```

### Day 2-3

**Feature Engineering (6-8 hours):**
1. Create form_metrics.py (last 5 games)
2. Create h2h_stats.py (head-to-head)
3. Create league_position.py (table position)
4. Create team_stats.py (overall stats)
5. Create feature_builder.py (orchestrates all)

### Day 4-5

**Dataset Creation (4-6 hours):**
1. Extract features for all 1,140 matches
2. Create labels (Home/Draw/Away)
3. Save to CSV
4. Analyze in Jupyter notebook
5. Validate quality

### Day 6-7

**Initial Model Training (6-8 hours):**
1. Train/test split
2. Train XGBoost model
3. Evaluate on test set
4. Document results
5. Iterate if needed

---

## Quick Start Commands

### Setup ML Service

```bash
# Navigate to project root
cd C:\Users\afolabi.opaleye\Desktop\builds\personal-builds\AI-builds\OddsIQ

# Create ml-service directory
mkdir ml-service
cd ml-service

# Create virtual environment
python -m venv venv
venv\Scripts\activate

# Create directory structure
mkdir -p app/database app/models app/features app/backtesting app/api
mkdir notebooks models data tests

# Create __init__.py files
type nul > app/__init__.py
type nul > app/database/__init__.py
type nul > app/models/__init__.py
type nul > app/features/__init__.py
type nul > app/backtesting/__init__.py
type nul > app/api/__init__.py

# Install dependencies
pip install pandas numpy scikit-learn xgboost fastapi uvicorn psycopg2-binary python-dotenv jupyter

# Save requirements
pip freeze > requirements.txt
```

### Test Database Connection

Create `app/database/connection.py`:
```python
import os
import psycopg2
from psycopg2.extras import RealDictCursor
from dotenv import load_dotenv

load_dotenv('../backend/.env')

def get_connection():
    return psycopg2.connect(
        os.getenv('DATABASE_URL'),
        cursor_factory=RealDictCursor
    )

def test_connection():
    conn = get_connection()
    cursor = conn.cursor()
    cursor.execute("SELECT COUNT(*) FROM fixtures")
    count = cursor.fetchone()['count']
    print(f"✅ Connected! Found {count} fixtures")
    conn.close()

if __name__ == "__main__":
    test_connection()
```

Test it:
```bash
python app/database/connection.py
```

---

## Dependencies (requirements.txt)

```
# Data processing
pandas==2.1.4
numpy==1.26.2

# Machine Learning
scikit-learn==1.3.2
xgboost==2.0.3

# API
fastapi==0.108.0
uvicorn==0.25.0
pydantic==2.5.3

# Database
psycopg2-binary==2.9.9

# Utilities
python-dotenv==1.0.0

# Development
jupyter==1.0.0
matplotlib==3.8.2
seaborn==0.13.1
```

---

## Expected Timeline

| Week | Days | Focus | Deliverable |
|------|------|-------|-------------|
| 3 | 1-2 | Setup + DB connection | Python service structure |
| 3 | 3-5 | Feature engineering | Feature extractors built |
| 3 | 6-7 | Dataset creation | training_data.csv (1,140 rows) |
| 4 | 8-9 | Data prep + split | Train/test sets ready |
| 4 | 10-11 | Model training | xgboost_v1.pkl trained |
| 4 | 12-13 | Validation | Model metrics documented |
| 4 | 14 | Backtesting | Backtest results |
| 5 | 15-16 | FastAPI service | API running on :8001 |
| 5 | 17-18 | Go integration | End-to-end flow works |
| 5 | 19-20 | Testing | Integration tests pass |
| 5 | 21 | Review | Phase 3 complete report |

---

## Risk Mitigation

### Risk 1: Model accuracy <55%
**Mitigation:**
- Try different features
- Ensemble methods
- More data (add 2021 season)
- Different algorithms (Random Forest, Neural Net)

### Risk 2: Can't extract features
**Mitigation:**
- Start simple (just form last 5)
- Add complexity gradually
- Manual feature engineering
- Ask for help on specific features

### Risk 3: Backtesting shows negative ROI
**Mitigation:**
- Adjust EV threshold (try 5%, 10%)
- Different betting strategies
- Focus on specific outcomes (only home/away, skip draws)
- Accept it - model needs more work

### Risk 4: Python/ML unfamiliar
**Mitigation:**
- Follow tutorials for XGBoost
- Use ChatGPT for code help
- Copy from examples online
- Start with simplest possible model

---

## Resources

### Learning XGBoost
- Official docs: https://xgboost.readthedocs.io/
- Tutorial: https://machinelearningmastery.com/xgboost-python/

### FastAPI
- Official docs: https://fastapi.tiangolo.com/
- Tutorial: https://fastapi.tiangolo.com/tutorial/

### Football Prediction
- Research papers on football prediction
- Kaggle competitions for football
- Reddit: r/sportsbook, r/MachineLearning

---

## Next Session Goals

**By end of next session:**
1. [ ] ML service directory created
2. [ ] Python environment set up
3. [ ] Database connection working
4. [ ] Started on feature engineering
5. [ ] Read sample fixtures from database

**Time needed:** 2-4 hours

---

**Ready to start building the ML model!** 🤖🚀
