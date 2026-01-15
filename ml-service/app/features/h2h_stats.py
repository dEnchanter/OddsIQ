"""
Head-to-Head (H2H) Statistics Feature Engineering

Calculates historical head-to-head performance between two teams
"""
from typing import List, Dict, Any
from datetime import datetime


def calculate_h2h_stats(
    home_team_id: int,
    away_team_id: int,
    current_date: datetime,
    all_fixtures: List[Dict[str, Any]],
    n_games: int = 5
) -> Dict[str, float]:
    """
    Calculate head-to-head statistics between two teams

    Args:
        home_team_id: Home team ID
        away_team_id: Away team ID
        current_date: Current match date (to exclude)
        all_fixtures: All historical fixtures
        n_games: Number of recent H2H games to consider

    Returns:
        H2H feature dictionary
    """
    # Find all previous H2H matches
    h2h_fixtures = [
        f for f in all_fixtures
        if ((f['home_team_id'] == home_team_id and f['away_team_id'] == away_team_id) or
            (f['home_team_id'] == away_team_id and f['away_team_id'] == home_team_id))
        and f['match_date'] < current_date
        and f['status'] == 'Match Finished'
        and f['home_score'] is not None
    ]

    # Sort by date descending
    h2h_fixtures.sort(key=lambda x: x['match_date'], reverse=True)
    recent_h2h = h2h_fixtures[:n_games]

    if len(recent_h2h) == 0:
        # No H2H history
        return {
            'h2h_games_played': 0.0,
            'h2h_home_wins': 0.0,
            'h2h_away_wins': 0.0,
            'h2h_draws': 0.0,
            'h2h_home_goals_scored': 0.0,
            'h2h_away_goals_scored': 0.0,
            'h2h_home_win_pct': 0.0,
            'h2h_avg_total_goals': 0.0,
            'h2h_home_as_home_wins': 0.0,
            'h2h_home_as_home_games': 0.0,
        }

    # Calculate H2H statistics
    home_wins = 0
    away_wins = 0
    draws = 0
    home_goals_total = 0
    away_goals_total = 0
    home_as_home_wins = 0
    home_as_home_games = 0

    for fixture in recent_h2h:
        # Check if current home team was home in this fixture
        if fixture['home_team_id'] == home_team_id:
            # Home team was home
            home_as_home_games += 1
            home_goals = fixture['home_score']
            away_goals = fixture['away_score']

            if home_goals > away_goals:
                home_wins += 1
                home_as_home_wins += 1
            elif home_goals < away_goals:
                away_wins += 1
            else:
                draws += 1

            home_goals_total += home_goals
            away_goals_total += away_goals
        else:
            # Home team was away
            home_goals = fixture['away_score']
            away_goals = fixture['home_score']

            if home_goals > away_goals:
                home_wins += 1
            elif home_goals < away_goals:
                away_wins += 1
            else:
                draws += 1

            home_goals_total += home_goals
            away_goals_total += away_goals

    games_played = len(recent_h2h)
    total_goals = home_goals_total + away_goals_total

    return {
        'h2h_games_played': float(games_played),
        'h2h_home_wins': float(home_wins),
        'h2h_away_wins': float(away_wins),
        'h2h_draws': float(draws),
        'h2h_home_goals_scored': float(home_goals_total),
        'h2h_away_goals_scored': float(away_goals_total),
        'h2h_goal_diff': float(home_goals_total - away_goals_total),
        'h2h_home_win_pct': float(home_wins) / games_played,
        'h2h_away_win_pct': float(away_wins) / games_played,
        'h2h_draw_pct': float(draws) / games_played,
        'h2h_avg_total_goals': float(total_goals) / games_played,
        'h2h_avg_home_goals': float(home_goals_total) / games_played,
        'h2h_avg_away_goals': float(away_goals_total) / games_played,
        'h2h_home_as_home_wins': float(home_as_home_wins),
        'h2h_home_as_home_games': float(home_as_home_games),
        'h2h_home_as_home_win_pct': float(home_as_home_wins) / home_as_home_games if home_as_home_games > 0 else 0.0,
    }


def get_h2h_features(
    home_team_id: int,
    away_team_id: int,
    match_date: datetime,
    all_fixtures: List[Dict[str, Any]]
) -> Dict[str, float]:
    """
    Get all H2H features

    Args:
        home_team_id: Home team ID
        away_team_id: Away team ID
        match_date: Match date
        all_fixtures: All fixtures

    Returns:
        H2H features dictionary
    """
    return calculate_h2h_stats(home_team_id, away_team_id, match_date, all_fixtures, n_games=5)
