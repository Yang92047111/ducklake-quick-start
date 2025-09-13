# DuckLake 

A modern, cloud-native Go application for enterprise-grade exercise data management. DuckLake combines the simplicity of traditional data loading with advanced **Lakehouse** capabilities for comprehensive data governance, versioning, and analytics.

## ✨ Features

### 🏆 Core Features
- **Multi-Format Data Loading**: CSV and JSON exercise data parsing
- **Data Validation**: Comprehensive validation engine for data quality
- **Flexible Storage**: PostgreSQL, in-memory, or **Lakehouse** storage options
- **REST API**: Complete API with standard and enterprise endpoints
- **Docker Support**: Production-ready containerized deployment
- **Comprehensive Testing**: Unit tests, integration tests, and coverage reporting

### 🚀 Enterprise Lakehouse Features
- **🕐 Data Versioning**: Track all changes with automatic versioning
- **⏰ Time Travel**: Query historical versions of your data
- **🔄 Schema Evolution**: Add/modify columns with backward compatibility
- **🔒 ACID Transactions**: Full transaction support with rollback capabilities
- **📋 Data Constraints**: NOT NULL, range, and custom validation rules
- **🔍 Advanced Querying**: Complex filtering, sorting, and aggregation
- **📊 Data Quality Metrics**: Monitor and validate data quality
- **🚀 Performance Optimization**: Indexing, compaction, and query optimization
- **📈 Change Tracking**: Real-time change streams and audit logs
- **🗂️ Metadata Management**: Rich metadata with custom properties

## 🏗️ Architecture Overview

DuckLake follows a clean, modular architecture that scales from simple data loading to enterprise data management:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Data Sources  │    │   Processing    │    │    Storage      │
│                 │    │                 │    │                 │
│ • CSV Files     │───▶│ • Data Loading  │───▶│ • PostgreSQL    │
│ • JSON Files    │    │ • Validation    │    │ • In-Memory     │
│ • API Imports   │    │ • Transformation│    │ • Lakehouse     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   REST API      │
                       │                 │
                       │ • Standard APIs │
                       │ • Lakehouse APIs│
                       │ • Query Engine  │
                       └─────────────────┘
```

## 📁 Project Structure

```
ducklake-quick-start/
├── cmd/
│   ├── ducklake-loader/        # Original application entry point
│   └── ducklake-lakehouse/     # 🆕 Lakehouse-enabled entry point
├── internal/
│   ├── loader/                 # Data loading and validation
│   ├── storage/                # Database repositories + Lakehouse
│   │   ├── repository.go       # Base repository interface  
│   │   ├── lakehouse.go        # 🆕 Lakehouse repository interface
│   │   ├── deltalake.go        # 🆕 Delta Lake implementation
│   │   └── lakehouse_features.go # 🆕 Advanced lakehouse features
│   └── api/                    # REST API handlers + Lakehouse APIs
│       ├── handlers.go         # Standard API endpoints
│       └── lakehouse_handlers.go # 🆕 Lakehouse API endpoints
├── test/                       # Test files and sample data
├── TUTORIAL.md                 # 🆕 Complete tutorial guide
├── test_lakehouse.sh          # 🆕 Lakehouse integration test
├── Dockerfile                 # Container build configuration
├── docker-compose.yml         # Multi-service deployment
└── Makefile                   # Build automation
```

## 🚀 Quick Start

### Prerequisites
- **Go 1.21+**: [Download here](https://golang.org/dl/)
- **Docker & Docker Compose**: For containerized deployment (optional)
- **PostgreSQL**: For database storage (optional, can use in-memory or Docker)

### 1. Standard DuckLake (Traditional)

Perfect for getting started or simple data processing needs:

```bash
# Clone and setup
git clone https://github.com/Yang92047111/ducklake-quick-start.git
cd ducklake-quick-start

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

### 2. 🆕 Enterprise Lakehouse (Advanced)

For production environments requiring data governance and advanced analytics:

```bash
# Build lakehouse binary
go build -o bin/ducklake-lakehouse cmd/ducklake-lakehouse/main.go

# Start lakehouse server
./bin/ducklake-lakehouse -server -lakehouse -lakehouse-path ./ducklake_data

# Load sample data via API
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

# Test lakehouse features
./test_lakehouse.sh
```

### 3. Docker Deployment

```bash
# Start complete stack with PostgreSQL
make docker-run

# Access API
curl http://localhost:8080/exercises

# Stop services
make docker-down
```

## 📊 Data Model

### Exercise Entity

```go
type Exercise struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`        // cardio, strength, flexibility
    Duration    int       `json:"duration"`    // minutes
    Calories    int       `json:"calories"`
    Date        time.Time `json:"date"`
    Description string    `json:"description"`
}
```

### Supported Data Formats

**CSV Format:**
```csv
id,name,type,duration,calories,date,description
1,Morning Run,cardio,30,300,2024-01-15,Easy morning jog
```

**JSON Format:**
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

## 🔗 API Endpoints

### Standard API (Backward Compatible)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/exercises` | Get all exercises |
| GET | `/exercises/{id}` | Get exercise by ID |
| GET | `/exercises/type/{type}` | Get exercises by type |
| GET | `/exercises/date-range?start=YYYY-MM-DD&end=YYYY-MM-DD` | Get exercises by date range |
| GET | `/health` | Health check |

### 🆕 Lakehouse API (Enterprise Features)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/lakehouse/version` | Get current data version |
| POST | `/api/v1/lakehouse/version` | Create new version |
| GET | `/api/v1/lakehouse/time-travel/{version}` | Query specific version |
| GET | `/api/v1/lakehouse/schema` | Get current schema |
| PUT | `/api/v1/lakehouse/schema` | Evolve schema |
| POST | `/api/v1/lakehouse/transactions` | Begin transaction |
| POST | `/api/v1/lakehouse/transactions/{id}/commit` | Commit transaction |
| POST | `/api/v1/lakehouse/transactions/{id}/rollback` | Rollback transaction |
| GET | `/api/v1/lakehouse/constraints` | List constraints |
| POST | `/api/v1/lakehouse/constraints` | Add constraint |
| GET | `/api/v1/lakehouse/data-quality` | Get data quality metrics |
| POST | `/api/v1/lakehouse/query` | Advanced query with filtering |
| POST | `/api/v1/lakehouse/indexes` | Create index |
| POST | `/api/v1/lakehouse/compact` | Compact data files |

## 💡 Usage Examples

### Standard API Usage
```bash
# Get all exercises
curl http://localhost:8080/exercises

# Get cardio exercises only
curl http://localhost:8080/exercises/type/cardio

# Get exercises in date range
curl "http://localhost:8080/exercises/date-range?start=2024-01-15&end=2024-01-16"
```

### Lakehouse API Usage
```bash
# Advanced query with filtering
curl -X POST http://localhost:8080/api/v1/lakehouse/query \
  -H "Content-Type: application/json" \
  -d '{
    "conditions": [
      {"field": "type", "operator": "equal", "value": "Cardio"},
      {"field": "duration", "operator": "greater_than", "value": 20}
    ],
    "sort_by": [{"field": "calories", "order": "desc"}],
    "limit": 10
  }'

# Add data constraint
curl -X POST http://localhost:8080/api/v1/lakehouse/constraints \
  -H "Content-Type: application/json" \
  -d '{
    "name": "positive_duration",
    "type": "range",
    "columns": ["duration"],
    "expression": "duration > 0",
    "enabled": true
  }'

# Time travel to previous version
curl http://localhost:8080/api/v1/lakehouse/time-travel/1
```

## ⚙️ Configuration

### Command Line Options

**Standard DuckLake:**
- `-csv <file>` - Load data from CSV file
- `-json <file>` - Load data from JSON file  
- `-memory` - Use in-memory storage (default: PostgreSQL)
- `-server` - Start REST API server
- `-port <port>` - Server port (default: 8080)

**Lakehouse DuckLake:**
- `-lakehouse` - Enable lakehouse (Delta Lake) storage
- `-lakehouse-path <path>` - Path for lakehouse data storage
- `-server` - Start REST API server with lakehouse endpoints
- `-port <port>` - Server port (default: 8080)

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_USER` | Database user | ducklake |
| `DB_PASSWORD` | Database password | password |
| `DB_NAME` | Database name | ducklake_db |
| `DB_SSLMODE` | SSL mode | disable |

## 🧪 Development & Testing

### Development Setup
```bash
# Clone repository
git clone https://github.com/Yang92047111/ducklake-quick-start.git
cd ducklake-quick-start

# Install dependencies
make deps

# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration
```

### Available Make Targets
```bash
make help           # Show all available commands
make build          # Build the application
make test           # Run all tests
make demo           # Quick demo with sample data
make docker-run     # Start with Docker Compose
make clean          # Clean build artifacts
```

## 🏆 Lakehouse Features Deep Dive

### Data Storage Structure (Delta Lake-like)
```
ducklake_data/
├── _delta_log/
│   ├── metadata.json                 # Table schema and metadata
│   ├── 00000000000000000000.json    # Transaction log entries
│   └── transactions/                 # Transaction records
├── part-00000-00001.json            # Data files (versioned)
├── part-00000-00002.json
└── indexes/                         # Performance indexes
```

### Key Capabilities

**🕐 Version Management**
- Every operation creates a new data version
- Complete audit trail of all changes
- Point-in-time recovery capabilities

**⏰ Time Travel Queries**
- Query any historical version of data
- Compare changes between versions
- Rollback to previous states

**🔄 Schema Evolution**
- Add new columns without breaking existing data
- Backward compatibility validation
- Type compatibility checking

**🔒 ACID Transactions**
- Begin/commit/rollback transaction support
- Isolation between concurrent operations
- Consistency guarantees

**📋 Data Quality**
- Built-in data validation rules
- Quality metrics and monitoring
- Constraint enforcement

## 🌟 Use Cases

### 1. Fitness & Health Applications
- **Personal Fitness Apps**: Workout tracking with history
- **Gym Management**: Member activity and equipment usage
- **Corporate Wellness**: Employee fitness programs
- **Healthcare**: Patient exercise prescription tracking

### 2. Data Analytics & Research
- **Sports Science**: Exercise pattern analysis
- **Performance Studies**: Training effectiveness research
- **Population Health**: Large-scale fitness data analysis
- **Machine Learning**: Training data for ML models

### 3. Enterprise Integration
- **ETL Pipelines**: Exercise data transformation workflows
- **API Integration**: Microservices data exchange
- **Data Warehousing**: Exercise data consolidation
- **Real-time Analytics**: Live fitness data processing

## 📈 Migration Path

DuckLake provides a smooth upgrade journey:

1. **🚀 Start Simple**: Use basic DuckLake with PostgreSQL or in-memory storage
2. **🔗 Add APIs**: Enable REST endpoints for data access
3. **⬆️ Upgrade Storage**: Switch to lakehouse mode (`-lakehouse`)
4. **🎯 Leverage Features**: Use versioning, constraints, advanced querying
5. **🏢 Scale Enterprise**: Full data governance and optimization

**No breaking changes when upgrading!** Both storage backends can coexist.

## 🔒 Production Readiness

### Security Features
- Input validation and sanitization
- SQL injection prevention
- Type safety throughout
- Resource limits and DoS protection

### Reliability Features
- Comprehensive error handling
- Graceful degradation
- Health check endpoints
- Connection pooling ready

### Performance Features
- Efficient memory usage
- Concurrent operation support
- Database query optimization
- Horizontal scaling ready

## Documentation

- **[TUTORIAL.md](TUTORIAL.md)** - Complete step-by-step tutorial
- **[DEVELOPMENT.md](DEVELOPMENT.md)** - Development and contribution guide
- **[LAKEHOUSE_DEMO.md](LAKEHOUSE_DEMO.md)** - Lakehouse features demonstration
- **[test_lakehouse.sh](test_lakehouse.sh)** - Automated integration testing

## 🤝 Contributing

We welcome contributions! Please see our [development guide](DEVELOPMENT.md) for:
- Setting up the development environment
- Running tests and quality checks
- Submitting pull requests
- Code style guidelines

## License

MIT License - see [LICENSE](LICENSE) file for details.

## 🎯 What's Next?

### Upcoming Features
- **Cloud Storage**: S3, GCS, Azure Blob support
- **Advanced Analytics**: Machine learning integration
- **Real-time Streaming**: Change data capture (CDC)
- **Multi-tenancy**: Tenant isolation and resource management
- **Enhanced Security**: RBAC, encryption, audit logging

### Production Enhancements
- **Monitoring**: Prometheus/Grafana integration
- **Logging**: Structured logging with correlation IDs
- **Deployment**: Kubernetes manifests and Helm charts
- **Scaling**: Auto-scaling and load balancing

---

## 🎉 Get Started Today!

Ready to transform your exercise data management? 

1. **Quick Start**: Follow the [Quick Start](#-quick-start) guide above
2. **Learn More**: Read the complete [TUTORIAL.md](TUTORIAL.md)
3. **Explore Features**: Try the [lakehouse demo](LAKEHOUSE_DEMO.md)
4. **Join Development**: Check the [development guide](DEVELOPMENT.md)

**DuckLake: From simple data loading to enterprise data lakehouse in minutes!** 🚀

---

*For questions, support, or contributions, visit the [GitHub repository](https://github.com/Yang92047111/ducklake-quick-start) or open an issue.*
