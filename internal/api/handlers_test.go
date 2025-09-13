package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Yang92047111/ducklake-quick-start/internal/loader"
	"github.com/Yang92047111/ducklake-quick-start/internal/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestHandler() *Handler {
	repo := storage.NewMemoryRepository()

	// Add test data
	exercises := []loader.Exercise{
		{
			Name:        "Running",
			Type:        "cardio",
			Duration:    30,
			Calories:    300,
			Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Description: "Morning run",
		},
		{
			Name:        "Push-ups",
			Type:        "strength",
			Duration:    15,
			Calories:    100,
			Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Description: "Strength training",
		},
	}

	repo.InsertBatch(exercises)
	return NewHandler(repo)
}

func TestHandler_GetExercises(t *testing.T) {
	handler := setupTestHandler()

	req, err := http.NewRequest("GET", "/exercises", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetExercises(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var exercises []loader.Exercise
	err = json.Unmarshal(rr.Body.Bytes(), &exercises)
	require.NoError(t, err)
	assert.Len(t, exercises, 2)
}

func TestHandler_GetExerciseByID(t *testing.T) {
	handler := setupTestHandler()
	router := mux.NewRouter()
	router.HandleFunc("/exercises/{id:[0-9]+}", handler.GetExerciseByID).Methods("GET")

	// Test valid ID
	req, err := http.NewRequest("GET", "/exercises/1", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var exercise loader.Exercise
	err = json.Unmarshal(rr.Body.Bytes(), &exercise)
	require.NoError(t, err)
	assert.Equal(t, "Running", exercise.Name)
}

func TestHandler_GetExerciseByID_NotFound(t *testing.T) {
	handler := setupTestHandler()
	router := mux.NewRouter()
	router.HandleFunc("/exercises/{id:[0-9]+}", handler.GetExerciseByID).Methods("GET")

	req, err := http.NewRequest("GET", "/exercises/999", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestHandler_GetExercisesByType(t *testing.T) {
	handler := setupTestHandler()
	router := mux.NewRouter()
	router.HandleFunc("/exercises/type/{type}", handler.GetExercisesByType).Methods("GET")

	req, err := http.NewRequest("GET", "/exercises/type/cardio", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var exercises []loader.Exercise
	err = json.Unmarshal(rr.Body.Bytes(), &exercises)
	require.NoError(t, err)
	assert.Len(t, exercises, 1)
	assert.Equal(t, "cardio", exercises[0].Type)
}

func TestHandler_GetExercisesByDateRange(t *testing.T) {
	handler := setupTestHandler()

	req, err := http.NewRequest("GET", "/exercises/date-range?start=2024-01-15&end=2024-01-15", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetExercisesByDateRange(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var exercises []loader.Exercise
	err = json.Unmarshal(rr.Body.Bytes(), &exercises)
	require.NoError(t, err)
	assert.Len(t, exercises, 2)
}

func TestHandler_GetExercisesByDateRange_MissingParams(t *testing.T) {
	handler := setupTestHandler()

	req, err := http.NewRequest("GET", "/exercises/date-range", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetExercisesByDateRange(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_SetupRoutes(t *testing.T) {
	handler := setupTestHandler()
	router := handler.SetupRoutes()

	// Test that routes are properly configured
	req, err := http.NewRequest("GET", "/exercises", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

// New comprehensive tests for improved functionality

func TestHandler_GetExerciseByID_InvalidID(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name       string
		id         string
		expectCode int
		useRouter  bool
	}{
		{"negative ID", "-1", http.StatusNotFound, true},           // Router will reject negative numbers
		{"zero ID", "0", http.StatusBadRequest, false},             // Handler will reject zero
		{"non-numeric ID", "abc", http.StatusNotFound, true},       // Router will reject non-numeric
		{"float ID", "1.5", http.StatusNotFound, true},             // Router will reject float
		{"positive invalid ID", "999", http.StatusNotFound, false}, // Valid format but doesn't exist
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useRouter {
				// Test with router that has regex pattern
				req, err := http.NewRequest("GET", "/exercises/"+tt.id, nil)
				require.NoError(t, err)

				router := mux.NewRouter()
				router.HandleFunc("/exercises/{id:[0-9]+}", handler.GetExerciseByID).Methods("GET")

				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				assert.Equal(t, tt.expectCode, rr.Code)
			} else {
				// Test direct handler call for valid format but non-existent ID
				req, err := http.NewRequest("GET", "/exercises/"+tt.id, nil)
				require.NoError(t, err)
				req = mux.SetURLVars(req, map[string]string{"id": tt.id})

				rr := httptest.NewRecorder()
				handler.GetExerciseByID(rr, req)

				assert.Equal(t, tt.expectCode, rr.Code)
			}
		})
	}
}

func TestHandler_GetExercisesByType_Validation(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name         string
		exerciseType string
		expectCode   int
		useRouter    bool
	}{
		{"valid cardio type", "cardio", http.StatusOK, true},
		{"valid strength type", "strength", http.StatusOK, true},
		{"valid flexibility type", "flexibility", http.StatusOK, true},
		{"valid sports type", "sports", http.StatusOK, true},
		{"valid other type", "other", http.StatusOK, true},
		{"uppercase type", "CARDIO", http.StatusOK, true},
		{"mixed case type", "Strength", http.StatusOK, true},
		{"invalid type", "invalid", http.StatusBadRequest, false},
		{"empty type", "", http.StatusNotFound, true}, // Router will reject empty paths
		{"spaces type", "   ", http.StatusBadRequest, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useRouter {
				req, err := http.NewRequest("GET", "/exercises/type/"+tt.exerciseType, nil)
				require.NoError(t, err)

				router := mux.NewRouter()
				router.HandleFunc("/exercises/type/{type}", handler.GetExercisesByType).Methods("GET")

				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				assert.Equal(t, tt.expectCode, rr.Code)
			} else {
				// Test direct handler call
				req, err := http.NewRequest("GET", "/exercises/type/"+tt.exerciseType, nil)
				require.NoError(t, err)
				req = mux.SetURLVars(req, map[string]string{"type": tt.exerciseType})

				rr := httptest.NewRecorder()
				handler.GetExercisesByType(rr, req)

				assert.Equal(t, tt.expectCode, rr.Code)

				if tt.expectCode != http.StatusOK {
					var errorResp ErrorResponse
					err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
					assert.NoError(t, err)
					assert.NotEmpty(t, errorResp.Message)
				}
			}
		})
	}
}

func TestHandler_GetExercisesByDateRange_Validation(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name       string
		start      string
		end        string
		expectCode int
	}{
		{"valid date range", "2024-01-01", "2024-01-31", http.StatusOK},
		{"start after end", "2024-01-31", "2024-01-01", http.StatusBadRequest},
		{"invalid start date", "invalid", "2024-01-31", http.StatusBadRequest},
		{"invalid end date", "2024-01-01", "invalid", http.StatusBadRequest},
		{"range too large", "2020-01-01", "2024-01-01", http.StatusBadRequest},
		{"empty start", "", "2024-01-31", http.StatusBadRequest},
		{"empty end", "2024-01-01", "", http.StatusBadRequest},
		{"same date", "2024-01-15", "2024-01-15", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/exercises/date-range"
			if tt.start != "" || tt.end != "" {
				url += "?start=" + tt.start + "&end=" + tt.end
			}

			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.GetExercisesByDateRange(rr, req)

			assert.Equal(t, tt.expectCode, rr.Code)

			if tt.expectCode != http.StatusOK {
				var errorResp ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.NotEmpty(t, errorResp.Message)
			}
		})
	}
}

func TestHandler_Health(t *testing.T) {
	handler := setupTestHandler()

	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.Health(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var healthResp HealthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &healthResp)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", healthResp.Status)
	assert.Equal(t, "ducklake-loader", healthResp.Service)
	assert.False(t, healthResp.Timestamp.IsZero())
}

func TestHandler_Ready(t *testing.T) {
	handler := setupTestHandler()

	req, err := http.NewRequest("GET", "/ready", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.Ready(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var readyResp HealthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &readyResp)
	assert.NoError(t, err)
	assert.Equal(t, "ready", readyResp.Status)
	assert.Equal(t, "ducklake-loader", readyResp.Service)
}

func TestHandler_SetupRoutes_NewEndpoints(t *testing.T) {
	handler := setupTestHandler()
	router := handler.SetupRoutes()

	tests := []struct {
		name     string
		method   string
		path     string
		expected int
	}{
		{"health endpoint", "GET", "/health", http.StatusOK},
		{"ready endpoint", "GET", "/ready", http.StatusOK},
		{"exercises endpoint", "GET", "/exercises", http.StatusOK},
		{"exercise by id endpoint", "GET", "/exercises/999", http.StatusNotFound}, // ID 999 doesn't exist in test data
		{"exercise by type endpoint", "GET", "/exercises/type/cardio", http.StatusOK},
		{"date range endpoint", "GET", "/exercises/date-range?start=2024-01-01&end=2024-01-31", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expected, rr.Code)
		})
	}
}

func TestHandler_JSONErrorResponses(t *testing.T) {
	handler := setupTestHandler()

	// Test that error responses are properly formatted JSON
	req, err := http.NewRequest("GET", "/exercises/type/invalid", nil)
	require.NoError(t, err)

	router := mux.NewRouter()
	router.HandleFunc("/exercises/type/{type}", handler.GetExercisesByType).Methods("GET")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var errorResp ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Bad Request", errorResp.Error)
	assert.Equal(t, 400, errorResp.Code)
	assert.Contains(t, errorResp.Message, "Invalid exercise type")
}
