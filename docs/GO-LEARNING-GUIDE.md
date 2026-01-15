# OddsIQ Backend - Go Learning Guide

**Purpose:** Learn Go practically by understanding the OddsIQ backend codebase step by step.

## 📚 Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Entry Points - Where to Start](#entry-points)
3. [Core Go Concepts Used](#core-go-concepts)
4. [Code Flow - How It All Works](#code-flow)
5. [Learning Path - What to Study First](#learning-path)
6. [Hands-On Exercises](#hands-on-exercises)

---

## Architecture Overview

### The Big Picture

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Request                         │
│                         ↓                               │
│  cmd/api/main.go → API Handlers → Services → Repos     │
│         ↓              ↓             ↓          ↓        │
│     Router         Business      Database    Models     │
│                     Logic        Access                  │
└─────────────────────────────────────────────────────────┘
```

### Folder Structure Explained

```
backend/
├── cmd/                    # Executables (Entry Points)
│   ├── api/               # Main API server
│   └── backfill/          # Data loading tool
│
├── internal/              # Private application code
│   ├── models/           # Data structures
│   ├── repository/       # Database operations
│   ├── services/         # Business logic
│   └── api/              # HTTP handlers
│
├── pkg/                   # Public/reusable packages
│   ├── apifootball/      # External API client
│   ├── oddsapi/          # External API client
│   └── database/         # Database connection
│
└── config/               # Configuration management
```

**Golden Rule:**
- `cmd/` = Programs you can run
- `internal/` = Your business code (can't be imported by other projects)
- `pkg/` = Libraries (can be imported by other projects)

---

## Entry Points - Where to Start

### 🔹 Entry Point #1: `cmd/api/main.go` (API Server)

**What it does:** Starts the HTTP server that handles requests.

**Read it line by line:**

```go
// File: backend/cmd/api/main.go

package main  // ← This makes it an executable program

import (
    "github.com/gin-gonic/gin"  // ← Web framework
    "github.com/dEnchanter/OddsIQ/backend/config"
    "github.com/dEnchanter/OddsIQ/backend/internal/api"
)

func main() {  // ← Program starts here
    // 1. Load configuration (.env file)
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Config failed: %v", err)
    }

    // 2. Connect to database
    db, err := database.New(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Database failed: %v", err)
    }
    defer db.Close()  // ← Close when program ends

    // 3. Create HTTP router
    router := gin.Default()

    // 4. Setup all routes (GET /api/fixtures, etc.)
    api.SetupRoutes(router, db.Pool, cfg)

    // 5. Start server on port 8000
    router.Run(":8000")
}
```

**Key Go Concepts Here:**
- `package main` + `func main()` = Executable program
- `import` statements
- Error handling with `if err != nil`
- `defer` = Run this when function exits
- Pointers (`*db`)

**Try This:**
1. Open `cmd/api/main.go`
2. Read each line slowly
3. Ask: "What does this line do?"

---

### 🔹 Entry Point #2: `cmd/backfill/main.go` (Data Loader)

**What it does:** Loads historical data into the database.

```go
// File: backend/cmd/backfill/main.go

func main() {
    // 1. Parse command-line flags
    seasonsFlag := flag.String("seasons", "2022,2023,2024", "Seasons to load")
    flag.Parse()

    // 2. Connect to database
    db, err := database.New(cfg.DatabaseURL)

    // 3. Create API client
    apiFootballClient := apifootball.NewClient(cfg.APIFootballKey)

    // 4. Create repositories (database access)
    teamsRepo := repository.NewTeamsRepository(db.Pool)
    fixturesRepo := repository.NewFixturesRepository(db.Pool)

    // 5. Create sync service (business logic)
    syncService := services.NewFixtureSyncService(
        apiFootballClient,
        teamsRepo,
        fixturesRepo,
    )

    // 6. Load data for each season
    for _, season := range seasons {
        syncService.SyncTeams(ctx, season)
        syncService.SyncFixturesBySeason(ctx, season)
    }
}
```

**Key Go Concepts Here:**
- Command-line flags (`flag` package)
- Loops (`for _, season := range seasons`)
- Dependency injection (passing objects to functions)
- Context (`ctx` for cancellation/timeouts)

---

## Core Go Concepts Used

### 1. Structs (Data Structures)

**What:** Groups of related data.

**Example from `internal/models/models.go`:**

```go
// Team represents a football team
type Team struct {
    ID            int       `json:"id"`           // ← Field name
    Name          string    `json:"name"`         // ← Field type
    Founded       int       `json:"founded"`
    CreatedAt     time.Time `json:"created_at"`
}
```

**Usage:**
```go
// Create a team
team := Team{
    ID:      1,
    Name:    "Arsenal",
    Founded: 1886,
}

// Access fields
fmt.Println(team.Name)  // Output: Arsenal
```

**Where we use it:** Models for Team, Fixture, Odds, etc.

---

### 2. Methods (Functions on Structs)

**What:** Functions that belong to a struct.

**Example from `internal/repository/teams.go`:**

```go
// TeamsRepository handles team database operations
type TeamsRepository struct {
    db *pgxpool.Pool  // ← Holds database connection
}

// Method: Belongs to TeamsRepository
func (r *TeamsRepository) GetByID(ctx context.Context, id int) (*Team, error) {
    //    ↑                                                        ↑
    // Receiver (like "self" in Python)                      Returns Team or error

    query := "SELECT id, name FROM teams WHERE id = $1"

    var team Team
    err := r.db.QueryRow(ctx, query, id).Scan(&team.ID, &team.Name)

    return &team, err
}
```

**Usage:**
```go
repo := NewTeamsRepository(db)
team, err := repo.GetByID(ctx, 5)  // Get team with ID 5
```

**Where we use it:** All repositories (teams, fixtures, odds)

---

### 3. Interfaces (Contracts)

**What:** Defines what methods a type must have.

**Example:**
```go
// Any type with Ping() method satisfies this interface
type Database interface {
    Ping(ctx context.Context) error
}

// PostgreSQL implementation
type PostgresDB struct {
    pool *pgxpool.Pool
}

func (p *PostgresDB) Ping(ctx context.Context) error {
    return p.pool.Ping(ctx)
}
```

**Why useful:** Can swap implementations without changing code.

**Where we use it:** Less explicit in our code, but Gin uses interfaces heavily.

---

### 4. Error Handling

**Go's way:** Functions return errors explicitly.

```go
// Bad: Crashes program
team := repo.GetByID(5)

// Good: Handle errors
team, err := repo.GetByID(ctx, 5)
if err != nil {
    log.Printf("Error getting team: %v", err)
    return err
}

// Or wrap errors with context
if err != nil {
    return fmt.Errorf("failed to get team %d: %w", id, err)
}
```

**Where we use it:** EVERYWHERE. Every function that can fail returns `error`.

---

### 5. Pointers

**What:** Reference to memory location.

```go
// Value: Creates a copy
team := Team{Name: "Arsenal"}
updateTeam(team)  // Changes won't persist

// Pointer: Passes reference
team := &Team{Name: "Arsenal"}
updateTeam(team)  // Changes persist
```

**When to use:**
- `*Type` = Pointer (can be nil)
- `Type` = Value (never nil)

**Example from our code:**
```go
func (r *TeamsRepository) GetByID(ctx, id) (*Team, error) {
    //                                        ↑
    //                               Returns pointer to Team

    var team Team
    // ... query database ...
    return &team, nil  // ← & creates pointer
}
```

---

### 6. Packages and Imports

**What:** Organize code into modules.

```go
// Declare package
package repository  // ← This file belongs to "repository" package

// Import other packages
import (
    "fmt"                    // ← Standard library
    "github.com/gin-gonic/gin"  // ← External package
    "github.com/dEnchanter/OddsIQ/backend/internal/models"  // ← Our package
)
```

**Rules:**
- Lowercase package name = private (can't import from outside)
- Uppercase function name = public (can be called from outside)

```go
package example

func PublicFunction() {}   // ← Can be called: example.PublicFunction()
func privateFunction() {}  // ← Can't be called from outside package
```

---

## Code Flow - How It All Works

### Example: GET /api/fixtures?season=2024

**Step-by-step flow:**

#### 1. Request arrives at Router

```go
// File: internal/api/routes.go

fixtures.GET("", api.getFixtures())
//            ↑         ↑
//          Path     Handler function
```

#### 2. Handler processes request

```go
// File: internal/api/handlers.go

func (api *API) getFixtures() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get query parameter
        seasonStr := c.Query("season")  // "2024"
        season, _ := strconv.Atoi(seasonStr)  // Convert to int

        // Call repository
        fixtures, err := api.fixturesRepo.GetBySeason(ctx, season)

        // Return JSON response
        c.JSON(http.StatusOK, gin.H{
            "fixtures": fixtures,
            "total":    len(fixtures),
        })
    }
}
```

#### 3. Repository queries database

```go
// File: internal/repository/fixtures.go

func (r *FixturesRepository) GetBySeason(ctx, season int) ([]Fixture, error) {
    query := `
        SELECT id, home_team_id, away_team_id, match_date
        FROM fixtures
        WHERE season = $1
        ORDER BY match_date
    `

    rows, err := r.db.Query(ctx, query, season)
    // ... scan rows into fixtures ...

    return fixtures, nil
}
```

#### 4. Response sent to client

```json
{
    "fixtures": [
        {"id": 1, "home_team_id": 5, "away_team_id": 10},
        {"id": 2, "home_team_id": 3, "away_team_id": 8}
    ],
    "total": 2
}
```

### Visual Flow

```
HTTP GET /api/fixtures?season=2024
         ↓
    routes.go (matches route)
         ↓
    handlers.go (getFixtures)
         ↓
    repository/fixtures.go (GetBySeason)
         ↓
    PostgreSQL Database
         ↓
    JSON Response
```

---

## Learning Path - What to Study First

### Week 1: Foundations

**Day 1-2: Understand the Entry Point**
- [ ] Read `cmd/api/main.go` line by line
- [ ] Understand: `package main`, `func main()`, `import`
- [ ] Run the server: `go run cmd/api/main.go`

**Day 3-4: Models (Data Structures)**
- [ ] Read `internal/models/models.go`
- [ ] Understand: `type`, `struct`, fields, tags
- [ ] Try: Create your own struct

```go
// Try this in a new file
type Player struct {
    Name     string
    Position string
    Goals    int
}

player := Player{
    Name:     "Saka",
    Position: "Forward",
    Goals:    20,
}
fmt.Println(player.Name)
```

**Day 5-7: Configuration**
- [ ] Read `config/config.go`
- [ ] Understand: Environment variables, `godotenv`
- [ ] Create your own `.env` file

---

### Week 2: Database Layer

**Day 1-3: Database Connection**
- [ ] Read `pkg/database/database.go`
- [ ] Understand: Connection pooling, `pgxpool`
- [ ] Try: Connect to database

```go
db, err := database.New("postgresql://...")
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

**Day 4-7: Repositories**
- [ ] Read `internal/repository/teams.go`
- [ ] Understand: Methods, receivers, SQL queries
- [ ] Try: Add a new method

```go
// Add to teams.go
func (r *TeamsRepository) GetByName(ctx context.Context, name string) (*models.Team, error) {
    query := "SELECT id, name FROM teams WHERE name = $1"

    var team models.Team
    err := r.db.QueryRow(ctx, query, name).Scan(&team.ID, &team.Name)

    return &team, err
}
```

---

### Week 3: Business Logic

**Day 1-4: Services**
- [ ] Read `internal/services/fixture_sync.go`
- [ ] Understand: How services orchestrate multiple operations
- [ ] Trace: Follow `SyncTeams()` method

**Day 5-7: API Clients**
- [ ] Read `pkg/apifootball/client.go`
- [ ] Understand: HTTP requests, JSON parsing
- [ ] Try: Make your own API call

---

### Week 4: HTTP Layer

**Day 1-3: Handlers**
- [ ] Read `internal/api/handlers.go`
- [ ] Understand: Gin context, JSON responses
- [ ] Try: Add a new endpoint

**Day 4-7: Routes**
- [ ] Read `internal/api/routes.go`
- [ ] Understand: Route groups, middleware
- [ ] Build: Complete request-response cycle

---

## Hands-On Exercises

### Exercise 1: Add a New Model Field

**Goal:** Add `nickname` field to Team.

```go
// 1. Update model (internal/models/models.go)
type Team struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Nickname string `json:"nickname"`  // ← Add this
}

// 2. Update database migration
ALTER TABLE teams ADD COLUMN nickname VARCHAR(50);

// 3. Update repository (internal/repository/teams.go)
// Add nickname to all SQL queries

// 4. Test it!
```

---

### Exercise 2: Create a New Endpoint

**Goal:** Add `GET /api/teams` endpoint.

```go
// 1. Add handler (internal/api/handlers.go)
func (api *API) getTeams() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()

        teams, err := api.teamsRepo.GetAll(ctx)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }

        c.JSON(200, gin.H{
            "teams": teams,
            "total": len(teams),
        })
    }
}

// 2. Add route (internal/api/routes.go)
teams := v1.Group("/teams")
{
    teams.GET("", api.getTeams())
}

// 3. Test: curl http://localhost:8000/api/teams
```

---

### Exercise 3: Add Query Parameters

**Goal:** Filter teams by founded year.

```go
// 1. Update handler
func (api *API) getTeams() gin.HandlerFunc {
    return func(c *gin.Context) {
        foundedStr := c.Query("founded")  // Get ?founded=1886

        if foundedStr != "" {
            founded, _ := strconv.Atoi(foundedStr)
            teams, _ := api.teamsRepo.GetByFounded(ctx, founded)
            c.JSON(200, gin.H{"teams": teams})
            return
        }

        // Default: all teams
        teams, _ := api.teamsRepo.GetAll(ctx)
        c.JSON(200, gin.H{"teams": teams})
    }
}

// 2. Add repository method
func (r *TeamsRepository) GetByFounded(ctx context.Context, founded int) ([]models.Team, error) {
    query := "SELECT * FROM teams WHERE founded = $1"
    rows, err := r.db.Query(ctx, query, founded)
    // ... scan rows ...
    return teams, err
}

// 3. Test: curl http://localhost:8000/api/teams?founded=1886
```

---

## Common Go Patterns in Our Code

### Pattern 1: Constructor Functions

```go
// Instead of: repo := &TeamsRepository{db: pool}
// We use:
func NewTeamsRepository(db *pgxpool.Pool) *TeamsRepository {
    return &TeamsRepository{db: db}
}

// Why: Encapsulates initialization logic
```

### Pattern 2: Error Wrapping

```go
// Bad
if err != nil {
    return err
}

// Good
if err != nil {
    return fmt.Errorf("failed to sync teams for season %d: %w", season, err)
}

// Why: Adds context to errors
```

### Pattern 3: Defer for Cleanup

```go
func ProcessFile() error {
    file, err := os.Open("data.txt")
    if err != nil {
        return err
    }
    defer file.Close()  // ← Always closes, even if panic

    // ... process file ...
    return nil
}
```

### Pattern 4: Context for Cancellation

```go
func LongRunningTask(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():  // ← Check if cancelled
            return ctx.Err()
        default:
            // Do work
        }
    }
}
```

---

## Debugging Tips

### 1. Add Print Statements

```go
fmt.Printf("DEBUG: season = %d\n", season)
log.Printf("Got %d fixtures from API", len(fixtures))
```

### 2. Use Go's Built-in Debugger

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug your program
dlv debug cmd/api/main.go
```

### 3. Check Error Messages

```go
if err != nil {
    log.Printf("ERROR: %+v", err)  // ← Prints full error trace
    return err
}
```

---

## Next Steps

**Once you're comfortable:**

1. **Read official Go tour:** https://go.dev/tour/
2. **Effective Go:** https://go.dev/doc/effective_go
3. **Add features to OddsIQ:**
   - New endpoints
   - New data models
   - New business logic

**Remember:**
- Start small (read one file at a time)
- Run the code (don't just read)
- Break things (best way to learn)
- Ask questions (use comments or ChatGPT)

---

## Quick Reference

### Common Go Commands

```bash
go run cmd/api/main.go       # Run program
go build cmd/api             # Compile binary
go test ./...                # Run tests
go mod tidy                  # Clean dependencies
go fmt ./...                 # Format code
```

### Project-Specific Commands

```bash
# Start API server
cd backend
go run cmd/api/main.go

# Run backfill
go run cmd/backfill/main.go -seasons 2024

# Build binaries
go build -o bin/api.exe ./cmd/api
go build -o bin/backfill.exe ./cmd/backfill
```

---

**You've got this!** Start with `cmd/api/main.go` and work your way through. Learning by doing is the best approach! 🚀
