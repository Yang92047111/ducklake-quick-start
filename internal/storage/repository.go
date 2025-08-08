package storage

import (
	"time"

	"github.com/yourname/ducklake-loader/internal/loader"
)

// ExerciseRepository defines the interface for exercise data persistence
type ExerciseRepository interface {
	Insert(exercise loader.Exercise) error
	InsertBatch(exercises []loader.Exercise) error
	GetByID(id int) (*loader.Exercise, error)
	GetByDateRange(start, end time.Time) ([]loader.Exercise, error)
	GetByType(exerciseType string) ([]loader.Exercise, error)
	GetAll() ([]loader.Exercise, error)
	Update(exercise loader.Exercise) error
	Delete(id int) error
	Close() error
}
