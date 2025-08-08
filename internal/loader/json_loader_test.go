package loader

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONLoader_LoadFromJSON(t *testing.T) {
	loader := NewJSONLoader()

	exercises, err := loader.LoadFromJSON("../../test/testdata/sample_exercises.json")
	require.NoError(t, err)
	require.Len(t, exercises, 3)

	// Test first exercise
	first := exercises[0]
	assert.Equal(t, 1, first.ID)
	assert.Equal(t, "Morning Run", first.Name)
	assert.Equal(t, "cardio", first.Type)
	assert.Equal(t, 30, first.Duration)
	assert.Equal(t, 300, first.Calories)
	expectedDate, _ := time.Parse(time.RFC3339, "2024-01-15T00:00:00Z")
	assert.Equal(t, expectedDate, first.Date)
	assert.Equal(t, "Easy morning jog around the park", first.Description)
}

func TestJSONLoader_LoadFromJSON_FileNotFound(t *testing.T) {
	loader := NewJSONLoader()

	_, err := loader.LoadFromJSON("nonexistent.json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open JSON file")
}
