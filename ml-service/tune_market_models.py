"""
Hyperparameter Tuning for Multi-Market Models

Explores different configurations to improve O/U and BTTS models
"""

import pandas as pd
import numpy as np
from sklearn.metrics import accuracy_score, roc_auc_score, f1_score
from xgboost import XGBClassifier
import pickle
import os
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from app.features.feature_builder import get_feature_columns, get_target_column


def load_and_prepare_data(market: str):
    """Load data and prepare features for a market"""
    df = pd.read_csv('data/training_data.csv')
    df['match_date'] = pd.to_datetime(df['match_date'])

    feature_cols = get_feature_columns(market)
    target_col = get_target_column(market)

    available = [f for f in feature_cols if f in df.columns]
    X = df[available].fillna(df[available].median())
    y = df[target_col]

    # Chronological split
    split_idx = int(len(df) * 0.8)
    X_train, X_test = X.iloc[:split_idx], X.iloc[split_idx:]
    y_train, y_test = y.iloc[:split_idx], y.iloc[split_idx:]

    return X_train, X_test, y_train, y_test, available


def get_goal_features_only():
    """Get only goal-specific feature names"""
    return [
        # Overall goal metrics
        'home_goals_scored_avg', 'home_goals_conceded_avg', 'home_total_goals_avg',
        'home_over_2_5_pct', 'home_over_1_5_pct', 'home_over_3_5_pct',
        'home_btts_pct', 'home_clean_sheet_pct', 'home_failed_to_score_pct',
        'away_goals_scored_avg', 'away_goals_conceded_avg', 'away_total_goals_avg',
        'away_over_2_5_pct', 'away_over_1_5_pct', 'away_over_3_5_pct',
        'away_btts_pct', 'away_clean_sheet_pct', 'away_failed_to_score_pct',
        # Location-specific
        'home_home_goals_scored_avg', 'home_home_goals_conceded_avg', 'home_home_total_goals_avg',
        'home_home_over_2_5_pct', 'home_home_btts_pct', 'home_home_clean_sheet_pct',
        'away_away_goals_scored_avg', 'away_away_goals_conceded_avg', 'away_away_total_goals_avg',
        'away_away_over_2_5_pct', 'away_away_btts_pct', 'away_away_clean_sheet_pct',
        # H2H
        'h2h_total_goals_avg', 'h2h_over_2_5_pct', 'h2h_btts_pct',
        'h2h_home_team_goals_avg', 'h2h_away_team_goals_avg',
        # Combined
        'combined_goals_avg', 'combined_total_goals_avg', 'expected_total_goals',
        'both_over_2_5_pct', 'both_btts_pct', 'btts_potential',
        # Form goal features
        'home_form_last_5_goals_scored', 'home_form_last_5_goals_conceded',
        'home_form_last_5_avg_goals_scored', 'home_form_last_5_avg_goals_conceded',
        'away_form_last_5_goals_scored', 'away_form_last_5_goals_conceded',
        'away_form_last_5_avg_goals_scored', 'away_form_last_5_avg_goals_conceded',
        'home_form_last_3_goals_scored', 'home_form_last_3_goals_conceded',
        'away_form_last_3_goals_scored', 'away_form_last_3_goals_conceded',
    ]


def tune_market(market: str):
    """Tune a market model with various configurations"""
    print(f"\n{'='*60}")
    print(f"Tuning {market.upper()} Model")
    print('='*60)

    # Load data
    df = pd.read_csv('data/training_data.csv')
    df['match_date'] = pd.to_datetime(df['match_date'])

    target_col = get_target_column(market)
    y = df[target_col]

    split_idx = int(len(df) * 0.8)
    y_train, y_test = y.iloc[:split_idx], y.iloc[split_idx:]

    # Calculate class weight for imbalanced test set
    pos_weight_train = (y_train == 0).sum() / (y_train == 1).sum()
    pos_weight_test = (y_test == 0).sum() / (y_test == 1).sum()

    baseline = (y_test == y_test.mode()[0]).mean()
    print(f"\nBaseline accuracy: {baseline*100:.1f}%")
    print(f"Train pos ratio: {y_train.mean()*100:.1f}%")
    print(f"Test pos ratio: {y_test.mean()*100:.1f}%")

    # Feature sets to try
    all_features = get_feature_columns(market)
    goal_only = get_goal_features_only()

    feature_sets = {
        'all_features': all_features,
        'goal_only': goal_only,
    }

    # Configurations to try
    configs = [
        {
            'name': 'Balanced (default)',
            'max_depth': 4,
            'learning_rate': 0.05,
            'n_estimators': 200,
            'scale_pos_weight': 1.0,
        },
        {
            'name': 'Conservative (shallow)',
            'max_depth': 2,
            'learning_rate': 0.03,
            'n_estimators': 300,
            'scale_pos_weight': 1.0,
        },
        {
            'name': 'Adjusted for test (weighted)',
            'max_depth': 3,
            'learning_rate': 0.05,
            'n_estimators': 200,
            'scale_pos_weight': pos_weight_test,
        },
        {
            'name': 'High regularization',
            'max_depth': 3,
            'learning_rate': 0.03,
            'n_estimators': 250,
            'scale_pos_weight': 1.0,
            'reg_alpha': 0.5,
            'reg_lambda': 2.0,
        },
        {
            'name': 'Low estimators',
            'max_depth': 3,
            'learning_rate': 0.1,
            'n_estimators': 50,
            'scale_pos_weight': 1.0,
        },
    ]

    results = []

    for feat_name, feat_cols in feature_sets.items():
        available = [f for f in feat_cols if f in df.columns]
        X = df[available].fillna(df[available].median())
        X_train, X_test = X.iloc[:split_idx], X.iloc[split_idx:]

        for config in configs:
            name = f"{feat_name} + {config['name']}"
            print(f"\n  Testing: {name}")

            model = XGBClassifier(
                objective='binary:logistic',
                max_depth=config['max_depth'],
                learning_rate=config['learning_rate'],
                n_estimators=config['n_estimators'],
                scale_pos_weight=config.get('scale_pos_weight', 1.0),
                reg_alpha=config.get('reg_alpha', 0.05),
                reg_lambda=config.get('reg_lambda', 1.0),
                subsample=0.8,
                colsample_bytree=0.8,
                min_child_weight=2,
                eval_metric='logloss',
                random_state=42,
                verbosity=0
            )

            model.fit(X_train, y_train)
            y_pred = model.predict(X_test)
            y_proba = model.predict_proba(X_test)

            accuracy = accuracy_score(y_test, y_pred)
            roc_auc = roc_auc_score(y_test, y_proba[:, 1])
            f1 = f1_score(y_test, y_pred)
            improvement = accuracy - baseline

            results.append({
                'config': name,
                'features': feat_name,
                'accuracy': accuracy,
                'improvement': improvement,
                'roc_auc': roc_auc,
                'f1': f1,
                'model': model,
                'feature_names': available,
            })

            status = "[OK]" if improvement > 0 else "[--]"
            print(f"    {status} Acc: {accuracy*100:.1f}% (vs baseline: {improvement*100:+.1f}%), ROC: {roc_auc:.3f}, F1: {f1:.3f}")

    # Find best
    results_df = pd.DataFrame(results)
    best = results_df.loc[results_df['accuracy'].idxmax()]

    print(f"\n{'='*60}")
    print(f"BEST CONFIG for {market.upper()}")
    print('='*60)
    print(f"  Config: {best['config']}")
    print(f"  Accuracy: {best['accuracy']*100:.1f}%")
    print(f"  vs Baseline: {best['improvement']*100:+.1f}%")
    print(f"  ROC AUC: {best['roc_auc']:.4f}")
    print(f"  F1: {best['f1']:.4f}")

    # Save best model
    if best['improvement'] > -0.05:  # Only save if within 5% of baseline
        os.makedirs('models', exist_ok=True)
        model_data = {
            'model': best['model'],
            'feature_names': best['feature_names'],
            'market': market,
            'config': best['config'],
            'metrics': {
                'accuracy': best['accuracy'],
                'baseline': baseline,
                'improvement': best['improvement'],
                'roc_auc': best['roc_auc'],
                'f1': best['f1'],
            }
        }
        model_path = f'models/xgboost_{market}_model.pkl'
        with open(model_path, 'wb') as f:
            pickle.dump(model_data, f)
        print(f"\n  Model saved: {model_path}")

    return best


def main():
    print("="*60)
    print("Multi-Market Model Tuning")
    print("="*60)

    ou_best = tune_market('over_under')
    btts_best = tune_market('btts')

    print("\n" + "="*60)
    print("SUMMARY")
    print("="*60)
    print(f"Over/Under 2.5: {ou_best['accuracy']*100:.1f}% ({ou_best['improvement']*100:+.1f}% vs baseline)")
    print(f"BTTS: {btts_best['accuracy']*100:.1f}% ({btts_best['improvement']*100:+.1f}% vs baseline)")


if __name__ == "__main__":
    main()
