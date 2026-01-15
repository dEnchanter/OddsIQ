"""
Database connection module for ML service
"""
import psycopg
from psycopg.rows import dict_row
from typing import List, Dict, Any, Optional
from config.config import config


class Database:
    """Database connection manager"""

    def __init__(self):
        self.connection_string = config.DATABASE_URL
        self._connection = None

    def get_connection(self):
        """Get database connection"""
        if self._connection is None or self._connection.closed:
            self._connection = psycopg.connect(
                self.connection_string,
                row_factory=dict_row
            )
        return self._connection

    def close(self):
        """Close database connection"""
        if self._connection and not self._connection.closed:
            self._connection.close()

    def execute_query(self, query: str, params: tuple = None) -> List[Dict[str, Any]]:
        """Execute a SELECT query and return results"""
        conn = self.get_connection()
        cursor = conn.cursor()
        try:
            cursor.execute(query, params)
            results = cursor.fetchall()
            return results  # Already dicts with dict_row factory
        finally:
            cursor.close()

    def execute_one(self, query: str, params: tuple = None) -> Optional[Dict[str, Any]]:
        """Execute a SELECT query and return single result"""
        conn = self.get_connection()
        cursor = conn.cursor()
        try:
            cursor.execute(query, params)
            result = cursor.fetchone()
            return result if result else None  # Already dict with dict_row factory
        finally:
            cursor.close()


# Global database instance
db = Database()


def get_all_fixtures(season: Optional[int] = None) -> List[Dict[str, Any]]:
    """
    Get all fixtures from database

    Args:
        season: Optional season filter

    Returns:
        List of fixture dictionaries
    """
    if season:
        query = """
            SELECT
                f.id,
                f.api_football_id,
                f.season,
                f.round,
                f.match_date,
                f.home_team_id,
                f.away_team_id,
                ht.name as home_team_name,
                at.name as away_team_name,
                f.home_score,
                f.away_score,
                f.status,
                f.venue_name,
                f.referee
            FROM fixtures f
            JOIN teams ht ON f.home_team_id = ht.id
            JOIN teams at ON f.away_team_id = at.id
            WHERE f.season = %s
            ORDER BY f.match_date
        """
        return db.execute_query(query, (season,))
    else:
        query = """
            SELECT
                f.id,
                f.api_football_id,
                f.season,
                f.round,
                f.match_date,
                f.home_team_id,
                f.away_team_id,
                ht.name as home_team_name,
                at.name as away_team_name,
                f.home_score,
                f.away_score,
                f.status,
                f.venue_name,
                f.referee
            FROM fixtures f
            JOIN teams ht ON f.home_team_id = ht.id
            JOIN teams at ON f.away_team_id = at.id
            ORDER BY f.season, f.match_date
        """
        return db.execute_query(query)


def get_fixtures_for_training(seasons: List[int]) -> List[Dict[str, Any]]:
    """
    Get fixtures for training (completed matches only)

    Args:
        seasons: List of seasons to include

    Returns:
        List of completed fixture dictionaries
    """
    query = """
        SELECT
            f.id,
            f.api_football_id,
            f.season,
            f.round,
            f.match_date,
            f.home_team_id,
            f.away_team_id,
            ht.name as home_team_name,
            at.name as away_team_name,
            f.home_score,
            f.away_score,
            f.status
        FROM fixtures f
        JOIN teams ht ON f.home_team_id = ht.id
        JOIN teams at ON f.away_team_id = at.id
        WHERE f.season = ANY(%s)
          AND f.status = 'FT'
          AND f.home_score IS NOT NULL
          AND f.away_score IS NOT NULL
        ORDER BY f.match_date
    """
    return db.execute_query(query, (seasons,))


def get_team_info(team_id: int) -> Optional[Dict[str, Any]]:
    """Get team information"""
    query = """
        SELECT
            id,
            api_football_id,
            name,
            code,
            venue_name,
            venue_city
        FROM teams
        WHERE id = %s
    """
    return db.execute_one(query, (team_id,))


def get_all_teams() -> List[Dict[str, Any]]:
    """Get all teams"""
    query = """
        SELECT
            id,
            api_football_id,
            name,
            code,
            venue_name,
            venue_city
        FROM teams
        ORDER BY name
    """
    return db.execute_query(query)


def test_connection():
    """Test database connection"""
    try:
        fixtures = get_all_fixtures()
        teams = get_all_teams()
        print(f"[OK] Database connection successful!")
        print(f"   Found {len(fixtures)} fixtures")
        print(f"   Found {len(teams)} teams")
        return True
    except Exception as e:
        print(f"[ERROR] Database connection failed: {e}")
        return False


if __name__ == "__main__":
    # Test the connection
    test_connection()
