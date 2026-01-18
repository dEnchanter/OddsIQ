"""
Goal-Specific Features for Over/Under and BTTS Models

Features focused on scoring patterns, clean sheets, and goal metrics
"""
from typing import List, Dict, Any
from datetime import datetime


def calculate_goal_metrics(
    team_id: int,
    current_date: datetime,
    all_fixtures: List[Dict[str, Any]],
    n_games: int = 10
) -> Dict[str, float]:
    """
    Calculate goal-related metrics for a team

    Args:
        team_id: Team ID
        current_date: Date to calculate before
        all_fixtures: All fixtures
        n_games: Number of recent games to consider

    Returns:
        Dictionary of goal metrics
    """
    # Filter completed fixtures for this team before current date
    team_fixtures = [
        f for f in all_fixtures
        if (f['home_team_id'] == team_id or f['away_team_id'] == team_id)
        and f['match_date'] < current_date
        and f['status'] == 'FT'
        and f['home_score'] is not None
    ]

    # Sort by date descending and take last n games
    team_fixtures.sort(key=lambda x: x['match_date'], reverse=True)
    recent_fixtures = team_fixtures[:n_games]

    if len(recent_fixtures) == 0:
        return {
            'goals_scored_avg': 0.0,
            'goals_conceded_avg': 0.0,
            'total_goals_avg': 0.0,
            'over_2_5_pct': 0.0,
            'over_1_5_pct': 0.0,
            'over_3_5_pct': 0.0,
            'btts_pct': 0.0,
            'clean_sheet_pct': 0.0,
            'failed_to_score_pct': 0.0,
            'scored_first_half_pct': 0.0,
            'conceded_first_half_pct': 0.0,
            'games_analyzed': 0.0,
        }

    total_scored = 0
    total_conceded = 0
    over_2_5_count = 0
    over_1_5_count = 0
    over_3_5_count = 0
    btts_count = 0
    clean_sheets = 0
    failed_to_score = 0

    for fixture in recent_fixtures:
        is_home = fixture['home_team_id'] == team_id
        team_goals = fixture['home_score'] if is_home else fixture['away_score']
        opponent_goals = fixture['away_score'] if is_home else fixture['home_score']
        total_goals = team_goals + opponent_goals

        total_scored += team_goals
        total_conceded += opponent_goals

        # Over/Under thresholds
        if total_goals > 2.5:
            over_2_5_count += 1
        if total_goals > 1.5:
            over_1_5_count += 1
        if total_goals > 3.5:
            over_3_5_count += 1

        # BTTS
        if team_goals > 0 and opponent_goals > 0:
            btts_count += 1

        # Clean sheets and failed to score
        if opponent_goals == 0:
            clean_sheets += 1
        if team_goals == 0:
            failed_to_score += 1

    games = len(recent_fixtures)

    return {
        'goals_scored_avg': total_scored / games,
        'goals_conceded_avg': total_conceded / games,
        'total_goals_avg': (total_scored + total_conceded) / games,
        'over_2_5_pct': over_2_5_count / games,
        'over_1_5_pct': over_1_5_count / games,
        'over_3_5_pct': over_3_5_count / games,
        'btts_pct': btts_count / games,
        'clean_sheet_pct': clean_sheets / games,
        'failed_to_score_pct': failed_to_score / games,
        'games_analyzed': float(games),
    }


def calculate_home_away_goal_metrics(
    team_id: int,
    current_date: datetime,
    all_fixtures: List[Dict[str, Any]],
    is_home: bool,
    n_games: int = 5
) -> Dict[str, float]:
    """
    Calculate goal metrics for home or away games only

    Args:
        team_id: Team ID
        current_date: Date to calculate before
        all_fixtures: All fixtures
        is_home: True for home games, False for away
        n_games: Number of games to consider

    Returns:
        Location-specific goal metrics
    """
    # Filter for home or away games
    if is_home:
        team_fixtures = [
            f for f in all_fixtures
            if f['home_team_id'] == team_id
            and f['match_date'] < current_date
            and f['status'] == 'FT'
            and f['home_score'] is not None
        ]
    else:
        team_fixtures = [
            f for f in all_fixtures
            if f['away_team_id'] == team_id
            and f['match_date'] < current_date
            and f['status'] == 'FT'
            and f['away_score'] is not None
        ]

    team_fixtures.sort(key=lambda x: x['match_date'], reverse=True)
    recent_fixtures = team_fixtures[:n_games]

    location = 'home' if is_home else 'away'

    if len(recent_fixtures) == 0:
        return {
            f'{location}_goals_scored_avg': 0.0,
            f'{location}_goals_conceded_avg': 0.0,
            f'{location}_total_goals_avg': 0.0,
            f'{location}_over_2_5_pct': 0.0,
            f'{location}_btts_pct': 0.0,
            f'{location}_clean_sheet_pct': 0.0,
            f'{location}_failed_to_score_pct': 0.0,
        }

    total_scored = 0
    total_conceded = 0
    over_2_5_count = 0
    btts_count = 0
    clean_sheets = 0
    failed_to_score = 0

    for fixture in recent_fixtures:
        team_goals = fixture['home_score'] if is_home else fixture['away_score']
        opponent_goals = fixture['away_score'] if is_home else fixture['home_score']
        total_goals = team_goals + opponent_goals

        total_scored += team_goals
        total_conceded += opponent_goals

        if total_goals > 2.5:
            over_2_5_count += 1
        if team_goals > 0 and opponent_goals > 0:
            btts_count += 1
        if opponent_goals == 0:
            clean_sheets += 1
        if team_goals == 0:
            failed_to_score += 1

    games = len(recent_fixtures)

    return {
        f'{location}_goals_scored_avg': total_scored / games,
        f'{location}_goals_conceded_avg': total_conceded / games,
        f'{location}_total_goals_avg': (total_scored + total_conceded) / games,
        f'{location}_over_2_5_pct': over_2_5_count / games,
        f'{location}_btts_pct': btts_count / games,
        f'{location}_clean_sheet_pct': clean_sheets / games,
        f'{location}_failed_to_score_pct': failed_to_score / games,
    }


def calculate_h2h_goal_metrics(
    home_team_id: int,
    away_team_id: int,
    current_date: datetime,
    all_fixtures: List[Dict[str, Any]],
    n_games: int = 5
) -> Dict[str, float]:
    """
    Calculate goal-related head-to-head metrics

    Args:
        home_team_id: Home team ID
        away_team_id: Away team ID
        current_date: Date to calculate before
        all_fixtures: All fixtures
        n_games: Number of H2H games to consider

    Returns:
        H2H goal metrics
    """
    # Find H2H matches
    h2h_fixtures = [
        f for f in all_fixtures
        if ((f['home_team_id'] == home_team_id and f['away_team_id'] == away_team_id) or
            (f['home_team_id'] == away_team_id and f['away_team_id'] == home_team_id))
        and f['match_date'] < current_date
        and f['status'] == 'FT'
        and f['home_score'] is not None
    ]

    h2h_fixtures.sort(key=lambda x: x['match_date'], reverse=True)
    recent_h2h = h2h_fixtures[:n_games]

    if len(recent_h2h) == 0:
        return {
            'h2h_total_goals_avg': 0.0,
            'h2h_over_2_5_pct': 0.0,
            'h2h_btts_pct': 0.0,
            'h2h_home_team_goals_avg': 0.0,
            'h2h_away_team_goals_avg': 0.0,
        }

    total_goals = 0
    home_team_goals = 0  # Goals scored by current home team
    away_team_goals = 0  # Goals scored by current away team
    over_2_5_count = 0
    btts_count = 0

    for fixture in recent_h2h:
        match_total = fixture['home_score'] + fixture['away_score']
        total_goals += match_total

        # Track goals by current home/away team
        if fixture['home_team_id'] == home_team_id:
            home_team_goals += fixture['home_score']
            away_team_goals += fixture['away_score']
        else:
            home_team_goals += fixture['away_score']
            away_team_goals += fixture['home_score']

        if match_total > 2.5:
            over_2_5_count += 1
        if fixture['home_score'] > 0 and fixture['away_score'] > 0:
            btts_count += 1

    games = len(recent_h2h)

    return {
        'h2h_total_goals_avg': total_goals / games,
        'h2h_over_2_5_pct': over_2_5_count / games,
        'h2h_btts_pct': btts_count / games,
        'h2h_home_team_goals_avg': home_team_goals / games,
        'h2h_away_team_goals_avg': away_team_goals / games,
    }


def get_goal_features(
    home_team_id: int,
    away_team_id: int,
    match_date: datetime,
    all_fixtures: List[Dict[str, Any]]
) -> Dict[str, float]:
    """
    Get all goal-related features for a fixture

    Args:
        home_team_id: Home team ID
        away_team_id: Away team ID
        match_date: Match date
        all_fixtures: All fixtures

    Returns:
        Combined goal features dictionary
    """
    features = {}

    # Overall goal metrics (last 10 games)
    home_goals = calculate_goal_metrics(home_team_id, match_date, all_fixtures, n_games=10)
    away_goals = calculate_goal_metrics(away_team_id, match_date, all_fixtures, n_games=10)

    for key, value in home_goals.items():
        features[f'home_{key}'] = value
    for key, value in away_goals.items():
        features[f'away_{key}'] = value

    # Home/away specific metrics (home team when playing at home, away team when playing away)
    home_at_home = calculate_home_away_goal_metrics(home_team_id, match_date, all_fixtures, is_home=True, n_games=5)
    away_at_away = calculate_home_away_goal_metrics(away_team_id, match_date, all_fixtures, is_home=False, n_games=5)

    # Prefix with 'home_' for home team's home stats, 'away_' for away team's away stats
    for key, value in home_at_home.items():
        features[f'home_{key}'] = value
    for key, value in away_at_away.items():
        features[f'away_{key}'] = value

    # H2H goal metrics
    h2h_goals = calculate_h2h_goal_metrics(home_team_id, away_team_id, match_date, all_fixtures, n_games=5)
    features.update(h2h_goals)

    # Combined/differential features
    features['combined_goals_avg'] = features['home_goals_scored_avg'] + features['away_goals_scored_avg']
    features['combined_conceded_avg'] = features['home_goals_conceded_avg'] + features['away_goals_conceded_avg']
    features['combined_total_goals_avg'] = features['home_total_goals_avg'] + features['away_total_goals_avg']

    # Expected goals in match (simple estimate)
    features['expected_total_goals'] = (
        features['home_goals_scored_avg'] + features['away_goals_conceded_avg'] +
        features['away_goals_scored_avg'] + features['home_goals_conceded_avg']
    ) / 2

    # Over 2.5 indicators (using overall stats)
    features['both_over_2_5_pct'] = (features['home_over_2_5_pct'] + features['away_over_2_5_pct']) / 2

    # BTTS indicators (using overall stats)
    features['both_btts_pct'] = (features['home_btts_pct'] + features['away_btts_pct']) / 2
    features['btts_potential'] = (
        (1 - features['home_clean_sheet_pct']) * (1 - features['away_failed_to_score_pct']) +
        (1 - features['away_clean_sheet_pct']) * (1 - features['home_failed_to_score_pct'])
    ) / 2

    return features
