package loader

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

type Validator struct {
	ValidTypes []string
}

func NewValidator() *Validator {
	return &Validator{
		ValidTypes: []string{"cardio", "strength", "flexibility", "sports", "other"},
	}
}

// ValidationError provides structured validation error information
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

func (ve ValidationError) Error() string {
	if ve.Value != "" {
		return fmt.Sprintf("%s: %s (got: %s)", ve.Field, ve.Message, ve.Value)
	}
	return fmt.Sprintf("%s: %s", ve.Field, ve.Message)
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(messages, "; "))
}

func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

func (v *Validator) Validate(exercise Exercise) error {
	var errors ValidationErrors

	// Validate ID (allow 0 for new records)
	if exercise.ID < 0 {
		errors = append(errors, ValidationError{
			Field:   "id",
			Message: "must be non-negative",
			Value:   fmt.Sprintf("%d", exercise.ID),
		})
	}

	// Validate name
	errors = append(errors, v.validateName(exercise.Name)...)

	// Validate type
	errors = append(errors, v.validateType(exercise.Type)...)

	// Validate duration
	errors = append(errors, v.validateDuration(exercise.Duration)...)

	// Validate calories
	errors = append(errors, v.validateCalories(exercise.Calories)...)

	// Validate date
	errors = append(errors, v.validateDate(exercise.Date)...)

	// Validate description
	errors = append(errors, v.validateDescription(exercise.Description)...)

	if errors.HasErrors() {
		return errors
	}

	return nil
}

func (v *Validator) validateName(name string) []ValidationError {
	var errors []ValidationError

	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "is required and cannot be empty",
		})
		return errors
	}

	if len(trimmed) < 2 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "must be at least 2 characters long",
			Value:   trimmed,
		})
	}

	if len(trimmed) > 100 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "must be no more than 100 characters long",
			Value:   fmt.Sprintf("%d chars", len(trimmed)),
		})
	}

	// Check for valid characters (letters, numbers, spaces, basic punctuation)
	if !isValidName(trimmed) {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "contains invalid characters (only letters, numbers, spaces, hyphens, and apostrophes allowed)",
			Value:   trimmed,
		})
	}

	return errors
}

func (v *Validator) validateType(exerciseType string) []ValidationError {
	var errors []ValidationError

	trimmed := strings.TrimSpace(strings.ToLower(exerciseType))
	if trimmed == "" {
		errors = append(errors, ValidationError{
			Field:   "type",
			Message: "is required and cannot be empty",
		})
		return errors
	}

	// Check if type is in valid list
	isValid := false
	for _, validType := range v.ValidTypes {
		if trimmed == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		errors = append(errors, ValidationError{
			Field:   "type",
			Message: fmt.Sprintf("must be one of: %s", strings.Join(v.ValidTypes, ", ")),
			Value:   trimmed,
		})
	}

	return errors
}

func (v *Validator) validateDuration(duration int) []ValidationError {
	var errors []ValidationError

	if duration <= 0 {
		errors = append(errors, ValidationError{
			Field:   "duration",
			Message: "must be positive (in minutes)",
			Value:   fmt.Sprintf("%d", duration),
		})
	}

	// Reasonable upper limit (24 hours = 1440 minutes)
	if duration > 1440 {
		errors = append(errors, ValidationError{
			Field:   "duration",
			Message: "cannot exceed 24 hours (1440 minutes)",
			Value:   fmt.Sprintf("%d", duration),
		})
	}

	return errors
}

func (v *Validator) validateCalories(calories int) []ValidationError {
	var errors []ValidationError

	if calories < 0 {
		errors = append(errors, ValidationError{
			Field:   "calories",
			Message: "cannot be negative",
			Value:   fmt.Sprintf("%d", calories),
		})
	}

	// Reasonable upper limit (10,000 calories is extreme but possible)
	if calories > 10000 {
		errors = append(errors, ValidationError{
			Field:   "calories",
			Message: "seems unreasonably high (max 10,000)",
			Value:   fmt.Sprintf("%d", calories),
		})
	}

	return errors
}

func (v *Validator) validateDate(date time.Time) []ValidationError {
	var errors []ValidationError

	if date.IsZero() {
		errors = append(errors, ValidationError{
			Field:   "date",
			Message: "is required",
		})
		return errors
	}

	now := time.Now()
	if date.After(now) {
		errors = append(errors, ValidationError{
			Field:   "date",
			Message: "cannot be in the future",
			Value:   date.Format("2006-01-02"),
		})
	}

	// Don't allow dates too far in the past (e.g., before 1900)
	if date.Year() < 1900 {
		errors = append(errors, ValidationError{
			Field:   "date",
			Message: "cannot be before year 1900",
			Value:   date.Format("2006-01-02"),
		})
	}

	return errors
}

func (v *Validator) validateDescription(description string) []ValidationError {
	var errors []ValidationError

	if len(description) > 1000 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "must be no more than 1000 characters long",
			Value:   fmt.Sprintf("%d chars", len(description)),
		})
	}

	// Check for potentially malicious content (basic check)
	if containsSuspiciousContent(description) {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "contains potentially unsafe content",
		})
	}

	return errors
}

// Helper functions

func isValidName(name string) bool {
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsSpace(r) && r != '-' && r != '\'' {
			return false
		}
	}
	return true
}

func containsSuspiciousContent(text string) bool {
	suspicious := []string{"<script", "javascript:", "vbscript:", "onload=", "onerror=", "eval("}
	lower := strings.ToLower(text)
	for _, pattern := range suspicious {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

// ValidatePartial validates only non-zero/non-empty fields, useful for updates
func (v *Validator) ValidatePartial(exercise Exercise) error {
	var errors ValidationErrors

	if exercise.ID < 0 {
		errors = append(errors, ValidationError{
			Field:   "id",
			Message: "must be non-negative",
			Value:   fmt.Sprintf("%d", exercise.ID),
		})
	}

	if exercise.Name != "" {
		errors = append(errors, v.validateName(exercise.Name)...)
	}

	if exercise.Type != "" {
		errors = append(errors, v.validateType(exercise.Type)...)
	}

	if exercise.Duration != 0 {
		errors = append(errors, v.validateDuration(exercise.Duration)...)
	}

	if exercise.Calories != 0 {
		errors = append(errors, v.validateCalories(exercise.Calories)...)
	}

	if !exercise.Date.IsZero() {
		errors = append(errors, v.validateDate(exercise.Date)...)
	}

	if exercise.Description != "" {
		errors = append(errors, v.validateDescription(exercise.Description)...)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
