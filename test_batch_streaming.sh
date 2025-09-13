#!/bin/bash

# Test script for DuckLake Batch and Streaming features
set -e

API_BASE="http://localhost:8080"
LAKEHOUSE_BASE="$API_BASE/api/v1/lakehouse"

echo "🦆 Testing DuckLake Batch and Streaming Features"
echo "================================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper function for colored output
print_step() {
    echo -e "${BLUE}➤ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Wait for server to be ready
wait_for_server() {
    print_step "Waiting for server to be ready..."
    local retries=30
    while [ $retries -gt 0 ]; do
        if curl -s "$API_BASE/health" > /dev/null 2>&1; then
            print_success "Server is ready"
            return 0
        fi
        ((retries--))
        sleep 1
    done
    print_error "Server failed to start"
    exit 1
}

# Test batch insert
test_batch_insert() {
    print_step "Testing batch insert..."
    
    cat > /tmp/batch_exercises.json << 'EOF'
{
    "exercises": [
        {
            "id": 100,
            "name": "Batch Running",
            "type": "cardio",
            "duration": 30,
            "calories": 300,
            "date": "2024-01-15T08:00:00Z",
            "description": "Morning batch run"
        },
        {
            "id": 101,
            "name": "Batch Cycling",
            "type": "cardio",
            "duration": 45,
            "calories": 400,
            "date": "2024-01-15T09:00:00Z",
            "description": "Batch cycling session"
        },
        {
            "id": 102,
            "name": "Batch Strength",
            "type": "strength",
            "duration": 60,
            "calories": 200,
            "date": "2024-01-15T10:00:00Z",
            "description": "Batch strength training"
        }
    ],
    "options": {
        "batch_size": 2,
        "skip_errors": true,
        "validate_first": true
    }
}
EOF

    response=$(curl -s -X POST "$LAKEHOUSE_BASE/batch/insert" \
        -H "Content-Type: application/json" \
        -d @/tmp/batch_exercises.json)
    
    if echo "$response" | grep -q "processed_count"; then
        print_success "Batch insert completed"
        echo "Response: $response"
    else
        print_error "Batch insert failed: $response"
    fi
}

# Test batch update
test_batch_update() {
    print_step "Testing batch update..."
    
    cat > /tmp/batch_update.json << 'EOF'
[
    {
        "id": 100,
        "name": "Updated Batch Running",
        "type": "cardio",
        "duration": 35,
        "calories": 350,
        "date": "2024-01-15T08:00:00Z",
        "description": "Updated morning batch run"
    },
    {
        "id": 101,
        "name": "Updated Batch Cycling",
        "type": "cardio",
        "duration": 50,
        "calories": 450,
        "date": "2024-01-15T09:00:00Z",
        "description": "Updated batch cycling session"
    }
]
EOF

    response=$(curl -s -X PUT "$LAKEHOUSE_BASE/batch/update" \
        -H "Content-Type: application/json" \
        -d @/tmp/batch_update.json)
    
    if echo "$response" | grep -q "processed_count"; then
        print_success "Batch update completed"
        echo "Response: $response"
    else
        print_error "Batch update failed: $response"
    fi
}

# Test bulk load
test_bulk_load() {
    print_step "Testing bulk load..."
    
    cat > /tmp/bulk_load_request.json << 'EOF'
{
    "data_source": {
        "type": "file",
        "location": "test/testdata/sample_exercises.csv",
        "format": "csv"
    },
    "options": {
        "batch_size": 10,
        "skip_errors": true,
        "parallel_jobs": 2
    }
}
EOF

    response=$(curl -s -X POST "$LAKEHOUSE_BASE/bulk-load" \
        -H "Content-Type: application/json" \
        -d @/tmp/bulk_load_request.json)
    
    if echo "$response" | grep -q "records_loaded"; then
        print_success "Bulk load completed"
        echo "Response: $response"
    else
        print_error "Bulk load failed: $response"
    fi
}

# Test stream creation
test_stream_creation() {
    print_step "Testing stream creation..."
    
    cat > /tmp/stream_config.json << 'EOF'
{
    "config": {
        "name": "exercise-stream",
        "type": "ingestion",
        "buffer_size": 1000,
        "flush_interval": "5s"
    }
}
EOF

    response=$(curl -s -X POST "$LAKEHOUSE_BASE/streams" \
        -H "Content-Type: application/json" \
        -d @/tmp/stream_config.json)
    
    if echo "$response" | grep -q "exercise-stream"; then
        print_success "Stream created successfully"
        echo "Response: $response"
    else
        print_error "Stream creation failed: $response"
    fi
}

# Test stream publishing
test_stream_publishing() {
    print_step "Testing stream publishing..."
    
    cat > /tmp/stream_data.json << 'EOF'
{
    "exercises": [
        {
            "id": 200,
            "name": "Streaming Yoga",
            "type": "flexibility",
            "duration": 60,
            "calories": 150,
            "date": "2024-01-15T11:00:00Z",
            "description": "Live streaming yoga session"
        },
        {
            "id": 201,
            "name": "Streaming HIIT",
            "type": "cardio",
            "duration": 20,
            "calories": 300,
            "date": "2024-01-15T12:00:00Z",
            "description": "High intensity streaming workout"
        }
    ]
}
EOF

    response=$(curl -s -X POST "$LAKEHOUSE_BASE/streams/exercise-stream/publish" \
        -H "Content-Type: application/json" \
        -d @/tmp/stream_data.json)
    
    if echo "$response" | grep -q "events_published"; then
        print_success "Stream publishing completed"
        echo "Response: $response"
    else
        print_error "Stream publishing failed: $response"
    fi
}

# Test stream status
test_stream_status() {
    print_step "Testing stream status..."
    
    response=$(curl -s -X GET "$LAKEHOUSE_BASE/streams")
    
    if echo "$response" | grep -q "streams"; then
        print_success "Stream status retrieved"
        echo "Response: $response"
    else
        print_error "Stream status failed: $response"
    fi
}

# Test batch delete
test_batch_delete() {
    print_step "Testing batch delete..."
    
    response=$(curl -s -X DELETE "$LAKEHOUSE_BASE/batch/delete" \
        -H "Content-Type: application/json" \
        -d '[100, 101, 102, 200, 201]')
    
    if echo "$response" | grep -q "processed_count"; then
        print_success "Batch delete completed"
        echo "Response: $response"
    else
        print_error "Batch delete failed: $response"
    fi
}

# Main test execution
main() {
    echo "Starting batch and streaming tests..."
    
    wait_for_server
    
    echo ""
    echo "🔄 Testing Batch Operations"
    echo "=========================="
    test_batch_insert
    echo ""
    test_batch_update
    echo ""
    test_bulk_load
    echo ""
    
    echo "🌊 Testing Streaming Operations"
    echo "=============================="
    test_stream_creation
    echo ""
    test_stream_publishing
    echo ""
    test_stream_status
    echo ""
    
    echo "🗑️  Testing Cleanup"
    echo "=================="
    test_batch_delete
    echo ""
    
    print_success "All batch and streaming tests completed!"
    
    # Cleanup
    rm -f /tmp/batch_exercises.json /tmp/batch_update.json /tmp/bulk_load_request.json
    rm -f /tmp/stream_config.json /tmp/stream_data.json
}

# Run the tests
main
