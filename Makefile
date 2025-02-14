VERSION ?= $(shell git describe --tags --always --dirty)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: all build test clean release validate ai-metadata

all: build test ai-metadata

# Standard build for CLI users
build:
	go build -v -ldflags "-X github.com/Nexlayer/nexlayer-cli/pkg/version.Version=$(VERSION)" \
		-o bin/nexlayer .

# Generate LLM-optimized metadata
ai-metadata:
	@echo "Generating LLM-optimized metadata..."
	@mkdir -p ai_training/metadata
	@go run tools/llm/main.go
	@echo "Copying example templates with LLM annotations..."
	@mkdir -p ai_training/examples
	@for file in examples/templates/*.yaml; do \
		python tools/llm/annotate.py $$file ai_training/examples/$$(basename $$file); \
	done
	@echo "Generating semantic search index..."
	@python tools/llm/index.py ai_training/metadata/llm_metadata.json ai_training/metadata/semantic_index.json

# Test all packages
test:
	go test -v -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/ build/
	go clean

# Validate Nexlayer configuration
validate:
	@echo "Validating nexlayer.yaml..."
	./bin/nexlayer validate

# Build release packages
release:
	@echo "Building release for $(GOOS)/$(GOARCH)..."
	mkdir -p dist/$(GOOS)_$(GOARCH)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-X github.com/Nexlayer/nexlayer-cli/pkg/version.Version=$(VERSION)" \
		-o dist/$(GOOS)_$(GOARCH)/nexlayer .

# Validate templates
template-validate:
	@echo "Validating template schema..."
	@for file in templates/*.yaml; do \
		echo "Validating $$file..."; \
		./bin/nexlayer validate -f $$file || exit 1; \
	done

# Build test container
docker-build:
	docker build -t nexlayer/cli-test:$(VERSION) \
		--build-arg VERSION=$(VERSION) \
		-f Dockerfile .

.PHONY: docker-build template-validate

