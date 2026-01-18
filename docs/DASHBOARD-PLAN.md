# OddsIQ Dashboard Plan

## Complete API Endpoint Reference

### Go Backend (Port 8000)

| Endpoint | Method | Description | Dashboard Feature |
|----------|--------|-------------|-------------------|
| `/health` | GET | Service health check | System Status |
| `/api/teams` | GET | List all teams | Manual Entry (dropdowns) |
| `/api/fixtures` | GET | List fixtures (with filters) | Fixtures List |
| `/api/fixtures/upcoming` | GET | Upcoming fixtures with odds status | Manual Entry, Predictions |
| `/api/fixtures/:id` | GET | Single fixture details | Fixture Detail |
| `/api/fixtures/:id/odds` | GET | Fixture odds | Fixture Detail |
| `/api/fixtures/manual` | POST | Create fixture manually | Manual Entry |
| `/api/fixtures/:id` | DELETE | Delete fixture | Manual Entry |
| `/api/odds/manual` | POST | Add single odds entry | Manual Entry |
| `/api/odds/manual/batch` | POST | Add multiple odds at once | Manual Entry |
| `/api/picks/weekly` | GET | Legacy 1X2 picks | - |
| `/api/picks/multi` | GET | Smart Market Selector picks | **Weekly Picks** |
| `/api/accumulators/weekly` | GET | Accumulator recommendations | **Accumulators** |
| `/api/accumulators/config` | GET | Accumulator configuration | Settings |
| `/api/predictions/fixture/:id` | GET | Single fixture prediction | Fixture Detail |
| `/api/predictions/fixture/:id/evaluate` | GET | Evaluate all markets | Fixture Analysis |
| `/api/model/metrics` | GET | 1X2 model metrics | Model Performance |
| `/api/model/metrics/all` | GET | All market model metrics | **Model Performance** |
| `/api/model/health` | GET | ML service health | System Status |
| `/api/bets` | GET | List placed bets | **Bet Tracker** |
| `/api/bets` | POST | Record a bet | Bet Tracker |
| `/api/bets/:id/settle` | PUT | Settle a bet | Bet Tracker |
| `/api/performance/summary` | GET | Performance metrics | **Dashboard Home** |
| `/api/performance/daily` | GET | Daily performance | Performance Charts |
| `/api/bankroll/history` | GET | Bankroll history | **Bankroll Chart** |

### Python ML Service (Port 8001)

| Endpoint | Method | Description | Used By |
|----------|--------|-------------|---------|
| `/api/predict` | POST | Single 1X2 prediction | Backend |
| `/api/predict/multi` | POST | Multi-market prediction | Backend |
| `/api/predict/batch` | POST | Batch predictions | Backend |
| `/api/model/metrics` | GET | Model metrics | Backend |
| `/api/model/metrics/all` | GET | All models metrics | Backend |
| `/api/markets` | GET | List available markets | Backend |
| `/api/model/reload` | POST | Reload model cache | Admin |

---

## Dashboard Pages

### 1. Dashboard Home (`/`)
**Purpose:** Overview of betting performance and key metrics

**Endpoints Used:**
- `GET /api/performance/summary` - ROI, win rate, profit
- `GET /api/bankroll/history` - Bankroll trend chart
- `GET /api/picks/multi?limit=5` - Top 5 current picks preview
- `GET /api/model/health` - System status indicator

**Components:**
```
┌─────────────────────────────────────────────────────────────────┐
│  OddsIQ Dashboard                              [System: Online] │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐           │
│  │ Bankroll │ │   ROI    │ │ Win Rate │ │  Profit  │           │
│  │  $1,250  │ │  +12.5%  │ │   58%    │ │  +$250   │           │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘           │
│                                                                 │
│  ┌─────────────────────────────────┐ ┌─────────────────────────┐│
│  │     Bankroll History Chart      │ │   Top 5 Picks Today    ││
│  │         [Line Chart]            │ │  1. Arsenal vs Chelsea ││
│  │                                 │ │     BTTS Yes @ 1.85    ││
│  │                                 │ │     EV: +8.2%          ││
│  │                                 │ │  2. Liverpool vs...    ││
│  └─────────────────────────────────┘ └─────────────────────────┘│
│                                                                 │
│  Quick Actions: [Add Fixture] [View All Picks] [Record Bet]    │
└─────────────────────────────────────────────────────────────────┘
```

---

### 2. Manual Entry (`/entry`)
**Purpose:** Enter upcoming fixtures and odds for predictions

**Endpoints Used:**
- `GET /api/teams` - Team dropdown options
- `GET /api/fixtures/upcoming` - List entered fixtures
- `POST /api/fixtures/manual` - Create new fixture
- `POST /api/odds/manual/batch` - Add odds to fixture
- `DELETE /api/fixtures/:id` - Remove fixture

**Components:**
```
┌─────────────────────────────────────────────────────────────────┐
│  Manual Entry                                                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────── Add New Fixture ────────────────────┐    │
│  │                                                         │    │
│  │  Home Team: [Arsenal        ▼]  Away Team: [Chelsea ▼] │    │
│  │  Match Date: [2025-01-25]  Time: [15:00]               │    │
│  │  Season: [2024]  Round: [Matchweek 22]                 │    │
│  │                                                         │    │
│  │  [Create Fixture]                                       │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  ┌─────────────────── Add Odds ───────────────────────────┐    │
│  │  Fixture: [Arsenal vs Chelsea - Jan 25 ▼]              │    │
│  │  Bookmaker: [Bet365          ]                         │    │
│  │                                                         │    │
│  │  1X2 Market:        Over/Under 2.5:      BTTS:         │    │
│  │  Home: [1.85]       Over:  [1.90]        Yes: [1.80]   │    │
│  │  Draw: [3.60]       Under: [1.95]        No:  [1.95]   │    │
│  │  Away: [4.20]                                          │    │
│  │                                                         │    │
│  │  [Add All Odds]                                        │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  ┌─────────────────── Upcoming Fixtures ──────────────────┐    │
│  │  Fixture              Date       Odds   Status  Action │    │
│  │  Arsenal vs Chelsea   Jan 25     7/7    Ready   [Del]  │    │
│  │  Liverpool vs Man U   Jan 26     7/7    Ready   [Del]  │    │
│  │  Everton vs Brighton  Jan 26     0/7    No Odds [Del]  │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

---

### 3. Weekly Picks (`/picks`)
**Purpose:** View recommended bets across all markets

**Endpoints Used:**
- `GET /api/picks/multi?bankroll=1000&limit=15` - Multi-market picks
- `GET /api/predictions/fixture/:id/evaluate` - Detailed fixture analysis

**Components:**
```
┌─────────────────────────────────────────────────────────────────┐
│  Weekly Picks                          Bankroll: [$1,000    ]   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Summary: 12 picks | Total Stake: $245 | Expected Value: +$32  │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ #1  Arsenal vs Chelsea              Sat, Jan 25 @ 15:00 │   │
│  │     ─────────────────────────────────────────────────── │   │
│  │     Best Pick: BTTS Yes @ 1.85                          │   │
│  │     Model Prob: 62%  |  EV: +14.7%  |  Confidence: High │   │
│  │     Suggested Stake: $35                                │   │
│  │                                                         │   │
│  │     All Value Bets:                                     │   │
│  │     • BTTS Yes      62%   1.85   +14.7%  $35           │   │
│  │     • Over 2.5      58%   1.90   +10.2%  $28           │   │
│  │     • Home Win      45%   1.85   -16.8%  --            │   │
│  │                                                         │   │
│  │     [View Details] [Record Bet]                         │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ #2  Liverpool vs Man Utd            Sun, Jan 26 @ 14:00 │   │
│  │     Best Pick: Home Win @ 1.65                          │   │
│  │     ...                                                 │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Filter: [All Markets ▼] [Min EV: 5% ▼] [Sort: EV ▼]          │
└─────────────────────────────────────────────────────────────────┘
```

---

### 4. Accumulators (`/accumulators`)
**Purpose:** View accumulator/parlay recommendations

**Endpoints Used:**
- `GET /api/accumulators/weekly?bankroll=1000` - Accumulator picks
- `GET /api/accumulators/config` - Configuration details

**Components:**
```
┌─────────────────────────────────────────────────────────────────┐
│  Accumulators                          Bankroll: [$1,000    ]   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Summary: 3 accumulators | Total Stake: $95 | Potential: $485  │
│                                                                 │
│  ┌────────────────────────── Double #1 ────────────────────┐   │
│  │  Combined Odds: 3.42  |  Win Prob: 35%  |  EV: +19.7%   │   │
│  │  Stake: $40  |  Potential Return: $137                  │   │
│  │                                                          │   │
│  │  Leg 1: Arsenal vs Chelsea                              │   │
│  │         BTTS Yes @ 1.85  (Prob: 62%)                    │   │
│  │                                                          │   │
│  │  Leg 2: Brighton vs Everton                             │   │
│  │         Over 2.5 @ 1.85  (Prob: 58%)                    │   │
│  │                                                          │   │
│  │  [Record Bet]                                           │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌────────────────────────── Treble #1 ────────────────────┐   │
│  │  Combined Odds: 5.28  |  Win Prob: 22%  |  EV: +16.2%   │   │
│  │  Stake: $30  |  Potential Return: $158                  │   │
│  │  ...                                                    │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Config: Min Legs: 2 | Max Legs: 3 | Min EV: 5%               │
└─────────────────────────────────────────────────────────────────┘
```

---

### 5. Bet Tracker (`/bets`)
**Purpose:** Record and track placed bets

**Endpoints Used:**
- `GET /api/bets` - List all bets
- `POST /api/bets` - Record new bet
- `PUT /api/bets/:id/settle` - Settle bet (won/lost)

**Components:**
```
┌─────────────────────────────────────────────────────────────────┐
│  Bet Tracker                                    [+ Record Bet]  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Filter: [All ▼] [Pending ▼] [This Week ▼]                     │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Date       Fixture            Bet         Odds  Stake  │   │
│  │  ────────────────────────────────────────────────────── │   │
│  │  Jan 20     Arsenal v Chelsea  BTTS Yes    1.85  $35    │   │
│  │             Status: PENDING    EV: +14.7%               │   │
│  │             [Settle: Won] [Settle: Lost]                │   │
│  │  ────────────────────────────────────────────────────── │   │
│  │  Jan 18     Liverpool v Spurs  Home Win    1.55  $50    │   │
│  │             Status: WON        Profit: +$27.50          │   │
│  │  ────────────────────────────────────────────────────── │   │
│  │  Jan 18     Man City v Wolves  Over 2.5    1.75  $40    │   │
│  │             Status: LOST       Profit: -$40.00          │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─── Record New Bet ──────────────────────────────────────┐   │
│  │  Fixture: [Arsenal vs Chelsea ▼]                        │   │
│  │  Bet Type: [BTTS Yes ▼]  Odds: [1.85]  Stake: [$35]    │   │
│  │  Bookmaker: [Bet365    ]                                │   │
│  │  [Record Bet]                                           │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

### 6. Performance (`/performance`)
**Purpose:** Detailed performance analytics

**Endpoints Used:**
- `GET /api/performance/summary` - Overall metrics
- `GET /api/performance/daily` - Daily breakdown
- `GET /api/bankroll/history` - Bankroll chart data
- `GET /api/model/metrics/all` - Model accuracy

**Components:**
```
┌─────────────────────────────────────────────────────────────────┐
│  Performance Analytics                                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐           │
│  │ Total    │ │ Win Rate │ │   ROI    │ │  Profit  │           │
│  │ 156 Bets │ │   58.3%  │ │  +12.5%  │ │  +$312   │           │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘           │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              Bankroll Over Time                         │   │
│  │  $1300 ┤                                    ●──●        │   │
│  │  $1200 ┤                          ●───●───●            │   │
│  │  $1100 ┤               ●────●────●                      │   │
│  │  $1000 ┼───●────●────●                                  │   │
│  │        └───────────────────────────────────────────────│   │
│  │         Week 1    Week 2    Week 3    Week 4           │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─── Performance by Market ───┐ ┌─── Model Accuracy ─────┐   │
│  │  Market      Bets  ROI      │ │  Model       Accuracy  │   │
│  │  1X2         62    +8.2%    │ │  1X2         56.1%     │   │
│  │  Over/Under  54    +15.3%   │ │  Over/Under  56.1%     │   │
│  │  BTTS        40    +18.7%   │ │  BTTS        52.6%     │   │
│  └─────────────────────────────┘ └─────────────────────────┘   │
│                                                                 │
│  ┌─── Daily Performance ───────────────────────────────────┐   │
│  │  Date       Bets   Won    Lost   Profit   ROI           │   │
│  │  Jan 20     8      5      3      +$45     +11.2%        │   │
│  │  Jan 19     6      4      2      +$32     +13.3%        │   │
│  │  Jan 18     10     5      5      -$12     -2.4%         │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

### 7. Fixture Detail (`/fixtures/:id`)
**Purpose:** Deep dive into a single fixture

**Endpoints Used:**
- `GET /api/fixtures/:id` - Fixture info
- `GET /api/fixtures/:id/odds` - All odds
- `GET /api/predictions/fixture/:id/evaluate` - Full market evaluation

**Components:**
```
┌─────────────────────────────────────────────────────────────────┐
│  Arsenal vs Chelsea                     Saturday, Jan 25 @ 3PM  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─── Team Stats ──────────────────────────────────────────┐   │
│  │        Arsenal              vs           Chelsea        │   │
│  │        Form: WWDWW                      Form: WLWDL     │   │
│  │        Pos: 2nd                         Pos: 6th        │   │
│  │        GF: 42  GA: 18                   GF: 35  GA: 28  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─── Market Analysis ─────────────────────────────────────┐   │
│  │  Market        Prediction    Odds    EV      Verdict    │   │
│  │  ─────────────────────────────────────────────────────  │   │
│  │  1X2           Home (45%)    1.85    -16.8%  NO VALUE   │   │
│  │                Draw (28%)    3.60    +0.8%   MARGINAL   │   │
│  │                Away (27%)    4.20    +13.4%  VALUE ✓    │   │
│  │  ─────────────────────────────────────────────────────  │   │
│  │  Over/Under    Over (58%)    1.90    +10.2%  VALUE ✓    │   │
│  │                Under (42%)   1.95    -18.1%  NO VALUE   │   │
│  │  ─────────────────────────────────────────────────────  │   │
│  │  BTTS          Yes (62%)     1.85    +14.7%  VALUE ✓    │   │
│  │                No (38%)      1.95    -25.9%  NO VALUE   │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Best Bet: BTTS Yes @ 1.85 (EV: +14.7%, Stake: $35)           │
│  [Record This Bet]                                             │
└─────────────────────────────────────────────────────────────────┘
```

---

## Dashboard Navigation

```
┌─────────────────────────────────────────────────────────────────┐
│  OddsIQ    [Dashboard] [Entry] [Picks] [Accas] [Bets] [Stats]  │
└─────────────────────────────────────────────────────────────────┘
```

| Nav Item | Route | Primary Purpose |
|----------|-------|-----------------|
| Dashboard | `/` | Overview & quick stats |
| Entry | `/entry` | Manual fixture & odds entry |
| Picks | `/picks` | Weekly betting recommendations |
| Accas | `/accumulators` | Accumulator recommendations |
| Bets | `/bets` | Record & track bets |
| Stats | `/performance` | Performance analytics |

---

## User Workflow

### Weekly Betting Workflow

```
1. ENTRY: Add upcoming fixtures
   └─► GET /api/teams (dropdown)
   └─► POST /api/fixtures/manual

2. ENTRY: Add odds from bookmaker
   └─► POST /api/odds/manual/batch

3. PICKS: View recommendations
   └─► GET /api/picks/multi?bankroll=1000

4. ACCAS: View accumulator suggestions
   └─► GET /api/accumulators/weekly?bankroll=1000

5. BETS: Record placed bets
   └─► POST /api/bets

6. BETS: Settle completed bets
   └─► PUT /api/bets/:id/settle

7. STATS: Review performance
   └─► GET /api/performance/summary
```

---

## Frontend Tech Stack (Recommended)

- **Framework**: Next.js 14 (App Router)
- **Styling**: Tailwind CSS
- **Charts**: Recharts
- **State**: React Query (for API caching)
- **Forms**: React Hook Form
- **UI Components**: shadcn/ui

---

## API Client Structure

```typescript
// lib/api.ts
const API_BASE = 'http://localhost:8000/api';

export const api = {
  // Teams
  getTeams: () => fetch(`${API_BASE}/teams`),

  // Fixtures
  getUpcomingFixtures: () => fetch(`${API_BASE}/fixtures/upcoming`),
  getFixture: (id: number) => fetch(`${API_BASE}/fixtures/${id}`),
  createFixture: (data: FixtureInput) => fetch(`${API_BASE}/fixtures/manual`, {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  deleteFixture: (id: number) => fetch(`${API_BASE}/fixtures/${id}`, {
    method: 'DELETE',
  }),

  // Odds
  addOddsBatch: (data: OddsBatchInput) => fetch(`${API_BASE}/odds/manual/batch`, {
    method: 'POST',
    body: JSON.stringify(data),
  }),

  // Picks
  getMultiMarketPicks: (bankroll: number, limit?: number) =>
    fetch(`${API_BASE}/picks/multi?bankroll=${bankroll}&limit=${limit || 15}`),

  // Accumulators
  getAccumulators: (bankroll: number) =>
    fetch(`${API_BASE}/accumulators/weekly?bankroll=${bankroll}`),

  // Predictions
  evaluateFixture: (id: number, bankroll: number) =>
    fetch(`${API_BASE}/predictions/fixture/${id}/evaluate?bankroll=${bankroll}`),

  // Bets
  getBets: () => fetch(`${API_BASE}/bets`),
  createBet: (data: BetInput) => fetch(`${API_BASE}/bets`, {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  settleBet: (id: number, result: 'won' | 'lost') =>
    fetch(`${API_BASE}/bets/${id}/settle`, {
      method: 'PUT',
      body: JSON.stringify({ result }),
    }),

  // Performance
  getPerformanceSummary: () => fetch(`${API_BASE}/performance/summary`),
  getDailyPerformance: () => fetch(`${API_BASE}/performance/daily`),
  getBankrollHistory: () => fetch(`${API_BASE}/bankroll/history`),

  // Model
  getModelMetrics: () => fetch(`${API_BASE}/model/metrics/all`),
  getMLHealth: () => fetch(`${API_BASE}/model/health`),
};
```

---

## Summary: Endpoints by Dashboard Feature

| Dashboard Page | Endpoints Used |
|----------------|----------------|
| **Home** | `/performance/summary`, `/bankroll/history`, `/picks/multi`, `/model/health` |
| **Manual Entry** | `/teams`, `/fixtures/upcoming`, `/fixtures/manual`, `/odds/manual/batch`, `/fixtures/:id` (DELETE) |
| **Weekly Picks** | `/picks/multi`, `/predictions/fixture/:id/evaluate` |
| **Accumulators** | `/accumulators/weekly`, `/accumulators/config` |
| **Bet Tracker** | `/bets` (GET, POST), `/bets/:id/settle` |
| **Performance** | `/performance/summary`, `/performance/daily`, `/bankroll/history`, `/model/metrics/all` |
| **Fixture Detail** | `/fixtures/:id`, `/fixtures/:id/odds`, `/predictions/fixture/:id/evaluate` |
