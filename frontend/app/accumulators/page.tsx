'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Accumulator, AccumulatorSummary, AccumulatorConfig } from '@/types';
import { getAccumulators } from '@/lib/api';

export default function AccumulatorsPage() {
  const [accumulators, setAccumulators] = useState<Accumulator[]>([]);
  const [summary, setSummary] = useState<AccumulatorSummary | null>(null);
  const [config, setConfig] = useState<AccumulatorConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [bankroll, setBankroll] = useState<string>('1000');

  // Load accumulators
  const loadAccumulators = async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await getAccumulators(parseFloat(bankroll) || 1000);
      setAccumulators(result.accumulators || []);
      setSummary(result.summary || null);
      setConfig(result.config || null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load accumulators');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAccumulators();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

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

  // Get market display name
  const getMarketName = (market: string) => {
    const names: Record<string, string> = {
      '1x2': '1X2',
      'over_under': 'O/U 2.5',
      'btts': 'BTTS',
    };
    return names[market] || market;
  };

  // Get confidence badge color
  const getConfidenceBadge = (confidence: string) => {
    switch (confidence.toLowerCase()) {
      case 'high':
        return <Badge className="bg-green-500">High</Badge>;
      case 'medium':
        return <Badge className="bg-yellow-500">Medium</Badge>;
      default:
        return <Badge className="bg-blue-500">Low</Badge>;
    }
  };

  // Get accumulator type name
  const getAccaType = (numLegs: number) => {
    switch (numLegs) {
      case 2:
        return 'Double';
      case 3:
        return 'Treble';
      case 4:
        return 'Fourfold';
      case 5:
        return 'Fivefold';
      default:
        return `${numLegs}-fold`;
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">Loading accumulators...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Accumulators</h1>
          <p className="text-muted-foreground">AI-generated multi-bet recommendations</p>
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
          <Button onClick={loadAccumulators}>Refresh</Button>
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
            <div className="grid grid-cols-2 md:grid-cols-6 gap-4 text-center">
              <div>
                <p className="text-2xl font-bold">{summary.total_accumulators}</p>
                <p className="text-sm text-muted-foreground">Total Accas</p>
              </div>
              <div>
                <p className="text-2xl font-bold">{summary.total_doubles}</p>
                <p className="text-sm text-muted-foreground">Doubles</p>
              </div>
              <div>
                <p className="text-2xl font-bold">{summary.total_trebles}</p>
                <p className="text-sm text-muted-foreground">Trebles</p>
              </div>
              <div>
                <p className="text-2xl font-bold">${summary.total_suggested_stake.toFixed(0)}</p>
                <p className="text-sm text-muted-foreground">Total Stake</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-green-600">
                  ${summary.total_potential_return.toFixed(0)}
                </p>
                <p className="text-sm text-muted-foreground">Potential Return</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-green-600">
                  +{(summary.average_ev * 100).toFixed(1)}%
                </p>
                <p className="text-sm text-muted-foreground">Avg EV</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Config Info */}
      {config && (
        <Card>
          <CardHeader className="py-3">
            <CardDescription className="text-sm">
              Legs: {config.min_legs}-{config.max_legs} | Min EV: {(config.min_ev_threshold * 100).toFixed(0)}% |
              Kelly: {config.kelly_fraction}x | Max Stake: {(config.max_stake_percent * 100).toFixed(0)}% of bankroll
            </CardDescription>
          </CardHeader>
        </Card>
      )}

      {/* Accumulators List */}
      {accumulators.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-muted-foreground">
              No accumulators available. Add fixtures and odds on the Entry page first.
            </p>
            <p className="text-sm text-muted-foreground mt-2">
              Accumulators require multiple fixtures with positive EV outcomes.
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          {accumulators.map((acca, index) => (
            <Card key={acca.id}>
              <CardHeader className="pb-2">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <span className="text-lg font-bold text-muted-foreground">
                      #{index + 1}
                    </span>
                    <div>
                      <CardTitle className="text-lg">
                        {getAccaType(acca.num_legs)} - {acca.num_legs} Legs
                      </CardTitle>
                      <CardDescription>
                        Combined Odds: {acca.combined_odds.toFixed(2)}
                      </CardDescription>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {getConfidenceBadge(acca.confidence)}
                    <Badge className="bg-green-500">
                      EV: +{(acca.ev_percent).toFixed(1)}%
                    </Badge>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                {/* Legs */}
                <div className="space-y-2 mb-4">
                  {acca.legs.map((leg, legIndex) => (
                    <div
                      key={legIndex}
                      className="flex items-center justify-between p-3 bg-muted/50 rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <span className="text-sm font-medium text-muted-foreground">
                          {legIndex + 1}.
                        </span>
                        <div>
                          <p className="font-medium">{leg.description}</p>
                          <p className="text-sm text-muted-foreground">
                            {leg.fixture.match_date ? formatDate(leg.fixture.match_date) : 'TBD'}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center gap-4">
                        <Badge variant="outline" className="text-xs">
                          {getMarketName(leg.market)}
                        </Badge>
                        <span className="text-sm text-muted-foreground">
                          {(leg.probability * 100).toFixed(0)}%
                        </span>
                        <span className="font-medium">{leg.odds.toFixed(2)}</span>
                        <span className="text-green-600 text-sm">
                          +{(leg.single_ev * 100).toFixed(1)}%
                        </span>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Summary Row */}
                <div className="flex items-center justify-between pt-4 border-t">
                  <div className="grid grid-cols-3 gap-8 text-sm">
                    <div>
                      <span className="text-muted-foreground">Combined Prob: </span>
                      <span className="font-medium">{(acca.combined_probability * 100).toFixed(1)}%</span>
                    </div>
                    <div>
                      <span className="text-muted-foreground">Combined Odds: </span>
                      <span className="font-medium">{acca.combined_odds.toFixed(2)}</span>
                    </div>
                    <div>
                      <span className="text-muted-foreground">Expected Value: </span>
                      <span className="font-medium text-green-600">+{(acca.expected_value * 100).toFixed(1)}%</span>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="flex items-center gap-4">
                      <div>
                        <p className="text-sm text-muted-foreground">Stake</p>
                        <p className="font-bold text-lg">${acca.suggested_stake.toFixed(2)}</p>
                      </div>
                      <div className="text-green-600">
                        <p className="text-sm text-muted-foreground">Potential Return</p>
                        <p className="font-bold text-lg">${acca.potential_return.toFixed(2)}</p>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Info Card */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">About Accumulators</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2 text-sm text-muted-foreground">
            <p>
              Accumulators combine 2-3 uncorrelated betting picks into a single bet with higher potential returns.
            </p>
            <p>
              <strong>Correlation Rules:</strong> No same fixture, no same team across legs.
            </p>
            <p>
              <strong>Stake Strategy:</strong> Uses 1/8 Kelly (more conservative than singles) with max 20% of weekly stake.
            </p>
            <p>
              <strong>EV Threshold:</strong> Only includes accumulators with 5%+ expected value.
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
