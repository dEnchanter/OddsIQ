package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/dEnchanter/OddsIQ/backend/internal/models"
)

// OddsRepository handles odds database operations
type OddsRepository struct {
	db *pgxpool.Pool
}

// NewOddsRepository creates a new odds repository
func NewOddsRepository(db *pgxpool.Pool) *OddsRepository {
	return &OddsRepository{db: db}
}

// Create inserts new odds
func (r *OddsRepository) Create(ctx context.Context, odds *models.Odds) error {
	query := `
		INSERT INTO odds (
			fixture_id, bookmaker, market_type, outcome, odds_value, timestamp, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		odds.FixtureID,
		odds.Bookmaker,
		odds.MarketType,
		odds.Outcome,
		odds.OddsValue,
		odds.Timestamp,
		now,
	).Scan(&odds.ID)

	if err != nil {
		return fmt.Errorf("failed to create odds: %w", err)
	}

	odds.CreatedAt = now

	return nil
}

// CreateBatch inserts multiple odds in a single transaction
func (r *OddsRepository) CreateBatch(ctx context.Context, oddsList []models.Odds) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO odds (
			fixture_id, bookmaker, market_type, outcome, odds_value, timestamp, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	for _, odds := range oddsList {
		_, err := tx.Exec(ctx, query,
			odds.FixtureID,
			odds.Bookmaker,
			odds.MarketType,
			odds.Outcome,
			odds.OddsValue,
			odds.Timestamp,
			now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert odds: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByFixture retrieves all odds for a specific fixture
func (r *OddsRepository) GetByFixture(ctx context.Context, fixtureID int) ([]models.Odds, error) {
	query := `
		SELECT id, fixture_id, bookmaker, market_type, outcome, odds_value, timestamp, created_at
		FROM odds
		WHERE fixture_id = $1
		ORDER BY timestamp DESC, bookmaker, market_type, outcome
	`

	rows, err := r.db.Query(ctx, query, fixtureID)
	if err != nil {
		return nil, fmt.Errorf("failed to query odds: %w", err)
	}
	defer rows.Close()

	return r.scanOdds(rows)
}

// GetLatestByFixture retrieves the latest odds for each market/outcome combination for a fixture
func (r *OddsRepository) GetLatestByFixture(ctx context.Context, fixtureID int) ([]models.Odds, error) {
	query := `
		SELECT DISTINCT ON (bookmaker, market_type, outcome)
			id, fixture_id, bookmaker, market_type, outcome, odds_value, timestamp, created_at
		FROM odds
		WHERE fixture_id = $1
		ORDER BY bookmaker, market_type, outcome, timestamp DESC
	`

	rows, err := r.db.Query(ctx, query, fixtureID)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest odds: %w", err)
	}
	defer rows.Close()

	return r.scanOdds(rows)
}

// GetByFixtureAndMarket retrieves odds for a specific fixture and market type
func (r *OddsRepository) GetByFixtureAndMarket(ctx context.Context, fixtureID int, marketType string) ([]models.Odds, error) {
	query := `
		SELECT id, fixture_id, bookmaker, market_type, outcome, odds_value, timestamp, created_at
		FROM odds
		WHERE fixture_id = $1 AND market_type = $2
		ORDER BY timestamp DESC, bookmaker, outcome
	`

	rows, err := r.db.Query(ctx, query, fixtureID, marketType)
	if err != nil {
		return nil, fmt.Errorf("failed to query odds by market: %w", err)
	}
	defer rows.Close()

	return r.scanOdds(rows)
}

// GetLatestByFixtureAndMarket retrieves the latest odds for a specific fixture and market
func (r *OddsRepository) GetLatestByFixtureAndMarket(ctx context.Context, fixtureID int, marketType string) ([]models.Odds, error) {
	query := `
		SELECT DISTINCT ON (bookmaker, outcome)
			id, fixture_id, bookmaker, market_type, outcome, odds_value, timestamp, created_at
		FROM odds
		WHERE fixture_id = $1 AND market_type = $2
		ORDER BY bookmaker, outcome, timestamp DESC
	`

	rows, err := r.db.Query(ctx, query, fixtureID, marketType)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest odds by market: %w", err)
	}
	defer rows.Close()

	return r.scanOdds(rows)
}

// GetBestOdds retrieves the best (highest) odds for a specific fixture, market, and outcome
func (r *OddsRepository) GetBestOdds(ctx context.Context, fixtureID int, marketType, outcome string) (*models.Odds, error) {
	query := `
		SELECT id, fixture_id, bookmaker, market_type, outcome, odds_value, timestamp, created_at
		FROM odds
		WHERE fixture_id = $1 AND market_type = $2 AND outcome = $3
		ORDER BY odds_value DESC, timestamp DESC
		LIMIT 1
	`

	odds := &models.Odds{}
	err := r.db.QueryRow(ctx, query, fixtureID, marketType, outcome).Scan(
		&odds.ID,
		&odds.FixtureID,
		&odds.Bookmaker,
		&odds.MarketType,
		&odds.Outcome,
		&odds.OddsValue,
		&odds.Timestamp,
		&odds.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("odds not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get best odds: %w", err)
	}

	return odds, nil
}

// GetByBookmaker retrieves all odds from a specific bookmaker
func (r *OddsRepository) GetByBookmaker(ctx context.Context, bookmaker string) ([]models.Odds, error) {
	query := `
		SELECT id, fixture_id, bookmaker, market_type, outcome, odds_value, timestamp, created_at
		FROM odds
		WHERE bookmaker = $1
		ORDER BY timestamp DESC
		LIMIT 1000
	`

	rows, err := r.db.Query(ctx, query, bookmaker)
	if err != nil {
		return nil, fmt.Errorf("failed to query odds by bookmaker: %w", err)
	}
	defer rows.Close()

	return r.scanOdds(rows)
}

// GetByDateRange retrieves odds within a date range
func (r *OddsRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]models.Odds, error) {
	query := `
		SELECT id, fixture_id, bookmaker, market_type, outcome, odds_value, timestamp, created_at
		FROM odds
		WHERE timestamp >= $1 AND timestamp <= $2
		ORDER BY timestamp DESC
		LIMIT 5000
	`

	rows, err := r.db.Query(ctx, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query odds by date range: %w", err)
	}
	defer rows.Close()

	return r.scanOdds(rows)
}

// DeleteOldOdds deletes odds older than a specific date
func (r *OddsRepository) DeleteOldOdds(ctx context.Context, before time.Time) (int64, error) {
	query := `DELETE FROM odds WHERE timestamp < $1`

	result, err := r.db.Exec(ctx, query, before)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old odds: %w", err)
	}

	return result.RowsAffected(), nil
}

// GetMarketTypes retrieves all distinct market types
func (r *OddsRepository) GetMarketTypes(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT market_type FROM odds ORDER BY market_type`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query market types: %w", err)
	}
	defer rows.Close()

	var marketTypes []string
	for rows.Next() {
		var marketType string
		if err := rows.Scan(&marketType); err != nil {
			return nil, fmt.Errorf("failed to scan market type: %w", err)
		}
		marketTypes = append(marketTypes, marketType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return marketTypes, nil
}

// GetBookmakers retrieves all distinct bookmakers
func (r *OddsRepository) GetBookmakers(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT bookmaker FROM odds ORDER BY bookmaker`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookmakers: %w", err)
	}
	defer rows.Close()

	var bookmakers []string
	for rows.Next() {
		var bookmaker string
		if err := rows.Scan(&bookmaker); err != nil {
			return nil, fmt.Errorf("failed to scan bookmaker: %w", err)
		}
		bookmakers = append(bookmakers, bookmaker)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return bookmakers, nil
}

// GetAverageOdds calculates average odds for a specific fixture, market, and outcome
func (r *OddsRepository) GetAverageOdds(ctx context.Context, fixtureID int, marketType, outcome string) (float64, error) {
	query := `
		SELECT AVG(odds_value)
		FROM (
			SELECT DISTINCT ON (bookmaker) odds_value
			FROM odds
			WHERE fixture_id = $1 AND market_type = $2 AND outcome = $3
			ORDER BY bookmaker, timestamp DESC
		) latest_odds
	`

	var avgOdds float64
	err := r.db.QueryRow(ctx, query, fixtureID, marketType, outcome).Scan(&avgOdds)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate average odds: %w", err)
	}

	return avgOdds, nil
}

// UpsertOdds inserts or updates odds (based on fixture, bookmaker, market, outcome, and recent timestamp)
func (r *OddsRepository) UpsertOdds(ctx context.Context, odds *models.Odds) error {
	// For odds, we generally want to insert new records to track changes over time
	// rather than update existing ones, so we just insert
	return r.Create(ctx, odds)
}

// Helper function to scan odds from rows
func (r *OddsRepository) scanOdds(rows pgx.Rows) ([]models.Odds, error) {
	var oddsList []models.Odds
	for rows.Next() {
		var odds models.Odds
		err := rows.Scan(
			&odds.ID,
			&odds.FixtureID,
			&odds.Bookmaker,
			&odds.MarketType,
			&odds.Outcome,
			&odds.OddsValue,
			&odds.Timestamp,
			&odds.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan odds: %w", err)
		}
		oddsList = append(oddsList, odds)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return oddsList, nil
}
