# API-Football Free Tier Limitations - Detailed Analysis

**Date:** 2026-01-15
**Account:** afolabiopaleye@gmail.com
**Plan:** Free (expires 2027-01-15)

---

## Summary of Findings

After extensive testing, we've discovered significant limitations on the free tier that impact our ability to fetch current odds:

### ✅ What Works (Free Tier)

1. **Historical Fixtures**
   - ✅ Can fetch fixtures from completed seasons (2022, 2023, 2024)
   - ✅ Can fetch by season parameter
   - ✅ Can fetch by date range
   - ✅ All fixture details included (teams, scores, venue, etc.)

2. **Teams Data**
   - ✅ All teams for any season
   - ✅ Full team details (venue, logo, etc.)

3. **Standings**
   - ✅ Historical standings available

4. **Bookmakers & Bet Types**
   - ✅ Can list all 31 bookmakers
   - ✅ Can list all 327 bet types

### ❌ What Doesn't Work (Free Tier)

1. **Current Season Fixtures**
   - ❌ Season 2025 returns 0 fixtures
   - ❌ `next` parameter not supported ("Free plans do not have access")
   - ❌ `last` parameter not supported ("Free plans do not have access")
   - ❌ Current date queries return 0 results

2. **Odds Data**
   - ❌ Historical odds (2022-2024): Not available (coverage shows `"odds": false`)
   - ❌ Current season odds (2025): Not available (0 results despite coverage showing `"odds": true`)
   - ❌ `/odds` endpoint returns 0 results for all queries tested

3. **Live Data**
   - ❌ Cannot access current/upcoming fixtures
   - ❌ Cannot access live odds
   - ❌ Real-time features restricted

---

## Test Results

### Test 1: Season-Based Queries
```bash
# Historical seasons - Works
GET /fixtures?league=39&season=2024
Result: ✅ 380 fixtures

GET /fixtures?league=39&season=2023
Result: ✅ 380 fixtures

# Current season - Fails
GET /fixtures?league=39&season=2025
Result: ❌ 0 fixtures
```

### Test 2: Date-Based Queries
```bash
# Historical dates - Works
GET /fixtures?league=39&season=2024&from=2025-05-01&to=2025-05-31
Result: ✅ 41 fixtures (end of 2024 season)

# Current dates - Fails
GET /fixtures?league=39&date=2026-01-15
Result: ❌ 0 fixtures

GET /fixtures?league=39&from=2026-01-16&to=2026-01-22
Result: ❌ 0 fixtures
```

### Test 3: Special Parameters
```bash
# Next parameter - Not allowed
GET /fixtures?league=39&next=10
Error: "Free plans do not have access to the Next parameter"

# Last parameter - Not allowed
GET /fixtures?league=39&last=5
Error: "Free plans do not have access to the Last parameter"
```

### Test 4: Odds Queries
```bash
# By fixture
GET /odds?fixture=1208386
Result: ❌ 0 results

# By league/season (historical)
GET /odds?league=39&season=2024
Result: ❌ 0 results (coverage shows "odds": false)

# By league/season (current)
GET /odds?league=39&season=2025
Result: ❌ 0 results (even though coverage shows "odds": true)
```

---

## API Coverage Analysis

From `/leagues?id=39` response:

### Season 2024 (Completed)
```json
{
  "year": 2024,
  "current": false,
  "coverage": {
    "odds": false,  // ❌ No odds
    "fixtures": {
      "events": true,
      "lineups": true,
      "statistics_fixtures": true,
      "statistics_players": true
    }
  }
}
```

### Season 2025 (Current)
```json
{
  "year": 2025,
  "current": true,
  "start": "2025-08-15",
  "end": "2026-05-24",
  "coverage": {
    "odds": true,  // ✅ Shows as available
    "fixtures": {
      "events": true,
      "lineups": true,
      "statistics_fixtures": true,
      "statistics_players": true
    }
  }
}
```

**Discrepancy:** Coverage shows `"odds": true` for season 2025, but all queries return 0 results.

---

## Rate Limits

**Daily Quota:** 100 requests/day
**Used Today:** 18/100 (after all testing)
**Remaining:** 82 requests

---

## Conclusions

### Finding 1: No Current Season Access
The free tier does **NOT** provide access to:
- Current/upcoming fixtures
- Live match data
- Current season odds

Even though the API coverage indicates odds are available for season 2025, actual queries return 0 results.

### Finding 2: Limited to Historical Data Only
The free tier is effectively **historical data only**:
- ✅ Completed seasons (fixtures, results, teams)
- ❌ Current season (no fixtures, no odds)
- ❌ Live/upcoming data

### Finding 3: No Odds Data at All
Despite API documentation suggesting odds are available:
- ❌ No historical odds (seasons 2022-2024)
- ❌ No current odds (season 2025)
- The free tier appears to have **no odds access whatsoever**

---

## Impact on OddsIQ MVP

### Original Plan
1. Train ML model on historical matches WITH odds (2022-2024)
2. Apply predictions to current fixtures and odds
3. Paper trade with live odds
4. Build track record for investors

### Reality with Free Tier
1. ✅ Can train ML model on historical match outcomes (1,140 matches)
2. ❌ Cannot access current fixtures for predictions
3. ❌ Cannot access any odds data (historical or current)
4. ❌ Cannot paper trade without odds

### Critical Blocker
**We cannot build the MVP with API-Football free tier alone** because:
- No access to current season data
- No access to odds data whatsoever
- Cannot generate live predictions or track results

---

## Alternative Solutions

### Option 1: Use API-Football for History + Web Scraping for Current
**Approach:**
- Use API-Football for historical training data (what we have now)
- Scrape current fixtures and odds from betting sites
- Build predictions and track results

**Pros:**
- Free solution
- Can still build and test ML model
- Can paper trade with real odds

**Cons:**
- Web scraping is fragile (sites change, can block)
- Legal/ToS concerns
- Maintenance burden

---

### Option 2: Upgrade API-Football Plan
**Cost:** Check pricing at https://dashboard.api-football.com/

**Pro Plan Benefits:**
- Current season fixtures
- Live odds data
- More bookmakers
- Higher rate limits

**Decision:** Depends on budget and priorities

---

### Option 3: Use The Odds API (Paid)
**Alternative odds provider:** https://the-odds-api.com/

**Pricing:**
- $59/month for live odds
- Multiple sports including football/soccer
- Good coverage of bookmakers

**Combine with:**
- API-Football (free) for historical fixtures
- The Odds API for current odds

---

### Option 4: Manual Data Collection for Demo
**For MVP/Demo:**
- Use historical data to train and backtest model
- Manually track 2-3 weekends of predictions
- Screenshot odds from betting sites
- Build proof-of-concept without live API

**Then:**
- Use demo to raise investment
- Upgrade to paid APIs once funded

---

### Option 5: Focus on Backtesting Only (Recommended for Now)
**Approach:**
1. ✅ Train ML model on 1,140 historical matches
2. ✅ Create synthetic odds based on market averages
3. ✅ Backtest strategy on historical data
4. ✅ Build dashboard showing theoretical performance
5. ⏸️ Defer live trading until paid API secured

**Benefits:**
- Can complete 80% of MVP with current data
- Proves concept works
- Shows potential to investors
- Minimal additional cost

**Missing:**
- Live predictions (can add later)
- Real odds comparison (can simulate)
- Paper trading (can do manually for demo)

---

## Recommended Path Forward

### Phase 2 Completion (Modified)
Since we cannot access current odds with free tier:

**✅ Completed:**
- Database setup
- Historical data loaded (1,140 matches, 24 teams)
- API infrastructure built
- Data sync services created

**❌ Cannot Complete:**
- Current odds sync (no access to current odds)
- Live fixtures sync (no access to current season)
- Daily automation for live data (nothing to sync)

**✅ Can Complete:**
- Create synthetic odds generator for backtesting
- Build ML model training pipeline
- Implement backtesting framework

### Phase 3: Move Forward with ML Model
**What we CAN do with current data:**
1. Train XGBoost model on 1,140 historical matches
2. Engineer features from historical data
3. Create backtesting framework
4. Generate theoretical performance metrics
5. Build dashboard with backtest results

**What we CANNOT do:**
1. Live predictions on current fixtures
2. Real-time odds comparison
3. Paper trading with live odds

### Recommendation: **Option 5** - Focus on Backtesting
- Complete ML model development
- Use synthetic/average odds for backtesting
- Build compelling demo with historical performance
- Upgrade to paid API later (after funding or validation)

---

## Next Steps

### Immediate (This Session)
1. **Document free tier limitations** (✅ This file)
2. **Update Phase 2 status** to reflect what's possible
3. **Decide on path forward** (recommend Option 5)
4. **Move to Phase 3** if continuing with backtesting approach

### Short Term (Next Session)
1. Create synthetic odds generator
2. Begin ML model development
3. Build backtesting framework
4. Measure historical performance

### Long Term (Future)
1. Evaluate paid API options
2. Consider web scraping for live data
3. Implement live prediction system when ready
4. Deploy to production with proper APIs

---

## Questions for Decision

1. **Budget:** Can you afford $59-100/month for live odds API?
2. **Timeline:** Is backtesting-only MVP acceptable for now?
3. **Strategy:** Should we build with free tier limitations or wait for funding?

**My recommendation:** Proceed with backtesting-focused MVP using historical data. This proves the concept works and creates compelling material for investors, without requiring paid APIs upfront.

---

## Files Updated
- `docs/API-FREE-TIER-LIMITATIONS.md` (this file)
- `docs/PHASE-2-STATUS.md` (needs update)
- `docs/implementation-plan.md` (needs update)
