# Docker Development Setup

## Quick Start

The fastest way to get the Fitness Coach API running locally:

```bash
# 1. Clone and navigate to backend directory
cd backend

# 2. Copy environment file
cp .env.example .env

# 3. Start all services with one command
make quickstart
```

That's it! The API will be available at http://localhost:8080

## Prerequisites

- Docker Desktop (Windows/Mac) or Docker Engine (Linux)
- Docker Compose v2.0+
- Make (usually pre-installed on Mac/Linux, Windows users can use WSL or Git Bash)
- Go 1.21+ (for local development without Docker)

## Project Structure

```
backend/
├── cmd/api/              # Application entry point
├── migrations/           # Database migrations
├── docker-compose.yml    # Docker services configuration
├── Dockerfile           # Multi-stage Docker build
├── Makefile             # Development commands
└── .env.example         # Environment template
```

## Docker Services

### PostgreSQL Database
- **Image**: postgres:15-alpine
- **Port**: 5432
- **Database**: fitness_coach
- **User**: fitness_user
- **Password**: fitness_password (changeable via .env)
- **Volume**: Persistent data storage
- **Health Check**: Automatic readiness verification

### API Service (Optional)
- **Build**: Multi-stage Dockerfile
- **Port**: 8080
- **Hot Reload**: Enabled with volume mounts
- **Dependencies**: Waits for database to be healthy

## Available Make Commands

### Quick Commands
```bash
make help          # Show all available commands
make quickstart    # Complete setup for new developers
make dev           # Start development environment
make dev-down      # Stop development environment
```

### Setup & Dependencies
```bash
make setup         # Install Go dependencies and tools
make install-tools # Install development tools (golangci-lint, goose, air)
```

### Database Operations
```bash
make migrate-up              # Run all pending migrations
make migrate-down            # Rollback last migration
make migrate-create NAME=... # Create new migration file
make migrate-status          # Show migration status
make seed                    # Seed database with sample data
make db-reset                # Reset database (drop, recreate, migrate)
make db-shell                # Open PostgreSQL shell
```

### Application
```bash
make run           # Run application locally (without Docker)
make build         # Build production binary
make watch         # Run with hot reload (requires air)
```

### Testing
```bash
make test          # Run all tests with coverage
make test-coverage # Generate HTML coverage report
```

### Docker
```bash
make docker-up       # Start Docker services
make docker-down     # Stop Docker services
make docker-logs     # View all logs
make docker-logs-api # View API logs only
make docker-logs-db  # View database logs only
make docker-rebuild  # Rebuild images from scratch
make docker-clean    # Remove containers and volumes
```

### Code Quality
```bash
make lint  # Run golangci-lint
make fmt   # Format code
make vet   # Run go vet
```

## Development Workflows

### Option 1: Full Docker Development (Recommended for beginners)

Everything runs in containers, no local Go installation needed:

```bash
# Start services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

Access:
- API: http://localhost:8080
- Database: localhost:5432

### Option 2: Hybrid Development (Recommended for active development)

Database in Docker, API runs locally for faster iteration:

```bash
# Start only database
docker-compose up -d postgres

# Run migrations
make migrate-up

# Run API locally with hot reload
make watch
# OR without hot reload
make run
```

### Option 3: Full Local Development

Both database and API run locally:

```bash
# Requires PostgreSQL installed locally
# Update .env with local database connection

# Run migrations
make migrate-up

# Run application
make run
```

## Environment Configuration

### Required Environment Variables

Copy `.env.example` to `.env` and update:

```bash
# Database (adjust for your setup)
DB_HOST=postgres        # Use 'localhost' for local development
DB_PORT=5432
DB_USER=fitness_user
DB_PASSWORD=fitness_password
DB_NAME=fitness_coach

# Application
PORT=8080
JWT_SECRET=change-this-to-a-secure-secret

# AI Features (LangChain)
OPENAI_API_KEY=your-openai-api-key  # Required for AI features
```

### Docker vs Local Development

**For Docker**: Use `DB_HOST=postgres` (container name)
**For Local**: Use `DB_HOST=localhost`

## Database Migrations

### Creating Migrations

```bash
# Create a new migration
make migrate-create NAME=add_user_preferences

# This creates two files:
# migrations/YYYYMMDDHHMMSS_add_user_preferences.up.sql
# migrations/YYYYMMDDHHMMSS_add_user_preferences.down.sql
```

### Running Migrations

```bash
# Apply all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check migration status
make migrate-status
```

### Migration Best Practices

1. Always create both up and down migrations
2. Test rollback before committing
3. Keep migrations atomic and reversible
4. Don't modify existing migrations in production

## Hot Reload Development

For fastest development cycle, use Air for hot reload:

```bash
# Install air
make install-tools

# Start with hot reload
make watch
```

Air configuration is in `.air.toml` and will:
- Watch for file changes
- Automatically rebuild
- Restart the application
- Exclude test files

## Testing

### Run Tests

```bash
# Run all tests
make test

# Generate coverage report
make test-coverage
# Opens coverage.html in browser
```

### Test Database

Tests use a separate database configuration:
- Database: fitness_tracker_test
- Port: 5433 (to avoid conflicts)

Configure in `.env`:
```bash
TEST_DB_HOST=localhost
TEST_DB_PORT=5433
TEST_DB_NAME=fitness_tracker_test
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using port 8080
lsof -i :8080  # Mac/Linux
netstat -ano | findstr :8080  # Windows

# Use different port in .env
PORT=8081
```

### Database Connection Failed

```bash
# Check database is running
docker-compose ps

# View database logs
make docker-logs-db

# Verify health check
docker-compose exec postgres pg_isready -U fitness_user
```

### Migration Errors

```bash
# Check migration status
make migrate-status

# Reset database (WARNING: loses all data)
make db-reset
```

### Docker Build Issues

```bash
# Clean rebuild
make docker-clean
make docker-rebuild

# Clear Docker cache
docker system prune -a
```

### Permission Issues (Linux)

```bash
# Fix file permissions
sudo chown -R $USER:$USER .

# Or run Docker commands with sudo
sudo make docker-up
```

## Production Deployment

### Building Production Image

```bash
# Build optimized production image
docker build --target production -t fitness-coach-api:latest .

# Run production container
docker run -d \
  -p 8080:8080 \
  --env-file .env.production \
  fitness-coach-api:latest
```

### Production Checklist

- [ ] Use strong JWT_SECRET
- [ ] Set ENVIRONMENT=production
- [ ] Enable SSL for database (DB_SSLMODE=require)
- [ ] Use managed database (not container)
- [ ] Set appropriate CORS origins
- [ ] Configure proper logging
- [ ] Set up health check monitoring
- [ ] Use secrets management (not .env file)
- [ ] Enable rate limiting
- [ ] Configure backups

## Health Checks

### Application Health

```bash
# Check API health
curl http://localhost:8080/health

# Expected response:
{"status": "healthy", "database": "connected"}
```

### Database Health

```bash
# Check database
docker-compose exec postgres pg_isready -U fitness_user

# Connect to database
make db-shell
```

## Performance Optimization

### Database Connection Pooling

Configured in `.env`:
```bash
DB_MAX_OPEN_CONNS=25      # Maximum open connections
DB_MAX_IDLE_CONNS=5       # Idle connections to keep
DB_CONN_MAX_LIFETIME=5m   # Maximum connection lifetime
```

### Docker Resource Limits

Add to `docker-compose.yml`:
```yaml
services:
  api:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
```

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [Goose Migrations](https://github.com/pressly/goose)

## Support

For issues or questions:
1. Check logs: `make docker-logs`
2. Verify environment: `.env` configuration
3. Check database: `make db-shell`
4. Review migrations: `make migrate-status`

## Next Steps

After setup:
1. Review API documentation at `/api/docs`
2. Explore database schema in `migrations/`
3. Check example requests in `docs/API_EXAMPLES.md`
4. Read architecture guide in `docs/ARCHITECTURE.md`
