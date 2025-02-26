# Validation Package

This package provides validation utilities for Nexlayer configurations.

## Migration Notice

This package has been consolidated with the schema package. Please:

1. Replace imports from `github.com/Nexlayer/nexlayer-cli/pkg/validation` with `github.com/Nexlayer/nexlayer-cli/pkg/core/schema`
2. Use the validation functions provided by the schema package

## Example

```go
import "github.com/Nexlayer/nexlayer-cli/pkg/core/schema"

// Create a new schema service
service := schema.NewService()

// Validate a configuration
err := service.Validate(config)
if err != nil {
    log.Fatal(err)
}
```
