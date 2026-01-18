"""
Feature Builder

Orchestrates all feature extraction for model training.
Supports multiple markets: 1X2, Over/Under 2.5, BTTS
"""
from typing import List, Dict, Any
import pandas as pd
from datetime import datetime

from app.features.form_metrics import get_form_features
from app.features.h2h_stats import get_h2h_features
from app.features.league_position import get_league_position_features
from app.features.goal_features import get_goal_features


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

    # Goal-specific features (for O/U and BTTS models)
    goal_features = get_goal_features(home_team_id, away_team_id, match_date, all_fixtures)
    features.update(goal_features)

    # Create labels for all markets
    if fixture['home_score'] is not None and fixture['away_score'] is not None:
        home_score = fixture['home_score']
        away_score = fixture['away_score']
        total_goals = home_score + away_score

        # 1X2 Market (Match Result)
        if home_score > away_score:
            features['outcome'] = 'home_win'
            features['outcome_encoded'] = 0
        elif home_score < away_score:
            features['outcome'] = 'away_win'
            features['outcome_encoded'] = 2
        else:
            features['outcome'] = 'draw'
            features['outcome_encoded'] = 1

        # Over/Under 2.5 Goals Market
        features['over_2_5'] = 1 if total_goals > 2.5 else 0
        features['over_1_5'] = 1 if total_goals > 1.5 else 0
        features['over_3_5'] = 1 if total_goals > 3.5 else 0

        # BTTS (Both Teams To Score) Market
        features['btts'] = 1 if (home_score > 0 and away_score > 0) else 0

        # Also include actual scores for analysis
        features['home_score'] = home_score
        features['away_score'] = away_score
        features['total_goals'] = total_goals

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


def get_feature_columns(market: str = 'all') -> List[str]:
    """
    Get list of feature column names (excluding id, date, label columns)

    Args:
        market: Which market's features to return
                'all' - all features
                '1x2' - features for match result prediction
                'over_under' - features optimized for O/U prediction
                'btts' - features optimized for BTTS prediction

    Returns:
        List of feature column names
    """
    # Base features (form, H2H, position)
    base_features = [
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

    # Goal-specific features (for O/U and BTTS)
    goal_features = [
        # Home team overall goal metrics (last 10 games, any venue)
        'home_goals_scored_avg', 'home_goals_conceded_avg', 'home_total_goals_avg',
        'home_over_2_5_pct', 'home_over_1_5_pct', 'home_over_3_5_pct',
        'home_btts_pct', 'home_clean_sheet_pct', 'home_failed_to_score_pct',
        'home_games_analyzed',

        # Away team overall goal metrics (last 10 games, any venue)
        'away_goals_scored_avg', 'away_goals_conceded_avg', 'away_total_goals_avg',
        'away_over_2_5_pct', 'away_over_1_5_pct', 'away_over_3_5_pct',
        'away_btts_pct', 'away_clean_sheet_pct', 'away_failed_to_score_pct',
        'away_games_analyzed',

        # Home team at home specific (last 5 home games)
        'home_home_goals_scored_avg', 'home_home_goals_conceded_avg', 'home_home_total_goals_avg',
        'home_home_over_2_5_pct', 'home_home_btts_pct', 'home_home_clean_sheet_pct',
        'home_home_failed_to_score_pct',

        # Away team away specific (last 5 away games)
        'away_away_goals_scored_avg', 'away_away_goals_conceded_avg', 'away_away_total_goals_avg',
        'away_away_over_2_5_pct', 'away_away_btts_pct', 'away_away_clean_sheet_pct',
        'away_away_failed_to_score_pct',

        # H2H goal metrics
        'h2h_total_goals_avg', 'h2h_over_2_5_pct', 'h2h_btts_pct',
        'h2h_home_team_goals_avg', 'h2h_away_team_goals_avg',

        # Combined/differential goal features
        'combined_goals_avg', 'combined_conceded_avg', 'combined_total_goals_avg',
        'expected_total_goals', 'both_over_2_5_pct',
        'both_btts_pct', 'btts_potential',
    ]

    if market == '1x2':
        return base_features
    elif market == 'over_under':
        # Combine base + goal features, prioritizing goal-related
        return base_features + goal_features
    elif market == 'btts':
        # Same as over_under - all features available
        return base_features + goal_features
    else:  # 'all'
        return base_features + goal_features


def get_target_column(market: str = '1x2') -> str:
    """
    Get the target column name for a specific market

    Args:
        market: '1x2', 'over_under', or 'btts'

    Returns:
        Target column name
    """
    targets = {
        '1x2': 'outcome_encoded',
        'over_under': 'over_2_5',
        'over_1_5': 'over_1_5',
        'over_3_5': 'over_3_5',
        'btts': 'btts',
    }
    return targets.get(market, 'outcome_encoded')


def get_all_targets() -> Dict[str, str]:
    """
    Get all available target columns

    Returns:
        Dictionary mapping market name to target column
    """
    return {
        '1x2': 'outcome_encoded',
        'over_2_5': 'over_2_5',
        'over_1_5': 'over_1_5',
        'over_3_5': 'over_3_5',
        'btts': 'btts',
    }


if __name__ == "__main__":
    # Test feature extraction
    print("Testing feature extraction...")
    print(f"Base feature count (1X2): {len(get_feature_columns('1x2'))}")
    print(f"Goal feature count (O/U, BTTS): {len(get_feature_columns('over_under'))}")
    print(f"All features: {len(get_feature_columns('all'))}")
    print(f"\nTarget columns:")
    for market, target in get_all_targets().items():
        print(f"  {market}: {target}")
