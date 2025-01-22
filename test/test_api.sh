#!/bin/bash

# Start mock API server
echo "Starting mock API server..."
go run test/mock_api.go &
MOCK_PID=$!

# Wait for server to start
sleep 2

# Build the CLI
go build -o nexlayer

# Set up mock environment
export NEXLAYER_AUTH_TOKEN="test-token"

# Test service configure command
echo "Testing service configure command..."
./nexlayer service configure --app test-app --service test-service --env "DB_URL=postgres://localhost:5432/db" --api-url "http://localhost:8080"

# Test service deploy command
echo -e "\nTesting service deploy command..."
./nexlayer service deploy --app test-app --service test-service --env "DB_URL=postgres://localhost:5432/db" --api-url "http://localhost:8080"

# Clean up
kill $MOCK_PID
rm nexlayer
