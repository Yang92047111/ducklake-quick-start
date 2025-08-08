package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourname/ducklake-loader/internal/loader"
	"github.com/yourname/ducklake-loader/internal/storage"
)

func TestIntegration_LoadAndStore(t *testing.T) {
	// Use memory repository for integration test
	repo := storage.NewMemoryRepository()
	defer repo.Close()

	// Test CSV loading
	csvLoader := loader.NewCSVLoader()
	exercises, err := csvLoader.LoadFromCSV("testdata/sample_exercises.csv")
	require.NoError(t, err)
	require.Len(t, exercises, 5)

	// Validate exercises
	validator := loader.NewValidator()
	for _, exercise := range exercises {
		err := validator.Validate(exercise)
		require.NoError(t, err)
	}

	// Store exercises
	err = repo.InsertBatch(exercises)
	require.NoError(t, err)

	// Verify storage
	allExercises, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, allExercises, 5)

	// Test queries
	cardioExercises, err := repo.GetByType("cardio")
	require.NoError(t, err)
	assert.Len(t, cardioExercises, 2) // Morning Run and Swimming

	// Test date range query
	start, _ := time.Parse("2006-01-02", "2024-01-15")
	end, _ := time.Parse("2006-01-02", "2024-01-16")
	dateRangeExercises, err := repo.GetByDateRange(start, end)
	require.NoError(t, err)
	assert.Len(t, dateRangeExercises, 3) // Exercises on 2024-01-15 and 2024-01-16
}
