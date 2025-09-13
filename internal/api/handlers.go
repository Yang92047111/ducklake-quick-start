package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Yang92047111/ducklake-quick-start/internal/loader"
	"github.com/Yang92047111/ducklake-quick-start/internal/storage"
	"github.com/gorilla/mux"
)

type Handler struct {
	repo storage.ExerciseRepository
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

func NewHandler(repo storage.ExerciseRepository) *Handler {
	return &Handler{repo: repo}
}

// writeJSONError writes a structured JSON error response
func (h *Handler) writeJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorResp := ErrorResponse{
		Error:   http.StatusText(code),
		Code:    code,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

// writeJSONResponse writes a successful JSON response
func (h *Handler) writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		h.writeJSONError(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.repo.GetAll()
	if err != nil {
		log.Printf("Failed to get all exercises: %v", err)
		h.writeJSONError(w, "Failed to retrieve exercises", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, exercises)
}

func (h *Handler) GetExerciseByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		h.writeJSONError(w, "Exercise ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.writeJSONError(w, "Invalid exercise ID format", http.StatusBadRequest)
		return
	}

	if id <= 0 {
		h.writeJSONError(w, "Exercise ID must be a positive integer", http.StatusBadRequest)
		return
	}

	exercise, err := h.repo.GetByID(id)
	if err != nil {
		log.Printf("Failed to get exercise by ID %d: %v", id, err)
		h.writeJSONError(w, "Failed to retrieve exercise", http.StatusInternalServerError)
		return
	}

	if exercise == nil {
		h.writeJSONError(w, fmt.Sprintf("Exercise with ID %d not found", id), http.StatusNotFound)
		return
	}

	h.writeJSONResponse(w, exercise)
}

func (h *Handler) GetExercisesByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	exerciseType, exists := vars["type"]
	if !exists {
		h.writeJSONError(w, "Exercise type is required", http.StatusBadRequest)
		return
	}

	// Validate exercise type
	exerciseType = strings.TrimSpace(strings.ToLower(exerciseType))
	if exerciseType == "" {
		h.writeJSONError(w, "Exercise type cannot be empty", http.StatusBadRequest)
		return
	}

	// Validate against known exercise types
	validTypes := []string{"cardio", "strength", "flexibility", "sports", "other"}
	isValid := false
	for _, validType := range validTypes {
		if exerciseType == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		h.writeJSONError(w, fmt.Sprintf("Invalid exercise type. Valid types are: %s", strings.Join(validTypes, ", ")), http.StatusBadRequest)
		return
	}

	exercises, err := h.repo.GetByType(exerciseType)
	if err != nil {
		log.Printf("Failed to get exercises by type %s: %v", exerciseType, err)
		h.writeJSONError(w, "Failed to retrieve exercises", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, exercises)
}

func (h *Handler) GetExercisesByDateRange(w http.ResponseWriter, r *http.Request) {
	startStr := strings.TrimSpace(r.URL.Query().Get("start"))
	endStr := strings.TrimSpace(r.URL.Query().Get("end"))

	if startStr == "" || endStr == "" {
		h.writeJSONError(w, "Both start and end date parameters are required (format: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		h.writeJSONError(w, "Invalid start date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		h.writeJSONError(w, "Invalid end date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Validate date range
	if start.After(end) {
		h.writeJSONError(w, "Start date must be before or equal to end date", http.StatusBadRequest)
		return
	}

	// Prevent overly large date ranges (more than 2 years)
	maxRange := 2 * 365 * 24 * time.Hour
	if end.Sub(start) > maxRange {
		h.writeJSONError(w, "Date range cannot exceed 2 years", http.StatusBadRequest)
		return
	}

	exercises, err := h.repo.GetByDateRange(start, end)
	if err != nil {
		log.Printf("Failed to get exercises by date range %s to %s: %v", startStr, endStr, err)
		h.writeJSONError(w, "Failed to retrieve exercises", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, exercises)
}

// Health check endpoint
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Service:   "ducklake-loader",
	}

	h.writeJSONResponse(w, response)
}

// Readiness check endpoint - checks if database is accessible
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	// Try to perform a simple database operation to check readiness
	_, err := h.repo.GetAll()
	if err != nil {
		log.Printf("Readiness check failed: %v", err)
		response := HealthResponse{
			Status:    "not ready",
			Timestamp: time.Now().UTC(),
			Service:   "ducklake-loader",
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		h.writeJSONResponse(w, response)
		return
	}

	response := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now().UTC(),
		Service:   "ducklake-loader",
	}

	h.writeJSONResponse(w, response)
}

// CreateExercise handles POST requests to create a new exercise
func (h *Handler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	var exercise loader.Exercise

	if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
		h.writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate the exercise
	validator := loader.NewValidator()
	if err := validator.Validate(exercise); err != nil {
		h.writeJSONError(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// Insert the exercise
	if err := h.repo.Insert(exercise); err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to create exercise: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.writeJSONResponse(w, map[string]interface{}{
		"message":  "Exercise created successfully",
		"exercise": exercise,
	})
}

func (h *Handler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Exercise endpoints
	r.HandleFunc("/exercises", h.GetExercises).Methods("GET")
	r.HandleFunc("/exercises", h.CreateExercise).Methods("POST")
	r.HandleFunc("/exercises/{id:[0-9]+}", h.GetExerciseByID).Methods("GET")
	r.HandleFunc("/exercises/type/{type}", h.GetExercisesByType).Methods("GET")
	r.HandleFunc("/exercises/date-range", h.GetExercisesByDateRange).Methods("GET")

	// API v1 endpoints
	r.HandleFunc("/api/v1/exercises", h.GetExercises).Methods("GET")
	r.HandleFunc("/api/v1/exercises", h.CreateExercise).Methods("POST")
	r.HandleFunc("/api/v1/exercises/{id:[0-9]+}", h.GetExerciseByID).Methods("GET")

	// Health check endpoints
	r.HandleFunc("/health", h.Health).Methods("GET")
	r.HandleFunc("/ready", h.Ready).Methods("GET")

	// Register batch and streaming routes if lakehouse repository is available
	h.RegisterTestRoutes(r)

	// Add middleware for logging requests
	r.Use(loggingMiddleware)

	return r
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request
		log.Printf("%s %s %s - %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}
