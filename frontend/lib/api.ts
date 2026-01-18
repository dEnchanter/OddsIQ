import {
  Team,
  Fixture,
  EnrichedFixture,
  Odds,
  MultiMarketPick,
  PicksSummary,
  Accumulator,
  AccumulatorSummary,
  AccumulatorConfig,
  Bet,
  PerformanceMetrics,
  DailyPerformance,
  BankrollSnapshot,
  ModelMetrics,
  ManualFixtureRequest,
  ManualOddsBatchRequest,
  CreateBetRequest,
  SettleBetRequest,
  TeamsResponse,
  UpcomingFixturesResponse,
  MultiMarketPicksResponse,
  AccumulatorsResponse,
  BetsResponse,
  PerformanceSummaryResponse,
  DailyPerformanceResponse,
  BankrollHistoryResponse,
  ModelMetricsResponse,
  HealthResponse,
} from '@/types';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000/api';

// Helper for API calls with error handling
async function fetchApi<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}${endpoint}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    ...options,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(error.error || `HTTP ${response.status}`);
  }

  return response.json();
}

// ============= TEAMS =============

export async function getTeams(): Promise<TeamsResponse> {
  return fetchApi<TeamsResponse>('/teams');
}

// ============= FIXTURES =============

export async function getUpcomingFixtures(): Promise<UpcomingFixturesResponse> {
  return fetchApi<UpcomingFixturesResponse>('/fixtures/upcoming');
}

export async function getFixture(id: number): Promise<{
  fixture: Fixture;
  home_team: Team;
  away_team: Team;
}> {
  return fetchApi(`/fixtures/${id}`);
}

export async function getFixtureOdds(id: number): Promise<{
  fixture_id: number;
  odds: Odds[];
  market_types: string[];
  total: number;
}> {
  return fetchApi(`/fixtures/${id}/odds`);
}

export async function createManualFixture(data: ManualFixtureRequest): Promise<{
  fixture: Fixture;
  home_team: Team;
  away_team: Team;
  message: string;
}> {
  return fetchApi('/fixtures/manual', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function deleteFixture(id: number): Promise<{ message: string; fixture_id: number }> {
  return fetchApi(`/fixtures/${id}`, {
    method: 'DELETE',
  });
}

// ============= ODDS =============

export async function addOddsBatch(data: ManualOddsBatchRequest): Promise<{
  odds_count: number;
  fixture: Fixture;
  message: string;
}> {
  return fetchApi('/odds/manual/batch', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

// ============= PICKS =============

export async function getMultiMarketPicks(
  bankroll: number = 1000,
  limit: number = 15
): Promise<MultiMarketPicksResponse> {
  return fetchApi<MultiMarketPicksResponse>(`/picks/multi?bankroll=${bankroll}&limit=${limit}`);
}

// ============= ACCUMULATORS =============

export async function getAccumulators(bankroll: number = 1000): Promise<AccumulatorsResponse> {
  return fetchApi<AccumulatorsResponse>(`/accumulators/weekly?bankroll=${bankroll}`);
}

export async function getAccumulatorConfig(): Promise<{
  config: AccumulatorConfig;
  description: Record<string, string>;
}> {
  return fetchApi('/accumulators/config');
}

// ============= PREDICTIONS =============

export async function evaluateFixture(
  id: number,
  bankroll: number = 1000
): Promise<{
  fixture: Fixture;
  home_team: Team;
  away_team: Team;
  evaluation: {
    all_outcomes: Array<{
      market: string;
      outcome: string;
      description: string;
      probability: number;
      best_odds: number;
      ev: number;
      ev_percent: number;
      kelly_stake: number;
    }>;
    best_outcome: any;
    value_outcomes: any[];
  };
}> {
  return fetchApi(`/predictions/fixture/${id}/evaluate?bankroll=${bankroll}`);
}

// ============= BETS =============

export async function getBets(): Promise<BetsResponse> {
  return fetchApi<BetsResponse>('/bets');
}

export async function createBet(data: CreateBetRequest): Promise<{
  id: number;
  status: string;
}> {
  return fetchApi('/bets', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function settleBet(id: number, data: SettleBetRequest): Promise<{
  id: string;
  status: string;
}> {
  return fetchApi(`/bets/${id}/settle`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

// ============= PERFORMANCE =============

export async function getPerformanceSummary(): Promise<PerformanceSummaryResponse> {
  return fetchApi<PerformanceSummaryResponse>('/performance/summary');
}

export async function getDailyPerformance(): Promise<DailyPerformanceResponse> {
  return fetchApi<DailyPerformanceResponse>('/performance/daily');
}

// ============= BANKROLL =============

export async function getBankrollHistory(): Promise<BankrollHistoryResponse> {
  return fetchApi<BankrollHistoryResponse>('/bankroll/history');
}

// ============= MODEL =============

export async function getModelMetrics(): Promise<ModelMetricsResponse> {
  return fetchApi<ModelMetricsResponse>('/model/metrics/all');
}

export async function getMLHealth(): Promise<HealthResponse> {
  return fetchApi<HealthResponse>('/model/health');
}

// ============= HEALTH =============

export async function healthCheck(): Promise<HealthResponse> {
  const response = await fetch(`${API_URL.replace('/api', '')}/health`);
  if (!response.ok) {
    throw new Error('Health check failed');
  }
  return response.json();
}

// ============= LEGACY EXPORTS (for backward compatibility) =============

export interface WeeklyPick {
  fixture: Fixture;
  recommendation: {
    bet_type: string;
    outcome: string;
    model_probability: number;
    best_odds: number;
    bookmaker: string;
    expected_value: number;
    ev_percentage: number;
    suggested_stake: number;
    confidence: string;
  };
}

export async function getWeeklyPicks(): Promise<{ picks: WeeklyPick[] }> {
  return fetchApi('/picks/weekly');
}

export async function getFixtures(params?: {
  from_date?: string;
  to_date?: string;
  status?: string;
}): Promise<{ fixtures: Fixture[]; total: number }> {
  const queryParams = new URLSearchParams(params as Record<string, string>);
  return fetchApi(`/fixtures?${queryParams}`);
}
