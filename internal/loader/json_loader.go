package loader

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONLoader struct{}

func NewJSONLoader() *JSONLoader {
	return &JSONLoader{}
}

func (j *JSONLoader) LoadFromJSON(filename string) ([]Exercise, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var exercises []Exercise
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&exercises); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return exercises, nil
}
