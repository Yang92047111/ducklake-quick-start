package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Yang92047111/ducklake-quick-start/internal/storage"
)

func main() {
	log.Println("Testing DeltaLakeRepository initialization...")

	// Create test directory
	testPath := "./test_lakehouse"
	os.RemoveAll(testPath) // Clean up any existing test data

	// Initialize repository
	repo, err := storage.NewDeltaLakeRepository(testPath, nil)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	fmt.Println("✓ Repository initialized successfully")

	// Test basic functionality
	exercises, err := repo.GetAll()
	if err != nil {
		log.Fatalf("Failed to get all exercises: %v", err)
	}

	fmt.Printf("✓ Retrieved %d exercises\n", len(exercises))

	// Close repository
	err = repo.Close()
	if err != nil {
		log.Fatalf("Failed to close repository: %v", err)
	}

	fmt.Println("✓ Repository closed successfully")

	// Clean up
	os.RemoveAll(testPath)
	fmt.Println("✓ Test completed successfully")
}
