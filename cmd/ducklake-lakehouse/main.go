package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Yang92047111/ducklake-quick-start/internal/api"
	"github.com/Yang92047111/ducklake-quick-start/internal/loader"
	"github.com/Yang92047111/ducklake-quick-start/internal/storage"
)

func main() {
	var (
		csvFile       = flag.String("csv", "", "Path to CSV file to load")
		jsonFile      = flag.String("json", "", "Path to JSON file to load")
		useMemory     = flag.Bool("memory", false, "Use in-memory storage instead of PostgreSQL or lakehouse")
		useLakehouse  = flag.Bool("lakehouse", false, "Use lakehouse (Delta Lake) storage")
		lakehousePath = flag.String("lakehouse-path", "./ducklake_data", "Path for lakehouse data storage")
		serverMode    = flag.Bool("server", false, "Run in server mode")
		port          = flag.String("port", "8080", "Server port")
	)
	flag.Parse()

	log.Println("DuckLake Loader starting...")

	// Initialize repository based on storage type
	var repo storage.ExerciseRepository
	var lakehouseRepo storage.LakehouseRepository
	var err error

	if *useMemory {
		log.Println("Using in-memory storage")
		repo = storage.NewMemoryRepository()
	} else if *useLakehouse {
		log.Printf("Using lakehouse storage at: %s", *lakehousePath)

		// Ensure lakehouse path exists
		if err := os.MkdirAll(*lakehousePath, 0755); err != nil {
			log.Fatalf("Failed to create lakehouse directory: %v", err)
		}

		// Create lakehouse configuration
		config := &storage.DeltaConfig{
			MaxFileSize:        100 * 1024 * 1024, // 100MB
			MinFileSize:        1 * 1024 * 1024,   // 1MB
			CheckpointInterval: 10,
			EnableOptimization: true,
			AutoCompact:        true,
			CompressionCodec:   "snappy",
		}

		lakehouseRepo, err = storage.NewDeltaLakeRepository(*lakehousePath, config)
		if err != nil {
			log.Fatalf("Failed to create lakehouse repository: %v", err)
		}
		repo = lakehouseRepo
		log.Println("Lakehouse initialized successfully")
	} else {
		// Build PostgreSQL connection string from environment variables
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "5432")
		dbUser := getEnv("DB_USER", "ducklake")
		dbPassword := getEnv("DB_PASSWORD", "password")
		dbName := getEnv("DB_NAME", "ducklake_db")
		dbSSLMode := getEnv("DB_SSLMODE", "disable")

		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

		log.Println("Connecting to PostgreSQL...")
		repo, err = storage.NewPostgresRepository(connStr)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}
	defer repo.Close()

	// Load data if files are specified
	if *csvFile != "" {
		if err := loadCSVData(repo, *csvFile); err != nil {
			log.Fatalf("Failed to load CSV data: %v", err)
		}
	}

	if *jsonFile != "" {
		if err := loadJSONData(repo, *jsonFile); err != nil {
			log.Fatalf("Failed to load JSON data: %v", err)
		}
	}

	// Start server if requested
	if *serverMode {
		log.Printf("Starting server on port %s", *port)

		var router http.Handler

		if lakehouseRepo != nil {
			// Use lakehouse handler with advanced features
			log.Println("Starting server with lakehouse features enabled")
			handler := api.NewLakehouseHandler(lakehouseRepo)
			router = handler.SetupLakehouseRoutes()

			// Display available lakehouse endpoints
			displayLakehouseEndpoints()
		} else {
			// Use standard handler
			log.Println("Starting server with standard features")
			handler := api.NewHandler(repo)
			router = handler.SetupRoutes()
		}

		log.Fatal(http.ListenAndServe(":"+*port, router))
	}

	log.Println("DuckLake Loader completed successfully")
}

func loadCSVData(repo storage.ExerciseRepository, filename string) error {
	log.Printf("Loading CSV data from %s", filename)

	csvLoader := loader.NewCSVLoader()
	exercises, err := csvLoader.LoadFromCSV(filename)
	if err != nil {
		return err
	}

	validator := loader.NewValidator()
	for i, exercise := range exercises {
		if err := validator.Validate(exercise); err != nil {
			log.Printf("Validation error for exercise %d: %v", i+1, err)
			continue
		}
	}

	if err := repo.InsertBatch(exercises); err != nil {
		return err
	}

	log.Printf("Successfully loaded %d exercises from CSV", len(exercises))
	return nil
}

func loadJSONData(repo storage.ExerciseRepository, filename string) error {
	log.Printf("Loading JSON data from %s", filename)

	jsonLoader := loader.NewJSONLoader()
	exercises, err := jsonLoader.LoadFromJSON(filename)
	if err != nil {
		return err
	}

	validator := loader.NewValidator()
	for i, exercise := range exercises {
		if err := validator.Validate(exercise); err != nil {
			log.Printf("Validation error for exercise %d: %v", i+1, err)
			continue
		}
	}

	if err := repo.InsertBatch(exercises); err != nil {
		return err
	}

	log.Printf("Successfully loaded %d exercises from JSON", len(exercises))
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func displayLakehouseEndpoints() {
	log.Println("\n=== DuckLake Lakehouse API Endpoints ===")
	log.Println("Standard Exercise Endpoints:")
	log.Println("  GET    /exercises")
	log.Println("  GET    /exercises/{id}")
	log.Println("  GET    /exercises/type/{type}")
	log.Println("  GET    /exercises/date-range?start=YYYY-MM-DD&end=YYYY-MM-DD")
	log.Println()
	log.Println("Lakehouse Version & Time Travel:")
	log.Println("  GET    /api/v1/versions                    - Get version history")
	log.Println("  GET    /api/v1/versions/{version}          - Get data at specific version")
	log.Println("  POST   /api/v1/versions                    - Create new version")
	log.Println("  GET    /api/v1/time-travel?timestamp=...   - Time travel queries")
	log.Println()
	log.Println("Schema Management:")
	log.Println("  GET    /api/v1/schema                      - Get current schema")
	log.Println("  PUT    /api/v1/schema                      - Evolve schema")
	log.Println("  GET    /api/v1/schema/history              - Get schema history")
	log.Println("  POST   /api/v1/schema/validate             - Validate schema compatibility")
	log.Println()
	log.Println("Metadata & Catalog:")
	log.Println("  GET    /api/v1/metadata                    - Get table metadata")
	log.Println("  PUT    /api/v1/metadata/properties         - Update table properties")
	log.Println("  GET    /api/v1/partitions                  - Get partition information")
	log.Println()
	log.Println("Transaction Management:")
	log.Println("  POST   /api/v1/transactions                - Begin transaction")
	log.Println("  POST   /api/v1/transactions/{id}/commit    - Commit transaction")
	log.Println("  POST   /api/v1/transactions/{id}/rollback  - Rollback transaction")
	log.Println("  GET    /api/v1/transactions/{id}/status    - Get transaction status")
	log.Println()
	log.Println("Optimization & Maintenance:")
	log.Println("  POST   /api/v1/optimize                    - Optimize table")
	log.Println("  POST   /api/v1/compact                     - Compact table")
	log.Println("  POST   /api/v1/vacuum                      - Vacuum old files")
	log.Println()
	log.Println("Data Quality & Constraints:")
	log.Println("  GET    /api/v1/constraints                 - Get constraints")
	log.Println("  POST   /api/v1/constraints                 - Add constraint")
	log.Println("  DELETE /api/v1/constraints/{name}          - Remove constraint")
	log.Println("  GET    /api/v1/data-quality                - Get data quality metrics")
	log.Println()
	log.Println("Change Data Capture:")
	log.Println("  GET    /api/v1/changes                     - Get changelog")
	log.Println("  GET    /api/v1/changes/stream              - Stream changes")
	log.Println()
	log.Println("Advanced Querying:")
	log.Println("  POST   /api/v1/query/sql                   - SQL queries")
	log.Println("  POST   /api/v1/query/filter                - Advanced filtering")
	log.Println("  POST   /api/v1/query/aggregate             - Time window aggregations")
	log.Println()
	log.Println("Performance & Statistics:")
	log.Println("  GET    /api/v1/stats/query                 - Query statistics")
	log.Println("  GET    /api/v1/indexes                     - Get indexes")
	log.Println("  POST   /api/v1/indexes                     - Create index")
	log.Println("  DELETE /api/v1/indexes/{name}              - Drop index")
	log.Println("\n=========================================")
	log.Println("Example Time Travel Query:")
	log.Println("  curl \"http://localhost:8080/api/v1/time-travel?timestamp=2024-01-15T10:00:00Z\"")
	log.Println()
	log.Println("Example Version Query:")
	log.Println("  curl \"http://localhost:8080/api/v1/versions/1\"")
	log.Println()
	log.Println("Example Schema Evolution:")
	log.Println("  curl -X PUT -d '{\"fields\":[...]}' http://localhost:8080/api/v1/schema")
	log.Println("=========================================")
}
