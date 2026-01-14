# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**OddsIQ** is an AI/ML-powered betting prediction system designed as a SaaS platform. This repository contains the planning documentation and will house the implementation of the prediction engine and associated services.

## Current Repository State

**Phase 2 Complete** - Data infrastructure fully implemented:
- ✅ **Phase 1**: Project structure and foundation
  - Multi-service architecture (Go backend, Python ML service, Next.js frontend)
  - Database schema with 10 tables (including accumulators)
  - Docker Compose development environment

- ✅ **Phase 2**: Data ingestion pipeline
  - API-Football client (fixtures, teams, standings)
  - The Odds API client (multi-market odds: H2H, Totals, BTTS)
  - Repository layer (Teams, Fixtures, Odds, TeamStats)
  - Data sync services (fixture sync, odds sync)
  - Automated scheduling (cron jobs)
  - Historical backfill script
  - API endpoints for data querying

**Next:** Phase 3 - Feature Engineering & ML Model (Weeks 3-5)

**See:** `docs/PHASE-2-COMPLETE.md` for detailed Phase 2 summary

## Architecture

### Multi-Service Design

```
┌─────────────────┐
│  Next.js        │
│  Frontend       │ ← Investor dashboard (Port 3000)
└────────┬────────┘
         │
    ┌────▼─────┐
    │   Go     │
    │  Backend │ ← API, data ingestion (Port 8000)
    │  (Gin)   │
    └─┬──────┬─┘
      │      │
      │      └──────► PostgreSQL (fixtures, odds, bets)
      │
      │ HTTP
      ▼
┌─────────────────┐
│   Python ML     │ ← XGBoost predictions (Port 8001)
│   (FastAPI)     │
└─────────────────┘
```

### Tech Stack

- **Backend**: Go 1.21+, Gin web framework, pgx (PostgreSQL driver)
- **ML Service**: Python 3.11+, FastAPI, XGBoost, scikit-learn, pandas, numpy
- **Frontend**: Next.js 14+, TypeScript, Tailwind CSS, Recharts
- **Database**: PostgreSQL 15+

## Development Standards

**Pattern:** Clean Architecture + Selective TDD

### Architecture Layers
```
API Handlers (Presentation)
    ↓
Services (Business Logic)
    ↓
Repository (Data Access)
    ↓
Models (Domain)
```

### Testing Strategy

**MUST TEST (TDD - Write tests FIRST):**
- ✅ Kelly Criterion stake sizing
- ✅ Expected Value (EV) calculations
- ✅ Profit/loss calculations
- ✅ Bet placement/settlement logic
- ✅ ML model predictions
- ✅ Feature engineering
- ✅ Accumulator combinations
- ✅ Correlation detection

**SHOULD TEST (Test after implementation):**
- ⚠️ Data sync workflows
- ⚠️ Fixture matching algorithms
- ⚠️ Odds parsing

**CAN SKIP (Low priority for MVP):**
- ❌ HTTP API clients (mock external APIs)
- ❌ Simple CRUD operations
- ❌ UI components

**See:** `docs/DEVELOPMENT-STANDARDS.md` for complete guidelines

## Development Commands

### Start All Services (Docker)
```bash
docker-compose up -d          # Start all services
docker-compose logs -f        # View logs
docker-compose down           # Stop all services
```

### Local Development

**Backend (Go)**
```bash
cd backend
go mod download               # Install dependencies
go run cmd/api/main.go       # Start API server (port 8000)
go test ./...                # Run tests
```

**ML Service (Python)**
```bash
cd ml-service
python -m venv venv          # Create virtual environment
source venv/bin/activate     # Activate (use venv\Scripts\activate on Windows)
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8001  # Start ML service
pytest                       # Run tests
```

**Frontend (Next.js)**
```bash
cd frontend
npm install                  # Install dependencies
npm run dev                  # Start dev server (port 3000)
npm run build                # Production build
npm test                     # Run tests
```

### Database Operations

**Setup PostgreSQL**
```bash
# Via Docker
docker-compose up -d postgres

# Via psql
psql -U oddsiq_user -d oddsiq
```

**Run Migrations**
```bash
# Migrations are in database/migrations/
# 001_initial_schema.up.sql - Teams, fixtures, odds, team_stats
# 002_add_betting_tables.up.sql - Predictions, bets, bankroll

# Apply manually via psql
psql -U oddsiq_user -d oddsiq -f database/migrations/001_initial_schema.up.sql
psql -U oddsiq_user -d oddsiq -f database/migrations/002_add_betting_tables.up.sql
```

## Project Structure

```
OddsIQ/
├── backend/                 # Go backend service
│   ├── cmd/
│   │   ├── api/            # Main API server
│   │   ├── migrate/        # Migration tool
│   │   ├── backfill/       # Historical data loader
│   │   └── sync/           # Data sync jobs
│   ├── internal/
│   │   ├── api/            # HTTP handlers and routes
│   │   ├── models/         # Data models
│   │   ├── repository/     # Database access layer
│   │   └── services/       # Business logic
│   ├── pkg/
│   │   ├── database/       # DB connection
│   │   ├── apifootball/    # API-Football client
│   │   └── oddsapi/        # The Odds API client
│   ├── config/             # Configuration management
│   └── go.mod
├── ml-service/             # Python ML service
│   ├── app/
│   │   ├── api/            # FastAPI endpoints (predictions.py)
│   │   ├── models/         # ML models (xgboost_model.py)
│   │   └── features/       # Feature engineering
│   ├── config/             # ML config
│   └── requirements.txt
├── frontend/               # Next.js dashboard
│   ├── app/                # Next.js 14 App Router
│   │   ├── layout.tsx      # Root layout
│   │   ├── page.tsx        # Dashboard home
│   │   └── globals.css
│   ├── components/         # React components
│   ├── lib/                # API client (api.ts)
│   └── package.json
├── database/
│   └── migrations/         # SQL migration files
├── docs/                   # Documentation
│   ├── setup-guide.md
│   ├── database-schema.md
│   ├── api-specification.md
│   └── architecture-decisions.md
└── docker-compose.yml
```

## Key Implementation Details

### Database Schema

**Core Tables:**
- `teams` - Premier League teams (20 teams)
- `fixtures` - Match fixtures (~1,140 for 3 seasons)
- `odds` - Bookmaker odds with timestamp tracking
- `team_stats` - Team performance metrics for features
- `predictions` - ML model predictions
- `bets` - Bet tracking and results
- `bankroll` - Bankroll snapshots
- `model_performance` - Model evaluation metrics

### API Endpoints

**Go Backend (`:8000/api`)**
- `GET /fixtures` - List fixtures with filters
- `GET /fixtures/:id` - Single fixture details
- `GET /fixtures/:id/odds` - Fixture odds
- `GET /picks/weekly` - Weekly betting recommendations
- `GET /bets` - List bets
- `POST /bets` - Record bet placement
- `PUT /bets/:id/settle` - Settle bet result
- `GET /performance/summary` - Performance metrics
- `GET /bankroll/history` - Bankroll tracking

**Python ML Service (`:8001/api`)**
- `POST /predict` - Single fixture prediction
- `POST /predict/batch` - Batch predictions
- `GET /model/metrics` - Model performance
- `POST /model/train` - Trigger retraining

### Configuration

**Environment Variables Required:**
- `DATABASE_URL` - PostgreSQL connection string
- `API_FOOTBALL_KEY` - API-Football API key
- `ODDS_API_KEY` - The Odds API key
- `ML_SERVICE_URL` - ML service URL (default: http://localhost:8001)

See `.env.example` files in each service directory.

## MVP Implementation Phases

### ✅ Phase 1: Project Setup (Week 1) - COMPLETED
- Project structure established
- Database schema created
- Basic API servers running
- Docker Compose configured

### 🚧 Phase 2: Data Infrastructure (Weeks 1-2) - NEXT
- Implement API-Football integration (`backend/pkg/apifootball/`)
- Implement The Odds API integration (`backend/pkg/oddsapi/`)
- **Fetch odds for multiple markets** (1X2, Over/Under, BTTS, etc.)
- Create data sync jobs (`backend/cmd/sync/`)
- Build backfill script for 3 seasons (`backend/cmd/backfill/`)
- Implement repository layer (`backend/internal/repository/`)

### Phase 3: ML Models - MULTI-MARKET (Weeks 3-5)
- **Week 3**: 1X2 model (Home/Draw/Away) - 55-60% accuracy target
- **Week 4**: Over/Under 2.5 Goals model - 60-65% accuracy target
- **Week 5**: BTTS model (Both Teams to Score) - 58-62% accuracy target
- Feature engineering pipeline (`ml-service/app/features/`)
- Market-specific models (`ml-service/app/models/`)
- Model evaluation and backtesting per market
- Multi-market prediction API

### Phase 4: Betting Engine & Market Selector (Weeks 6-7)
- **Smart Market Selector** - evaluates ALL markets, picks highest EV
- Expected Value calculator (works across all markets)
- Kelly Criterion stake sizing (configurable per market)
- Multi-market weekly picks generation
- Bet tracking with market categorization

### Phase 5: Accumulator Builder (Week 7)
- **Smart Accumulator Generator** - combines 2-3 uncorrelated picks
- Correlation detection (avoid same fixture, same team)
- Accumulator EV calculation (combined probabilities)
- Conservative Kelly (1/8 for accumulators vs 1/4 for singles)
- 2-3 accumulator recommendations per week
- Max 20% of weekly stake on accumulators

### Phase 6: Dashboard (Weeks 7-8)
- Singles and accumulators display
- Performance charts (Recharts)
- Accumulator legs breakdown
- Bet management UI (singles + parlays)
- Real-time metrics

### Phase 7: Testing (Weeks 8-9)
- Integration tests
- Paper trading (singles + accumulators)
- Investor demo package

## Important Notes

### Betting Strategy

**Singles (15-20 per week):**
- **Kelly Fraction**: 1/4 Kelly (configurable per market)
- **Min EV Threshold**: 3% (filters low-value bets)
- **Max Bet Size**: 5% of bankroll (safety cap)
- **Market Focus**: Multi-market approach
  - **Week 3**: 1X2 (Home/Draw/Away)
  - **Week 4**: + Over/Under 2.5 Goals
  - **Week 5**: + Both Teams to Score (BTTS)
  - **Week 6+**: Smart market selector picks highest EV across all markets
- **Market Diversification**: Max 40% of weekly picks from single market type

**Accumulators (2-3 per week, starting Week 7):**
- **Kelly Fraction**: 1/8 Kelly (more conservative than singles)
- **Min EV Threshold**: 5% (higher than singles due to variance)
- **Num Legs**: 2-3 picks per accumulator
- **Max Stake Allocation**: 20% of weekly total stake
- **Correlation Detection**: Avoid same fixture, same team
- **Example**: Arsenal Over 2.5 + Brighton BTTS + Liverpool Away
  - Combined odds: 1.80 × 1.70 × 2.20 = 6.73
  - Combined probability: 24.3%
  - Stake: $75 (1/8 Kelly)

### Data Sources
- **API-Football**: Fixtures, standings, team stats (~$50-100/month)
- **The Odds API**: Bookmaker odds (~$200/month)
- **Total API Cost**: ~$250-300/month

### Model Details

**Multi-Market Models:**

1. **1X2 Model** (Home/Draw/Away)
   - Algorithm: XGBoost multi-class classification (3 outcomes)
   - Features: 15-20 features (form, H2H, position, goals, xG)
   - Target Accuracy: 55-60%
   - Training Data: 3 seasons (~1,140 matches)

2. **Over/Under 2.5 Model**
   - Algorithm: XGBoost binary classification
   - Features: Goal averages, scoring patterns, defensive metrics
   - Target Accuracy: 60-65%
   - Most popular totals market

3. **BTTS Model** (Both Teams to Score)
   - Algorithm: XGBoost binary classification
   - Features: Clean sheets, offensive/defensive strength
   - Target Accuracy: 58-62%
   - Independent of match outcome

**Shared Requirements:**
- Training Data: 3 seasons (~1,140 matches)
- Probability Calibration: Required for accurate EV calculations
- Time-series cross-validation

**Future Markets:**
- Double Chance (can derive from 1X2 model)
- Other totals (O/U 1.5, 3.5, 4.5)
- Half-time markets
- Handicaps
- Correct score
- See `docs/market-expansion-roadmap.md` for full plan

### Development Workflow
1. Data flows from API-Football/Odds API → PostgreSQL
   - Fixtures, team stats, odds for ALL markets
2. ML service reads from PostgreSQL for training
   - Trains separate models for each market
3. Go backend calls ML service for multi-market predictions
   - Requests predictions for: 1X2, O/U 2.5, BTTS, etc.
4. **Smart Market Selector** evaluates all markets
   - Calculates EV for every market/outcome combination
   - Ranks by EV, filters by 3% minimum threshold
   - Returns top recommendations across ALL markets
5. Frontend displays multi-market picks and performance
   - Shows best picks regardless of market type
   - Performance breakdown by market

## Next Steps

1. Implement API integrations (API-Football, The Odds API)
2. Build data sync jobs and backfill historical data
3. Develop feature engineering pipeline
4. Train and evaluate XGBoost model
5. Implement betting recommendation engine
6. Build investor dashboard UI

For detailed setup instructions, see [docs/setup-guide.md](docs/setup-guide.md).
