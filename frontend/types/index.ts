// Team types
export interface Team {
  id: number;
  api_football_id: number;
  name: string;
  code: string;
  logo_url: string;
  venue_name: string;
  venue_city: string;
}

// Fixture types
export interface Fixture {
  id: number;
  api_football_id: number;
  season: number;
  round: string;
  match_date: string;
  home_team_id: number;
  away_team_id: number;
  home_score: number | null;
  away_score: number | null;
  status: string;
  venue_name: string;
}

export interface EnrichedFixture extends Fixture {
  home_team_name: string;
  away_team_name: string;
  has_odds: boolean;
  odds_count: number;
}

// Odds types
export interface Odds {
  id: number;
  fixture_id: number;
  bookmaker: string;
  market_type: string;
  outcome: string;
  odds_value: number;
  recorded_at: string;
}

// Prediction types
export interface BetOutcome {
  market: string;
  outcome: string;
  description: string;
  probability: number;
  best_odds: number;
  bookmaker: string;
  ev: number;
  ev_percent: number;
  kelly_stake: number;
  confidence: number;
}

export interface MultiMarketPick {
  fixture: Fixture;
  all_outcomes: BetOutcome[];
  best_outcome: BetOutcome | null;
  value_outcomes: BetOutcome[];
  suggested_stake: number;
  total_ev: number;
  evaluated_at: string;
}

export interface PicksSummary {
  total_picks: number;
  total_value_bets: number;
  total_suggested_stake: number;
  total_expected_value: number;
  picks_by_market: Record<string, number>;
  average_ev: number;
  bankroll: number;
}

// Accumulator types
export interface AccumulatorLeg {
  fixture_id: number;
  fixture: Fixture;
  market: string;
  outcome: string;
  description: string;
  probability: number;
  odds: number;
  bookmaker: string;
  single_ev: number;
}

export interface Accumulator {
  id: string;
  legs: AccumulatorLeg[];
  num_legs: number;
  combined_probability: number;
  combined_odds: number;
  expected_value: number;
  ev_percent: number;
  suggested_stake: number;
  potential_return: number;
  confidence: string;
  generated_at: string;
}

export interface AccumulatorSummary {
  total_accumulators: number;
  total_doubles: number;
  total_trebles: number;
  total_suggested_stake: number;
  total_potential_return: number;
  average_ev: number;
  best_ev: number;
  bankroll: number;
  max_stake_allocation: number;
}

export interface AccumulatorConfig {
  min_legs: number;
  max_legs: number;
  min_ev_threshold: number;
  min_leg_ev: number;
  min_leg_probability: number;
  kelly_fraction: number;
  max_stake_percent: number;
  allow_same_team: boolean;
  allow_same_fixture: boolean;
}

// Bet types
export interface Bet {
  id: number;
  fixture_id: number;
  fixture?: Fixture;
  prediction_id: number | null;
  bet_type: string;
  stake: number;
  odds: number;
  expected_value: number;
  bookmaker: string;
  placed_at: string;
  status: 'pending' | 'won' | 'lost';
  payout: number | null;
  profit_loss: number | null;
  settled_at: string | null;
  notes: string;
}

// Performance types
export interface PerformanceMetrics {
  total_bets: number;
  total_staked: number;
  total_returned: number;
  total_profit: number;
  roi_percentage: number;
  win_rate: number;
  avg_odds: number;
  avg_stake: number;
  num_wins: number;
  num_losses: number;
  biggest_win: number;
  biggest_loss: number;
}

export interface DailyPerformance {
  date: string;
  bets: number;
  won: number;
  lost: number;
  profit: number;
  roi: number;
}

export interface BankrollSnapshot {
  id: number;
  balance: number;
  total_staked: number;
  total_returned: number;
  total_profit_loss: number;
  roi_percentage: number;
  num_bets: number;
  num_wins: number;
  num_losses: number;
  win_rate: number;
  recorded_at: string;
}

// Model metrics types
export interface ModelMetrics {
  model_version: string;
  training_date: string;
  accuracy: number;
  baseline_accuracy: number;
  improvement: number;
  config_name: string;
  feature_count: number;
}

// API Request types
export interface ManualFixtureRequest {
  home_team_id: number;
  away_team_id: number;
  match_date: string;
  season: number;
  round?: string;
  venue_name?: string;
}

export interface ManualOddsRequest {
  fixture_id: number;
  bookmaker: string;
  market_type: string;
  outcome: string;
  odds_value: number;
}

export interface ManualOddsBatchRequest {
  fixture_id: number;
  bookmaker: string;
  odds: Array<{
    market_type: string;
    outcome: string;
    odds_value: number;
  }>;
}

export interface CreateBetRequest {
  fixture_id: number;
  bet_type: string;
  stake: number;
  odds: number;
  expected_value: number;
  bookmaker: string;
  notes?: string;
}

export interface SettleBetRequest {
  result: 'won' | 'lost';
}

// API Response types
export interface TeamsResponse {
  teams: Team[];
  total: number;
}

export interface FixturesResponse {
  fixtures: Fixture[];
  total: number;
}

export interface UpcomingFixturesResponse {
  fixtures: EnrichedFixture[];
  total: number;
}

export interface MultiMarketPicksResponse {
  picks: MultiMarketPick[];
  summary: PicksSummary;
}

export interface AccumulatorsResponse {
  accumulators: Accumulator[];
  summary: AccumulatorSummary;
  config: AccumulatorConfig;
  generated_at: string;
}

export interface BetsResponse {
  bets: Bet[];
  total: number;
}

export interface PerformanceSummaryResponse {
  metrics: PerformanceMetrics;
}

export interface DailyPerformanceResponse {
  daily_performance: DailyPerformance[];
}

export interface BankrollHistoryResponse {
  history: BankrollSnapshot[];
}

export interface ModelMetricsResponse {
  markets: Record<string, ModelMetrics>;
  available_markets: string[];
}

export interface HealthResponse {
  status: string;
  service: string;
  version?: string;
}
