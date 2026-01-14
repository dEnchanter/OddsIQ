package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/dEnchanter/OddsIQ/backend/internal/models"
)

// TeamStatsRepository handles team statistics database operations
type TeamStatsRepository struct {
	db *pgxpool.Pool
}

// NewTeamStatsRepository creates a new team stats repository
func NewTeamStatsRepository(db *pgxpool.Pool) *TeamStatsRepository {
	return &TeamStatsRepository{db: db}
}

// Create inserts new team stats
func (r *TeamStatsRepository) Create(ctx context.Context, stats *models.TeamStats) error {
	query := `
		INSERT INTO team_stats (
			team_id, season, matches_played, wins, draws, losses,
			goals_for, goals_against, goal_difference, points,
			home_wins, home_draws, home_losses, away_wins, away_draws, away_losses,
			form, clean_sheets, failed_to_score,
			avg_goals_scored, avg_goals_conceded, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		stats.TeamID,
		stats.Season,
		stats.MatchesPlayed,
		stats.Wins,
		stats.Draws,
		stats.Losses,
		stats.GoalsFor,
		stats.GoalsAgainst,
		stats.GoalDifference,
		stats.Points,
		stats.HomeWins,
		stats.HomeDraws,
		stats.HomeLosses,
		stats.AwayWins,
		stats.AwayDraws,
		stats.AwayLosses,
		stats.Form,
		stats.CleanSheets,
		stats.FailedToScore,
		stats.AvgGoalsScored,
		stats.AvgGoalsConceded,
		now,
		now,
	).Scan(&stats.ID)

	if err != nil {
		return fmt.Errorf("failed to create team stats: %w", err)
	}

	stats.CreatedAt = now
	stats.UpdatedAt = now

	return nil
}

// GetByID retrieves team stats by ID
func (r *TeamStatsRepository) GetByID(ctx context.Context, id int) (*models.TeamStats, error) {
	query := `
		SELECT id, team_id, season, matches_played, wins, draws, losses,
			goals_for, goals_against, goal_difference, points,
			home_wins, home_draws, home_losses, away_wins, away_draws, away_losses,
			form, clean_sheets, failed_to_score,
			avg_goals_scored, avg_goals_conceded, created_at, updated_at
		FROM team_stats
		WHERE id = $1
	`

	stats := &models.TeamStats{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&stats.ID,
		&stats.TeamID,
		&stats.Season,
		&stats.MatchesPlayed,
		&stats.Wins,
		&stats.Draws,
		&stats.Losses,
		&stats.GoalsFor,
		&stats.GoalsAgainst,
		&stats.GoalDifference,
		&stats.Points,
		&stats.HomeWins,
		&stats.HomeDraws,
		&stats.HomeLosses,
		&stats.AwayWins,
		&stats.AwayDraws,
		&stats.AwayLosses,
		&stats.Form,
		&stats.CleanSheets,
		&stats.FailedToScore,
		&stats.AvgGoalsScored,
		&stats.AvgGoalsConceded,
		&stats.CreatedAt,
		&stats.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("team stats not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team stats: %w", err)
	}

	return stats, nil
}

// GetByTeamAndSeason retrieves team stats for a specific team and season
func (r *TeamStatsRepository) GetByTeamAndSeason(ctx context.Context, teamID, season int) (*models.TeamStats, error) {
	query := `
		SELECT id, team_id, season, matches_played, wins, draws, losses,
			goals_for, goals_against, goal_difference, points,
			home_wins, home_draws, home_losses, away_wins, away_draws, away_losses,
			form, clean_sheets, failed_to_score,
			avg_goals_scored, avg_goals_conceded, created_at, updated_at
		FROM team_stats
		WHERE team_id = $1 AND season = $2
	`

	stats := &models.TeamStats{}
	err := r.db.QueryRow(ctx, query, teamID, season).Scan(
		&stats.ID,
		&stats.TeamID,
		&stats.Season,
		&stats.MatchesPlayed,
		&stats.Wins,
		&stats.Draws,
		&stats.Losses,
		&stats.GoalsFor,
		&stats.GoalsAgainst,
		&stats.GoalDifference,
		&stats.Points,
		&stats.HomeWins,
		&stats.HomeDraws,
		&stats.HomeLosses,
		&stats.AwayWins,
		&stats.AwayDraws,
		&stats.AwayLosses,
		&stats.Form,
		&stats.CleanSheets,
		&stats.FailedToScore,
		&stats.AvgGoalsScored,
		&stats.AvgGoalsConceded,
		&stats.CreatedAt,
		&stats.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("team stats not found for team %d season %d", teamID, season)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team stats: %w", err)
	}

	return stats, nil
}

// GetBySeason retrieves all team stats for a specific season
func (r *TeamStatsRepository) GetBySeason(ctx context.Context, season int) ([]models.TeamStats, error) {
	query := `
		SELECT id, team_id, season, matches_played, wins, draws, losses,
			goals_for, goals_against, goal_difference, points,
			home_wins, home_draws, home_losses, away_wins, away_draws, away_losses,
			form, clean_sheets, failed_to_score,
			avg_goals_scored, avg_goals_conceded, created_at, updated_at
		FROM team_stats
		WHERE season = $1
		ORDER BY points DESC, goal_difference DESC
	`

	rows, err := r.db.Query(ctx, query, season)
	if err != nil {
		return nil, fmt.Errorf("failed to query team stats: %w", err)
	}
	defer rows.Close()

	return r.scanTeamStats(rows)
}

// GetByTeam retrieves all stats for a specific team across all seasons
func (r *TeamStatsRepository) GetByTeam(ctx context.Context, teamID int) ([]models.TeamStats, error) {
	query := `
		SELECT id, team_id, season, matches_played, wins, draws, losses,
			goals_for, goals_against, goal_difference, points,
			home_wins, home_draws, home_losses, away_wins, away_draws, away_losses,
			form, clean_sheets, failed_to_score,
			avg_goals_scored, avg_goals_conceded, created_at, updated_at
		FROM team_stats
		WHERE team_id = $1
		ORDER BY season DESC
	`

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query team stats: %w", err)
	}
	defer rows.Close()

	return r.scanTeamStats(rows)
}

// Update updates existing team stats
func (r *TeamStatsRepository) Update(ctx context.Context, stats *models.TeamStats) error {
	query := `
		UPDATE team_stats
		SET matches_played = $1, wins = $2, draws = $3, losses = $4,
			goals_for = $5, goals_against = $6, goal_difference = $7, points = $8,
			home_wins = $9, home_draws = $10, home_losses = $11,
			away_wins = $12, away_draws = $13, away_losses = $14,
			form = $15, clean_sheets = $16, failed_to_score = $17,
			avg_goals_scored = $18, avg_goals_conceded = $19, updated_at = $20
		WHERE id = $21
	`

	now := time.Now()
	result, err := r.db.Exec(ctx, query,
		stats.MatchesPlayed,
		stats.Wins,
		stats.Draws,
		stats.Losses,
		stats.GoalsFor,
		stats.GoalsAgainst,
		stats.GoalDifference,
		stats.Points,
		stats.HomeWins,
		stats.HomeDraws,
		stats.HomeLosses,
		stats.AwayWins,
		stats.AwayDraws,
		stats.AwayLosses,
		stats.Form,
		stats.CleanSheets,
		stats.FailedToScore,
		stats.AvgGoalsScored,
		stats.AvgGoalsConceded,
		now,
		stats.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update team stats: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("team stats not found with id %d", stats.ID)
	}

	stats.UpdatedAt = now

	return nil
}

// Upsert inserts or updates team stats
func (r *TeamStatsRepository) Upsert(ctx context.Context, stats *models.TeamStats) error {
	query := `
		INSERT INTO team_stats (
			team_id, season, matches_played, wins, draws, losses,
			goals_for, goals_against, goal_difference, points,
			home_wins, home_draws, home_losses, away_wins, away_draws, away_losses,
			form, clean_sheets, failed_to_score,
			avg_goals_scored, avg_goals_conceded, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
		ON CONFLICT (team_id, season)
		DO UPDATE SET
			matches_played = EXCLUDED.matches_played,
			wins = EXCLUDED.wins,
			draws = EXCLUDED.draws,
			losses = EXCLUDED.losses,
			goals_for = EXCLUDED.goals_for,
			goals_against = EXCLUDED.goals_against,
			goal_difference = EXCLUDED.goal_difference,
			points = EXCLUDED.points,
			home_wins = EXCLUDED.home_wins,
			home_draws = EXCLUDED.home_draws,
			home_losses = EXCLUDED.home_losses,
			away_wins = EXCLUDED.away_wins,
			away_draws = EXCLUDED.away_draws,
			away_losses = EXCLUDED.away_losses,
			form = EXCLUDED.form,
			clean_sheets = EXCLUDED.clean_sheets,
			failed_to_score = EXCLUDED.failed_to_score,
			avg_goals_scored = EXCLUDED.avg_goals_scored,
			avg_goals_conceded = EXCLUDED.avg_goals_conceded,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		stats.TeamID,
		stats.Season,
		stats.MatchesPlayed,
		stats.Wins,
		stats.Draws,
		stats.Losses,
		stats.GoalsFor,
		stats.GoalsAgainst,
		stats.GoalDifference,
		stats.Points,
		stats.HomeWins,
		stats.HomeDraws,
		stats.HomeLosses,
		stats.AwayWins,
		stats.AwayDraws,
		stats.AwayLosses,
		stats.Form,
		stats.CleanSheets,
		stats.FailedToScore,
		stats.AvgGoalsScored,
		stats.AvgGoalsConceded,
		now,
		now,
	).Scan(&stats.ID)

	if err != nil {
		return fmt.Errorf("failed to upsert team stats: %w", err)
	}

	stats.UpdatedAt = now

	return nil
}

// Delete deletes team stats
func (r *TeamStatsRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM team_stats WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete team stats: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("team stats not found with id %d", id)
	}

	return nil
}

// GetTopTeams retrieves top N teams by points for a specific season
func (r *TeamStatsRepository) GetTopTeams(ctx context.Context, season, limit int) ([]models.TeamStats, error) {
	query := `
		SELECT id, team_id, season, matches_played, wins, draws, losses,
			goals_for, goals_against, goal_difference, points,
			home_wins, home_draws, home_losses, away_wins, away_draws, away_losses,
			form, clean_sheets, failed_to_score,
			avg_goals_scored, avg_goals_conceded, created_at, updated_at
		FROM team_stats
		WHERE season = $1
		ORDER BY points DESC, goal_difference DESC, goals_for DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, season, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top teams: %w", err)
	}
	defer rows.Close()

	return r.scanTeamStats(rows)
}

// Helper function to scan team stats from rows
func (r *TeamStatsRepository) scanTeamStats(rows pgx.Rows) ([]models.TeamStats, error) {
	var statsList []models.TeamStats
	for rows.Next() {
		var stats models.TeamStats
		err := rows.Scan(
			&stats.ID,
			&stats.TeamID,
			&stats.Season,
			&stats.MatchesPlayed,
			&stats.Wins,
			&stats.Draws,
			&stats.Losses,
			&stats.GoalsFor,
			&stats.GoalsAgainst,
			&stats.GoalDifference,
			&stats.Points,
			&stats.HomeWins,
			&stats.HomeDraws,
			&stats.HomeLosses,
			&stats.AwayWins,
			&stats.AwayDraws,
			&stats.AwayLosses,
			&stats.Form,
			&stats.CleanSheets,
			&stats.FailedToScore,
			&stats.AvgGoalsScored,
			&stats.AvgGoalsConceded,
			&stats.CreatedAt,
			&stats.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team stats: %w", err)
		}
		statsList = append(statsList, stats)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return statsList, nil
}
