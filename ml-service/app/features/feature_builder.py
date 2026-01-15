"""
Feature Builder

Orchestrates all feature extraction for model training
"""
from typing import List, Dict, Any
import pandas as pd
from datetime import datetime

from app.features.form_metrics import get_form_features
from app.features.h2h_stats import get_h2h_features
from app.features.league_position import get_league_position_features


def extract_features_for_fixture(
    fixture: Dict[str, Any],
    all_fixtures: List[Dict[str, Any]]
) -> Dict[str, Any]:
    """
    Extract all features for a single fixture

    Args:
        fixture: The fixture to extract features for
        all_fixtures: All historical fixtures (for context)

    Returns:
        Dictionary of features
    """
    home_team_id = fixture['home_team_id']
    away_team_id = fixture['away_team_id']
    match_date = fixture['match_date']
    season = fixture['season']

    features = {
        'fixture_id': fixture['id'],
        'season': season,
        'match_date': match_date,
        'home_team_id': home_team_id,
        'away_team_id': away_team_id,
    }

    # Form features
    form_features = get_form_features(home_team_id, away_team_id, match_date, all_fixtures)
    features.update(form_features)

    # Head-to-head features
    h2h_features = get_h2h_features(home_team_id, away_team_id, match_date, all_fixtures)
    features.update(h2h_features)

    # League position features
    position_features = get_league_position_features(
        home_team_id, away_team_id, season, match_date, all_fixtures
    )
    features.update(position_features)

    # Create label (outcome)
    if fixture['home_score'] is not None and fixture['away_score'] is not None:
        if fixture['home_score'] > fixture['away_score']:
            features['outcome'] = 'home_win'
            features['outcome_encoded'] = 0
        elif fixture['home_score'] < fixture['away_score']:
            features['outcome'] = 'away_win'
            features['outcome_encoded'] = 2
        else:
            features['outcome'] = 'draw'
            features['outcome_encoded'] = 1

        # Also include actual scores for analysis
        features['home_score'] = fixture['home_score']
        features['away_score'] = fixture['away_score']
        features['total_goals'] = fixture['home_score'] + fixture['away_score']

    return features


def build_training_dataset(
    fixtures: List[Dict[str, Any]],
    verbose: bool = True
) -> pd.DataFrame:
    """
    Build complete training dataset with features for all fixtures

    Args:
        fixtures: List of fixtures to extract features for
        verbose: Print progress

    Returns:
        DataFrame with features and labels
    """
    if verbose:
        print(f"Building features for {len(fixtures)} fixtures...")

    all_features = []

    for i, fixture in enumerate(fixtures):
        if verbose and (i + 1) % 100 == 0:
            print(f"  Processed {i + 1}/{len(fixtures)} fixtures...")

        try:
            features = extract_features_for_fixture(fixture, fixtures)
            all_features.append(features)
        except Exception as e:
            if verbose:
                print(f"  Warning: Failed to extract features for fixture {fixture.get('id')}: {e}")
            continue

    if verbose:
        print(f"[OK] Feature extraction complete! Built {len(all_features)} samples")

    df = pd.DataFrame(all_features)

    # Sort by date
    df = df.sort_values('match_date').reset_index(drop=True)

    return df


def get_feature_columns() -> List[str]:
    """
    Get list of feature column names (excluding id, date, label columns)

    Returns:
        List of feature column names
    """
    # These are all the feature columns we extract
    feature_columns = [
        # Form features (last 5)
        'home_form_last_5_points', 'home_form_last_5_wins', 'home_form_last_5_draws',
        'home_form_last_5_losses', 'home_form_last_5_goals_scored', 'home_form_last_5_goals_conceded',
        'home_form_last_5_goal_diff', 'home_form_last_5_clean_sheets', 'home_form_last_5_failed_to_score',
        'home_form_last_5_avg_points', 'home_form_last_5_avg_goals_scored', 'home_form_last_5_avg_goals_conceded',

        'away_form_last_5_points', 'away_form_last_5_wins', 'away_form_last_5_draws',
        'away_form_last_5_losses', 'away_form_last_5_goals_scored', 'away_form_last_5_goals_conceded',
        'away_form_last_5_goal_diff', 'away_form_last_5_clean_sheets', 'away_form_last_5_failed_to_score',
        'away_form_last_5_avg_points', 'away_form_last_5_avg_goals_scored', 'away_form_last_5_avg_goals_conceded',

        # Home/Away specific form (last 3)
        'home_form_last_3_points', 'home_form_last_3_wins', 'home_form_last_3_goals_scored',
        'home_form_last_3_goals_conceded', 'home_form_last_3_avg_points', 'home_form_last_3_avg_goals_scored',

        'away_form_last_3_points', 'away_form_last_3_wins', 'away_form_last_3_goals_scored',
        'away_form_last_3_goals_conceded', 'away_form_last_3_avg_points', 'away_form_last_3_avg_goals_scored',

        # Form differentials
        'form_points_diff', 'form_goals_scored_diff', 'form_goal_diff_diff',

        # H2H features
        'h2h_games_played', 'h2h_home_wins', 'h2h_away_wins', 'h2h_draws',
        'h2h_home_goals_scored', 'h2h_away_goals_scored', 'h2h_goal_diff',
        'h2h_home_win_pct', 'h2h_away_win_pct', 'h2h_draw_pct',
        'h2h_avg_total_goals', 'h2h_avg_home_goals', 'h2h_avg_away_goals',
        'h2h_home_as_home_wins', 'h2h_home_as_home_games', 'h2h_home_as_home_win_pct',

        # League position features
        'home_position', 'away_position', 'position_diff',
        'home_points', 'away_points', 'points_diff',
        'home_ppg', 'away_ppg', 'ppg_diff',
        'home_season_goals_for', 'away_season_goals_for',
        'home_season_goals_against', 'away_season_goals_against',
        'home_avg_goals_scored', 'away_avg_goals_scored',
        'home_avg_goals_conceded', 'away_avg_goals_conceded',
        'home_goal_diff', 'away_goal_diff', 'goal_diff_diff',
        'home_games_played', 'away_games_played',
        'home_win_pct', 'away_win_pct',
    ]

    return feature_columns


def get_target_column() -> str:
    """Get the target column name"""
    return 'outcome_encoded'


if __name__ == "__main__":
    # Test feature extraction
    print("Testing feature extraction...")
    print(f"Total feature count: {len(get_feature_columns())}")
    print(f"Target column: {get_target_column()}")
