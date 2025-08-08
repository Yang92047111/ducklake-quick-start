# DuckLake Loader

A Go application that reads exercise data from DuckLake, processes it, and stores it in a database with optional REST API access.

## Features

- **Data Loading**: Parse CSV and JSON exercise data files
- **Data Validation**: Validate exercise records for completeness and correctness
- **Storage Options**: PostgreSQL database or in-memory storage for development
- **REST API**: Query exercise data via HTTP endpoints
- **Docker Support**: Containerized deployment with Docker Compose
- **Comprehensive Testing**: Unit tests, integration tests, and coverage reporting

## Project Structure

```
ducklake-loader/
├── cmd/ducklake-loader/     # Main application entry point
├── internal/
│   ├── loader/              # Data loading and validation
│   ├── storage/             # Database repositories
│   └── api/                 # REST API handlers
├── test/                    # Test files and sample data
├── Dockerfile               # Container build configuration
├── docker-compose.yml       # Multi-service deployment
└── Makefile                 # Build automation
```

## Quick Start

### Prerequisites
- Go 1.21 or later
- Docker and Docker Compose (for containerized deployment)
- PostgreSQL (if running without Docker)

### Local Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Build the application
make build

# Load sample CSV data (in-memory storage)
./bin/ducklake-loader -csv test/testdata/sample_exercises.csv -memory

# Start API server with sample data
./bin/ducklake-loader -json test/testdata/sample_exercises.json -memory -server
```

### Docker Deployment

```bash
# Start application with PostgreSQL
make docker-run

# Run integration tests
make test-integration

# Stop services
make docker-down
```

## Usage

### Command Line Options

- `-csv <file>` - Load data from CSV file
- `-json <file>` - Load data from JSON file  
- `-memory` - Use in-memory storage (default: PostgreSQL)
- `-server` - Start REST API server
- `-port <port>` - Server port (default: 8080)

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/exercises` | Get all exercises |
| GET | `/exercises/{id}` | Get exercise by ID |
| GET | `/exercises/type/{type}` | Get exercises by type |
| GET | `/exercises/date-range?start=YYYY-MM-DD&end=YYYY-MM-DD` | Get exercises by date range |

### Example API Calls

```bash
# Get all exercises
curl http://localhost:8080/exercises

# Get exercise by ID
curl http://localhost:8080/exercises/1

# Get cardio exercises
curl http://localhost:8080/exercises/type/cardio

# Get exercises in date range
curl "http://localhost:8080/exercises/date-range?start=2024-01-15&end=2024-01-16"
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_USER` | Database user | ducklake |
| `DB_PASSWORD` | Database password | password |
| `DB_NAME` | Database name | ducklake_db |
| `DB_SSLMODE` | SSL mode | disable |

## Development

### Running Tests

```bash
# Unit tests
make test

# Tests with coverage
make test-coverage

# Integration tests with Docker
make test-integration
```

### Data Format

#### CSV Format
```csv
id,name,type,duration,calories,date,description
1,Morning Run,cardio,30,300,2024-01-15,Easy morning jog
```

#### JSON Format
```json
[
  {
    "id": 1,
    "name": "Morning Run",
    "type": "cardio", 
    "duration": 30,
    "calories": 300,
    "date": "2024-01-15T00:00:00Z",
    "description": "Easy morning jog"
  }
]
```

## License

MIT License
