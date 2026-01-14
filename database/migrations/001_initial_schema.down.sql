-- Drop triggers
DROP TRIGGER IF EXISTS update_team_stats_updated_at ON team_stats;
DROP TRIGGER IF EXISTS update_fixtures_updated_at ON fixtures;
DROP TRIGGER IF EXISTS update_teams_updated_at ON teams;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order
DROP TABLE IF EXISTS team_stats;
DROP TABLE IF EXISTS odds;
DROP TABLE IF EXISTS fixtures;
DROP TABLE IF EXISTS teams;
