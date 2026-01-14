package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/dEnchanter/OddsIQ/backend/internal/models"
)

// FixturesRepository handles fixture database operations
type FixturesRepository struct {
	db *pgxpool.Pool
}

// NewFixturesRepository creates a new fixtures repository
func NewFixturesRepository(db *pgxpool.Pool) *FixturesRepository {
	return &FixturesRepository{db: db}
}

// Create inserts a new fixture
func (r *FixturesRepository) Create(ctx context.Context, fixture *models.Fixture) error {
	query := `
		INSERT INTO fixtures (
			api_football_id, season, match_date, round, home_team_id, away_team_id,
			status, home_score, away_score, venue_name, referee, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		fixture.APIFootballID,
		fixture.Season,
		fixture.MatchDate,
		fixture.Round,
		fixture.HomeTeamID,
		fixture.AwayTeamID,
		fixture.Status,
		fixture.HomeScore,
		fixture.AwayScore,
		fixture.VenueName,
		fixture.Referee,
		now,
		now,
	).Scan(&fixture.ID)

	if err != nil {
		return fmt.Errorf("failed to create fixture: %w", err)
	}

	fixture.CreatedAt = now
	fixture.UpdatedAt = now

	return nil
}

// GetByID retrieves a fixture by ID
func (r *FixturesRepository) GetByID(ctx context.Context, id int) (*models.Fixture, error) {
	query := `
		SELECT id, api_football_id, season, match_date, round, home_team_id, away_team_id,
			status, home_score, away_score, venue_name, referee, created_at, updated_at
		FROM fixtures
		WHERE id = $1
	`

	fixture := &models.Fixture{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&fixture.ID,
		&fixture.APIFootballID,
		&fixture.Season,
		&fixture.MatchDate,
		&fixture.Round,
		&fixture.HomeTeamID,
		&fixture.AwayTeamID,
		&fixture.Status,
		&fixture.HomeScore,
		&fixture.AwayScore,
		&fixture.VenueName,
		&fixture.Referee,
		&fixture.CreatedAt,
		&fixture.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("fixture not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get fixture: %w", err)
	}

	return fixture, nil
}

// GetByAPIFootballID retrieves a fixture by API-Football ID
func (r *FixturesRepository) GetByAPIFootballID(ctx context.Context, apiFootballID int) (*models.Fixture, error) {
	query := `
		SELECT id, api_football_id, season, match_date, round, home_team_id, away_team_id,
			status, home_score, away_score, venue_name, referee, created_at, updated_at
		FROM fixtures
		WHERE api_football_id = $1
	`

	fixture := &models.Fixture{}
	err := r.db.QueryRow(ctx, query, apiFootballID).Scan(
		&fixture.ID,
		&fixture.APIFootballID,
		&fixture.Season,
		&fixture.MatchDate,
		&fixture.Round,
		&fixture.HomeTeamID,
		&fixture.AwayTeamID,
		&fixture.Status,
		&fixture.HomeScore,
		&fixture.AwayScore,
		&fixture.VenueName,
		&fixture.Referee,
		&fixture.CreatedAt,
		&fixture.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("fixture not found with api_football_id %d", apiFootballID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get fixture: %w", err)
	}

	return fixture, nil
}

// GetBySeason retrieves all fixtures for a specific season
func (r *FixturesRepository) GetBySeason(ctx context.Context, season int) ([]models.Fixture, error) {
	query := `
		SELECT id, api_football_id, season, match_date, round, home_team_id, away_team_id,
			status, home_score, away_score, venue_name, referee, created_at, updated_at
		FROM fixtures
		WHERE season = $1
		ORDER BY match_date
	`

	rows, err := r.db.Query(ctx, query, season)
	if err != nil {
		return nil, fmt.Errorf("failed to query fixtures: %w", err)
	}
	defer rows.Close()

	return r.scanFixtures(rows)
}

// GetByDateRange retrieves fixtures within a date range
func (r *FixturesRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]models.Fixture, error) {
	query := `
		SELECT id, api_football_id, season, match_date, round, home_team_id, away_team_id,
			status, home_score, away_score, venue_name, referee, created_at, updated_at
		FROM fixtures
		WHERE match_date >= $1 AND match_date <= $2
		ORDER BY match_date
	`

	rows, err := r.db.Query(ctx, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query fixtures: %w", err)
	}
	defer rows.Close()

	return r.scanFixtures(rows)
}

// GetUpcoming retrieves upcoming fixtures (not yet played)
func (r *FixturesRepository) GetUpcoming(ctx context.Context, limit int) ([]models.Fixture, error) {
	query := `
		SELECT id, api_football_id, season, match_date, round, home_team_id, away_team_id,
			status, home_score, away_score, venue_name, referee, created_at, updated_at
		FROM fixtures
		WHERE status = 'NS' AND match_date > NOW()
		ORDER BY match_date
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query upcoming fixtures: %w", err)
	}
	defer rows.Close()

	return r.scanFixtures(rows)
}

// GetByStatus retrieves fixtures by status
func (r *FixturesRepository) GetByStatus(ctx context.Context, status string) ([]models.Fixture, error) {
	query := `
		SELECT id, api_football_id, season, match_date, round, home_team_id, away_team_id,
			status, home_score, away_score, venue_name, referee, created_at, updated_at
		FROM fixtures
		WHERE status = $1
		ORDER BY match_date DESC
	`

	rows, err := r.db.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query fixtures by status: %w", err)
	}
	defer rows.Close()

	return r.scanFixtures(rows)
}

// GetByTeam retrieves all fixtures for a specific team
func (r *FixturesRepository) GetByTeam(ctx context.Context, teamID int) ([]models.Fixture, error) {
	query := `
		SELECT id, api_football_id, season, match_date, round, home_team_id, away_team_id,
			status, home_score, away_score, venue_name, referee, created_at, updated_at
		FROM fixtures
		WHERE home_team_id = $1 OR away_team_id = $1
		ORDER BY match_date DESC
	`

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query fixtures by team: %w", err)
	}
	defer rows.Close()

	return r.scanFixtures(rows)
}

// GetRecentByTeam retrieves recent fixtures for a team
func (r *FixturesRepository) GetRecentByTeam(ctx context.Context, teamID int, limit int) ([]models.Fixture, error) {
	query := `
		SELECT id, api_football_id, season, match_date, round, home_team_id, away_team_id,
			status, home_score, away_score, venue_name, referee, created_at, updated_at
		FROM fixtures
		WHERE (home_team_id = $1 OR away_team_id = $1) AND status = 'FT'
		ORDER BY match_date DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, teamID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent fixtures: %w", err)
	}
	defer rows.Close()

	return r.scanFixtures(rows)
}

// Update updates an existing fixture
func (r *FixturesRepository) Update(ctx context.Context, fixture *models.Fixture) error {
	query := `
		UPDATE fixtures
		SET season = $1, match_date = $2, round = $3, home_team_id = $4, away_team_id = $5,
			status = $6, home_score = $7, away_score = $8, venue_name = $9, referee = $10, updated_at = $11
		WHERE id = $12
	`

	now := time.Now()
	result, err := r.db.Exec(ctx, query,
		fixture.Season,
		fixture.MatchDate,
		fixture.Round,
		fixture.HomeTeamID,
		fixture.AwayTeamID,
		fixture.Status,
		fixture.HomeScore,
		fixture.AwayScore,
		fixture.VenueName,
		fixture.Referee,
		now,
		fixture.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update fixture: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("fixture not found with id %d", fixture.ID)
	}

	fixture.UpdatedAt = now

	return nil
}

// UpdateScore updates fixture score and status
func (r *FixturesRepository) UpdateScore(ctx context.Context, id int, homeScore, awayScore *int, status string) error {
	query := `
		UPDATE fixtures
		SET home_score = $1, away_score = $2, status = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.Exec(ctx, query, homeScore, awayScore, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update fixture score: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("fixture not found with id %d", id)
	}

	return nil
}

// Upsert inserts or updates a fixture based on API-Football ID
func (r *FixturesRepository) Upsert(ctx context.Context, fixture *models.Fixture) error {
	query := `
		INSERT INTO fixtures (
			api_football_id, season, match_date, round, home_team_id, away_team_id,
			status, home_score, away_score, venue_name, referee, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (api_football_id)
		DO UPDATE SET
			season = EXCLUDED.season,
			match_date = EXCLUDED.match_date,
			round = EXCLUDED.round,
			home_team_id = EXCLUDED.home_team_id,
			away_team_id = EXCLUDED.away_team_id,
			status = EXCLUDED.status,
			home_score = EXCLUDED.home_score,
			away_score = EXCLUDED.away_score,
			venue_name = EXCLUDED.venue_name,
			referee = EXCLUDED.referee,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		fixture.APIFootballID,
		fixture.Season,
		fixture.MatchDate,
		fixture.Round,
		fixture.HomeTeamID,
		fixture.AwayTeamID,
		fixture.Status,
		fixture.HomeScore,
		fixture.AwayScore,
		fixture.VenueName,
		fixture.Referee,
		now,
		now,
	).Scan(&fixture.ID)

	if err != nil {
		return fmt.Errorf("failed to upsert fixture: %w", err)
	}

	fixture.UpdatedAt = now

	return nil
}

// Delete deletes a fixture
func (r *FixturesRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM fixtures WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete fixture: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("fixture not found with id %d", id)
	}

	return nil
}

// Helper function to scan fixtures from rows
func (r *FixturesRepository) scanFixtures(rows pgx.Rows) ([]models.Fixture, error) {
	var fixtures []models.Fixture
	for rows.Next() {
		var fixture models.Fixture
		err := rows.Scan(
			&fixture.ID,
			&fixture.APIFootballID,
			&fixture.Season,
			&fixture.MatchDate,
			&fixture.Round,
			&fixture.HomeTeamID,
			&fixture.AwayTeamID,
			&fixture.Status,
			&fixture.HomeScore,
			&fixture.AwayScore,
			&fixture.VenueName,
			&fixture.Referee,
			&fixture.CreatedAt,
			&fixture.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fixture: %w", err)
		}
		fixtures = append(fixtures, fixture)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return fixtures, nil
}
