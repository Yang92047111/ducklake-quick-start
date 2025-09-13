package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Yang92047111/ducklake-quick-start/internal/loader"
	"github.com/Yang92047111/ducklake-quick-start/internal/storage"
	"github.com/gorilla/mux"
)

// Batch Processing Endpoints

// BatchInsertRequest represents a batch insert request
type BatchInsertRequest struct {
	Exercises []loader.Exercise    `json:"exercises"`
	Options   storage.BatchOptions `json:"options,omitempty"`
}

// BulkLoadRequest represents a bulk load request
type BulkLoadRequest struct {
	DataSource storage.DataSource      `json:"data_source"`
	Options    storage.BulkLoadOptions `json:"options,omitempty"`
}

// StreamCreateRequest represents a stream creation request
type StreamCreateRequest struct {
	Config storage.StreamConfig `json:"config"`
}

// StreamPublishRequest represents a stream publish request
type StreamPublishRequest struct {
	Exercises []loader.Exercise `json:"exercises"`
}

// HandleBatchInsert handles batch insertion of exercises
func (h *Handler) HandleBatchInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BatchInsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if len(req.Exercises) == 0 {
		http.Error(w, "No exercises provided", http.StatusBadRequest)
		return
	}

	// Get lakehouse repository
	lakeRepo, ok := h.repo.(storage.LakehouseRepository)
	if !ok {
		http.Error(w, "Lakehouse operations not supported", http.StatusNotImplemented)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Minute)
	defer cancel()

	result, err := lakeRepo.InsertBatchWithOptions(ctx, req.Exercises, req.Options)
	if err != nil {
		http.Error(w, fmt.Sprintf("Batch insert failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// HandleBatchUpdate handles batch updates of exercises
func (h *Handler) HandleBatchUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var exercises []loader.Exercise
	if err := json.NewDecoder(r.Body).Decode(&exercises); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if len(exercises) == 0 {
		http.Error(w, "No exercises provided", http.StatusBadRequest)
		return
	}

	lakeRepo, ok := h.repo.(storage.LakehouseRepository)
	if !ok {
		http.Error(w, "Lakehouse operations not supported", http.StatusNotImplemented)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Minute)
	defer cancel()

	result, err := lakeRepo.UpdateBatch(ctx, exercises)
	if err != nil {
		http.Error(w, fmt.Sprintf("Batch update failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// HandleBatchDelete handles batch deletion of exercises
func (h *Handler) HandleBatchDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var ids []int
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if len(ids) == 0 {
		http.Error(w, "No IDs provided", http.StatusBadRequest)
		return
	}

	lakeRepo, ok := h.repo.(storage.LakehouseRepository)
	if !ok {
		http.Error(w, "Lakehouse operations not supported", http.StatusNotImplemented)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Minute)
	defer cancel()

	result, err := lakeRepo.DeleteBatch(ctx, ids)
	if err != nil {
		http.Error(w, fmt.Sprintf("Batch delete failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// HandleBulkLoad handles bulk loading from various data sources
func (h *Handler) HandleBulkLoad(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BulkLoadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if req.DataSource.Location == "" {
		http.Error(w, "Data source location is required", http.StatusBadRequest)
		return
	}

	lakeRepo, ok := h.repo.(storage.LakehouseRepository)
	if !ok {
		http.Error(w, "Lakehouse operations not supported", http.StatusNotImplemented)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Hour)
	defer cancel()

	result, err := lakeRepo.BulkLoad(ctx, req.DataSource, req.Options)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bulk load failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// Streaming Endpoints

// HandleStreamCreate creates a new stream
func (h *Handler) HandleStreamCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StreamCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if req.Config.Name == "" {
		http.Error(w, "Stream name is required", http.StatusBadRequest)
		return
	}

	lakeRepo, ok := h.repo.(storage.LakehouseRepository)
	if !ok {
		http.Error(w, "Lakehouse operations not supported", http.StatusNotImplemented)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	stream, err := lakeRepo.StartStream(ctx, req.Config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create stream: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"name":      stream.Name(),
		"type":      stream.Type(),
		"is_active": stream.IsActive(),
		"config":    stream.Config(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// HandleStreamPublish publishes data to a stream
func (h *Handler) HandleStreamPublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	streamName := vars["streamName"]
	if streamName == "" {
		http.Error(w, "Stream name is required", http.StatusBadRequest)
		return
	}

	var req StreamPublishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if len(req.Exercises) == 0 {
		http.Error(w, "No exercises provided", http.StatusBadRequest)
		return
	}

	lakeRepo, ok := h.repo.(storage.LakehouseRepository)
	if !ok {
		http.Error(w, "Lakehouse operations not supported", http.StatusNotImplemented)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	err := lakeRepo.PublishToStream(ctx, streamName, req.Exercises)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to publish to stream: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"stream":           streamName,
		"events_published": len(req.Exercises),
		"timestamp":        time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleStreamSubscribe creates a WebSocket connection for stream subscription
func (h *Handler) HandleStreamSubscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamName := vars["streamName"]
	if streamName == "" {
		http.Error(w, "Stream name is required", http.StatusBadRequest)
		return
	}

	lakeRepo, ok := h.repo.(storage.LakehouseRepository)
	if !ok {
		http.Error(w, "Lakehouse operations not supported", http.StatusNotImplemented)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	eventChan, err := lakeRepo.SubscribeToStream(ctx, streamName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to subscribe to stream: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Set up Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send initial connection message
	fmt.Fprintf(w, "data: {\"type\":\"connected\",\"stream\":\"%s\",\"timestamp\":\"%s\"}\n\n",
		streamName, time.Now().Format(time.RFC3339))
	flusher.Flush()

	// Stream events
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				fmt.Fprintf(w, "data: {\"type\":\"stream_closed\",\"stream\":\"%s\",\"timestamp\":\"%s\"}\n\n",
					streamName, time.Now().Format(time.RFC3339))
				flusher.Flush()
				return
			}

			eventData, err := json.Marshal(event)
			if err != nil {
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", eventData)
			flusher.Flush()

		case <-ctx.Done():
			return
		}
	}
}

// HandleStreamsStatus returns status of all active streams
func (h *Handler) HandleStreamsStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lakeRepo, ok := h.repo.(storage.LakehouseRepository)
	if !ok {
		http.Error(w, "Lakehouse operations not supported", http.StatusNotImplemented)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	streams, err := lakeRepo.GetActiveStreams(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get streams: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"streams": streams,
		"count":   len(streams),
	})
}

// Register batch and streaming routes
func (h *Handler) RegisterBatchStreamingRoutes(router *mux.Router) {
	// Batch operations
	router.HandleFunc("/api/v1/lakehouse/batch/insert", h.HandleBatchInsert).Methods("POST")
	router.HandleFunc("/api/v1/lakehouse/batch/update", h.HandleBatchUpdate).Methods("PUT")
	router.HandleFunc("/api/v1/lakehouse/batch/delete", h.HandleBatchDelete).Methods("DELETE")
	router.HandleFunc("/api/v1/lakehouse/bulk-load", h.HandleBulkLoad).Methods("POST")

	// Streaming operations
	router.HandleFunc("/api/v1/lakehouse/streams", h.HandleStreamCreate).Methods("POST")
	router.HandleFunc("/api/v1/lakehouse/streams", h.HandleStreamsStatus).Methods("GET")
	router.HandleFunc("/api/v1/lakehouse/streams/{streamName}/publish", h.HandleStreamPublish).Methods("POST")
	router.HandleFunc("/api/v1/lakehouse/streams/{streamName}/subscribe", h.HandleStreamSubscribe).Methods("GET")
}
