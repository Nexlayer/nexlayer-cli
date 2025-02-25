# Schema Validation Package

This package provides a consolidated approach to validating Nexlayer YAML configurations.

## Core Components

1. **ValidationRegistry**: Maintains a registry of validation functions and rules.
   - Register custom validators
   - Retrieve validators by name
   - Execute validation rules on configuration values

2. **Validator**: Main component for validating configurations.
   - Combines structural validation (JSON Schema) with semantic validation
   - Provides detailed error messages with suggestions

3. **Validation Functions**: Individual validators for specific fields.
   - Pod name validator
   - URL validator
   - Image name validator
   - Volume size validator
   - Environment variable name validator
   - Filename validator

## Usage Example

```go
// Create a default validator with built-in schema
validator := schema.NewDefaultValidator()

// Validate a configuration
config := loadConfig() // Load your configuration
errors := validator.ValidateYAML(config)

// Handle validation errors
if len(errors) > 0 {
    for _, err := range errors {
        fmt.Printf("Error: %s\n", err.Error())
    }
}
```

## Extending with Custom Validators

```go
// Create a custom validator
func validateCustomField(field, value string, ctx *schema.ValidationContext) []schema.ValidationError {
    // Implement custom validation logic
    // Return validation errors if any
}

// Register with the validation registry
registry := schema.NewValidationRegistry()
registry.Register("custom", validateCustomField)

// Create a validator with the custom registry
validator := schema.NewValidator(true, schema.NewStringSchemaSource(schema.SchemaV2))
validator.registry = registry
```

## Notes

This validation package consolidates previously redundant validation code scattered across the codebase, including:

- `pkg/schema/validator.go`
- `pkg/schema/validation.go`
- `pkg/validation/schema.go`
- Other validation-related files

By centralizing validation logic, we ensure consistent validation behavior across the application and reduce code duplication. 