# Nexlayer CLI Refactoring: Phase 3 - API Client Consolidation

## Overview

This document outlines the plan for consolidating the API client implementation in the Nexlayer CLI. The goal is to create a unified API client that provides a consistent interface for interacting with the Nexlayer API, eliminating redundancies and improving maintainability.

## Current Issues

1. **Interface Duplication**:
   - Multiple interfaces defined in `pkg/core/api/client.go`: `ClientAPI`, `APIClient`, and `APIClientForCommands`.
   - These interfaces have overlapping methods and responsibilities.

2. **API Client Usage**:
   - Commands directly use the API client, sometimes with redundant error handling.
   - Some commands implement their own API interaction logic instead of using the centralized client.

3. **Error Handling Duplication**:
   - Error handling for API responses is duplicated across different command implementations.
   - Inconsistent error messages and recovery strategies.

## Consolidation Plan

### Step 1: Define a Unified API Client Interface

1. **Create a Single Interface**:
   - Define a comprehensive `APIClient` interface in `pkg/core/api/client.go`.
   - Ensure the interface covers all required API operations.
   - Add clear documentation for each method.

2. **Organize Methods by Domain**:
   - Group methods by domain (e.g., deployments, feedback, domains).
   - Ensure consistent naming and parameter conventions.

### Step 2: Implement the Unified API Client

1. **Enhance the Client Implementation**:
   - Update the `Client` struct to implement the unified interface.
   - Ensure consistent error handling across all methods.
   - Add comprehensive logging and telemetry.

2. **Create Domain-Specific Clients**:
   - Implement domain-specific clients (e.g., `DeploymentClient`, `FeedbackClient`).
   - Compose these clients within the main `Client` struct.
   - This allows for better separation of concerns and testability.

### Step 3: Standardize Error Handling

1. **Create API Error Types**:
   - Define a comprehensive set of API error types in `pkg/core/api/errors.go`.
   - Implement error categorization (e.g., authentication, validation, server).
   - Add support for error recovery and retry strategies.

2. **Implement Consistent Error Handling**:
   - Ensure all API client methods handle errors consistently.
   - Add context to errors for better debugging.
   - Implement retry logic for transient errors.

### Step 4: Update Command Implementations

1. **Update Command Dependencies**:
   - Modify commands to use the unified API client interface.
   - Ensure consistent error handling across commands.
   - Remove redundant API interaction logic.

2. **Implement Command-Specific Wrappers**:
   - Create command-specific wrappers for API client methods if needed.
   - These wrappers can add command-specific context or behavior.
   - Ensure these wrappers are thin and focused.

### Step 5: Add Comprehensive Testing

1. **Unit Tests**:
   - Test each API client method individually.
   - Use mocks to simulate API responses.
   - Test error handling and recovery.

2. **Integration Tests**:
   - Test API client in the context of command execution.
   - Use a mock server to simulate the Nexlayer API.
   - Test with real-world scenarios.

### Step 6: Update Documentation

1. **Update API Client Documentation**:
   - Document the unified API client interface.
   - Add examples for common use cases.
   - Document error handling and recovery strategies.

2. **Update Command Documentation**:
   - Update command documentation to reflect the changes.
   - Add examples for using the API client in commands.

## Implementation Details

### Unified API Client Interface

The unified API client interface will be defined as follows:

```go
// APIClient defines the interface for interacting with the Nexlayer API.
type APIClient interface {
    // Deployment operations
    StartDeployment(ctx context.Context, appID string, configPath string) (*schema.APIResponse[schema.DeploymentResponse], error)
    GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*schema.APIResponse[schema.Deployment], error)
    ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error)
    GetLogs(ctx context.Context, namespace string, appID string, follow bool, tail int) ([]string, error)
    
    // Domain operations
    SaveCustomDomain(ctx context.Context, appID string, domain string) (*schema.APIResponse[struct{}], error)
    
    // Feedback operations
    SendFeedback(ctx context.Context, text string) error
    
    // Authentication operations
    Login(ctx context.Context, username, password string) (*schema.APIResponse[schema.AuthResponse], error)
    Logout(ctx context.Context) error
    RefreshToken(ctx context.Context) (*schema.APIResponse[schema.AuthResponse], error)
}
```

### Domain-Specific Clients

The domain-specific clients will be implemented as follows:

```go
// DeploymentClient handles deployment-related API operations.
type DeploymentClient struct {
    baseURL    string
    httpClient *http.Client
    token      string
}

// FeedbackClient handles feedback-related API operations.
type FeedbackClient struct {
    baseURL    string
    httpClient *http.Client
    token      string
}

// DomainClient handles domain-related API operations.
type DomainClient struct {
    baseURL    string
    httpClient *http.Client
    token      string
}

// Client implements the APIClient interface by composing domain-specific clients.
type Client struct {
    baseURL    string
    httpClient *http.Client
    token      string
    
    deployment *DeploymentClient
    feedback   *FeedbackClient
    domain     *DomainClient
}
```

### Standardized Error Handling

The error handling will be standardized as follows:

```go
// APIError represents an error returned by the Nexlayer API.
type APIError struct {
    StatusCode int    `json:"status_code"`
    Message    string `json:"message"`
    Code       string `json:"code"`
    Details    any    `json:"details,omitempty"`
}

// Error categories
const (
    ErrorCategoryAuth      = "auth"
    ErrorCategoryValidation = "validation"
    ErrorCategoryServer    = "server"
    ErrorCategoryNetwork   = "network"
)

// Error codes
const (
    ErrorCodeUnauthorized      = "unauthorized"
    ErrorCodeInvalidRequest    = "invalid_request"
    ErrorCodeResourceNotFound  = "resource_not_found"
    ErrorCodeInternalServer    = "internal_server"
    ErrorCodeNetworkError      = "network_error"
)

// IsAuthError returns true if the error is an authentication error.
func IsAuthError(err error) bool {
    var apiErr *APIError
    if errors.As(err, &apiErr) {
        return apiErr.Code == ErrorCodeUnauthorized
    }
    return false
}

// IsValidationError returns true if the error is a validation error.
func IsValidationError(err error) bool {
    var apiErr *APIError
    if errors.As(err, &apiErr) {
        return apiErr.StatusCode == http.StatusBadRequest
    }
    return false
}

// IsRetryableError returns true if the error is retryable.
func IsRetryableError(err error) bool {
    var apiErr *APIError
    if errors.As(err, &apiErr) {
        return apiErr.StatusCode >= 500 || apiErr.StatusCode == http.StatusTooManyRequests
    }
    return false
}
```

## Testing Strategy

1. **Unit Tests**:
   - Test each API client method individually.
   - Test error handling and recovery.
   - Test with various API responses.

2. **Integration Tests**:
   - Test API client in the context of command execution.
   - Test with real-world scenarios.
   - Test error handling and recovery in integration scenarios.

3. **Mock Server**:
   - Implement a mock server for testing.
   - Simulate various API responses and errors.
   - Test retry logic and error recovery.

## Migration Strategy

1. **Phase 1: Implement Unified Interface**:
   - Define the unified API client interface.
   - Implement the interface in the existing client.
   - Add comprehensive tests.

2. **Phase 2: Update Command Dependencies**:
   - Update commands to use the unified interface.
   - Ensure consistent error handling.
   - Add tests for command-API client integration.

3. **Phase 3: Implement Domain-Specific Clients**:
   - Refactor the client implementation to use domain-specific clients.
   - Update tests to reflect the changes.
   - Ensure backward compatibility.

4. **Phase 4: Deprecate Old Interfaces**:
   - Mark old interfaces as deprecated.
   - Add forwarding methods for backward compatibility.
   - Update documentation to reflect the changes.

## Timeline

- **Week 1**: Define unified interface and update client implementation.
- **Week 2**: Standardize error handling and update command dependencies.
- **Week 3**: Implement domain-specific clients and add tests.
- **Week 4**: Update documentation and finalize.

## Conclusion

Consolidating the API client implementation will create a unified interface for interacting with the Nexlayer API. This will improve maintainability, reduce the risk of inconsistencies, and provide a better developer experience. The migration strategy ensures backward compatibility while moving towards a more unified and robust system. 