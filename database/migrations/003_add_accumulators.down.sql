-- Drop view
DROP VIEW IF EXISTS accumulator_details;

-- Drop tables
DROP TABLE IF EXISTS accumulator_legs;

-- Drop trigger
DROP TRIGGER IF EXISTS update_accumulators_updated_at ON accumulators;

DROP TABLE IF EXISTS accumulators;
