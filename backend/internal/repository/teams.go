package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/dEnchanter/OddsIQ/backend/internal/models"
)

// TeamsRepository handles team database operations
type TeamsRepository struct {
	db *pgxpool.Pool
}

// NewTeamsRepository creates a new teams repository
func NewTeamsRepository(db *pgxpool.Pool) *TeamsRepository {
	return &TeamsRepository{db: db}
}

// Create inserts a new team
func (r *TeamsRepository) Create(ctx context.Context, team *models.Team) error {
	query := `
		INSERT INTO teams (api_football_id, name, code, logo_url, founded, venue_name, venue_city, venue_capacity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		team.APIFootballID,
		team.Name,
		team.Code,
		team.LogoURL,
		team.Founded,
		team.VenueName,
		team.VenueCity,
		team.VenueCapacity,
		now,
		now,
	).Scan(&team.ID)

	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	team.CreatedAt = now
	team.UpdatedAt = now

	return nil
}

// GetByID retrieves a team by ID
func (r *TeamsRepository) GetByID(ctx context.Context, id int) (*models.Team, error) {
	query := `
		SELECT id, api_football_id, name, code, logo_url, founded, venue_name, venue_city, venue_capacity, created_at, updated_at
		FROM teams
		WHERE id = $1
	`

	team := &models.Team{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&team.ID,
		&team.APIFootballID,
		&team.Name,
		&team.Code,
		&team.LogoURL,
		&team.Founded,
		&team.VenueName,
		&team.VenueCity,
		&team.VenueCapacity,
		&team.CreatedAt,
		&team.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("team not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return team, nil
}

// GetByAPIFootballID retrieves a team by API-Football ID
func (r *TeamsRepository) GetByAPIFootballID(ctx context.Context, apiFootballID int) (*models.Team, error) {
	query := `
		SELECT id, api_football_id, name, code, logo_url, founded, venue_name, venue_city, venue_capacity, created_at, updated_at
		FROM teams
		WHERE api_football_id = $1
	`

	team := &models.Team{}
	err := r.db.QueryRow(ctx, query, apiFootballID).Scan(
		&team.ID,
		&team.APIFootballID,
		&team.Name,
		&team.Code,
		&team.LogoURL,
		&team.Founded,
		&team.VenueName,
		&team.VenueCity,
		&team.VenueCapacity,
		&team.CreatedAt,
		&team.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("team not found with api_football_id %d", apiFootballID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return team, nil
}

// GetAll retrieves all teams
func (r *TeamsRepository) GetAll(ctx context.Context) ([]models.Team, error) {
	query := `
		SELECT id, api_football_id, name, code, logo_url, founded, venue_name, venue_city, venue_capacity, created_at, updated_at
		FROM teams
		ORDER BY name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query teams: %w", err)
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		err := rows.Scan(
			&team.ID,
			&team.APIFootballID,
			&team.Name,
			&team.Code,
			&team.LogoURL,
			&team.Founded,
			&team.VenueName,
			&team.VenueCity,
			&team.VenueCapacity,
			&team.CreatedAt,
			&team.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team: %w", err)
		}
		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return teams, nil
}

// Update updates an existing team
func (r *TeamsRepository) Update(ctx context.Context, team *models.Team) error {
	query := `
		UPDATE teams
		SET name = $1, code = $2, logo_url = $3, founded = $4, venue_name = $5, venue_city = $6, venue_capacity = $7, updated_at = $8
		WHERE id = $9
	`

	now := time.Now()
	result, err := r.db.Exec(ctx, query,
		team.Name,
		team.Code,
		team.LogoURL,
		team.Founded,
		team.VenueName,
		team.VenueCity,
		team.VenueCapacity,
		now,
		team.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("team not found with id %d", team.ID)
	}

	team.UpdatedAt = now

	return nil
}

// Delete deletes a team
func (r *TeamsRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM teams WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("team not found with id %d", id)
	}

	return nil
}

// Upsert inserts or updates a team based on API-Football ID
func (r *TeamsRepository) Upsert(ctx context.Context, team *models.Team) error {
	query := `
		INSERT INTO teams (api_football_id, name, code, logo_url, founded, venue_name, venue_city, venue_capacity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (api_football_id)
		DO UPDATE SET
			name = EXCLUDED.name,
			code = EXCLUDED.code,
			logo_url = EXCLUDED.logo_url,
			founded = EXCLUDED.founded,
			venue_name = EXCLUDED.venue_name,
			venue_city = EXCLUDED.venue_city,
			venue_capacity = EXCLUDED.venue_capacity,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		team.APIFootballID,
		team.Name,
		team.Code,
		team.LogoURL,
		team.Founded,
		team.VenueName,
		team.VenueCity,
		team.VenueCapacity,
		now,
		now,
	).Scan(&team.ID)

	if err != nil {
		return fmt.Errorf("failed to upsert team: %w", err)
	}

	team.UpdatedAt = now

	return nil
}

// GetPremierLeagueTeams retrieves all Premier League teams (convenience method)
func (r *TeamsRepository) GetPremierLeagueTeams(ctx context.Context) ([]models.Team, error) {
	// This assumes we have a way to identify Premier League teams
	// For now, we'll just return all teams
	return r.GetAll(ctx)
}
