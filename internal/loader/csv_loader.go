package loader

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type CSVLoader struct{}

func NewCSVLoader() *CSVLoader {
	return &CSVLoader{}
}

func (c *CSVLoader) LoadFromCSV(filename string) ([]Exercise, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// Skip header row
	exercises := make([]Exercise, 0, len(records)-1)
	for i, record := range records[1:] {
		exercise, err := c.parseCSVRecord(record, i+2)
		if err != nil {
			return nil, fmt.Errorf("error parsing row %d: %w", i+2, err)
		}
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (c *CSVLoader) parseCSVRecord(record []string, rowNum int) (Exercise, error) {
	if len(record) < 7 {
		return Exercise{}, fmt.Errorf("insufficient columns in row %d", rowNum)
	}

	id, err := strconv.Atoi(record[0])
	if err != nil {
		return Exercise{}, fmt.Errorf("invalid ID: %w", err)
	}

	duration, err := strconv.Atoi(record[3])
	if err != nil {
		return Exercise{}, fmt.Errorf("invalid duration: %w", err)
	}

	calories, err := strconv.Atoi(record[4])
	if err != nil {
		return Exercise{}, fmt.Errorf("invalid calories: %w", err)
	}

	date, err := time.Parse("2006-01-02", record[5])
	if err != nil {
		return Exercise{}, fmt.Errorf("invalid date format: %w", err)
	}

	return Exercise{
		ID:          id,
		Name:        record[1],
		Type:        record[2],
		Duration:    duration,
		Calories:    calories,
		Date:        date,
		Description: record[6],
	}, nil
}
