"""
Multi-Market Model Training Script

Trains XGBoost models for different betting markets:
- 1x2 (Match Result): 3-class classification
- over_under (Over/Under 2.5 Goals): Binary classification
- btts (Both Teams to Score): Binary classification
"""

import os
import sys
import json
import pickle
import argparse
from datetime import datetime
from pathlib import Path

import pandas as pd
import numpy as np
from sklearn.metrics import (
    accuracy_score, precision_score, recall_score, f1_score,
    classification_report, confusion_matrix, roc_auc_score
)
from xgboost import XGBClassifier

# Add parent to path for imports
sys.path.insert(0, str(Path(__file__).parent))

from app.features.feature_builder import get_feature_columns, get_target_column


def load_data(filepath: str = "data/training_data.csv") -> pd.DataFrame:
    """Load training data from CSV"""
    print(f"\n[1] Loading data from {filepath}...")
    df = pd.read_csv(filepath)

    # Convert match_date to datetime
    df['match_date'] = pd.to_datetime(df['match_date'])

    print(f"    Loaded {len(df)} samples")
    print(f"    Date range: {df['match_date'].min()} to {df['match_date'].max()}")

    return df


def prepare_features(df: pd.DataFrame, market: str) -> tuple:
    """
    Prepare features and target for a specific market

    Args:
        df: DataFrame with all features
        market: '1x2', 'over_under', or 'btts'

    Returns:
        X (features), y (target), feature_names
    """
    print(f"\n[2] Preparing features for {market} market...")

    # Get feature columns for this market
    feature_cols = get_feature_columns(market)
    target_col = get_target_column(market)

    # Find available features
    available_features = [f for f in feature_cols if f in df.columns]
    missing_features = [f for f in feature_cols if f not in df.columns]

    if missing_features:
        print(f"    [WARN] Missing {len(missing_features)} features: {missing_features[:5]}...")

    print(f"    Using {len(available_features)} features")
    print(f"    Target: {target_col}")

    X = df[available_features].copy()
    y = df[target_col].copy()

    # Handle NaN values
    nan_count = X.isna().sum().sum()
    if nan_count > 0:
        print(f"    [WARN] Filling {nan_count} NaN values with column medians")
        X = X.fillna(X.median())

    # Verify target
    print(f"    Target distribution:")
    for val, count in y.value_counts().sort_index().items():
        pct = count / len(y) * 100
        print(f"      {val}: {count} ({pct:.1f}%)")

    return X, y, available_features


def train_test_split_chronological(df: pd.DataFrame, X: pd.DataFrame, y: pd.Series,
                                   test_ratio: float = 0.2) -> tuple:
    """
    Split data chronologically (earlier = train, later = test)
    This is important for time-series data to prevent data leakage
    """
    print(f"\n[3] Splitting data chronologically ({int((1-test_ratio)*100)}/{int(test_ratio*100)} train/test)...")

    split_idx = int(len(df) * (1 - test_ratio))

    X_train = X.iloc[:split_idx]
    X_test = X.iloc[split_idx:]
    y_train = y.iloc[:split_idx]
    y_test = y.iloc[split_idx:]

    train_dates = df['match_date'].iloc[:split_idx]
    test_dates = df['match_date'].iloc[split_idx:]

    print(f"    Train: {len(X_train)} samples ({train_dates.min().date()} to {train_dates.max().date()})")
    print(f"    Test:  {len(X_test)} samples ({test_dates.min().date()} to {test_dates.max().date()})")

    return X_train, X_test, y_train, y_test


def get_model_config(market: str) -> dict:
    """
    Get optimized XGBoost configuration for each market
    """
    configs = {
        '1x2': {
            'objective': 'multi:softprob',
            'num_class': 3,
            'max_depth': 3,
            'learning_rate': 0.05,
            'n_estimators': 200,
            'subsample': 0.7,
            'colsample_bytree': 0.7,
            'min_child_weight': 3,
            'reg_alpha': 0.1,
            'reg_lambda': 1.0,
            'eval_metric': 'mlogloss',
            'random_state': 42,
            'verbosity': 0,
        },
        'over_under': {
            'objective': 'binary:logistic',
            'max_depth': 4,
            'learning_rate': 0.05,
            'n_estimators': 200,
            'subsample': 0.8,
            'colsample_bytree': 0.8,
            'min_child_weight': 2,
            'reg_alpha': 0.05,
            'reg_lambda': 1.0,
            'scale_pos_weight': 1.0,  # Classes are balanced
            'eval_metric': 'logloss',
            'random_state': 42,
            'verbosity': 0,
        },
        'btts': {
            'objective': 'binary:logistic',
            'max_depth': 4,
            'learning_rate': 0.05,
            'n_estimators': 200,
            'subsample': 0.8,
            'colsample_bytree': 0.8,
            'min_child_weight': 2,
            'reg_alpha': 0.05,
            'reg_lambda': 1.0,
            'scale_pos_weight': 1.0,  # Classes are balanced
            'eval_metric': 'logloss',
            'random_state': 42,
            'verbosity': 0,
        }
    }
    return configs.get(market, configs['over_under'])


def train_model(X_train: pd.DataFrame, y_train: pd.Series, market: str) -> XGBClassifier:
    """
    Train XGBoost model for specified market
    """
    print(f"\n[4] Training {market} model...")

    config = get_model_config(market)
    model = XGBClassifier(**config)

    start_time = datetime.now()
    model.fit(X_train, y_train)
    train_time = (datetime.now() - start_time).total_seconds()

    print(f"    Training completed in {train_time:.2f} seconds")

    return model


def evaluate_model(model: XGBClassifier, X_test: pd.DataFrame, y_test: pd.Series,
                   market: str, feature_names: list) -> dict:
    """
    Evaluate model performance
    """
    print(f"\n[5] Evaluating {market} model...")

    y_pred = model.predict(X_test)
    y_proba = model.predict_proba(X_test)

    accuracy = accuracy_score(y_test, y_pred)

    # Calculate baseline (always predicting majority class)
    majority_class = y_test.mode()[0]
    baseline_accuracy = (y_test == majority_class).mean()
    improvement = accuracy - baseline_accuracy

    print(f"\n    Results:")
    print(f"    ========")
    print(f"    Accuracy:  {accuracy:.4f} ({accuracy*100:.1f}%)")
    print(f"    Baseline:  {baseline_accuracy:.4f} ({baseline_accuracy*100:.1f}%)")
    print(f"    Improvement: +{improvement:.4f} ({improvement*100:.1f}%)")

    if market == '1x2':
        # Multi-class metrics
        print(f"\n    Classification Report:")
        labels = ['Home Win', 'Draw', 'Away Win']
        print(classification_report(y_test, y_pred, target_names=labels))
    else:
        # Binary metrics
        precision = precision_score(y_test, y_pred, zero_division=0)
        recall = recall_score(y_test, y_pred, zero_division=0)
        f1 = f1_score(y_test, y_pred, zero_division=0)

        # ROC AUC for binary
        roc_auc = roc_auc_score(y_test, y_proba[:, 1])

        print(f"    Precision: {precision:.4f}")
        print(f"    Recall:    {recall:.4f}")
        print(f"    F1 Score:  {f1:.4f}")
        print(f"    ROC AUC:   {roc_auc:.4f}")

        label_map = {
            'over_under': ['Under 2.5', 'Over 2.5'],
            'btts': ['No (FTS)', 'Yes (BTTS)']
        }
        labels = label_map.get(market, ['0', '1'])
        print(f"\n    Classification Report:")
        print(classification_report(y_test, y_pred, target_names=labels))

    # Confusion matrix
    print(f"\n    Confusion Matrix:")
    cm = confusion_matrix(y_test, y_pred)
    print(cm)

    # Feature importance
    print(f"\n    Top 15 Feature Importances:")
    importance = model.feature_importances_
    feature_importance = pd.DataFrame({
        'feature': feature_names,
        'importance': importance
    }).sort_values('importance', ascending=False)

    for _, row in feature_importance.head(15).iterrows():
        print(f"      {row['feature']}: {row['importance']:.4f}")

    metrics = {
        'market': market,
        'accuracy': accuracy,
        'baseline_accuracy': baseline_accuracy,
        'improvement': improvement,
        'train_samples': len(y_test) * 4,  # Approximate
        'test_samples': len(y_test),
    }

    if market != '1x2':
        metrics.update({
            'precision': precision,
            'recall': recall,
            'f1_score': f1,
            'roc_auc': roc_auc,
        })

    return metrics


def save_model(model: XGBClassifier, feature_names: list, metrics: dict,
               market: str, models_dir: str = "models"):
    """
    Save trained model and metadata
    """
    print(f"\n[6] Saving {market} model...")

    os.makedirs(models_dir, exist_ok=True)

    # Model filename includes market type
    model_path = os.path.join(models_dir, f"xgboost_{market}_model.pkl")

    model_data = {
        'model': model,
        'feature_names': feature_names,
        'market': market,
        'training_date': datetime.now().isoformat(),
        'metrics': metrics,
        'version': f"1.0.0-{market}",
    }

    with open(model_path, 'wb') as f:
        pickle.dump(model_data, f)

    print(f"    Model saved to: {model_path}")

    # Save metrics separately
    metrics_path = os.path.join(models_dir, f"metrics_{market}.json")
    with open(metrics_path, 'w') as f:
        json.dump(metrics, f, indent=2)

    print(f"    Metrics saved to: {metrics_path}")

    return model_path


def train_market(market: str, data_path: str = "data/training_data.csv"):
    """
    Train a model for a specific market
    """
    print("\n" + "=" * 60)
    print(f"Training Model: {market.upper()}")
    print("=" * 60)

    # Load data
    df = load_data(data_path)

    # Prepare features
    X, y, feature_names = prepare_features(df, market)

    # Split chronologically
    X_train, X_test, y_train, y_test = train_test_split_chronological(df, X, y)

    # Train model
    model = train_model(X_train, y_train, market)

    # Evaluate
    metrics = evaluate_model(model, X_test, y_test, market, feature_names)

    # Save
    model_path = save_model(model, feature_names, metrics, market)

    print("\n" + "=" * 60)
    print(f"[OK] {market.upper()} MODEL TRAINING COMPLETE!")
    print("=" * 60)
    print(f"  Accuracy: {metrics['accuracy']*100:.1f}%")
    print(f"  Baseline: {metrics['baseline_accuracy']*100:.1f}%")
    print(f"  Improvement: +{metrics['improvement']*100:.1f}%")
    print(f"  Model: {model_path}")

    return metrics


def main():
    parser = argparse.ArgumentParser(description='Train betting market prediction models')
    parser.add_argument('--market', type=str, default='all',
                       choices=['1x2', 'over_under', 'btts', 'all'],
                       help='Which market model to train')
    parser.add_argument('--data', type=str, default='data/training_data.csv',
                       help='Path to training data CSV')

    args = parser.parse_args()

    markets_to_train = []
    if args.market == 'all':
        markets_to_train = ['1x2', 'over_under', 'btts']
    else:
        markets_to_train = [args.market]

    results = {}
    for market in markets_to_train:
        metrics = train_market(market, args.data)
        results[market] = metrics

    # Summary
    if len(results) > 1:
        print("\n" + "=" * 60)
        print("TRAINING SUMMARY")
        print("=" * 60)
        for market, metrics in results.items():
            print(f"\n  {market.upper()}:")
            print(f"    Accuracy: {metrics['accuracy']*100:.1f}%")
            print(f"    vs Baseline: +{metrics['improvement']*100:.1f}%")


if __name__ == "__main__":
    main()
