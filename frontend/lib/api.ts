const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000/api';

export interface Fixture {
  id: number;
  match_date: string;
  home_team: {
    id: number;
    name: string;
    code: string;
  };
  away_team: {
    id: number;
    name: string;
    code: string;
  };
  status: string;
}

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

export interface PerformanceMetrics {
  total_bets: number;
  total_staked: number;
  total_returned: number;
  total_profit: number;
  roi_percentage: number;
  win_rate: number;
  avg_odds: number;
}

/**
 * Fetch weekly betting picks
 */
export async function getWeeklyPicks(): Promise<{ picks: WeeklyPick[] }> {
  const response = await fetch(`${API_URL}/picks/weekly`);
  if (!response.ok) {
    throw new Error('Failed to fetch weekly picks');
  }
  return response.json();
}

/**
 * Fetch performance summary
 */
export async function getPerformanceSummary(): Promise<{ metrics: PerformanceMetrics }> {
  const response = await fetch(`${API_URL}/performance/summary`);
  if (!response.ok) {
    throw new Error('Failed to fetch performance summary');
  }
  return response.json();
}

/**
 * Fetch fixtures
 */
export async function getFixtures(params?: {
  from_date?: string;
  to_date?: string;
  status?: string;
}): Promise<{ fixtures: Fixture[]; total: number }> {
  const queryParams = new URLSearchParams(params as any);
  const response = await fetch(`${API_URL}/fixtures?${queryParams}`);
  if (!response.ok) {
    throw new Error('Failed to fetch fixtures');
  }
  return response.json();
}

/**
 * Health check
 */
export async function healthCheck(): Promise<{ status: string; service: string }> {
  const response = await fetch(`${API_URL.replace('/api', '')}/health`);
  if (!response.ok) {
    throw new Error('Health check failed');
  }
  return response.json();
}
