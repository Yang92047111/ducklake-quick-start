package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/Yang92047111/ducklake-quick-start/internal/loader"
)

type MemoryRepository struct {
	exercises map[int]loader.Exercise
	nextID    int
	mutex     sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		exercises: make(map[int]loader.Exercise),
		nextID:    1,
	}
}

func (m *MemoryRepository) Insert(exercise loader.Exercise) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exercise.ID = m.nextID
	m.exercises[exercise.ID] = exercise
	m.nextID++
	return nil
}

func (m *MemoryRepository) InsertBatch(exercises []loader.Exercise) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, exercise := range exercises {
		exercise.ID = m.nextID
		m.exercises[exercise.ID] = exercise
		m.nextID++
	}
	return nil
}

func (m *MemoryRepository) GetByID(id int) (*loader.Exercise, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	exercise, exists := m.exercises[id]
	if !exists {
		return nil, nil
	}
	return &exercise, nil
}

func (m *MemoryRepository) GetByDateRange(start, end time.Time) ([]loader.Exercise, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var result []loader.Exercise
	for _, exercise := range m.exercises {
		if (exercise.Date.Equal(start) || exercise.Date.After(start)) &&
			(exercise.Date.Equal(end) || exercise.Date.Before(end)) {
			result = append(result, exercise)
		}
	}
	return result, nil
}

func (m *MemoryRepository) GetByType(exerciseType string) ([]loader.Exercise, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var result []loader.Exercise
	for _, exercise := range m.exercises {
		if exercise.Type == exerciseType {
			result = append(result, exercise)
		}
	}
	return result, nil
}

func (m *MemoryRepository) GetAll() ([]loader.Exercise, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]loader.Exercise, 0, len(m.exercises))
	for _, exercise := range m.exercises {
		result = append(result, exercise)
	}
	return result, nil
}

func (m *MemoryRepository) Update(exercise loader.Exercise) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.exercises[exercise.ID]; !exists {
		return fmt.Errorf("exercise with ID %d not found", exercise.ID)
	}
	m.exercises[exercise.ID] = exercise
	return nil
}

func (m *MemoryRepository) Delete(id int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.exercises[id]; !exists {
		return fmt.Errorf("exercise with ID %d not found", id)
	}
	delete(m.exercises, id)
	return nil
}

func (m *MemoryRepository) Close() error {
	return nil
}
