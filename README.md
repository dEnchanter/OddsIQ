# OddsIQ - AI/ML Sports Betting Prediction System

An AI-powered sports betting prediction system focused on Premier League football, designed to identify value betting opportunities using machine learning and advanced analytics.

## Project Structure

```
OddsIQ/
├── backend/              # Go backend service
├── ml-service/          # Python ML service
├── frontend/            # Next.js dashboard
├── database/            # Database migrations
├── docs/                # Documentation
└── docker-compose.yml   # Docker orchestration
```

## Tech Stack

- **Backend**: Go 1.21+, Gin, pgx
- **ML Service**: Python 3.11+, FastAPI, XGBoost, scikit-learn
- **Frontend**: Next.js 14+, TypeScript, Tailwind CSS
- **Database**: PostgreSQL 15+

## Quick Start

### Prerequisites

- Docker & Docker Compose (recommended)
- OR: Go 1.21+, Python 3.11+, Node.js 18+, PostgreSQL 15+

### Option 1: Docker (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd OddsIQ

# Start all services
docker-compose up -d

# Check service health
docker-compose ps

# View logs
docker-compose logs -f
```

Services will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8000
- ML Service: http://localhost:8001
- PostgreSQL: localhost:5432

### Option 2: Local Development

See [docs/setup-guide.md](docs/setup-guide.md) for detailed setup instructions.

## Documentation

- [Setup Guide](docs/setup-guide.md) - Complete setup instructions
- [Database Schema](docs/database-schema.md) - Database design and queries
- [API Specification](docs/api-specification.md) - API endpoints and contracts
- [Architecture Decisions](docs/architecture-decisions.md) - Key technical decisions
- [CLAUDE.md](CLAUDE.md) - AI assistant guidance

## Development Workflow

### Running Locally

```bash
# Terminal 1 - PostgreSQL
docker-compose up postgres

# Terminal 2 - Go Backend
cd backend
go run cmd/api/main.go

# Terminal 3 - Python ML Service
cd ml-service
source venv/bin/activate
uvicorn app.main:app --reload --port 8001

# Terminal 4 - Next.js Frontend
cd frontend
npm run dev
```

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# ML Service tests
cd ml-service
pytest

# Frontend tests
cd frontend
npm test
```

## MVP Roadmap

### Phase 1: Project Setup (Week 1) ✅
- [x] Project structure
- [x] Database schema
- [x] Development environment

### Phase 2: Data Infrastructure (Weeks 1-2) 🚧
- [ ] API-Football integration
- [ ] The Odds API integration
- [ ] Historical data backfill
- [ ] Data sync jobs

### Phase 3: ML Models - Multi-Market (Weeks 3-5)
- [ ] Week 3: 1X2 model (Home/Draw/Away) - 55-60% accuracy
- [ ] Week 4: Over/Under 2.5 Goals model - 60-65% accuracy
- [ ] Week 5: BTTS model (Both Teams to Score) - 58-62% accuracy
- [ ] Feature engineering pipeline
- [ ] Model evaluation and backtesting per market

### Phase 4: Betting Engine & Market Selector (Weeks 6-7)
- [ ] Smart Market Selector (evaluates ALL markets)
- [ ] EV calculation across markets
- [ ] Kelly Criterion stake sizing (per market)
- [ ] Multi-market weekly picks generation
- [ ] Bet tracking with market categorization

### Phase 5: Dashboard (Weeks 5-6)
- [ ] Picks display
- [ ] Performance metrics
- [ ] Charts and visualization
- [ ] Bet management UI

### Phase 6: Testing (Weeks 7-8)
- [ ] Integration testing
- [ ] Paper trading
- [ ] Investor demo package

## Environment Variables

Create `.env` files in each service directory based on `.env.example`:

```bash
# Backend
cp backend/.env.example backend/.env

# ML Service
cp ml-service/.env.example ml-service/.env

# Frontend
cp frontend/.env.local.example frontend/.env.local
```

Update with your API keys and configuration.

## API Keys Required

- **API-Football**: Get key from https://www.api-football.com/
- **The Odds API**: Get key from https://the-odds-api.com/

## Contributing

1. Create feature branch: `git checkout -b feature/your-feature`
2. Make changes and test locally
3. Commit: `git commit -m 'Add your feature'`
4. Push: `git push origin feature/your-feature`
5. Create Pull Request

## License

Proprietary - All rights reserved

## MVP Success Criteria

- [ ] 3 seasons of Premier League data loaded (1,140+ matches)
- [ ] Odds data for multiple markets (1X2, O/U, BTTS) stored
- [ ] **3 ML models trained:**
  - [ ] 1X2: 55-60% accuracy
  - [ ] Over/Under 2.5: 60-65% accuracy
  - [ ] BTTS: 58-62% accuracy
- [ ] **Smart Market Selector operational** (picks best EV across all markets)
- [ ] **Accumulator Builder operational** (2-3 leg parlays)
- [ ] Positive theoretical ROI: +6-10% (singles + accumulators)
- [ ] Weekly picks:
  - [ ] 15-20 single bets from multiple markets
  - [ ] 2-3 accumulators (2-3 legs each)
- [ ] Functional investor dashboard:
  - [ ] Singles and accumulators display
  - [ ] Performance breakdown by bet type
  - [ ] Market breakdown
- [ ] 8+ weeks of paper trading results (singles + accumulators)

## Support

For setup issues, see [docs/setup-guide.md](docs/setup-guide.md).

For architecture questions, see [docs/architecture-decisions.md](docs/architecture-decisions.md).
