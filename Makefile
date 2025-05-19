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

# Help target
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build the application
build: ## Build the Go application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@templ generate
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Run the application
run: build ## Run the application
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run with hot-reloading using Air
dev: ## Run the application with hot-reloading (Air)
	@echo "Starting development server with hot-reloading..."
	@air

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
	@templ generate

# Install dependencies
deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go install github.com/a-h/templ/cmd/templ@latest 