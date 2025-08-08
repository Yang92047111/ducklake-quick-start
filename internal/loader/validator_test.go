package loader

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidator_Validate_ValidExercise(t *testing.T) {
	validator := NewValidator()

	exercise := Exercise{
		ID:          1,
		Name:        "Running",
		Type:        "cardio",
		Duration:    30,
		Calories:    300,
		Date:        time.Now().AddDate(0, 0, -1), // Yesterday
		Description: "Morning run",
	}

	err := validator.Validate(exercise)
	assert.NoError(t, err)
}

func TestValidator_Validate_InvalidExercise(t *testing.T) {
	validator := NewValidator()

	exercise := Exercise{
		ID:          0,           // Invalid: should be positive
		Name:        "",          // Invalid: empty name
		Type:        "",          // Invalid: empty type
		Duration:    -10,         // Invalid: negative duration
		Calories:    -50,         // Invalid: negative calories
		Date:        time.Time{}, // Invalid: zero date
		Description: "Test",
	}

	err := validator.Validate(exercise)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID must be positive")
	assert.Contains(t, err.Error(), "name is required")
	assert.Contains(t, err.Error(), "type is required")
	assert.Contains(t, err.Error(), "duration must be positive")
	assert.Contains(t, err.Error(), "calories cannot be negative")
	assert.Contains(t, err.Error(), "date is required")
}

func TestValidator_Validate_FutureDate(t *testing.T) {
	validator := NewValidator()

	exercise := Exercise{
		ID:          1,
		Name:        "Running",
		Type:        "cardio",
		Duration:    30,
		Calories:    300,
		Date:        time.Now().AddDate(0, 0, 1), // Tomorrow
		Description: "Future run",
	}

	err := validator.Validate(exercise)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "date cannot be in the future")
}
