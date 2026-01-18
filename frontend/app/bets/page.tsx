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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Bet, EnrichedFixture } from '@/types';
import { getBets, createBet, settleBet, getUpcomingFixtures } from '@/lib/api';

export default function BetsPage() {
  const [bets, setBets] = useState<Bet[]>([]);
  const [fixtures, setFixtures] = useState<EnrichedFixture[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [dialogOpen, setDialogOpen] = useState(false);

  // New bet form state
  const [newBetFixture, setNewBetFixture] = useState<string>('');
  const [newBetType, setNewBetType] = useState<string>('');
  const [newBetOdds, setNewBetOdds] = useState<string>('');
  const [newBetStake, setNewBetStake] = useState<string>('');
  const [newBetBookmaker, setNewBetBookmaker] = useState<string>('');
  const [submitting, setSubmitting] = useState(false);

  // Load bets and fixtures
  const loadData = async () => {
    try {
      setLoading(true);
      const [betsRes, fixturesRes] = await Promise.all([
        getBets(),
        getUpcomingFixtures(),
      ]);
      setBets(betsRes.bets || []);
      setFixtures(fixturesRes.fixtures || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  // Filter bets by status
  const filteredBets = bets.filter((bet) => {
    if (statusFilter === 'all') return true;
    return bet.status === statusFilter;
  });

  // Calculate summary stats
  const pendingBets = bets.filter((b) => b.status === 'pending');
  const wonBets = bets.filter((b) => b.status === 'won');
  const lostBets = bets.filter((b) => b.status === 'lost');
  const totalStaked = bets.reduce((sum, b) => sum + b.stake, 0);
  const totalProfit = bets.reduce((sum, b) => sum + (b.profit_loss || 0), 0);

  // Handle create bet
  const handleCreateBet = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    if (!newBetFixture || !newBetType || !newBetOdds || !newBetStake || !newBetBookmaker) {
      setError('Please fill in all fields');
      return;
    }

    try {
      setSubmitting(true);
      await createBet({
        fixture_id: parseInt(newBetFixture),
        bet_type: newBetType,
        odds: parseFloat(newBetOdds),
        stake: parseFloat(newBetStake),
        expected_value: 0, // Will be calculated server-side
        bookmaker: newBetBookmaker,
      });

      setSuccess('Bet recorded successfully');
      setDialogOpen(false);

      // Reset form
      setNewBetFixture('');
      setNewBetType('');
      setNewBetOdds('');
      setNewBetStake('');
      setNewBetBookmaker('');

      // Reload bets
      loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create bet');
    } finally {
      setSubmitting(false);
    }
  };

  // Handle settle bet
  const handleSettleBet = async (betId: number, result: 'won' | 'lost') => {
    try {
      await settleBet(betId, { result });
      setSuccess(`Bet marked as ${result}`);
      loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to settle bet');
    }
  };

  // Format date
  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('en-GB', {
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  // Get status badge variant
  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'won':
        return <Badge className="bg-green-500">Won</Badge>;
      case 'lost':
        return <Badge variant="destructive">Lost</Badge>;
      default:
        return <Badge variant="secondary">Pending</Badge>;
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">Loading bets...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Bet Tracker</h1>
          <p className="text-muted-foreground">Record and track your bets</p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button>+ Record Bet</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Record New Bet</DialogTitle>
              <DialogDescription>Enter the details of your placed bet</DialogDescription>
            </DialogHeader>
            <form onSubmit={handleCreateBet} className="space-y-4">
              <div className="space-y-2">
                <Label>Fixture</Label>
                <Select value={newBetFixture} onValueChange={setNewBetFixture}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select fixture" />
                  </SelectTrigger>
                  <SelectContent>
                    {fixtures.map((fixture) => (
                      <SelectItem key={fixture.id} value={fixture.id.toString()}>
                        {fixture.home_team_name} vs {fixture.away_team_name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label>Bet Type</Label>
                <Select value={newBetType} onValueChange={setNewBetType}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select bet type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="home_win">Home Win</SelectItem>
                    <SelectItem value="draw">Draw</SelectItem>
                    <SelectItem value="away_win">Away Win</SelectItem>
                    <SelectItem value="over_2_5">Over 2.5 Goals</SelectItem>
                    <SelectItem value="under_2_5">Under 2.5 Goals</SelectItem>
                    <SelectItem value="btts_yes">BTTS Yes</SelectItem>
                    <SelectItem value="btts_no">BTTS No</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Odds</Label>
                  <Input
                    type="number"
                    step="0.01"
                    placeholder="1.85"
                    value={newBetOdds}
                    onChange={(e) => setNewBetOdds(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label>Stake ($)</Label>
                  <Input
                    type="number"
                    step="0.01"
                    placeholder="25.00"
                    value={newBetStake}
                    onChange={(e) => setNewBetStake(e.target.value)}
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label>Bookmaker</Label>
                <Input
                  placeholder="e.g., Bet365"
                  value={newBetBookmaker}
                  onChange={(e) => setNewBetBookmaker(e.target.value)}
                />
              </div>

              <Button type="submit" className="w-full" disabled={submitting}>
                {submitting ? 'Recording...' : 'Record Bet'}
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {success && (
        <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">
          {success}
        </div>
      )}

      {/* Summary Cards */}
      <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
        <Card>
          <CardContent className="pt-6 text-center">
            <p className="text-2xl font-bold">{bets.length}</p>
            <p className="text-sm text-muted-foreground">Total Bets</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6 text-center">
            <p className="text-2xl font-bold">{pendingBets.length}</p>
            <p className="text-sm text-muted-foreground">Pending</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6 text-center">
            <p className="text-2xl font-bold text-green-600">{wonBets.length}</p>
            <p className="text-sm text-muted-foreground">Won</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6 text-center">
            <p className="text-2xl font-bold text-red-600">{lostBets.length}</p>
            <p className="text-sm text-muted-foreground">Lost</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6 text-center">
            <p className={`text-2xl font-bold ${totalProfit >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              {totalProfit >= 0 ? '+' : ''}${totalProfit.toFixed(2)}
            </p>
            <p className="text-sm text-muted-foreground">Profit/Loss</p>
          </CardContent>
        </Card>
      </div>

      {/* Filter and Table */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Bet History</CardTitle>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-32">
                <SelectValue placeholder="Filter" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="won">Won</SelectItem>
                <SelectItem value="lost">Lost</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          {filteredBets.length === 0 ? (
            <p className="text-muted-foreground text-center py-8">
              No bets recorded yet. Click &quot;Record Bet&quot; to add your first bet.
            </p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Date</TableHead>
                  <TableHead>Fixture</TableHead>
                  <TableHead>Bet Type</TableHead>
                  <TableHead>Odds</TableHead>
                  <TableHead>Stake</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Profit</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredBets.map((bet) => (
                  <TableRow key={bet.id}>
                    <TableCell>{formatDate(bet.placed_at)}</TableCell>
                    <TableCell className="font-medium">
                      Fixture #{bet.fixture_id}
                    </TableCell>
                    <TableCell>{bet.bet_type}</TableCell>
                    <TableCell>{bet.odds.toFixed(2)}</TableCell>
                    <TableCell>${bet.stake.toFixed(2)}</TableCell>
                    <TableCell>{getStatusBadge(bet.status)}</TableCell>
                    <TableCell>
                      {bet.profit_loss !== null && (
                        <span className={bet.profit_loss >= 0 ? 'text-green-600' : 'text-red-600'}>
                          {bet.profit_loss >= 0 ? '+' : ''}${bet.profit_loss.toFixed(2)}
                        </span>
                      )}
                    </TableCell>
                    <TableCell className="text-right">
                      {bet.status === 'pending' && (
                        <div className="flex gap-2 justify-end">
                          <Button
                            size="sm"
                            variant="outline"
                            className="text-green-600 border-green-600"
                            onClick={() => handleSettleBet(bet.id, 'won')}
                          >
                            Won
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            className="text-red-600 border-red-600"
                            onClick={() => handleSettleBet(bet.id, 'lost')}
                          >
                            Lost
                          </Button>
                        </div>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
