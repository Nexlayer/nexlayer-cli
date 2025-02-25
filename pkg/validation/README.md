# Nexlayer YAML Schema Validation [DEPRECATED]

⚠️ **IMPORTANT: This package is deprecated** ⚠️

The validation functionality has been consolidated into the `pkg/schema` package, which now serves as the single source of truth for all Nexlayer YAML schema definitions, validation, and utilities.

## Migration Guide

If your code currently depends on this package:

1. Replace imports from `github.com/Nexlayer/nexlayer-cli/pkg/validation` with `github.com/Nexlayer/nexlayer-cli/pkg/schema`
2. Use the `schema.Validator` and related types for validation 
3. Use the schema types directly instead of defining your own

Example:

```go
// Old way (deprecated)
import "github.com/Nexlayer/nexlayer-cli/pkg/validation"

// New way
import "github.com/Nexlayer/nexlayer-cli/pkg/schema"
```

Please refer to [pkg/schema/README.md](../schema/README.md) for the definitive schema documentation.
