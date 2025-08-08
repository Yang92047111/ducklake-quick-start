package loader

import (
	"fmt"
	"strings"
	"time"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(exercise Exercise) error {
	var errors []string

	if exercise.ID <= 0 {
		errors = append(errors, "ID must be positive")
	}

	if strings.TrimSpace(exercise.Name) == "" {
		errors = append(errors, "name is required")
	}

	if strings.TrimSpace(exercise.Type) == "" {
		errors = append(errors, "type is required")
	}

	if exercise.Duration <= 0 {
		errors = append(errors, "duration must be positive")
	}

	if exercise.Calories < 0 {
		errors = append(errors, "calories cannot be negative")
	}

	if exercise.Date.IsZero() {
		errors = append(errors, "date is required")
	}

	if exercise.Date.After(time.Now()) {
		errors = append(errors, "date cannot be in the future")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, ", "))
	}

	return nil
}
