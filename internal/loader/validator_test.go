package loader

import (
	"strings"
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
		ID:          -1,          // Invalid: negative ID
		Name:        "",          // Invalid: empty name
		Type:        "",          // Invalid: empty type
		Duration:    -10,         // Invalid: negative duration
		Calories:    -50,         // Invalid: negative calories
		Date:        time.Time{}, // Invalid: zero date
		Description: "Test",
	}

	err := validator.Validate(exercise)
	assert.Error(t, err)

	// Check new error message format
	errMsg := err.Error()
	assert.Contains(t, errMsg, "id: must be non-negative")
	assert.Contains(t, errMsg, "name: is required and cannot be empty")
	assert.Contains(t, errMsg, "type: is required and cannot be empty")
	assert.Contains(t, errMsg, "duration: must be positive")
	assert.Contains(t, errMsg, "calories: cannot be negative")
	assert.Contains(t, errMsg, "date: is required")
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
	assert.Contains(t, err.Error(), "date: cannot be in the future")
}

// New comprehensive tests for enhanced validation

func TestValidator_ValidateName(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name         string
		exerciseName string
		expectError  bool
		errorMsg     string
	}{
		{"valid name", "Morning Run", false, ""},
		{"empty name", "", true, "is required and cannot be empty"},
		{"short name", "A", true, "must be at least 2 characters long"},
		{"long name", strings.Repeat("A", 101), true, "must be no more than 100 characters long"},
		{"name with hyphen", "High-Intensity Training", false, ""},
		{"name with apostrophe", "John's Workout", false, ""},
		{"name with numbers", "Workout 123", false, ""},
		{"name with invalid chars", "Workout@#$", true, "contains invalid characters"},
		{"spaces only", "   ", true, "is required and cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exercise := Exercise{
				ID:       1,
				Name:     tt.exerciseName,
				Type:     "cardio",
				Duration: 30,
				Calories: 300,
				Date:     time.Now().AddDate(0, 0, -1),
			}

			err := validator.Validate(exercise)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateType(t *testing.T) {
	validator := NewValidator()

	validTypes := []string{"cardio", "strength", "flexibility", "sports", "other"}

	for _, validType := range validTypes {
		t.Run("valid type: "+validType, func(t *testing.T) {
			exercise := Exercise{
				ID:       1,
				Name:     "Test Exercise",
				Type:     validType,
				Duration: 30,
				Calories: 300,
				Date:     time.Now().AddDate(0, 0, -1),
			}

			err := validator.Validate(exercise)
			assert.NoError(t, err)
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		exercise := Exercise{
			ID:       1,
			Name:     "Test Exercise",
			Type:     "invalid",
			Duration: 30,
			Calories: 300,
			Date:     time.Now().AddDate(0, 0, -1),
		}

		err := validator.Validate(exercise)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be one of")
	})

	t.Run("case insensitive type", func(t *testing.T) {
		exercise := Exercise{
			ID:       1,
			Name:     "Test Exercise",
			Type:     "CARDIO",
			Duration: 30,
			Calories: 300,
			Date:     time.Now().AddDate(0, 0, -1),
		}

		err := validator.Validate(exercise)
		assert.NoError(t, err)
	})
}

func TestValidator_ValidateDuration(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		duration    int
		expectError bool
		errorMsg    string
	}{
		{"valid duration", 30, false, ""},
		{"zero duration", 0, true, "must be positive"},
		{"negative duration", -10, true, "must be positive"},
		{"maximum valid duration", 1440, false, ""}, // 24 hours
		{"excessive duration", 1441, true, "cannot exceed 24 hours"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exercise := Exercise{
				ID:       1,
				Name:     "Test Exercise",
				Type:     "cardio",
				Duration: tt.duration,
				Calories: 300,
				Date:     time.Now().AddDate(0, 0, -1),
			}

			err := validator.Validate(exercise)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateCalories(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		calories    int
		expectError bool
		errorMsg    string
	}{
		{"valid calories", 300, false, ""},
		{"zero calories", 0, false, ""}, // Zero calories is valid
		{"negative calories", -10, true, "cannot be negative"},
		{"maximum valid calories", 10000, false, ""},
		{"excessive calories", 10001, true, "seems unreasonably high"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exercise := Exercise{
				ID:       1,
				Name:     "Test Exercise",
				Type:     "cardio",
				Duration: 30,
				Calories: tt.calories,
				Date:     time.Now().AddDate(0, 0, -1),
			}

			err := validator.Validate(exercise)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateDate(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		date        time.Time
		expectError bool
		errorMsg    string
	}{
		{"valid date yesterday", time.Now().AddDate(0, 0, -1), false, ""},
		{"valid date today", time.Now(), false, ""}, // Should be valid if not future
		{"future date", time.Now().AddDate(0, 0, 1), true, "cannot be in the future"},
		{"zero date", time.Time{}, true, "is required"},
		{"very old date", time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC), true, "cannot be before year 1900"},
		{"edge case 1900", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exercise := Exercise{
				ID:       1,
				Name:     "Test Exercise",
				Type:     "cardio",
				Duration: 30,
				Calories: 300,
				Date:     tt.date,
			}

			err := validator.Validate(exercise)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateDescription(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		description string
		expectError bool
		errorMsg    string
	}{
		{"valid description", "A normal workout description", false, ""},
		{"empty description", "", false, ""}, // Empty description is valid
		{"long description", strings.Repeat("A", 1001), true, "must be no more than 1000 characters"},
		{"max valid description", strings.Repeat("A", 1000), false, ""},
		{"suspicious content script", "<script>alert('xss')</script>", true, "contains potentially unsafe content"},
		{"suspicious content javascript", "javascript:alert('xss')", true, "contains potentially unsafe content"},
		{"normal HTML-like text", "This is <good> content", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exercise := Exercise{
				ID:          1,
				Name:        "Test Exercise",
				Type:        "cardio",
				Duration:    30,
				Calories:    300,
				Date:        time.Now().AddDate(0, 0, -1),
				Description: tt.description,
			}

			err := validator.Validate(exercise)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidatePartial(t *testing.T) {
	validator := NewValidator()

	// Test partial validation where only some fields are set
	exercise := Exercise{
		ID:   1,
		Name: "Updated Name", // Only name is updated
		// Other fields are zero values, should be ignored in partial validation
	}

	err := validator.ValidatePartial(exercise)
	assert.NoError(t, err)

	// Test partial validation with invalid field
	exercise = Exercise{
		ID:   1,
		Name: "A", // Too short
	}

	err = validator.ValidatePartial(exercise)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be at least 2 characters long")
}

func TestValidationError_Error(t *testing.T) {
	// Test ValidationError error message formatting
	err := ValidationError{
		Field:   "name",
		Message: "is required",
		Value:   "empty",
	}

	expected := "name: is required (got: empty)"
	assert.Equal(t, expected, err.Error())

	// Test without value
	err = ValidationError{
		Field:   "date",
		Message: "is required",
	}

	expected = "date: is required"
	assert.Equal(t, expected, err.Error())
}

func TestValidationErrors_Error(t *testing.T) {
	errors := ValidationErrors{
		{Field: "name", Message: "is required"},
		{Field: "type", Message: "is invalid"},
	}

	errMsg := errors.Error()
	assert.Contains(t, errMsg, "validation failed:")
	assert.Contains(t, errMsg, "name: is required")
	assert.Contains(t, errMsg, "type: is invalid")
}

func TestValidationErrors_HasErrors(t *testing.T) {
	var errors ValidationErrors
	assert.False(t, errors.HasErrors())

	errors = append(errors, ValidationError{Field: "test", Message: "error"})
	assert.True(t, errors.HasErrors())
}
