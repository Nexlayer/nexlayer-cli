# Phase 3: API Client Consolidation

**Status: 🔄 IN PROGRESS**  
**Duration: 2 days**

## Overview

Phase 3 focuses on consolidating the API client implementations across the codebase to eliminate redundancy and improve maintainability. Currently, there are multiple implementations of API client functionality in different packages, leading to code duplication and potential inconsistencies.

## Current Issues

- Duplicate API client methods in:
  - `pkg/core/api/client.go`
  - `pkg/commands/deploy/api.go`
  - `pkg/commands/feedback/api.go`
- Inconsistent error handling across different implementations
- Lack of standardized approach to API client usage
- No centralized logging or monitoring of API calls

## Consolidation Plan

- **Enhance Core API Client**
  - Ensure `pkg/core/api/client.go` contains all needed functionality
  - Add missing methods from command-specific implementations
  - Improve error handling and logging

- **Create Command-Specific Wrappers**
  - Replace direct API implementations in commands with wrappers
  - Use dependency injection for the API client

- **Implement Middleware Pattern**
  - Create middleware for common functionality (logging, error handling)
  - Allow commands to add specific middleware as needed

## Implementation Steps

### 1. Enhance Core API Client

```go
// In pkg/core/api/client.go
package api

// Client is the interface for interacting with the Nexlayer API
type Client interface {
    // Core methods
    CreateDeployment(ctx context.Context, request *DeploymentRequest) (*DeploymentResponse, error)
    GetDeployment(ctx context.Context, id string) (*DeploymentResponse, error)
    ListDeployments(ctx context.Context) ([]DeploymentResponse, error)
    DeleteDeployment(ctx context.Context, id string) error
    
    // Additional methods from command-specific implementations
    SendFeedback(ctx context.Context, feedback *FeedbackRequest) error
    GetStatus(ctx context.Context) (*StatusResponse, error)
    
    // Configuration methods
    WithToken(token string) Client
    WithBaseURL(url string) Client
    WithMiddleware(middleware Middleware) Client
}

// DefaultClient is the default implementation of the Client interface
type DefaultClient struct {
    baseURL    string
    httpClient *http.Client
    token      string
    middlewares []Middleware
}

// NewClient creates a new API client with the provided options
func NewClient(options ...Option) Client {
    // Implementation details...
}
```

### 2. Create Command-Specific Wrappers

```go
// In pkg/commands/deploy/deploy.go
package deploy

import (
    "github.com/Nexlayer/nexlayer-cli/pkg/core/api"
    "github.com/spf13/cobra"
)

// NewCommand creates a new deploy command
func NewCommand(apiClient api.Client) *cobra.Command {
    // Use the injected client instead of creating a new one
}

// In cmd/root.go
func initCommands() {
    // Create API client
    apiClient := api.NewClient(
        api.WithToken(config.GetToken()),
        api.WithBaseURL(config.GetAPIURL()),
        api.WithMiddleware(api.LoggingMiddleware(logger)),
    )
    
    // Initialize commands with the API client
    rootCmd.AddCommand(deploy.NewCommand(apiClient))
    rootCmd.AddCommand(feedback.NewCommand(apiClient))
    // Other commands...
}
```

### 3. Implement Middleware Pattern

```go
// In pkg/core/api/middleware.go
package api

// Middleware represents a function that wraps an HTTP handler
type Middleware func(http.RoundTripper) http.RoundTripper

// LoggingMiddleware logs information about API requests and responses
func LoggingMiddleware(logger Logger) Middleware {
    // Implementation details...
}

// RetryMiddleware implements retry logic for failed API requests
func RetryMiddleware(maxRetries int, backoff BackoffStrategy) Middleware {
    // Implementation details...
}
```

## Tasks

- [ ] Audit existing API client implementations
  - [ ] Identify all API endpoints used across the codebase
  - [ ] Document parameters, return types, and behaviors
  - [ ] Identify common patterns and inconsistencies

- [ ] Enhance core API client
  - [ ] Define comprehensive Client interface
  - [ ] Implement DefaultClient with all required methods
  - [ ] Add support for configuration options
  - [ ] Implement robust error handling

- [ ] Create middleware pattern
  - [ ] Define Middleware interface
  - [ ] Implement common middleware (logging, retry, etc.)
  - [ ] Add middleware support to DefaultClient

- [ ] Update command implementations
  - [ ] Refactor commands to use dependency injection
  - [ ] Replace direct API implementations with core client
  - [ ] Update error handling to use standardized approach

- [ ] Test and validate changes
  - [ ] Write unit tests for core API client
  - [ ] Write integration tests for API functionality
  - [ ] Verify commands work with the new API client

## Deliverables

- Enhanced core API client in `pkg/core/api/client.go`
- Middleware implementation for cross-cutting concerns
- Updated command implementations using the new API client
- Comprehensive test suite for API client functionality
- Updated documentation and examples

## Acceptance Criteria

- All command implementations use the centralized API client
- No duplicate API client implementations in the codebase
- Consistent error handling and logging across all API calls
- Comprehensive test coverage for API client functionality
- Documentation is updated to reflect the new API client usage

## Risks and Mitigation

- **Risk**: Breaking changes to command implementations
  - **Mitigation**: Thorough testing and gradual migration

- **Risk**: Performance impact of additional middleware
  - **Mitigation**: Performance testing and optimizations

- **Risk**: Incomplete API endpoint coverage
  - **Mitigation**: Comprehensive audit of current API usage

## Next Steps

Proceed to [Phase 4: Configuration Loading Consolidation](./Phase4-ConfigLoading.md) after completion 