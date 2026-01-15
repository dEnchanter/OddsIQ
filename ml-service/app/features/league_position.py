"""
League Position Feature Engineering

Calculates league position and points-based features
"""
from typing import List, Dict, Any
from datetime import datetime
from collections import defaultdict


def calculate_league_table(
    season: int,
    up_to_date: datetime,
    all_fixtures: List[Dict[str, Any]]
) -> Dict[int, Dict[str, Any]]:
    """
    Calculate league table standings up to a specific date

    Args:
        season: Season year
        up_to_date: Calculate table up to this date (exclusive)
        all_fixtures: All fixtures

    Returns:
        Dictionary mapping team_id to standings dict
    """
    # Filter fixtures for this season up to date
    season_fixtures = [
        f for f in all_fixtures
        if f['season'] == season
        and f['match_date'] < up_to_date
        and f['status'] == 'FT'
        and f['home_score'] is not None
    ]

    # Initialize team stats
    table = defaultdict(lambda: {
        'played': 0,
        'wins': 0,
        'draws': 0,
        'losses': 0,
        'goals_for': 0,
        'goals_against': 0,
        'goal_diff': 0,
        'points': 0,
    })

    # Calculate standings
    for fixture in season_fixtures:
        home_id = fixture['home_team_id']
        away_id = fixture['away_team_id']
        home_score = fixture['home_score']
        away_score = fixture['away_score']

        # Home team
        table[home_id]['played'] += 1
        table[home_id]['goals_for'] += home_score
        table[home_id]['goals_against'] += away_score
        table[home_id]['goal_diff'] += (home_score - away_score)

        # Away team
        table[away_id]['played'] += 1
        table[away_id]['goals_for'] += away_score
        table[away_id]['goals_against'] += home_score
        table[away_id]['goal_diff'] += (away_score - home_score)

        # Points
        if home_score > away_score:
            table[home_id]['wins'] += 1
            table[home_id]['points'] += 3
            table[away_id]['losses'] += 1
        elif home_score < away_score:
            table[away_id]['wins'] += 1
            table[away_id]['points'] += 3
            table[home_id]['losses'] += 1
        else:
            table[home_id]['draws'] += 1
            table[home_id]['points'] += 1
            table[away_id]['draws'] += 1
            table[away_id]['points'] += 1

    # Calculate positions
    teams_sorted = sorted(
        table.items(),
        key=lambda x: (x[1]['points'], x[1]['goal_diff'], x[1]['goals_for']),
        reverse=True
    )

    # Add position to each team
    position = 1
    for team_id, stats in teams_sorted:
        table[team_id]['position'] = position
        position += 1

    return dict(table)


def get_league_position_features(
    home_team_id: int,
    away_team_id: int,
    season: int,
    match_date: datetime,
    all_fixtures: List[Dict[str, Any]]
) -> Dict[str, float]:
    """
    Get league position features for both teams

    Args:
        home_team_id: Home team ID
        away_team_id: Away team ID
        season: Season year
        match_date: Match date
        all_fixtures: All fixtures

    Returns:
        League position features
    """
    # Calculate league table
    table = calculate_league_table(season, match_date, all_fixtures)

    # Get team stats
    home_stats = table.get(home_team_id, {
        'position': 10,  # Mid-table default
        'points': 0,
        'played': 0,
        'goals_for': 0,
        'goals_against': 0,
        'goal_diff': 0,
        'wins': 0,
    })

    away_stats = table.get(away_team_id, {
        'position': 10,
        'points': 0,
        'played': 0,
        'goals_for': 0,
        'goals_against': 0,
        'goal_diff': 0,
        'wins': 0,
    })

    # Calculate features
    home_played = home_stats['played'] if home_stats['played'] > 0 else 1
    away_played = away_stats['played'] if away_stats['played'] > 0 else 1

    features = {
        # Positions
        'home_position': float(home_stats['position']),
        'away_position': float(away_stats['position']),
        'position_diff': float(home_stats['position'] - away_stats['position']),

        # Points
        'home_points': float(home_stats['points']),
        'away_points': float(away_stats['points']),
        'points_diff': float(home_stats['points'] - away_stats['points']),

        # Points per game
        'home_ppg': float(home_stats['points']) / home_played,
        'away_ppg': float(away_stats['points']) / away_played,
        'ppg_diff': (float(home_stats['points']) / home_played) - (float(away_stats['points']) / away_played),

        # Goals
        'home_season_goals_for': float(home_stats['goals_for']),
        'away_season_goals_for': float(away_stats['goals_for']),
        'home_season_goals_against': float(home_stats['goals_against']),
        'away_season_goals_against': float(away_stats['goals_against']),

        # Goal averages
        'home_avg_goals_scored': float(home_stats['goals_for']) / home_played,
        'away_avg_goals_scored': float(away_stats['goals_for']) / away_played,
        'home_avg_goals_conceded': float(home_stats['goals_against']) / home_played,
        'away_avg_goals_conceded': float(away_stats['goals_against']) / away_played,

        # Goal difference
        'home_goal_diff': float(home_stats['goal_diff']),
        'away_goal_diff': float(away_stats['goal_diff']),
        'goal_diff_diff': float(home_stats['goal_diff'] - away_stats['goal_diff']),

        # Games played
        'home_games_played': float(home_stats['played']),
        'away_games_played': float(away_stats['played']),

        # Win percentage
        'home_win_pct': float(home_stats['wins']) / home_played,
        'away_win_pct': float(away_stats['wins']) / away_played,
    }

    return features
