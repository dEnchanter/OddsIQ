"""
Train XGBoost Model for Match Outcome Prediction

Trains a multi-class classifier to predict Home Win / Draw / Away Win
"""
import sys
import os
import pickle
from datetime import datetime

import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split, TimeSeriesSplit
from sklearn.metrics import accuracy_score, classification_report, confusion_matrix
from sklearn.preprocessing import StandardScaler
import xgboost as xgb

# Add app to path
sys.path.insert(0, os.path.abspath(os.path.dirname(__file__)))

from config.config import config


def load_data(filepath: str) -> pd.DataFrame:
    """Load training data from CSV"""
    print(f"Loading data from {filepath}...")
    df = pd.read_csv(filepath)
    print(f"  Loaded {len(df)} samples with {len(df.columns)} columns")
    return df


def prepare_features(df: pd.DataFrame) -> tuple:
    """
    Prepare features and target for training

    Returns:
        X: Feature matrix
        y: Target vector
        feature_names: List of feature column names
    """
    # Define feature columns (exclude metadata and target)
    exclude_cols = [
        'fixture_id', 'season', 'match_date', 'home_team_id', 'away_team_id',
        'outcome', 'outcome_encoded', 'home_score', 'away_score', 'total_goals'
    ]

    feature_cols = [col for col in df.columns if col not in exclude_cols]

    print(f"  Using {len(feature_cols)} features")

    X = df[feature_cols].copy()
    y = df['outcome_encoded'].copy()

    # Handle any remaining NaN values
    nan_count = X.isnull().sum().sum()
    if nan_count > 0:
        print(f"  Filling {nan_count} NaN values with 0")
        X = X.fillna(0)

    return X, y, feature_cols


def train_test_split_chronological(df: pd.DataFrame, test_size: float = 0.2):
    """
    Split data chronologically (important for time series)
    Later matches go to test set
    """
    # Sort by date
    df_sorted = df.sort_values('match_date').reset_index(drop=True)

    split_idx = int(len(df_sorted) * (1 - test_size))

    train_df = df_sorted.iloc[:split_idx]
    test_df = df_sorted.iloc[split_idx:]

    print(f"  Train: {len(train_df)} samples (up to {train_df['match_date'].max()})")
    print(f"  Test: {len(test_df)} samples (from {test_df['match_date'].min()})")

    return train_df, test_df


def train_xgboost(X_train, y_train, X_test, y_test):
    """
    Train XGBoost classifier
    """
    print("\nTraining XGBoost model...")

    # XGBoost parameters
    params = {
        'objective': 'multi:softprob',
        'num_class': 3,
        'max_depth': config.XGBOOST_MAX_DEPTH,
        'learning_rate': config.XGBOOST_LEARNING_RATE,
        'n_estimators': config.XGBOOST_N_ESTIMATORS,
        'subsample': 0.8,
        'colsample_bytree': 0.8,
        'random_state': config.RANDOM_STATE,
        'eval_metric': 'mlogloss',
        'early_stopping_rounds': 20,
    }

    print(f"  Parameters: max_depth={params['max_depth']}, "
          f"learning_rate={params['learning_rate']}, "
          f"n_estimators={params['n_estimators']}")

    model = xgb.XGBClassifier(**params)

    # Train with early stopping
    model.fit(
        X_train, y_train,
        eval_set=[(X_test, y_test)],
        verbose=False
    )

    print(f"  Best iteration: {model.best_iteration}")

    return model


def evaluate_model(model, X_test, y_test, feature_names):
    """
    Evaluate model performance
    """
    print("\n" + "=" * 60)
    print("MODEL EVALUATION")
    print("=" * 60)

    # Predictions
    y_pred = model.predict(X_test)
    y_pred_proba = model.predict_proba(X_test)

    # Accuracy
    accuracy = accuracy_score(y_test, y_pred)
    print(f"\nOverall Accuracy: {accuracy:.1%}")

    # Baseline comparison (always predict most common class)
    most_common = y_test.mode()[0]
    baseline_accuracy = (y_test == most_common).mean()
    print(f"Baseline Accuracy (always predict class {most_common}): {baseline_accuracy:.1%}")
    print(f"Improvement over baseline: {(accuracy - baseline_accuracy):.1%}")

    # Classification report
    print("\nClassification Report:")
    print("-" * 40)
    target_names = ['Home Win', 'Draw', 'Away Win']
    print(classification_report(y_test, y_pred, target_names=target_names))

    # Confusion matrix
    print("Confusion Matrix:")
    print("-" * 40)
    cm = confusion_matrix(y_test, y_pred)
    print(f"              Predicted")
    print(f"              Home  Draw  Away")
    print(f"Actual Home   {cm[0][0]:4d}  {cm[0][1]:4d}  {cm[0][2]:4d}")
    print(f"       Draw   {cm[1][0]:4d}  {cm[1][1]:4d}  {cm[1][2]:4d}")
    print(f"       Away   {cm[2][0]:4d}  {cm[2][1]:4d}  {cm[2][2]:4d}")

    # Feature importance
    print("\nTop 15 Most Important Features:")
    print("-" * 40)
    importance = pd.DataFrame({
        'feature': feature_names,
        'importance': model.feature_importances_
    }).sort_values('importance', ascending=False)

    for i, row in importance.head(15).iterrows():
        print(f"  {row['feature']}: {row['importance']:.4f}")

    # Class distribution in predictions
    print("\nPrediction Distribution:")
    print("-" * 40)
    pred_dist = pd.Series(y_pred).value_counts().sort_index()
    actual_dist = y_test.value_counts().sort_index()
    print(f"  Home Wins - Predicted: {pred_dist.get(0, 0):3d}, Actual: {actual_dist.get(0, 0):3d}")
    print(f"  Draws     - Predicted: {pred_dist.get(1, 0):3d}, Actual: {actual_dist.get(1, 0):3d}")
    print(f"  Away Wins - Predicted: {pred_dist.get(2, 0):3d}, Actual: {actual_dist.get(2, 0):3d}")

    return {
        'accuracy': accuracy,
        'baseline_accuracy': baseline_accuracy,
        'classification_report': classification_report(y_test, y_pred, target_names=target_names, output_dict=True),
        'confusion_matrix': cm.tolist(),
        'feature_importance': importance.to_dict('records')
    }


def save_model(model, feature_names, metrics, filepath: str):
    """
    Save trained model and metadata
    """
    print(f"\nSaving model to {filepath}...")

    # Create models directory if it doesn't exist
    os.makedirs(os.path.dirname(filepath), exist_ok=True)

    model_data = {
        'model': model,
        'feature_names': feature_names,
        'metrics': metrics,
        'training_date': datetime.now().isoformat(),
        'model_version': config.MODEL_VERSION,
        'config': {
            'max_depth': config.XGBOOST_MAX_DEPTH,
            'learning_rate': config.XGBOOST_LEARNING_RATE,
            'n_estimators': config.XGBOOST_N_ESTIMATORS,
        }
    }

    with open(filepath, 'wb') as f:
        pickle.dump(model_data, f)

    file_size_mb = os.path.getsize(filepath) / (1024 * 1024)
    print(f"  Model saved ({file_size_mb:.2f} MB)")


def main():
    print("=" * 60)
    print("XGBoost Model Training")
    print("=" * 60)
    print(f"Started at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print()

    # Step 1: Load data
    print("Step 1: Loading training data...")
    df = load_data('data/training_data.csv')
    print()

    # Step 2: Split data chronologically
    print("Step 2: Splitting data (chronological)...")
    train_df, test_df = train_test_split_chronological(df, test_size=config.TEST_SIZE)
    print()

    # Step 3: Prepare features
    print("Step 3: Preparing features...")
    X_train, y_train, feature_names = prepare_features(train_df)
    X_test, y_test, _ = prepare_features(test_df)

    print(f"\n  Training target distribution:")
    print(f"    Home Wins: {(y_train == 0).sum()} ({(y_train == 0).mean():.1%})")
    print(f"    Draws: {(y_train == 1).sum()} ({(y_train == 1).mean():.1%})")
    print(f"    Away Wins: {(y_train == 2).sum()} ({(y_train == 2).mean():.1%})")
    print()

    # Step 4: Train model
    print("Step 4: Training XGBoost classifier...")
    model = train_xgboost(X_train, y_train, X_test, y_test)

    # Step 5: Evaluate model
    print("\nStep 5: Evaluating model...")
    metrics = evaluate_model(model, X_test, y_test, feature_names)

    # Step 6: Save model
    print("\nStep 6: Saving model...")
    model_path = os.path.join(config.MODEL_PATH, 'xgboost_v1.pkl')
    save_model(model, feature_names, metrics, model_path)

    # Summary
    print("\n" + "=" * 60)
    print("TRAINING COMPLETE!")
    print("=" * 60)
    print()
    print("Results Summary:")
    print(f"  - Accuracy: {metrics['accuracy']:.1%}")
    print(f"  - Baseline: {metrics['baseline_accuracy']:.1%}")
    print(f"  - Improvement: {(metrics['accuracy'] - metrics['baseline_accuracy']):.1%}")
    print(f"  - Model saved to: {model_path}")
    print()

    # Success criteria check
    target_accuracy = 0.55
    if metrics['accuracy'] >= target_accuracy:
        print(f"[OK] Model meets target accuracy (>={target_accuracy:.0%})")
    else:
        print(f"[WARN] Model below target accuracy ({metrics['accuracy']:.1%} < {target_accuracy:.0%})")
        print("       Consider: more features, hyperparameter tuning, or more data")

    print()
    print(f"Completed at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")


if __name__ == "__main__":
    main()
