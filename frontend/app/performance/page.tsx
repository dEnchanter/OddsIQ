'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
} from 'recharts';
import { PerformanceMetrics, DailyPerformance, BankrollSnapshot, ModelMetrics } from '@/types';
import {
  getPerformanceSummary,
  getDailyPerformance,
  getBankrollHistory,
  getModelMetrics,
} from '@/lib/api';

export default function PerformancePage() {
  const [metrics, setMetrics] = useState<PerformanceMetrics | null>(null);
  const [dailyPerformance, setDailyPerformance] = useState<DailyPerformance[]>([]);
  const [bankrollHistory, setBankrollHistory] = useState<BankrollSnapshot[]>([]);
  const [modelMetrics, setModelMetrics] = useState<Record<string, ModelMetrics>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadData() {
      try {
        setLoading(true);
        const [perfRes, dailyRes, bankrollRes, modelRes] = await Promise.all([
          getPerformanceSummary().catch(() => ({ metrics: null })),
          getDailyPerformance().catch(() => ({ daily_performance: [] })),
          getBankrollHistory().catch(() => ({ history: [] })),
          getModelMetrics().catch(() => ({ markets: {} })),
        ]);

        setMetrics(perfRes.metrics);
        setDailyPerformance(dailyRes.daily_performance || []);
        setBankrollHistory(bankrollRes.history || []);
        setModelMetrics(modelRes.markets || {});
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load data');
      } finally {
        setLoading(false);
      }
    }
    loadData();
  }, []);

  // Sample data for charts when no real data
  const sampleBankrollData = bankrollHistory.length > 0
    ? bankrollHistory.map((b, i) => ({
        name: `Day ${i + 1}`,
        balance: b.balance,
      }))
    : [
        { name: 'Week 1', balance: 1000 },
        { name: 'Week 2', balance: 1050 },
        { name: 'Week 3', balance: 1120 },
        { name: 'Week 4', balance: 1085 },
        { name: 'Week 5', balance: 1180 },
        { name: 'Week 6', balance: 1250 },
      ];

  const sampleDailyData = dailyPerformance.length > 0
    ? dailyPerformance.map((d) => ({
        date: d.date,
        profit: d.profit,
        bets: d.bets,
      }))
    : [
        { date: 'Mon', profit: 25, bets: 3 },
        { date: 'Tue', profit: -15, bets: 2 },
        { date: 'Wed', profit: 42, bets: 4 },
        { date: 'Thu', profit: 18, bets: 3 },
        { date: 'Fri', profit: -8, bets: 2 },
        { date: 'Sat', profit: 55, bets: 5 },
        { date: 'Sun', profit: 32, bets: 4 },
      ];

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">Loading performance data...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Performance Analytics</h1>
        <p className="text-muted-foreground">Track your betting performance over time</p>
      </div>

      {error && (
        <div className="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded">
          Using sample data - {error}
        </div>
      )}

      {/* KPI Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <p className="text-sm text-muted-foreground">Total Bets</p>
            <p className="text-3xl font-bold">{metrics?.total_bets || 0}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <p className="text-sm text-muted-foreground">Win Rate</p>
            <p className="text-3xl font-bold">
              {((metrics?.win_rate || 0) * 100).toFixed(1)}%
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
            <p className="text-sm text-muted-foreground">Total Profit</p>
            <p className={`text-3xl font-bold ${(metrics?.total_profit || 0) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              {(metrics?.total_profit || 0) >= 0 ? '+' : ''}${(metrics?.total_profit || 0).toFixed(2)}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Bankroll Chart */}
      <Card>
        <CardHeader>
          <CardTitle>Bankroll Over Time</CardTitle>
          <CardDescription>Track your bankroll growth</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={sampleBankrollData}>
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

      {/* Daily Performance Chart */}
      <Card>
        <CardHeader>
          <CardTitle>Daily Profit/Loss</CardTitle>
          <CardDescription>Performance breakdown by day</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[250px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={sampleDailyData}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                <XAxis dataKey="date" className="text-xs" />
                <YAxis className="text-xs" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: 'hsl(var(--background))',
                    border: '1px solid hsl(var(--border))',
                  }}
                />
                <Bar
                  dataKey="profit"
                  fill="hsl(var(--primary))"
                  radius={[4, 4, 0, 0]}
                />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </CardContent>
      </Card>

      {/* Model Performance */}
      <Card>
        <CardHeader>
          <CardTitle>Model Accuracy</CardTitle>
          <CardDescription>Performance of prediction models by market</CardDescription>
        </CardHeader>
        <CardContent>
          {Object.keys(modelMetrics).length === 0 ? (
            <div className="text-center py-8">
              <p className="text-muted-foreground">Model metrics not available</p>
              <p className="text-sm text-muted-foreground mt-1">
                Start the ML service to view model performance
              </p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Market</TableHead>
                  <TableHead>Accuracy</TableHead>
                  <TableHead>Baseline</TableHead>
                  <TableHead>Improvement</TableHead>
                  <TableHead>Features</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {Object.entries(modelMetrics).map(([market, metrics]) => (
                  <TableRow key={market}>
                    <TableCell className="font-medium uppercase">{market}</TableCell>
                    <TableCell>{(metrics.accuracy * 100).toFixed(1)}%</TableCell>
                    <TableCell>{(metrics.baseline_accuracy * 100).toFixed(1)}%</TableCell>
                    <TableCell className={metrics.improvement >= 0 ? 'text-green-600' : 'text-red-600'}>
                      {metrics.improvement >= 0 ? '+' : ''}{(metrics.improvement * 100).toFixed(1)}%
                    </TableCell>
                    <TableCell>{metrics.feature_count}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {/* Additional Stats */}
      <div className="grid md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Betting Stats</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Total Staked</span>
                <span className="font-medium">${(metrics?.total_staked || 0).toFixed(2)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Total Returned</span>
                <span className="font-medium">${(metrics?.total_returned || 0).toFixed(2)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Average Odds</span>
                <span className="font-medium">{(metrics?.avg_odds || 0).toFixed(2)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Average Stake</span>
                <span className="font-medium">${(metrics?.avg_stake || 0).toFixed(2)}</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Win/Loss Breakdown</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Wins</span>
                <span className="font-medium text-green-600">{metrics?.num_wins || 0}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Losses</span>
                <span className="font-medium text-red-600">{metrics?.num_losses || 0}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Biggest Win</span>
                <span className="font-medium text-green-600">
                  +${(metrics?.biggest_win || 0).toFixed(2)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Biggest Loss</span>
                <span className="font-medium text-red-600">
                  -${Math.abs(metrics?.biggest_loss || 0).toFixed(2)}
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
