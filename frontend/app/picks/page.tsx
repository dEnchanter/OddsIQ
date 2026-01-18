'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { MultiMarketPick, PicksSummary, BetOutcome } from '@/types';
import { getMultiMarketPicks } from '@/lib/api';

export default function PicksPage() {
  const [picks, setPicks] = useState<MultiMarketPick[]>([]);
  const [summary, setSummary] = useState<PicksSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [bankroll, setBankroll] = useState<string>('1000');
  const [marketFilter, setMarketFilter] = useState<string>('all');

  // Load picks
  const loadPicks = async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await getMultiMarketPicks(parseFloat(bankroll) || 1000, 20);
      setPicks(result.picks || []);
      setSummary(result.summary || null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load picks');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPicks();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Filter picks by market
  const filteredPicks = picks.filter((pick) => {
    if (marketFilter === 'all') return true;
    return pick.best_outcome?.market === marketFilter;
  });

  // Format date
  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('en-GB', {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  // Get confidence color
  const getConfidenceColor = (ev: number) => {
    if (ev >= 0.15) return 'bg-green-500';
    if (ev >= 0.08) return 'bg-yellow-500';
    return 'bg-blue-500';
  };

  // Get market display name
  const getMarketName = (market: string) => {
    const names: Record<string, string> = {
      '1x2': '1X2',
      'over_under': 'O/U 2.5',
      'btts': 'BTTS',
    };
    return names[market] || market;
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">Loading picks...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Weekly Picks</h1>
          <p className="text-muted-foreground">AI-powered betting recommendations</p>
        </div>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <Label htmlFor="bankroll">Bankroll:</Label>
            <Input
              id="bankroll"
              type="number"
              value={bankroll}
              onChange={(e) => setBankroll(e.target.value)}
              className="w-24"
            />
          </div>
          <Button onClick={loadPicks}>Refresh</Button>
        </div>
      </div>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {/* Summary Card */}
      {summary && (
        <Card>
          <CardContent className="pt-6">
            <div className="grid grid-cols-2 md:grid-cols-5 gap-4 text-center">
              <div>
                <p className="text-2xl font-bold">{summary.total_picks}</p>
                <p className="text-sm text-muted-foreground">Total Picks</p>
              </div>
              <div>
                <p className="text-2xl font-bold">{summary.total_value_bets}</p>
                <p className="text-sm text-muted-foreground">Value Bets</p>
              </div>
              <div>
                <p className="text-2xl font-bold">${summary.total_suggested_stake.toFixed(0)}</p>
                <p className="text-sm text-muted-foreground">Total Stake</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-green-600">
                  +{(summary.average_ev * 100).toFixed(1)}%
                </p>
                <p className="text-sm text-muted-foreground">Avg EV</p>
              </div>
              <div>
                <p className="text-2xl font-bold">${summary.bankroll.toFixed(0)}</p>
                <p className="text-sm text-muted-foreground">Bankroll</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Filters */}
      <div className="flex gap-4">
        <Select value={marketFilter} onValueChange={setMarketFilter}>
          <SelectTrigger className="w-40">
            <SelectValue placeholder="Filter by market" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Markets</SelectItem>
            <SelectItem value="1x2">1X2</SelectItem>
            <SelectItem value="over_under">Over/Under</SelectItem>
            <SelectItem value="btts">BTTS</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Picks List */}
      {filteredPicks.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-muted-foreground">
              No picks available. Add fixtures and odds on the Entry page first.
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          {filteredPicks.map((pick, index) => (
            <Card key={pick.fixture.id}>
              <CardHeader className="pb-2">
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="text-lg">
                      #{index + 1} {pick.fixture.home_team_id} vs {pick.fixture.away_team_id}
                    </CardTitle>
                    <CardDescription>{formatDate(pick.fixture.match_date)}</CardDescription>
                  </div>
                  {pick.best_outcome && (
                    <Badge className={getConfidenceColor(pick.best_outcome.ev)}>
                      EV: +{(pick.best_outcome.ev * 100).toFixed(1)}%
                    </Badge>
                  )}
                </div>
              </CardHeader>
              <CardContent>
                {pick.best_outcome && (
                  <div className="bg-muted rounded-lg p-4 mb-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="font-semibold text-lg">{pick.best_outcome.description}</p>
                        <p className="text-sm text-muted-foreground">
                          {getMarketName(pick.best_outcome.market)} @ {pick.best_outcome.best_odds.toFixed(2)}
                        </p>
                      </div>
                      <div className="text-right">
                        <p className="font-bold text-lg">${pick.suggested_stake.toFixed(2)}</p>
                        <p className="text-sm text-muted-foreground">Suggested Stake</p>
                      </div>
                    </div>
                    <div className="mt-3 grid grid-cols-3 gap-4 text-sm">
                      <div>
                        <span className="text-muted-foreground">Probability: </span>
                        <span className="font-medium">{(pick.best_outcome.probability * 100).toFixed(0)}%</span>
                      </div>
                      <div>
                        <span className="text-muted-foreground">Kelly: </span>
                        <span className="font-medium">${pick.best_outcome.kelly_stake.toFixed(2)}</span>
                      </div>
                      <div>
                        <span className="text-muted-foreground">Bookmaker: </span>
                        <span className="font-medium">{pick.best_outcome.bookmaker}</span>
                      </div>
                    </div>
                  </div>
                )}

                {/* All Value Outcomes */}
                {pick.value_outcomes.length > 1 && (
                  <div>
                    <p className="text-sm font-medium mb-2">All Value Bets:</p>
                    <div className="space-y-2">
                      {pick.value_outcomes.map((outcome, i) => (
                        <div
                          key={i}
                          className="flex items-center justify-between text-sm py-1 border-b last:border-0"
                        >
                          <div className="flex items-center gap-2">
                            <Badge variant="outline" className="text-xs">
                              {getMarketName(outcome.market)}
                            </Badge>
                            <span>{outcome.description}</span>
                          </div>
                          <div className="flex items-center gap-4">
                            <span>{(outcome.probability * 100).toFixed(0)}%</span>
                            <span className="font-medium">{outcome.best_odds.toFixed(2)}</span>
                            <span className="text-green-600">+{(outcome.ev * 100).toFixed(1)}%</span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
