# OddsIQ MVP - Setup Guide

## Prerequisites

### Required Software
- **Go** 1.21 or higher
- **Python** 3.11 or higher
- **Node.js** 18 or higher
- **PostgreSQL** 15 or higher
- **Git**

### Recommended Tools
- **Docker** & Docker Compose (optional, for containerized PostgreSQL)
- **VS Code** or **GoLand** for Go development
- **PyCharm** or **VS Code** for Python development
- **TablePlus** or **pgAdmin** for database management

## Initial Setup Steps

### 1. Clone Repository
```bash
git clone <repository-url>
cd OddsIQ
```

### 2. Database Setup

#### Option A: Local PostgreSQL
```bash
# Install PostgreSQL 15+
# macOS
brew install postgresql@15

# Ubuntu/Debian
sudo apt install postgresql-15

# Windows
# Download from https://www.postgresql.org/download/windows/

# Start PostgreSQL
# macOS
brew services start postgresql@15

# Ubuntu/Debian
sudo systemctl start postgresql

# Create database and user
psql postgres
CREATE DATABASE oddsiq;
CREATE USER oddsiq_user WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE oddsiq TO oddsiq_user;
\q
```

#### Option B: Docker PostgreSQL
```bash
# Using docker-compose (create this file in project root)
docker-compose up -d postgres
```

Example `docker-compose.yml`:
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: oddsiq-postgres
    environment:
      POSTGRES_DB: oddsiq
      POSTGRES_USER: oddsiq_user
      POSTGRES_PASSWORD: your_secure_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./database/migrations:/docker-entrypoint-initdb.d

volumes:
  postgres_data:
```

### 3. Backend (Go) Setup

```bash
cd backend

# Initialize Go module (if not already done)
go mod init github.com/yourusername/oddsiq-backend

# Install dependencies (after creating initial files)
go mod tidy

# Copy environment template
cp .env.example .env

# Edit .env with your configuration
# DATABASE_URL=postgres://oddsiq_user:your_password@localhost:5432/oddsiq
# API_FOOTBALL_KEY=your_api_football_key
# ODDS_API_KEY=your_odds_api_key
# ML_SERVICE_URL=http://localhost:8001

# Run migrations
go run cmd/migrate/main.go up

# Run the server
go run cmd/api/main.go
```

Backend should now be running on `http://localhost:8000`

### 4. ML Service (Python) Setup

```bash
cd ml-service

# Create virtual environment
python -m venv venv

# Activate virtual environment
# macOS/Linux
source venv/bin/activate

# Windows
venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Or using poetry
poetry install

# Copy environment template
cp .env.example .env

# Edit .env with your configuration
# DATABASE_URL=postgresql://oddsiq_user:your_password@localhost:5432/oddsiq
# MODEL_PATH=./models
# FEATURE_STORE_PATH=./feature_store

# Run migrations (if separate from backend)
# alembic upgrade head

# Start the ML service
uvicorn app.main:app --host 0.0.0.0 --port 8001 --reload
```

ML Service should now be running on `http://localhost:8001`

### 5. Frontend (Next.js) Setup

```bash
cd frontend

# Install dependencies
npm install
# or
pnpm install
# or
yarn install

# Copy environment template
cp .env.local.example .env.local

# Edit .env.local with your configuration
# NEXT_PUBLIC_API_URL=http://localhost:8000/api

# Run development server
npm run dev
```

Frontend should now be running on `http://localhost:3000`

## API Keys Setup

### API-Football
1. Sign up at https://www.api-football.com/
2. Choose appropriate plan (Classic ~$50/month recommended for MVP)
3. Copy API key to backend `.env` as `API_FOOTBALL_KEY`

### The Odds API
1. Sign up at https://the-odds-api.com/
2. Choose appropriate plan (~$200/month for needed requests)
3. Copy API key to backend `.env` as `ODDS_API_KEY`

## Database Migrations

### Create New Migration
```bash
# Backend migrations (Go)
cd backend
migrate create -ext sql -dir database/migrations -seq migration_name

# ML Service migrations (Python with Alembic)
cd ml-service
alembic revision -m "migration_name"
```

### Run Migrations
```bash
# Go backend
go run cmd/migrate/main.go up

# Python ML service
alembic upgrade head
```

### Rollback Migrations
```bash
# Go backend
go run cmd/migrate/main.go down

# Python ML service
alembic downgrade -1
```

## Initial Data Load

### Load Historical Data
```bash
# Run data backfill script (from backend directory)
go run cmd/backfill/main.go --seasons 2021,2022,2023

# This will:
# 1. Fetch all fixtures for specified seasons
# 2. Load team data
# 3. Fetch historical odds (if available)
# 4. Calculate team statistics
```

Expected output: ~1,140 fixtures (380 per season × 3 seasons)

### Verify Data
```sql
-- Connect to database
psql oddsiq

-- Check fixtures count
SELECT season, COUNT(*)
FROM fixtures
GROUP BY season
ORDER BY season;

-- Check teams
SELECT COUNT(*) FROM teams;  -- Should be ~20

-- Check odds
SELECT COUNT(*) FROM odds;  -- Should be substantial

-- Check team stats
SELECT COUNT(*) FROM team_stats;
```

## Running the Full Stack

### Option 1: Manual (Development)
```bash
# Terminal 1 - PostgreSQL (if not using Docker)
# Already running as service

# Terminal 2 - Go Backend
cd backend
go run cmd/api/main.go

# Terminal 3 - Python ML Service
cd ml-service
source venv/bin/activate
uvicorn app.main:app --reload --port 8001

# Terminal 4 - Next.js Frontend
cd frontend
npm run dev
```

### Option 2: Docker Compose (Production-like)
```bash
# Build and start all services
docker-compose up --build

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

Example full `docker-compose.yml`:
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: oddsiq
      POSTGRES_USER: oddsiq_user
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  backend:
    build: ./backend
    ports:
      - "8000:8000"
    environment:
      DATABASE_URL: postgres://oddsiq_user:${DB_PASSWORD}@postgres:5432/oddsiq
      ML_SERVICE_URL: http://ml-service:8001
    depends_on:
      - postgres

  ml-service:
    build: ./ml-service
    ports:
      - "8001:8001"
    environment:
      DATABASE_URL: postgresql://oddsiq_user:${DB_PASSWORD}@postgres:5432/oddsiq
    depends_on:
      - postgres
    volumes:
      - ./ml-service/models:/app/models

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      NEXT_PUBLIC_API_URL: http://localhost:8000/api
    depends_on:
      - backend

volumes:
  postgres_data:
```

## Scheduled Jobs Setup

### Data Sync Jobs
```bash
# Setup cron jobs for daily data updates
# Edit crontab
crontab -e

# Add daily data sync (runs at 2 AM)
0 2 * * * cd /path/to/OddsIQ/backend && go run cmd/sync/main.go fixtures
30 2 * * * cd /path/to/OddsIQ/backend && go run cmd/sync/main.go odds
```

Or use built-in Go scheduler (recommended):
```go
// In backend/cmd/api/main.go
// Schedule daily jobs using go-cron or similar
```

## Verification Checklist

After setup, verify:
- [ ] PostgreSQL is running and accessible
- [ ] Database has `oddsiq` database created
- [ ] Go backend starts without errors on port 8000
- [ ] Python ML service starts without errors on port 8001
- [ ] Next.js frontend accessible at http://localhost:3000
- [ ] API endpoints respond correctly:
  - GET http://localhost:8000/api/health
  - GET http://localhost:8001/health
- [ ] Database contains historical data (fixtures, teams, odds)
- [ ] Frontend can fetch data from backend

## Common Issues & Troubleshooting

### PostgreSQL Connection Failed
```bash
# Check PostgreSQL is running
# macOS
brew services list | grep postgresql

# Ubuntu
sudo systemctl status postgresql

# Verify connection
psql -U oddsiq_user -d oddsiq -h localhost
```

### Go Dependency Issues
```bash
cd backend
go mod tidy
go mod download
```

### Python Package Conflicts
```bash
# Recreate virtual environment
rm -rf venv
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

### Port Already in Use
```bash
# Find process using port
# macOS/Linux
lsof -i :8000

# Windows
netstat -ano | findstr :8000

# Kill process
kill -9 <PID>
```

### API Keys Not Working
- Verify keys are in `.env` files
- Check key format (no quotes needed in .env)
- Ensure API subscription is active
- Check API rate limits

## Development Workflow

### Daily Development
1. Start PostgreSQL (if not running)
2. Start backend: `cd backend && go run cmd/api/main.go`
3. Start ML service: `cd ml-service && uvicorn app.main:app --reload`
4. Start frontend: `cd frontend && npm run dev`

### Making Changes
1. Create feature branch: `git checkout -b feature/your-feature`
2. Make changes
3. Test locally
4. Commit and push
5. Create pull request

### Testing
```bash
# Backend tests
cd backend
go test ./...

# ML Service tests
cd ml-service
pytest

# Frontend tests
cd frontend
npm test
```

## Next Steps

After successful setup:
1. Review API documentation in `docs/api-specification.md`
2. Review database schema in `docs/database-schema.md`
3. Start data backfill for historical matches
4. Train initial ML model
5. Test weekly picks generation workflow
