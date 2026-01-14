-- Create accumulators table
CREATE TABLE IF NOT EXISTS accumulators (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100), -- e.g., "Weekend Treble #1"
    num_legs INTEGER NOT NULL,
    stake DECIMAL(10, 2) NOT NULL,
    combined_odds DECIMAL(10, 2) NOT NULL,
    combined_probability DECIMAL(5, 4) NOT NULL,
    expected_value DECIMAL(10, 4) NOT NULL,
    potential_payout DECIMAL(10, 2),
    placed_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'won', 'lost', 'void'
    actual_payout DECIMAL(10, 2),
    profit_loss DECIMAL(10, 2),
    settled_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_accumulators_status ON accumulators(status);
CREATE INDEX idx_accumulators_placed_at ON accumulators(placed_at);

-- Add trigger for accumulators updated_at
CREATE TRIGGER update_accumulators_updated_at BEFORE UPDATE ON accumulators
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create accumulator_legs table (junction table)
CREATE TABLE IF NOT EXISTS accumulator_legs (
    id SERIAL PRIMARY KEY,
    accumulator_id INTEGER REFERENCES accumulators(id) ON DELETE CASCADE,
    bet_id INTEGER REFERENCES bets(id) ON DELETE CASCADE,
    leg_order INTEGER NOT NULL, -- 1, 2, 3 for ordering
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(accumulator_id, bet_id) -- Can't add same bet twice to accumulator
);

CREATE INDEX idx_accumulator_legs_accumulator ON accumulator_legs(accumulator_id);
CREATE INDEX idx_accumulator_legs_bet ON accumulator_legs(bet_id);

-- Add view for easy accumulator querying with legs
CREATE OR REPLACE VIEW accumulator_details AS
SELECT
    a.id as accumulator_id,
    a.name,
    a.num_legs,
    a.stake,
    a.combined_odds,
    a.combined_probability,
    a.expected_value,
    a.status as accumulator_status,
    a.potential_payout,
    a.profit_loss,
    a.placed_at,
    a.settled_at,
    al.leg_order,
    b.id as bet_id,
    b.fixture_id,
    b.bet_type,
    b.odds as leg_odds,
    b.status as leg_status,
    f.match_date,
    ht.name as home_team,
    at.name as away_team
FROM accumulators a
LEFT JOIN accumulator_legs al ON a.id = al.accumulator_id
LEFT JOIN bets b ON al.bet_id = b.id
LEFT JOIN fixtures f ON b.fixture_id = f.id
LEFT JOIN teams ht ON f.home_team_id = ht.id
LEFT JOIN teams at ON f.away_team_id = at.id
ORDER BY a.id, al.leg_order;
