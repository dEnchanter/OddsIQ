# OddsIQ MVP Roadmap - Revised Strategy

**Date Created:** 2026-01-15
**Status:** Active Development Plan
**Approach:** Build → Validate → Scale

---

## Executive Summary

### The Strategy

Build a **self-validated MVP** that proves profitability with real-world testing before approaching investors or paying for APIs.

**Three-Phase Approach:**

1. **Phase 1 (Weeks 1-5):** Build foundation with historical data - **$0 cost**
2. **Phase 2 (Weeks 6-18):** Real-world validation with web scraping - **$0 cost**
3. **Phase 3 (Week 19+):** Scale & automate with paid APIs - **$59-100/month**

### Why This Works

✅ **Proves concept before spending** - No paid APIs until validated
✅ **Builds real track record** - 12+ weeks of actual predictions
✅ **Tests market efficiency** - Does the model actually beat bookmakers?
✅ **Investor-ready proof** - Real P&L, not theoretical backtests
✅ **Low risk** - If model fails, spent $0 not $1000+

---

## Current Status

### What We Have ✅

- ✅ **Database:** PostgreSQL with all tables created
- ✅ **Historical Data:** 1,140 Premier League matches (2022-2024)
- ✅ **Teams:** 24 Premier League teams with full details
- ✅ **Infrastructure:** Go backend + API clients built
- ✅ **API:** Free API-Football account (100 requests/day)

### What We Discovered ⚠️

**API-Football Free Tier Limitations:**
- ❌ No current season fixtures (2025)
- ❌ No current or historical odds
- ❌ No live/upcoming match data
- ❌ Limited to historical completed matches only

**Impact:** Cannot build live prediction system with free tier alone

**Solution:** Use historical data for training, add web scraping for real-world testing

---

## Phase 1: Foundation (Weeks 1-5)

**Goal:** Build and validate ML model using historical data

**Cost:** $0

### Week 1-2: Feature Engineering

**Objective:** Extract predictive features from 1,140 historical matches

**Tasks:**
1. Set up Python ML service structure
2. Create feature engineering pipeline
3. Build feature extraction modules:
   - Team form (last 5 games, points, goals, xG)
   - Head-to-head statistics
   - Home/away performance splits
   - League position differentials
   - Goal scoring/conceding rates
   - Recent form trends

**Deliverables:**
- `ml-service/app/features/form_metrics.py`
- `ml-service/app/features/h2h_stats.py`
- `ml-service/app/features/league_position.py`
- `ml-service/app/features/team_stats.py`
- Feature dataset CSV with all 1,140 matches

**Success Criteria:**
- All features extracted for 1,140 matches
- No missing data
- Features validated for sanity (ranges, distributions)

---

### Week 3-4: Model Training & Backtesting

**Objective:** Train XGBoost model and validate on historical data

**Tasks:**
1. Train/test split (80/20 with time-series awareness)
2. Train XGBoost classifier for match outcomes
3. Probability calibration
4. Build synthetic odds generator (based on market averages)
5. Create backtesting framework
6. Run backtests on held-out data

**Deliverables:**
- `ml-service/app/models/xgboost_model.py`
- `ml-service/app/backtesting/synthetic_odds.py`
- `ml-service/app/backtesting/backtest_engine.py`
- Trained model file (`xgboost_v1.pkl`)
- Backtest results report

**Success Criteria:**
- Model accuracy: >55% on test set
- Backtested ROI: >5% over test period
- Sharpe ratio: >1.0
- Predictions are well-calibrated

---

### Week 5: ML API & Integration

**Objective:** Create prediction API and integrate with Go backend

**Tasks:**
1. Build FastAPI service for predictions
2. Create prediction endpoints:
   - `POST /predict` - Single fixture prediction
   - `POST /batch-predict` - Weekly batch predictions
   - `GET /model/metrics` - Model performance
3. Integrate Go backend with ML service
4. Build prediction storage in database
5. Create basic performance tracking

**Deliverables:**
- `ml-service/app/api/predictions.py`
- `backend/internal/services/ml_client.go`
- `backend/internal/services/prediction_tracker.go`
- API documentation
- Integration tests

**Success Criteria:**
- ML API running on port 8001
- Go backend can call ML service
- Predictions stored in database
- End-to-end prediction flow works

---

### Phase 1 Deliverables Summary

**Code:**
- ✅ Python ML service with trained model
- ✅ Feature engineering pipeline
- ✅ Backtesting framework
- ✅ FastAPI prediction service
- ✅ Go backend integration

**Documentation:**
- ✅ Model training report
- ✅ Backtest results analysis
- ✅ Feature importance analysis
- ✅ API documentation

**Validation:**
- ✅ Model achieves >55% accuracy on historical data
- ✅ Backtested strategy shows positive ROI
- ✅ System can generate predictions for new fixtures

---

## Phase 2: Real-World Testing (Weeks 6-18)

**Goal:** Validate model with 12+ weeks of real predictions on actual market odds

**Cost:** $0 (web scraping + manual tracking)

### Week 6: Build Scraping Infrastructure

**Objective:** Create web scrapers for current fixtures and odds

**Tasks:**
1. Research scraping targets:
   - **Fixtures:** BBC Sport, Premier League official site, ESPN
   - **Odds:** Oddschecker.com, Oddsportal.com
2. Build fixture scraper
3. Build odds scraper (multiple bookmakers)
4. Build results scraper
5. Create database storage for scraped data
6. Set up daily scraping schedule

**Deliverables:**
- `scraping-service/scrapers/fixtures_scraper.py`
- `scraping-service/scrapers/odds_scraper.py`
- `scraping-service/scrapers/results_scraper.py`
- `scraping-service/storage/scraper_db.py`
- `scraping-service/scheduler.py`
- Scraping documentation with legal considerations

**Technical Stack:**
- **BeautifulSoup** for static pages (BBC Sport)
- **Selenium** for dynamic pages (betting sites)
- **Rate limiting:** 2-5 seconds between requests
- **Error handling:** Robust fallbacks
- **Caching:** Aggressive to minimize requests

**Success Criteria:**
- Can scrape upcoming Premier League fixtures
- Can scrape odds from 3+ bookmakers
- Can update match results automatically
- Scraper runs daily without errors

---

### Weeks 7-18: Live Testing Period (12 weeks)

**Objective:** Generate real predictions, track against actual market odds, build proof

**Weekly Workflow:**

**Monday/Tuesday:**
1. Scraper fetches upcoming weekend fixtures
2. ML model generates predictions
3. Record predictions with confidence levels
4. Log opening odds from multiple bookmakers

**Wednesday-Friday:**
5. Scraper updates odds daily
6. Track odds movement
7. Identify value bets (prediction > odds implied probability)
8. Friday: Record closing odds

**Saturday-Sunday:**
9. Matches are played
10. Watch results come in

**Monday (next week):**
11. Scraper updates match results
12. Calculate P&L for each prediction
13. Update running bankroll
14. Analyze performance metrics
15. Adjust model if needed

**Tracking Spreadsheet Columns:**
```
| Week | Date | Home | Away | Prediction | Confidence | Opening_Odds | Closing_Odds | Result | Bet_Amount | P&L | Bankroll |
```

**What to Track:**
- **Predictions:** All model predictions (not just bets placed)
- **Accuracy:** Hit rate by confidence level
- **Value:** Did we identify mispriced odds?
- **Odds Movement:** How do odds change through the week?
- **Kelly Criterion:** What stake size would be optimal?
- **Results:** Actual P&L as if betting with $10,000 bankroll

**Deliverables (per week):**
- Weekly predictions CSV
- Odds tracking CSV
- Results CSV with P&L
- Weekly performance report

**Success Criteria (after 12 weeks):**
- ✅ 50+ predictions tracked
- ✅ >55% accuracy overall
- ✅ Positive ROI (>5%)
- ✅ Sharpe ratio >1.0
- ✅ Max drawdown <20%
- ✅ Consistent performance (not just lucky streak)

---

### Phase 2 Deliverables Summary

**Code:**
- ✅ Web scraping service
- ✅ Daily automation scripts
- ✅ Results tracking system

**Data:**
- ✅ 12+ weeks of fixtures
- ✅ Daily odds snapshots
- ✅ All predictions logged
- ✅ All results recorded

**Validation:**
- ✅ Real track record with actual market odds
- ✅ Proven profitability over 3 months
- ✅ Statistical confidence in model performance
- ✅ Understanding of model strengths/weaknesses

**Documentation:**
- ✅ Weekly performance reports
- ✅ 12-week summary analysis
- ✅ Investor-ready track record document
- ✅ Model refinement notes

---

## Phase 3: Scale & Automate (Week 19+)

**Goal:** Replace scrapers with paid APIs, automate everything, prepare for production

**Cost:** $59-100/month (paid APIs)

**When to Start Phase 3:**
- ✅ Model validated profitable over 12+ weeks
- ✅ Consistent positive results
- ✅ Ready to scale beyond manual tracking
- ✅ Approaching investors OR using own capital

### Tasks

**Replace Scrapers with APIs:**
1. Subscribe to API-Football Pro ($60-100/month)
   - OR The Odds API ($59/month)
2. Update services to use paid APIs instead of scrapers
3. Remove web scraping code
4. Increase automation frequency (daily → hourly)

**Automate Everything:**
1. Automated daily fixture sync
2. Automated hourly odds updates
3. Automated prediction generation
4. Automated bet recommendations
5. Automated result tracking
6. Automated performance reporting

**Build Production Dashboard:**
1. Real-time predictions display
2. Live odds comparison
3. Performance metrics dashboard
4. Bankroll tracking
5. Bet history
6. ROI charts

**Enhance System:**
1. Add more leagues (if model works)
2. Add more markets (O/U, BTTS, etc.)
3. Add accumulator builder
4. Add smart market selector
5. Add notifications (Telegram/Email)

### Phase 3 Deliverables

**Infrastructure:**
- ✅ Production-ready APIs (no scrapers)
- ✅ Automated end-to-end pipeline
- ✅ Real-time data updates
- ✅ Monitoring and alerts

**Dashboard:**
- ✅ Investor-grade dashboard
- ✅ Live predictions
- ✅ Performance tracking
- ✅ Historical track record display

**Scale:**
- ✅ Ready for real capital deployment
- ✅ Multiple leagues/markets
- ✅ API rate limits managed
- ✅ Error handling & resilience

---

## Technical Architecture

### Phase 1 Architecture (Weeks 1-5)

```
┌─────────────────────────────────────────────────────┐
│                   OddsIQ MVP v1                      │
│                (Historical Data Only)                │
└─────────────────────────────────────────────────────┘

┌──────────────────┐
│  PostgreSQL DB   │
│  - 1,140 matches │
│  - 24 teams      │
│  - Features      │
│  - Predictions   │
└────────┬─────────┘
         │
         │
┌────────▼─────────┐         ┌──────────────────┐
│   Go Backend     │◄────────┤  Python ML       │
│   (Port 8000)    │         │  Service         │
│                  │         │  (Port 8001)     │
│  - API handlers  │         │                  │
│  - Repositories  │         │  - XGBoost       │
│  - ML client     │         │  - Features      │
└──────────────────┘         │  - Backtesting   │
                             │  - FastAPI       │
                             └──────────────────┘
```

### Phase 2 Architecture (Weeks 6-18)

```
┌─────────────────────────────────────────────────────┐
│               OddsIQ MVP v2                          │
│        (Real-world Testing)                          │
└─────────────────────────────────────────────────────┘

┌──────────────────┐
│  PostgreSQL DB   │
│  + Scraped data  │
│  + Live results  │
│  + Tracking      │
└────────┬─────────┘
         │
         │
┌────────▼─────────┐         ┌──────────────────┐
│   Go Backend     │◄────────┤  Python ML       │
│   (Port 8000)    │         │  Service         │
└────────┬─────────┘         │  (Port 8001)     │
         │                   └──────────────────┘
         │
         │
┌────────▼─────────┐         ┌──────────────────┐
│  Web Scraping    │         │  Manual Tracking │
│  Service         │         │  (Excel/CSV)     │
│  (Python)        │         │                  │
│                  │         │  - Weekly picks  │
│  - Fixtures      │         │  - Odds history  │
│  - Odds          │         │  - Results       │
│  - Results       │         │  - P&L           │
└──────────────────┘         └──────────────────┘
```

### Phase 3 Architecture (Week 19+)

```
┌─────────────────────────────────────────────────────┐
│            OddsIQ Production                         │
│           (Fully Automated)                          │
└─────────────────────────────────────────────────────┘

┌──────────────────┐         ┌──────────────────┐
│  API-Football    │         │  The Odds API    │
│  Pro Tier        │         │  (Optional)      │
│  - Live fixtures │         │  - Live odds     │
│  - Live odds     │         │  - More books    │
└────────┬─────────┘         └────────┬─────────┘
         │                            │
         └────────────┬───────────────┘
                      │
              ┌───────▼────────┐
              │  PostgreSQL DB │
              │  - Live data   │
              │  - Predictions │
              │  - Results     │
              └───────┬────────┘
                      │
         ┌────────────┴────────────┐
         │                         │
┌────────▼─────────┐      ┌───────▼──────────┐
│   Go Backend     │◄─────┤  Python ML       │
│   (Port 8000)    │      │  Service         │
│                  │      │  (Port 8001)     │
│  - Automated     │      │                  │
│    sync jobs     │      │  - Predictions   │
│  - Cron tasks    │      │  - Retraining    │
└────────┬─────────┘      └──────────────────┘
         │
         │
┌────────▼─────────┐
│  Next.js         │
│  Dashboard       │
│  (Port 3000)     │
│                  │
│  - Live picks    │
│  - Performance   │
│  - Track record  │
└──────────────────┘
```

---

## File Structure

### Phase 1 Files (Create Now)

```
OddsIQ/
├── backend/                      # Go service (existing)
│   ├── cmd/
│   │   ├── api/                  # API server
│   │   └── backfill/             # Historical data loader (✅ done)
│   ├── internal/
│   │   ├── models/               # Data models (✅ done)
│   │   ├── repository/           # Database layer (✅ done)
│   │   ├── services/
│   │   │   ├── ml_client.go      # ⏳ Create - Call Python ML service
│   │   │   └── prediction_tracker.go  # ⏳ Create - Track predictions
│   │   └── api/                  # HTTP handlers (✅ done)
│   └── pkg/
│       ├── apifootball/          # API-Football client (✅ done)
│       └── database/             # DB connection (✅ done)
│
├── ml-service/                   # ⏳ CREATE - Python ML service
│   ├── app/
│   │   ├── __init__.py
│   │   ├── main.py               # FastAPI app
│   │   ├── models/
│   │   │   ├── __init__.py
│   │   │   └── xgboost_model.py  # Model training & prediction
│   │   ├── features/
│   │   │   ├── __init__.py
│   │   │   ├── form_metrics.py   # Team form features
│   │   │   ├── h2h_stats.py      # Head-to-head features
│   │   │   ├── league_position.py # Position features
│   │   │   └── team_stats.py     # General team stats
│   │   ├── backtesting/
│   │   │   ├── __init__.py
│   │   │   ├── synthetic_odds.py # Generate realistic odds
│   │   │   └── backtest_engine.py # Backtest strategy
│   │   └── api/
│   │       ├── __init__.py
│   │       └── predictions.py    # FastAPI endpoints
│   ├── notebooks/
│   │   └── model_exploration.ipynb
│   ├── requirements.txt
│   ├── pyproject.toml
│   └── README.md
│
├── database/                     # Database (✅ done)
│   └── migrations/
│
├── docs/                         # Documentation
│   ├── MVP-ROADMAP-REVISED.md   # ✅ This file
│   ├── API-FREE-TIER-LIMITATIONS.md  # ✅ Done
│   └── PHASE-2-STATUS.md        # ✅ Done
│
└── README.md
```

### Phase 2 Files (Add Later)

```
OddsIQ/
├── scraping-service/            # ⏳ CREATE (Week 6)
│   ├── scrapers/
│   │   ├── __init__.py
│   │   ├── fixtures_scraper.py  # Scrape current fixtures
│   │   ├── odds_scraper.py      # Scrape odds
│   │   └── results_scraper.py   # Scrape results
│   ├── storage/
│   │   ├── __init__.py
│   │   └── scraper_db.py        # Save to PostgreSQL
│   ├── scheduler.py             # Daily cron job
│   ├── config.py
│   ├── requirements.txt
│   └── README.md
│
├── tracking/                    # ⏳ CREATE (Week 7)
│   ├── weekly_picks.csv         # Predictions each week
│   ├── odds_history.csv         # Daily odds snapshots
│   ├── results.csv              # Match outcomes & P&L
│   └── analysis.ipynb           # Performance analysis
│
└── docs/
    ├── scraping-guide.md        # Legal/technical scraping docs
    └── testing-reports/         # Weekly performance reports
        ├── week-07.md
        ├── week-08.md
        └── ...
```

### Phase 3 Files (Add Later)

```
OddsIQ/
├── frontend/                    # ⏳ CREATE (Week 19+)
│   ├── src/
│   │   ├── app/
│   │   ├── components/
│   │   └── lib/
│   └── package.json
│
└── Remove scraping-service/     # Replace with paid APIs
```

---

## Success Metrics

### Phase 1 Success Criteria

**Model Performance:**
- [ ] Accuracy >55% on test set
- [ ] Backtested ROI >5%
- [ ] Sharpe ratio >1.0
- [ ] Predictions well-calibrated

**Technical:**
- [ ] ML service running
- [ ] Go backend integrated
- [ ] Predictions stored in DB
- [ ] End-to-end flow works

### Phase 2 Success Criteria

**Real-World Performance (12 weeks):**
- [ ] 50+ predictions tracked
- [ ] Accuracy >55% on real fixtures
- [ ] ROI >5% with real odds
- [ ] Sharpe ratio >1.0
- [ ] Max drawdown <20%
- [ ] Consistent week-over-week

**Validation:**
- [ ] Model beats closing line (value)
- [ ] Model beats naive strategies
- [ ] Performance not due to luck (statistical significance)

### Phase 3 Success Criteria

**Production Ready:**
- [ ] Fully automated pipeline
- [ ] Dashboard complete
- [ ] Monitoring & alerts set up
- [ ] Ready for real capital

**Business Ready:**
- [ ] Track record documented
- [ ] Investor materials prepared
- [ ] Scaling plan defined

---

## Cost Breakdown

### Phase 1 (Weeks 1-5)
**Cost:** $0
- Use existing free API-Football data
- Development only

### Phase 2 (Weeks 6-18)
**Cost:** $0
- Web scraping (free, just time)
- Manual tracking (free)

### Phase 3 (Week 19+)
**API Costs (Monthly):**
- **Option A:** API-Football Pro - $60-100/month
- **Option B:** The Odds API - $59/month
- **Option C:** Both - $119-159/month

**First 3 months of Phase 3:** $177-300
**Annual (if continuing):** $708-1,200

### Total Investment Timeline

**Weeks 1-18:** $0
**Weeks 19-30 (3 months):** $177-300
**Total to validated product:** $177-300

**Compare to:**
- Paying APIs from Day 1: $708+ per year
- Building without validation: Risk of wasting $708+ on unproven system

---

## Risk Mitigation

### Phase 1 Risks

**Risk:** Model doesn't work on historical data
- **Mitigation:** Iterate on features, try different models
- **Fallback:** Re-evaluate approach, adjust strategy

**Risk:** Can't achieve >55% accuracy
- **Mitigation:** Lower bar to >52%, focus on value betting
- **Fallback:** Consider ensemble models, external data

### Phase 2 Risks

**Risk:** Web scrapers break frequently
- **Mitigation:** Multiple scraping targets, robust error handling
- **Fallback:** Manual data entry for testing period

**Risk:** Model doesn't work with real odds
- **Mitigation:** Continuous monitoring, quick iterations
- **Fallback:** Extend testing period, refine model

**Risk:** Legal issues with scraping
- **Mitigation:** Personal use only, respect ToS, add delays
- **Fallback:** Manual tracking, accelerate to Phase 3

### Phase 3 Risks

**Risk:** APIs too expensive
- **Mitigation:** Validate ROI first, ensure profitability covers cost
- **Fallback:** Continue scraping, seek investment

**Risk:** Model stops working at scale
- **Mitigation:** Continuous retraining, monitoring
- **Fallback:** Reduce bet sizes, iterate on model

---

## Decision Points

### After Phase 1 (Week 5)

**Questions to answer:**
1. Does the model achieve >55% accuracy on historical data?
2. Is the backtested ROI positive?
3. Are predictions well-calibrated?

**Decision:**
- ✅ Yes to all → Proceed to Phase 2
- ❌ No → Iterate on model, extend Phase 1

### After Phase 2 (Week 18)

**Questions to answer:**
1. Did we achieve >55% accuracy over 12 weeks with real odds?
2. Is the ROI positive and consistent?
3. Is performance statistically significant (not luck)?
4. Do we beat the closing line?

**Decision:**
- ✅ Yes to all → **Model validated!** → Proceed to Phase 3
- ⚠️ Mixed results → Extend testing, refine model
- ❌ No → **Pivot or stop** → Model doesn't work in real markets

### Before Phase 3 (Week 19)

**Questions to answer:**
1. Is the model profitable enough to cover API costs?
2. Am I ready to deploy real capital?
3. Do I want to approach investors now?

**Decision:**
- ✅ Ready → Subscribe to paid APIs, automate
- ⏸️ Not yet → Continue Phase 2 testing, build more track record
- 🔄 Alternative → Seek investment to fund APIs

---

## Investor Pitch Timeline

### After Phase 1 (Week 5)
**Not ready yet**
- Only have backtested results
- No real-world validation
- Theoretical performance only

### After Phase 2 (Week 18)
**Ready to pitch!**

**Investor materials:**
1. **Track Record:** 12 weeks of real predictions vs actual odds
2. **Performance Report:** ROI, Sharpe ratio, max drawdown
3. **Model Explanation:** How it works, what it predicts
4. **Business Plan:** Scale to more leagues, markets
5. **Capital Ask:** Amount needed for APIs, operations, growth

**Pitch:**
*"I built a sports betting AI that achieved X% ROI over 12 weeks of live testing with actual bookmaker odds. I manually tracked 50+ predictions to prove it works. Now I need capital to automate and scale."*

**Much stronger than:**
*"I built a model that looks good in backtests. I think it might work in real markets."*

---

## Next Immediate Steps

### This Session (Now)

1. ✅ Document roadmap (this file)
2. ✅ Mark Phase 2 (data infrastructure) complete
3. ⏳ Begin Phase 3 → Phase 1 of revised roadmap (ML model)

### This Week

1. Set up Python ML service structure
2. Create feature engineering pipeline
3. Extract features from 1,140 matches
4. Begin model training

### Next 2 Weeks

1. Complete model training
2. Build backtesting framework
3. Validate on historical data
4. Create ML API service

### Week 5

1. Integrate ML service with Go backend
2. Test end-to-end prediction flow
3. Review backtest results
4. **Decision point:** Proceed to real-world testing?

---

## References

**Key Documents:**
- `docs/API-FREE-TIER-LIMITATIONS.md` - API capabilities & constraints
- `docs/PHASE-2-STATUS.md` - Current infrastructure status
- `docs/implementation-plan.md` - Original 9-week plan (now revised)

**Technical Specs:**
- `docs/database-schema.md` - Database structure
- `docs/api-specification.md` - API endpoints
- `docs/architecture-decisions.md` - Technical decisions

---

## Revision History

| Date | Version | Changes |
|------|---------|---------|
| 2026-01-15 | 1.0 | Initial roadmap created based on API limitations discovery |

---

## Appendix A: Web Scraping Legal Considerations

### What's Generally Acceptable

✅ **Public data for personal research/education**
✅ **Respecting robots.txt**
✅ **Reasonable rate limiting (2-5 seconds between requests)**
✅ **Not bypassing paywalls or authentication**
✅ **For non-commercial personal testing**

### What to Avoid

❌ **Excessive requests (DDoS-like behavior)**
❌ **Bypassing technical restrictions**
❌ **Commercial use without permission**
❌ **Violating explicit ToS**
❌ **Claiming scraped data as your own**

### Best Practices for Our Use Case

1. **Respect Rate Limits:** 2-5 seconds between requests
2. **Use Multiple Sources:** Don't hammer one site
3. **Cache Aggressively:** Don't re-scrape same data
4. **User Agent:** Identify as personal research
5. **Robots.txt:** Check and respect
6. **Fallback:** Always have manual option
7. **Personal Use Only:** For testing your own system
8. **Transition:** Move to paid APIs when commercializing

### Risk Assessment: Low

- ✅ Personal research/testing phase
- ✅ Reasonable request volume
- ✅ Public odds display data
- ✅ Not competing with scraped sites
- ✅ Plan to transition to paid APIs

---

## Appendix B: Alternative Data Sources

### Free/Low-Cost Odds Sources

1. **Oddsportal.com**
   - Historical odds
   - Closing lines
   - Multiple bookmakers
   - Good for research

2. **Oddschecker.com**
   - Current odds comparison
   - Multiple bookmakers
   - UK-focused
   - Good scraping target

3. **Football-Data.co.uk**
   - Historical match results with odds
   - CSV downloads
   - Free for research
   - Limited to completed matches

### Free Fixture Sources

1. **BBC Sport**
   - Reliable fixture data
   - Results
   - Match stats
   - Good scraping target

2. **Premier League Official Site**
   - Official fixtures
   - Accurate results
   - Team stats
   - Good scraping target

3. **ESPN**
   - Fixtures
   - Results
   - Comprehensive stats
   - Alternative source

---

## Appendix C: Technology Stack

### Phase 1 Stack

**Backend (Go):**
- Gin (web framework)
- pgx (PostgreSQL driver)
- godotenv (environment config)

**ML Service (Python):**
- FastAPI (API framework)
- XGBoost (ML algorithm)
- scikit-learn (ML utilities)
- pandas (data manipulation)
- numpy (numerical computing)
- joblib (model persistence)

**Database:**
- PostgreSQL 15+

### Phase 2 Additional Stack

**Scraping (Python):**
- BeautifulSoup4 (HTML parsing)
- Selenium (dynamic pages)
- requests (HTTP client)
- lxml (XML/HTML parser)
- schedule (task scheduling)

### Phase 3 Additional Stack

**Frontend (Next.js):**
- Next.js 14+ (React framework)
- TypeScript
- Tailwind CSS (styling)
- Recharts (charts)
- shadcn/ui (components)

**Deployment:**
- Docker (containerization)
- AWS EC2 or similar (hosting)
- AWS RDS (managed PostgreSQL)
- PM2 or systemd (process management)

---

**END OF ROADMAP**

---

**This is your guiding document. Reference it at each decision point. Update it as you learn and adapt.** 🚀
