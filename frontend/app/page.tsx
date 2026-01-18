'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import { PerformanceMetrics, MultiMarketPick, BankrollSnapshot } from '@/types';
import {
  getPerformanceSummary,
  getMultiMarketPicks,
  getBankrollHistory,
  healthCheck,
  getMLHealth,
} from '@/lib/api';

export default function DashboardPage() {
  const [metrics, setMetrics] = useState<PerformanceMetrics | null>(null);
  const [topPicks, setTopPicks] = useState<MultiMarketPick[]>([]);
  const [bankrollHistory, setBankrollHistory] = useState<BankrollSnapshot[]>([]);
  const [backendStatus, setBackendStatus] = useState<'online' | 'offline' | 'checking'>('checking');
  const [mlStatus, setMlStatus] = useState<'online' | 'offline' | 'checking'>('checking');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadData() {
      // Check backend health
      try {
        await healthCheck();
        setBackendStatus('online');
      } catch {
        setBackendStatus('offline');
      }

      // Check ML service health
      try {
        await getMLHealth();
        setMlStatus('online');
      } catch {
        setMlStatus('offline');
      }

      // Load dashboard data
      try {
        const [perfRes, picksRes, bankrollRes] = await Promise.all([
          getPerformanceSummary().catch(() => ({ metrics: null })),
          getMultiMarketPicks(1000, 5).catch(() => ({ picks: [] })),
          getBankrollHistory().catch(() => ({ history: [] })),
        ]);

        setMetrics(perfRes.metrics);
        setTopPicks(picksRes.picks || []);
        setBankrollHistory(bankrollRes.history || []);
      } catch (err) {
        console.error('Failed to load dashboard data:', err);
      } finally {
        setLoading(false);
      }
    }
    loadData();
  }, []);

  // Sample chart data
  const chartData = bankrollHistory.length > 0
    ? bankrollHistory.map((b, i) => ({
        name: `Day ${i + 1}`,
        balance: b.balance,
      }))
    : [
        { name: 'Start', balance: 1000 },
        { name: 'Week 1', balance: 1050 },
        { name: 'Week 2', balance: 1120 },
        { name: 'Week 3', balance: 1180 },
        { name: 'Week 4', balance: 1250 },
      ];

  // Format date
  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('en-GB', {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
    });
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <p className="text-muted-foreground">Welcome to OddsIQ</p>
        </div>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <span className={`h-2 w-2 rounded-full ${backendStatus === 'online' ? 'bg-green-500' : backendStatus === 'offline' ? 'bg-red-500' : 'bg-yellow-500'}`}></span>
            <span className="text-sm text-muted-foreground">Backend</span>
          </div>
          <div className="flex items-center gap-2">
            <span className={`h-2 w-2 rounded-full ${mlStatus === 'online' ? 'bg-green-500' : mlStatus === 'offline' ? 'bg-red-500' : 'bg-yellow-500'}`}></span>
            <span className="text-sm text-muted-foreground">ML Service</span>
          </div>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <p className="text-sm text-muted-foreground">Bankroll</p>
            <p className="text-3xl font-bold">
              ${bankrollHistory.length > 0
                ? bankrollHistory[bankrollHistory.length - 1].balance.toFixed(0)
                : '1,000'}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <p className="text-sm text-muted-foreground">ROI</p>
            <p className={`text-3xl font-bold ${(metrics?.roi_percentage || 0) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              {(metrics?.roi_percentage || 0) >= 0 ? '+' : ''}{(metrics?.roi_percentage || 0).toFixed(1)}%
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <p className="text-sm text-muted-foreground">Win Rate</p>
            <p className="text-3xl font-bold">
              {((metrics?.win_rate || 0) * 100).toFixed(0)}%
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <p className="text-sm text-muted-foreground">Profit</p>
            <p className={`text-3xl font-bold ${(metrics?.total_profit || 0) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              {(metrics?.total_profit || 0) >= 0 ? '+' : ''}${(metrics?.total_profit || 0).toFixed(0)}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Main Content */}
      <div className="grid md:grid-cols-2 gap-6">
        {/* Bankroll Chart */}
        <Card>
          <CardHeader>
            <CardTitle>Bankroll History</CardTitle>
            <CardDescription>Track your bankroll growth</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[250px]">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                  <XAxis dataKey="name" className="text-xs" />
                  <YAxis className="text-xs" />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: 'hsl(var(--background))',
                      border: '1px solid hsl(var(--border))',
                    }}
                  />
                  <Line
                    type="monotone"
                    dataKey="balance"
                    stroke="hsl(var(--primary))"
                    strokeWidth={2}
                    dot={{ fill: 'hsl(var(--primary))' }}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        {/* Top Picks */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>Top Picks Today</CardTitle>
              <CardDescription>Best value bets</CardDescription>
            </div>
            <Link href="/picks">
              <Button variant="outline" size="sm">View All</Button>
            </Link>
          </CardHeader>
          <CardContent>
            {topPicks.length === 0 ? (
              <div className="text-center py-8">
                <p className="text-muted-foreground">No picks available</p>
                <p className="text-sm text-muted-foreground mt-1">
                  Add fixtures on the Entry page
                </p>
                <Link href="/entry">
                  <Button variant="outline" className="mt-4">Go to Entry</Button>
                </Link>
              </div>
            ) : (
              <div className="space-y-3">
                {topPicks.map((pick, index) => (
                  <div
                    key={pick.fixture.id}
                    className="flex items-center justify-between p-3 rounded-lg bg-muted/50"
                  >
                    <div className="flex items-center gap-3">
                      <span className="text-lg font-bold text-muted-foreground">
                        #{index + 1}
                      </span>
                      <div>
                        <p className="font-medium text-sm">
                          {pick.best_outcome?.description || 'N/A'}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {formatDate(pick.fixture.match_date)}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <Badge className="bg-green-500">
                        +{((pick.best_outcome?.ev || 0) * 100).toFixed(1)}%
                      </Badge>
                      <p className="text-xs text-muted-foreground mt-1">
                        @ {pick.best_outcome?.best_odds.toFixed(2) || '-'}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-4">
            <Link href="/entry">
              <Button>Add Fixture</Button>
            </Link>
            <Link href="/picks">
              <Button variant="outline">View All Picks</Button>
            </Link>
            <Link href="/accumulators">
              <Button variant="outline">View Accumulators</Button>
            </Link>
            <Link href="/bets">
              <Button variant="outline">Record Bet</Button>
            </Link>
            <Link href="/performance">
              <Button variant="outline">View Stats</Button>
            </Link>
          </div>
        </CardContent>
      </Card>

      {/* Getting Started (show if no data) */}
      {!loading && topPicks.length === 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Getting Started</CardTitle>
            <CardDescription>Follow these steps to start using OddsIQ</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-start gap-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground text-sm font-bold">
                  1
                </div>
                <div>
                  <p className="font-medium">Add a Fixture</p>
                  <p className="text-sm text-muted-foreground">
                    Go to the Entry page and create a fixture for an upcoming match
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground text-sm font-bold">
                  2
                </div>
                <div>
                  <p className="font-medium">Add Odds</p>
                  <p className="text-sm text-muted-foreground">
                    Enter the odds from your bookmaker for all markets (1X2, O/U, BTTS)
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground text-sm font-bold">
                  3
                </div>
                <div>
                  <p className="font-medium">View Picks</p>
                  <p className="text-sm text-muted-foreground">
                    Check the Picks page for AI-powered betting recommendations
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground text-sm font-bold">
                  4
                </div>
                <div>
                  <p className="font-medium">Track Your Bets</p>
                  <p className="text-sm text-muted-foreground">
                    Record your bets and settle them to track performance
                  </p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
