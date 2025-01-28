#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Test counter
TOTAL_TESTS=0
PASSED_TESTS=0

# Timing function
time_cmd() {
    local start_time=$(date +%s)
    eval "$1"
    local status=$?
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    echo "$duration"
    return $status
}

# Test function
run_test() {
    local name=$1
    local cmd=$2
    local expected_status=$3
    
    echo "Running test: $name"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Run command and capture output and status
    local output
    local status
    local duration
    
    output=$(eval "$cmd" 2>&1)
    status=$?
    duration=0
    
    # Check if status matches expected
    if [ $status -eq $expected_status ]; then
        echo -e "${GREEN}✓ Test passed${NC} (${duration}s)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗ Test failed${NC} (${duration}s)"
        echo "Expected status: $expected_status, got: $status"
        echo "Output: $output"
    fi
    echo "----------------------------------------"
}

# Clean up previous test artifacts
rm -rf test_projects
rm -f nexlayer.yaml

# Create test directories
mkdir -p test_projects/nodejs
mkdir -p test_projects/python
mkdir -p test_projects/empty

# Create test files
cat > test_projects/nodejs/package.json << EOF
{
    "name": "test-node-app",
    "dependencies": {
        "next": "^13.0.0",
        "langchain": "^0.0.200"
    }
}
EOF

cat > test_projects/python/requirements.txt << EOF
langchain==0.1.0
fastapi==0.100.0
EOF

# Test Cases

# 1. Basic Commands
run_test "Help Command" "./nexlayer --help" 0
run_test "Init Help" "./nexlayer init --help" 0

# 2. Project Initialization
run_test "Init LangChain Next.js" "cd test_projects/nodejs && ../../nexlayer init test-next -t langchain-nextjs" 0
run_test "Init LangChain FastAPI" "cd test_projects/python && ../../nexlayer init test-fastapi -t langchain-fastapi" 0

# 3. Error Cases
run_test "Invalid Template" "./nexlayer init test -t invalid-template" 1
run_test "Missing Project Name" "./nexlayer init" 1
run_test "Missing Template" "./nexlayer init test" 0

# 4. Project Detection
run_test "Detect Node.js Project" "cd test_projects/nodejs && ../../nexlayer init auto-detect" 0
run_test "Detect Python Project" "cd test_projects/python && ../../nexlayer init auto-detect" 0

# 5. Performance Tests
for i in {1..5}; do
    run_test "Performance Test $i" "./nexlayer init perf-test-$i -t langchain-nextjs" 0
done

# 6. Concurrent Tests
for i in {1..3}; do
    ./nexlayer init concurrent-$i -t langchain-nextjs &
done
wait

# Print summary
echo "Test Summary:"
echo "Total tests: $TOTAL_TESTS"
echo "Passed tests: $PASSED_TESTS"
echo "Failed tests: $((TOTAL_TESTS - PASSED_TESTS))"
echo "Success rate: $(( (PASSED_TESTS * 100) / TOTAL_TESTS ))%"

# Clean up
rm -rf test_projects
