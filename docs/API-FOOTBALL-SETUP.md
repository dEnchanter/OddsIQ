# API-Football Setup & Testing Guide

**Status:** ✅ API Tested and Working (2026-01-15)

## Important Update: Odds Availability

**Free Tier Limitation Discovered:**
- ✅ Odds available for **current season only** (2025)
- ❌ Historical odds NOT available (2022-2024)
- ✅ This is sufficient for MVP with adjusted approach

**What API-Football Provides (Free Tier):**
- ✅ Fixtures (all seasons 2010-2025)
- ✅ Teams (all seasons)
- ✅ Standings (all seasons)
- ✅ **Pre-match odds** (current season only)
- ✅ **Live odds** (current season only)
- ✅ **Multiple bookmakers** (31 available)

---

## Step 1: Get Your API Key

1. Go to: https://dashboard.api-football.com/
2. Register/Login
3. Copy your API key
4. Check your quota (free tier usually gives 100 requests/day)

---

## Step 2: Create .env File

```bash
cd backend

# Copy the example
cp .env.example .env

# Edit .env and add your API key
DATABASE_URL=postgresql://user:pass@localhost:5432/oddsiq
API_FOOTBALL_KEY=your_actual_api_key_here
```

---

## Step 3: Test Your API Access

Run the test script to see what endpoints you have access to:

```bash
cd backend

# Run the test
go run cmd/test-api/main.go
```

**Expected Output:**
```
=== Testing API-Football Endpoints ===

1. Checking API Status & Quota...
✅ Success! Status: 200, Results: 1
   Remaining requests: 95/100

2. Getting Premier League info...
✅ Success! Status: 200, Results: 1

3. Getting upcoming Premier League fixtures...
✅ Success! Status: 200, Results: 5

4. Checking Odds endpoint...
✅ Success! Status: 200, Results: ...
   (This shows if odds are available)

5. Checking Bookmakers endpoint...
✅ Success! Status: 200, Results: ...

6. Checking Bets endpoint...
✅ Success! Status: 200, Results: ...
```

---

## Step 4: Understand What's Available

### API-Football Endpoints We Use

#### 1. **Fixtures** `/fixtures`
```bash
# Get upcoming Premier League fixtures
GET /fixtures?league=39&season=2024&next=10

# Get fixtures by date
GET /fixtures?date=2024-12-25

# Get fixture by ID
GET /fixtures?id=12345
```

#### 2. **Teams** `/teams`
```bash
# Get all Premier League teams
GET /teams?league=39&season=2024

# Get team by ID
GET /teams?id=42
```

#### 3. **Standings** `/standings`
```bash
# Get current Premier League table
GET /standings?league=39&season=2024
```

#### 4. **Odds** `/odds` ⭐ NEW!
```bash
# Get pre-match odds for fixtures
GET /odds?fixture=12345

# Get odds for league
GET /odds?league=39&season=2024

# Get live odds
GET /odds/live?fixture=12345
```

#### 5. **Bookmakers** `/odds/bookmakers`
```bash
# Get list of all bookmakers
GET /odds/bookmakers
```

#### 6. **Bet Types** `/odds/bets`
```bash
# Get available bet types
GET /odds/bets
```

---

## Step 5: Common Bet Types Available

API-Football typically provides these markets:

**Match Odds:**
- `Match Winner` - 1X2 (Home/Draw/Away)
- `Double Chance`
- `Draw No Bet`

**Goals:**
- `Over/Under 2.5 Goals`
- `Over/Under 1.5 Goals`
- `Over/Under 3.5 Goals`
- `Both Teams to Score` (BTTS)

**Handicap:**
- `Asian Handicap`
- `European Handicap`

**Others:**
- `Correct Score`
- `Half Time / Full Time`
- `Total Goals - Home Team`
- `Total Goals - Away Team`

---

## Step 6: Check Free Tier Limitations

**Typical Free Tier (may vary):**
- 100 requests per day
- All endpoints available
- All competitions available
- Historical data access
- Real-time updates (15-second delay)

**Run this to check YOUR limits:**
```bash
go run cmd/test-api/main.go
```

Look for the "paging" or "requests" section in the response.

---

## Step 7: Update Our Implementation

Since you have API-Football but not The Odds API, we need to:

### Option A: Use API-Football for Everything (Recommended)

**Pros:**
- ✅ One API key needed
- ✅ Simpler setup
- ✅ Fewer rate limit concerns
- ✅ Odds are matched to fixtures automatically

**Cons:**
- ⚠️ Fewer bookmakers than specialized odds APIs
- ⚠️ Higher request usage (fixtures + odds = 2 requests)

### Option B: Use Free Odds Sources Later

**For MVP:** Just use API-Football odds
**For Production:** Consider paid odds API if needed

---

## Updated Data Flow

### Old Plan (2 APIs):
```
API-Football → Fixtures, Teams, Standings
The Odds API → Odds (multiple bookmakers)
```

### New Plan (1 API):
```
API-Football → Fixtures, Teams, Standings, Odds
```

---

## Code Changes Needed

I'll update the codebase to:

1. ✅ Add odds endpoints to `pkg/apifootball/` client
2. ✅ Remove dependency on The Odds API (optional)
3. ✅ Update sync services to use API-Football odds
4. ✅ Update documentation

**These changes preserve the option to add The Odds API later if you get a key.**

---

## Rate Limit Management

With 100 requests/day on free tier:

**Daily Budget:**
- Status check: 1 request
- Get teams (once): 1 request
- Get fixtures (daily): 1 request
- Get standings (weekly): 1/7 request
- Get odds: ~50-60 requests

**Strategy:**
- Cache fixture and team data (doesn't change often)
- Only fetch odds for upcoming fixtures (next 7 days)
- Run sync once per day instead of hourly
- Use database to store historical data

---

## Testing Checklist

Run these tests to verify everything works:

```bash
# 1. Test API connection
go run cmd/test-api/main.go

# 2. Test backfill (with small dataset first)
go run cmd/backfill/main.go -seasons 2024 -teams-only

# 3. Check database
psql -d oddsiq -c "SELECT COUNT(*) FROM teams;"

# 4. Test fixtures backfill
go run cmd/backfill/main.go -seasons 2024 -fixtures-only

# 5. Check fixtures
psql -d oddsiq -c "SELECT COUNT(*) FROM fixtures WHERE season = 2024;"
```

---

## Next Steps

Once you confirm the test script works:

1. **I'll update the code** to add odds endpoints to API-Football client
2. **You test** the updated backfill with odds data
3. **We verify** odds are being stored correctly
4. **Continue** with Phase 3 (ML Model)

---

## Troubleshooting

### ❌ "Invalid API Key" Error
- Check your .env file has correct key
- Verify key at https://dashboard.api-football.com/
- No quotes around the key in .env

### ❌ "Rate Limit Exceeded"
- Wait 24 hours for quota reset
- Use cached data from database
- Consider upgrading plan if needed

### ❌ "Odds endpoint returns 0 results"
- Not all fixtures have odds available
- Odds appear closer to match time (usually 2-3 days before)
- Try fixtures in next 7 days

### ❌ "401 Unauthorized"
- Wrong header format
- Should be: `x-apisports-key: your_key`
- NOT: `x-rapidapi-key` (old v2 format)

---

## Useful Links

- **Dashboard:** https://dashboard.api-football.com/
- **Documentation:** https://www.api-football.com/documentation-v3
- **Support:** Check dashboard for support options
- **Status:** Check if API is down: https://status.api-football.com/

---

## Test Results (2026-01-15)

✅ **API test completed successfully!**

**Results:**
- ✅ API Status: Working
- ✅ Premier League: Accessible (League ID: 39)
- ✅ Fixtures: 16 seasons available (2010-2025)
- ✅ Bookmakers: 31 available
- ✅ Bet Types: 327 available
- ⚠️ **Odds: Current season only (2025)**

**Key Finding:**
- Historical seasons (2022-2024): `"odds": false`
- Current season (2025): `"odds": true`

**Implication:**
- Can train ML model on historical match outcomes (2022-2024)
- Can apply predictions to current odds (2025 season)
- Will build our own odds database going forward

See `docs/API-TEST-RESULTS.md` for detailed analysis and updated MVP approach.

---

## Summary

✅ You have everything you need with just API-Football!

**Your API is configured and working. Ready to proceed with Phase 2 implementation!** 🚀
