GOPATH := $(shell go env GOPATH)
PATH := $(GOPATH)/bin:$(PATH)

.PHONY: all build metadata clean install-tools ensure-deps check-deps

all: ensure-deps build metadata

build:
	go build -v ./...

metadata: ensure-deps
	@echo "Generating project metadata..."
	@mkdir -p build
	@PATH=$(GOPATH)/bin:$(PATH) go run tools/metadata/main.go

clean:
	rm -rf build/
	go clean

check-deps:
	@echo "Checking dependencies..."
	@if ! which dot > /dev/null; then \
		echo "graphviz not found"; \
		exit 1; \
	fi
	@if ! which go-callvis > /dev/null; then \
		echo "go-callvis not found"; \
		exit 1; \
	fi

# Ensure all dependencies are installed
ensure-deps:
	@echo "Checking and installing dependencies..."
	@if ! which dot > /dev/null; then \
		echo "Installing graphviz..." && \
		brew install graphviz || { echo "Failed to install graphviz"; exit 1; }; \
	else \
		echo "graphviz is already installed"; \
	fi
	@if ! which go-callvis > /dev/null; then \
		echo "Installing go-callvis..." && \
		go install github.com/ofabry/go-callvis@latest || { echo "Failed to install go-callvis"; exit 1; }; \
	else \
		echo "go-callvis is already installed"; \
	fi

.PHONY: install-tools
install-tools: ensure-deps
