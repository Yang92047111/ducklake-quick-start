# Development Guide

This guide provides comprehensive information for developers working on the DuckLake Loader project.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Development Setup](#development-setup)
3. [Project Structure](#project-structure)
4. [Building and Running](#building-and-running)
5. [Testing](#testing)
6. [Code Quality](#code-quality)
7. [Database Development](#database-development)
8. [API Development](#api-development)
9. [Docker Development](#docker-development)
10. [Debugging](#debugging)
11. [Contributing Guidelines](#contributing-guidelines)
12. [Troubleshooting](#troubleshooting)

## Prerequisites

### Required
- **Go 1.21+**: [Download here](https://golang.org/dl/)
- **Git**: For version control
- **Make**: For running build tasks

### Optional (for full development experience)
- **Docker & Docker Compose**: For containerized development and testing
- **PostgreSQL**: For local database development (or use Docker)
- **Air**: For hot reloading during development
- **VS Code** with Go extension: Recommended IDE setup

### Installing Optional Tools

```bash
# Install Air for hot reloading
go install github.com/cosmtrek/air@latest

# Install development tools
make install-tools
```

## Development Setup

### 1. Clone and Setup

```bash
# Clone the repository
git clone https://github.com/Yang92047111/ducklake-quick-start.git
cd ducklake-quick-start

# Install dependencies
make deps

# Install development tools
make install-tools

# Verify setup
make test
```

### 2. Environment Configuration

The application supports configuration via environment variables:

```bash
# Database configuration (for PostgreSQL mode)
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=ducklake
export DB_PASSWORD=password
export DB_NAME=ducklake_db
export DB_SSLMODE=disable

# Skip PostgreSQL tests if no database available
export SKIP_POSTGRES_TESTS=true
```

### 3. IDE Setup (VS Code)

Recommended VS Code extensions:
- Go (official Go extension)
- Go Test Explorer
- Docker
- YAML

Workspace settings (`.vscode/settings.json`):
```json
{
    "go.testFlags": ["-v"],
    "go.coverOnSave": true,
    "go.lintOnSave": "package",
    "go.formatTool": "gofmt",
    "go.useLanguageServer": true
}
```

## Project Structure

```
ducklake-quick-start/
├── cmd/
│   └── ducklake-loader/     # Main application entry point
│       ├── main.go
│       └── main_test.go
├── internal/
│   ├── api/                 # REST API handlers and routes
│   │   ├── handlers.go
│   │   └── handlers_test.go
│   ├── loader/              # Data loading and validation logic
│   │   ├── csv_loader.go
│   │   ├── csv_loader_test.go
│   │   ├── json_loader.go
│   │   ├── json_loader_test.go
│   │   ├── validator.go
│   │   ├── validator_test.go
│   │   └── exercise.go
│   └── storage/             # Database repositories and interfaces
│       ├── memory.go
│       ├── memory_test.go
│       ├── postgres.go
│       ├── postgres_test.go
│       └── repository.go
├── test/
│   ├── integration_test.go  # Integration tests
│   └── testdata/           # Sample data files
├── docker-compose.yml      # Production Docker setup
├── docker-compose.test.yml # Testing Docker setup
├── Dockerfile              # Production Docker image
├── Dockerfile.test         # Testing Docker image
├── Makefile               # Build automation
└── README.md              # User documentation
```

### Code Organization Principles

- **cmd/**: Application entry points (main packages)
- **internal/**: Private application code (not importable by other projects)
- **pkg/**: Public library code (if any) - not used in this project
- Each package has its own test files with `_test.go` suffix
- Integration tests are in the `test/` directory

## Building and Running

### Quick Start Commands

```bash
# Build the application
make build

# Run with sample data (in-memory)
make demo

# Run all tests
make test

# Run with coverage
make test-coverage

# Start development server with hot reload
make dev
```

### Manual Building

```bash
# Build for current platform
go build -o bin/ducklake-loader ./cmd/ducklake-loader

# Build for multiple platforms
make release
```

### Running Options

```bash
# Load CSV data with in-memory storage
./bin/ducklake-loader -csv test/testdata/sample_exercises.csv -memory

# Load JSON data and start API server
./bin/ducklake-loader -json test/testdata/sample_exercises.json -memory -server

# Start server with PostgreSQL (requires database setup)
./bin/ducklake-loader -server

# Custom port
./bin/ducklake-loader -server -port 9000
```

## Testing

### Test Categories

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **Database Tests**: Test database operations (PostgreSQL)
4. **API Tests**: Test HTTP endpoints

### Running Tests

```bash
# Run all tests
make test

# Run tests with race detection
make test-unit

# Run with coverage
make test-coverage

# Run integration tests (requires Docker)
make test-integration

# Run specific package tests
go test ./internal/api -v

# Run specific test
go test ./internal/api -run TestHandler_GetExercises -v

# Run benchmarks
make benchmark
```

### Test Environment Variables

```bash
# Skip PostgreSQL tests when no database available
export SKIP_POSTGRES_TESTS=true

# Use test database for PostgreSQL tests
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=test
export DB_PASSWORD=test
export DB_NAME=test_db
```

### Writing Tests

#### Unit Test Example

```go
func TestMyFunction(t *testing.T) {
    // Setup
    input := "test input"
    expected := "expected output"
    
    // Execute
    result := MyFunction(input)
    
    // Verify
    assert.Equal(t, expected, result)
}
```

#### Table-Driven Test Example

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expectError bool
        errorMsg    string
    }{
        {"valid input", "valid", false, ""},
        {"invalid input", "", true, "required"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.input)
            if tt.expectError {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errorMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Code Quality

### Code Quality Checks

```bash
# Run all quality checks
make quality-check

# Individual checks
make format      # Format code
make vet        # Go vet
make lint       # Golangci-lint
make security-scan  # Security analysis
```

### Code Style Guidelines

1. **Formatting**: Use `gofmt` (automated via `make format`)
2. **Naming**: Follow Go naming conventions
   - Use `camelCase` for variables and functions
   - Use `PascalCase` for exported types and functions
   - Use meaningful names
3. **Comments**: 
   - Document all exported functions and types
   - Use complete sentences
   - Start with the name of the thing being documented
4. **Error Handling**: Always handle errors explicitly
5. **Testing**: Maintain high test coverage (aim for >80%)

### Linting Configuration

The project uses `golangci-lint` with configuration in `.golangci.yml`:

- **Enabled linters**: govet, gocyclo, misspell, goimports, gocritic, etc.
- **Line length**: 140 characters max
- **Complexity**: Max cyclomatic complexity of 15

## Database Development

### Local PostgreSQL Setup

#### Option 1: Docker (Recommended)

```bash
# Start PostgreSQL database
make db-up

# Stop database
make db-down

# Reset database
make db-reset
```

#### Option 2: Local Installation

```bash
# Install PostgreSQL (macOS)
brew install postgresql
brew services start postgresql

# Create database and user
psql postgres
CREATE DATABASE ducklake_db;
CREATE USER ducklake WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE ducklake_db TO ducklake;
```

### Database Schema

The application auto-creates tables on startup:

```sql
CREATE TABLE IF NOT EXISTS exercises (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    duration INTEGER NOT NULL,
    calories INTEGER NOT NULL,
    date DATE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Repository Pattern

The project uses the repository pattern for data access:

```go
type ExerciseRepository interface {
    Insert(exercise Exercise) error
    InsertBatch(exercises []Exercise) error
    GetByID(id int) (*Exercise, error)
    GetByType(exerciseType string) ([]Exercise, error)
    GetByDateRange(start, end time.Time) ([]Exercise, error)
    GetAll() ([]Exercise, error)
    Update(exercise Exercise) error
    Delete(id int) error
    Close() error
}
```

## API Development

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/exercises` | Get all exercises |
| GET | `/exercises/{id}` | Get exercise by ID |
| GET | `/exercises/type/{type}` | Get exercises by type |
| GET | `/exercises/date-range?start=YYYY-MM-DD&end=YYYY-MM-DD` | Get exercises by date range |
| GET | `/health` | Health check |
| GET | `/ready` | Readiness check |

### Testing API Endpoints

```bash
# Start the server
make demo

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/exercises
curl http://localhost:8080/exercises/1
curl http://localhost:8080/exercises/type/cardio
curl "http://localhost:8080/exercises/date-range?start=2024-01-01&end=2024-01-31"
```

### Error Handling

The API returns structured JSON errors:

```json
{
  "error": "Bad Request",
  "code": 400,
  "message": "Exercise ID must be a positive integer"
}
```

### Adding New Endpoints

1. Define handler function in `internal/api/handlers.go`
2. Add route in `SetupRoutes()` method
3. Add comprehensive tests in `internal/api/handlers_test.go`
4. Update documentation

## Docker Development

### Development with Docker

```bash
# Build Docker image
make docker-build

# Start all services
make docker-run

# View logs
docker compose logs -f app

# Stop services
make docker-down
```

### Docker Configuration

- **Dockerfile**: Production-ready multi-stage build
- **docker-compose.yml**: Application + PostgreSQL
- **docker-compose.test.yml**: Testing environment

### Debugging Docker Issues

```bash
# Check container status
docker compose ps

# View logs
docker compose logs app
docker compose logs postgres

# Execute commands in container
docker compose exec app sh
docker compose exec postgres psql -U ducklake -d ducklake_db
```

## Debugging

### Go Debugging with Delve

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug main application
dlv debug ./cmd/ducklake-loader

# Debug tests
dlv test ./internal/api
```

### VS Code Debugging

Add to `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Program",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/ducklake-loader",
            "args": ["-memory", "-server"]
        },
        {
            "name": "Debug Test",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}/internal/api"
        }
    ]
}
```

### Logging and Monitoring

Enable verbose logging by setting environment variables:

```bash
export LOG_LEVEL=debug
export LOG_FORMAT=json
```

The application provides health endpoints for monitoring:
- `/health`: Basic health check
- `/ready`: Readiness check (includes database connectivity)

## Contributing Guidelines

### Development Workflow

1. **Fork and Clone**: Fork the repository and clone your fork
2. **Branch**: Create a feature branch (`git checkout -b feature/my-feature`)
3. **Develop**: Make your changes following the code style guidelines
4. **Test**: Ensure all tests pass (`make quality-check`)
5. **Commit**: Make atomic commits with clear messages
6. **Push**: Push to your fork (`git push origin feature/my-feature`)
7. **Pull Request**: Create a pull request with a clear description

### Commit Message Format

```
type(scope): brief description

Longer description if needed

- Key changes
- Breaking changes (if any)
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `chore`

### Pull Request Checklist

- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`make format`)
- [ ] Linting passes (`make lint`)
- [ ] Security scan passes (`make security-scan`)
- [ ] Documentation updated (if needed)
- [ ] Test coverage maintained or improved

### Code Review Guidelines

1. **Functionality**: Does the code work as intended?
2. **Testing**: Are there adequate tests?
3. **Performance**: Are there any performance implications?
4. **Security**: Are there any security concerns?
5. **Maintainability**: Is the code readable and maintainable?

## Troubleshooting

### Common Issues

#### Build Issues

**Problem**: `go build` fails with dependency errors
```bash
# Solution: Clean and reinstall dependencies
make clean
make deps
go mod verify
```

**Problem**: Binary not found after build
```bash
# Solution: Check build directory
ls -la bin/
# Rebuild
make build
```

#### Test Issues

**Problem**: PostgreSQL tests fail
```bash
# Solution: Skip PostgreSQL tests
export SKIP_POSTGRES_TESTS=true
make test
```

**Problem**: Tests timeout
```bash
# Solution: Run with verbose output to identify slow tests
go test -v -timeout 60s ./...
```

#### Database Issues

**Problem**: Cannot connect to PostgreSQL
```bash
# Check if PostgreSQL is running
docker compose ps
# Check logs
docker compose logs postgres
# Restart database
make db-reset
```

**Problem**: Permission denied errors
```bash
# Check database user permissions
docker compose exec postgres psql -U ducklake -d ducklake_db -c "\du"
```

#### Docker Issues

**Problem**: Docker build fails
```bash
# Check Docker daemon
docker version
# Clean Docker cache
docker system prune -f
# Rebuild
make docker-build
```

**Problem**: Port conflicts
```bash
# Check what's using the port
lsof -i :8080
# Use different port
./bin/ducklake-loader -server -port 9000
```

### Getting Help

1. **Check the logs**: Most issues are visible in application logs
2. **Read error messages**: Go provides detailed error messages
3. **Check the documentation**: This guide and README.md
4. **Run diagnostics**: Use `make health-check` to verify setup
5. **Ask for help**: Create an issue with detailed information

### Debugging Commands

```bash
# Check Go environment
go env

# Verify module dependencies
go mod verify

# List all available make targets
make help

# Check application health
make health-check

# View detailed build output
go build -v ./cmd/ducklake-loader

# Run specific test with verbose output
go test -v -run TestSpecificFunction ./internal/package
```

## Performance Considerations

### Optimization Tips

1. **Database**: Use appropriate indexes for query patterns
2. **Memory**: Consider using database storage for large datasets
3. **Concurrency**: The application is designed to be concurrent-safe
4. **Caching**: Consider adding caching layer for read-heavy workloads

### Monitoring

- Use `/health` and `/ready` endpoints for health checks
- Monitor response times and error rates
- Track database connection pool usage
- Set up proper logging in production

---

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Go Testing](https://golang.org/doc/tutorial/add-a-test)
- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Gorilla Mux Documentation](https://github.com/gorilla/mux)
- [Testify Documentation](https://github.com/stretchr/testify)

For questions or issues not covered in this guide, please create an issue in the repository.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Data Files    │    │   Application   │    │    Storage      │
│                 │    │                 │    │                 │
│ • CSV Files     │───▶│ • Data Loader   │───▶│ • PostgreSQL    │
│ • JSON Files    │    │ • Validator     │    │ • In-Memory     │
│ • API Sources   │    │ • API Handlers  │    │ • Repository    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   REST API      │
                       │                 │
                       │ • GET /exercises│
                       │ • GET /exercises/{id}
                       │ • GET /exercises/type/{type}
                       │ • GET /exercises/date-range
                       └─────────────────┘
```

## Performance & Scalability

- **Memory Usage**: Efficient with streaming JSON/CSV parsing
- **Database**: Connection pooling ready for PostgreSQL
- **Concurrency**: Thread-safe in-memory repository with mutex
- **API**: Stateless design, horizontally scalable
- **Docker**: Multi-stage builds for minimal image size

## Production Readiness Checklist

- ✅ Error handling and logging
- ✅ Input validation and sanitization
- ✅ Database connection management
- ✅ Configuration via environment variables
- ✅ Health checks (implicit via database ping)
- ✅ Graceful shutdown (via defer statements)
- ✅ Security best practices (no hardcoded secrets)
- ✅ Comprehensive testing
- ✅ Documentation and examples

## Next Steps (Optional Enhancements)

1. **Metrics & Monitoring**
   - Add Prometheus metrics
   - Health check endpoint
   - Structured logging with levels

2. **Advanced Features**
   - Batch processing with queues
   - Data transformation pipelines
   - API authentication/authorization
   - Rate limiting

3. **Database Enhancements**
   - Database migrations with golang-migrate
   - Connection pooling configuration
   - Read replicas support

4. **Deployment**
   - Kubernetes manifests
   - Helm charts
   - Production Docker Compose

## Conclusion

The DuckLake Loader is a production-ready Go application that successfully implements all the requirements from the original specification. It demonstrates best practices in Go development, testing, and containerization.