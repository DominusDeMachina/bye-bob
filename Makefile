.PHONY: build run dev docker-build docker-run docker-dev clean test lint help

# Default target
.DEFAULT_GOAL := help

# Binary name and build directory
BINARY_NAME=byebob
BUILD_DIR=./tmp

# Version info
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Paths to binaries
TEMPL=$(HOME)/go/bin/templ
AIR=$(HOME)/go/bin/air

# Help target
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build the application
build: ## Build the Go application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(TEMPL) generate
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Run the application
run: build ## Run the application
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run with hot-reloading using Air
dev: ## Run the application with hot-reloading (Air)
	@echo "Starting development server with hot-reloading..."
	@$(AIR)

# Build Docker image
docker-build: ## Build the Docker image
	@echo "Building Docker image..."
	@docker build -t byebob .

# Run with Docker
docker-run: docker-build ## Run the application in Docker
	@echo "Running Docker container..."
	@docker run -p 3000:3000 --env-file .env byebob

# Run development environment with Docker Compose
docker-dev: ## Run the development environment with Docker Compose
	@echo "Starting development environment with Docker Compose..."
	@docker-compose -f docker-compose.dev.yml up --build

# Run production environment with Docker Compose
docker-prod: ## Run the production environment with Docker Compose
	@echo "Starting production environment with Docker Compose..."
	@docker-compose up --build

# Stop Docker Compose services
docker-stop: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	@docker-compose down
	@docker-compose -f docker-compose.dev.yml down

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Run tests
test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

# Run linting
lint: ## Run linting
	@echo "Running linter..."
	@if [ -z "$$(which golangci-lint)" ]; then \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@golangci-lint run ./...

# Generate Templ templates
templ: ## Generate Templ templates
	@echo "Generating Templ templates..."
	@$(TEMPL) generate

# Install dependencies
deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Test database connection
test-db: ## Test database connection
	@echo "Testing database connection..."
	@go run ./scripts/cmd/test_db
	
# Test Railway connection specifically
test-railway: ## Test Railway connection
	@echo "Testing Railway connection (make sure your Railway credentials are set)..."
	@go run ./scripts/cmd/test_railway
	
# Test local database connection
test-local: ## Test local database connection
	@echo "Testing local database connection..."
	@go run ./scripts/cmd/test_local

# Setup database migrations directory
setup-migrations: ## Setup migrations directory structure
	@echo "Setting up migrations directory..."
	@mkdir -p migrations/postgres
	@echo "Created migrations/postgres directory for database migrations"

# Create a new migration
migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "Error: Migration name is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)..."
	@migrate create -ext sql -dir migrations/postgres -seq $(name)
	@echo "Created migration files in migrations/postgres"

# Run migrations up
migrate-up: ## Run all pending migrations
	@echo "Running migrations up..."
	@if [ -z "$(POSTGRESQL_URL)" ]; then \
		echo "Error: POSTGRESQL_URL environment variable is required"; \
		exit 1; \
	fi
	@migrate -path migrations/postgres -database "$(POSTGRESQL_URL)" up
	@echo "Migrations applied successfully"

# Roll back last migration
migrate-down: ## Roll back the last migration
	@echo "Rolling back the last migration..."
	@if [ -z "$(POSTGRESQL_URL)" ]; then \
		echo "Error: POSTGRESQL_URL environment variable is required"; \
		exit 1; \
	fi
	@migrate -path migrations/postgres -database "$(POSTGRESQL_URL)" down 1
	@echo "Last migration rolled back"

# Set up database security configuration
db-security-config: ## Set up database security configuration (roles, permissions, RLS)
	@echo "Applying database security configuration..."
	@if [ -z "$(POSTGRESQL_URL)" ]; then \
		echo "Error: POSTGRESQL_URL environment variable is required"; \
		exit 1; \
	fi
	@migrate -path migrations/postgres -database "$(POSTGRESQL_URL)" up 4
	@echo "Security configuration applied successfully!"
	@echo "IMPORTANT: Update database user passwords for your environment using scripts/db/update_db_passwords.sh"

# Update database user passwords
db-update-passwords: ## Update database user passwords (usage: make db-update-passwords env=dev|staging|prod)
	@if [ -z "$(env)" ]; then \
		echo "Error: Environment is required. Usage: make db-update-passwords env=dev|staging|prod"; \
		exit 1; \
	fi
	@echo "Updating database passwords for $(env) environment..."
	@scripts/db/update_db_passwords.sh $(env)