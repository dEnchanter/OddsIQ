# OddsIQ MVP Implementation Plan

## Architecture Overview

**Multi-service architecture:**
- **Go Backend Service** - Data ingestion, REST API, business logic
- **Python ML Service** - Model training, predictions, feature engineering
- **Next.js Frontend** - Investor dashboard, visualization
- **PostgreSQL Database** - Centralized data storage

## Tech Stack

### Backend (Go)
- Go 1.21+
- Gin or Fiber (web framework)
- pgx (PostgreSQL driver)
- go-cron (job scheduling)

### ML Service (Python)
- Python 3.11+
- FastAPI (for ML API endpoints)
- XGBoost, scikit-learn
- pandas, numpy
- psycopg2 (PostgreSQL)

### Frontend (Next.js)
- Next.js 14+ (App Router)
- TypeScript
- Recharts (for charts)
- Tailwind CSS
- shadcn/ui (components)

### Database
- PostgreSQL 15+ (local)

## Implementation Phases

### Phase 1: Project Structure & Database (Week 1)

**1.1 Repository Structure**
```
OddsIQ/
├── backend/              # Go service
│   ├── cmd/
│   │   └── api/
│   ├── internal/
│   │   ├── api/         # HTTP handlers
│   │   ├── models/      # Data models
│   │   ├── repository/  # Database access
│   │   └── services/    # Business logic
│   ├── pkg/             # Shared packages
│   └── go.mod
├── ml-service/          # Python service
│   ├── app/
│   │   ├── models/      # ML models
│   │   ├── features/    # Feature engineering
│   │   └── api/         # FastAPI endpoints
│   ├── notebooks/       # Jupyter notebooks
│   ├── requirements.txt
│   └── pyproject.toml
├── frontend/            # Next.js
│   ├── src/
│   │   ├── app/
│   │   ├── components/
│   │   └── lib/
│   ├── package.json
│   └── next.config.js
├── database/
│   └── migrations/
├── docs/
└── docker-compose.yml   # For local development
```

**1.2 Database Schema**

Core tables:
- `teams` - Premier League teams
- `fixtures` - Match fixtures (3 seasons = ~1,140 matches)
- `odds` - Historical and current odds from multiple bookmakers
- `team_stats` - Team statistics (form, xG, etc.)
- `predictions` - Model predictions
- `bets` - Bet tracking and results
- `bankroll` - Bankroll history

**1.3 Initial Setup Tasks**
- [x] Initialize Go module
- [x] Initialize Python project with poetry/pip
- [x] Initialize Next.js project
- [x] Create database migration structure
- [x] Set up docker-compose for local PostgreSQL
- [x] Create environment variable templates

### Phase 2: Data Infrastructure (Weeks 1-2)

**2.1 Go Backend - Data Ingestion**

Implement API clients:
- API-Football client (fixtures, standings, team stats)
- The Odds API client (bookmaker odds)

Data ingestion services:
- Fixture sync service (daily job)
- Odds sync service (regular updates)
- Historical data backfill script

**2.2 Database Layer**

Repository pattern for:
- Teams CRUD
- Fixtures management
- Odds storage and retrieval
- Stats aggregation

**2.3 Deliverables**
- Working data pipeline pulling live Premier League data
- 3 seasons of historical data loaded
- Automated daily refresh jobs
- API endpoints to query fixtures and odds

### Phase 3: Feature Engineering & ML Models (Weeks 3-5) - MULTI-MARKET

**3.1 Python ML Service - Feature Engineering**

Base feature modules (shared across all models):
- Form metrics (last 5 games, points, goals)
- Head-to-head statistics
- Home/away performance splits
- League position differentials
- Expected goals (xG) integration

Market-specific features:
- Goals features (for O/U, BTTS)
- Defensive/offensive strength metrics
- Scoring/conceding patterns

**3.2 Model Development (Incremental Rollout)**

**Week 3: 1X2 Model (Home/Draw/Away)**
- XGBoost multi-class classifier (3 outcomes)
- Train on 3 seasons of data
- Cross-validation with time-series splits
- Probability calibration
- Backtest: Target 55-60% accuracy, +3-5% ROI

**Week 4: Over/Under 2.5 Goals Model**
- XGBoost binary classifier
- Additional features: goal averages, scoring patterns
- Backtest: Target 60-65% accuracy
- Most popular totals market

**Week 5: BTTS Model (Both Teams to Score)**
- XGBoost binary classifier
- Features: clean sheets, defensive/offensive strength
- Backtest: Target 58-62% accuracy
- Independent of match outcome

**3.3 ML API Endpoints**

FastAPI endpoints:
- `POST /predict` - Get prediction for a single market
- `POST /predict/all-markets` - Get predictions across all markets
- `POST /batch-predict` - Weekly batch predictions
- `GET /model/metrics` - Model performance metrics per market
- `POST /train/:market` - Trigger market-specific retraining

**3.4 Deliverables**
- Feature engineering pipeline (shared + market-specific)
- 3 trained models (1X2, O/U 2.5, BTTS)
- Model performance reports per market
- Backtest showing theoretical ROI per market
- Multi-market prediction API operational

### Phase 4: Betting Engine & Market Selector (Weeks 6-7)

**4.1 Go Backend - Value Detection & Stake Sizing**

Betting logic services:
- Expected Value (EV) calculator (works across all markets)
- Kelly Criterion stake sizing (1/4 Kelly, configurable per market)
- Value bet filter (+3% EV minimum)
- **Smart Market Selector** - evaluates ALL markets, picks highest EV

**4.2 Multi-Market Integration**

- Go backend calls Python ML service for ALL market predictions
- Fetches odds for all available markets (1X2, O/U, BTTS, etc.)
- Calculates EV for EVERY market outcome
- Ranks recommendations by EV
- Can recommend multiple markets per fixture

**Example Flow:**
```
Fixture: Arsenal vs Liverpool
→ Get predictions: h2h, totals_2_5, btts
→ Get odds for all markets
→ Calculate EV:
  - Home Win: -2.3% ❌
  - Over 2.5: +8.1% ✅ (Highest EV)
  - BTTS Yes: +6.3% ✅
→ Recommend: Over 2.5 (primary) + BTTS (secondary)
```

**4.3 API Endpoints**

- `GET /api/picks/weekly` - Current week's recommendations (best EV across all markets)
- `GET /api/picks/weekly?market=totals` - Filter by market type
- `POST /api/bets` - Log bet placement (with market type)
- `PUT /api/bets/:id/settle` - Update bet result
- `GET /api/bankroll` - Bankroll history
- `GET /api/performance` - Performance metrics (overall + per market)
- `GET /api/performance/by-market` - Performance breakdown by market

**4.4 Deliverables**
- Multi-market betting recommendation engine
- Smart market selector (highest EV across all markets)
- Bet tracking system with market categorization
- Performance calculation per market type

### Phase 5: Accumulator Builder (Week 7)

**5.1 Go Backend - Smart Accumulator Generation**

Accumulator logic:
- Combine 2-3 uncorrelated picks into parlays
- Correlation detection (avoid same fixture, same team)
- Accumulator EV calculation (combined probability × combined odds)
- Conservative Kelly sizing (1/8 Kelly for parlays vs 1/4 for singles)
- Maximum 20% of weekly stake allocation on accumulators

**5.2 Correlation Filtering**

Avoid correlated events:
- Same fixture (can't combine Arsenal Win + Over 2.5 from same match)
- Same team in different fixtures (debatable, configurable)
- Highly correlated markets (research needed)

**5.3 Accumulator EV Calculation**

```
Example 3-leg accumulator:
Leg 1: Arsenal Over 2.5, Prob=0.65, Odds=1.80
Leg 2: Brighton BTTS Yes, Prob=0.68, Odds=1.70
Leg 3: Liverpool Away, Prob=0.55, Odds=2.20

Combined Probability: 0.65 × 0.68 × 0.55 = 0.243 (24.3%)
Combined Odds: 1.80 × 1.70 × 2.20 = 6.73
EV = (0.243 × 6.73) - 1 = 0.635 = 63.5% ✅

Minimum EV for accumulators: 5% (higher than 3% for singles)
```

**5.4 Database Schema Updates**

New tables:
- `accumulators` - Accumulator bet tracking
- `accumulator_legs` - Individual legs within accumulator

**5.5 API Endpoints**

- `GET /api/picks/weekly/accumulators` - Weekly accumulator recommendations
- `POST /api/accumulators` - Record accumulator placement
- `PUT /api/accumulators/:id/settle` - Settle accumulator (all legs must win)

**5.6 Deliverables**
- Smart accumulator builder (2-3 legs)
- Correlation detection logic
- Accumulator-specific Kelly sizing
- 2-3 accumulator recommendations per week
- Accumulator tracking and settlement

### Phase 6: Investor Dashboard (Weeks 7-8)

**6.1 Next.js Frontend**

Pages:
- Dashboard home (current week picks, P&L summary, accumulators)
- Picks page (singles + accumulators)
- Performance page (charts, metrics, ROI tracking)
- Bets history (singles and accumulators with results)

Components:
- Pick cards with model confidence
- P&L chart (Recharts)
- Performance metrics dashboard
- Bet entry form

**6.2 Key Features**

- Real-time performance metrics
- Weekly picks display (singles + accumulators)
- Odds and stake recommendations
- ROI, win rate, Sharpe ratio visualization
- Bankroll growth chart
- Model accuracy tracking
- Accumulator display with legs breakdown

**6.3 Deliverables**
- Fully functional investor dashboard
- Accumulator visualization
- Responsive design
- Connected to Go API
- Live data updates

### Phase 7: Testing & Live Validation (Weeks 8-9)

**6.1 Testing Strategy**

- Unit tests for critical Go services
- Python ML model tests
- Integration tests for API endpoints
- End-to-end testing workflow

**7.2 Paper Trading**

- Generate picks for 2-3 weekends
- Singles: 15-20 per week
- Accumulators: 2-3 per week (2-3 legs each)
- Track results in real-time
- Validate system performance
- Adjust thresholds if needed

**7.3 Investor Demo Package**

- System documentation
- Performance report (singles + accumulators)
- Demo walkthrough
- Deployment guide
- Accumulator strategy explanation

## Critical Files to Create

### Configuration
- `backend/config/config.go` - Configuration management
- `ml-service/config.py` - ML service config
- `.env.example` - Environment variables template

### Database
- `database/migrations/001_initial_schema.sql`
- `database/migrations/002_add_betting_tables.sql`

### Go Backend Core
- `backend/internal/repository/fixtures.go`
- `backend/internal/repository/odds.go`
- `backend/internal/services/data_sync.go`
- `backend/internal/services/betting.go`
- `backend/internal/api/handlers.go`

### Python ML Service Core
- `ml-service/app/features/form_metrics.py`
- `ml-service/app/features/h2h_stats.py`
- `ml-service/app/models/xgboost_model.py`
- `ml-service/app/api/predictions.py`

### Frontend Core
- `frontend/src/app/page.tsx` - Dashboard
- `frontend/src/components/PickCard.tsx`
- `frontend/src/components/PerformanceChart.tsx`
- `frontend/src/lib/api.ts` - API client

## Verification & Testing

### End-to-End Verification Flow

1. **Data Pipeline**: Verify fixtures and odds are syncing
   ```bash
   # Check database has data
   psql -d oddsiq -c "SELECT COUNT(*) FROM fixtures;"
   ```

2. **ML Service**: Test prediction endpoint
   ```bash
   curl -X POST http://localhost:8001/predict \
     -H "Content-Type: application/json" \
     -d '{"fixture_id": 123}'
   ```

3. **Go API**: Test picks generation
   ```bash
   curl http://localhost:8000/api/picks/weekly
   ```

4. **Frontend**: Access dashboard at http://localhost:3000

### Success Criteria

- [ ] Database populated with 3 seasons of Premier League data (1,140+ matches)
- [ ] Automated daily data refresh working
- [ ] Odds data for multiple markets (1X2, O/U, BTTS) stored
- [ ] **3 ML models trained and evaluated:**
  - [ ] 1X2 model: 55-60% accuracy, +3-5% backtest ROI
  - [ ] Over/Under 2.5: 60-65% accuracy
  - [ ] BTTS: 58-62% accuracy
- [ ] **Smart market selector operational** (picks best EV across all markets)
- [ ] **Accumulator builder operational** (2-3 leg parlays with correlation detection)
- [ ] Backtested theoretical ROI: +6-10% (singles + accumulators)
- [ ] Betting engine generating weekly picks:
  - [ ] 15-20 single bets from multiple markets
  - [ ] 2-3 accumulators (2-3 legs each)
- [ ] Dashboard displaying:
  - [ ] Multi-market single picks
  - [ ] Accumulator picks with legs breakdown
  - [ ] Performance breakdown by type (singles vs accumulators)
  - [ ] Overall metrics
- [ ] 2-3 weekends of paper trading documented (singles + accumulators)
- [ ] System ready for real-money testing across markets and bet types

## Next Steps After MVP

1. Deploy to AWS (EC2 for services, RDS for database)
2. Set up monitoring and alerts
3. Begin real-money testing with small bankroll
4. Document 8-12 weeks of live results for investors
5. Plan Phase 2 enhancements (ensemble models, multi-ticket strategy)
