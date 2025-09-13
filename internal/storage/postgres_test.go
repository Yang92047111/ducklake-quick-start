package storage

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Yang92047111/ducklake-quick-start/internal/loader"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *PostgresRepository {
	// Skip these tests if running in CI or if no database is available
	if os.Getenv("SKIP_POSTGRES_TESTS") == "true" {
		t.Skip("Skipping PostgreSQL tests")
	}

	// Use test database configuration
	dbHost := getTestEnv("DB_HOST", "localhost")
	dbPort := getTestEnv("DB_PORT", "5432")
	dbUser := getTestEnv("DB_USER", "test")
	dbPassword := getTestEnv("DB_PASSWORD", "test")
	dbName := getTestEnv("DB_NAME", "test_db")
	dbSSLMode := getTestEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	repo, err := NewPostgresRepository(connStr)
	if err != nil {
		t.Skipf("Could not connect to test database: %v", err)
	}

	// Clean up any existing test data
	_, err = repo.db.Exec("DELETE FROM exercises")
	require.NoError(t, err)

	return repo
}

func getTestEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func TestNewPostgresRepository(t *testing.T) {
	t.Run("connects successfully with valid connection string", func(t *testing.T) {
		repo := setupTestDB(t)
		defer repo.Close()

		assert.NotNil(t, repo)
		assert.NotNil(t, repo.db)
	})

	t.Run("returns error with invalid connection string", func(t *testing.T) {
		if os.Getenv("SKIP_POSTGRES_TESTS") == "true" {
			t.Skip("Skipping PostgreSQL tests")
		}

		_, err := NewPostgresRepository("invalid connection string")
		assert.Error(t, err)
	})

	t.Run("returns error when database is unreachable", func(t *testing.T) {
		if os.Getenv("SKIP_POSTGRES_TESTS") == "true" {
			t.Skip("Skipping PostgreSQL tests")
		}

		connStr := "host=nonexistent port=5432 user=test password=test dbname=test sslmode=disable"
		_, err := NewPostgresRepository(connStr)
		assert.Error(t, err)
	})
}

func TestPostgresRepository_Insert(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	exercise := loader.Exercise{
		Name:        "Test Exercise",
		Type:        "cardio",
		Duration:    30,
		Calories:    250,
		Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Description: "Test description",
	}

	err := repo.Insert(exercise)
	assert.NoError(t, err)

	// Verify the exercise was inserted
	exercises, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, exercises, 1)
	assert.Equal(t, exercise.Name, exercises[0].Name)
	assert.Equal(t, exercise.Type, exercises[0].Type)
}

func TestPostgresRepository_InsertBatch(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	exercises := []loader.Exercise{
		{
			Name:        "Morning Run",
			Type:        "cardio",
			Duration:    30,
			Calories:    300,
			Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Description: "Easy morning jog",
		},
		{
			Name:        "Push-ups",
			Type:        "strength",
			Duration:    15,
			Calories:    150,
			Date:        time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			Description: "Upper body workout",
		},
	}

	err := repo.InsertBatch(exercises)
	assert.NoError(t, err)

	// Verify all exercises were inserted
	allExercises, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, allExercises, 2)
}

func TestPostgresRepository_GetByID(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Insert a test exercise first
	exercise := loader.Exercise{
		Name:        "Test Exercise",
		Type:        "cardio",
		Duration:    30,
		Calories:    250,
		Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Description: "Test description",
	}
	err := repo.Insert(exercise)
	require.NoError(t, err)

	// Get all exercises to find the inserted ID
	exercises, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, exercises, 1)
	insertedID := exercises[0].ID

	t.Run("returns exercise when found", func(t *testing.T) {
		found, err := repo.GetByID(insertedID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, exercise.Name, found.Name)
		assert.Equal(t, exercise.Type, found.Type)
	})

	t.Run("returns nil when not found", func(t *testing.T) {
		found, err := repo.GetByID(99999)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestPostgresRepository_GetByType(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	exercises := []loader.Exercise{
		{
			Name:     "Morning Run",
			Type:     "cardio",
			Duration: 30,
			Calories: 300,
			Date:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:     "Push-ups",
			Type:     "strength",
			Duration: 15,
			Calories: 150,
			Date:     time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:     "Evening Run",
			Type:     "cardio",
			Duration: 45,
			Calories: 450,
			Date:     time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC),
		},
	}

	err := repo.InsertBatch(exercises)
	require.NoError(t, err)

	cardioExercises, err := repo.GetByType("cardio")
	assert.NoError(t, err)
	assert.Len(t, cardioExercises, 2)

	strengthExercises, err := repo.GetByType("strength")
	assert.NoError(t, err)
	assert.Len(t, strengthExercises, 1)

	nonExistentType, err := repo.GetByType("nonexistent")
	assert.NoError(t, err)
	assert.Len(t, nonExistentType, 0)
}

func TestPostgresRepository_GetByDateRange(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	exercises := []loader.Exercise{
		{
			Name:     "Exercise 1",
			Type:     "cardio",
			Duration: 30,
			Calories: 300,
			Date:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:     "Exercise 2",
			Type:     "strength",
			Duration: 15,
			Calories: 150,
			Date:     time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:     "Exercise 3",
			Type:     "cardio",
			Duration: 45,
			Calories: 450,
			Date:     time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
		},
	}

	err := repo.InsertBatch(exercises)
	require.NoError(t, err)

	// Test date range that includes all exercises
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	allInRange, err := repo.GetByDateRange(start, end)
	assert.NoError(t, err)
	assert.Len(t, allInRange, 3)

	// Test date range that includes only first two exercises
	start = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end = time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC)
	partialRange, err := repo.GetByDateRange(start, end)
	assert.NoError(t, err)
	assert.Len(t, partialRange, 2)

	// Test date range with no exercises
	start = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	end = time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)
	emptyRange, err := repo.GetByDateRange(start, end)
	assert.NoError(t, err)
	assert.Len(t, emptyRange, 0)
}

func TestPostgresRepository_Update(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Insert an exercise first
	exercise := loader.Exercise{
		Name:        "Original Exercise",
		Type:        "cardio",
		Duration:    30,
		Calories:    250,
		Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Description: "Original description",
	}
	err := repo.Insert(exercise)
	require.NoError(t, err)

	// Get the inserted exercise to get its ID
	exercises, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, exercises, 1)
	insertedExercise := exercises[0]

	// Update the exercise
	insertedExercise.Name = "Updated Exercise"
	insertedExercise.Type = "strength"
	insertedExercise.Duration = 45
	insertedExercise.Calories = 350
	insertedExercise.Description = "Updated description"

	err = repo.Update(insertedExercise)
	assert.NoError(t, err)

	// Verify the update
	updated, err := repo.GetByID(insertedExercise.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Exercise", updated.Name)
	assert.Equal(t, "strength", updated.Type)
	assert.Equal(t, 45, updated.Duration)
	assert.Equal(t, 350, updated.Calories)
	assert.Equal(t, "Updated description", updated.Description)
}

func TestPostgresRepository_Delete(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Insert an exercise first
	exercise := loader.Exercise{
		Name:     "To Delete",
		Type:     "cardio",
		Duration: 30,
		Calories: 250,
		Date:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	err := repo.Insert(exercise)
	require.NoError(t, err)

	// Get the inserted exercise to get its ID
	exercises, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, exercises, 1)
	exerciseID := exercises[0].ID

	// Delete the exercise
	err = repo.Delete(exerciseID)
	assert.NoError(t, err)

	// Verify it's deleted
	deleted, err := repo.GetByID(exerciseID)
	assert.NoError(t, err)
	assert.Nil(t, deleted)

	// Verify no exercises remain
	allExercises, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, allExercises, 0)
}

func TestPostgresRepository_Close(t *testing.T) {
	repo := setupTestDB(t)

	err := repo.Close()
	assert.NoError(t, err)

	// Verify that operations fail after closing
	err = repo.Insert(loader.Exercise{})
	assert.Error(t, err)
}

func TestPostgresRepository_DatabaseIntegration(t *testing.T) {
	// This test verifies the complete workflow
	repo := setupTestDB(t)
	defer repo.Close()

	// Step 1: Insert batch of exercises
	exercises := []loader.Exercise{
		{
			Name:        "Morning Yoga",
			Type:        "flexibility",
			Duration:    60,
			Calories:    200,
			Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Description: "Relaxing morning session",
		},
		{
			Name:        "HIIT Workout",
			Type:        "cardio",
			Duration:    30,
			Calories:    400,
			Date:        time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			Description: "High intensity interval training",
		},
	}

	err := repo.InsertBatch(exercises)
	require.NoError(t, err)

	// Step 2: Verify all exercises are retrievable
	allExercises, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, allExercises, 2)

	// Step 3: Query by type
	cardioExercises, err := repo.GetByType("cardio")
	require.NoError(t, err)
	assert.Len(t, cardioExercises, 1)
	assert.Equal(t, "HIIT Workout", cardioExercises[0].Name)

	// Step 4: Query by date range
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	dateRangeExercises, err := repo.GetByDateRange(start, end)
	require.NoError(t, err)
	assert.Len(t, dateRangeExercises, 2)

	// Step 5: Update an exercise
	exerciseToUpdate := allExercises[0]
	exerciseToUpdate.Name = "Updated Yoga Session"
	err = repo.Update(exerciseToUpdate)
	require.NoError(t, err)

	// Step 6: Verify the update
	updated, err := repo.GetByID(exerciseToUpdate.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Yoga Session", updated.Name)

	// Step 7: Delete an exercise
	err = repo.Delete(exerciseToUpdate.ID)
	require.NoError(t, err)

	// Step 8: Verify deletion
	finalExercises, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, finalExercises, 1)
	assert.Equal(t, "HIIT Workout", finalExercises[0].Name)
}
