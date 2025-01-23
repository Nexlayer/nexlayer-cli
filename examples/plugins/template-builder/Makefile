.PHONY: all build test test-integration clean install

BINARY_NAME=nexlayer
VERSION=$(shell git describe --tags --always --dirty)
BUILD_FLAGS=-ldflags "-X main.version=$(VERSION)"

all: clean build test

build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(BUILD_FLAGS) -o $(BINARY_NAME)

test:
	@echo "Running unit tests..."
	@go test -v -short ./...

test-integration:
	@echo "Running integration tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean

install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(BUILD_FLAGS)

lint:
	@echo "Running linters..."
	@golangci-lint run

coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

docs:
	@echo "Generating documentation..."
	@godoc -http=:6060

release:
	@echo "Creating release..."
	@goreleaser release --snapshot --rm-dist

.DEFAULT_GOAL := all
