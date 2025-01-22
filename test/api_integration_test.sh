#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Mock server port
PORT=8080

# Start mock API server
echo "Starting mock API server..."
go run test/mock_api.go &
MOCK_PID=$!

# Wait for server to start
sleep 2

# Build the CLI
echo "Building CLI..."
go build -o nexlayer

# Set up test environment
export NEXLAYER_AUTH_TOKEN="test-token"
TEST_APP="test-app"
TEST_SERVICE="test-service"
TEST_DOMAIN="example.com"

# Helper function to run a test
run_test() {
    local name=$1
    local cmd=$2
    local expected_output=$3
    
    echo -n "Testing $name... "
    output=$(eval "$cmd" 2>&1)
    
    if echo "$output" | grep -q "$expected_output"; then
        echo -e "${GREEN}PASSED${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}FAILED${NC}"
        echo "Expected output containing: $expected_output"
        echo "Got: $output"
        ((TESTS_FAILED++))
    fi
}

# Create test deployment YAML
cat > test-deploy.yaml << EOL
services:
  frontend:
    image: nginx:latest
    ports:
      - "80:80"
EOL

# Test cases

# 1. Test deployment
run_test "deployment" \
    "./nexlayer deploy --app $TEST_APP --file test-deploy.yaml --url http://localhost:$PORT" \
    "Deployment started"

# 2. Test custom domain
run_test "custom domain" \
    "./nexlayer domain --app $TEST_APP --domain $TEST_DOMAIN --url http://localhost:$PORT" \
    "Domain saved"

# 3. Test listing deployments
run_test "list deployments" \
    "./nexlayer list --app $TEST_APP --url http://localhost:$PORT" \
    "Template: K-d chat"

# 4. Test deployment info
run_test "deployment info" \
    "./nexlayer info --app $TEST_APP --url http://localhost:$PORT" \
    "Template: K-d chat"

# Clean up
echo "Cleaning up..."
kill $MOCK_PID
rm -f nexlayer test-deploy.yaml

# Print summary
echo "===================="
echo "Test Summary:"
echo "------------------"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"
echo "===================="

# Exit with failure if any tests failed
[ $TESTS_FAILED -eq 0 ] || exit 1
