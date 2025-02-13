#!/bin/bash

# Build the CLI
make build

# Create test directory
mkdir -p test-projects
cd test-projects

# Test Node.js project
echo "Testing Node.js project..."
mkdir -p node-test
cd node-test
npm init -y
../../bin/nexlayer init > nexlayer.yaml
../../bin/nexlayer validate
cd ..

# Test React project
echo "Testing React project..."
mkdir -p react-test
cd react-test
npx create-react-app . --template typescript
../../bin/nexlayer init > nexlayer.yaml
../../bin/nexlayer validate
cd ..

# Test Python project
echo "Testing Python project..."
mkdir -p python-test
cd python-test
touch requirements.txt
../../bin/nexlayer init > nexlayer.yaml
../../bin/nexlayer validate
cd ..

# Test Go project
echo "Testing Go project..."
mkdir -p go-test
cd go-test
go mod init test-app
../../bin/nexlayer init > nexlayer.yaml
../../bin/nexlayer validate
cd ..