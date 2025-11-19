# Fitness Coach Backend - Architecture Documentation

This document provides detailed information about the system architecture, design patterns, and technical implementation of the Fitness Coach Backend API.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Architectural Patterns](#architectural-patterns)
- [Layer Design](#layer-design)
- [Database Schema](#database-schema)
- [Service Interactions](#service-interactions)
- [External Integrations](#external-integrations)
- [Security Architecture](#security-architecture)
- [Performance Optimizations](#performance-optimizations)
- [Scalability Considerations](#scalability-considerations)

## Architecture Overview

The Fitness Coach Backend implements **Clean Architecture** (also known as Onion Architecture) combined with the **Hexagonal (Ports & Adapters)** pattern. This architecture promotes:

- **Separation of Concerns**: Clear boundaries between layers
- **Testability**: Business logic independent of frameworks
- **Maintainability**: Easy to understand and modify
- **Flexibility**: Easy to swap implementations
- **Scalability**: Support for horizontal scaling

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                            │
│         (Mobile Apps, Web Frontend, Third-party Apps)           │
└───────────────────────────┬─────────────────────────────────────┘
                            │ HTTP/REST
┌───────────────────────────┴─────────────────────────────────────┐
│                    HTTP Adapters Layer                          │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Gin Router │ Handlers │ Middleware │ DTOs │ Validation │  │
│  └──────────────────────────────────────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────────┘
                            │
┌───────────────────────────┴─────────────────────────────────────┐
│                      Services Layer                             │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Business Logic │ Use Cases │ Orchestration │ Validation │  │
│  │  Auth │ Meal Parser │ Food │ Activity │ Workout │ Chat   │  │
│  └──────────────────────────────────────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────────┘
                            │
┌───────────────────────────┴─────────────────────────────────────┐
│                      Core Domain Layer                          │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Domain Entities │ Business Rules │ Value Objects        │  │
│  │  Port Interfaces (Repositories, External Services)       │  │
│  └──────────────────────────────────────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────────┘
                            │
┌───────────────────────────┴─────────────────────────────────────┐
│                   Infrastructure Layer                          │
│  ┌──────────────┬──────────────┬────────────────────────────┐  │
│  │  PostgreSQL  │  OpenRouter  │   Gemini Vision API        │  │
│  │  Repositories│  (DeepSeek)  │   Supabase Storage         │  │
│  └──────────────┴──────────────┴────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Architectural Patterns

### 1. Clean Architecture (Onion Architecture)

The application is organized in concentric layers, with dependencies pointing inward:

```
External Systems → Adapters → Services → Domain Core
```

**Key Principles:**
- **Dependency Inversion**: Outer layers depend on inner layers, never vice versa
- **Domain Independence**: Core business logic has no external dependencies
- **Framework Independence**: Not tied to specific frameworks or libraries

### 2. Hexagonal Architecture (Ports & Adapters)

The core application defines **ports** (interfaces), and infrastructure provides **adapters** (implementations):

**Ports** (Interfaces):
- `UserRepository`
- `MealRepository`
- `FoodRepository`
- `AIClient`
- `VisionClient`
- `StorageClient`

**Adapters** (Implementations):
- PostgreSQL repositories
- OpenRouter HTTP client
- Gemini Vision client
- Supabase storage client

### 3. Repository Pattern

Abstracts data access, providing:
- Clean separation between data access and business logic
- Easy to mock for testing
- Swappable implementations (PostgreSQL, MongoDB, etc.)

**Example**:
```go
// Port (interface)
type FoodRepository interface {
    Create(ctx context.Context, food *domain.Food) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Food, error)
    SearchFoods(ctx context.Context, query string, limit int) ([]domain.Food, error)
    Update(ctx context.Context, food *domain.Food) error
    Delete(ctx context.Context, id uuid.UUID) error
}

// Adapter (PostgreSQL implementation)
type PostgresFoodRepository struct {
    db *gorm.DB
}
```

### 4. Dependency Injection

All dependencies are injected via constructors:

```go
// Service depends on repository interface, not concrete implementation
func NewMealService(repo ports.MealRepository) *MealService {
    return &MealService{
        repo: repo,
    }
}

// Handler depends on service interface
func NewMealHandler(service ports.MealService, logger *zap.Logger) *MealHandler {
    return &MealHandler{
        service: service,
        logger:  logger,
    }
}
```

## Layer Design

### 1. Domain Layer (`internal/core/domain/`)

**Responsibilities:**
- Define core business entities
- Contain domain business rules
- Define port interfaces
- Pure Go, no external dependencies

**Key Components:**
- **Entities**: User, Meal, Food, Activity, Workout, Exercise, Metric, Goal
- **Value Objects**: ParsedMeal, ParsedFoodItem
- **Interfaces**: Repository ports (when defined separately)

**Example Entity**:
```go
type Meal struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Name      string
    MealType  string
    ConsumedAt time.Time

    // Denormalized for performance
    TotalCalories     float64
    TotalProtein      float64
    TotalCarbohydrates float64
    TotalFat          float64

    FoodItems []MealFoodItem
}
```

### 2. Services Layer (`internal/services/`)

**Responsibilities:**
- Implement business use cases
- Orchestrate domain objects and repositories
- Handle transactions
- Enforce business rules

**Key Services:**
- **AuthService**: User registration, login, token management
- **MealService**: CRUD operations for meals
- **MealParserService**: AI-powered meal parsing (text & photos)
- **FoodService**: Food database management
- **ActivityService**: Cardio activity tracking
- **WorkoutService**: Strength training management
- **MetricService**: Body metrics tracking
- **GoalService**: Goal setting and progress
- **ChatService**: AI chat assistant
- **SummaryService**: Aggregate statistics

**Example Service**:
```go
type MealParserService struct {
    openRouterClient *external.OpenRouterClient
    visionClient     *external.VisionClient
    foodRepo         ports.FoodRepository
}

func (s *MealParserService) ParseText(ctx context.Context, userID uuid.UUID, text string) (*domain.ParsedMeal, error) {
    // 1. Call AI to extract food items
    // 2. Match foods against database
    // 3. Create AI-generated foods for unknowns
    // 4. Return structured ParsedMeal
}
```

### 3. Adapters Layer (`internal/adapters/`)

#### HTTP Adapters (`adapters/http/`)

**Responsibilities:**
- Handle HTTP requests/responses
- Route incoming requests
- Apply middleware
- Transform DTOs to domain models

**Components:**
- **Handlers**: Convert HTTP to service calls
- **Middleware**: Auth, logging, rate limiting, CORS, recovery
- **DTOs**: Request/response data structures
- **Router**: Route configuration

**Request Flow**:
```
Client Request
    → Router
    → Middleware Chain (Auth, Logger, Rate Limiter)
    → Handler
    → Service
    → Repository
    → Database
```

#### Repository Adapters (`adapters/repositories/postgres/`)

**Responsibilities:**
- Implement repository interfaces
- Translate between domain models and database models
- Execute database queries
- Handle database errors

**Example Repository**:
```go
type PostgresMealRepository struct {
    db *gorm.DB
}

func (r *PostgresMealRepository) Create(ctx context.Context, meal *domain.Meal) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Save meal
        if err := tx.Create(meal).Error; err != nil {
            return err
        }

        // Calculate totals (denormalized)
        meal.TotalCalories = calculateTotals(meal.FoodItems)

        return tx.Save(meal).Error
    })
}
```

#### External Adapters (`adapters/external/`)

**Responsibilities:**
- Integrate with external APIs
- Handle API authentication
- Transform external responses to domain models
- Retry logic and error handling

**Components:**
- **OpenRouterClient**: DeepSeek AI for text parsing
- **VisionClient**: Gemini Vision for photo analysis
- **SupabaseStorageClient**: File storage

### 4. Configuration Layer (`internal/config/`, `pkg/config/`)

**Responsibilities:**
- Load environment variables
- Validate configuration
- Provide typed configuration objects
- Database connection setup

## Database Schema

### Entity Relationship Diagram

```
┌─────────────┐
│    Users    │
└──────┬──────┘
       │
       ├──────────────┬──────────────┬──────────────┬─────────────┐
       │              │              │              │             │
       ▼              ▼              ▼              ▼             ▼
 ┌──────────┐  ┌────────────┐ ┌────────────┐ ┌─────────┐  ┌──────────┐
 │  Meals   │  │ Activities │ │  Workouts  │ │ Metrics │  │  Goals   │
 └────┬─────┘  └────────────┘ └─────┬──────┘ └─────────┘  └──────────┘
      │                              │
      │                              │
      ▼                              ▼
┌─────────────────┐       ┌──────────────────┐
│ MealFoodItems   │       │ WorkoutExercises │
└────┬────────────┘       └─────────┬────────┘
     │                              │
     │                              ▼
     │                    ┌─────────────────┐
     │                    │  WorkoutSets    │
     │                    └─────────────────┘
     │
     ▼
┌──────────┐
│  Foods   │
└────┬─────┘
     │
     ├──────────────┬─────────────────────┐
     │              │                     │
     ▼              ▼                     ▼
┌────────────┐ ┌──────────┐  ┌────────────────────────┐
│ Ingredients│ │ Servings │  │ ServingConversions     │
└────────────┘ └──────────┘  └────────────────────────┘
```

### Key Tables

#### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    date_of_birth DATE,
    gender VARCHAR(20),
    height_cm DECIMAL(5,2),
    weight_kg DECIMAL(5,2),
    activity_level VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

#### Foods Table
```sql
CREATE TABLE foods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fdc_id INT UNIQUE,  -- USDA FoodData Central ID
    name VARCHAR(255) NOT NULL,
    description TEXT,
    brand VARCHAR(255),
    category VARCHAR(100),
    serving_size DECIMAL(10,2) NOT NULL,
    serving_unit VARCHAR(50) NOT NULL,
    calories DECIMAL(10,2) NOT NULL,
    protein DECIMAL(10,2) NOT NULL,
    carbohydrates DECIMAL(10,2) NOT NULL,
    fat DECIMAL(10,2) NOT NULL,
    fiber DECIMAL(10,2),
    sugar DECIMAL(10,2),
    saturated_fat DECIMAL(10,2),
    trans_fat DECIMAL(10,2),
    cholesterol DECIMAL(10,2),
    sodium DECIMAL(10,2),
    potassium DECIMAL(10,2),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    source VARCHAR(100),  -- usda, user, ai_generated
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_food_search ON foods USING GIN(to_tsvector('english', name));
```

#### Meals Table
```sql
CREATE TABLE meals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    meal_type VARCHAR(50) NOT NULL,  -- breakfast, lunch, dinner, snack
    consumed_at TIMESTAMP NOT NULL,
    notes TEXT,
    -- Denormalized totals for performance
    total_calories DECIMAL(10,2) NOT NULL,
    total_protein DECIMAL(10,2) NOT NULL,
    total_carbohydrates DECIMAL(10,2) NOT NULL,
    total_fat DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_user_meals ON meals(user_id, consumed_at);
```

### Indexes

**Performance Indexes:**
- `idx_user_meals`: Fast meal lookups by user and date
- `idx_food_search`: Full-text search on food names (GIN index)
- `idx_activities_user_date`: Activity queries by user and date
- `idx_workouts_user_date`: Workout queries by user and date

**Relationship Indexes:**
- Foreign key indexes for all relationships
- Composite indexes for common query patterns

## Service Interactions

### Meal Logging Flow

```
User → MealHandler.Create()
          ↓
      MealService.Create()
          ↓
      ┌───┴────┐
      │        │
      ▼        ▼
  FoodRepo  MealRepo
    .Get()   .Create()
      │        │
      └───┬────┘
          ↓
      Database
```

### AI Meal Parsing Flow

```
User Input (Text/Photo)
          ↓
    MealHandler.ParseText/Photo()
          ↓
    MealParserService
          ↓
    ┌─────┴──────────┐
    │                │
    ▼                ▼
OpenRouter      VisionClient
(DeepSeek)      (Gemini)
    │                │
    └────┬───────────┘
         ↓
   Extract Food Items
         ↓
   ┌─────┴────────┐
   │              │
   ▼              ▼
FoodRepo      AIGeneration
.Search()     (if no match)
   │              │
   └──────┬───────┘
          ↓
    ParsedMeal
```

### Authentication Flow

```
Login Request
      ↓
  AuthHandler.Login()
      ↓
  AuthService.Authenticate()
      ↓
  ┌───┴────┐
  │        │
  ▼        ▼
UserRepo  BCrypt
.GetByEmail() .Compare()
  │        │
  └───┬────┘
      ↓
  JWT.Sign()
      ↓
  Tokens (Access + Refresh)
```

## External Integrations

### OpenRouter API (DeepSeek)

**Purpose**: Text-based meal parsing and nutrition estimation

**Integration Points:**
- Meal text parsing
- AI-generated food nutrition estimation
- Chat assistant responses

**Configuration:**
```go
type OpenRouterClient struct {
    apiKey     string
    baseURL    string
    model      string
    httpClient *http.Client
}
```

**Retry Logic:**
- 3 retry attempts
- Exponential backoff
- Timeout: 30 seconds

### Gemini Vision API

**Purpose**: Food photo analysis

**Integration Points:**
- Photo-based meal parsing
- Food identification from images

**Image Processing:**
- Maximum size: 10MB
- Supported formats: JPEG, PNG
- Preprocessing: Resize if needed

### Supabase Storage

**Purpose**: Photo storage for meal logging

**Integration Points:**
- Upload food photos
- Generate public URLs
- Manage user uploads

## Security Architecture

### Authentication & Authorization

**JWT-Based Authentication:**
```
Access Token (24h) + Refresh Token (7d)
```

**Token Structure:**
```json
{
  "sub": "user_id",
  "email": "user@example.com",
  "exp": 1700000000,
  "iat": 1699913600
}
```

**Authorization Middleware:**
```go
// Extracts and validates JWT from Authorization header
// Sets user context for downstream handlers
// Returns 401 if invalid/expired
```

### Input Validation

**Validation Layers:**
1. **DTO Validation**: Using `validator` tags
2. **Service Validation**: Business rule validation
3. **Repository Validation**: Data integrity checks

**Example:**
```go
type CreateMealRequest struct {
    Name      string    `json:"name" validate:"required,min=1,max=255"`
    MealType  string    `json:"meal_type" validate:"required,oneof=breakfast lunch dinner snack"`
    ConsumedAt time.Time `json:"consumed_at" validate:"required"`
}
```

### Data Protection

- **Password Hashing**: bcrypt (cost factor: 12)
- **SQL Injection**: Prevented by GORM parameterized queries
- **XSS Prevention**: Input sanitization
- **CORS**: Configurable allowed origins
- **Rate Limiting**: Per-user request limits

## Performance Optimizations

### Database Optimizations

**Denormalization:**
- Meal totals (calories, macros) stored in `meals` table
- Avoids expensive JOIN queries for listing

**Indexing Strategy:**
- GIN indexes for full-text search
- Composite indexes for common queries
- Partial indexes for filtered queries

**Connection Pooling:**
```go
sqlDB.SetMaxOpenConns(100)
sqlDB.SetMaxIdleConns(10)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Caching Strategy

**Potential Caching:**
- Food database queries (Redis)
- User profile data (in-memory)
- AI responses (cache similar queries)

### Query Optimization

**Eager Loading:**
```go
db.Preload("FoodItems.Food").Find(&meals)
```

**Pagination:**
- Offset/limit for all list endpoints
- Maximum limit: 100 records

## Scalability Considerations

### Horizontal Scaling

**Stateless Design:**
- No server-side session storage
- JWT for authentication (no session DB)
- Can run multiple API instances behind load balancer

**Database Scaling:**
- Read replicas for read-heavy operations
- Connection pooling per instance
- Prepared for database sharding if needed

### Microservices Path

**Potential Service Split:**
- **Auth Service**: User management, authentication
- **Meal Service**: Meal logging, food database
- **Activity Service**: Activity and workout tracking
- **AI Service**: Meal parsing, chat assistant
- **Analytics Service**: Summaries and insights

### Message Queue Integration

**Future Enhancement:**
- Asynchronous AI processing (RabbitMQ, Kafka)
- Background jobs (photo processing, notifications)
- Event-driven architecture

## Technology Choices Rationale

### Why Go?
- High performance
- Excellent concurrency support
- Strong typing
- Great tooling
- Fast compilation
- Small memory footprint

### Why PostgreSQL?
- ACID compliance
- Excellent full-text search (GIN indexes)
- JSON support
- Mature and reliable
- Strong community

### Why GORM?
- Type-safe
- Automatic migrations
- Relationship handling
- Hooks and callbacks
- Good performance

### Why Gin?
- Fast HTTP router
- Middleware support
- JSON validation
- Active community
- Great documentation

---

For implementation details, see:
- [Main Documentation](./README.md)
- [API Reference](./API.md)
- [Meal Parsing Details](./MEAL_PARSING.md)
