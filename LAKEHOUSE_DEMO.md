# DuckLake Lakehouse Demo

This document demonstrates the lakehouse features implemented in DuckLake.

## Quick Start

### 1. Build the Lakehouse Binary

```bash
go build -o bin/ducklake-lakehouse cmd/ducklake-lakehouse/main.go
```

### 2. Start the Lakehouse Server

```bash
./bin/ducklake-lakehouse -server -lakehouse -lakehouse-path ./demo_lakehouse
```

### 3. Load Sample Data

```bash
# Load sample data
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

## Lakehouse Features

### 1. Versioning and Time Travel

```bash
# Get current version
curl http://localhost:8080/api/v1/lakehouse/version

# Create a new version after data changes
curl -X POST http://localhost:8080/api/v1/lakehouse/version

# Time travel to specific version
curl http://localhost:8080/api/v1/lakehouse/time-travel/1
```

### 2. Schema Evolution

```bash
# Get current schema
curl http://localhost:8080/api/v1/lakehouse/schema

# Evolve schema (add new field)
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

### 3. ACID Transactions

```bash
# Begin transaction
curl -X POST http://localhost:8080/api/v1/lakehouse/transactions

# Commit transaction
curl -X POST http://localhost:8080/api/v1/lakehouse/transactions/{transaction_id}/commit

# Rollback transaction
curl -X POST http://localhost:8080/api/v1/lakehouse/transactions/{transaction_id}/rollback
```

### 4. Data Constraints

```bash
# Add constraint
curl -X POST http://localhost:8080/api/v1/lakehouse/constraints \
  -H "Content-Type: application/json" \
  -d '{
    "name": "positive_duration",
    "type": "range",
    "columns": ["duration"],
    "expression": "duration > 0",
    "enabled": true
  }'

# List constraints
curl http://localhost:8080/api/v1/lakehouse/constraints

# Remove constraint
curl -X DELETE http://localhost:8080/api/v1/lakehouse/constraints/positive_duration
```

### 5. Metadata Management

```bash
# Get table metadata
curl http://localhost:8080/api/v1/lakehouse/metadata

# Update table properties
curl -X PUT http://localhost:8080/api/v1/lakehouse/metadata/properties \
  -H "Content-Type: application/json" \
  -d '{
    "owner": "data-team",
    "description": "Exercise tracking data",
    "retention_days": "365"
  }'
```

### 6. Data Quality Metrics

```bash
# Get data quality metrics
curl http://localhost:8080/api/v1/lakehouse/data-quality
```

### 7. Query with Advanced Filtering

```bash
# Query with filtering
curl -X POST http://localhost:8080/api/v1/lakehouse/query \
  -H "Content-Type: application/json" \
  -d '{
    "conditions": [
      {"field": "type", "operator": "equal", "value": "Cardio"},
      {"field": "duration", "operator": "greater_than", "value": 20}
    ],
    "sort_by": [
      {"field": "calories", "order": "desc"}
    ],
    "limit": 10,
    "offset": 0
  }'
```

### 8. Performance Optimization

```bash
# Create index
curl -X POST http://localhost:8080/api/v1/lakehouse/indexes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "type_idx",
    "columns": ["type"],
    "type": "btree"
  }'

# Get query statistics
curl http://localhost:8080/api/v1/lakehouse/stats

# Compact data files
curl -X POST http://localhost:8080/api/v1/lakehouse/compact
```

## Data Storage Structure

The lakehouse implementation stores data in a Delta Lake-like structure:

```
./demo_lakehouse/
├── _delta_log/
│   ├── metadata.json          # Table metadata and schema
│   ├── 00000000000000000000.json  # Transaction log
│   └── transactions/          # Transaction records
├── part-00000-00001.json     # Data files (versioned)
└── part-00000-00002.json
```

## Features Implemented

✅ **Data Versioning**: Every operation creates a new version  
✅ **Time Travel**: Query historical versions of data  
✅ **Schema Evolution**: Add/modify columns with compatibility checks  
✅ **ACID Transactions**: Full transaction support with rollback  
✅ **Data Constraints**: NOT NULL, range, and custom constraints  
✅ **Metadata Management**: Rich metadata with properties  
✅ **Data Quality**: Metrics and validation  
✅ **Advanced Querying**: Filtering, sorting, aggregation  
✅ **Performance**: Indexing, compaction, query optimization  
✅ **Change Tracking**: Real-time change streams  

## Configuration

The lakehouse can be configured via command-line flags:

- `-lakehouse`: Enable lakehouse mode
- `-lakehouse-path`: Set data storage path (default: `./ducklake_data`)
- `-server`: Run in server mode
- `-port`: Set server port (default: 8080)

## Integration with Existing DuckLake

The lakehouse implementation extends the existing DuckLake project without breaking changes:

1. **Backward Compatible**: Existing API endpoints continue to work
2. **Storage Abstraction**: Uses the same repository pattern
3. **Modular Design**: Lakehouse features are optional
4. **Unified API**: Standard and lakehouse endpoints coexist

This provides a smooth migration path from simple storage to enterprise-grade lakehouse capabilities.
