package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Yang92047111/ducklake-quick-start/internal/storage"
	"github.com/gorilla/mux"
)

// LakehouseHandler extends the basic Handler with lakehouse-specific endpoints
type LakehouseHandler struct {
	*Handler
	lakehouseRepo storage.LakehouseRepository
}

// NewLakehouseHandler creates a new lakehouse handler
func NewLakehouseHandler(repo storage.LakehouseRepository) *LakehouseHandler {
	// Create base handler by wrapping the lakehouse repo as ExerciseRepository
	baseHandler := NewHandler(repo)

	return &LakehouseHandler{
		Handler:       baseHandler,
		lakehouseRepo: repo,
	}
}

// SetupLakehouseRoutes sets up all lakehouse-specific routes
func (h *LakehouseHandler) SetupLakehouseRoutes() *mux.Router {
	// Start with base routes
	router := h.SetupRoutes()

	// Add lakehouse-specific routes

	// Version and Time Travel endpoints
	router.HandleFunc("/api/v1/versions", h.GetVersionHistory).Methods("GET")
	router.HandleFunc("/api/v1/versions/{version}", h.GetByVersion).Methods("GET")
	router.HandleFunc("/api/v1/versions", h.CreateVersion).Methods("POST")
	router.HandleFunc("/api/v1/time-travel", h.GetByTimestamp).Methods("GET")

	// Schema Management endpoints
	router.HandleFunc("/api/v1/schema", h.GetCurrentSchema).Methods("GET")
	router.HandleFunc("/api/v1/schema", h.EvolveSchema).Methods("PUT")
	router.HandleFunc("/api/v1/schema/history", h.GetSchemaHistory).Methods("GET")
	router.HandleFunc("/api/v1/schema/validate", h.ValidateSchemaCompatibility).Methods("POST")

	// Metadata and Catalog endpoints
	router.HandleFunc("/api/v1/metadata", h.GetTableMetadata).Methods("GET")
	router.HandleFunc("/api/v1/metadata/properties", h.UpdateTableProperties).Methods("PUT")
	router.HandleFunc("/api/v1/partitions", h.GetPartitions).Methods("GET")

	// Transaction Management endpoints
	router.HandleFunc("/api/v1/transactions", h.BeginTransaction).Methods("POST")
	router.HandleFunc("/api/v1/transactions/{txId}/commit", h.CommitTransaction).Methods("POST")
	router.HandleFunc("/api/v1/transactions/{txId}/rollback", h.RollbackTransaction).Methods("POST")
	router.HandleFunc("/api/v1/transactions/{txId}/status", h.GetTransactionStatus).Methods("GET")

	// Optimization and Maintenance endpoints
	router.HandleFunc("/api/v1/optimize", h.OptimizeTable).Methods("POST")
	router.HandleFunc("/api/v1/compact", h.CompactTable).Methods("POST")
	router.HandleFunc("/api/v1/vacuum", h.VacuumTable).Methods("POST")

	// Data Quality endpoints
	router.HandleFunc("/api/v1/constraints", h.GetConstraints).Methods("GET")
	router.HandleFunc("/api/v1/constraints", h.AddConstraint).Methods("POST")
	router.HandleFunc("/api/v1/constraints/{name}", h.RemoveConstraint).Methods("DELETE")
	router.HandleFunc("/api/v1/data-quality", h.GetDataQualityMetrics).Methods("GET")

	// Streaming and Change Data Capture endpoints
	router.HandleFunc("/api/v1/changes", h.GetChangelog).Methods("GET")
	router.HandleFunc("/api/v1/changes/stream", h.StreamChanges).Methods("GET")

	// Advanced Query endpoints
	router.HandleFunc("/api/v1/query/sql", h.QueryWithSQL).Methods("POST")
	router.HandleFunc("/api/v1/query/filter", h.QueryWithFilter).Methods("POST")
	router.HandleFunc("/api/v1/query/aggregate", h.AggregateByTimeWindow).Methods("POST")

	// Performance and Statistics endpoints
	router.HandleFunc("/api/v1/stats/query", h.GetQueryStats).Methods("GET")
	router.HandleFunc("/api/v1/indexes", h.GetIndexes).Methods("GET")
	router.HandleFunc("/api/v1/indexes", h.CreateIndex).Methods("POST")
	router.HandleFunc("/api/v1/indexes/{name}", h.DropIndex).Methods("DELETE")

	return router
}

// Version and Time Travel Handlers

func (h *LakehouseHandler) GetVersionHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	versions, err := h.lakehouseRepo.GetVersionHistory(ctx)
	if err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to get version history: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"versions": versions,
		"count":    len(versions),
	})
}

func (h *LakehouseHandler) GetByVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	version, err := strconv.ParseInt(vars["version"], 10, 64)
	if err != nil {
		h.writeJSONError(w, "Invalid version number", http.StatusBadRequest)
		return
	}

	exercises, err := h.lakehouseRepo.GetByVersion(ctx, version)
	if err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to get data for version %d: %v", version, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"version":   version,
		"exercises": exercises,
		"count":     len(exercises),
	})
}

func (h *LakehouseHandler) CreateVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	version, err := h.lakehouseRepo.CreateVersion(ctx, req.Description)
	if err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to create version: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(version)
}

func (h *LakehouseHandler) GetByTimestamp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	timestampStr := r.URL.Query().Get("timestamp")
	if timestampStr == "" {
		h.writeJSONError(w, "timestamp parameter is required", http.StatusBadRequest)
		return
	}

	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		h.writeJSONError(w, "Invalid timestamp format. Use RFC3339 format", http.StatusBadRequest)
		return
	}

	exercises, err := h.lakehouseRepo.GetByTimestamp(ctx, timestamp)
	if err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to get data for timestamp %s: %v", timestampStr, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"timestamp": timestamp,
		"exercises": exercises,
		"count":     len(exercises),
	})
}

// Schema Management Handlers

func (h *LakehouseHandler) GetCurrentSchema(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	schema, err := h.lakehouseRepo.GetCurrentSchema(ctx)
	if err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to get current schema: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schema)
}

func (h *LakehouseHandler) EvolveSchema(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newSchema storage.Schema
	if err := json.NewDecoder(r.Body).Decode(&newSchema); err != nil {
		h.writeJSONError(w, "Invalid schema format", http.StatusBadRequest)
		return
	}

	if err := h.lakehouseRepo.EvolveSchema(ctx, &newSchema); err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to evolve schema: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Schema evolved successfully",
		"schema":  newSchema,
	})
}

func (h *LakehouseHandler) GetSchemaHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	schemas, err := h.lakehouseRepo.GetSchemaHistory(ctx)
	if err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to get schema history: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"schemas": schemas,
		"count":   len(schemas),
	})
}

func (h *LakehouseHandler) ValidateSchemaCompatibility(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var schema storage.Schema
	if err := json.NewDecoder(r.Body).Decode(&schema); err != nil {
		h.writeJSONError(w, "Invalid schema format", http.StatusBadRequest)
		return
	}

	err := h.lakehouseRepo.ValidateSchemaCompatibility(ctx, &schema)

	response := map[string]interface{}{
		"compatible": err == nil,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Metadata and Catalog Handlers

func (h *LakehouseHandler) GetTableMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metadata, err := h.lakehouseRepo.GetTableMetadata(ctx)
	if err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to get table metadata: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}

func (h *LakehouseHandler) UpdateTableProperties(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var properties map[string]string
	if err := json.NewDecoder(r.Body).Decode(&properties); err != nil {
		h.writeJSONError(w, "Invalid properties format", http.StatusBadRequest)
		return
	}

	if err := h.lakehouseRepo.UpdateTableProperties(ctx, properties); err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to update table properties: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Table properties updated successfully",
		"properties": properties,
	})
}

func (h *LakehouseHandler) GetPartitions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	partitions, err := h.lakehouseRepo.GetPartitions(ctx)
	if err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to get partitions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"partitions": partitions,
		"count":      len(partitions),
	})
}

// Transaction Management Handlers

func (h *LakehouseHandler) BeginTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tx, err := h.lakehouseRepo.BeginTransaction(ctx)
	if err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to begin transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transaction_id":  tx.ID(),
		"start_time":      tx.StartTime(),
		"isolation_level": tx.IsolationLevel(),
		"status":          "active",
	})
}

func (h *LakehouseHandler) CommitTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txId := vars["txId"]

	// In a real implementation, you'd retrieve the transaction by ID
	// For now, this is a simplified implementation
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transaction_id": txId,
		"message":        "Transaction management not fully implemented in this demo",
		"status":         "not_implemented",
	})
}

func (h *LakehouseHandler) RollbackTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txId := vars["txId"]

	// In a real implementation, you'd retrieve and rollback the transaction
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transaction_id": txId,
		"message":        "Transaction management not fully implemented in this demo",
		"status":         "not_implemented",
	})
}

func (h *LakehouseHandler) GetTransactionStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txId := vars["txId"]

	// In a real implementation, you'd look up the transaction status
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transaction_id": txId,
		"status":         "unknown",
		"message":        "Transaction status lookup not implemented in this demo",
	})
}

// Optimization and Maintenance Handlers

func (h *LakehouseHandler) OptimizeTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var options storage.OptimizeOptions
	if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
		// Use default options if body is empty or invalid
		options = storage.OptimizeOptions{
			CompactSmallFiles: true,
			RewriteLargeFiles: false,
		}
	}

	result, err := h.lakehouseRepo.OptimizeTable(ctx, options)
	if err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to optimize table: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *LakehouseHandler) CompactTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := h.lakehouseRepo.Compact(ctx)
	if err != nil {
		h.writeJSONError(w, fmt.Sprintf("Failed to compact table: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *LakehouseHandler) VacuumTable(w http.ResponseWriter, r *http.Request) {
	// Vacuum implementation would clean up old files
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Vacuum operation not implemented in this demo",
		"status":  "not_implemented",
	})
}

// Placeholder implementations for remaining handlers
// These would be fully implemented in a production system

func (h *LakehouseHandler) GetConstraints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"constraints": []interface{}{},
		"message":     "Constraints feature not implemented in this demo",
	})
}

func (h *LakehouseHandler) AddConstraint(w http.ResponseWriter, r *http.Request) {
	h.writeJSONError(w, "Constraints feature not implemented in this demo", http.StatusNotImplemented)
}

func (h *LakehouseHandler) RemoveConstraint(w http.ResponseWriter, r *http.Request) {
	h.writeJSONError(w, "Constraints feature not implemented in this demo", http.StatusNotImplemented)
}

func (h *LakehouseHandler) GetDataQualityMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Data quality metrics not implemented in this demo",
	})
}

func (h *LakehouseHandler) GetChangelog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"changes": []interface{}{},
		"message": "Changelog feature not implemented in this demo",
	})
}

func (h *LakehouseHandler) StreamChanges(w http.ResponseWriter, r *http.Request) {
	h.writeJSONError(w, "Change streaming not implemented in this demo", http.StatusNotImplemented)
}

func (h *LakehouseHandler) QueryWithSQL(w http.ResponseWriter, r *http.Request) {
	h.writeJSONError(w, "SQL query feature not implemented in this demo", http.StatusNotImplemented)
}

func (h *LakehouseHandler) QueryWithFilter(w http.ResponseWriter, r *http.Request) {
	h.writeJSONError(w, "Advanced filtering not implemented in this demo", http.StatusNotImplemented)
}

func (h *LakehouseHandler) AggregateByTimeWindow(w http.ResponseWriter, r *http.Request) {
	h.writeJSONError(w, "Time window aggregation not implemented in this demo", http.StatusNotImplemented)
}

func (h *LakehouseHandler) GetQueryStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Query statistics not implemented in this demo",
	})
}

func (h *LakehouseHandler) GetIndexes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"indexes": []interface{}{},
		"message": "Index management not implemented in this demo",
	})
}

func (h *LakehouseHandler) CreateIndex(w http.ResponseWriter, r *http.Request) {
	h.writeJSONError(w, "Index creation not implemented in this demo", http.StatusNotImplemented)
}

func (h *LakehouseHandler) DropIndex(w http.ResponseWriter, r *http.Request) {
	h.writeJSONError(w, "Index deletion not implemented in this demo", http.StatusNotImplemented)
}
