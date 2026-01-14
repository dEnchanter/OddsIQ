# OddsIQ Development Standards

**Established:** 2026-01-14
**Pattern:** Clean Architecture + Selective TDD

## Architecture Pattern

### Clean Architecture (Layered)

We organize code into distinct layers with clear dependencies:

```
┌─────────────────────────────────────────────────┐
│  Presentation Layer (API Handlers, Routes)     │
│  - HTTP endpoints                               │
│  - Request/response formatting                  │
│  - Input validation                             │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│  Service Layer (Business Logic)                │
│  - fixture_sync, odds_sync, betting_engine     │
│  - Kelly Criterion, EV calculations            │
│  - Market selection logic                       │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│  Repository Layer (Data Access)                │
│  - teams_repo, fixtures_repo, odds_repo        │
│  - CRUD operations                              │
│  - Query builders                               │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│  Infrastructure Layer (External Dependencies)  │
│  - API clients (apifootball, oddsapi)          │
│  - Database connection                          │
│  - Cron scheduler                               │
└─────────────────────────────────────────────────┘
```

### Dependency Rules

1. **Inner layers don't know about outer layers**
   - Repository doesn't import Service
   - Service doesn't import API handlers

2. **Dependencies point inward**
   - API → Service → Repository → Models

3. **Models are in the center**
   - Defined in `internal/models/`
   - Used by all layers
   - No external dependencies

## Testing Strategy: Selective TDD

### Critical Paths (MUST TEST)

**Financial Calculations:**
```go
// ALWAYS write tests FIRST for:
- Kelly Criterion stake sizing
- Expected Value (EV) calculations
- Profit/loss calculations
- Bankroll management
- ROI tracking
```

**ML Predictions:**
```go
// ALWAYS test:
- Feature engineering correctness
- Model prediction outputs
- Probability calibration
- Prediction API responses
```

**Betting Logic:**
```go
// ALWAYS test:
- Bet placement logic
- Bet settlement
- Value bet filtering
- Accumulator combinations
- Correlation detection
```

### Medium Priority (SHOULD TEST)

**Data Integrity:**
```go
// Test after implementation:
- Fixture matching algorithms
- Team name normalization
- Odds parsing
- Date/time handling
```

**Service Orchestration:**
```go
// Integration tests:
- Data sync workflows
- Scheduler jobs
- Multi-step operations
```

### Low Priority (TEST LATER OR SKIP)

**Infrastructure:**
```go
// Mock or skip:
- HTTP API clients (external dependencies)
- Database CRUD (use integration tests)
- Simple getters/setters
```

**Presentation:**
```go
// Can skip for MVP:
- HTTP handler input parsing
- JSON serialization
- UI components
```

## Directory Structure

```
backend/
├── cmd/
│   ├── api/              # Main API server
│   └── backfill/         # Data backfill tool
├── internal/
│   ├── models/           # Domain models (center)
│   ├── repository/       # Data access layer
│   ├── services/         # Business logic layer
│   └── api/              # HTTP handlers (outer)
├── pkg/
│   ├── apifootball/      # External API client
│   ├── oddsapi/          # External API client
│   └── database/         # Database utilities
├── config/               # Configuration
└── tests/                # Test files (new)
    ├── unit/             # Unit tests
    ├── integration/      # Integration tests
    └── fixtures/         # Test data fixtures

ml-service/
├── app/
│   ├── models/           # ML models
│   ├── features/         # Feature engineering
│   ├── api/              # FastAPI endpoints
│   └── services/         # Business logic
├── tests/                # Python tests
│   ├── unit/
│   ├── integration/
│   └── fixtures/
└── notebooks/            # Jupyter notebooks (exploration)

frontend/
├── src/
│   ├── app/              # Next.js pages
│   ├── components/       # React components
│   └── lib/              # Utilities
└── __tests__/            # Jest tests (if needed)
```

## Test File Conventions

### Go Tests

**File naming:**
```
service.go       → service_test.go
betting.go       → betting_test.go
kelly.go         → kelly_test.go
```

**Test naming:**
```go
func TestKellyCriterion_CalculateStake_PositiveEV(t *testing.T) {}
func TestKellyCriterion_CalculateStake_NegativeEV(t *testing.T) {}
func TestKellyCriterion_CalculateStake_ZeroEV(t *testing.T) {}
```

**Test structure:**
```go
func TestFunctionName_Scenario_ExpectedResult(t *testing.T) {
    // Arrange
    input := setupInput()
    expected := expectedOutput()

    // Act
    result := FunctionName(input)

    // Assert
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Python Tests

**File naming:**
```
feature_engineering.py  → test_feature_engineering.py
xgboost_model.py       → test_xgboost_model.py
```

**Test naming:**
```python
def test_calculate_form_metrics_last_5_games():
def test_calculate_form_metrics_less_than_5_games():
def test_calculate_form_metrics_no_games():
```

**Framework:** pytest

## TDD Workflow for Critical Code

### Example: Implementing Kelly Criterion

**Step 1: Write test FIRST**
```go
// backend/internal/services/kelly_test.go
func TestKellyCriterion_CalculateStake_PositiveEV(t *testing.T) {
    kelly := NewKellyCriterion(0.25) // 1/4 Kelly

    bankroll := 10000.0
    odds := 2.0
    probability := 0.55 // 55% win probability

    stake := kelly.CalculateStake(bankroll, odds, probability)

    // Expected: positive stake for +EV bet
    if stake <= 0 {
        t.Errorf("Expected positive stake for +EV bet, got %v", stake)
    }
    if stake > bankroll*0.25 {
        t.Errorf("Stake exceeds max Kelly fraction")
    }
}
```

**Step 2: Run test (should FAIL)**
```bash
go test ./internal/services/kelly_test.go
# FAIL: undefined: NewKellyCriterion
```

**Step 3: Write minimal code to pass**
```go
// backend/internal/services/kelly.go
type KellyCriterion struct {
    fraction float64
}

func NewKellyCriterion(fraction float64) *KellyCriterion {
    return &KellyCriterion{fraction: fraction}
}

func (k *KellyCriterion) CalculateStake(bankroll, odds, probability float64) float64 {
    // Kelly formula: f = (bp - q) / b
    // where b = odds - 1, p = probability, q = 1 - p

    b := odds - 1
    p := probability
    q := 1 - p

    kellyFraction := (b*p - q) / b

    // Apply fractional Kelly
    stake := bankroll * kellyFraction * k.fraction

    // Don't bet if negative EV
    if stake < 0 {
        return 0
    }

    return stake
}
```

**Step 4: Run test (should PASS)**
```bash
go test ./internal/services/kelly_test.go
# PASS
```

**Step 5: Add more tests, refactor, repeat**

## Code Review Checklist

### Before Committing

- [ ] Code follows clean architecture layers
- [ ] Critical paths have tests (if applicable)
- [ ] Tests pass: `go test ./...` or `pytest`
- [ ] No hardcoded values (use config)
- [ ] Error handling is present
- [ ] Logging for important operations
- [ ] Comments for complex logic only

### Pull Request Requirements

**For Critical Code (Financial, ML, Betting):**
- [ ] Tests written FIRST (TDD)
- [ ] 100% test coverage for critical functions
- [ ] Edge cases tested
- [ ] Integration tests if multi-layer

**For Non-Critical Code:**
- [ ] Tests optional (unless complex)
- [ ] Code follows architecture pattern
- [ ] No obvious bugs

## When to Write Tests

### Before Implementation (TDD)
✅ Kelly Criterion stake sizing
✅ Expected Value calculations
✅ Profit/loss tracking
✅ Bet settlement logic
✅ ML model predictions
✅ Feature engineering
✅ Accumulator combination logic
✅ Correlation detection

### After Implementation (Test After)
⚠️ Data sync services
⚠️ Fixture matching algorithms
⚠️ Odds parsing
⚠️ Repository queries (use integration tests)

### Skip for MVP (Test Later)
❌ HTTP API clients (mock external APIs)
❌ Simple CRUD operations
❌ Configuration loading
❌ Logging utilities
❌ UI components

## Testing Tools

### Go
- **Unit tests:** `testing` package (built-in)
- **Assertions:** `testify/assert` (optional)
- **Mocking:** `testify/mock` or interfaces
- **Coverage:** `go test -cover`

### Python
- **Unit tests:** `pytest`
- **Mocking:** `unittest.mock` or `pytest-mock`
- **ML testing:** `pytest` + numpy assertions
- **Coverage:** `pytest-cov`

### Integration Tests
- **Database:** Use test database or Docker
- **API:** Use `httptest` (Go) or `TestClient` (FastAPI)

## Example Test Organization

```go
// backend/tests/unit/kelly_test.go
package unit

import "testing"

func TestKellyCriterion_PositiveEV(t *testing.T) { /* ... */ }
func TestKellyCriterion_NegativeEV(t *testing.T) { /* ... */ }
func TestKellyCriterion_ZeroEV(t *testing.T) { /* ... */ }
```

```python
# ml-service/tests/unit/test_feature_engineering.py
import pytest
from app.features.form_metrics import calculate_form

def test_calculate_form_last_5_games():
    # Arrange
    fixtures = create_fixture_data()

    # Act
    form = calculate_form(fixtures, last_n=5)

    # Assert
    assert form['points'] == 12
    assert form['goals_for'] == 8
```

## Phase-Specific Testing Strategy

### Phase 1-2: Infrastructure (Done)
- ❌ No tests yet (clean architecture established)
- ✅ Code is testable (layered design)

### Phase 3: ML Model (Weeks 3-5)
- ✅ TDD for feature engineering
- ✅ Test model predictions
- ✅ Test probability calibration
- ❌ Skip API client tests

### Phase 4: Betting Engine (Weeks 5-6)
- ✅ TDD for ALL betting logic
- ✅ TDD for Kelly Criterion
- ✅ TDD for EV calculations
- ✅ Integration tests for bet workflow

### Phase 5: Dashboard (Weeks 5-6)
- ⚠️ Test after implementation
- ⚠️ Integration tests for API
- ❌ Skip UI unit tests (optional)

### Phase 6: Testing (Weeks 7-8)
- ✅ Full integration test suite
- ✅ End-to-end testing
- ✅ Paper trading validation

## Benefits of Our Approach

### Clean Architecture Benefits:
- ✅ Easy to test (isolated layers)
- ✅ Easy to understand (clear separation)
- ✅ Easy to change (swap implementations)
- ✅ Reusable code (service layer)

### Selective TDD Benefits:
- ✅ High confidence in critical paths
- ✅ Fast feature development
- ✅ Tests document critical behavior
- ✅ Catches financial bugs early

### Combined Benefits:
- ✅ Investor confidence (tested money logic)
- ✅ Fast MVP delivery
- ✅ Maintainable codebase
- ✅ Production-ready critical paths

## Anti-Patterns to Avoid

### ❌ Don't Do This:
- Writing tests for everything (too slow)
- Skipping tests for financial logic (too risky)
- Mixing layers (service calling handlers)
- Tight coupling to external APIs
- Hardcoding values in business logic
- Complex logic in handlers
- Database logic in services

### ✅ Do This Instead:
- Test critical paths only
- TDD for money calculations
- Keep layers separate
- Use repository pattern
- Use configuration files
- Keep handlers thin
- Use repositories for data access

## Summary

**Pattern:** Clean Architecture + Selective TDD

**Test Strategy:**
- 🔴 RED: Write test first for critical code
- 🟢 GREEN: Implement minimal code to pass
- 🔵 REFACTOR: Clean up code, tests stay green
- ⚪ SKIP: Non-critical code can skip tests for MVP

**Critical Paths:** Financial calculations, ML predictions, betting logic
**Non-Critical:** Infrastructure, CRUD, API clients, UI

**Goal:** Production-ready critical paths, fast MVP delivery

---

**This is our official development standard for OddsIQ.**
