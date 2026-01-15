# OddsIQ - Next Steps After API Setup

**Date:** 2026-01-15
**Current Status:** ✅ API-Football tested and working | ✅ Odds endpoints added

---

## What We Just Completed

### 1. API-Football Setup ✅
- Created `.env` file with your API key
- Fixed test script path issue
- Successfully tested all endpoints
- Discovered free tier capabilities and limitations

### 2. Test Results ✅
- **API Status**: Working perfectly
- **Premier League Access**: League ID 39, 16 seasons available
- **Bookmakers**: 31 available
- **Bet Types**: 327 available
- **Odds Availability**: Current season only (2025)

### 3. Code Updates ✅
- Added odds endpoints to API-Football client (`pkg/apifootball/odds.go`)
- Updated backfill script to include season 2025
- Created comprehensive test results documentation

---

## Updated MVP Strategy

Based on API test results, here's our revised approach:

### Training Data (Seasons 2022-2024)
- ✅ Historical match results available
- ✅ Team statistics available
- ✅ Standings available
- ❌ Historical odds NOT available (free tier limitation)
- **Strategy**: Train ML model on match outcomes without odds

### Live Data (Season 2025)
- ✅ Current fixtures available
- ✅ Live odds available
- ✅ 31 bookmakers accessible
- ✅ 327 bet types available
- **Strategy**: Apply predictions to current odds, start paper trading

### Going Forward
- Store all odds we fetch to build our own historical database
- Track closing lines for future model improvements
- Add The Odds API later if needed for more bookmakers

---

## Phase 2: Data Infrastructure - Immediate Next Steps

### Step 1: Setup Database (If Not Already Done)

```bash
# Create database
psql -U postgres -c "CREATE DATABASE oddsiq;"

# Run migrations
cd C:\Users\afolabi.opaleye\Desktop\builds\personal-builds\AI-builds\OddsIQ
psql -U postgres -d oddsiq -f database/migrations/001_initial_schema.up.sql
psql -U postgres -d oddsiq -f database/migrations/002_add_betting_tables.up.sql
psql -U postgres -d oddsiq -f database/migrations/003_add_accumulators.up.sql
```

**Verify database:**
```bash
psql -U postgres -d oddsiq -c "\dt"
```

You should see tables: teams, fixtures, odds, team_stats, predictions, bets, etc.

---

### Step 2: Test Data Backfill (Small Test First)

Let's start by testing with just teams and one season:

```bash
cd backend

# Test 1: Load teams only for 2025 season
go run cmd/backfill/main.go -seasons 2025 -teams-only

# Check what was loaded
psql -U postgres -d oddsiq -c "SELECT COUNT(*) FROM teams;"
psql -U postgres -d oddsiq -c "SELECT name, code, venue_name FROM teams LIMIT 10;"
```

**Expected**: Should see ~20 Premier League teams loaded.

---

### Step 3: Load Fixtures for One Season

```bash
# Test 2: Load fixtures for 2025 season
go run cmd/backfill/main.go -seasons 2025 -fixtures-only

# Check fixtures
psql -U postgres -d oddsiq -c "SELECT COUNT(*) FROM fixtures WHERE season = 2025;"
psql -U postgres -d oddsiq -c "SELECT home_team_name, away_team_name, match_date, status FROM fixtures WHERE season = 2025 LIMIT 10;"
```

**Expected**: Should see fixtures for 2025 season (season in progress).

---

### Step 4: Load All Historical Data

Once steps 2-3 work, load all seasons:

```bash
# Load all data (2022-2025)
go run cmd/backfill/main.go

# Check results
psql -U postgres -d oddsiq -c "
  SELECT
    season,
    COUNT(*) as fixture_count
  FROM fixtures
  GROUP BY season
  ORDER BY season;
"
```

**Expected Output:**
```
 season | fixture_count
--------+---------------
   2022 |           380
   2023 |           380
   2024 |           380
   2025 |           ~200  (season in progress)
```

---

### Step 5: Create Odds Sync Service

Now we need to create a service to fetch odds for current season fixtures.

**File to create:** `backend/internal/services/odds_sync_apifootball.go`

This service will:
1. Get upcoming fixtures for next 7 days (2025 season)
2. Fetch odds for each fixture from API-Football
3. Store odds in database with bookmaker and bet type details
4. Run daily via cron job

**Key considerations:**
- Only fetch odds for fixtures in next 7 days (rate limit management)
- Store all bookmaker odds (31 bookmakers available)
- Focus on key markets: Match Winner (1X2), Over/Under 2.5, BTTS
- Mark closing odds when match starts

---

### Step 6: Test Odds Fetching

```bash
# Get a fixture ID from the database
psql -U postgres -d oddsiq -c "
  SELECT id, api_football_id, home_team_name, away_team_name, match_date
  FROM fixtures
  WHERE season = 2025
    AND status = 'Not Started'
    AND match_date > NOW()
  LIMIT 1;
"

# Note the api_football_id, then test fetching odds
# (We'll create a test script for this)
```

---

## Files to Create Next

### 1. Odds Sync Service
**File:** `backend/internal/services/odds_sync_apifootball.go`

**Purpose:**
- Fetch current odds for upcoming fixtures
- Store in database with proper bookmaker/market tracking
- Handle multiple markets (1X2, O/U, BTTS, etc.)

### 2. Test Script for Odds
**File:** `backend/cmd/test-odds/main.go`

**Purpose:**
- Test fetching odds for a specific fixture
- Verify odds are being parsed correctly
- Check database storage

### 3. Daily Sync Script
**File:** `backend/cmd/daily-sync/main.go`

**Purpose:**
- Run daily to sync current fixtures
- Fetch odds for upcoming matches
- Update fixture statuses
- Can be scheduled via cron

---

## Rate Limit Management Strategy

With 100 requests/day on free tier:

**Daily Budget:**
- Status check: 1 request
- Get standings (weekly): ~1 request
- Get fixtures for next 7 days: 1 request
- **Get odds for each fixture: ~50-60 requests**
- Total: ~55-65 requests/day ✅ (within limit)

**Optimization:**
- Cache teams (update monthly)
- Cache standings (update weekly)
- Only fetch odds for fixtures in next 7 days
- Store everything in database
- Run once per day (not hourly)

---

## Success Criteria for Phase 2 Completion

- [ ] Database created with all tables
- [ ] Teams loaded for all seasons (2022-2025)
- [ ] Fixtures loaded for all seasons (2022-2025)
- [ ] Odds sync service created
- [ ] Odds being fetched for upcoming 2025 fixtures
- [ ] Odds stored in database with bookmaker details
- [ ] Daily sync running successfully
- [ ] API rate limits not exceeded

---

## Timeline

**This Week (Week 2):**
- ✅ API setup and testing (DONE)
- ✅ Odds endpoints added (DONE)
- [ ] Database setup
- [ ] Data backfill (teams + fixtures)
- [ ] Odds sync service implementation
- [ ] Daily sync testing

**Next Week (Week 3):**
- [ ] Begin Phase 3: ML Model Development
- [ ] Feature engineering
- [ ] Model training on historical data
- [ ] First predictions

---

## Commands Quick Reference

```bash
# Setup database
psql -U postgres -c "CREATE DATABASE oddsiq;"
psql -U postgres -d oddsiq -f database/migrations/001_initial_schema.up.sql

# Load data
cd backend
go run cmd/backfill/main.go                    # Load all seasons
go run cmd/backfill/main.go -seasons 2025      # Load just 2025
go run cmd/backfill/main.go -teams-only        # Load only teams

# Check data
psql -U postgres -d oddsiq -c "SELECT COUNT(*) FROM teams;"
psql -U postgres -d oddsiq -c "SELECT COUNT(*) FROM fixtures;"
psql -U postgres -d oddsiq -c "SELECT season, COUNT(*) FROM fixtures GROUP BY season;"

# Start API server
cd backend
go run cmd/api/main.go

# Test endpoints
curl http://localhost:8000/health
curl http://localhost:8000/api/fixtures?season=2025
```

---

## Need Help?

- **API Documentation**: `docs/API-FOOTBALL-SETUP.md`
- **Test Results**: `docs/API-TEST-RESULTS.md`
- **Go Learning Guide**: `docs/GO-LEARNING-GUIDE.md`
- **Database Schema**: `docs/database-schema.md`
- **API Specification**: `docs/api-specification.md`

---

## What to Do Right Now

1. **Setup Database** (if not already done)
   ```bash
   psql -U postgres -c "CREATE DATABASE oddsiq;"
   psql -U postgres -d oddsiq -f database/migrations/001_initial_schema.up.sql
   ```

2. **Test Backfill Script**
   ```bash
   cd backend
   go run cmd/backfill/main.go -seasons 2025 -teams-only
   ```

3. **Verify Data Loaded**
   ```bash
   psql -U postgres -d oddsiq -c "SELECT name FROM teams LIMIT 10;"
   ```

4. **Report Results**
   - Share any errors you encounter
   - Confirm number of teams loaded
   - Ready to proceed to odds sync implementation

---

**Let's get the data infrastructure fully operational!** 🚀
