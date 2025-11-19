# Fitness Coach Backend API - Complete Documentation

A robust, scalable Go-based REST API for fitness tracking with AI-powered meal parsing, workout tracking, and intelligent coaching insights.

## Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Running the Application](#running-the-application)
- [Running Tests](#running-tests)
- [Deployment](#deployment)
- [Project Structure](#project-structure)
- [API Documentation](#api-documentation)
- [Additional Documentation](#additional-documentation)

## Project Overview

The Fitness Coach Backend is a comprehensive API that powers a complete fitness tracking application. It provides:

- User authentication and profile management
- AI-powered meal logging (text and photo parsing)
- Food database with nutritional information
- Activity and workout tracking
- Exercise library management
- Progress metrics and goals
- AI chat assistant for fitness guidance
- Daily summary statistics

### Key Features

- **AI Meal Parsing**: Parse meals from natural language text or food photos using DeepSeek and Gemini Vision
- **Intelligent Food Matching**: Full-text search with fuzzy matching against food database
- **AI-Generated Foods**: Automatic nutrition estimation for unknown foods
- **Hexagonal Architecture**: Clean separation of concerns with ports and adapters pattern
- **Type-Safe**: Strongly typed Go with comprehensive validation
- **RESTful API**: Follows REST principles with proper HTTP semantics
- **JWT Authentication**: Secure token-based authentication
- **Database Migrations**: Version-controlled schema evolution
- **Comprehensive Testing**: Unit and integration tests with testcontainers
- **Docker Support**: Containerized deployment ready

## Architecture

This project implements **Clean Architecture** principles with the **Hexagonal (Ports & Adapters)** pattern for maximum maintainability and testability.

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     HTTP Layer (Adapters)                   │
│  ┌────────────┬────────────┬────────────┬─────────────┐    │
│  │  Handlers  │ Middleware │  Router    │     DTOs    │    │
│  └────────────┴────────────┴────────────┴─────────────┘    │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────┴──────────────────────────────────┐
│                      Services Layer                         │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Auth │ Meal │ Food │ Activity │ Chat │ Summary   │    │
│  └────────────────────────────────────────────────────┘    │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────┴──────────────────────────────────┐
│                      Core Domain                            │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Entities: User, Meal, Food, Activity, Exercise    │    │
│  │  Business Logic & Domain Rules                     │    │
│  │  Port Interfaces (Repositories, External Services) │    │
│  └────────────────────────────────────────────────────┘    │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────┴──────────────────────────────────┐
│                   Infrastructure Layer                      │
│  ┌──────────────┬──────────────┬──────────────────────┐   │
│  │  PostgreSQL  │  OpenRouter  │  Gemini Vision API   │   │
│  │  Repos       │  (DeepSeek)  │  (Photo Analysis)    │   │
│  └──────────────┴──────────────┴──────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

#### 1. Domain Layer (`internal/core/domain/`)
- Pure business entities with no external dependencies
- Domain models: User, Meal, Food, Activity, Workout, Exercise, Metric, Goal
- Business rules and validation logic
- Interface definitions (ports)

#### 2. Services Layer (`internal/services/`)
- Application business logic orchestration
- Use case implementation
- Coordination between domain and adapters
- Transaction management

#### 3. Adapters Layer (`internal/adapters/`)
- **HTTP Adapters**: REST API handlers, middleware, routing
- **Repository Adapters**: PostgreSQL database implementations
- **External Adapters**: AI services (OpenRouter, Vision API)

#### 4. Configuration Layer (`internal/config/`, `pkg/config/`)
- Environment variable loading
- Configuration validation
- Database connection setup

See [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed architecture documentation.

## Technology Stack

### Core Framework
- **Language**: Go 1.24.4
- **Web Framework**: Gin Web Framework
- **Configuration**: Viper
- **Logging**: Uber Zap

### Database
- **Primary Database**: PostgreSQL 15+
- **ORM**: GORM
- **Migrations**: Goose (golang-migrate)
- **Full-Text Search**: PostgreSQL GIN indexes

### Authentication & Security
- **JWT**: golang-jwt/jwt/v5
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator/v10

### AI & External Services
- **Text Parsing**: OpenRouter API (DeepSeek model)
- **Vision Parsing**: Google Gemini 2.0 Flash Exp (via OpenRouter)
- **LangChain**: tmc/langchaingo

### Testing
- **Testing Framework**: Go testing package
- **Assertions**: Testify
- **Mocking**: testify/mock, uber/mock
- **Integration Tests**: Testcontainers

### DevOps & Deployment
- **Containerization**: Docker, Docker Compose
- **Build Tool**: Make
- **Hot Reload**: Air
- **Linting**: golangci-lint

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.24.4 or higher** - [Download](https://golang.org/dl/)
- **PostgreSQL 15+** - [Download](https://www.postgresql.org/download/)
- **Docker & Docker Compose** (optional, for containerized deployment) - [Download](https://www.docker.com/get-started)
- **Make** (optional, for convenient commands) - Usually pre-installed on Unix systems

### Recommended Tools

- **Air** - For hot reload during development
- **Goose** - For database migrations
- **golangci-lint** - For code quality checks

Install development tools:
```bash
make install-tools
```

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/your-org/fitness-coach.git
cd fitness-coach/backend
```

### 2. Install Dependencies

```bash
go mod download
go mod verify
```

### 3. Environment Configuration

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` with your configuration (see [Configuration](#configuration) section).

### 4. Quick Start (Using Docker)

The fastest way to get started:

```bash
# Start all services (PostgreSQL + API)
make quickstart
```

This command will:
- Start PostgreSQL database
- Run database migrations
- Start the API server on port 8080

### 5. Manual Setup

If you prefer manual setup without Docker:

```bash
# Setup development environment
make setup

# Start PostgreSQL (if not using Docker)
# Make sure PostgreSQL is running and accessible

# Run database migrations
make migrate-up

# Start the API server
make run
```

## Configuration

The application is configured via environment variables. All configuration options are documented in `.env.example`.

### Critical Configuration Variables

#### Server Configuration
```env
SERVER_PORT=8080
SERVER_HOST=localhost
SERVER_ENV=development
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
```

#### Database Configuration
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=fitness_user
DB_PASSWORD=fitness_password
DB_NAME=fitness_coach
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
```

#### JWT Configuration
```env
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h
```

#### AI Configuration
```env
OPENAI_API_KEY=your-openai-api-key
AI_MODEL=gpt-4
AI_TEMPERATURE=0.7
AI_MAX_TOKENS=1000
```

#### CORS Configuration
```env
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,PATCH,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization,Accept
CORS_ALLOW_CREDENTIALS=true
```

#### Rate Limiting
```env
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m
```

### Environment-Specific Configurations

#### Development
- Detailed logging
- Hot reload enabled
- Debug mode
- Relaxed CORS

#### Production
- JSON logging
- Performance optimizations
- Strict security
- Rate limiting enforced

## Database Setup

### Using Docker (Recommended for Development)

```bash
# Start PostgreSQL container
make docker-up

# Run migrations
make migrate-up

# (Optional) Seed database with sample data
make seed
```

### Manual PostgreSQL Setup

1. Create database and user:
```sql
CREATE DATABASE fitness_coach;
CREATE USER fitness_user WITH PASSWORD 'fitness_password';
GRANT ALL PRIVILEGES ON DATABASE fitness_coach TO fitness_user;
```

2. Run migrations:
```bash
make migrate-up
```

### Database Migrations

The project uses Goose for database migrations.

#### Create a new migration:
```bash
make migrate-create NAME=add_new_feature
```

#### Run migrations:
```bash
make migrate-up
```

#### Rollback last migration:
```bash
make migrate-down
```

#### Check migration status:
```bash
make migrate-status
```

### Database Schema

The database includes the following main tables:

- `users` - User accounts and profiles
- `foods` - Food items with nutritional information
- `serving_units` - Common serving units (cup, tbsp, etc.)
- `food_ingredients` - Composite food ingredients
- `food_serving_conversions` - Unit conversion factors
- `meals` - User meal logs
- `meal_food_items` - Foods in each meal
- `activities` - Cardio and general activities
- `exercises` - Exercise library
- `workouts` - Workout sessions
- `workout_exercises` - Exercises in workouts
- `workout_sets` - Individual sets in exercises
- `metrics` - User health metrics (weight, body fat, etc.)
- `daily_summaries` - Aggregated daily statistics
- `goals` - User fitness goals
- `conversations` - Chat conversation threads
- `messages` - Individual chat messages

See [Database Schema Diagram](./ARCHITECTURE.md#database-schema) for detailed relationships.

## Running the Application

### Development Mode

#### Using Make (Recommended)
```bash
# Start the application
make run

# Or with hot reload (requires Air)
make watch
```

#### Direct Go Command
```bash
go run cmd/api/main.go
```

### Production Mode

#### Using Docker
```bash
# Build and start all services
docker-compose up -d

# View logs
make docker-logs

# Stop services
make docker-down
```

#### Binary Build
```bash
# Build production binary
make build

# Run the binary
./bin/fitness-coach-api
```

### Accessing the API

Once running, the API will be available at:
- **Base URL**: `http://localhost:8080`
- **Health Check**: `http://localhost:8080/health`
- **API v1**: `http://localhost:8080/api/v1`

Test the health endpoint:
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "fitness-coach-api",
  "version": "1.0.0",
  "time": "2025-11-19T00:00:00Z"
}
```

## Running Tests

### Unit Tests

Run all unit tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

This generates:
- Terminal coverage report
- `coverage.html` for detailed visualization

### Integration Tests

Integration tests use testcontainers to spin up real PostgreSQL instances:

```bash
go test ./tests/integration -v
```

### Test Structure

```
tests/
├── integration/          # Integration tests with real database
│   ├── auth_test.go
│   ├── meal_test.go
│   └── ...
└── unit/                 # Unit tests (alongside code)
    └── services/
        ├── auth_service_test.go
        └── ...
```

### Code Coverage Goals

- Minimum coverage: 80%
- Critical paths: 95%+
- Service layer: 90%+

## Deployment

### Docker Deployment

#### Production Docker Compose

```bash
# Build and start
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose logs -f api

# Stop
docker-compose down
```

#### Environment Variables

Ensure production `.env` is properly configured:
- Strong `JWT_SECRET`
- Production database credentials
- Valid API keys
- Appropriate `CORS_ALLOWED_ORIGINS`

### Cloud Deployment Options

#### AWS ECS/Fargate
1. Build Docker image
2. Push to ECR
3. Create ECS task definition
4. Deploy service

#### Google Cloud Run
```bash
gcloud run deploy fitness-coach-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

#### Kubernetes
Use the provided Kubernetes manifests (if available) or create your own based on Docker deployment.

### Health Checks

Configure load balancers to use:
- **Health Endpoint**: `GET /health`
- **Expected Status**: 200
- **Timeout**: 5s
- **Interval**: 30s

### Monitoring

Recommended monitoring setup:
- Application logs (Zap JSON output)
- Database connection pool metrics
- HTTP request metrics (latency, status codes)
- Error rate tracking
- AI API usage and costs

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── core/
│   │   ├── domain/                 # Domain entities and models
│   │   │   ├── user.go
│   │   │   ├── meal.go
│   │   │   ├── food.go
│   │   │   ├── activity.go
│   │   │   ├── workout.go
│   │   │   ├── exercise.go
│   │   │   ├── metric.go
│   │   │   ├── goal.go
│   │   │   └── conversation.go
│   │   └── ports/                  # Interface definitions (if separate)
│   ├── services/                   # Business logic implementation
│   │   ├── auth_service.go
│   │   ├── meal_service.go
│   │   ├── meal_parser_service.go  # AI meal parsing
│   │   ├── food_service.go
│   │   ├── activity_service.go
│   │   ├── workout_service.go
│   │   ├── metric_service.go
│   │   ├── goal_service.go
│   │   └── summary_service.go
│   ├── adapters/
│   │   ├── http/
│   │   │   ├── handlers/           # HTTP request handlers
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── meal_handler.go
│   │   │   │   ├── food_handler.go
│   │   │   │   ├── activity_handler.go
│   │   │   │   ├── workout_handler.go
│   │   │   │   ├── exercise_handler.go
│   │   │   │   ├── metric_handler.go
│   │   │   │   ├── goal_handler.go
│   │   │   │   ├── chat_handler.go
│   │   │   │   └── summary_handler.go
│   │   │   ├── middleware/         # HTTP middleware
│   │   │   │   ├── auth.go
│   │   │   │   ├── cors.go
│   │   │   │   ├── logger.go
│   │   │   │   ├── rate_limiter.go
│   │   │   │   ├── recovery.go
│   │   │   │   ├── request_id.go
│   │   │   │   └── validator.go
│   │   │   ├── dto/                # Data Transfer Objects
│   │   │   │   ├── request_models.go
│   │   │   │   └── response_models.go
│   │   │   └── router.go           # Route definitions
│   │   ├── repositories/postgres/   # Database implementations
│   │   │   ├── user_repo.go
│   │   │   ├── meal_repo.go
│   │   │   ├── food_repo.go
│   │   │   ├── activity_repo.go
│   │   │   ├── workout_repo.go
│   │   │   ├── metric_repo.go
│   │   │   ├── goal_repo.go
│   │   │   └── conversation_repo.go
│   │   └── external/               # External API integrations
│   │       ├── openrouter_client.go  # DeepSeek text parsing
│   │       ├── vision_client.go      # Gemini vision parsing
│   │       └── supabase_storage_client.go
│   ├── config/                     # Configuration management
│   │   ├── config.go
│   │   └── database.go
│   └── pkg/                        # Shared packages
│       ├── errors/                 # Custom error types
│       └── utils/                  # Utility functions
├── migrations/                     # Database migration files
│   ├── 001_create_users.up.sql
│   ├── 001_create_users.down.sql
│   ├── 002_create_foods.up.sql
│   └── ...
├── tests/
│   ├── integration/                # Integration tests
│   └── unit/                       # Unit tests
├── scripts/                        # Build and deployment scripts
│   └── seed.sql                    # Database seed data
├── docs/                           # Documentation
│   ├── README.md                   # This file
│   ├── API.md                      # API documentation
│   ├── ARCHITECTURE.md             # Architecture details
│   ├── MEAL_PARSING.md             # Meal parsing documentation
│   └── AI_AGENT.md                 # AI agent capabilities
├── .air.toml                       # Hot reload configuration
├── .env.example                    # Environment variables template
├── .gitignore
├── docker-compose.yml              # Docker Compose configuration
├── Dockerfile                      # Docker image definition
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
├── Makefile                        # Build and deployment commands
└── README.md                       # Quick start guide
```

## API Documentation

Complete API documentation is available in [API.md](./API.md).

### Quick Reference

#### Authentication Endpoints
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh access token

#### Meal Endpoints
- `POST /api/v1/meals` - Create meal
- `GET /api/v1/meals` - List meals
- `GET /api/v1/meals/:id` - Get meal by ID
- `PUT /api/v1/meals/:id` - Update meal
- `DELETE /api/v1/meals/:id` - Delete meal

#### Food Endpoints
- `POST /api/v1/foods` - Create food
- `GET /api/v1/foods` - List foods
- `GET /api/v1/foods/:id` - Get food by ID
- `PUT /api/v1/foods/:id` - Update food
- `DELETE /api/v1/foods/:id` - Delete food

See [API.md](./API.md) for complete endpoint documentation with request/response examples.

## Additional Documentation

### Detailed Documentation

- **[API.md](./API.md)** - Complete API reference with examples
- **[ARCHITECTURE.md](./ARCHITECTURE.md)** - System architecture and design patterns
- **[MEAL_PARSING.md](./MEAL_PARSING.md)** - AI meal parsing implementation
- **[AI_AGENT.md](./AI_AGENT.md)** - AI agent capabilities and tools
- **[DOCKER_SETUP.md](./DOCKER_SETUP.md)** - Docker deployment guide

### External Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## Development Guidelines

### Code Style
- Follow Go idioms and best practices
- Use `gofmt` for consistent formatting
- Keep functions small and focused (max 50 lines)
- Write meaningful comments for exported functions
- Use descriptive variable names

### Testing Requirements
- Write tests for all business logic
- Use table-driven tests where applicable
- Maintain >80% code coverage
- Test edge cases and error conditions
- Use testcontainers for integration tests

### API Design Principles
- Follow RESTful conventions
- Use proper HTTP status codes
- Implement consistent error responses
- Version APIs (`/api/v1/`)
- Document all endpoints

### Security Best Practices
- Never commit secrets to version control
- Use environment variables for sensitive data
- Implement rate limiting on all endpoints
- Validate all input data
- Use prepared statements (GORM handles this)
- Hash passwords with bcrypt
- Use HTTPS in production

### Git Workflow
1. Create feature branch from `main`
2. Write tests first (TDD)
3. Implement feature
4. Ensure all tests pass
5. Run linting: `make lint`
6. Format code: `make fmt`
7. Submit pull request

## Troubleshooting

### Common Issues

#### Database Connection Failed
```
Error: failed to connect to database
```
**Solution**: Check PostgreSQL is running and credentials in `.env` are correct.

#### Port Already in Use
```
Error: bind: address already in use
```
**Solution**: Change `SERVER_PORT` in `.env` or kill process using port 8080.

#### Migration Failed
```
Error: migration failed
```
**Solution**: Check database permissions and migration file syntax. Roll back: `make migrate-down`

#### AI Parsing Not Working
```
Error: OpenRouter API request failed
```
**Solution**: Verify `OPENROUTER_API_KEY` is set correctly in `.env`.

### Getting Help

- Check the [API.md](./API.md) for endpoint documentation
- Review [ARCHITECTURE.md](./ARCHITECTURE.md) for design details
- Check GitHub Issues for known problems
- Contact the development team

## Performance Optimization

### Database Optimization
- Indexes on frequently queried columns
- Connection pooling configured
- Full-text search with GIN indexes
- Denormalized totals for performance

### API Performance
- Request/response compression
- Response caching (where appropriate)
- Efficient pagination
- Database query optimization

### Monitoring Recommendations
- Track API response times
- Monitor database query performance
- Watch AI API usage and costs
- Track error rates
- Monitor memory and CPU usage

## License

TBD

## Contributors

See CONTRIBUTORS.md for the list of contributors.

## Version History

- **v1.0.0** (2025-11-19) - Initial release
  - User authentication and profiles
  - AI-powered meal parsing
  - Food database management
  - Activity and workout tracking
  - AI chat assistant
  - Daily summaries and goals

---

For quick start instructions, see the root [README.md](../README.md).
