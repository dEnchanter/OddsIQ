-- Create predictions table
CREATE TABLE IF NOT EXISTS predictions (
    id SERIAL PRIMARY KEY,
    fixture_id INTEGER REFERENCES fixtures(id) ON DELETE CASCADE,
    model_version VARCHAR(20) NOT NULL,
    home_win_prob DECIMAL(5, 4) NOT NULL,
    draw_prob DECIMAL(5, 4) NOT NULL,
    away_win_prob DECIMAL(5, 4) NOT NULL,
    predicted_outcome VARCHAR(10),
    confidence_score DECIMAL(5, 4),
    features JSONB,
    predicted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_predictions_fixture ON predictions(fixture_id);
CREATE INDEX idx_predictions_model_version ON predictions(model_version);
CREATE INDEX idx_predictions_predicted_at ON predictions(predicted_at);

-- Create bets table
CREATE TABLE IF NOT EXISTS bets (
    id SERIAL PRIMARY KEY,
    fixture_id INTEGER REFERENCES fixtures(id),
    prediction_id INTEGER REFERENCES predictions(id),
    bet_type VARCHAR(50) NOT NULL,
    stake DECIMAL(10, 2) NOT NULL,
    odds DECIMAL(10, 2) NOT NULL,
    expected_value DECIMAL(10, 4) NOT NULL,
    bookmaker VARCHAR(50),
    placed_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    payout DECIMAL(10, 2),
    profit_loss DECIMAL(10, 2),
    settled_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bets_fixture ON bets(fixture_id);
CREATE INDEX idx_bets_status ON bets(status);
CREATE INDEX idx_bets_placed_at ON bets(placed_at);

-- Add trigger for bets updated_at
CREATE TRIGGER update_bets_updated_at BEFORE UPDATE ON bets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create bankroll table
CREATE TABLE IF NOT EXISTS bankroll (
    id SERIAL PRIMARY KEY,
    balance DECIMAL(10, 2) NOT NULL,
    total_staked DECIMAL(10, 2) DEFAULT 0,
    total_returned DECIMAL(10, 2) DEFAULT 0,
    total_profit_loss DECIMAL(10, 2) DEFAULT 0,
    roi_percentage DECIMAL(10, 4),
    num_bets INTEGER DEFAULT 0,
    num_wins INTEGER DEFAULT 0,
    num_losses INTEGER DEFAULT 0,
    win_rate DECIMAL(5, 4),
    recorded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bankroll_recorded_at ON bankroll(recorded_at);

-- Create model_performance table
CREATE TABLE IF NOT EXISTS model_performance (
    id SERIAL PRIMARY KEY,
    model_version VARCHAR(20) NOT NULL,
    evaluation_date DATE NOT NULL,
    accuracy DECIMAL(5, 4),
    precision_score DECIMAL(5, 4),
    recall DECIMAL(5, 4),
    f1_score DECIMAL(5, 4),
    brier_score DECIMAL(5, 4),
    log_loss DECIMAL(10, 6),
    roc_auc DECIMAL(5, 4),
    num_predictions INTEGER,
    backtest_roi DECIMAL(10, 4),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_model_performance_version ON model_performance(model_version);
CREATE INDEX idx_model_performance_date ON model_performance(evaluation_date);
