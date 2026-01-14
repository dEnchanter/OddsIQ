# Multi-Market Strategy - Update Summary

**Date:** 2026-01-12
**Change:** Upgraded from single-market (1X2 only) to multi-market approach

## What Changed

### Strategy Shift

**Before:**
- Focus only on 1X2 (Home/Draw/Away) predictions
- ~10 picks per week
- Single XGBoost model
- Limited to match outcome market

**After:**
- Multi-market intelligent selection
- **3 models** covering major markets:
  1. **1X2** (Home/Draw/Away) - Week 3
  2. **Over/Under 2.5 Goals** - Week 4
  3. **BTTS** (Both Teams to Score) - Week 5
- **Smart Market Selector** - Week 6
- ~20-30 picks per week
- Picks highest EV across ALL markets

## Why This Matters

### Real Example
```
Fixture: Arsenal vs Liverpool

OLD APPROACH (1X2 only):
✓ Analyze: Home/Draw/Away
✓ Best option: Draw @ 3.50, EV = +1.8%
✗ Below 3% threshold - NO BET

NEW APPROACH (Multi-market):
✓ Analyze: 1X2, Over/Under, BTTS
✓ Calculate EV for all options:
  - Draw (1X2): +1.8% ❌
  - Over 2.5: +8.3% ✅ RECOMMENDED
  - BTTS Yes: +6.1% ✅ ALTERNATIVE
✓ Multiple high-EV bets found!
```

**Result:** Old approach would have missed 2 profitable opportunities.

## Updated Documents

### Core Documentation
1. **`docs/market-expansion-roadmap.md`** ⭐ NEW
   - Complete market rollout plan
   - 12+ future markets documented
   - Technical implementation details
   - Success metrics per market

2. **`docs/implementation-plan.md`** ✏️ UPDATED
   - Phase 3: Multi-market model development (Weeks 3-5)
   - Phase 4: Smart Market Selector (Weeks 6-7)
   - Updated success criteria

3. **`docs/architecture-decisions.md`** ✏️ UPDATED
   - ADR-010: Multi-Market Strategy (NEW)
   - Updated ADR-009 with multi-market targets

4. **`CLAUDE.md`** ✏️ UPDATED
   - Multi-market architecture
   - 3-model approach documented
   - Smart Market Selector workflow
   - Future expansion roadmap referenced

5. **`README.md`** ✏️ UPDATED
   - MVP roadmap reflects 3 models
   - Success criteria updated
   - 20-30 picks per week target

## Technical Changes Needed

### Database (No Changes Required! ✅)
The existing schema already supports multi-market:
```sql
-- odds table already has market_type field
CREATE TABLE odds (
    ...
    market_type VARCHAR(50) NOT NULL,  -- 'h2h', 'totals_2_5', 'btts'
    ...
);
```

### ML Service (To Be Implemented)
```
ml-service/app/models/
├── base_model.py        # Abstract base class
├── h2h_model.py         # 1X2 model (Week 3)
├── totals_model.py      # O/U 2.5 model (Week 4)
├── btts_model.py        # BTTS model (Week 5)
└── model_manager.py     # Orchestrates all models
```

### Go Backend (To Be Implemented)
```
backend/internal/services/
└── market_selector.go   # Smart Market Selector (Week 6)
```

### API Endpoints (To Be Added)
```
POST /api/predict/all-markets        # ML Service
GET /api/picks/weekly?market=totals  # Backend
GET /api/performance/by-market       # Backend
```

## Implementation Timeline

| Week | Milestone | Deliverable |
|------|-----------|-------------|
| 3 | 1X2 Model | 55-60% accuracy, proves pipeline works |
| 4 | O/U 2.5 Model | 60-65% accuracy, totals market covered |
| 5 | BTTS Model | 58-62% accuracy, 3rd market operational |
| 6 | Market Selector | Smart selector picks best EV across all markets |
| 7 | Integration Testing | End-to-end multi-market workflow validated |
| 8 | Paper Trading | 2 weekends with 20-30 picks/week |

## Success Metrics

### MVP Targets (Week 8)
- **3 models** trained and operational
- **Combined accuracy:** 58-63% (weighted average)
- **ROI target:** 6-10% (vs 3-5% single market)
- **Weekly picks:** 20-30 (vs 10)
- **Market coverage:** 85% of available betting volume

### Performance Tracking
Track separately per market:
- ROI percentage
- Win rate
- Average odds
- Closing Line Value (CLV)
- Stake distribution

**Overall performance = weighted average across all markets**

## Future Markets (Post-MVP)

See `docs/market-expansion-roadmap.md` for full details.

**Month 3-4:**
- Double Chance
- Over/Under 1.5, 3.5, 4.5
- Half-time markets

**Month 5-6:**
- Handicaps
- Correct Score
- First/Last Goal

**Month 7+:**
- Player props
- Corners/Cards
- Exotic markets

**Target:** 15-20 markets by Month 6

## Risk Management

### Diversification Rules
- **Max 40%** of weekly picks from any single market
- **Independent models** - failure of one doesn't affect others
- **Market-specific Kelly fractions** - conservative on new markets
- **CLV tracking** per market - disable underperforming markets

### Validation Requirements
Before enabling each market:
1. Backtest on 3 seasons of historical data
2. Achieve target accuracy on holdout set
3. Positive theoretical ROI > 5%
4. Stable probability calibration
5. At least 100 test samples

## Implementation Priority

**Must Have (Week 8):**
- ✅ 1X2 model
- ✅ Over/Under 2.5 model
- ✅ BTTS model
- ✅ Smart Market Selector

**Nice to Have (Month 3):**
- Double Chance (can derive from 1X2)
- Over/Under 1.5, 3.5

**Future (Month 4+):**
- Everything else in expansion roadmap

## Questions & Answers

### Q: Will this delay MVP?
**A:** Timeline extends from 6 weeks to 8 weeks (2 weeks longer), but delivers significantly better product.

### Q: What if one model performs poorly?
**A:** Can disable specific markets. System works fine with just 1-2 markets active.

### Q: How does this affect API costs?
**A:** The Odds API already provides multi-market data in single request. No significant cost increase.

### Q: Is the database ready?
**A:** Yes! Schema already supports `market_type` field. No migrations needed.

### Q: Can we add more markets later?
**A:** Absolutely! See `docs/market-expansion-roadmap.md` - designed for incremental expansion.

## Next Steps

1. **Continue Phase 2** (Data Infrastructure) as planned
   - Fetch odds for multiple markets when syncing
   - Store with correct `market_type` values

2. **Week 3**: Build 1X2 model (as originally planned)

3. **Week 4**: Add Over/Under 2.5 model

4. **Week 5**: Add BTTS model

5. **Week 6**: Implement Smart Market Selector

6. **Week 7-8**: Integration testing and paper trading

---

**Status:** ✅ All documentation updated
**Impact:** Higher ROI potential, more picks, better diversification
**Risk:** Low (incremental rollout, can rollback)
**Timeline:** +2 weeks to MVP (8 weeks total instead of 6)
**Investor Appeal:** High (demonstrates sophistication and comprehensive coverage)
