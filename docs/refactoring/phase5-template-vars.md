# Nexlayer CLI Refactoring: Phase 5 - Template Variable Processing Consolidation

## Overview

This document outlines the plan for consolidating the template variable processing implementation in the Nexlayer CLI. The goal is to create a unified system for handling template variables, eliminating redundancies and improving maintainability.

## Current Issues

1. **Variable Substitution Duplication**:
   - Template variable substitution is implemented in multiple places.
   - `pkg/schema/generator.go` and other files contain hardcoded template variables.

2. **Environment Variable Handling Duplication**:
   - Environment variable handling is duplicated across packages.
   - Inconsistent handling of default values and validation.

3. **Lack of Centralized Registry**:
   - No centralized registry of supported template variables.
   - Adding new template variables requires updates in multiple places.

## Consolidation Plan

### Step 1: Create a Unified Variables Package

1. **Enhance `pkg/vars` Package**:
   - Implement a comprehensive variables package.
   - Define clear interfaces for variable processing and substitution.
   - Ensure backward compatibility with existing code.

2. **Implement Variable Processor**:
   - Create a `VariableProcessor` struct in `pkg/vars/processor.go`.
   - Implement methods for variable substitution and validation.
   - Add support for different variable formats and sources.

### Step 2: Implement Variable Registry

1. **Create Variable Registry**:
   - Implement a registry of supported template variables in `pkg/vars/registry.go`.
   - Define variable categories (e.g., system, user, environment).
   - Add support for variable documentation and validation.

2. **Standardize Variable Formats**:
   - Define standard variable formats (e.g., `<% VAR_NAME %>`, `${VAR_NAME}`).
   - Ensure consistent handling of variable formats.
   - Add support for default values and transformations.

### Step 3: Implement Environment Variable Integration

1. **Create Environment Variable Handler**:
   - Implement an environment variable handler in `pkg/vars/env.go`.
   - Add support for environment variable validation and transformation.
   - Implement fallback mechanisms for missing environment variables.

2. **Add Secret Management**:
   - Implement secure handling of sensitive variables.
   - Add support for secret storage and retrieval.
   - Ensure secure logging of sensitive variables.

### Step 4: Update Dependent Code

1. **Update Schema Generator**:
   - Modify `pkg/schema/generator.go` to use the unified variables package.
   - Remove hardcoded template variables.
   - Ensure consistent variable handling.

2. **Update Command Implementations**:
   - Modify commands to use the unified variables package.
   - Ensure consistent variable handling across commands.
   - Remove redundant variable processing logic.

### Step 5: Add Comprehensive Testing

1. **Unit Tests**:
   - Test each variable processing operation individually.
   - Test with various variable formats and sources.
   - Test error handling and recovery.

2. **Integration Tests**:
   - Test variable processing in the context of command execution.
   - Test with real-world configurations.
   - Test error handling and recovery in integration scenarios.

### Step 6: Update Documentation

1. **Update Variable Documentation**:
   - Document the unified variables package.
   - Add examples for common use cases.
   - Document supported variables and their usage.

2. **Update Command Documentation**:
   - Update command documentation to reflect the changes.
   - Add examples for using variables in commands.

## Implementation Details

### Variable Processor Interface

The variable processor interface will be defined as follows:

```go
// VariableProcessor provides methods for processing template variables.
type VariableProcessor interface {
    // Process substitutes variables in the input string.
    Process(input string) (string, error)
    
    // ProcessBytes substitutes variables in the input byte slice.
    ProcessBytes(input []byte) ([]byte, error)
    
    // ProcessMap substitutes variables in the input map.
    ProcessMap(input map[string]string) (map[string]string, error)
    
    // ProcessYAML substitutes variables in the input YAML.
    ProcessYAML(input []byte) ([]byte, error)
    
    // SetVariable sets a variable value.
    SetVariable(name string, value string)
    
    // GetVariable gets a variable value.
    GetVariable(name string) (string, bool)
}
```

### Variable Registry Interface

The variable registry interface will be defined as follows:

```go
// VariableRegistry provides methods for managing template variables.
type VariableRegistry interface {
    // Register registers a variable with the registry.
    Register(name string, category string, description string, validator VariableValidator)
    
    // Get gets a variable from the registry.
    Get(name string) (Variable, bool)
    
    // List lists all variables in the registry.
    List() []Variable
    
    // ListByCategory lists variables by category.
    ListByCategory(category string) []Variable
}

// Variable represents a template variable.
type Variable struct {
    Name        string
    Category    string
    Description string
    Validator   VariableValidator
}

// VariableValidator validates a variable value.
type VariableValidator func(value string) error
```

### Environment Variable Handler Interface

The environment variable handler interface will be defined as follows:

```go
// EnvHandler provides methods for handling environment variables.
type EnvHandler interface {
    // Get gets an environment variable.
    Get(name string) (string, bool)
    
    // Set sets an environment variable.
    Set(name string, value string) error
    
    // GetWithDefault gets an environment variable with a default value.
    GetWithDefault(name string, defaultValue string) string
    
    // GetRequired gets a required environment variable.
    GetRequired(name string) (string, error)
    
    // GetWithValidator gets an environment variable with validation.
    GetWithValidator(name string, validator VariableValidator) (string, error)
}
```

### Implementation Classes

The implementation classes will be defined as follows:

```go
// DefaultVariableProcessor implements the VariableProcessor interface.
type DefaultVariableProcessor struct {
    registry VariableRegistry
    env      EnvHandler
    vars     map[string]string
}

// DefaultVariableRegistry implements the VariableRegistry interface.
type DefaultVariableRegistry struct {
    variables map[string]Variable
}

// DefaultEnvHandler implements the EnvHandler interface.
type DefaultEnvHandler struct {
    prefix string
}
```

## Variable Categories and Examples

The variable registry will include the following categories:

1. **System Variables**:
   - `REGISTRY`: Container registry URL
   - `URL`: Application URL
   - `NAMESPACE`: Kubernetes namespace
   - `VERSION`: Application version

2. **User Variables**:
   - `USER_NAME`: User name
   - `USER_EMAIL`: User email
   - `USER_TOKEN`: User authentication token

3. **Environment Variables**:
   - `NODE_ENV`: Node.js environment
   - `PORT`: Application port
   - `DATABASE_URL`: Database connection URL

4. **Secret Variables**:
   - `API_KEY`: API key
   - `DATABASE_PASSWORD`: Database password
   - `JWT_SECRET`: JWT secret

## Testing Strategy

1. **Unit Tests**:
   - Test each variable processing operation individually.
   - Test with various variable formats and sources.
   - Test error handling and recovery.

2. **Integration Tests**:
   - Test variable processing in the context of command execution.
   - Test with real-world configurations.
   - Test error handling and recovery in integration scenarios.

3. **Security Tests**:
   - Test secure handling of sensitive variables.
   - Test variable validation and sanitization.
   - Test protection against injection attacks.

## Migration Strategy

1. **Phase 1: Implement Unified Package**:
   - Enhance the `pkg/vars` package.
   - Implement the core interfaces and classes.
   - Add comprehensive tests.

2. **Phase 2: Update Dependent Code**:
   - Update schema generator to use the unified package.
   - Update commands to use the unified package.
   - Ensure consistent variable handling.

3. **Phase 3: Enhance Functionality**:
   - Add support for variable categories and documentation.
   - Implement secure handling of sensitive variables.
   - Ensure backward compatibility.

4. **Phase 4: Deprecate Old Code**:
   - Mark old variable processing code as deprecated.
   - Add forwarding functions for backward compatibility.
   - Update documentation to reflect the changes.

## Timeline

- **Week 1**: Enhance `pkg/vars` package and implement core interfaces.
- **Week 2**: Implement variable registry and environment variable integration.
- **Week 3**: Update dependent code and add tests.
- **Week 4**: Enhance functionality, update documentation, and finalize.

## Conclusion

Consolidating the template variable processing implementation will create a unified system for handling template variables. This will improve maintainability, reduce the risk of inconsistencies, and provide a better developer experience. The migration strategy ensures backward compatibility while moving towards a more unified and robust system. 