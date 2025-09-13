# DuckLake Tutorial: Complete Guide

Welcome to DuckLake! This comprehensive tutorial will guide you through everything from basic setup to advanced enterprise features. Whether you're a beginner or an experienced developer, this guide will help you unlock DuckLake's full potential.

## Table of Contents

1. [Getting Started](#1-getting-started)
2. [Basic Usage](#2-basic-usage)
3. [API Fundamentals](#3-api-fundamentals)
4. [Data Management](#4-data-management)
5. [Enterprise Lakehouse Features](#5-enterprise-lakehouse-features)
6. [Advanced Scenarios](#6-advanced-scenarios)
7. [Production Deployment](#7-production-deployment)
8. [Development Guide](#8-development-guide)
9. [Troubleshooting](#9-troubleshooting)
10. [Best Practices](#10-best-practices)

---

## 1. Getting Started

### What You'll Learn
- Install and set up DuckLake
- Understand the project structure
- Run your first DuckLake commands
- Load sample data

### Prerequisites

Before we begin, ensure you have:
- **Go 1.21+**: [Download here](https://golang.org/dl/)
- **Git**: For cloning the repository
- **Docker & Docker Compose**: For containerized deployment (optional)
- **Basic terminal/command-line knowledge**

### Step 1: Installation

```bash
# Clone the repository
git clone https://github.com/Yang92047111/ducklake-quick-start.git
cd ducklake-quick-start

# Install dependencies
make deps
```

### Step 2: Verify Installation

```bash
# Run tests to ensure everything is working
make test

# Build the application
make build

# Verify the binary was created
ls -la bin/
```

You should see `ducklake-loader` in the `bin/` directory.

### Step 3: Your First DuckLake Command

Let's start with a simple data loading example:

```bash
# Load sample CSV data using in-memory storage
./bin/ducklake-loader -csv test/testdata/sample_exercises.csv -memory
```

🎉 **Congratulations!** You've successfully loaded your first data into DuckLake.

### Understanding the Output

DuckLake will show you:
- Number of records loaded
- Validation results
- Any errors encountered

---

## 2. Basic Usage

### What You'll Learn
- Load different data formats
- Use various storage options
- Understand command-line flags
- Work with sample data

### Data Formats

DuckLake supports two primary data formats:

#### CSV Format
```csv
id,name,type,duration,calories,date,description
1,Morning Run,cardio,30,300,2024-01-15,Easy morning jog
2,Push-ups,strength,15,100,2024-01-15,3 sets of 15 push-ups
3,Yoga,flexibility,45,150,2024-01-15,Morning yoga session
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
  },
  {
    "id": 2,
    "name": "Push-ups",
    "type": "strength",
    "duration": 15,
    "calories": 100,
    "date": "2024-01-15T00:00:00Z",
    "description": "3 sets of 15 push-ups"
  }
]
```

### Loading Data

#### From CSV Files
```bash
# Load CSV with in-memory storage (fastest, great for testing)
./bin/ducklake-loader -csv test/testdata/sample_exercises.csv -memory

# Load CSV with PostgreSQL storage (requires database setup)
./bin/ducklake-loader -csv test/testdata/sample_exercises.csv
```

#### From JSON Files
```bash
# Load JSON with in-memory storage
./bin/ducklake-loader -json test/testdata/sample_exercises.json -memory

# Load JSON with PostgreSQL storage
./bin/ducklake-loader -json test/testdata/sample_exercises.json
```

### Storage Options

DuckLake offers three storage backends:

1. **In-Memory Storage** (`-memory`)
   - Fastest performance
   - Perfect for development and testing
   - Data lost when application closes
   - No external dependencies

2. **PostgreSQL Storage** (default)
   - Production-ready persistence
   - ACID compliance
   - Advanced querying capabilities
   - Requires PostgreSQL setup

3. **Lakehouse Storage** (`-lakehouse`)
   - Enterprise-grade features
   - Version control and time travel
   - Schema evolution
   - ACID transactions with Delta Lake-like storage

### Command-Line Flags Reference

```bash
# Data source flags
-csv <file>          # Load from CSV file
-json <file>         # Load from JSON file

# Storage flags
-memory              # Use in-memory storage
-lakehouse           # Use lakehouse storage
-lakehouse-path <path> # Set lakehouse data directory

# Server flags
-server              # Start REST API server
-port <port>         # Set server port (default: 8080)

# Examples
./bin/ducklake-loader -csv data.csv -memory -server -port 9000
./bin/ducklake-loader -json data.json -lakehouse -server
```

---

## 3. API Fundamentals

### What You'll Learn
- Start the REST API server
- Make your first API calls
- Understand API endpoints
- Work with exercise data via HTTP

### Starting the API Server

```bash
# Start server with in-memory storage and sample data
./bin/ducklake-loader -json test/testdata/sample_exercises.json -memory -server

# The server will start on http://localhost:8080
```

You should see output like:
```
Loading JSON data from test/testdata/sample_exercises.json
Successfully loaded 3 exercises
Starting server on :8080
```

### Basic API Endpoints

#### 1. Health Check
```bash
# Verify the server is running
curl http://localhost:8080/health

# Expected response: 200 OK
```

#### 2. Get All Exercises
```bash
# Retrieve all exercises
curl http://localhost:8080/exercises

# Pretty print JSON response
curl http://localhost:8080/exercises | jq
```

#### 3. Get Exercise by ID
```bash
# Get a specific exercise
curl http://localhost:8080/exercises/1

# Expected response
{
  "id": 1,
  "name": "Morning Run",
  "type": "cardio",
  "duration": 30,
  "calories": 300,
  "date": "2024-01-15T00:00:00Z",
  "description": "Easy morning jog"
}
```

#### 4. Filter by Exercise Type
```bash
# Get all cardio exercises
curl http://localhost:8080/exercises/type/cardio

# Get all strength exercises
curl http://localhost:8080/exercises/type/strength

# Get all flexibility exercises
curl http://localhost:8080/exercises/type/flexibility
```

#### 5. Date Range Queries
```bash
# Get exercises within a date range
curl "http://localhost:8080/exercises/date-range?start=2024-01-15&end=2024-01-16"

# Get exercises for January 2024
curl "http://localhost:8080/exercises/date-range?start=2024-01-01&end=2024-01-31"
```

### API Response Format

All API responses follow a consistent format:

**Success Response:**
```json
{
  "id": 1,
  "name": "Morning Run",
  "type": "cardio",
  "duration": 30,
  "calories": 300,
  "date": "2024-01-15T00:00:00Z",
  "description": "Easy morning jog"
}
```

**Error Response:**
```json
{
  "error": "Bad Request",
  "code": 400,
  "message": "Exercise ID must be a positive integer"
}
```

### Testing with curl

Create a test script for easy API testing:

```bash
#!/bin/bash
# save as test_api.sh

BASE_URL="http://localhost:8080"

echo "Testing DuckLake API..."

echo "1. Health check:"
curl -s $BASE_URL/health
echo -e "\n"

echo "2. Get all exercises:"
curl -s $BASE_URL/exercises | jq '.[0]'
echo -e "\n"

echo "3. Get exercise by ID:"
curl -s $BASE_URL/exercises/1 | jq
echo -e "\n"

echo "4. Get cardio exercises:"
curl -s $BASE_URL/exercises/type/cardio | jq 'length'
echo " cardio exercises found"
```

Run it:
```bash
chmod +x test_api.sh
./test_api.sh
```

---

## 4. Data Management

### What You'll Learn
- Create and validate exercise data
- Understand data validation rules
- Handle different exercise types
- Work with dates and times

### Exercise Data Structure

Every exercise in DuckLake has these fields:

```go
type Exercise struct {
    ID          int       // Unique identifier
    Name        string    // Exercise name (e.g., "Morning Run")
    Type        string    // Exercise type: "cardio", "strength", or "flexibility"
    Duration    int       // Duration in minutes
    Calories    int       // Calories burned
    Date        time.Time // Exercise date
    Description string    // Optional description
}
```

### Exercise Types

DuckLake supports three exercise categories:

#### 1. Cardio
- **Purpose**: Cardiovascular health and endurance
- **Examples**: Running, swimming, cycling, walking
- **Typical Duration**: 20-60 minutes
- **Calorie Range**: 200-800 calories

#### 2. Strength
- **Purpose**: Muscle building and strength training
- **Examples**: Weight training, push-ups, squats, resistance exercises
- **Typical Duration**: 15-90 minutes
- **Calorie Range**: 100-500 calories

#### 3. Flexibility
- **Purpose**: Mobility, stretching, and relaxation
- **Examples**: Yoga, stretching, pilates, meditation
- **Typical Duration**: 10-90 minutes
- **Calorie Range**: 50-300 calories

### Data Validation Rules

DuckLake enforces strict validation to ensure data quality:

#### Required Fields
- **ID**: Must be a positive integer and unique
- **Name**: Cannot be empty
- **Type**: Must be "cardio", "strength", or "flexibility" (case-insensitive)
- **Duration**: Must be positive (minutes)
- **Calories**: Must be positive
- **Date**: Must be a valid date

#### Validation Examples

**✅ Valid Exercise:**
```json
{
  "id": 1,
  "name": "Morning Run",
  "type": "cardio",
  "duration": 30,
  "calories": 300,
  "date": "2024-01-15T00:00:00Z",
  "description": "Easy 5K jog"
}
```

**❌ Invalid Exercise (missing required fields):**
```json
{
  "id": 0,
  "name": "",
  "type": "invalid",
  "duration": -10,
  "calories": 0
}
```

### Creating Sample Data

#### CSV Sample (`my_exercises.csv`)
```csv
id,name,type,duration,calories,date,description
1,Morning Run,cardio,30,300,2024-01-15,5K jog around the park
2,Weight Training,strength,45,250,2024-01-15,Upper body workout
3,Evening Yoga,flexibility,60,150,2024-01-15,Relaxing yoga session
4,Cycling,cardio,90,600,2024-01-16,Weekend bike ride
5,Push-up Challenge,strength,20,120,2024-01-16,100 push-ups total
```

#### JSON Sample (`my_exercises.json`)
```json
[
  {
    "id": 1,
    "name": "Morning Run",
    "type": "cardio",
    "duration": 30,
    "calories": 300,
    "date": "2024-01-15T08:00:00Z",
    "description": "5K jog around the park"
  },
  {
    "id": 2,
    "name": "Weight Training",
    "type": "strength",
    "duration": 45,
    "calories": 250,
    "date": "2024-01-15T18:00:00Z",
    "description": "Upper body workout"
  },
  {
    "id": 3,
    "name": "Evening Yoga",
    "type": "flexibility",
    "duration": 60,
    "calories": 150,
    "date": "2024-01-15T20:00:00Z",
    "description": "Relaxing yoga session"
  }
]
```

### Loading Your Own Data

```bash
# Load your CSV file
./bin/ducklake-loader -csv my_exercises.csv -memory -server

# Load your JSON file
./bin/ducklake-loader -json my_exercises.json -memory -server
```

### Date Formats

DuckLake accepts various date formats:

- **ISO 8601**: `2024-01-15T08:00:00Z` (recommended)
- **Date only**: `2024-01-15`
- **US format**: `01/15/2024`
- **European format**: `15/01/2024`

---

## 5. Enterprise Lakehouse Features

### What You'll Learn
- Enable lakehouse mode
- Use versioning and time travel
- Manage schemas and transactions
- Implement data constraints
- Monitor data quality

### Introduction to Lakehouse

The DuckLake Lakehouse provides enterprise-grade data management capabilities:

- **🕐 Data Versioning**: Every change creates a new version
- **⏰ Time Travel**: Query historical data states
- **🔄 Schema Evolution**: Add columns without breaking compatibility
- **🔒 ACID Transactions**: Full transaction support
- **📋 Data Constraints**: Enforce data quality rules
- **📊 Data Quality Metrics**: Monitor and validate data

### Enabling Lakehouse Mode

#### Build the Lakehouse Binary
```bash
# Build the lakehouse-enabled binary
go build -o bin/ducklake-lakehouse cmd/ducklake-lakehouse/main.go
```

#### Start Lakehouse Server
```bash
# Start with lakehouse storage
./bin/ducklake-lakehouse -server -lakehouse -lakehouse-path ./my_lakehouse

# You should see:
# Initializing lakehouse storage at ./my_lakehouse
# Starting server on :8080
```

### Understanding Lakehouse Storage Structure

When you enable lakehouse mode, DuckLake creates a Delta Lake-like structure:

```
my_lakehouse/
├── _delta_log/
│   ├── metadata.json                # Table schema and metadata
│   ├── 00000000000000000000.json   # Transaction log
│   └── transactions/                # Active transactions
├── part-00000-00001.json           # Data files (immutable)
├── part-00000-00002.json
└── indexes/                        # Performance indexes
```

### Basic Lakehouse Operations

#### 1. Load Data
```bash
# Add an exercise via API
curl -X POST http://localhost:8080/api/v1/exercises \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "Morning Run",
    "type": "Cardio",
    "duration": 30,
    "calories": 300,
    "date": "2024-01-15",
    "description": "Morning jog around the park"
  }'
```

#### 2. Check Current Version
```bash
# Get current data version
curl http://localhost:8080/api/v1/lakehouse/version

# Response:
{
  "current_version": 1,
  "total_versions": 1,
  "latest_timestamp": "2024-01-15T10:30:00Z"
}
```

#### 3. Create New Version
```bash
# Add more data to create version 2
curl -X POST http://localhost:8080/api/v1/exercises \
  -H "Content-Type: application/json" \
  -d '{
    "id": 2,
    "name": "Weight Training",
    "type": "Strength",
    "duration": 45,
    "calories": 250,
    "date": "2024-01-15"
  }'

# Check version again
curl http://localhost:8080/api/v1/lakehouse/version
# Now shows version 2
```

### Time Travel Queries

#### View Historical Data
```bash
# Query version 1 (before adding weight training)
curl http://localhost:8080/api/v1/lakehouse/time-travel/1

# Query version 2 (after adding weight training)
curl http://localhost:8080/api/v1/lakehouse/time-travel/2

# Compare the results to see data evolution
```

#### Time Travel Use Cases
- **Auditing**: See what data looked like at any point
- **Recovery**: Rollback to previous versions
- **Analysis**: Compare data changes over time
- **Compliance**: Maintain historical records

### Schema Evolution

#### View Current Schema
```bash
curl http://localhost:8080/api/v1/lakehouse/schema

# Response shows current table structure
{
  "name": "exercises",
  "version": 1,
  "fields": [
    {"name": "id", "type": "integer", "nullable": false},
    {"name": "name", "type": "string", "nullable": false},
    {"name": "type", "type": "string", "nullable": false},
    {"name": "duration", "type": "integer", "nullable": false},
    {"name": "calories", "type": "integer", "nullable": false},
    {"name": "date", "type": "string", "nullable": false},
    {"name": "description", "type": "string", "nullable": true}
  ]
}
```

#### Evolve Schema (Add New Field)
```bash
# Add a new field "intensity" to track workout intensity
curl -X PUT http://localhost:8080/api/v1/lakehouse/schema \
  -H "Content-Type: application/json" \
  -d '{
    "fields": [
      {"name": "id", "type": "integer", "nullable": false},
      {"name": "name", "type": "string", "nullable": false},
      {"name": "type", "type": "string", "nullable": false},
      {"name": "duration", "type": "integer", "nullable": false},
      {"name": "calories", "type": "integer", "nullable": false},
      {"name": "date", "type": "string", "nullable": false},
      {"name": "description", "type": "string", "nullable": true},
      {"name": "intensity", "type": "string", "nullable": true}
    ]
  }'
```

#### Backward Compatibility
DuckLake ensures existing data remains accessible even after schema changes:
- New fields are optional (nullable)
- Existing queries continue to work
- Old data shows `null` for new fields

### ACID Transactions

#### Begin Transaction
```bash
# Start a new transaction
curl -X POST http://localhost:8080/api/v1/lakehouse/transactions

# Response includes transaction ID
{
  "transaction_id": "tx_1642234567890",
  "status": "active",
  "started_at": "2024-01-15T10:30:00Z"
}
```

#### Work Within Transaction
```bash
# Add multiple exercises within the transaction
# (Implementation depends on your specific lakehouse setup)
```

#### Commit Transaction
```bash
# Commit the transaction
curl -X POST http://localhost:8080/api/v1/lakehouse/transactions/tx_1642234567890/commit

# Response
{
  "transaction_id": "tx_1642234567890",
  "status": "committed",
  "committed_at": "2024-01-15T10:35:00Z"
}
```

#### Rollback Transaction
```bash
# If something goes wrong, rollback
curl -X POST http://localhost:8080/api/v1/lakehouse/transactions/tx_1642234567890/rollback

# Response
{
  "transaction_id": "tx_1642234567890",
  "status": "rolled_back",
  "rolled_back_at": "2024-01-15T10:35:00Z"
}
```

### Data Constraints

#### Add Data Quality Constraints
```bash
# Ensure duration is always positive
curl -X POST http://localhost:8080/api/v1/lakehouse/constraints \
  -H "Content-Type: application/json" \
  -d '{
    "name": "positive_duration",
    "type": "range",
    "columns": ["duration"],
    "expression": "duration > 0",
    "enabled": true
  }'

# Ensure calories are reasonable (0-2000)
curl -X POST http://localhost:8080/api/v1/lakehouse/constraints \
  -H "Content-Type: application/json" \
  -d '{
    "name": "reasonable_calories",
    "type": "range",
    "columns": ["calories"],
    "expression": "calories >= 0 AND calories <= 2000",
    "enabled": true
  }'
```

#### List Active Constraints
```bash
curl http://localhost:8080/api/v1/lakehouse/constraints

# Response shows all active constraints
[
  {
    "name": "positive_duration",
    "type": "range",
    "columns": ["duration"],
    "expression": "duration > 0",
    "enabled": true,
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "name": "reasonable_calories",
    "type": "range",
    "columns": ["calories"],
    "expression": "calories >= 0 AND calories <= 2000",
    "enabled": true,
    "created_at": "2024-01-15T10:31:00Z"
  }
]
```

#### Remove Constraint
```bash
# Remove a constraint
curl -X DELETE http://localhost:8080/api/v1/lakehouse/constraints/positive_duration
```

### Data Quality Monitoring

#### Get Data Quality Metrics
```bash
curl http://localhost:8080/api/v1/lakehouse/data-quality

# Response shows comprehensive quality metrics
{
  "total_records": 1000,
  "valid_records": 985,
  "invalid_records": 15,
  "completeness": 0.985,
  "constraints_passed": 970,
  "constraints_failed": 30,
  "quality_score": 0.97,
  "metrics": {
    "null_values": {
      "id": 0,
      "name": 2,
      "type": 0,
      "duration": 1,
      "calories": 3,
      "date": 0,
      "description": 120
    },
    "constraint_violations": [
      {
        "constraint": "positive_duration",
        "violations": 5,
        "examples": ["Exercise ID 501: duration = -10"]
      }
    ]
  }
}
```

---

## 6. Advanced Scenarios

### What You'll Learn
- Complex querying and filtering
- Performance optimization
- Indexing strategies
- Data analytics workflows

### Advanced Querying

The lakehouse provides powerful querying capabilities beyond basic CRUD operations.

#### Complex Filtering
```bash
# Find high-intensity cardio workouts over 30 minutes
curl -X POST http://localhost:8080/api/v1/lakehouse/query \
  -H "Content-Type: application/json" \
  -d '{
    "conditions": [
      {"field": "type", "operator": "equal", "value": "Cardio"},
      {"field": "duration", "operator": "greater_than", "value": 30},
      {"field": "calories", "operator": "greater_than", "value": 250}
    ],
    "sort_by": [
      {"field": "calories", "order": "desc"}
    ],
    "limit": 10
  }'
```

#### Date-Based Queries
```bash
# Find exercises from the last 7 days
curl -X POST http://localhost:8080/api/v1/lakehouse/query \
  -H "Content-Type: application/json" \
  -d '{
    "conditions": [
      {"field": "date", "operator": "greater_than", "value": "2024-01-08"},
      {"field": "date", "operator": "less_than_or_equal", "value": "2024-01-15"}
    ],
    "sort_by": [
      {"field": "date", "order": "desc"}
    ]
  }'
```

#### Aggregation Queries
```bash
# Get exercise statistics by type
curl -X POST http://localhost:8080/api/v1/lakehouse/query \
  -H "Content-Type: application/json" \
  -d '{
    "aggregation": {
      "group_by": ["type"],
      "functions": [
        {"field": "duration", "function": "sum", "alias": "total_duration"},
        {"field": "calories", "function": "avg", "alias": "avg_calories"},
        {"field": "id", "function": "count", "alias": "exercise_count"}
      ]
    }
  }'
```

### Performance Optimization

#### Creating Indexes
```bash
# Create index on exercise type for faster filtering
curl -X POST http://localhost:8080/api/v1/lakehouse/indexes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "type_index",
    "columns": ["type"],
    "type": "btree"
  }'

# Create composite index for date range queries
curl -X POST http://localhost:8080/api/v1/lakehouse/indexes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "date_type_index",
    "columns": ["date", "type"],
    "type": "btree"
  }'

# Create hash index for exact ID lookups
curl -X POST http://localhost:8080/api/v1/lakehouse/indexes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "id_hash_index",
    "columns": ["id"],
    "type": "hash"
  }'
```

#### Query Statistics
```bash
# Get query performance statistics
curl http://localhost:8080/api/v1/lakehouse/stats

# Response shows query performance metrics
{
  "total_queries": 1250,
  "avg_query_time_ms": 45,
  "slow_queries": 12,
  "index_usage": {
    "type_index": 450,
    "date_type_index": 200,
    "id_hash_index": 800
  },
  "most_common_queries": [
    {"pattern": "type = ?", "count": 300},
    {"pattern": "id = ?", "count": 500},
    {"pattern": "date BETWEEN ? AND ?", "count": 150}
  ]
}
```

#### Data Compaction
```bash
# Compact data files for better performance
curl -X POST http://localhost:8080/api/v1/lakehouse/compact

# Response shows compaction results
{
  "files_before": 50,
  "files_after": 12,
  "size_reduction": "68%",
  "compaction_time_ms": 2500,
  "status": "completed"
}
```

### Analytics Workflows

#### Weekly Exercise Summary
```bash
#!/bin/bash
# weekly_summary.sh - Generate weekly exercise analytics

# Get total exercises this week
TOTAL=$(curl -s -X POST http://localhost:8080/api/v1/lakehouse/query \
  -H "Content-Type: application/json" \
  -d '{
    "conditions": [
      {"field": "date", "operator": "greater_than", "value": "2024-01-08"}
    ],
    "aggregation": {
      "functions": [
        {"field": "id", "function": "count", "alias": "total"}
      ]
    }
  }' | jq '.results[0].total')

# Get calories burned by type
curl -s -X POST http://localhost:8080/api/v1/lakehouse/query \
  -H "Content-Type: application/json" \
  -d '{
    "conditions": [
      {"field": "date", "operator": "greater_than", "value": "2024-01-08"}
    ],
    "aggregation": {
      "group_by": ["type"],
      "functions": [
        {"field": "calories", "function": "sum", "alias": "total_calories"},
        {"field": "duration", "function": "sum", "alias": "total_minutes"}
      ]
    }
  }' | jq

echo "Total exercises this week: $TOTAL"
```

---

## 7. Production Deployment

### What You'll Learn
- Deploy with Docker Compose
- Configure for production
- Set up monitoring and logging
- Scale horizontally

### Docker Deployment

#### Production Docker Compose Setup
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  ducklake:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_USER=ducklake
      - DB_PASSWORD=secure_password
      - DB_NAME=ducklake_prod
      - DB_SSLMODE=require
    depends_on:
      - postgres
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=ducklake_prod
      - POSTGRES_USER=ducklake
      - POSTGRES_PASSWORD=secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ducklake"]
      interval: 30s
      timeout: 10s
      retries: 3

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - ducklake
    restart: unless-stopped

volumes:
  postgres_data:
```

#### Deploy to Production
```bash
# Deploy the stack
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f ducklake
```

### Environment Configuration

#### Production Environment Variables
```bash
# Database configuration
export DB_HOST=prod-postgres.example.com
export DB_PORT=5432
export DB_USER=ducklake_prod
export DB_PASSWORD=your_secure_password
export DB_NAME=ducklake_production
export DB_SSLMODE=require

# Application configuration
export PORT=8080
export LOG_LEVEL=info
export LOG_FORMAT=json

# Lakehouse configuration
export LAKEHOUSE_PATH=/data/lakehouse
export ENABLE_METRICS=true
export METRICS_PORT=9090
```

#### Configuration File (config.yaml)
```yaml
database:
  host: ${DB_HOST:localhost}
  port: ${DB_PORT:5432}
  user: ${DB_USER:ducklake}
  password: ${DB_PASSWORD:password}
  name: ${DB_NAME:ducklake_db}
  sslmode: ${DB_SSLMODE:disable}
  max_connections: 25
  connection_timeout: 30s

server:
  port: ${PORT:8080}
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s

lakehouse:
  enabled: ${LAKEHOUSE_ENABLED:false}
  path: ${LAKEHOUSE_PATH:./ducklake_data}
  compaction_interval: 1h
  max_versions: 100

logging:
  level: ${LOG_LEVEL:info}
  format: ${LOG_FORMAT:text}
  
monitoring:
  enabled: ${ENABLE_METRICS:false}
  port: ${METRICS_PORT:9090}
```

### Load Balancing and Scaling

#### Multiple Instance Deployment
```yaml
# docker-compose.scale.yml
version: '3.8'

services:
  ducklake:
    build: .
    environment:
      - DB_HOST=postgres
      - DB_USER=ducklake
      - DB_PASSWORD=secure_password
    depends_on:
      - postgres
    restart: unless-stopped
    
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx-lb.conf:/etc/nginx/nginx.conf
    depends_on:
      - ducklake
    restart: unless-stopped

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=ducklake_prod
      - POSTGRES_USER=ducklake
      - POSTGRES_PASSWORD=secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

#### NGINX Load Balancer Configuration
```nginx
# nginx-lb.conf
upstream ducklake_backend {
    server ducklake_1:8080;
    server ducklake_2:8080;
    server ducklake_3:8080;
}

server {
    listen 80;
    
    location / {
        proxy_pass http://ducklake_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        
        # Health check
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
    }
    
    location /health {
        access_log off;
        proxy_pass http://ducklake_backend;
    }
}
```

#### Scale the Application
```bash
# Scale to 3 instances
docker-compose -f docker-compose.scale.yml up -d --scale ducklake=3

# Check running instances
docker-compose -f docker-compose.scale.yml ps
```

### Monitoring and Observability

#### Health Checks
```bash
# Application health
curl http://localhost:8080/health

# Database connectivity
curl http://localhost:8080/ready

# Lakehouse status
curl http://localhost:8080/api/v1/lakehouse/health
```

#### Metrics Collection
DuckLake can expose Prometheus metrics:

```bash
# Enable metrics
export ENABLE_METRICS=true
export METRICS_PORT=9090

# Start with metrics enabled
./bin/ducklake-lakehouse -server -lakehouse -metrics

# View metrics
curl http://localhost:9090/metrics
```

#### Log Aggregation
```bash
# Structured JSON logging
export LOG_FORMAT=json
export LOG_LEVEL=info

# View logs with structured data
docker-compose logs ducklake | jq
```

---

## 8. Development Guide

### What You'll Learn
- Set up development environment
- Run tests and quality checks
- Contribute to the project
- Debug common issues

### Development Environment Setup

#### 1. Clone and Setup
```bash
git clone https://github.com/Yang92047111/ducklake-quick-start.git
cd ducklake-quick-start

# Install Go dependencies
make deps

# Install development tools
make install-tools
```

#### 2. Development Dependencies
```bash
# Install additional development tools
go install github.com/cosmtrek/air@latest          # Hot reloading
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Linting
go install github.com/go-delve/delve/cmd/dlv@latest  # Debugging
```

#### 3. IDE Setup (VS Code)
```json
// .vscode/settings.json
{
    "go.testFlags": ["-v"],
    "go.coverOnSave": true,
    "go.lintOnSave": "package",
    "go.formatTool": "gofmt",
    "go.useLanguageServer": true,
    "files.exclude": {
        "**/bin": true,
        "**/.git": true
    }
}
```

### Development Workflow

#### 1. Hot Reloading Development
```bash
# Install Air for hot reloading
go install github.com/cosmtrek/air@latest

# Create .air.toml configuration
cat > .air.toml << EOF
root = "."
testdata_dir = "test"
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/ducklake-loader"
bin = "tmp/main"
full_bin = "./tmp/main -memory -server"
include_ext = ["go", "mod"]
exclude_dir = ["tmp", "bin", "vendor"]

[log]
time = true
EOF

# Start development server with hot reloading
air
```

#### 2. Testing Workflow
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/api -v

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

#### 3. Code Quality Checks
```bash
# Format code
make format

# Run linter
make lint

# Vet code
make vet

# Security scan
make security-scan

# Run all quality checks
make quality-check
```

### Writing Tests

#### Unit Test Example
```go
// internal/loader/validator_test.go
package loader

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
)

func TestValidateExercise(t *testing.T) {
    tests := []struct {
        name        string
        exercise    Exercise
        expectError bool
        errorMsg    string
    }{
        {
            name: "valid exercise",
            exercise: Exercise{
                ID:       1,
                Name:     "Test Exercise",
                Type:     "cardio",
                Duration: 30,
                Calories: 300,
                Date:     time.Now(),
            },
            expectError: false,
        },
        {
            name: "invalid type",
            exercise: Exercise{
                ID:       1,
                Name:     "Test Exercise",
                Type:     "invalid",
                Duration: 30,
                Calories: 300,
                Date:     time.Now(),
            },
            expectError: true,
            errorMsg:    "invalid exercise type",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateExercise(tt.exercise)
            
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

#### Integration Test Example
```go
// test/integration_test.go
package test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAPIEndpoints(t *testing.T) {
    // Setup test server
    repo := storage.NewMemoryRepository()
    handler := api.NewHandler(repo)
    server := httptest.NewServer(handler)
    defer server.Close()

    // Test health endpoint
    resp, err := http.Get(server.URL + "/health")
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // Test exercises endpoint
    resp, err = http.Get(server.URL + "/exercises")
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    var exercises []Exercise
    err = json.NewDecoder(resp.Body).Decode(&exercises)
    require.NoError(t, err)
    assert.IsType(t, []Exercise{}, exercises)
}
```

### Debugging

#### 1. Using Delve Debugger
```bash
# Debug main application
dlv debug ./cmd/ducklake-loader -- -memory -server

# Debug specific test
dlv test ./internal/api -- -test.run TestSpecificFunction

# Debug with breakpoints
dlv debug ./cmd/ducklake-loader
(dlv) break main.main
(dlv) continue
```

#### 2. VS Code Debugging
```json
// .vscode/launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug DuckLake",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/ducklake-loader",
            "args": ["-memory", "-server"],
            "env": {"LOG_LEVEL": "debug"}
        },
        {
            "name": "Debug Tests",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}/internal/api",
            "args": ["-test.v"]
        }
    ]
}
```

#### 3. Logging for Debug
```go
// Enable debug logging
import "log"

func debugFunction() {
    log.Printf("Debug: Processing exercise ID %d", exerciseID)
    log.Printf("Debug: Database query: %s", query)
    log.Printf("Debug: Response size: %d bytes", len(response))
}
```

### Contributing Guidelines

#### 1. Code Style
- Follow standard Go conventions
- Use meaningful variable and function names
- Write comprehensive tests for new features
- Document exported functions and types
- Maintain backward compatibility

#### 2. Pull Request Process
```bash
# 1. Create feature branch
git checkout -b feature/my-new-feature

# 2. Make changes and test
make test
make quality-check

# 3. Commit with clear message
git commit -m "feat: add exercise filtering by intensity level"

# 4. Push and create PR
git push origin feature/my-new-feature
```

#### 3. Commit Message Format
```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

---

## 9. Troubleshooting

### Common Issues and Solutions

#### 1. Build Issues

**Problem**: `go build` fails with dependency errors
```bash
# Solution: Clean and reinstall dependencies
go clean -modcache
go mod download
go mod tidy
make deps
```

**Problem**: Binary not found after build
```bash
# Check if binary was created
ls -la bin/

# Rebuild if necessary
make clean
make build
```

#### 2. Database Issues

**Problem**: Cannot connect to PostgreSQL
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Reset database
make db-reset
```

**Problem**: Permission denied errors
```bash
# Check database user permissions
docker-compose exec postgres psql -U ducklake -d ducklake_db -c "\du"

# Recreate database with proper permissions
docker-compose down
docker volume rm ducklake-quick-start_postgres_data
docker-compose up -d
```

#### 3. API Issues

**Problem**: Server fails to start
```bash
# Check if port is already in use
lsof -i :8080

# Use different port
./bin/ducklake-loader -server -port 9000

# Check for binding issues
./bin/ducklake-loader -server -help
```

**Problem**: API returns 500 errors
```bash
# Enable debug logging
LOG_LEVEL=debug ./bin/ducklake-loader -server

# Check application logs
tail -f server.log

# Test with curl verbose mode
curl -v http://localhost:8080/exercises
```

#### 4. Test Issues

**Problem**: Tests fail with timeout
```bash
# Increase test timeout
go test -timeout 60s ./...

# Run specific failing test
go test -run TestSpecificFunction ./internal/package -v
```

**Problem**: PostgreSQL tests fail
```bash
# Skip PostgreSQL tests if no database
export SKIP_POSTGRES_TESTS=true
make test
```

#### 5. Docker Issues

**Problem**: Docker build fails
```bash
# Check Docker daemon
docker version

# Clean Docker cache
docker system prune -f

# Rebuild images
make docker-build
```

**Problem**: Services won't start
```bash
# Check service status
docker-compose ps

# View detailed logs
docker-compose logs --details

# Restart services
docker-compose restart
```

#### 6. Lakehouse Issues

**Problem**: Lakehouse initialization fails
```bash
# Check directory permissions
ls -la ./ducklake_data/

# Create directory manually
mkdir -p ./ducklake_data
chmod 755 ./ducklake_data

# Start with explicit path
./bin/ducklake-lakehouse -lakehouse-path /tmp/test_lakehouse
```

**Problem**: Version conflicts
```bash
# Check current lakehouse version
curl http://localhost:8080/api/v1/lakehouse/version

# Reset lakehouse data
rm -rf ./ducklake_data
./bin/ducklake-lakehouse -lakehouse-path ./ducklake_data
```

### Debug Commands

#### System Information
```bash
# Go version and environment
go version
go env

# System resources
df -h
free -h
ps aux | grep ducklake
```

#### Application Diagnostics
```bash
# Check binary info
file bin/ducklake-loader
ldd bin/ducklake-loader

# Network diagnostics
netstat -tulpn | grep :8080
curl -I http://localhost:8080/health
```

#### Performance Profiling
```bash
# CPU profiling
go tool pprof http://localhost:8080/debug/pprof/profile

# Memory profiling
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine profiling
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

---

## 10. Best Practices

### Data Management Best Practices

#### 1. Data Validation
```go
// Always validate input data
func ValidateExercise(e Exercise) error {
    if e.ID <= 0 {
        return errors.New("exercise ID must be positive")
    }
    if strings.TrimSpace(e.Name) == "" {
        return errors.New("exercise name cannot be empty")
    }
    if !isValidExerciseType(e.Type) {
        return errors.New("invalid exercise type")
    }
    return nil
}
```

#### 2. Error Handling
```go
// Comprehensive error handling
func ProcessExercise(e Exercise) error {
    if err := ValidateExercise(e); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    if err := repo.Insert(e); err != nil {
        return fmt.Errorf("database insert failed: %w", err)
    }
    
    log.Printf("Successfully processed exercise: %s", e.Name)
    return nil
}
```

#### 3. Consistent Data Formats
```bash
# Use ISO 8601 for dates
"date": "2024-01-15T08:00:00Z"

# Use consistent exercise types
"type": "cardio"  # not "Cardio" or "CARDIO"

# Use positive integers for IDs
"id": 1  # not 0 or negative numbers
```

### API Design Best Practices

#### 1. RESTful Endpoints
```
GET    /exercises         # List all exercises
GET    /exercises/{id}    # Get specific exercise
POST   /exercises         # Create new exercise
PUT    /exercises/{id}    # Update exercise
DELETE /exercises/{id}    # Delete exercise
```

#### 2. Error Responses
```json
{
  "error": "Bad Request",
  "code": 400,
  "message": "Exercise ID must be a positive integer",
  "details": {
    "field": "id",
    "provided": "abc",
    "expected": "positive integer"
  }
}
```

#### 3. Pagination
```bash
# Include pagination for large datasets
curl "http://localhost:8080/exercises?limit=50&offset=100"

# Response includes pagination metadata
{
  "exercises": [...],
  "pagination": {
    "limit": 50,
    "offset": 100,
    "total": 1000,
    "has_more": true
  }
}
```

### Performance Best Practices

#### 1. Database Optimization
```sql
-- Create indexes for common queries
CREATE INDEX idx_exercises_type ON exercises(type);
CREATE INDEX idx_exercises_date ON exercises(date);
CREATE INDEX idx_exercises_date_type ON exercises(date, type);
```

#### 2. Connection Management
```go
// Use connection pooling
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

#### 3. Caching Strategies
```go
// Cache frequently accessed data
type CachedRepository struct {
    repo  ExerciseRepository
    cache map[int]Exercise
    mutex sync.RWMutex
}

func (c *CachedRepository) GetByID(id int) (*Exercise, error) {
    c.mutex.RLock()
    if exercise, found := c.cache[id]; found {
        c.mutex.RUnlock()
        return &exercise, nil
    }
    c.mutex.RUnlock()
    
    // Cache miss - fetch from database
    exercise, err := c.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    c.mutex.Lock()
    c.cache[id] = *exercise
    c.mutex.Unlock()
    
    return exercise, nil
}
```

### Security Best Practices

#### 1. Input Sanitization
```go
// Sanitize all user inputs
func SanitizeExerciseName(name string) string {
    // Remove potentially dangerous characters
    name = strings.TrimSpace(name)
    name = html.EscapeString(name)
    
    // Limit length
    if len(name) > 255 {
        name = name[:255]
    }
    
    return name
}
```

#### 2. SQL Injection Prevention
```go
// Always use parameterized queries
func (r *PostgresRepository) GetByType(exerciseType string) ([]Exercise, error) {
    query := `SELECT id, name, type, duration, calories, date, description 
              FROM exercises WHERE type = $1`
    
    rows, err := r.db.Query(query, exerciseType)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    // Process results...
}
```

#### 3. Environment Configuration
```bash
# Never hardcode secrets
export DB_PASSWORD="$(cat /run/secrets/db_password)"
export API_KEY="$(cat /run/secrets/api_key)"

# Use strong passwords
export DB_PASSWORD="$(openssl rand -base64 32)"
```

### Deployment Best Practices

#### 1. Health Checks
```go
// Implement comprehensive health checks
func HealthHandler(repo ExerciseRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        health := map[string]interface{}{
            "status": "healthy",
            "timestamp": time.Now(),
            "version": "1.0.0",
        }
        
        // Check database connectivity
        if err := repo.Ping(); err != nil {
            health["status"] = "unhealthy"
            health["database"] = "disconnected"
            w.WriteHeader(http.StatusServiceUnavailable)
        }
        
        json.NewEncoder(w).Encode(health)
    }
}
```

#### 2. Graceful Shutdown
```go
// Handle shutdown signals properly
func main() {
    server := &http.Server{Addr: ":8080", Handler: handler}
    
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()
    
    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    
    log.Println("Server exited")
}
```

#### 3. Monitoring and Logging
```go
// Structured logging
import "github.com/sirupsen/logrus"

func LogRequest(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create a response writer wrapper to capture status code
        ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        handler.ServeHTTP(ww, r)
        
        logrus.WithFields(logrus.Fields{
            "method":     r.Method,
            "path":       r.URL.Path,
            "status":     ww.statusCode,
            "duration":   time.Since(start).Milliseconds(),
            "user_agent": r.UserAgent(),
            "remote_ip":  r.RemoteAddr,
        }).Info("Request processed")
    })
}
```

---

## Conclusion

🎉 **Congratulations!** You've completed the comprehensive DuckLake tutorial. You now have the knowledge to:

- ✅ Set up and run DuckLake in various configurations
- ✅ Load and manage exercise data effectively
- ✅ Use both standard and enterprise lakehouse features
- ✅ Deploy DuckLake in production environments
- ✅ Develop and contribute to the project
- ✅ Troubleshoot common issues
- ✅ Follow best practices for security and performance

### Next Steps

1. **Explore Advanced Features**: Try the enterprise lakehouse capabilities
2. **Build Your Own Dataset**: Create exercise data for your specific use case
3. **Integrate with Other Systems**: Use DuckLake's APIs in your applications
4. **Contribute**: Join the development community and help improve DuckLake
5. **Scale Up**: Deploy DuckLake in your production environment

### Additional Resources

- **[README.md](README.md)** - Project overview and quick reference
- **[DEVELOPMENT.md](DEVELOPMENT.md)** - Detailed development guide
- **[LAKEHOUSE_DEMO.md](LAKEHOUSE_DEMO.md)** - Enterprise feature showcase
- **[GitHub Issues](https://github.com/Yang92047111/ducklake-quick-start/issues)** - Report bugs or request features

### Getting Help

If you encounter issues or have questions:

1. **Check the troubleshooting section** in this tutorial
2. **Review the documentation** links above
3. **Search existing issues** on GitHub
4. **Create a new issue** with detailed information
5. **Join the community** discussions

**Happy data management with DuckLake!** 🦆

---

*This tutorial is continuously updated. For the latest version, visit the [GitHub repository](https://github.com/Yang92047111/ducklake-quick-start).*
