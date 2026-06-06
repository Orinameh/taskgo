# Load .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: build start dev clean test install run help migrate-up migrate-down migrate-status migrate-create

BINARY=bin/api
MAIN=./cmd/api

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed"

# Build production binary
build: install
	@echo "🔨 Building production binary..."
	go build -ldflags="-s -w" -o $(BINARY) $(MAIN)
	@echo "✅ Build complete: $(BINARY)"

# Run production binary (requires build first)
start:
	@echo "🚀 Starting production server..."
	@if [ ! -f $(BINARY) ]; then \
		echo "❌ Binary not found. Run 'make build' first."; \
		exit 1; \
	fi
	./$(BINARY)

# Development mode (runs source directly)
dev:
	@echo "🔄 Starting development server..."
	go run $(MAIN)

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache
	@echo "✅ Clean complete"

# Build and run (one command)
run: build start

# Run tests with coverage
test:
	@echo "🧪 Running tests..."
	go test -v ./...

test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Run tests with race detector
test-race:
	@echo "🏃 Running tests with race detector..."
	go test -race -v ./...

# Run integration tests
test-integration:
	@echo "🔗 Running integration tests..."
	go test -v -tags=integration ./test/...

# Run specific package tests
test-auth:
	go test -v ./internal/auth/...

test-repository:
	go test -v ./internal/repository/...

test-service:
	go test -v ./internal/service/...

test-handler:
	go test -v ./internal/handler/...

# Goose migration commands
migrate-up:
	@echo "📦 Running migrations..."
	goose -dir internal/migrations postgres "$(DATABASE_URL)" up
	@echo "✅ Migrations complete"

migrate-down:
	@echo "⚠️ Rolling back last migration..."
	goose -dir internal/migrations postgres "$(DATABASE_URL)" down
	@echo "✅ Rollback complete"

migrate-status:
	@echo "📊 Migration status..."
	goose -dir internal/migrations postgres "$(DATABASE_URL)" status

migrate-create:
	@echo "📝 Creating new migration..."
	goose -dir internal/migrations create $(NAME) sql


# Show help
help:
	@echo "Available commands:"
	@echo "  make install      - Install dependencies"
	@echo "  make build        - Build production binary"
	@echo "  make start        - Run production binary"
	@echo "  make dev          - Run in development mode"
	@echo "  make test         - Run all tests"
	@echo "  make test-coverage- Run tests with coverage"
	@echo "  make test-race    - Run tests with race detector"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-auth    - Run auth package tests"
	@echo "  make test-repository - Run repository tests"
	@echo "  make test-service - Run service tests"
	@echo "  make test-handler - Run handler tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make migrate-up   - Run all pending migrations"
	@echo "  make migrate-down - Rollback last migration"
	@echo "  make migrate-status - Show migration status"
	@echo "  make migrate-create NAME=name - Create new migration"
	@echo "  make run          - Build and start"
	@echo "  make help         - Show this help"

# Create a new migration file
# make migrate-create NAME=add_tasks_index

# Edit the generated file in migrations/ folder
# Then run it
# make migrate-up