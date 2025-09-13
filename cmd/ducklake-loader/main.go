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
		csvFile      = flag.String("csv", "", "Path to CSV file to load")
		jsonFile     = flag.String("json", "", "Path to JSON file to load")
		useMemory    = flag.Bool("memory", false, "Use in-memory storage instead of PostgreSQL")
		useLakehouse = flag.Bool("lakehouse", false, "Use lakehouse (Delta Lake) storage")
		serverMode   = flag.Bool("server", false, "Run in server mode")
		port         = flag.String("port", "8080", "Server port")
	)
	flag.Parse()

	log.Println("DuckLake Loader starting...")

	// Initialize repository
	var repo storage.ExerciseRepository
	var err error

	if *useLakehouse {
		log.Println("Using lakehouse (Delta Lake) storage")
		lakehousePath := getEnv("LAKEHOUSE_PATH", "./lakehouse_data")
		repo, err = storage.NewDeltaLakeRepository(lakehousePath, nil)
		if err != nil {
			log.Fatalf("Failed to initialize lakehouse: %v", err)
		}
	} else if *useMemory {
		log.Println("Using in-memory storage")
		repo = storage.NewMemoryRepository()
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
		handler := api.NewHandler(repo)
		router := handler.SetupRoutes()

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
