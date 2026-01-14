-- Drop tables
DROP TABLE IF EXISTS model_performance;
DROP TABLE IF EXISTS bankroll;

-- Drop trigger
DROP TRIGGER IF EXISTS update_bets_updated_at ON bets;

DROP TABLE IF EXISTS bets;
DROP TABLE IF EXISTS predictions;
