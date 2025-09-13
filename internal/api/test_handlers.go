package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Simple test endpoints for batch and streaming functionality

// HandleBatchTest handles basic batch testing
func (h *Handler) HandleBatchTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"message":  "Batch endpoint is working",
		"version":  "1.0",
		"features": []string{"batch_insert", "batch_update", "bulk_load"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleStreamTest handles basic streaming testing
func (h *Handler) HandleStreamTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"message":  "Streaming endpoint is working",
		"version":  "1.0",
		"features": []string{"stream_create", "stream_publish", "stream_subscribe"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RegisterTestRoutes registers simple test routes
func (h *Handler) RegisterTestRoutes(router *mux.Router) {
	// Simple test endpoints
	router.HandleFunc("/api/v1/lakehouse/batch/test", h.HandleBatchTest).Methods("POST")
	router.HandleFunc("/api/v1/lakehouse/streams/test", h.HandleStreamTest).Methods("GET")
}
