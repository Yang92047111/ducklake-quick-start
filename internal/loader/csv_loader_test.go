package loader

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSVLoader_LoadFromCSV(t *testing.T) {
	loader := NewCSVLoader()

	exercises, err := loader.LoadFromCSV("../../test/testdata/sample_exercises.csv")
	require.NoError(t, err)
	require.Len(t, exercises, 5)

	// Test first exercise
	first := exercises[0]
	assert.Equal(t, 1, first.ID)
	assert.Equal(t, "Morning Run", first.Name)
	assert.Equal(t, "cardio", first.Type)
	assert.Equal(t, 30, first.Duration)
	assert.Equal(t, 300, first.Calories)
	expectedDate, _ := time.Parse("2006-01-02", "2024-01-15")
	assert.Equal(t, expectedDate, first.Date)
	assert.Equal(t, "Easy morning jog around the park", first.Description)
}

func TestCSVLoader_LoadFromCSV_FileNotFound(t *testing.T) {
	loader := NewCSVLoader()

	_, err := loader.LoadFromCSV("nonexistent.csv")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open CSV file")
}
