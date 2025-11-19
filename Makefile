.PHONY: help setup migrate-up migrate-down migrate-create seed run test build clean docker-up docker-down docker-logs docker-clean lint fmt vet

# Variables
APP_NAME=fitness-coach-api
BINARY_NAME=bin/$(APP_NAME)
MIGRATION_DIR=migrations
DB_URL=postgres://fitness_user:fitness_password@localhost:5432/fitness_coach?sslmode=disable

# Default target
help:
	@echo "Available targets:"
	@echo "  setup         - Install dependencies and tools"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback last migration"
	@echo "  migrate-create NAME=<name> - Create new migration file"
	@echo "  seed          - Run database seed scripts"
	@echo "  run           - Run the application"
	@echo "  test          - Run all tests"
	@echo "  build         - Build the application binary"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-up     - Start Docker services"
	@echo "  docker-down   - Stop Docker services"
	@echo "  docker-logs   - View Docker logs"
	@echo "  docker-clean  - Remove Docker volumes and containers"
	@echo "  lint          - Run golangci-lint"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"

# Setup development environment
setup:
	@echo "Installing dependencies..."
	go mod download
	go mod verify
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Setup complete!"

# Database migrations
migrate-up:
	@echo "Running migrations..."
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" up
	@echo "Migrations complete!"

migrate-down:
	@echo "Rolling back last migration..."
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" down
	@echo "Rollback complete!"

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=<migration_name>"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)..."
	goose -dir $(MIGRATION_DIR) create $(NAME) sql
	@echo "Migration created!"

migrate-status:
	@echo "Migration status:"
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" status

# Seed database
seed:
	@echo "Seeding database..."
	@if [ -f "scripts/seed.sql" ]; then \
		psql "$(DB_URL)" -f scripts/seed.sql; \
		echo "Database seeded successfully!"; \
	else \
		echo "No seed file found at scripts/seed.sql"; \
	fi

# Run application
run:
	@echo "Starting application..."
	go run cmd/api/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Test coverage:"
	go tool cover -func=coverage.out

test-coverage:
	@echo "Generating coverage report..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build application
build:
	@echo "Building application..."
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-w -s" \
		-o $(BINARY_NAME) \
		./cmd/api
	@echo "Build complete: $(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Docker commands
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Services started! API available at http://localhost:8080"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down
	@echo "Services stopped!"

docker-logs:
	docker-compose logs -f

docker-logs-api:
	docker-compose logs -f api

docker-logs-db:
	docker-compose logs -f postgres

docker-clean:
	@echo "Removing Docker containers and volumes..."
	docker-compose down -v
	@echo "Docker cleanup complete!"

docker-rebuild:
	@echo "Rebuilding Docker images..."
	docker-compose build --no-cache
	docker-compose up -d
	@echo "Rebuild complete!"

# Code quality
lint:
	@echo "Running golangci-lint..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Formatting complete!"

vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "Vet complete!"

# Database access
db-shell:
	@echo "Connecting to database..."
	docker-compose exec postgres psql -U fitness_user -d fitness_coach

db-reset: docker-down docker-clean docker-up
	@echo "Waiting for database to be ready..."
	sleep 5
	$(MAKE) migrate-up
	@echo "Database reset complete!"

# Development helpers
dev: docker-up
	@echo "Starting development environment..."
	@echo "Database: localhost:5432"
	@echo "API: localhost:8080"
	@echo "Run 'make docker-logs' to view logs"

dev-down: docker-down

# Install tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/air-verse/air@latest
	@echo "Tools installed!"

# Hot reload (requires air)
watch:
	@echo "Starting hot reload..."
	air

# Quick start for new developers
quickstart: setup docker-up
	@echo "Waiting for services to start..."
	sleep 5
	$(MAKE) migrate-up
	@echo ""
	@echo "Quick start complete!"
	@echo "API available at http://localhost:8080"
	@echo "Run 'make run' to start the application"
