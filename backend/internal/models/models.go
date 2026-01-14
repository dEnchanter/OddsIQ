package models

import "time"

// Team represents a football team
type Team struct {
	ID             int       `json:"id"`
	APIFootballID  int       `json:"api_football_id"`
	Name           string    `json:"name"`
	Code           string    `json:"code"`
	LogoURL        string    `json:"logo_url"`
	Founded        int       `json:"founded"`
	VenueName      string    `json:"venue_name"`
	VenueCity      string    `json:"venue_city"`
	VenueCapacity  int       `json:"venue_capacity"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Fixture represents a match fixture
type Fixture struct {
	ID             int       `json:"id"`
	APIFootballID  int       `json:"api_football_id"`
	Season         int       `json:"season"`
	Round          string    `json:"round"`
	MatchDate      time.Time `json:"match_date"`
	HomeTeamID     int       `json:"home_team_id"`
	AwayTeamID     int       `json:"away_team_id"`
	HomeTeam       *Team     `json:"home_team,omitempty"`
	AwayTeam       *Team     `json:"away_team,omitempty"`
	HomeScore      *int      `json:"home_score"`
	AwayScore      *int      `json:"away_score"`
	Status         string    `json:"status"`
	VenueName      string    `json:"venue"`
	Referee        string    `json:"referee"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Odds represents bookmaker odds for a fixture
type Odds struct {
	ID            int       `json:"id"`
	FixtureID     int       `json:"fixture_id"`
	Bookmaker     string    `json:"bookmaker"`
	MarketType    string    `json:"market_type"`
	Outcome       string    `json:"outcome"`
	OddsValue     float64   `json:"odds_value"`
	Timestamp     time.Time `json:"recorded_at"`
	IsClosingLine bool      `json:"is_closing_line"`
	CreatedAt     time.Time `json:"created_at"`
}

// TeamStats represents team statistics at a point in time
type TeamStats struct {
	ID               int       `json:"id"`
	TeamID           int       `json:"team_id"`
	Season           int       `json:"season"`
	MatchesPlayed    int       `json:"matches_played"`
	Wins             int       `json:"wins"`
	Draws            int       `json:"draws"`
	Losses           int       `json:"losses"`
	GoalsFor         int       `json:"goals_for"`
	GoalsAgainst     int       `json:"goals_against"`
	GoalDifference   int       `json:"goal_difference"`
	Points           int       `json:"points"`
	HomeWins         int       `json:"home_wins"`
	HomeDraws        int       `json:"home_draws"`
	HomeLosses       int       `json:"home_losses"`
	AwayWins         int       `json:"away_wins"`
	AwayDraws        int       `json:"away_draws"`
	AwayLosses       int       `json:"away_losses"`
	Form             string    `json:"form"`
	CleanSheets      int       `json:"clean_sheets"`
	FailedToScore    int       `json:"failed_to_score"`
	AvgGoalsScored   float64   `json:"avg_goals_scored"`
	AvgGoalsConceded float64   `json:"avg_goals_conceded"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Prediction represents a model prediction for a fixture
type Prediction struct {
	ID               int                    `json:"id"`
	FixtureID        int                    `json:"fixture_id"`
	ModelVersion     string                 `json:"model_version"`
	HomeWinProb      float64                `json:"home_win_prob"`
	DrawProb         float64                `json:"draw_prob"`
	AwayWinProb      float64                `json:"away_win_prob"`
	PredictedOutcome string                 `json:"predicted_outcome"`
	ConfidenceScore  float64                `json:"confidence_score"`
	Features         map[string]interface{} `json:"features"`
	PredictedAt      time.Time              `json:"predicted_at"`
	CreatedAt        time.Time              `json:"created_at"`
}

// Bet represents a placed bet
type Bet struct {
	ID            int       `json:"id"`
	FixtureID     int       `json:"fixture_id"`
	Fixture       *Fixture  `json:"fixture,omitempty"`
	PredictionID  *int      `json:"prediction_id"`
	BetType       string    `json:"bet_type"`
	Stake         float64   `json:"stake"`
	Odds          float64   `json:"odds"`
	ExpectedValue float64   `json:"expected_value"`
	Bookmaker     string    `json:"bookmaker"`
	PlacedAt      time.Time `json:"placed_at"`
	Status        string    `json:"status"`
	Payout        *float64  `json:"payout"`
	ProfitLoss    *float64  `json:"profit_loss"`
	SettledAt     *time.Time `json:"settled_at"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Bankroll represents bankroll snapshot
type Bankroll struct {
	ID              int       `json:"id"`
	Balance         float64   `json:"balance"`
	TotalStaked     float64   `json:"total_staked"`
	TotalReturned   float64   `json:"total_returned"`
	TotalProfitLoss float64   `json:"total_profit_loss"`
	ROIPercentage   float64   `json:"roi_percentage"`
	NumBets         int       `json:"num_bets"`
	NumWins         int       `json:"num_wins"`
	NumLosses       int       `json:"num_losses"`
	WinRate         float64   `json:"win_rate"`
	RecordedAt      time.Time `json:"recorded_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// WeeklyPick represents a betting recommendation
type WeeklyPick struct {
	Fixture      Fixture    `json:"fixture"`
	Prediction   Prediction `json:"prediction"`
	BestOdds     float64    `json:"best_odds"`
	Bookmaker    string     `json:"bookmaker"`
	ExpectedValue float64   `json:"expected_value"`
	EVPercentage float64    `json:"ev_percentage"`
	SuggestedStake float64  `json:"suggested_stake"`
	KellyFraction float64   `json:"kelly_fraction"`
	BetType      string     `json:"bet_type"`
	Confidence   string     `json:"confidence"`
}

// PerformanceMetrics represents performance summary
type PerformanceMetrics struct {
	TotalBets      int       `json:"total_bets"`
	TotalStaked    float64   `json:"total_staked"`
	TotalReturned  float64   `json:"total_returned"`
	TotalProfit    float64   `json:"total_profit"`
	ROIPercentage  float64   `json:"roi_percentage"`
	WinRate        float64   `json:"win_rate"`
	AvgOdds        float64   `json:"avg_odds"`
	AvgStake       float64   `json:"avg_stake"`
	NumWins        int       `json:"num_wins"`
	NumLosses      int       `json:"num_losses"`
	BiggestWin     float64   `json:"biggest_win"`
	BiggestLoss    float64   `json:"biggest_loss"`
	MaxDrawdown    float64   `json:"max_drawdown"`
	SharpeRatio    float64   `json:"sharpe_ratio"`
	CLVAverage     float64   `json:"clv_average"`
	FromDate       time.Time `json:"from_date"`
	ToDate         time.Time `json:"to_date"`
}
