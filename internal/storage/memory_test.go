package storage

import (
	"testing"
	"time"

	"github.com/Yang92047111/ducklake-quick-start/internal/loader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRepository_Insert(t *testing.T) {
	repo := NewMemoryRepository()

	exercise := loader.Exercise{
		Name:        "Running",
		Type:        "cardio",
		Duration:    30,
		Calories:    300,
		Date:        time.Now(),
		Description: "Morning run",
	}

	err := repo.Insert(exercise)
	require.NoError(t, err)

	// Verify the exercise was inserted with an ID
	exercises, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, exercises, 1)
	assert.Equal(t, 1, exercises[0].ID)
	assert.Equal(t, "Running", exercises[0].Name)
}

func TestMemoryRepository_GetByID(t *testing.T) {
	repo := NewMemoryRepository()

	exercise := loader.Exercise{
		Name:        "Running",
		Type:        "cardio",
		Duration:    30,
		Calories:    300,
		Date:        time.Now(),
		Description: "Morning run",
	}

	err := repo.Insert(exercise)
	require.NoError(t, err)

	// Get by ID
	retrieved, err := repo.GetByID(1)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, "Running", retrieved.Name)

	// Get non-existent ID
	notFound, err := repo.GetByID(999)
	require.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestMemoryRepository_GetByType(t *testing.T) {
	repo := NewMemoryRepository()

	exercises := []loader.Exercise{
		{Name: "Running", Type: "cardio", Duration: 30, Calories: 300, Date: time.Now()},
		{Name: "Push-ups", Type: "strength", Duration: 15, Calories: 100, Date: time.Now()},
		{Name: "Swimming", Type: "cardio", Duration: 40, Calories: 350, Date: time.Now()},
	}

	err := repo.InsertBatch(exercises)
	require.NoError(t, err)

	// Get cardio exercises
	cardioExercises, err := repo.GetByType("cardio")
	require.NoError(t, err)
	assert.Len(t, cardioExercises, 2)

	// Get strength exercises
	strengthExercises, err := repo.GetByType("strength")
	require.NoError(t, err)
	assert.Len(t, strengthExercises, 1)
	assert.Equal(t, "Push-ups", strengthExercises[0].Name)
}
func TestMemoryRepository_GetByDateRange(t *testing.T) {
	repo := NewMemoryRepository()

	date1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC)

	exercises := []loader.Exercise{
		{Name: "Exercise1", Type: "cardio", Duration: 30, Calories: 300, Date: date1},
		{Name: "Exercise2", Type: "strength", Duration: 15, Calories: 100, Date: date2},
		{Name: "Exercise3", Type: "cardio", Duration: 40, Calories: 350, Date: date3},
	}

	err := repo.InsertBatch(exercises)
	require.NoError(t, err)

	// Get exercises in date range
	rangeExercises, err := repo.GetByDateRange(date1, date2)
	require.NoError(t, err)
	assert.Len(t, rangeExercises, 2)
}

func TestMemoryRepository_Update(t *testing.T) {
	repo := NewMemoryRepository()

	exercise := loader.Exercise{
		Name:        "Running",
		Type:        "cardio",
		Duration:    30,
		Calories:    300,
		Date:        time.Now(),
		Description: "Morning run",
	}

	err := repo.Insert(exercise)
	require.NoError(t, err)

	// Update the exercise
	updatedExercise := loader.Exercise{
		ID:          1,
		Name:        "Updated Running",
		Type:        "cardio",
		Duration:    45,
		Calories:    400,
		Date:        time.Now(),
		Description: "Updated morning run",
	}

	err = repo.Update(updatedExercise)
	require.NoError(t, err)

	// Verify update
	retrieved, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, "Updated Running", retrieved.Name)
	assert.Equal(t, 45, retrieved.Duration)
}

func TestMemoryRepository_Update_NotFound(t *testing.T) {
	repo := NewMemoryRepository()

	exercise := loader.Exercise{
		ID:          999,
		Name:        "Non-existent",
		Type:        "cardio",
		Duration:    30,
		Calories:    300,
		Date:        time.Now(),
		Description: "Does not exist",
	}

	err := repo.Update(exercise)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestMemoryRepository_Delete(t *testing.T) {
	repo := NewMemoryRepository()

	exercise := loader.Exercise{
		Name:        "Running",
		Type:        "cardio",
		Duration:    30,
		Calories:    300,
		Date:        time.Now(),
		Description: "Morning run",
	}

	err := repo.Insert(exercise)
	require.NoError(t, err)

	// Delete the exercise
	err = repo.Delete(1)
	require.NoError(t, err)

	// Verify deletion
	retrieved, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestMemoryRepository_Delete_NotFound(t *testing.T) {
	repo := NewMemoryRepository()

	err := repo.Delete(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestMemoryRepository_Close(t *testing.T) {
	repo := NewMemoryRepository()

	err := repo.Close()
	assert.NoError(t, err)
}
