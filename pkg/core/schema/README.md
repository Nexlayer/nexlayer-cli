# Schema Package

This package provides centralized schema management for Nexlayer YAML configurations. It handles validation, parsing, and processing of configuration files.

## Overview

The schema package is responsible for:
- Defining the structure of Nexlayer YAML configurations
- Validating configuration files against the schema
- Processing and generating configuration files
- Managing configuration templates and examples

## Components

- `types.go`: Core data structures for Nexlayer configurations
- `validation.go`: Configuration validation logic
- `generator.go`: Configuration generation utilities
- `processor.go`: Configuration processing and transformation
- `jsonschema.go`: JSON Schema definitions
- `errors.go`: Error types and handling
- `service.go`: High-level schema management services

## Usage

```go
import "github.com/Nexlayer/nexlayer-cli/pkg/core/schema"

// Create a new schema service
service := schema.NewService()

// Process a configuration file
config, err := service.ProcessFile("config.yaml")
if err != nil {
    log.Fatal(err)
}
```

For more examples, see the `examples` directory.
