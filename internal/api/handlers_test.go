package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourname/ducklake-loader/internal/loader"
	"github.com/yourname/ducklake-loader/internal/storage"
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
