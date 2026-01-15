"""
Build Training Dataset

Extracts features from all fixtures and saves to CSV
"""
import sys
import os
from datetime import datetime

# Add app to path
sys.path.insert(0, os.path.abspath(os.path.dirname(__file__)))

print("=" * 60)
print("Building Training Dataset")
print("=" * 60)
print(f"Started at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
print()

# Load fixtures
print("Step 1: Loading fixtures from database...")
try:
    from app.database.connection import get_fixtures_for_training

    seasons = [2022, 2023, 2024]
    fixtures = get_fixtures_for_training(seasons)

    print(f"[OK] Loaded {len(fixtures)} fixtures")
    print(f"   Seasons: {seasons}")
    print(f"   Date range: {fixtures[0]['match_date']} to {fixtures[-1]['match_date']}")
    print()

except Exception as e:
    print(f"[ERROR] Error loading fixtures: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Build feature dataset
print("Step 2: Extracting features (this may take 2-5 minutes)...")
try:
    from app.features.feature_builder import build_training_dataset

    df = build_training_dataset(fixtures, verbose=True)

    print()
    print(f"[OK] Feature extraction complete!")
    print(f"   Rows: {len(df)}")
    print(f"   Columns: {len(df.columns)}")
    print()

except Exception as e:
    print(f"[ERROR] Error extracting features: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Analyze dataset
print("Step 3: Analyzing dataset...")
try:
    from app.features.feature_builder import get_feature_columns, get_target_column

    feature_cols = get_feature_columns()
    target_col = get_target_column()

    # Check columns
    print(f"   Feature columns defined: {len(feature_cols)}")
    print(f"   Target column: {target_col}")

    # Check for missing features
    missing = [f for f in feature_cols if f not in df.columns]
    if missing:
        print(f"   [WARN]  Missing {len(missing)} features:")
        for f in missing[:5]:
            print(f"      - {f}")

    # Check for NaN values
    nan_counts = df[feature_cols].isnull().sum()
    total_nans = nan_counts.sum()

    print(f"   NaN values: {total_nans}")
    if total_nans > 0:
        print(f"   Columns with NaNs:")
        for col, count in nan_counts[nan_counts > 0].items():
            print(f"      - {col}: {count}")

    # Check target distribution
    if target_col in df.columns:
        outcome_counts = df[target_col].value_counts()
        print(f"   Target distribution:")
        print(f"      Home wins (0): {outcome_counts.get(0, 0)} ({outcome_counts.get(0, 0)/len(df)*100:.1f}%)")
        print(f"      Draws (1): {outcome_counts.get(1, 0)} ({outcome_counts.get(1, 0)/len(df)*100:.1f}%)")
        print(f"      Away wins (2): {outcome_counts.get(2, 0)} ({outcome_counts.get(2, 0)/len(df)*100:.1f}%)")

    print()

except Exception as e:
    print(f"[WARN]  Warning during analysis: {e}")

# Save to CSV
print("Step 4: Saving to CSV...")
try:
    # Create data directory if it doesn't exist
    os.makedirs('data', exist_ok=True)

    output_file = 'data/training_data.csv'
    df.to_csv(output_file, index=False)

    # Get file size
    file_size_mb = os.path.getsize(output_file) / (1024 * 1024)

    print(f"[OK] Dataset saved to: {output_file}")
    print(f"   File size: {file_size_mb:.2f} MB")
    print()

except Exception as e:
    print(f"[ERROR] Error saving dataset: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Summary
print("=" * 60)
print("[OK] DATASET BUILD COMPLETE!")
print("=" * 60)
print()
print("Dataset Summary:")
print(f"  - Samples: {len(df)}")
print(f"  - Features: {len(feature_cols)}")
print(f"  - File: {output_file}")
print(f"  - Size: {file_size_mb:.2f} MB")
print()
print("Next steps:")
print("  1. Review dataset: df = pd.read_csv('data/training_data.csv')")
print("  2. Train model: python train_model.py")
print("  3. Or open notebook for exploration")
print()
print(f"Completed at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
