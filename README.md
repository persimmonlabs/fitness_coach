# Fitness Coach Backend API

A robust, scalable Go-based REST API for fitness tracking with AI-powered meal parsing, workout tracking, and intelligent coaching.

## Quick Start

### Prerequisites

- Go 1.24.4+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Get Started in 3 Steps

1. **Clone and configure**:
```bash
git clone https://github.com/your-org/fitness-coach.git
cd fitness-coach/backend
cp .env.example .env
```

2. **Start with Docker** (recommended):
```bash
make quickstart
```

3. **Access the API**:
```bash
curl http://localhost:8080/health
```

That's it! Your API is running on `http://localhost:8080`

### Manual Setup (without Docker)

```bash
# Install dependencies
go mod download

# Run database migrations
make migrate-up

# Start the server
make run
```

## Features

- User authentication & JWT tokens
- AI-powered meal parsing (text & photos)
- Food database with nutrition info
- Activity & workout tracking
- Progress metrics and goals
- AI chat assistant
- Daily summary stats

## Quick Commands

```bash
# Development
make run              # Start the API server
make watch            # Start with hot reload
make test             # Run tests
make test-coverage    # Generate coverage report

# Database
make migrate-up       # Run migrations
make migrate-down     # Rollback migration
make migrate-create NAME=my_migration  # Create migration
make seed             # Seed database

# Docker
make docker-up        # Start all services
make docker-down      # Stop all services
make docker-logs      # View logs
make docker-clean     # Remove containers & volumes

# Code Quality
make lint             # Run linter
make fmt              # Format code
make vet              # Run go vet

# Build
make build            # Build production binary
make clean            # Clean build artifacts
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh token

### Meals
- `POST /api/v1/meals` - Create meal
- `GET /api/v1/meals` - List meals
- `GET /api/v1/meals/:id` - Get meal
- `PUT /api/v1/meals/:id` - Update meal
- `DELETE /api/v1/meals/:id` - Delete meal

### Foods
- `POST /api/v1/foods` - Create food
- `GET /api/v1/foods` - Search foods
- `GET /api/v1/foods/:id` - Get food
- `PUT /api/v1/foods/:id` - Update food
- `DELETE /api/v1/foods/:id` - Delete food

### Activities
- `POST /api/v1/activities` - Log activity
- `GET /api/v1/activities` - List activities
- CRUD operations similar to meals

### Workouts
- `POST /api/v1/workouts` - Create workout
- `GET /api/v1/workouts` - List workouts
- CRUD operations with exercises and sets

### Other Endpoints
- Exercises, Metrics, Goals, Chat, Summary
- See [API Documentation](./docs/API.md) for complete reference

## Project Structure

```
backend/
├── cmd/api/            # Application entry point
├── internal/
│   ├── core/domain/    # Domain entities
│   ├── services/       # Business logic
│   ├── adapters/       # HTTP, DB, External APIs
│   └── config/         # Configuration
├── migrations/         # Database migrations
├── docs/               # Documentation
├── tests/              # Integration tests
└── scripts/            # Build scripts
```

## Technology Stack

- **Language**: Go 1.24.4
- **Framework**: Gin
- **Database**: PostgreSQL + GORM
- **Auth**: JWT (golang-jwt)
- **AI**: DeepSeek (text), Gemini Vision (photos)
- **Testing**: Testify + Testcontainers
- **Deployment**: Docker + Docker Compose

## Architecture

Clean Architecture (Hexagonal Pattern):
- **Domain Layer**: Business entities and rules
- **Services Layer**: Use case orchestration
- **Adapters Layer**: HTTP, DB, External APIs
- **Independent**: Framework and infrastructure agnostic

See [Architecture Documentation](./docs/ARCHITECTURE.md) for details.

## Configuration

Key environment variables (`.env`):

```env
# Server
SERVER_PORT=8080
SERVER_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=fitness_user
DB_PASSWORD=fitness_password
DB_NAME=fitness_coach

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h

# AI
OPENAI_API_KEY=your-openai-api-key
```

See `.env.example` for all options.

## Testing

```bash
# Run all tests
make test

# With coverage
make test-coverage

# Integration tests
go test ./tests/integration -v
```

## Documentation

Comprehensive documentation is available in the `docs/` directory:

- **[Complete Documentation](./docs/README.md)** - Full project documentation
- **[API Reference](./docs/API.md)** - All endpoints with examples
- **[Architecture Guide](./docs/ARCHITECTURE.md)** - System design and patterns
- **[Meal Parsing](./docs/MEAL_PARSING.md)** - AI meal parsing details
- **[AI Agent](./docs/AI_AGENT.md)** - AI capabilities
- **[Docker Setup](./docs/DOCKER_SETUP.md)** - Docker deployment

## Example Usage

### Register and Login

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123!",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123!"
  }'
```

### Log a Meal

```bash
curl -X POST http://localhost:8080/api/v1/meals \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Breakfast",
    "meal_type": "breakfast",
    "consumed_at": "2025-11-19T08:00:00Z",
    "food_items": [
      {
        "food_id": "123e4567-e89b-12d3-a456-426614174000",
        "quantity": 2.0,
        "unit": "piece"
      }
    ]
  }'
```

## Development Workflow

1. Create feature branch
2. Write tests (TDD)
3. Implement feature
4. Run tests: `make test`
5. Format code: `make fmt`
6. Run linter: `make lint`
7. Create pull request

## Deployment

### Docker

```bash
# Production build
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose logs -f api
```

### Binary

```bash
# Build
make build

# Run
./bin/fitness-coach-api
```

See [Complete Documentation](./docs/README.md) for detailed deployment instructions.

## Troubleshooting

### Database Connection Failed
Check PostgreSQL is running:
```bash
docker-compose ps postgres
```

### Port Already in Use
Change `SERVER_PORT` in `.env` or kill process:
```bash
lsof -ti:8080 | xargs kill -9
```

### Migration Errors
Reset database:
```bash
make db-reset
```

## Performance

- **Full-text search** with GIN indexes
- **Connection pooling** (100 max connections)
- **Denormalized totals** for fast queries
- **Rate limiting** (100 req/min per user)
- **Stateless design** (horizontal scaling ready)

## Security

- JWT authentication
- bcrypt password hashing
- Input validation
- SQL injection prevention (GORM)
- CORS configured
- Rate limiting

## Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open pull request

## License

TBD

## Support

- **Documentation**: See `docs/` directory
- **Issues**: GitHub Issues
- **Contact**: your-email@example.com

## Version

**v1.0.0** - Initial Release (2025-11-19)

Features:
- User authentication
- AI meal parsing (text & photos)
- Food database
- Activity & workout tracking
- Metrics & goals
- Chat assistant
- Daily summaries

---

For detailed documentation, see [docs/README.md](./docs/README.md)
