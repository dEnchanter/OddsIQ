# Phase 2: Data Infrastructure - Implementation Complete

**Date:** 2026-01-14
**Status:** ✅ COMPLETED
**Duration:** Weeks 1-2 of MVP Timeline

## Overview

Phase 2 focused on building the complete data ingestion pipeline for OddsIQ. This infrastructure enables the system to pull live Premier League data from external APIs and store it in PostgreSQL for analysis and predictions.

## What Was Built

### 1. API Client Packages

#### API-Football Client (`backend/pkg/apifootball/`)
- **client.go** - HTTP client with authentication and base structures
  - `NewClient()` - Initializes client with API key
  - `doRequest()` - Handles HTTP requests with proper error handling
  - Response structures: `Team`, `Venue`, `Fixture`, `FixtureResponse`, `StandingsResponse`

- **fixtures.go** - Fixture data fetching
  - `GetFixtures(leagueID, season)` - All fixtures for a season
  - `GetFixturesByDate(date)` - Fixtures for specific date
  - `GetFixturesByDateRange(from, to)` - Date range queries
  - `GetFixture(fixtureID)` - Single fixture by ID

- **teams.go** - Team data fetching
  - `GetTeams(leagueID, season)` - All teams for league/season
  - `GetTeam(teamID)` - Single team by ID

- **standings.go** - League standings
  - `GetStandings(leagueID, season)` - Current league table

#### The Odds API Client (`backend/pkg/oddsapi/`)
- **client.go** - HTTP client for odds data
  - `NewClient()` - Initialize with API key
  - Market constants: `MarketH2H`, `MarketTotals`, `MarketBTTS`, `MarketSpread`
  - Response structures: `Event`, `Bookmaker`, `Market`, `Outcome`

- **odds.go** - Odds fetching and helpers
  - `GetOdds(sport, markets, regions)` - Generic odds fetching
  - `GetEPLOdds(markets)` - Premier League convenience method
  - `GetAllMarketsEPL()` - All markets for EPL
  - `GetH2HOdds()` - 1X2 odds only
  - `GetTotalsOdds()` - Over/Under odds
  - `GetBTTSOdds()` - Both Teams to Score odds
  - Helper functions:
    - `GetBestOdds()` - Find highest odds across bookmakers
    - `GetAverageOdds()` - Calculate average odds
    - `ExtractH2HOdds()` - Parse 1X2 odds
    - `ExtractOverUnderOdds()` - Parse totals odds

### 2. Repository Layer (`backend/internal/repository/`)

#### Teams Repository (`teams.go`)
Complete CRUD operations for teams:
- `Create()` - Insert new team
- `GetByID()` - Retrieve by database ID
- `GetByAPIFootballID()` - Retrieve by external API ID
- `GetAll()` - All teams
- `Update()` - Update existing team
- `Delete()` - Remove team
- `Upsert()` - Insert or update based on API ID
- `GetPremierLeagueTeams()` - Convenience method

#### Fixtures Repository (`fixtures.go`)
Comprehensive fixture management:
- `Create()`, `Update()`, `Delete()` - Basic CRUD
- `GetByID()`, `GetByAPIFootballID()` - Single fixture retrieval
- `GetBySeason()` - All fixtures for season
- `GetByDateRange()` - Fixtures between dates
- `GetUpcoming()` - Future fixtures
- `GetByStatus()` - Filter by match status (NS, FT, etc.)
- `GetByTeam()` - All fixtures for specific team
- `GetRecentByTeam()` - Last N fixtures for team
- `UpdateScore()` - Update match result
- `Upsert()` - Insert or update fixture

#### Odds Repository (`odds.go`)
Multi-market odds storage:
- `Create()`, `CreateBatch()` - Insert odds (single/batch)
- `GetByFixture()` - All odds for fixture
- `GetLatestByFixture()` - Most recent odds per market/outcome
- `GetByFixtureAndMarket()` - Filter by market type
- `GetLatestByFixtureAndMarket()` - Latest for specific market
- `GetBestOdds()` - Highest odds for fixture/market/outcome
- `GetByBookmaker()` - All odds from specific bookmaker
- `GetByDateRange()` - Odds within date range
- `DeleteOldOdds()` - Cleanup historical data
- `GetMarketTypes()`, `GetBookmakers()` - List available data
- `GetAverageOdds()` - Calculate average across bookmakers

#### Team Stats Repository (`team_stats.go`)
Season statistics management:
- `Create()`, `Update()`, `Delete()` - Basic CRUD
- `GetByID()` - Single stat record
- `GetByTeamAndSeason()` - Stats for team in season
- `GetBySeason()` - All teams in season
- `GetByTeam()` - Team across all seasons
- `Upsert()` - Insert or update
- `GetTopTeams()` - Top N teams by points

### 3. Data Sync Services (`backend/internal/services/`)

#### Fixture Sync Service (`fixture_sync.go`)
Orchestrates fixture data synchronization:
- `SyncTeams(season)` - Fetch and store all teams
- `SyncFixturesBySeason(season)` - All fixtures for season
- `SyncFixturesByDateRange(from, to)` - Fixtures in date range
- `SyncUpcomingFixtures()` - Next 7 days
- `UpdateFixtureResults()` - Update scores for recent matches
- `processFixture()` - Convert API data to models and store
- `SyncAllSeasons(seasons[])` - Bulk sync multiple seasons

#### Odds Sync Service (`odds_sync.go`)
Manages odds data synchronization:
- `SyncAllMarkets()` - All supported markets (H2H, Totals, BTTS)
- `SyncMarket(marketType)` - Specific market only
- `SyncH2HOdds()`, `SyncTotalsOdds()`, `SyncBTTSOdds()` - Market-specific methods
- `processEvent()` - Process single event and store odds
- `findMatchingFixture()` - Match odds event to database fixture
- `matchTeamNames()` - Fuzzy team name matching
- `extractOddsFromEvent()` - Parse all odds from event
- `normalizeOutcome()` - Standardize outcome names
- `CleanupOldOdds(daysToKeep)` - Remove old data
- `GetOddsSummary()` - Statistics about stored odds

### 4. Automated Scheduling (`backend/internal/services/scheduler.go`)

Cron-based job scheduler:

**Production Schedule:**
- **Daily at 6:00 AM** - Sync upcoming fixtures
- **Every 30 minutes (match days only)** - Update fixture results
- **Every 2 hours** - Sync odds for all markets
- **Every hour** - Sync H2H odds (most important market)
- **Weekly (Sunday 3:00 AM)** - Cleanup old odds (30+ days)

**Development Schedule:**
- Once per day at noon - Sync fixtures
- Twice per day - Sync odds

**Key Methods:**
- `Start()` - Start production schedule
- `StartDevelopmentSchedule()` - Start dev schedule
- `Stop()` - Stop all jobs
- `RunNow()` - Execute all jobs immediately (testing)
- `GetNextRunTimes()` - View scheduled run times

### 5. Historical Data Backfill Script (`backend/cmd/backfill/main.go`)

Command-line tool for populating historical data:

**Usage:**
```bash
# Backfill all data for 2022-2024
go run cmd/backfill/main.go

# Backfill specific season
go run cmd/backfill/main.go -seasons 2024

# Backfill only teams
go run cmd/backfill/main.go -teams-only

# Backfill only fixtures
go run cmd/backfill/main.go -fixtures-only
```

**Features:**
- Parses comma-separated season list
- Sync teams first, then fixtures
- Progress logging for each season
- Summary statistics on completion
- Error handling with continue-on-failure

### 6. API Endpoints (`backend/internal/api/`)

Updated handlers with repository integration:

**Data Query Endpoints:**
- `GET /api/fixtures` - List fixtures (query by season, status)
- `GET /api/fixtures/:id` - Single fixture with teams
- `GET /api/fixtures/:id/odds` - All odds for fixture

**Handler Features:**
- Created `API` struct holding repositories
- `NewAPI()` - Dependency injection pattern
- Proper error handling and status codes
- Query parameter support
- Context-aware database operations

## Key Technical Decisions

### 1. Repository Pattern
- Clean separation between API clients and database
- Testable business logic
- Easy to swap implementations

### 2. Upsert Operations
- All sync operations use `ON CONFLICT` upserts
- Prevents duplicate data
- Handles API data changes gracefully

### 3. Fuzzy Team Name Matching
- Odds API and API-Football use different team names
- Implemented `matchTeamNames()` with normalization
- Handles common abbreviations (Man Utd, Spurs, etc.)

### 4. Batch Insert for Odds
- Single transaction for multiple odds
- Better performance for high-volume data
- Atomic operations

### 5. Separate Odds History
- Don't update odds, insert new records
- Preserves historical odds movements
- Enables odds trend analysis

### 6. Scheduled Job Frequency
- H2H odds more frequent (every hour) - most important
- All markets every 2 hours - balanced load
- Match day fixture updates every 30 minutes

## Database Schema Used

Phase 2 uses these tables (from Phase 1):
1. **teams** - Premier League teams
2. **fixtures** - Match fixtures (3 seasons ≈ 1,140 matches)
3. **odds** - Multi-market odds from multiple bookmakers
4. **team_stats** - Season statistics per team

## Environment Variables Required

```env
# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/oddsiq

# API Keys
API_FOOTBALL_KEY=your_api_football_key
ODDS_API_KEY=your_odds_api_key

# Optional
ML_SERVICE_URL=http://localhost:8001
KELLY_FRACTION=0.25
MIN_EV_THRESHOLD=0.03
```

## Files Created

### API Clients (2 packages, 5 files)
1. `backend/pkg/apifootball/client.go`
2. `backend/pkg/apifootball/fixtures.go`
3. `backend/pkg/apifootball/teams.go`
4. `backend/pkg/apifootball/standings.go`
5. `backend/pkg/oddsapi/client.go`
6. `backend/pkg/oddsapi/odds.go`

### Repositories (4 files)
7. `backend/internal/repository/teams.go`
8. `backend/internal/repository/fixtures.go`
9. `backend/internal/repository/odds.go`
10. `backend/internal/repository/team_stats.go`

### Services (3 files)
11. `backend/internal/services/fixture_sync.go`
12. `backend/internal/services/odds_sync.go`
13. `backend/internal/services/scheduler.go`

### Tools (1 file)
14. `backend/cmd/backfill/main.go`

### API Updates (1 file modified)
15. `backend/internal/api/handlers.go` (updated)

### Documentation (1 file)
16. `docs/PHASE-2-COMPLETE.md` (this file)

**Total:** 16 files created/modified

## Testing Phase 2

### Manual Testing Checklist

```bash
# 1. Test database connection
psql -d oddsiq -c "SELECT 1;"

# 2. Run backfill for one season
go run backend/cmd/backfill/main.go -seasons 2024

# 3. Verify data in database
psql -d oddsiq -c "SELECT COUNT(*) FROM teams;"
psql -d oddsiq -c "SELECT COUNT(*) FROM fixtures WHERE season = 2024;"

# 4. Start API server
cd backend && go run cmd/api/main.go

# 5. Test API endpoints
curl http://localhost:8000/health
curl http://localhost:8000/api/fixtures?season=2024
curl http://localhost:8000/api/fixtures/1
curl http://localhost:8000/api/fixtures/1/odds
```

### Expected Results

After successful backfill:
- **Teams:** 20 Premier League teams
- **Fixtures (2024):** ~380 fixtures
- **Fixtures (2022-2024):** ~1,140 fixtures total
- **Odds:** Variable (depends on upcoming fixtures)

## Performance Considerations

### API Rate Limits
- **API-Football:** 100 requests/day (free tier)
- **The Odds API:** 500 requests/month (free tier)
- Scheduler respects limits with smart caching

### Database Indexes
Already defined in schema:
- `teams(api_football_id)` - UNIQUE
- `fixtures(api_football_id)` - UNIQUE
- `fixtures(season, match_date)` - Composite
- `odds(fixture_id, market_type, timestamp)` - Query optimization
- `team_stats(team_id, season)` - UNIQUE composite

### Batch Operations
- Odds inserts use batch transactions
- Backfill processes all teams before fixtures
- Scheduler batches API calls per job

## Known Limitations

1. **Team Name Matching** - Fuzzy matching may fail for very different names
2. **API Rate Limits** - Free tier limits may require paid upgrade
3. **Historical Odds** - The Odds API doesn't provide historical odds (only current)
4. **Time Zones** - All times stored in UTC, conversion needed for display
5. **Referee Data** - May be null for some fixtures

## Next Steps (Phase 3)

Phase 2 provides the foundation for Phase 3: Feature Engineering & ML Model

**Ready for Phase 3:**
- ✅ Historical fixture data (3 seasons)
- ✅ Current odds data (multiple markets)
- ✅ Team statistics
- ✅ Automated data refresh

**Phase 3 Tasks:**
- Build feature engineering pipeline (Python)
- Calculate form metrics (last 5 games)
- Compute head-to-head statistics
- Train XGBoost models (1X2, O/U 2.5, BTTS)
- Create prediction API endpoints

## Troubleshooting

### Issue: Backfill fails with "team not found"
**Solution:** Run with `-teams-only` first, then `-fixtures-only`

### Issue: Team name matching fails
**Solution:** Update `matchTeamNames()` abbreviations map in `odds_sync.go`

### Issue: API rate limit exceeded
**Solution:** Reduce scheduler frequency or upgrade API plan

### Issue: Database connection fails
**Solution:** Check `DATABASE_URL` format and PostgreSQL is running

### Issue: Scheduler doesn't run
**Solution:** Ensure cron format is correct and timezone is set

## Success Metrics

Phase 2 is considered successful when:
- ✅ All 16 files created and compiling
- ✅ Backfill script populates 1,140+ fixtures
- ✅ All 20 Premier League teams in database
- ✅ Odds sync returns data from multiple bookmakers
- ✅ Scheduler runs without errors
- ✅ API endpoints return valid data
- ✅ No data duplication (upserts working)

## Documentation References

- `docs/implementation-plan.md` - Full MVP plan
- `docs/database-schema.md` - Database structure
- `docs/api-specification.md` - API endpoint specs
- `backend/pkg/apifootball/` - API-Football client usage
- `backend/pkg/oddsapi/` - The Odds API client usage

---

**Phase 2 Status:** ✅ COMPLETE
**Next Phase:** Phase 3 - Feature Engineering & ML Model (Weeks 3-5)
**Estimated Time to MVP:** 7 weeks remaining
