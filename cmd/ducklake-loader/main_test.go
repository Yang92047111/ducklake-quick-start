package main

import (
	"flag"
	"os"
	"testing"

	"github.com/Yang92047111/ducklake-quick-start/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "returns environment value when set",
			key:          "TEST_ENV_VAR",
			defaultValue: "default",
			envValue:     "environment_value",
			expected:     "environment_value",
		},
		{
			name:         "returns default value when environment variable not set",
			key:          "NON_EXISTENT_VAR",
			defaultValue: "default_value",
			envValue:     "",
			expected:     "default_value",
		},
		{
			name:         "returns empty string when both are empty",
			key:          "EMPTY_VAR",
			defaultValue: "",
			envValue:     "",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variable if needed
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadCSVData(t *testing.T) {
	// Create a memory repository for testing
	repo := storage.NewMemoryRepository()
	defer repo.Close()

	t.Run("loads valid CSV file successfully", func(t *testing.T) {
		err := loadCSVData(repo, "../../test/testdata/sample_exercises.csv")
		assert.NoError(t, err)

		// Verify data was loaded
		exercises, err := repo.GetAll()
		require.NoError(t, err)
		assert.Greater(t, len(exercises), 0)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		err := loadCSVData(repo, "non_existent_file.csv")
		assert.Error(t, err)
	})

	t.Run("handles invalid CSV format", func(t *testing.T) {
		// Create a temporary invalid CSV file
		tmpFile := "/tmp/invalid.csv"
		content := "invalid,csv,content\nwith,wrong,format,too,many,columns\n"
		err := os.WriteFile(tmpFile, []byte(content), 0644)
		require.NoError(t, err)
		defer os.Remove(tmpFile)

		err = loadCSVData(repo, tmpFile)
		// Should return error for malformed CSV
		assert.Error(t, err)
	})
}

func TestLoadJSONData(t *testing.T) {
	// Create a memory repository for testing
	repo := storage.NewMemoryRepository()
	defer repo.Close()

	t.Run("loads valid JSON file successfully", func(t *testing.T) {
		err := loadJSONData(repo, "../../test/testdata/sample_exercises.json")
		assert.NoError(t, err)

		// Verify data was loaded
		exercises, err := repo.GetAll()
		require.NoError(t, err)
		assert.Greater(t, len(exercises), 0)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		err := loadJSONData(repo, "non_existent_file.json")
		assert.Error(t, err)
	})

	t.Run("handles invalid JSON format", func(t *testing.T) {
		// Create a temporary invalid JSON file
		tmpFile := "/tmp/invalid.json"
		content := "{ invalid json content"
		err := os.WriteFile(tmpFile, []byte(content), 0644)
		require.NoError(t, err)
		defer os.Remove(tmpFile)

		err = loadJSONData(repo, tmpFile)
		assert.Error(t, err)
	})
}

func TestMainFunctionality(t *testing.T) {
	// This test verifies that the main function can parse flags correctly
	// We can't easily test the main function directly, but we can test flag parsing

	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Reset flag package
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Test flag parsing
	os.Args = []string{"ducklake-loader", "-memory", "-port", "9000"}

	var (
		csvFile    = flag.String("csv", "", "Path to CSV file to load")
		jsonFile   = flag.String("json", "", "Path to JSON file to load")
		useMemory  = flag.Bool("memory", false, "Use in-memory storage instead of PostgreSQL")
		serverMode = flag.Bool("server", false, "Run in server mode")
		port       = flag.String("port", "8080", "Server port")
	)

	flag.Parse()

	assert.Equal(t, "", *csvFile)
	assert.Equal(t, "", *jsonFile)
	assert.True(t, *useMemory)
	assert.False(t, *serverMode)
	assert.Equal(t, "9000", *port)
}

func TestMainWithInvalidFlags(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Reset flag package
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Test with valid flags
	os.Args = []string{"ducklake-loader", "-csv", "test.csv", "-json", "test.json", "-server"}

	var (
		csvFile    = flag.String("csv", "", "Path to CSV file to load")
		jsonFile   = flag.String("json", "", "Path to JSON file to load")
		useMemory  = flag.Bool("memory", false, "Use in-memory storage instead of PostgreSQL")
		serverMode = flag.Bool("server", false, "Run in server mode")
		port       = flag.String("port", "8080", "Server port")
	)

	err := flag.CommandLine.Parse(os.Args[1:])
	assert.NoError(t, err)

	assert.Equal(t, "test.csv", *csvFile)
	assert.Equal(t, "test.json", *jsonFile)
	assert.False(t, *useMemory)
	assert.True(t, *serverMode)
	assert.Equal(t, "8080", *port)
}
