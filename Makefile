# Makefile for Nexlayer CLI

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_NAME=nexlayer
VERSION?=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X github.com/Nexlayer/nexlayer-cli/pkg/version.Version=$(VERSION) \
                  -X github.com/Nexlayer/nexlayer-cli/pkg/version.Commit=$(COMMIT) \
                  -X github.com/Nexlayer/nexlayer-cli/pkg/version.BuildDate=$(DATE)"

# Build directories
BUILD_DIR=build
DIST_DIR=dist

# Test parameters
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html
TEST_FLAGS=-v -race -parallel 4
BENCH_FLAGS=-benchmem -bench=.

# Linting
GOLINT=golangci-lint
GOSEC=gosec

# Performance
GOMAXPROCS?=$(shell nproc)
export GOMAXPROCS

# Vendor management
VENDOR_SCRIPT=scripts/update-vendor.sh
GOFLAGS=-mod=vendor

.PHONY: all build build-dev clean test coverage lint fmt vet install uninstall help bench security docker vendor vendor-update

all: lint test build ## Run lint, test, and build

vendor: ## Update and optimize vendor directory
	@echo "Updating vendor directory..."
	@$(VENDOR_SCRIPT)

vendor-update: ## Force update all dependencies and vendor directory
	@echo "Force updating all dependencies..."
	@$(GOCMD) get -u ./...
	@$(VENDOR_SCRIPT)

build: vendor ## Build the binary with optimizations
	@echo "Building Nexlayer CLI..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GOBUILD) -trimpath -a -tags netgo,osusergo \
		-installsuffix netgo $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/nexlayer

build-dev: vendor ## Build for development without optimizations
	@echo "Building development version..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/nexlayer

build-debug: ## Build with debug information
    @echo "Building debug version..."
    @mkdir -p $(BUILD_DIR)
    $(GOBUILD) -gcflags="all=-N -l" $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-debug ./cmd/nexlayer

clean: ## Clean build artifacts
    @echo "Cleaning build artifacts..."
    @rm -rf $(BUILD_DIR) $(DIST_DIR)
    @rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
    $(GOCMD) clean -cache -testcache -modcache

test: ## Run tests with race detection
    @echo "Running tests..."
    $(GOTEST) $(TEST_FLAGS) ./...

test-short: ## Run tests in short mode
    @echo "Running short tests..."
    $(GOTEST) -short $(TEST_FLAGS) ./...

coverage: ## Run tests with coverage and show percentage
    @echo "Running tests with coverage..."
    $(GOTEST) $(TEST_FLAGS) -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
    @go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print "Coverage: " $$3}'
    @go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
    @echo "Coverage report generated at $(COVERAGE_HTML)"

bench: ## Run benchmarks
    @echo "Running benchmarks..."
    $(GOTEST) $(TEST_FLAGS) $(BENCH_FLAGS) ./...

lint: ## Run linters
    @echo "Running linters..."
    $(GOLINT) run --timeout=5m

security: ## Run security checks
    @echo "Running security checks..."
    $(GOSEC) ./...

fmt: ## Format code
    @echo "Formatting code..."
    $(GOCMD) fmt ./...

vet: ## Run go vet
    @echo "Running go vet..."
    $(GOVET) -composites=false ./...

deps: ## Manage dependencies
    @if [ ! -d "vendor" ]; then \
        echo "Downloading dependencies..."; \
        $(GOCMD) mod download; \
    fi
    @echo "Verifying dependencies..."
    $(GOCMD) mod verify
    @echo "Tidying dependencies..."
    $(GOCMD) mod tidy
    @echo "Checking for vulnerable dependencies..."
    $(GOCMD) list -json -m all | nancy sleuth

install: build ## Install the CLI locally
    @echo "Installing Nexlayer CLI..."
    @cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
    @echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

uninstall: ## Uninstall the CLI
    @echo "Uninstalling Nexlayer CLI..."
    @rm -f $(GOPATH)/bin/$(BINARY_NAME)

release: ## Create a release build
    @echo "Creating release build..."
    @mkdir -p $(DIST_DIR)
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -trimpath -a -tags netgo,osusergo \
        -installsuffix netgo $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/nexlayer
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -trimpath -a -tags netgo,osusergo \
        -installsuffix netgo $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/nexlayer
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GOBUILD) -trimpath -a -tags netgo,osusergo \
        -installsuffix netgo $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/nexlayer
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -trimpath -a -tags netgo,osusergo \
        -installsuffix netgo $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/nexlayer
    @cd $(DIST_DIR) && \
        shasum -a 256 * > checksums.txt && \
        gpg --detach-sign --armor checksums.txt

docker: ## Build multi-arch Docker image
    @echo "Building multi-arch Docker image..."
    docker buildx build --platform linux/amd64,linux/arm64 -t nexlayer/cli:$(VERSION) \
        --build-arg VERSION=$(VERSION) \
        --build-arg COMMIT=$(COMMIT) \
        --build-arg BUILD_DATE=$(DATE) \
        -f Dockerfile .

ci: deps security lint test build ## Run all CI tasks

setup: ## Set up development environment
    @echo "Setting up development environment..."
    @if ! command -v $(GOLINT) > /dev/null; then \
        echo "Installing golangci-lint..." && \
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
    fi
    @if ! command -v $(GOSEC) > /dev/null; then \
        echo "Installing gosec..." && \
        curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
    fi
    @if ! command -v nancy > /dev/null; then \
        echo "Installing nancy..." && \
        curl -L https://github.com/sonatype-nexus-community/nancy/releases/latest/download/nancy-$(shell uname -s | tr '[:upper:]' '[:lower:]')-amd64 -o $$(go env GOPATH)/bin/nancy && \
        chmod +x $$(go env GOPATH)/bin/nancy; \
    fi
    $(GOCMD) mod download
    @echo "Development environment setup complete"

help: ## Show this help message
    @echo "Nexlayer CLI Makefile Commands:"
    @echo "  all            Run lint, test, and build"
    @echo "  build          Build the binary with optimizations"
    @echo "  build-dev      Build for development without optimizations"
    @echo "  build-debug    Build with debug information"
    @echo "  clean          Clean build artifacts"
    @echo "  test           Run tests with race detection"
    @echo "  test-short     Run tests in short mode"
    @echo "  coverage       Run tests with coverage and show percentage"
    @echo "  bench          Run benchmarks"
    @echo "  lint           Run linters"
    @echo "  security       Run security checks"
    @echo "  fmt            Format code"
    @echo "  vet            Run go vet"
    @echo "  deps           Manage dependencies"
    @echo "  install        Install the CLI locally"
    @echo "  uninstall      Uninstall the CLI"
    @echo "  release        Create a release build for multiple platforms"
    @echo "  docker         Build multi-arch Docker image"
    @echo "  ci             Run all CI tasks"
    @echo "  setup          Set up development environment"