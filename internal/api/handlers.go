package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourname/ducklake-loader/internal/storage"
)

type Handler struct {
	repo storage.ExerciseRepository
}

func NewHandler(repo storage.ExerciseRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercises)
}

func (h *Handler) GetExerciseByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	exercise, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exercise == nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercise)
}

func (h *Handler) GetExercisesByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	exerciseType := vars["type"]

	exercises, err := h.repo.GetByType(exerciseType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercises)
}

func (h *Handler) GetExercisesByDateRange(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		http.Error(w, "start and end date parameters are required", http.StatusBadRequest)
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		http.Error(w, "Invalid start date format (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		http.Error(w, "Invalid end date format (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	exercises, err := h.repo.GetByDateRange(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercises)
}

func (h *Handler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/exercises", h.GetExercises).Methods("GET")
	r.HandleFunc("/exercises/{id:[0-9]+}", h.GetExerciseByID).Methods("GET")
	r.HandleFunc("/exercises/type/{type}", h.GetExercisesByType).Methods("GET")
	r.HandleFunc("/exercises/date-range", h.GetExercisesByDateRange).Methods("GET")

	return r
}
