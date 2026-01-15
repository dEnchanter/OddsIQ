"""
Form Metrics Feature Engineering

Calculates team form features based on last N games
"""
from typing import List, Dict, Any, Tuple
from datetime import datetime, timedelta


def calculate_team_form_last_n(
    team_id: int,
    current_date: datetime,
    all_fixtures: List[Dict[str, Any]],
    n_games: int = 5
) -> Dict[str, float]:
    """
    Calculate team form metrics for last N games before current_date

    Args:
        team_id: Team ID to calculate form for
        current_date: Date to calculate form before (exclusive)
        all_fixtures: All fixtures sorted by date
        n_games: Number of recent games to consider

    Returns:
        Dictionary of form metrics
    """
    # Filter fixtures for this team before current date
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
        # No previous games - return zeros
        return {
            f'form_last_{n_games}_points': 0.0,
            f'form_last_{n_games}_wins': 0.0,
            f'form_last_{n_games}_draws': 0.0,
            f'form_last_{n_games}_losses': 0.0,
            f'form_last_{n_games}_goals_scored': 0.0,
            f'form_last_{n_games}_goals_conceded': 0.0,
            f'form_last_{n_games}_goal_diff': 0.0,
            f'form_last_{n_games}_clean_sheets': 0.0,
            f'form_last_{n_games}_failed_to_score': 0.0,
            f'form_last_{n_games}_games_played': 0.0,
        }

    # Calculate metrics
    points = 0
    wins = 0
    draws = 0
    losses = 0
    goals_scored = 0
    goals_conceded = 0
    clean_sheets = 0
    failed_to_score = 0

    for fixture in recent_fixtures:
        is_home = fixture['home_team_id'] == team_id
        team_score = fixture['home_score'] if is_home else fixture['away_score']
        opponent_score = fixture['away_score'] if is_home else fixture['home_score']

        goals_scored += team_score
        goals_conceded += opponent_score

        if team_score > opponent_score:
            wins += 1
            points += 3
        elif team_score == opponent_score:
            draws += 1
            points += 1
        else:
            losses += 1

        if opponent_score == 0:
            clean_sheets += 1

        if team_score == 0:
            failed_to_score += 1

    games_played = len(recent_fixtures)

    return {
        f'form_last_{n_games}_points': float(points),
        f'form_last_{n_games}_wins': float(wins),
        f'form_last_{n_games}_draws': float(draws),
        f'form_last_{n_games}_losses': float(losses),
        f'form_last_{n_games}_goals_scored': float(goals_scored),
        f'form_last_{n_games}_goals_conceded': float(goals_conceded),
        f'form_last_{n_games}_goal_diff': float(goals_scored - goals_conceded),
        f'form_last_{n_games}_clean_sheets': float(clean_sheets),
        f'form_last_{n_games}_failed_to_score': float(failed_to_score),
        f'form_last_{n_games}_games_played': float(games_played),
        # Averages
        f'form_last_{n_games}_avg_points': float(points) / games_played,
        f'form_last_{n_games}_avg_goals_scored': float(goals_scored) / games_played,
        f'form_last_{n_games}_avg_goals_conceded': float(goals_conceded) / games_played,
    }


def calculate_home_away_form(
    team_id: int,
    current_date: datetime,
    all_fixtures: List[Dict[str, Any]],
    is_home: bool,
    n_games: int = 5
) -> Dict[str, float]:
    """
    Calculate form metrics for home or away games only

    Args:
        team_id: Team ID
        current_date: Date to calculate form before
        all_fixtures: All fixtures
        is_home: True for home form, False for away form
        n_games: Number of games to consider

    Returns:
        Home or away specific form metrics
    """
    # Filter for home or away games only
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

    if len(recent_fixtures) == 0:
        location = 'home' if is_home else 'away'
        return {
            f'{location}_form_last_{n_games}_points': 0.0,
            f'{location}_form_last_{n_games}_wins': 0.0,
            f'{location}_form_last_{n_games}_goals_scored': 0.0,
            f'{location}_form_last_{n_games}_goals_conceded': 0.0,
            f'{location}_form_last_{n_games}_games_played': 0.0,
        }

    points = 0
    wins = 0
    goals_scored = 0
    goals_conceded = 0

    for fixture in recent_fixtures:
        team_score = fixture['home_score'] if is_home else fixture['away_score']
        opponent_score = fixture['away_score'] if is_home else fixture['home_score']

        goals_scored += team_score
        goals_conceded += opponent_score

        if team_score > opponent_score:
            wins += 1
            points += 3
        elif team_score == opponent_score:
            points += 1

    games_played = len(recent_fixtures)
    location = 'home' if is_home else 'away'

    return {
        f'{location}_form_last_{n_games}_points': float(points),
        f'{location}_form_last_{n_games}_wins': float(wins),
        f'{location}_form_last_{n_games}_goals_scored': float(goals_scored),
        f'{location}_form_last_{n_games}_goals_conceded': float(goals_conceded),
        f'{location}_form_last_{n_games}_games_played': float(games_played),
        f'{location}_form_last_{n_games}_avg_points': float(points) / games_played if games_played > 0 else 0.0,
        f'{location}_form_last_{n_games}_avg_goals_scored': float(goals_scored) / games_played if games_played > 0 else 0.0,
    }


def get_form_features(
    home_team_id: int,
    away_team_id: int,
    match_date: datetime,
    all_fixtures: List[Dict[str, Any]]
) -> Dict[str, float]:
    """
    Get all form features for both teams

    Args:
        home_team_id: Home team ID
        away_team_id: Away team ID
        match_date: Match date
        all_fixtures: All fixtures

    Returns:
        Combined form features dictionary
    """
    features = {}

    # Overall form (last 5 games)
    home_form = calculate_team_form_last_n(home_team_id, match_date, all_fixtures, n_games=5)
    away_form = calculate_team_form_last_n(away_team_id, match_date, all_fixtures, n_games=5)

    # Add home_ and away_ prefixes
    for key, value in home_form.items():
        features[f'home_{key}'] = value
    for key, value in away_form.items():
        features[f'away_{key}'] = value

    # Home/Away specific form
    home_home_form = calculate_home_away_form(home_team_id, match_date, all_fixtures, is_home=True, n_games=3)
    away_away_form = calculate_home_away_form(away_team_id, match_date, all_fixtures, is_home=False, n_games=3)

    features.update(home_home_form)
    features.update(away_away_form)

    # Form differentials
    features['form_points_diff'] = features['home_form_last_5_points'] - features['away_form_last_5_points']
    features['form_goals_scored_diff'] = features['home_form_last_5_goals_scored'] - features['away_form_last_5_goals_scored']
    features['form_goal_diff_diff'] = features['home_form_last_5_goal_diff'] - features['away_form_last_5_goal_diff']

    return features
