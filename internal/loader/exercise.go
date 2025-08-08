package loader

import (
	"time"
)

// Exercise represents a single exercise record from DuckLake
type Exercise struct {
	ID          int       `json:"id" csv:"id"`
	Name        string    `json:"name" csv:"name"`
	Type        string    `json:"type" csv:"type"`
	Duration    int       `json:"duration" csv:"duration"` // in minutes
	Calories    int       `json:"calories" csv:"calories"`
	Date        time.Time `json:"date" csv:"date"`
	Description string    `json:"description" csv:"description"`
}

// ExerciseLoader defines the interface for loading exercise data
type ExerciseLoader interface {
	LoadFromCSV(filename string) ([]Exercise, error)
	LoadFromJSON(filename string) ([]Exercise, error)
	Validate(exercise Exercise) error
}
