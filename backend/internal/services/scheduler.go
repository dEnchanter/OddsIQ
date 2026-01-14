package services

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// Scheduler manages scheduled tasks for data synchronization
type Scheduler struct {
	cron              *cron.Cron
	fixtureSyncService *FixtureSyncService
	oddsSyncService    *OddsSyncService
}

// NewScheduler creates a new scheduler
func NewScheduler(
	fixtureSyncService *FixtureSyncService,
	oddsSyncService *OddsSyncService,
) *Scheduler {
	// Create cron with second precision
	c := cron.New(cron.WithSeconds())

	return &Scheduler{
		cron:              c,
		fixtureSyncService: fixtureSyncService,
		oddsSyncService:    oddsSyncService,
	}
}

// Start starts the scheduler and all jobs
func (s *Scheduler) Start() error {
	log.Println("Starting scheduler...")

	ctx := context.Background()

	// Job 1: Sync upcoming fixtures daily at 6:00 AM
	_, err := s.cron.AddFunc("0 0 6 * * *", func() {
		log.Println("Running scheduled job: Sync upcoming fixtures")
		if err := s.fixtureSyncService.SyncUpcomingFixtures(ctx); err != nil {
			log.Printf("Error syncing upcoming fixtures: %v", err)
		}
	})
	if err != nil {
		return err
	}

	// Job 2: Update fixture results every 30 minutes during match days
	_, err = s.cron.AddFunc("0 */30 * * * *", func() {
		// Only run on match days (Friday-Monday)
		now := time.Now()
		weekday := now.Weekday()
		if weekday >= time.Friday || weekday <= time.Monday {
			log.Println("Running scheduled job: Update fixture results")
			if err := s.fixtureSyncService.UpdateFixtureResults(ctx); err != nil {
				log.Printf("Error updating fixture results: %v", err)
			}
		}
	})
	if err != nil {
		return err
	}

	// Job 3: Sync odds every 2 hours
	_, err = s.cron.AddFunc("0 0 */2 * * *", func() {
		log.Println("Running scheduled job: Sync odds for all markets")
		if err := s.oddsSyncService.SyncAllMarkets(ctx); err != nil {
			log.Printf("Error syncing odds: %v", err)
		}
	})
	if err != nil {
		return err
	}

	// Job 4: Sync H2H odds every hour (more frequent for main market)
	_, err = s.cron.AddFunc("0 0 * * * *", func() {
		log.Println("Running scheduled job: Sync H2H odds")
		if err := s.oddsSyncService.SyncH2HOdds(ctx); err != nil {
			log.Printf("Error syncing H2H odds: %v", err)
		}
	})
	if err != nil {
		return err
	}

	// Job 5: Cleanup old odds weekly (Sunday at 3:00 AM)
	_, err = s.cron.AddFunc("0 0 3 * * 0", func() {
		log.Println("Running scheduled job: Cleanup old odds")
		if err := s.oddsSyncService.CleanupOldOdds(ctx, 30); err != nil {
			log.Printf("Error cleaning up old odds: %v", err)
		}
	})
	if err != nil {
		return err
	}

	// Start the cron scheduler
	s.cron.Start()
	log.Println("Scheduler started successfully")

	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

// RunNow executes all jobs immediately (useful for testing)
func (s *Scheduler) RunNow() error {
	ctx := context.Background()

	log.Println("Running all jobs immediately...")

	// Sync upcoming fixtures
	log.Println("1/4: Syncing upcoming fixtures...")
	if err := s.fixtureSyncService.SyncUpcomingFixtures(ctx); err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Println("✓ Upcoming fixtures synced")
	}

	// Update fixture results
	log.Println("2/4: Updating fixture results...")
	if err := s.fixtureSyncService.UpdateFixtureResults(ctx); err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Println("✓ Fixture results updated")
	}

	// Sync odds
	log.Println("3/4: Syncing odds for all markets...")
	if err := s.oddsSyncService.SyncAllMarkets(ctx); err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Println("✓ Odds synced")
	}

	// Print summary
	log.Println("4/4: Getting odds summary...")
	summary, err := s.oddsSyncService.GetOddsSummary(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("✓ Odds summary: %+v", summary)
	}

	log.Println("All jobs completed")
	return nil
}

// GetNextRunTimes returns the next run times for all scheduled jobs
func (s *Scheduler) GetNextRunTimes() []time.Time {
	entries := s.cron.Entries()
	nextRuns := make([]time.Time, len(entries))
	for i, entry := range entries {
		nextRuns[i] = entry.Next
	}
	return nextRuns
}

// Custom job schedules for different environments

// StartDevelopmentSchedule starts a development-friendly schedule (less frequent)
func (s *Scheduler) StartDevelopmentSchedule() error {
	log.Println("Starting development scheduler...")

	ctx := context.Background()

	// Sync fixtures once per day at noon
	_, err := s.cron.AddFunc("0 0 12 * * *", func() {
		log.Println("[DEV] Syncing upcoming fixtures")
		if err := s.fixtureSyncService.SyncUpcomingFixtures(ctx); err != nil {
			log.Printf("Error: %v", err)
		}
	})
	if err != nil {
		return err
	}

	// Sync odds twice per day
	_, err = s.cron.AddFunc("0 0 10,18 * * *", func() {
		log.Println("[DEV] Syncing odds")
		if err := s.oddsSyncService.SyncAllMarkets(ctx); err != nil {
			log.Printf("Error: %v", err)
		}
	})
	if err != nil {
		return err
	}

	s.cron.Start()
	log.Println("Development scheduler started")

	return nil
}

// StartProductionSchedule starts a production schedule (more frequent, optimized)
func (s *Scheduler) StartProductionSchedule() error {
	return s.Start() // Use the default production schedule
}
