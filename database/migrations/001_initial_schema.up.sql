-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    api_football_id INTEGER UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10),
    logo_url TEXT,
    founded INTEGER,
    venue_name VARCHAR(100),
    venue_city VARCHAR(100),
    venue_capacity INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_teams_api_football_id ON teams(api_football_id);

-- Create fixtures table
CREATE TABLE IF NOT EXISTS fixtures (
    id SERIAL PRIMARY KEY,
    api_football_id INTEGER UNIQUE NOT NULL,
    season INTEGER NOT NULL,
    round VARCHAR(50),
    match_date TIMESTAMP NOT NULL,
    home_team_id INTEGER REFERENCES teams(id),
    away_team_id INTEGER REFERENCES teams(id),
    home_score INTEGER,
    away_score INTEGER,
    status VARCHAR(20) DEFAULT 'scheduled',
    venue VARCHAR(100),
    referee VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_fixtures_match_date ON fixtures(match_date);
CREATE INDEX idx_fixtures_season ON fixtures(season);
CREATE INDEX idx_fixtures_status ON fixtures(status);
CREATE INDEX idx_fixtures_home_team ON fixtures(home_team_id);
CREATE INDEX idx_fixtures_away_team ON fixtures(away_team_id);

-- Create odds table
CREATE TABLE IF NOT EXISTS odds (
    id SERIAL PRIMARY KEY,
    fixture_id INTEGER REFERENCES fixtures(id) ON DELETE CASCADE,
    bookmaker VARCHAR(50) NOT NULL,
    market_type VARCHAR(50) NOT NULL,
    outcome VARCHAR(20),
    odds_value DECIMAL(10, 2) NOT NULL,
    recorded_at TIMESTAMP NOT NULL,
    is_closing_line BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_odds_fixture ON odds(fixture_id);
CREATE INDEX idx_odds_bookmaker ON odds(bookmaker);
CREATE INDEX idx_odds_recorded_at ON odds(recorded_at);
CREATE INDEX idx_odds_closing_line ON odds(is_closing_line) WHERE is_closing_line = TRUE;

-- Create team_stats table
CREATE TABLE IF NOT EXISTS team_stats (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id),
    season INTEGER NOT NULL,
    match_date DATE NOT NULL,
    games_played INTEGER DEFAULT 0,
    wins INTEGER DEFAULT 0,
    draws INTEGER DEFAULT 0,
    losses INTEGER DEFAULT 0,
    goals_for INTEGER DEFAULT 0,
    goals_against INTEGER DEFAULT 0,
    points INTEGER DEFAULT 0,
    position INTEGER,
    form_last_5 VARCHAR(5),
    home_wins INTEGER DEFAULT 0,
    home_draws INTEGER DEFAULT 0,
    home_losses INTEGER DEFAULT 0,
    away_wins INTEGER DEFAULT 0,
    away_draws INTEGER DEFAULT 0,
    away_losses INTEGER DEFAULT 0,
    xg_for DECIMAL(5, 2),
    xg_against DECIMAL(5, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(team_id, season, match_date)
);

CREATE INDEX idx_team_stats_team_season ON team_stats(team_id, season);
CREATE INDEX idx_team_stats_match_date ON team_stats(match_date);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Add triggers for updated_at
CREATE TRIGGER update_teams_updated_at BEFORE UPDATE ON teams
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_fixtures_updated_at BEFORE UPDATE ON fixtures
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_team_stats_updated_at BEFORE UPDATE ON team_stats
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
