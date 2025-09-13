package storage

import (
	"context"
	"time"

	"github.com/Yang92047111/ducklake-quick-start/internal/loader"
)

// LakehouseRepository extends ExerciseRepository with lakehouse-specific features
// It provides ACID transactions, schema evolution, time travel, and metadata management
type LakehouseRepository interface {
	ExerciseRepository

	// Transaction Management
	BeginTransaction(ctx context.Context) (Transaction, error)
	CommitTransaction(ctx context.Context, tx Transaction) error
	RollbackTransaction(ctx context.Context, tx Transaction) error

	// Versioning and Time Travel
	GetByVersion(ctx context.Context, version int64) ([]loader.Exercise, error)
	GetByTimestamp(ctx context.Context, timestamp time.Time) ([]loader.Exercise, error)
	GetVersionHistory(ctx context.Context) ([]Version, error)
	CreateVersion(ctx context.Context, description string) (*Version, error)

	// Schema Evolution
	GetCurrentSchema(ctx context.Context) (*Schema, error)
	EvolveSchema(ctx context.Context, newSchema *Schema) error
	GetSchemaHistory(ctx context.Context) ([]Schema, error)
	ValidateSchemaCompatibility(ctx context.Context, newSchema *Schema) error

	// Metadata and Catalog
	GetTableMetadata(ctx context.Context) (*TableMetadata, error)
	UpdateTableProperties(ctx context.Context, properties map[string]string) error
	GetPartitions(ctx context.Context) ([]Partition, error)
	OptimizeTable(ctx context.Context, options OptimizeOptions) (*OptimizeResult, error)

	// Data Quality and Constraints
	AddConstraint(ctx context.Context, constraint Constraint) error
	RemoveConstraint(ctx context.Context, constraintName string) error
	ValidateConstraints(ctx context.Context, exercises []loader.Exercise) error
	GetDataQualityMetrics(ctx context.Context) (*DataQualityMetrics, error)

	// Streaming and Change Data Capture
	WatchChanges(ctx context.Context, from time.Time) (<-chan ChangeEvent, error)
	GetChangelog(ctx context.Context, fromVersion, toVersion int64) ([]ChangeEvent, error)

	// Performance and Optimization
	CreateIndex(ctx context.Context, indexName string, columns []string) error
	DropIndex(ctx context.Context, indexName string) error
	GetQueryStats(ctx context.Context) (*QueryStats, error)
	Compact(ctx context.Context) (*CompactionResult, error)

	// Advanced Querying
	QueryWithSQL(ctx context.Context, sql string, params ...interface{}) ([]loader.Exercise, error)
	QueryWithFilter(ctx context.Context, filter Filter) ([]loader.Exercise, error)
	AggregateByTimeWindow(ctx context.Context, window TimeWindow, aggregations []Aggregation) ([]AggregationResult, error)

	// Batch Processing
	InsertBatchWithOptions(ctx context.Context, exercises []loader.Exercise, options BatchOptions) (*BatchResult, error)
	UpdateBatch(ctx context.Context, exercises []loader.Exercise) (*BatchResult, error)
	DeleteBatch(ctx context.Context, ids []int) (*BatchResult, error)
	BulkLoad(ctx context.Context, dataSource DataSource, options BulkLoadOptions) (*BulkLoadResult, error)

	// Streaming Support
	StartStream(ctx context.Context, config StreamConfig) (Stream, error)
	PublishToStream(ctx context.Context, streamName string, exercises []loader.Exercise) error
	SubscribeToStream(ctx context.Context, streamName string) (<-chan StreamEvent, error)
	GetActiveStreams(ctx context.Context) ([]StreamInfo, error)
}

// Transaction represents a lakehouse transaction
type Transaction interface {
	ID() string
	StartTime() time.Time
	IsolationLevel() IsolationLevel
	GetOperations() []Operation

	// Transaction-specific operations
	Insert(exercise loader.Exercise) error
	InsertBatch(exercises []loader.Exercise) error
	Update(exercise loader.Exercise) error
	Delete(id int) error

	// Transaction state
	IsActive() bool
	GetConflicts() []Conflict
}

// Version represents a table version with metadata
type Version struct {
	ID          int64             `json:"id"`
	Timestamp   time.Time         `json:"timestamp"`
	Description string            `json:"description,omitempty"`
	SchemaID    int64             `json:"schema_id"`
	RecordCount int64             `json:"record_count"`
	FileCount   int               `json:"file_count"`
	SizeBytes   int64             `json:"size_bytes"`
	Properties  map[string]string `json:"properties,omitempty"`
	Operations  []Operation       `json:"operations"`
	ParentID    *int64            `json:"parent_id,omitempty"`
}

// Schema represents the table schema with evolution capabilities
type Schema struct {
	ID          int64             `json:"id"`
	Version     int               `json:"version"`
	Fields      []Field           `json:"fields"`
	CreatedAt   time.Time         `json:"created_at"`
	Properties  map[string]string `json:"properties,omitempty"`
	Constraints []Constraint      `json:"constraints,omitempty"`
}

// Field represents a column in the schema
type Field struct {
	Name         string            `json:"name"`
	Type         FieldType         `json:"type"`
	Nullable     bool              `json:"nullable"`
	DefaultValue interface{}       `json:"default_value,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// FieldType represents column data types
type FieldType string

const (
	FieldTypeInt       FieldType = "int"
	FieldTypeString    FieldType = "string"
	FieldTypeFloat     FieldType = "float"
	FieldTypeBoolean   FieldType = "boolean"
	FieldTypeTimestamp FieldType = "timestamp"
	FieldTypeDate      FieldType = "date"
	FieldTypeArray     FieldType = "array"
	FieldTypeMap       FieldType = "map"
	FieldTypeStruct    FieldType = "struct"
)

// TableMetadata contains comprehensive table information
type TableMetadata struct {
	Name            string            `json:"name"`
	Location        string            `json:"location"`
	Format          string            `json:"format"`
	CreatedAt       time.Time         `json:"created_at"`
	LastModified    time.Time         `json:"last_modified"`
	CurrentVersion  int64             `json:"current_version"`
	RecordCount     int64             `json:"record_count"`
	FileCount       int               `json:"file_count"`
	SizeBytes       int64             `json:"size_bytes"`
	Properties      map[string]string `json:"properties"`
	PartitionFields []string          `json:"partition_fields,omitempty"`
	SortFields      []string          `json:"sort_fields,omitempty"`
	Owner           string            `json:"owner,omitempty"`
	Description     string            `json:"description,omitempty"`
}

// Partition represents a table partition
type Partition struct {
	Values       map[string]interface{} `json:"values"`
	RecordCount  int64                  `json:"record_count"`
	FileCount    int                    `json:"file_count"`
	SizeBytes    int64                  `json:"size_bytes"`
	Location     string                 `json:"location"`
	LastModified time.Time              `json:"last_modified"`
}

// Constraint represents data quality constraints
type Constraint struct {
	Name        string         `json:"name"`
	Type        ConstraintType `json:"type"`
	Expression  string         `json:"expression"`
	Columns     []string       `json:"columns,omitempty"`
	Enabled     bool           `json:"enabled"`
	CreatedAt   time.Time      `json:"created_at"`
	Description string         `json:"description,omitempty"`
}

// ConstraintType defines types of constraints
type ConstraintType string

const (
	ConstraintTypeNotNull    ConstraintType = "not_null"
	ConstraintTypeUnique     ConstraintType = "unique"
	ConstraintTypePrimaryKey ConstraintType = "primary_key"
	ConstraintTypeForeignKey ConstraintType = "foreign_key"
	ConstraintTypeCheck      ConstraintType = "check"
	ConstraintTypeRange      ConstraintType = "range"
	ConstraintTypeRegex      ConstraintType = "regex"
)

// DataQualityMetrics provides insights into data quality
type DataQualityMetrics struct {
	TotalRecords         int64            `json:"total_records"`
	NullValues           map[string]int64 `json:"null_values"`
	UniqueValues         map[string]int64 `json:"unique_values"`
	DuplicateRecords     int64            `json:"duplicate_records"`
	ConstraintViolations map[string]int64 `json:"constraint_violations"`
	DataTypeMismatches   map[string]int64 `json:"data_type_mismatches"`
	OutlierCounts        map[string]int64 `json:"outlier_counts"`
	CompletenessScore    float64          `json:"completeness_score"`
	ValidityScore        float64          `json:"validity_score"`
	ConsistencyScore     float64          `json:"consistency_score"`
	LastUpdated          time.Time        `json:"last_updated"`
}

// ChangeEvent represents a change to the table
type ChangeEvent struct {
	Type      ChangeType             `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Version   int64                  `json:"version"`
	Operation Operation              `json:"operation"`
	RecordID  *int                   `json:"record_id,omitempty"`
	Before    *loader.Exercise       `json:"before,omitempty"`
	After     *loader.Exercise       `json:"after,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ChangeType defines types of changes
type ChangeType string

const (
	ChangeTypeInsert   ChangeType = "insert"
	ChangeTypeUpdate   ChangeType = "update"
	ChangeTypeDelete   ChangeType = "delete"
	ChangeTypeSchema   ChangeType = "schema"
	ChangeTypeOptimize ChangeType = "optimize"
)

// Operation represents a single operation in a transaction
type Operation struct {
	Type           OperationType          `json:"type"`
	Timestamp      time.Time              `json:"timestamp"`
	User           string                 `json:"user,omitempty"`
	Details        map[string]interface{} `json:"details"`
	RecordsRead    int64                  `json:"records_read"`
	RecordsWritten int64                  `json:"records_written"`
	Duration       time.Duration          `json:"duration"`
}

// OperationType defines types of operations
type OperationType string

const (
	OperationTypeWrite    OperationType = "write"
	OperationTypeDelete   OperationType = "delete"
	OperationTypeOptimize OperationType = "optimize"
	OperationTypeRestore  OperationType = "restore"
	OperationTypeClone    OperationType = "clone"
	OperationTypeVacuum   OperationType = "vacuum"
	OperationTypeSchema   OperationType = "schema"
)

// IsolationLevel defines transaction isolation levels
type IsolationLevel string

const (
	IsolationReadUncommitted IsolationLevel = "read_uncommitted"
	IsolationReadCommitted   IsolationLevel = "read_committed"
	IsolationRepeatableRead  IsolationLevel = "repeatable_read"
	IsolationSerializable    IsolationLevel = "serializable"
)

// Conflict represents a transaction conflict
type Conflict struct {
	Type          ConflictType `json:"type"`
	ResourceID    string       `json:"resource_id"`
	ConflictingTx string       `json:"conflicting_tx"`
	Description   string       `json:"description"`
	Timestamp     time.Time    `json:"timestamp"`
}

// ConflictType defines types of conflicts
type ConflictType string

const (
	ConflictTypeWrite  ConflictType = "write_write"
	ConflictTypeRead   ConflictType = "read_write"
	ConflictTypeSchema ConflictType = "schema"
)

// OptimizeOptions configures table optimization
type OptimizeOptions struct {
	MaxFileSize       int64    `json:"max_file_size,omitempty"`
	MinFileSize       int64    `json:"min_file_size,omitempty"`
	TargetFileCount   int      `json:"target_file_count,omitempty"`
	PartitionFilters  []string `json:"partition_filters,omitempty"`
	ZOrderColumns     []string `json:"z_order_columns,omitempty"`
	CompactSmallFiles bool     `json:"compact_small_files"`
	RewriteLargeFiles bool     `json:"rewrite_large_files"`
}

// OptimizeResult contains optimization results
type OptimizeResult struct {
	FilesAdded          int                    `json:"files_added"`
	FilesRemoved        int                    `json:"files_removed"`
	PartitionsOptimized int                    `json:"partitions_optimized"`
	RecordsRewritten    int64                  `json:"records_rewritten"`
	BytesWritten        int64                  `json:"bytes_written"`
	BytesRemoved        int64                  `json:"bytes_removed"`
	Duration            time.Duration          `json:"duration"`
	Metrics             map[string]interface{} `json:"metrics"`
}

// CompactionResult contains compaction results
type CompactionResult struct {
	FilesCompacted   int           `json:"files_compacted"`
	FilesCreated     int           `json:"files_created"`
	RecordsProcessed int64         `json:"records_processed"`
	SpaceReclaimed   int64         `json:"space_reclaimed"`
	Duration         time.Duration `json:"duration"`
}

// QueryStats provides query performance metrics
type QueryStats struct {
	TotalQueries     int64            `json:"total_queries"`
	AverageLatency   time.Duration    `json:"average_latency"`
	P95Latency       time.Duration    `json:"p95_latency"`
	CacheHitRate     float64          `json:"cache_hit_rate"`
	IndexUsage       map[string]int64 `json:"index_usage"`
	PartitionsPruned int64            `json:"partitions_pruned"`
	RecordsScanned   int64            `json:"records_scanned"`
	RecordsReturned  int64            `json:"records_returned"`
	TopQueries       []QueryInfo      `json:"top_queries"`
	LastUpdated      time.Time        `json:"last_updated"`
}

// QueryInfo contains individual query information
type QueryInfo struct {
	SQL           string        `json:"sql"`
	ExecutionTime time.Duration `json:"execution_time"`
	RecordsRead   int64         `json:"records_read"`
	Timestamp     time.Time     `json:"timestamp"`
}

// Filter represents query filtering options
type Filter struct {
	Conditions []Condition `json:"conditions"`
	SortBy     []SortField `json:"sort_by,omitempty"`
	Limit      *int        `json:"limit,omitempty"`
	Offset     *int        `json:"offset,omitempty"`
	GroupBy    []string    `json:"group_by,omitempty"`
	Having     []Condition `json:"having,omitempty"`
}

// Condition represents a filter condition
type Condition struct {
	Field    string      `json:"field"`
	Operator Operator    `json:"operator"`
	Value    interface{} `json:"value"`
}

// Operator defines filter operators
type Operator string

const (
	OperatorEqual              Operator = "eq"
	OperatorNotEqual           Operator = "ne"
	OperatorGreaterThan        Operator = "gt"
	OperatorGreaterThanOrEqual Operator = "gte"
	OperatorLessThan           Operator = "lt"
	OperatorLessThanOrEqual    Operator = "lte"
	OperatorIn                 Operator = "in"
	OperatorNotIn              Operator = "not_in"
	OperatorLike               Operator = "like"
	OperatorNotLike            Operator = "not_like"
	OperatorIsNull             Operator = "is_null"
	OperatorIsNotNull          Operator = "is_not_null"
	OperatorBetween            Operator = "between"
)

// SortField represents sorting configuration
type SortField struct {
	Field string    `json:"field"`
	Order SortOrder `json:"order"`
}

// SortOrder defines sort directions
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// TimeWindow represents time-based aggregation windows
type TimeWindow struct {
	Size  time.Duration `json:"size"`
	Slide time.Duration `json:"slide,omitempty"`
	Start time.Time     `json:"start,omitempty"`
	End   time.Time     `json:"end,omitempty"`
}

// Aggregation represents aggregation functions
type Aggregation struct {
	Function AggregationFunction `json:"function"`
	Field    string              `json:"field"`
	Alias    string              `json:"alias,omitempty"`
}

// AggregationFunction defines aggregation functions
type AggregationFunction string

const (
	AggregationCount AggregationFunction = "count"
	AggregationSum   AggregationFunction = "sum"
	AggregationAvg   AggregationFunction = "avg"
	AggregationMin   AggregationFunction = "min"
	AggregationMax   AggregationFunction = "max"
)

// AggregationResult represents the result of an aggregation
type AggregationResult struct {
	Window time.Time              `json:"window"`
	Values map[string]interface{} `json:"values"`
}

// BatchOptions configures batch operations
type BatchOptions struct {
	BatchSize     int           `json:"batch_size"`
	Timeout       time.Duration `json:"timeout"`
	ParallelJobs  int           `json:"parallel_jobs"`
	SkipErrors    bool          `json:"skip_errors"`
	ValidateFirst bool          `json:"validate_first"`
}

// BatchResult contains the result of batch operations
type BatchResult struct {
	ProcessedCount int           `json:"processed_count"`
	SuccessCount   int           `json:"success_count"`
	ErrorCount     int           `json:"error_count"`
	Duration       time.Duration `json:"duration"`
	Errors         []BatchError  `json:"errors,omitempty"`
	Version        int64         `json:"version,omitempty"`
}

// BatchError represents an error in batch processing
type BatchError struct {
	Index  int    `json:"index"`
	Error  string `json:"error"`
	Record string `json:"record,omitempty"`
}

// DataSource represents a data source for bulk loading
type DataSource struct {
	Type       DataSourceType    `json:"type"`
	Location   string            `json:"location"`
	Format     DataFormat        `json:"format"`
	Schema     *Schema           `json:"schema,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

// DataSourceType defines data source types
type DataSourceType string

const (
	DataSourceFile   DataSourceType = "file"
	DataSourceS3     DataSourceType = "s3"
	DataSourceGCS    DataSourceType = "gcs"
	DataSourceHTTP   DataSourceType = "http"
	DataSourceStream DataSourceType = "stream"
)

// DataFormat defines data formats
type DataFormat string

const (
	DataFormatJSON    DataFormat = "json"
	DataFormatCSV     DataFormat = "csv"
	DataFormatParquet DataFormat = "parquet"
	DataFormatAvro    DataFormat = "avro"
)

// BulkLoadOptions configures bulk loading
type BulkLoadOptions struct {
	BatchSize      int           `json:"batch_size"`
	ParallelJobs   int           `json:"parallel_jobs"`
	Timeout        time.Duration `json:"timeout"`
	SkipHeader     bool          `json:"skip_header"`
	SkipErrors     bool          `json:"skip_errors"`
	ValidateSchema bool          `json:"validate_schema"`
	Compression    string        `json:"compression,omitempty"`
}

// BulkLoadResult contains the result of bulk loading
type BulkLoadResult struct {
	RecordsLoaded  int64         `json:"records_loaded"`
	RecordsSkipped int64         `json:"records_skipped"`
	RecordsErrored int64         `json:"records_errored"`
	Duration       time.Duration `json:"duration"`
	Version        int64         `json:"version"`
	Errors         []BulkError   `json:"errors,omitempty"`
}

// BulkError represents an error in bulk loading
type BulkError struct {
	Line    int64  `json:"line"`
	Column  string `json:"column,omitempty"`
	Error   string `json:"error"`
	Content string `json:"content,omitempty"`
}

// StreamConfig configures streaming
type StreamConfig struct {
	Name          string            `json:"name"`
	Type          StreamType        `json:"type"`
	BufferSize    int               `json:"buffer_size"`
	FlushInterval time.Duration     `json:"flush_interval"`
	Partitions    int               `json:"partitions"`
	Properties    map[string]string `json:"properties,omitempty"`
}

// StreamType defines stream types
type StreamType string

const (
	StreamTypeChangeData StreamType = "change_data"
	StreamTypeInserts    StreamType = "inserts"
	StreamTypeUpdates    StreamType = "updates"
	StreamTypeDeletes    StreamType = "deletes"
	StreamTypeAll        StreamType = "all"
)

// Stream represents an active stream
type Stream interface {
	Name() string
	Type() StreamType
	Config() StreamConfig
	IsActive() bool
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Publish(ctx context.Context, events []StreamEvent) error
	Subscribe(ctx context.Context) (<-chan StreamEvent, error)
	GetStats() StreamStats
}

// StreamEvent represents a streaming event
type StreamEvent struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Type      StreamEventType        `json:"type"`
	Version   int64                  `json:"version"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]string      `json:"metadata,omitempty"`
}

// StreamEventType defines stream event types
type StreamEventType string

const (
	StreamEventInsert StreamEventType = "insert"
	StreamEventUpdate StreamEventType = "update"
	StreamEventDelete StreamEventType = "delete"
	StreamEventSchema StreamEventType = "schema"
)

// StreamInfo provides information about a stream
type StreamInfo struct {
	Name         string      `json:"name"`
	Type         StreamType  `json:"type"`
	IsActive     bool        `json:"is_active"`
	CreatedAt    time.Time   `json:"created_at"`
	LastActivity time.Time   `json:"last_activity"`
	Stats        StreamStats `json:"stats"`
}

// StreamStats provides streaming statistics
type StreamStats struct {
	EventsPublished   int64     `json:"events_published"`
	EventsConsumed    int64     `json:"events_consumed"`
	ActiveSubscribers int       `json:"active_subscribers"`
	LastEventTime     time.Time `json:"last_event_time"`
	BytesTransferred  int64     `json:"bytes_transferred"`
}
