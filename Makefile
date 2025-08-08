.PHONY: build test test-coverage clean run docker-build docker-run

# Build the application
build:
	go build -o bin/ducklake-loader ./cmd/ducklake-loader

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run the application
run: build
	./bin/ducklake-loader

# Install dependencies
deps:
	go mod tidy
	go mod download

# Lint code
lint:
	golangci-lint run

# Docker build
docker-build:
	docker build -t ducklake-loader .

# Docker run with compose
docker-run:
	docker compose up --build -d

# Docker down
docker-down:
	docker compose down

# Run integration tests
test-integration:
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker compose -f docker-compose.test.yml down