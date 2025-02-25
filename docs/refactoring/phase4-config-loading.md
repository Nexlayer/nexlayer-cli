# Nexlayer CLI Refactoring: Phase 4 - Configuration Loading Consolidation

## Overview

This document outlines the plan for consolidating the configuration loading implementation in the Nexlayer CLI. The goal is to create a unified configuration loading system that provides consistent handling of YAML configurations, eliminating redundancies and improving maintainability.

## Current Issues

1. **Configuration Loading Duplication**:
   - Configuration loading logic is duplicated across multiple commands.
   - `pkg/commands/initcmd/init.go` and `pkg/commands/deploy/deploy.go` both implement their own logic for loading and parsing YAML files.

2. **Inconsistent Validation**:
   - Validation is sometimes performed in commands rather than using the centralized validator.
   - Different commands may apply different validation rules.

3. **Project Detection Duplication**:
   - Project detection logic is duplicated across commands.
   - Caching of detection results is inconsistent.

## Consolidation Plan

### Step 1: Create a Unified Configuration Package

1. **Create `pkg/config` Package**:
   - Implement a comprehensive configuration package.
   - Define clear interfaces for configuration loading and validation.
   - Ensure backward compatibility with existing code.

2. **Implement Configuration Loader**:
   - Create a `ConfigLoader` struct in `pkg/config/loader.go`.
   - Implement methods for loading configurations from files, strings, and bytes.
   - Add support for environment variable substitution.

### Step 2: Implement Project Detection

1. **Create Project Detection Module**:
   - Move project detection logic to `pkg/config/detection.go`.
   - Implement a caching system for detection results.
   - Add support for custom detection rules.

2. **Standardize Project Type Definitions**:
   - Define standard project types in `pkg/config/types.go`.
   - Ensure consistent naming and categorization.
   - Add support for custom project types.

### Step 3: Implement Configuration Generation

1. **Create Configuration Generator**:
   - Implement a `ConfigGenerator` struct in `pkg/config/generator.go`.
   - Add methods for generating configurations based on project types.
   - Support customization of generated configurations.

2. **Add Template Support**:
   - Implement template-based configuration generation.
   - Support for different project templates.
   - Add validation of generated configurations.

### Step 4: Update Command Implementations

1. **Update Command Dependencies**:
   - Modify commands to use the unified configuration package.
   - Ensure consistent configuration handling across commands.
   - Remove redundant configuration loading logic.

2. **Implement Command-Specific Wrappers**:
   - Create command-specific wrappers for configuration operations if needed.
   - These wrappers can add command-specific context or behavior.
   - Ensure these wrappers are thin and focused.

### Step 5: Add Comprehensive Testing

1. **Unit Tests**:
   - Test each configuration operation individually.
   - Test with various configuration formats and sources.
   - Test error handling and recovery.

2. **Integration Tests**:
   - Test configuration loading in the context of command execution.
   - Test with real-world configurations.
   - Test error handling and recovery in integration scenarios.

### Step 6: Update Documentation

1. **Update Configuration Documentation**:
   - Document the unified configuration package.
   - Add examples for common use cases.
   - Document error handling and recovery strategies.

2. **Update Command Documentation**:
   - Update command documentation to reflect the changes.
   - Add examples for using the configuration package in commands.

## Implementation Details

### Configuration Loader Interface

The configuration loader interface will be defined as follows:

```go
// ConfigLoader provides methods for loading and parsing Nexlayer configurations.
type ConfigLoader interface {
    // LoadFromFile loads a configuration from a file.
    LoadFromFile(path string) (*schema.NexlayerYAML, error)
    
    // LoadFromString loads a configuration from a string.
    LoadFromString(content string) (*schema.NexlayerYAML, error)
    
    // LoadFromBytes loads a configuration from a byte slice.
    LoadFromBytes(content []byte) (*schema.NexlayerYAML, error)
    
    // FindConfigFile finds a configuration file in the specified directory.
    FindConfigFile(dir string) (string, error)
    
    // ValidateConfig validates a configuration.
    ValidateConfig(config *schema.NexlayerYAML) []schema.ValidationError
}
```

### Project Detection Interface

The project detection interface will be defined as follows:

```go
// ProjectDetector provides methods for detecting project types.
type ProjectDetector interface {
    // DetectProject detects the project type in the specified directory.
    DetectProject(dir string) (*types.ProjectInfo, error)
    
    // DetectProjectWithCache detects the project type, using cache if available.
    DetectProjectWithCache(dir string, force bool) (*types.ProjectInfo, error)
    
    // RegisterDetector registers a custom detector for a specific project type.
    RegisterDetector(projectType string, detector ProjectTypeDetector)
}

// ProjectTypeDetector is a function that detects a specific project type.
type ProjectTypeDetector func(dir string) (bool, *types.ProjectInfo, error)
```

### Configuration Generator Interface

The configuration generator interface will be defined as follows:

```go
// ConfigGenerator provides methods for generating Nexlayer configurations.
type ConfigGenerator interface {
    // GenerateFromProjectInfo generates a configuration from project info.
    GenerateFromProjectInfo(info *types.ProjectInfo) (*schema.NexlayerYAML, error)
    
    // GenerateFromTemplate generates a configuration from a template.
    GenerateFromTemplate(templateName string, vars map[string]string) (*schema.NexlayerYAML, error)
    
    // RegisterTemplate registers a custom template.
    RegisterTemplate(name string, template string)
}
```

### Implementation Classes

The implementation classes will be defined as follows:

```go
// DefaultConfigLoader implements the ConfigLoader interface.
type DefaultConfigLoader struct {
    validator schema.Validator
}

// DefaultProjectDetector implements the ProjectDetector interface.
type DefaultProjectDetector struct {
    cacheDir string
    detectors map[string]ProjectTypeDetector
}

// DefaultConfigGenerator implements the ConfigGenerator interface.
type DefaultConfigGenerator struct {
    templates map[string]string
}
```

## Testing Strategy

1. **Unit Tests**:
   - Test each configuration operation individually.
   - Test with various configuration formats and sources.
   - Test error handling and recovery.

2. **Integration Tests**:
   - Test configuration loading in the context of command execution.
   - Test with real-world configurations.
   - Test error handling and recovery in integration scenarios.

3. **Benchmarks**:
   - Benchmark configuration loading performance.
   - Compare performance before and after refactoring.
   - Identify and address performance bottlenecks.

## Migration Strategy

1. **Phase 1: Implement Unified Package**:
   - Create the unified configuration package.
   - Implement the core interfaces and classes.
   - Add comprehensive tests.

2. **Phase 2: Update Command Dependencies**:
   - Update commands to use the unified package.
   - Ensure consistent configuration handling.
   - Add tests for command-configuration integration.

3. **Phase 3: Enhance Functionality**:
   - Add support for templates and custom project types.
   - Implement caching and performance optimizations.
   - Ensure backward compatibility.

4. **Phase 4: Deprecate Old Code**:
   - Mark old configuration loading code as deprecated.
   - Add forwarding functions for backward compatibility.
   - Update documentation to reflect the changes.

## Timeline

- **Week 1**: Implement unified configuration package and core interfaces.
- **Week 2**: Implement project detection and configuration generation.
- **Week 3**: Update command dependencies and add tests.
- **Week 4**: Enhance functionality, update documentation, and finalize.

## Conclusion

Consolidating the configuration loading implementation will create a unified system for handling Nexlayer configurations. This will improve maintainability, reduce the risk of inconsistencies, and provide a better developer experience. The migration strategy ensures backward compatibility while moving towards a more unified and robust system. 