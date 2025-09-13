.PHONY: build test test-coverage test-unit test-integration clean run docker-build docker-run docker-down deps lint format vet security-scan help

# Build configuration
BINARY_NAME=ducklake-loader
BUILD_DIR=bin
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Go configuration
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOFMT=gofmt

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ducklake-loader

# Run all tests
test: test-unit

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) ./... -v -race

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) ./... -v -race -coverprofile=$(COVERAGE_FILE)
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# Run integration tests with Docker
test-integration:
	@echo "Running integration tests with Docker..."
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker compose -f docker-compose.test.yml down

# Clean build artifacts and coverage files
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)/
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -f cmd_coverage.out api_coverage.out loader_coverage.out

# Run the application
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Install and update dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download
	$(GOMOD) verify

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2
	@which gosec > /dev/null || $(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Lint code using golangci-lint
lint: install-tools
	@echo "Running golangci-lint..."
	golangci-lint run

# Format code
format:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	$(GOCMD) mod tidy

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Run security analysis
security-scan: install-tools
	@echo "Running security scan..."
	gosec ./...

# Run all quality checks
quality-check: format vet lint security-scan test-coverage
	@echo "All quality checks completed successfully!"

# Docker targets
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

# Start application with Docker Compose
docker-run:
	@echo "Starting application with Docker Compose..."
	docker compose up --build -d

# Stop Docker Compose services
docker-down:
	@echo "Stopping Docker Compose services..."
	docker compose down

# Development server (with file watching if available)
dev: build
	@echo "Starting development server..."
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		echo "For hot reload, install air: go install github.com/cosmtrek/air@latest"; \
		./$(BUILD_DIR)/$(BINARY_NAME) -memory -server; \
	fi

# Database operations
db-up:
	@echo "Starting PostgreSQL database..."
	docker compose up postgres -d

db-down:
	@echo "Stopping PostgreSQL database..."
	docker compose stop postgres

db-reset: db-down db-up
	@echo "Database reset completed"

# Load sample data
load-sample-csv: build
	@echo "Loading sample CSV data..."
	./$(BUILD_DIR)/$(BINARY_NAME) -csv test/testdata/sample_exercises.csv -memory

load-sample-json: build
	@echo "Loading sample JSON data..."
	./$(BUILD_DIR)/$(BINARY_NAME) -json test/testdata/sample_exercises.json -memory

# Start API server with sample data
demo: build
	@echo "Starting demo server with sample data..."
	./$(BUILD_DIR)/$(BINARY_NAME) -json test/testdata/sample_exercises.json -memory -server

# Performance benchmarks
benchmark:
	@echo "Running performance benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Generate documentation
docs:
	@echo "Generating documentation..."
	@which godoc > /dev/null || $(GOGET) golang.org/x/tools/cmd/godoc@latest
	@echo "Documentation server will be available at http://localhost:6060"
	godoc -http=:6060

# Release build (optimized)
release: clean
	@echo "Building release version..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/ducklake-loader
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/ducklake-loader
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/ducklake-loader

# Health check for running application
health-check:
	@echo "Checking application health..."
	@curl -f http://localhost:8080/health || echo "Application not running or unhealthy"

# Show help
help:
	@echo "DuckLake Loader Makefile"
	@echo "========================"
	@echo ""
	@echo "Available targets:"
	@echo "  build             Build the application"
	@echo "  test              Run all tests"
	@echo "  test-unit         Run unit tests only"
	@echo "  test-integration  Run integration tests with Docker"
	@echo "  test-coverage     Run tests with coverage report"
	@echo "  clean             Clean build artifacts"
	@echo "  run               Build and run the application"
	@echo "  deps              Install and update dependencies"
	@echo "  install-tools     Install development tools"
	@echo "  lint              Run golangci-lint"
	@echo "  format            Format code with gofmt"
	@echo "  vet               Run go vet"
	@echo "  security-scan     Run security analysis"
	@echo "  quality-check     Run all quality checks"
	@echo "  docker-build      Build Docker image"
	@echo "  docker-run        Start with Docker Compose"
	@echo "  docker-down       Stop Docker Compose"
	@echo "  dev               Start development server"
	@echo "  db-up             Start PostgreSQL database"
	@echo "  db-down           Stop PostgreSQL database"
	@echo "  db-reset          Reset database"
	@echo "  load-sample-csv   Load sample CSV data"
	@echo "  load-sample-json  Load sample JSON data"
	@echo "  demo              Start demo server with sample data"
	@echo "  benchmark         Run performance benchmarks"
	@echo "  docs              Start documentation server"
	@echo "  release           Build release binaries for multiple platforms"
	@echo "  health-check      Check if application is running"
	@echo "  help              Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make test-coverage"
	@echo "  make demo"
	@echo "  make quality-check"