"""
Test ML Service Setup

Verifies database connection and feature extraction
"""
import sys
import os

# Add app to path
sys.path.insert(0, os.path.abspath(os.path.dirname(__file__)))

print("=" * 60)
print("Testing ML Service Setup")
print("=" * 60)

# Test 1: Database connection
print("\n1. Testing database connection...")
try:
    from app.database.connection import test_connection
    if test_connection():
        print("   [OK] Database connection successful!")
    else:
        print("   [ERROR] Database connection failed!")
        sys.exit(1)
except Exception as e:
    print(f"   [ERROR] Error: {e}")
    sys.exit(1)

# Test 2: Load fixtures
print("\n2. Loading fixtures from database...")
try:
    from app.database.connection import get_fixtures_for_training
    fixtures = get_fixtures_for_training([2022, 2023, 2024])
    print(f"   [OK] Loaded {len(fixtures)} fixtures for training")

    if len(fixtures) == 0:
        print("   [WARN]  Warning: No fixtures found!")
        sys.exit(1)

    # Show sample
    print(f"   Sample fixture: {fixtures[0]['home_team_name']} vs {fixtures[0]['away_team_name']}")
    print(f"     Date: {fixtures[0]['match_date']}")
    print(f"     Score: {fixtures[0]['home_score']}-{fixtures[0]['away_score']}")

except Exception as e:
    print(f"   [ERROR] Error: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Test 3: Feature extraction
print("\n3. Testing feature extraction on sample fixture...")
try:
    from app.features.feature_builder import extract_features_for_fixture

    # Get a fixture from later in the season (so it has history)
    test_fixture = fixtures[len(fixtures) // 2]  # Middle of dataset

    print(f"   Testing on: {test_fixture['home_team_name']} vs {test_fixture['away_team_name']}")

    features = extract_features_for_fixture(test_fixture, fixtures)

    print(f"   [OK] Extracted {len(features)} features!")
    print(f"   Sample features:")
    print(f"     Home form (last 5): {features.get('home_form_last_5_points', 0)} points")
    print(f"     Away form (last 5): {features.get('away_form_last_5_points', 0)} points")
    print(f"     H2H games played: {features.get('h2h_games_played', 0)}")
    print(f"     Position diff: {features.get('position_diff', 0)}")
    print(f"     Outcome: {features.get('outcome', 'unknown')}")

except Exception as e:
    print(f"   [ERROR] Error: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Test 4: Build small dataset
print("\n4. Building feature dataset for first 10 fixtures...")
try:
    from app.features.feature_builder import build_training_dataset

    small_sample = fixtures[100:110]  # 10 fixtures from middle of dataset
    df = build_training_dataset(small_sample, verbose=False)

    print(f"   [OK] Built dataset with {len(df)} rows and {len(df.columns)} columns")
    print(f"   Columns: {', '.join(df.columns[:10])}...")

    # Check for NaN values
    nan_count = df.isnull().sum().sum()
    print(f"   NaN values: {nan_count}")

    if nan_count > 0:
        print(f"   [WARN]  Warning: Dataset has {nan_count} NaN values")

except Exception as e:
    print(f"   [ERROR] Error: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Test 5: Check feature count
print("\n5. Verifying feature configuration...")
try:
    from app.features.feature_builder import get_feature_columns, get_target_column

    feature_cols = get_feature_columns()
    target_col = get_target_column()

    print(f"   [OK] {len(feature_cols)} features defined")
    print(f"   [OK] Target column: {target_col}")

    # Check if features exist in our dataframe
    missing_features = [f for f in feature_cols if f not in df.columns]
    if missing_features:
        print(f"   [WARN]  Warning: {len(missing_features)} features missing from dataset:")
        for f in missing_features[:5]:
            print(f"      - {f}")
    else:
        print(f"   [OK] All features present in dataset!")

except Exception as e:
    print(f"   [ERROR] Error: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

print("\n" + "=" * 60)
print("[OK] ALL TESTS PASSED!")
print("=" * 60)
print("\nYou're ready to build the training dataset and train the model!")
print("\nNext steps:")
print("  1. Run: python build_dataset.py")
print("  2. Check: data/training_data.csv")
print("  3. Train model in notebook or script")
