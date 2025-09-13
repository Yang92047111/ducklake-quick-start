package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/Yang92047111/ducklake-quick-start/internal/loader"
)

// Batch Processing Implementation for DeltaLakeRepository

// InsertBatchWithOptions inserts a batch of exercises with configurable options
func (d *DeltaLakeRepository) InsertBatchWithOptions(ctx context.Context, exercises []loader.Exercise, options BatchOptions) (*BatchResult, error) {
	startTime := time.Now()
	result := &BatchResult{
		ProcessedCount: len(exercises),
	}

	// Set default options
	if options.BatchSize == 0 {
		options.BatchSize = 1000
	}
	if options.ParallelJobs == 0 {
		options.ParallelJobs = 1
	}
	if options.Timeout == 0 {
		options.Timeout = 30 * time.Minute
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, options.Timeout)
	defer cancel()

	// Validate first if requested
	if options.ValidateFirst {
		if err := d.validateBatch(exercises); err != nil {
			return result, fmt.Errorf("batch validation failed: %w", err)
		}
	}

	// Begin transaction
	tx, err := d.BeginTransaction(timeoutCtx)
	if err != nil {
		return result, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Process in batches
	for i := 0; i < len(exercises); i += options.BatchSize {
		end := i + options.BatchSize
		if end > len(exercises) {
			end = len(exercises)
		}

		batch := exercises[i:end]
		if err := d.processBatch(timeoutCtx, tx, batch, result, options.SkipErrors); err != nil {
			if !options.SkipErrors {
				d.RollbackTransaction(timeoutCtx, tx)
				return result, fmt.Errorf("batch processing failed: %w", err)
			}
		}
	}

	// Commit transaction
	if err := d.CommitTransaction(timeoutCtx, tx); err != nil {
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	result.Duration = time.Since(startTime)
	result.Version = d.currentVersion
	return result, nil
}

// UpdateBatch updates a batch of exercises
func (d *DeltaLakeRepository) UpdateBatch(ctx context.Context, exercises []loader.Exercise) (*BatchResult, error) {
	startTime := time.Now()
	result := &BatchResult{
		ProcessedCount: len(exercises),
	}

	tx, err := d.BeginTransaction(ctx)
	if err != nil {
		return result, fmt.Errorf("failed to begin transaction: %w", err)
	}

	for i, exercise := range exercises {
		if err := tx.Update(exercise); err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, BatchError{
				Index:  i,
				Error:  err.Error(),
				Record: fmt.Sprintf("ID: %d", exercise.ID),
			})
		} else {
			result.SuccessCount++
		}
	}

	if result.ErrorCount > 0 && result.SuccessCount == 0 {
		d.RollbackTransaction(ctx, tx)
		return result, fmt.Errorf("all updates failed")
	}

	if err := d.CommitTransaction(ctx, tx); err != nil {
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	result.Duration = time.Since(startTime)
	result.Version = d.currentVersion
	return result, nil
}

// DeleteBatch deletes a batch of exercises by IDs
func (d *DeltaLakeRepository) DeleteBatch(ctx context.Context, ids []int) (*BatchResult, error) {
	startTime := time.Now()
	result := &BatchResult{
		ProcessedCount: len(ids),
	}

	tx, err := d.BeginTransaction(ctx)
	if err != nil {
		return result, fmt.Errorf("failed to begin transaction: %w", err)
	}

	for i, id := range ids {
		if err := tx.Delete(id); err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, BatchError{
				Index:  i,
				Error:  err.Error(),
				Record: fmt.Sprintf("ID: %d", id),
			})
		} else {
			result.SuccessCount++
		}
	}

	if err := d.CommitTransaction(ctx, tx); err != nil {
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	result.Duration = time.Since(startTime)
	result.Version = d.currentVersion
	return result, nil
}

// BulkLoad loads data from various sources
func (d *DeltaLakeRepository) BulkLoad(ctx context.Context, dataSource DataSource, options BulkLoadOptions) (*BulkLoadResult, error) {
	result := &BulkLoadResult{}

	// Set default options
	if options.BatchSize == 0 {
		options.BatchSize = 5000
	}
	if options.ParallelJobs == 0 {
		options.ParallelJobs = 1
	}
	if options.Timeout == 0 {
		options.Timeout = 1 * time.Hour
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, options.Timeout)
	defer cancel()

	switch dataSource.Type {
	case DataSourceFile:
		return d.bulkLoadFromFile(timeoutCtx, dataSource, options)
	case DataSourceHTTP:
		return d.bulkLoadFromHTTP(timeoutCtx, dataSource, options)
	default:
		return result, fmt.Errorf("unsupported data source type: %s", dataSource.Type)
	}
}

// Helper methods for batch processing

func (d *DeltaLakeRepository) validateBatch(exercises []loader.Exercise) error {
	validator := loader.NewValidator()
	for i, exercise := range exercises {
		if err := validator.Validate(exercise); err != nil {
			return fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}
	return nil
}

func (d *DeltaLakeRepository) processBatch(ctx context.Context, tx Transaction, batch []loader.Exercise, result *BatchResult, skipErrors bool) error {
	for i, exercise := range batch {
		if err := tx.Insert(exercise); err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, BatchError{
				Index:  i,
				Error:  err.Error(),
				Record: fmt.Sprintf("ID: %d, Name: %s", exercise.ID, exercise.Name),
			})

			if !skipErrors {
				return fmt.Errorf("failed to insert exercise at index %d: %w", i, err)
			}
		} else {
			result.SuccessCount++
		}
	}
	return nil
}

func (d *DeltaLakeRepository) bulkLoadFromFile(ctx context.Context, dataSource DataSource, options BulkLoadOptions) (*BulkLoadResult, error) {
	result := &BulkLoadResult{}

	file, err := os.Open(dataSource.Location)
	if err != nil {
		return result, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	switch dataSource.Format {
	case DataFormatJSON:
		return d.bulkLoadJSON(ctx, file, options)
	case DataFormatCSV:
		return d.bulkLoadCSV(ctx, file, options)
	default:
		return result, fmt.Errorf("unsupported format: %s", dataSource.Format)
	}
}

func (d *DeltaLakeRepository) bulkLoadFromHTTP(ctx context.Context, dataSource DataSource, options BulkLoadOptions) (*BulkLoadResult, error) {
	// Implementation for HTTP data sources
	return &BulkLoadResult{}, fmt.Errorf("HTTP data source not implemented yet")
}

func (d *DeltaLakeRepository) bulkLoadJSON(ctx context.Context, reader io.Reader, options BulkLoadOptions) (*BulkLoadResult, error) {
	result := &BulkLoadResult{}

	var exercises []loader.Exercise
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&exercises); err != nil {
		return result, fmt.Errorf("failed to decode JSON: %w", err)
	}

	batchResult, err := d.InsertBatchWithOptions(ctx, exercises, BatchOptions{
		BatchSize:    options.BatchSize,
		ParallelJobs: options.ParallelJobs,
		Timeout:      options.Timeout,
		SkipErrors:   options.SkipErrors,
	})

	if err != nil {
		return result, err
	}

	result.RecordsLoaded = int64(batchResult.SuccessCount)
	result.RecordsErrored = int64(batchResult.ErrorCount)
	result.Duration = batchResult.Duration
	result.Version = batchResult.Version

	return result, nil
}

func (d *DeltaLakeRepository) bulkLoadCSV(ctx context.Context, reader io.Reader, options BulkLoadOptions) (*BulkLoadResult, error) {
	// Use the existing CSV loader
	csvLoader := loader.NewCSVLoader()

	// Create a temporary file from the reader
	tempFile, err := os.CreateTemp("", "bulk_load_*.csv")
	if err != nil {
		return &BulkLoadResult{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, reader); err != nil {
		return &BulkLoadResult{}, fmt.Errorf("failed to copy data: %w", err)
	}

	exercises, err := csvLoader.LoadFromCSV(tempFile.Name())
	if err != nil {
		return &BulkLoadResult{}, fmt.Errorf("failed to load CSV: %w", err)
	}

	batchResult, err := d.InsertBatchWithOptions(ctx, exercises, BatchOptions{
		BatchSize:    options.BatchSize,
		ParallelJobs: options.ParallelJobs,
		Timeout:      options.Timeout,
		SkipErrors:   options.SkipErrors,
	})

	if err != nil {
		return &BulkLoadResult{}, err
	}

	result := &BulkLoadResult{
		RecordsLoaded:  int64(batchResult.SuccessCount),
		RecordsErrored: int64(batchResult.ErrorCount),
		Duration:       batchResult.Duration,
		Version:        batchResult.Version,
	}

	return result, nil
}

// Streaming Implementation

// deltaStream implements the Stream interface
type deltaStream struct {
	name       string
	streamType StreamType
	config     StreamConfig
	isActive   bool
	events     chan StreamEvent
	stats      StreamStats
	mutex      sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

func (s *deltaStream) Name() string {
	return s.name
}

func (s *deltaStream) Type() StreamType {
	return s.streamType
}

func (s *deltaStream) Config() StreamConfig {
	return s.config
}

func (s *deltaStream) IsActive() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.isActive
}

func (s *deltaStream) Start(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isActive {
		return fmt.Errorf("stream %s is already active", s.name)
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	s.isActive = true
	s.events = make(chan StreamEvent, s.config.BufferSize)

	return nil
}

func (s *deltaStream) Stop(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isActive {
		return fmt.Errorf("stream %s is not active", s.name)
	}

	s.cancel()
	close(s.events)
	s.isActive = false

	return nil
}

func (s *deltaStream) Publish(ctx context.Context, events []StreamEvent) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.isActive {
		return fmt.Errorf("stream %s is not active", s.name)
	}

	for _, event := range events {
		select {
		case s.events <- event:
			s.stats.EventsPublished++
			s.stats.LastEventTime = event.Timestamp
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (s *deltaStream) Subscribe(ctx context.Context) (<-chan StreamEvent, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.isActive {
		return nil, fmt.Errorf("stream %s is not active", s.name)
	}

	s.stats.ActiveSubscribers++
	return s.events, nil
}

func (s *deltaStream) GetStats() StreamStats {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.stats
}

// StartStream creates and starts a new stream
func (d *DeltaLakeRepository) StartStream(ctx context.Context, config StreamConfig) (Stream, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check if stream already exists
	if _, exists := d.streams[config.Name]; exists {
		return nil, fmt.Errorf("stream %s already exists", config.Name)
	}

	if config.BufferSize == 0 {
		config.BufferSize = 1000
	}
	if config.FlushInterval == 0 {
		config.FlushInterval = 5 * time.Second
	}

	stream := &deltaStream{
		name:       config.Name,
		streamType: config.Type,
		config:     config,
		stats:      StreamStats{},
	}

	if err := stream.Start(ctx); err != nil {
		return nil, err
	}

	// Store stream in repository
	d.streams[config.Name] = stream
	return stream, nil
}

// PublishToStream publishes exercises to a stream
func (d *DeltaLakeRepository) PublishToStream(ctx context.Context, streamName string, exercises []loader.Exercise) error {
	d.mutex.RLock()
	stream, exists := d.streams[streamName]
	d.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("stream %s not found", streamName)
	}

	// Convert exercises to stream events
	events := make([]StreamEvent, len(exercises))
	for i, exercise := range exercises {
		data := map[string]interface{}{
			"id":          exercise.ID,
			"name":        exercise.Name,
			"type":        exercise.Type,
			"duration":    exercise.Duration,
			"calories":    exercise.Calories,
			"date":        exercise.Date,
			"description": exercise.Description,
		}

		events[i] = StreamEvent{
			ID:        fmt.Sprintf("%s-%d-%d", streamName, time.Now().UnixNano(), exercise.ID),
			Timestamp: time.Now(),
			Type:      StreamEventInsert,
			Version:   d.currentVersion,
			Data:      data,
			Metadata: map[string]string{
				"stream": streamName,
				"source": "lakehouse",
			},
		}
	}

	// Publish to the stream
	return stream.Publish(ctx, events)
}

// SubscribeToStream subscribes to a stream
func (d *DeltaLakeRepository) SubscribeToStream(ctx context.Context, streamName string) (<-chan StreamEvent, error) {
	d.mutex.RLock()
	stream, exists := d.streams[streamName]
	d.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("stream %s not found", streamName)
	}

	return stream.Subscribe(ctx)
}

// GetActiveStreams returns information about active streams
func (d *DeltaLakeRepository) GetActiveStreams(ctx context.Context) ([]StreamInfo, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	streamInfos := make([]StreamInfo, 0, len(d.streams))
	for _, stream := range d.streams {
		stats := stream.GetStats()
		streamInfos = append(streamInfos, StreamInfo{
			Name:         stream.Name(),
			Type:         stream.Type(),
			IsActive:     stream.IsActive(),
			CreatedAt:    time.Now(), // Would need to track creation time
			LastActivity: stats.LastEventTime,
			Stats:        stats,
		})
	}

	return streamInfos, nil
}
