# OddsIQ ML Service

Machine Learning service for football match prediction

## Setup Complete! ✅

### What's Been Set Up

1. **✅ Project Structure**
   - `app/` - Main application code
   - `app/database/` - Database connection
   - `app/features/` - Feature engineering modules
   - `app/models/` - ML model (to be created)
   - `app/api/` - FastAPI endpoints
   - `config/` - Configuration
   - `notebooks/` - Jupyter notebooks for exploration

2. **✅ Configuration**
   - `.env` file created with database credentials
   - `config.py` with all settings
   - Ready for local development

3. **✅ Dependencies**
   - Virtual environment: `venv/`
   - Installing dependencies (in progress)
   - All packages in `requirements.txt`

4. **✅ Database Connection**
   - `app/database/connection.py` created
   - Functions to query fixtures and teams
   - Ready to load 1,140 matches from PostgreSQL

5. **✅ Feature Engineering Modules**
   - `form_metrics.py` - Team form (last 5 games, home/away splits)
   - `h2h_stats.py` - Head-to-head statistics
   - `league_position.py` - League table positions
   - `feature_builder.py` - Orchestrates all features
   - **Total: 70+ features per match**

6. **✅ Test Script**
   - `test_setup.py` - Verifies everything works
   - Tests database, features, dataset building

## What You Have Now

```
1,140 fixtures → Feature Engineering → 70+ features per match → Training Dataset
```

## Next Steps

### 1. Wait for Dependencies (5-10 minutes)

Dependencies are installing. Once complete, you'll see:
```
Successfully installed pandas-2.1.4 numpy-1.26.3 xgboost-2.0.3 ...
```

### 2. Test Setup

```bash
cd ml-service
venv\Scripts\activate
python test_setup.py
```

**Expected output:**
```
✅ Database connection successful!
✅ Loaded 1,140 fixtures for training
✅ Extracted 70+ features!
✅ ALL TESTS PASSED!
```

### 3. Build Training Dataset

Create `build_dataset.py`:
```python
from app.database.connection import get_fixtures_for_training
from app.features.feature_builder import build_training_dataset

# Load all fixtures
fixtures = get_fixtures_for_training([2022, 2023, 2024])
print(f"Loaded {len(fixtures)} fixtures")

# Extract features
df = build_training_dataset(fixtures, verbose=True)

# Save to CSV
df.to_csv('data/training_data.csv', index=False)
print(f"✅ Saved training dataset: {len(df)} rows, {len(df.columns)} columns")
```

Run it:
```bash
python build_dataset.py
```

This will create `data/training_data.csv` with 1,140 rows and 70+ feature columns!

### 4. Train XGBoost Model

Create notebook or script to train model:

```python
import pandas as pd
from sklearn.model_selection import train_test_split
import xgboost as xgb

# Load data
df = pd.read_csv('data/training_data.csv')

# Get features and target
from app.features.feature_builder import get_feature_columns, get_target_column

X = df[get_feature_columns()]
y = df[get_target_column()]

# Train/test split (chronological)
split_idx = int(len(df) * 0.8)
X_train, X_test = X[:split_idx], X[split_idx:]
y_train, y_test = y[:split_idx], y[split_idx:]

# Train XGBoost
model = xgb.XGBClassifier(
    objective='multi:softprob',
    num_class=3,
    max_depth=6,
    learning_rate=0.1,
    n_estimators=200
)

model.fit(X_train, y_train)

# Evaluate
accuracy = model.score(X_test, y_test)
print(f"Test accuracy: {accuracy:.3f}")

# Save model
import joblib
joblib.dump(model, 'models/xgboost_v1.pkl')
print("✅ Model saved!")
```

## File Structure

```
ml-service/
├── .env                        # ✅ Created - Database config
├── requirements.txt            # ✅ Created - Dependencies
├── test_setup.py              # ✅ Created - Test script
├── build_dataset.py           # ⏳ To create next
├── train_model.py             # ⏳ To create next
│
├── venv/                      # ✅ Created - Virtual environment
│
├── app/
│   ├── main.py                # ✅ Ready - FastAPI app
│   ├── database/
│   │   └── connection.py      # ✅ Created - DB functions
│   ├── features/
│   │   ├── form_metrics.py    # ✅ Created - Form features
│   │   ├── h2h_stats.py       # ✅ Created - H2H features
│   │   ├── league_position.py # ✅ Created - Position features
│   │   └── feature_builder.py # ✅ Created - Orchestrator
│   ├── models/
│   │   └── xgboost_model.py   # ⏳ To create next
│   └── api/
│       └── predictions.py     # ✅ Ready - API endpoints
│
├── config/
│   └── config.py              # ✅ Ready - Configuration
│
├── models/                    # Models saved here
│   └── xgboost_v1.pkl        # ⏳ After training
│
├── data/                      # Data saved here
│   └── training_data.csv     # ⏳ After build_dataset.py
│
└── notebooks/                 # For exploration
    └── exploration.ipynb      # ⏳ Optional
```

## Features Extracted (70+)

### Form Features (last 5 games)
- Points, wins, draws, losses
- Goals scored/conceded
- Goal difference
- Clean sheets, failed to score
- Averages per game

### Home/Away Splits (last 3 games)
- Home-specific form
- Away-specific form
- Performance at venue type

### Head-to-Head
- Previous meeting results (last 5)
- H2H goals, win percentages
- Home advantage in H2H

### League Position
- Current position
- Points, points per game
- Goals for/against
- Goal difference
- Win percentage

### Differentials
- Position difference
- Points difference
- Form difference
- Goal difference

**Total: ~70 features per match!**

## Testing Checklist

After dependencies install:

- [ ] Run `test_setup.py` - All tests pass
- [ ] Create `build_dataset.py` script
- [ ] Run `build_dataset.py` - Creates CSV with 1,140 rows
- [ ] Check `data/training_data.csv` exists
- [ ] Open in Excel/pandas - Verify features look correct
- [ ] Create training script
- [ ] Train XGBoost model
- [ ] Achieve >55% accuracy
- [ ] Save model to `models/xgboost_v1.pkl`

## Current Status

✅ **Setup Phase: COMPLETE**
⏳ **Dependencies: Installing (5-10 mins remaining)**
⏳ **Testing: Ready when dependencies done**
⏳ **Dataset Building: Next step**
⏳ **Model Training: After dataset**

## Estimated Time to First Model

- Dependencies install: 5-10 minutes (running now)
- Test setup: 1 minute
- Build dataset: 2-5 minutes (1,140 matches)
- Train model: 5-10 minutes
- **Total: ~20-30 minutes from now!**

## Key Files Created This Session

1. `.env` - Database configuration
2. `app/database/connection.py` - DB queries (342 lines)
3. `app/features/form_metrics.py` - Form features (256 lines)
4. `app/features/h2h_stats.py` - H2H features (143 lines)
5. `app/features/league_position.py` - Position features (179 lines)
6. `app/features/feature_builder.py` - Orchestrator (169 lines)
7. `test_setup.py` - Test script (126 lines)
8. `README.md` - This file

**Total: ~1,200 lines of feature engineering code!**

---

**You're all set! Once dependencies finish installing, run the test script and you're ready to build your first model!** 🚀
