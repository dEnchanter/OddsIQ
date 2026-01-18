'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
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
import { Badge } from '@/components/ui/badge';
import { Team, EnrichedFixture } from '@/types';
import {
  getTeams,
  getUpcomingFixtures,
  createManualFixture,
  addOddsBatch,
  deleteFixture,
} from '@/lib/api';

export default function EntryPage() {
  // State for teams dropdown
  const [teams, setTeams] = useState<Team[]>([]);
  const [fixtures, setFixtures] = useState<EnrichedFixture[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // Fixture form state
  const [homeTeamId, setHomeTeamId] = useState<string>('');
  const [awayTeamId, setAwayTeamId] = useState<string>('');
  const [matchDate, setMatchDate] = useState<string>('');
  const [matchTime, setMatchTime] = useState<string>('15:00');
  const [round, setRound] = useState<string>('');
  const [fixtureSubmitting, setFixtureSubmitting] = useState(false);

  // Odds form state
  const [selectedFixtureId, setSelectedFixtureId] = useState<string>('');
  const [bookmaker, setBookmaker] = useState<string>('');
  const [homeOdds, setHomeOdds] = useState<string>('');
  const [drawOdds, setDrawOdds] = useState<string>('');
  const [awayOdds, setAwayOdds] = useState<string>('');
  const [overOdds, setOverOdds] = useState<string>('');
  const [underOdds, setUnderOdds] = useState<string>('');
  const [bttsYesOdds, setBttsYesOdds] = useState<string>('');
  const [bttsNoOdds, setBttsNoOdds] = useState<string>('');
  const [oddsSubmitting, setOddsSubmitting] = useState(false);

  // Load teams and fixtures on mount
  useEffect(() => {
    async function loadData() {
      try {
        setLoading(true);
        const [teamsRes, fixturesRes] = await Promise.all([
          getTeams(),
          getUpcomingFixtures(),
        ]);
        setTeams(teamsRes.teams || []);
        setFixtures(fixturesRes.fixtures || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load data');
      } finally {
        setLoading(false);
      }
    }
    loadData();
  }, []);

  // Get current season
  const getCurrentSeason = () => {
    const now = new Date();
    return now.getMonth() >= 7 ? now.getFullYear() : now.getFullYear() - 1;
  };

  // Handle fixture creation
  const handleCreateFixture = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    if (!homeTeamId || !awayTeamId || !matchDate) {
      setError('Please fill in all required fields');
      return;
    }

    if (homeTeamId === awayTeamId) {
      setError('Home and away teams must be different');
      return;
    }

    try {
      setFixtureSubmitting(true);
      const dateTime = `${matchDate}T${matchTime}:00Z`;

      const result = await createManualFixture({
        home_team_id: parseInt(homeTeamId),
        away_team_id: parseInt(awayTeamId),
        match_date: dateTime,
        season: getCurrentSeason(),
        round: round || undefined,
      });

      setSuccess(`Fixture created: ${result.home_team.name} vs ${result.away_team.name}`);

      // Reset form
      setHomeTeamId('');
      setAwayTeamId('');
      setMatchDate('');
      setRound('');

      // Reload fixtures
      const fixturesRes = await getUpcomingFixtures();
      setFixtures(fixturesRes.fixtures || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create fixture');
    } finally {
      setFixtureSubmitting(false);
    }
  };

  // Handle odds submission
  const handleAddOdds = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    if (!selectedFixtureId || !bookmaker) {
      setError('Please select a fixture and enter bookmaker name');
      return;
    }

    const oddsEntries = [];

    // Add 1X2 odds
    if (homeOdds) oddsEntries.push({ market_type: 'h2h', outcome: 'Home', odds_value: parseFloat(homeOdds) });
    if (drawOdds) oddsEntries.push({ market_type: 'h2h', outcome: 'Draw', odds_value: parseFloat(drawOdds) });
    if (awayOdds) oddsEntries.push({ market_type: 'h2h', outcome: 'Away', odds_value: parseFloat(awayOdds) });

    // Add Over/Under odds
    if (overOdds) oddsEntries.push({ market_type: 'totals', outcome: 'Over', odds_value: parseFloat(overOdds) });
    if (underOdds) oddsEntries.push({ market_type: 'totals', outcome: 'Under', odds_value: parseFloat(underOdds) });

    // Add BTTS odds
    if (bttsYesOdds) oddsEntries.push({ market_type: 'btts', outcome: 'Yes', odds_value: parseFloat(bttsYesOdds) });
    if (bttsNoOdds) oddsEntries.push({ market_type: 'btts', outcome: 'No', odds_value: parseFloat(bttsNoOdds) });

    if (oddsEntries.length === 0) {
      setError('Please enter at least one odds value');
      return;
    }

    // Validate all odds are > 1
    for (const entry of oddsEntries) {
      if (entry.odds_value <= 1) {
        setError('All odds must be greater than 1.0');
        return;
      }
    }

    try {
      setOddsSubmitting(true);
      await addOddsBatch({
        fixture_id: parseInt(selectedFixtureId),
        bookmaker,
        odds: oddsEntries,
      });

      setSuccess(`Added ${oddsEntries.length} odds entries`);

      // Reset odds form
      setHomeOdds('');
      setDrawOdds('');
      setAwayOdds('');
      setOverOdds('');
      setUnderOdds('');
      setBttsYesOdds('');
      setBttsNoOdds('');

      // Reload fixtures
      const fixturesRes = await getUpcomingFixtures();
      setFixtures(fixturesRes.fixtures || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add odds');
    } finally {
      setOddsSubmitting(false);
    }
  };

  // Handle fixture deletion
  const handleDeleteFixture = async (id: number) => {
    if (!confirm('Are you sure you want to delete this fixture?')) return;

    try {
      await deleteFixture(id);
      setSuccess('Fixture deleted');
      const fixturesRes = await getUpcomingFixtures();
      setFixtures(fixturesRes.fixtures || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete fixture');
    }
  };

  // Format date for display
  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('en-GB', {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Manual Entry</h1>
        <p className="text-muted-foreground">Add fixtures and odds for upcoming matches</p>
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

      <div className="grid gap-6 md:grid-cols-2">
        {/* Add Fixture Card */}
        <Card>
          <CardHeader>
            <CardTitle>Add New Fixture</CardTitle>
            <CardDescription>Create a fixture for an upcoming match</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreateFixture} className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="homeTeam">Home Team</Label>
                  <Select value={homeTeamId} onValueChange={setHomeTeamId}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select home team" />
                    </SelectTrigger>
                    <SelectContent>
                      {teams.map((team) => (
                        <SelectItem key={team.id} value={team.id.toString()}>
                          {team.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="awayTeam">Away Team</Label>
                  <Select value={awayTeamId} onValueChange={setAwayTeamId}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select away team" />
                    </SelectTrigger>
                    <SelectContent>
                      {teams.map((team) => (
                        <SelectItem key={team.id} value={team.id.toString()}>
                          {team.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="matchDate">Match Date</Label>
                  <Input
                    id="matchDate"
                    type="date"
                    value={matchDate}
                    onChange={(e) => setMatchDate(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="matchTime">Kick-off Time</Label>
                  <Input
                    id="matchTime"
                    type="time"
                    value={matchTime}
                    onChange={(e) => setMatchTime(e.target.value)}
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="round">Round (optional)</Label>
                <Input
                  id="round"
                  placeholder="e.g., Matchweek 22"
                  value={round}
                  onChange={(e) => setRound(e.target.value)}
                />
              </div>

              <Button type="submit" className="w-full" disabled={fixtureSubmitting}>
                {fixtureSubmitting ? 'Creating...' : 'Create Fixture'}
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Add Odds Card */}
        <Card>
          <CardHeader>
            <CardTitle>Add Odds</CardTitle>
            <CardDescription>Enter odds from your bookmaker</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleAddOdds} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="fixture">Select Fixture</Label>
                <Select value={selectedFixtureId} onValueChange={setSelectedFixtureId}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select a fixture" />
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
                <Label htmlFor="bookmaker">Bookmaker</Label>
                <Input
                  id="bookmaker"
                  placeholder="e.g., Bet365"
                  value={bookmaker}
                  onChange={(e) => setBookmaker(e.target.value)}
                />
              </div>

              <div className="space-y-3">
                <Label>1X2 Market</Label>
                <div className="grid grid-cols-3 gap-2">
                  <div>
                    <Label className="text-xs text-muted-foreground">Home</Label>
                    <Input
                      type="number"
                      step="0.01"
                      placeholder="1.85"
                      value={homeOdds}
                      onChange={(e) => setHomeOdds(e.target.value)}
                    />
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">Draw</Label>
                    <Input
                      type="number"
                      step="0.01"
                      placeholder="3.60"
                      value={drawOdds}
                      onChange={(e) => setDrawOdds(e.target.value)}
                    />
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">Away</Label>
                    <Input
                      type="number"
                      step="0.01"
                      placeholder="4.20"
                      value={awayOdds}
                      onChange={(e) => setAwayOdds(e.target.value)}
                    />
                  </div>
                </div>
              </div>

              <div className="space-y-3">
                <Label>Over/Under 2.5</Label>
                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <Label className="text-xs text-muted-foreground">Over</Label>
                    <Input
                      type="number"
                      step="0.01"
                      placeholder="1.90"
                      value={overOdds}
                      onChange={(e) => setOverOdds(e.target.value)}
                    />
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">Under</Label>
                    <Input
                      type="number"
                      step="0.01"
                      placeholder="1.95"
                      value={underOdds}
                      onChange={(e) => setUnderOdds(e.target.value)}
                    />
                  </div>
                </div>
              </div>

              <div className="space-y-3">
                <Label>Both Teams to Score</Label>
                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <Label className="text-xs text-muted-foreground">Yes</Label>
                    <Input
                      type="number"
                      step="0.01"
                      placeholder="1.80"
                      value={bttsYesOdds}
                      onChange={(e) => setBttsYesOdds(e.target.value)}
                    />
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">No</Label>
                    <Input
                      type="number"
                      step="0.01"
                      placeholder="1.95"
                      value={bttsNoOdds}
                      onChange={(e) => setBttsNoOdds(e.target.value)}
                    />
                  </div>
                </div>
              </div>

              <Button type="submit" className="w-full" disabled={oddsSubmitting}>
                {oddsSubmitting ? 'Adding...' : 'Add All Odds'}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>

      {/* Upcoming Fixtures Table */}
      <Card>
        <CardHeader>
          <CardTitle>Upcoming Fixtures</CardTitle>
          <CardDescription>Fixtures ready for predictions</CardDescription>
        </CardHeader>
        <CardContent>
          {fixtures.length === 0 ? (
            <p className="text-muted-foreground text-center py-8">
              No upcoming fixtures. Add a fixture above to get started.
            </p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Fixture</TableHead>
                  <TableHead>Date</TableHead>
                  <TableHead>Odds</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {fixtures.map((fixture) => (
                  <TableRow key={fixture.id}>
                    <TableCell className="font-medium">
                      {fixture.home_team_name} vs {fixture.away_team_name}
                    </TableCell>
                    <TableCell>{formatDate(fixture.match_date)}</TableCell>
                    <TableCell>
                      <span className="text-sm">
                        {fixture.odds_count}/7
                      </span>
                    </TableCell>
                    <TableCell>
                      {fixture.has_odds ? (
                        <Badge variant="default">Ready</Badge>
                      ) : (
                        <Badge variant="secondary">No Odds</Badge>
                      )}
                    </TableCell>
                    <TableCell className="text-right">
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDeleteFixture(fixture.id)}
                      >
                        Delete
                      </Button>
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
