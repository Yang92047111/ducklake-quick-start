package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/Yang92047111/ducklake-quick-start/internal/loader"
)

// Lakehouse feature implementations for DeltaLakeRepository

// GetByVersion retrieves data as it existed at a specific version
func (d *DeltaLakeRepository) GetByVersion(ctx context.Context, version int64) ([]loader.Exercise, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if _, exists := d.versions[version]; !exists {
		return nil, fmt.Errorf("version %d does not exist", version)
	}

	dataPath := filepath.Join(d.basePath, fmt.Sprintf("version_%d.json", version))
	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read version %d data: %w", version, err)
	}

	var exercises []loader.Exercise
	if err := json.Unmarshal(data, &exercises); err != nil {
		return nil, fmt.Errorf("failed to unmarshal version %d data: %w", version, err)
	}

	return exercises, nil
}

// GetByTimestamp retrieves data as it existed at a specific timestamp
func (d *DeltaLakeRepository) GetByTimestamp(ctx context.Context, timestamp time.Time) ([]loader.Exercise, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Find the latest version before or at the timestamp
	var targetVersion int64 = -1
	var closestTime time.Time

	for version, versionInfo := range d.versions {
		if versionInfo.Timestamp.Before(timestamp) || versionInfo.Timestamp.Equal(timestamp) {
			if targetVersion == -1 || versionInfo.Timestamp.After(closestTime) {
				targetVersion = version
				closestTime = versionInfo.Timestamp
			}
		}
	}

	if targetVersion == -1 {
		return []loader.Exercise{}, nil // No data existed at that time
	}

	return d.GetByVersion(ctx, targetVersion)
}

// GetVersionHistory returns the history of all versions
func (d *DeltaLakeRepository) GetVersionHistory(ctx context.Context) ([]Version, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	versions := make([]Version, 0, len(d.versions))
	for _, version := range d.versions {
		versions = append(versions, *version)
	}

	// Sort by version ID
	for i := 0; i < len(versions)-1; i++ {
		for j := i + 1; j < len(versions); j++ {
			if versions[i].ID > versions[j].ID {
				versions[i], versions[j] = versions[j], versions[i]
			}
		}
	}

	return versions, nil
}

// CreateVersion creates a new version with the current state
func (d *DeltaLakeRepository) CreateVersion(ctx context.Context, description string) (*Version, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	exercises, err := d.getAllFromFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get current data: %w", err)
	}

	newVersionID := d.currentVersion + 1

	// Save current data as new version
	if err := d.saveVersionData(newVersionID, exercises); err != nil {
		return nil, fmt.Errorf("failed to save version data: %w", err)
	}

	version := &Version{
		ID:          newVersionID,
		Timestamp:   time.Now(),
		Description: description,
		SchemaID:    d.currentSchema.ID,
		RecordCount: int64(len(exercises)),
		FileCount:   1, // Simplified for this implementation
		SizeBytes:   0, // Would calculate actual size in production
		Operations:  []Operation{},
		ParentID:    &d.currentVersion,
	}

	d.versions[newVersionID] = version
	d.currentVersion = newVersionID

	// Update metadata
	d.metadata.CurrentVersion = newVersionID
	d.metadata.LastModified = time.Now()

	if err := d.saveMetadata(); err != nil {
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	return version, nil
}

// GetCurrentSchema returns the current schema
func (d *DeltaLakeRepository) GetCurrentSchema(ctx context.Context) (*Schema, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Deep copy to prevent modifications
	schemaCopy := *d.currentSchema
	schemaCopy.Fields = make([]Field, len(d.currentSchema.Fields))
	copy(schemaCopy.Fields, d.currentSchema.Fields)

	return &schemaCopy, nil
}

// EvolveSchema evolves the table schema
func (d *DeltaLakeRepository) EvolveSchema(ctx context.Context, newSchema *Schema) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Validate schema compatibility
	if err := d.validateSchemaCompatibilityInternal(newSchema); err != nil {
		return fmt.Errorf("schema is not compatible: %w", err)
	}

	// Create new schema version
	newSchema.ID = d.currentSchema.ID + 1
	newSchema.Version = d.currentSchema.Version + 1
	newSchema.CreatedAt = time.Now()

	// Save previous schema for history
	oldSchema := *d.currentSchema

	// Update current schema
	d.currentSchema = newSchema

	// Log schema change operation
	operation := Operation{
		Type:      OperationTypeSchema,
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"old_schema_id": oldSchema.ID,
			"new_schema_id": newSchema.ID,
			"changes":       d.calculateSchemaChanges(&oldSchema, newSchema),
		},
	}

	// Create new version for schema change
	d.currentVersion++
	version := &Version{
		ID:          d.currentVersion,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Schema evolved to version %d", newSchema.Version),
		SchemaID:    newSchema.ID,
		Operations:  []Operation{operation},
		ParentID:    &d.currentVersion,
	}

	d.versions[d.currentVersion] = version

	// Update metadata
	d.metadata.CurrentVersion = d.currentVersion
	d.metadata.LastModified = time.Now()

	return d.saveMetadata()
}

// GetSchemaHistory returns the history of schema changes
func (d *DeltaLakeRepository) GetSchemaHistory(ctx context.Context) ([]Schema, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// In a full implementation, this would read from schema history files
	// For now, return current schema
	return []Schema{*d.currentSchema}, nil
}

// ValidateSchemaCompatibility validates if a new schema is compatible
func (d *DeltaLakeRepository) ValidateSchemaCompatibility(ctx context.Context, newSchema *Schema) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.validateSchemaCompatibilityInternal(newSchema)
}

// validateSchemaCompatibilityInternal internal schema validation
func (d *DeltaLakeRepository) validateSchemaCompatibilityInternal(newSchema *Schema) error {
	currentFields := make(map[string]Field)
	for _, field := range d.currentSchema.Fields {
		currentFields[field.Name] = field
	}

	for _, newField := range newSchema.Fields {
		if currentField, exists := currentFields[newField.Name]; exists {
			// Check if type change is compatible
			if !d.isTypeCompatible(currentField.Type, newField.Type) {
				return fmt.Errorf("field %s type change from %s to %s is not compatible",
					newField.Name, currentField.Type, newField.Type)
			}

			// Check if nullable change is compatible
			if currentField.Nullable && !newField.Nullable {
				return fmt.Errorf("field %s cannot change from nullable to non-nullable", newField.Name)
			}
		} else {
			// New field must be nullable or have a default value
			if !newField.Nullable && newField.DefaultValue == nil {
				return fmt.Errorf("new field %s must be nullable or have a default value", newField.Name)
			}
		}
	}

	return nil
}

// isTypeCompatible checks if type changes are compatible
func (d *DeltaLakeRepository) isTypeCompatible(oldType, newType FieldType) bool {
	if oldType == newType {
		return true
	}

	// Define compatible type transitions
	compatibleTransitions := map[FieldType][]FieldType{
		FieldTypeInt:    {FieldTypeFloat, FieldTypeString},
		FieldTypeFloat:  {FieldTypeString},
		FieldTypeString: {}, // String can't be converted to other types safely
		FieldTypeDate:   {FieldTypeTimestamp, FieldTypeString},
	}

	allowedTypes, exists := compatibleTransitions[oldType]
	if !exists {
		return false
	}

	for _, allowedType := range allowedTypes {
		if newType == allowedType {
			return true
		}
	}

	return false
}

// calculateSchemaChanges calculates the differences between schemas
func (d *DeltaLakeRepository) calculateSchemaChanges(oldSchema, newSchema *Schema) map[string]interface{} {
	changes := map[string]interface{}{
		"added_fields":    []string{},
		"removed_fields":  []string{},
		"modified_fields": []string{},
	}

	oldFields := make(map[string]Field)
	for _, field := range oldSchema.Fields {
		oldFields[field.Name] = field
	}

	newFields := make(map[string]Field)
	for _, field := range newSchema.Fields {
		newFields[field.Name] = field
	}

	// Find added and modified fields
	addedFields := []string{}
	modifiedFields := []string{}
	for name, newField := range newFields {
		if oldField, exists := oldFields[name]; exists {
			if oldField.Type != newField.Type || oldField.Nullable != newField.Nullable {
				modifiedFields = append(modifiedFields, name)
			}
		} else {
			addedFields = append(addedFields, name)
		}
	}

	// Find removed fields
	removedFields := []string{}
	for name := range oldFields {
		if _, exists := newFields[name]; !exists {
			removedFields = append(removedFields, name)
		}
	}

	changes["added_fields"] = addedFields
	changes["removed_fields"] = removedFields
	changes["modified_fields"] = modifiedFields

	return changes
}

// GetTableMetadata returns comprehensive table metadata
func (d *DeltaLakeRepository) GetTableMetadata(ctx context.Context) (*TableMetadata, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Create a copy to prevent modifications
	metadata := *d.metadata
	metadata.Properties = make(map[string]string)
	for k, v := range d.metadata.Properties {
		metadata.Properties[k] = v
	}

	return &metadata, nil
}

// UpdateTableProperties updates table properties
func (d *DeltaLakeRepository) UpdateTableProperties(ctx context.Context, properties map[string]string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.metadata.Properties == nil {
		d.metadata.Properties = make(map[string]string)
	}

	for k, v := range properties {
		d.metadata.Properties[k] = v
	}

	d.metadata.LastModified = time.Now()
	return d.saveMetadata()
}

// GetPartitions returns partition information (simplified implementation)
func (d *DeltaLakeRepository) GetPartitions(ctx context.Context) ([]Partition, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// In a real implementation, this would read partition information
	// For now, return a single partition representing all data
	exercises, err := d.getAllFromFiles()
	if err != nil {
		return nil, err
	}

	partition := Partition{
		Values:       map[string]interface{}{},
		RecordCount:  int64(len(exercises)),
		FileCount:    1,
		SizeBytes:    0, // Would calculate actual size
		Location:     d.basePath,
		LastModified: d.metadata.LastModified,
	}

	return []Partition{partition}, nil
}

// OptimizeTable optimizes the table layout and files
func (d *DeltaLakeRepository) OptimizeTable(ctx context.Context, options OptimizeOptions) (*OptimizeResult, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	startTime := time.Now()

	// In a real implementation, this would:
	// 1. Analyze file sizes and layout
	// 2. Compact small files
	// 3. Rewrite large files
	// 4. Apply Z-ordering if specified
	// 5. Update statistics

	// For this implementation, we'll simulate optimization
	exercises, err := d.getAllFromFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to read current data: %w", err)
	}

	// Create optimized version
	d.currentVersion++
	if err := d.saveVersionData(d.currentVersion, exercises); err != nil {
		return nil, fmt.Errorf("failed to save optimized data: %w", err)
	}

	// Update metadata
	d.metadata.CurrentVersion = d.currentVersion
	d.metadata.LastModified = time.Now()

	result := &OptimizeResult{
		FilesAdded:          1,
		FilesRemoved:        1,
		PartitionsOptimized: 1,
		RecordsRewritten:    int64(len(exercises)),
		BytesWritten:        0, // Would calculate actual bytes
		BytesRemoved:        0, // Would calculate actual bytes
		Duration:            time.Since(startTime),
		Metrics: map[string]interface{}{
			"optimization_type": "full_table",
			"compression_ratio": 1.0,
		},
	}

	if err := d.saveMetadata(); err != nil {
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	return result, nil
}

// Data Quality and Constraints Implementation

// AddConstraint adds a data quality constraint
func (d *DeltaLakeRepository) AddConstraint(ctx context.Context, constraint Constraint) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check if constraint with same name already exists
	for _, existing := range d.constraints {
		if existing.Name == constraint.Name {
			return fmt.Errorf("constraint with name %s already exists", constraint.Name)
		}
	}

	constraint.CreatedAt = time.Now()
	constraint.Enabled = true
	d.constraints = append(d.constraints, constraint)

	return d.saveMetadata()
}

// RemoveConstraint removes a data quality constraint
func (d *DeltaLakeRepository) RemoveConstraint(ctx context.Context, constraintName string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for i, constraint := range d.constraints {
		if constraint.Name == constraintName {
			// Remove constraint by replacing with last element and truncating
			d.constraints[i] = d.constraints[len(d.constraints)-1]
			d.constraints = d.constraints[:len(d.constraints)-1]
			return d.saveMetadata()
		}
	}

	return fmt.Errorf("constraint %s not found", constraintName)
}

// ValidateConstraints validates data against defined constraints
func (d *DeltaLakeRepository) ValidateConstraints(ctx context.Context, exercises []loader.Exercise) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, constraint := range d.constraints {
		if !constraint.Enabled {
			continue
		}

		for _, exercise := range exercises {
			if err := d.validateConstraint(constraint, exercise); err != nil {
				return fmt.Errorf("constraint %s violated: %w", constraint.Name, err)
			}
		}
	}

	return nil
}

// validateConstraint validates a single exercise against a constraint
func (d *DeltaLakeRepository) validateConstraint(constraint Constraint, exercise loader.Exercise) error {
	switch constraint.Type {
	case ConstraintTypeNotNull:
		return d.validateNotNull(constraint, exercise)
	case ConstraintTypeRange:
		return d.validateRange(constraint, exercise)
	case ConstraintTypeCheck:
		return d.validateCheck(constraint, exercise)
	default:
		return fmt.Errorf("unsupported constraint type: %s", constraint.Type)
	}
}

// validateNotNull validates not null constraints
func (d *DeltaLakeRepository) validateNotNull(constraint Constraint, exercise loader.Exercise) error {
	for _, column := range constraint.Columns {
		switch column {
		case "name":
			if exercise.Name == "" {
				return fmt.Errorf("name cannot be null")
			}
		case "type":
			if exercise.Type == "" {
				return fmt.Errorf("type cannot be null")
			}
		case "description":
			if exercise.Description == "" && len(constraint.Columns) == 1 {
				return fmt.Errorf("description cannot be null")
			}
		}
	}
	return nil
}

// validateRange validates range constraints
func (d *DeltaLakeRepository) validateRange(constraint Constraint, exercise loader.Exercise) error {
	// Simplified range validation - in production, this would parse the expression
	for _, column := range constraint.Columns {
		switch column {
		case "duration":
			if exercise.Duration < 0 || exercise.Duration > 1440 { // 24 hours max
				return fmt.Errorf("duration must be between 0 and 1440 minutes")
			}
		case "calories":
			if exercise.Calories < 0 || exercise.Calories > 10000 {
				return fmt.Errorf("calories must be between 0 and 10000")
			}
		}
	}
	return nil
}

// validateCheck validates check constraints
func (d *DeltaLakeRepository) validateCheck(constraint Constraint, exercise loader.Exercise) error {
	// Simplified check validation - in production, this would evaluate the expression
	if constraint.Expression == "type IN ('cardio', 'strength', 'flexibility')" {
		validTypes := map[string]bool{"cardio": true, "strength": true, "flexibility": true}
		if !validTypes[exercise.Type] {
			return fmt.Errorf("invalid exercise type: %s", exercise.Type)
		}
	}
	return nil
}

// GetDataQualityMetrics returns data quality metrics
func (d *DeltaLakeRepository) GetDataQualityMetrics(ctx context.Context) (*DataQualityMetrics, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	exercises, err := d.getAllFromFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get data for quality analysis: %w", err)
	}

	metrics := &DataQualityMetrics{
		TotalRecords:         int64(len(exercises)),
		NullValues:           make(map[string]int64),
		UniqueValues:         make(map[string]int64),
		DuplicateRecords:     0,
		ConstraintViolations: make(map[string]int64),
		DataTypeMismatches:   make(map[string]int64),
		OutlierCounts:        make(map[string]int64),
		LastUpdated:          time.Now(),
	}

	// Calculate metrics
	uniqueNames := make(map[string]bool)
	uniqueTypes := make(map[string]bool)
	seenExercises := make(map[string]bool)

	for _, exercise := range exercises {
		// Check for nulls
		if exercise.Name == "" {
			metrics.NullValues["name"]++
		} else {
			uniqueNames[exercise.Name] = true
		}

		if exercise.Type == "" {
			metrics.NullValues["type"]++
		} else {
			uniqueTypes[exercise.Type] = true
		}

		if exercise.Description == "" {
			metrics.NullValues["description"]++
		}

		// Check for duplicates (simplified - by name and date)
		key := fmt.Sprintf("%s-%s", exercise.Name, exercise.Date.Format("2006-01-02"))
		if seenExercises[key] {
			metrics.DuplicateRecords++
		} else {
			seenExercises[key] = true
		}

		// Check for outliers
		if exercise.Duration > 240 { // More than 4 hours
			metrics.OutlierCounts["duration"]++
		}
		if exercise.Calories > 2000 {
			metrics.OutlierCounts["calories"]++
		}
	}

	metrics.UniqueValues["name"] = int64(len(uniqueNames))
	metrics.UniqueValues["type"] = int64(len(uniqueTypes))

	// Calculate quality scores
	totalFields := int64(len(exercises) * 7) // 7 fields per exercise
	nullCount := metrics.NullValues["name"] + metrics.NullValues["type"] + metrics.NullValues["description"]

	metrics.CompletenessScore = float64(totalFields-nullCount) / float64(totalFields) * 100
	metrics.ValidityScore = 100.0 // Would calculate based on constraint violations
	metrics.ConsistencyScore = float64(int64(len(exercises))-metrics.DuplicateRecords) / float64(len(exercises)) * 100

	return metrics, nil
}

// Streaming and Change Data Capture Implementation

// WatchChanges watches for changes and returns a channel of change events
func (d *DeltaLakeRepository) WatchChanges(ctx context.Context, from time.Time) (<-chan ChangeEvent, error) {
	changeChan := make(chan ChangeEvent, 100)

	// This is a simplified implementation
	// In production, this would use file system watchers or database triggers
	go func() {
		defer close(changeChan)

		// Send existing changes since 'from' time
		for _, change := range d.changeLog {
			if change.Timestamp.After(from) {
				select {
				case changeChan <- change:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return changeChan, nil
}

// GetChangelog returns change events between versions
func (d *DeltaLakeRepository) GetChangelog(ctx context.Context, fromVersion, toVersion int64) ([]ChangeEvent, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var changes []ChangeEvent

	for _, change := range d.changeLog {
		if change.Version >= fromVersion && change.Version <= toVersion {
			changes = append(changes, change)
		}
	}

	return changes, nil
}

// Performance and Optimization Implementation

// CreateIndex creates an index on specified columns
func (d *DeltaLakeRepository) CreateIndex(ctx context.Context, indexName string, columns []string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, exists := d.indexes[indexName]; exists {
		return fmt.Errorf("index %s already exists", indexName)
	}

	index := &Index{
		Name:      indexName,
		Columns:   columns,
		Type:      IndexTypeBTree, // Default to B-tree
		CreatedAt: time.Now(),
		Stats: &IndexStats{
			Uses:        0,
			Selectivity: 1.0, // Would calculate based on data
			Size:        0,
		},
	}

	d.indexes[indexName] = index
	return d.saveMetadata()
}

// DropIndex drops an index
func (d *DeltaLakeRepository) DropIndex(ctx context.Context, indexName string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, exists := d.indexes[indexName]; !exists {
		return fmt.Errorf("index %s does not exist", indexName)
	}

	delete(d.indexes, indexName)
	return d.saveMetadata()
}

// GetQueryStats returns query performance statistics
func (d *DeltaLakeRepository) GetQueryStats(ctx context.Context) (*QueryStats, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Return current stats (would be populated by actual queries)
	return d.queryStats, nil
}

// Compact compacts table files
func (d *DeltaLakeRepository) Compact(ctx context.Context) (*CompactionResult, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	startTime := time.Now()

	// Simplified compaction - in production, this would merge small files
	exercises, err := d.getAllFromFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to read data for compaction: %w", err)
	}

	// Create new compacted version
	d.currentVersion++
	if err := d.saveVersionData(d.currentVersion, exercises); err != nil {
		return nil, fmt.Errorf("failed to save compacted data: %w", err)
	}

	result := &CompactionResult{
		FilesCompacted:   1,
		FilesCreated:     1,
		RecordsProcessed: int64(len(exercises)),
		SpaceReclaimed:   0, // Would calculate actual space savings
		Duration:         time.Since(startTime),
	}

	// Update metadata
	d.metadata.CurrentVersion = d.currentVersion
	d.metadata.LastModified = time.Now()

	if err := d.saveMetadata(); err != nil {
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	return result, nil
}

// Advanced Querying Implementation (simplified stubs)

// QueryWithSQL executes SQL queries (not implemented in this demo)
func (d *DeltaLakeRepository) QueryWithSQL(ctx context.Context, sql string, params ...interface{}) ([]loader.Exercise, error) {
	return nil, fmt.Errorf("SQL querying not implemented in this demo")
}

// QueryWithFilter executes queries with advanced filtering
func (d *DeltaLakeRepository) QueryWithFilter(ctx context.Context, filter Filter) ([]loader.Exercise, error) {
	exercises, err := d.getAllFromFiles()
	if err != nil {
		return nil, err
	}

	// Apply basic filtering (simplified implementation)
	var result []loader.Exercise
	for _, exercise := range exercises {
		if d.matchesFilter(exercise, filter) {
			result = append(result, exercise)
		}
	}

	// Apply sorting, limit, offset
	result = d.applySorting(result, filter.SortBy)
	result = d.applyPagination(result, filter.Limit, filter.Offset)

	return result, nil
}

// matchesFilter checks if an exercise matches filter conditions
func (d *DeltaLakeRepository) matchesFilter(exercise loader.Exercise, filter Filter) bool {
	for _, condition := range filter.Conditions {
		if !d.matchesCondition(exercise, condition) {
			return false
		}
	}
	return true
}

// matchesCondition checks if an exercise matches a single condition
func (d *DeltaLakeRepository) matchesCondition(exercise loader.Exercise, condition Condition) bool {
	var fieldValue interface{}

	switch condition.Field {
	case "id":
		fieldValue = exercise.ID
	case "name":
		fieldValue = exercise.Name
	case "type":
		fieldValue = exercise.Type
	case "duration":
		fieldValue = exercise.Duration
	case "calories":
		fieldValue = exercise.Calories
	case "date":
		fieldValue = exercise.Date
	case "description":
		fieldValue = exercise.Description
	default:
		return false
	}

	// Simplified condition matching
	switch condition.Operator {
	case OperatorEqual:
		return fieldValue == condition.Value
	case OperatorNotEqual:
		return fieldValue != condition.Value
	case OperatorGreaterThan:
		if intVal, ok := fieldValue.(int); ok {
			if targetVal, ok := condition.Value.(int); ok {
				return intVal > targetVal
			}
		}
	case OperatorLike:
		if strVal, ok := fieldValue.(string); ok {
			if pattern, ok := condition.Value.(string); ok {
				return contains(strVal, pattern)
			}
		}
	}

	return false
}

// contains checks if a string contains a pattern (simplified)
func contains(str, pattern string) bool {
	return len(str) >= len(pattern) && str[:len(pattern)] == pattern
}

// applySorting applies sorting to results
func (d *DeltaLakeRepository) applySorting(exercises []loader.Exercise, sortFields []SortField) []loader.Exercise {
	// Simplified sorting implementation
	if len(sortFields) == 0 {
		return exercises
	}

	// For demo purposes, just sort by first field
	if len(sortFields) > 0 {
		sortField := sortFields[0]
		sort.Slice(exercises, func(i, j int) bool {
			var left, right interface{}

			switch sortField.Field {
			case "id":
				left, right = exercises[i].ID, exercises[j].ID
			case "name":
				left, right = exercises[i].Name, exercises[j].Name
			case "type":
				left, right = exercises[i].Type, exercises[j].Type
			case "duration":
				left, right = exercises[i].Duration, exercises[j].Duration
			case "calories":
				left, right = exercises[i].Calories, exercises[j].Calories
			case "date":
				left, right = exercises[i].Date, exercises[j].Date
			default:
				return false
			}

			// Basic comparison logic
			if sortField.Order == SortOrderDesc {
				return fmt.Sprintf("%v", left) > fmt.Sprintf("%v", right)
			}
			return fmt.Sprintf("%v", left) < fmt.Sprintf("%v", right)
		})
	}

	return exercises
}

// applyPagination applies limit and offset to results
func (d *DeltaLakeRepository) applyPagination(exercises []loader.Exercise, limit, offset *int) []loader.Exercise {
	start := 0
	if offset != nil {
		start = *offset
	}

	if start >= len(exercises) {
		return []loader.Exercise{}
	}

	end := len(exercises)
	if limit != nil {
		end = start + *limit
		if end > len(exercises) {
			end = len(exercises)
		}
	}

	return exercises[start:end]
}

// AggregateByTimeWindow performs time-based aggregations
func (d *DeltaLakeRepository) AggregateByTimeWindow(ctx context.Context, window TimeWindow, aggregations []Aggregation) ([]AggregationResult, error) {
	return nil, fmt.Errorf("time window aggregation not implemented in this demo")
}
