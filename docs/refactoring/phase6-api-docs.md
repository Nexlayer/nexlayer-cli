# Nexlayer CLI Refactoring: Phase 6 - API Reference Documentation Consolidation

## Overview

This document outlines the plan for consolidating the API reference documentation in the Nexlayer CLI. The goal is to create a single source of truth for API documentation, eliminating redundancies and improving maintainability.

## Current Issues

1. **Documentation Duplication**:
   - API endpoint documentation is duplicated between code comments and documentation files.
   - Different parts of the codebase may describe API endpoints differently.

2. **Inconsistent Documentation**:
   - Documentation format and level of detail vary across the codebase.
   - Some endpoints may be documented in multiple places with different information.

3. **Lack of Automated Generation**:
   - Documentation is manually maintained, leading to inconsistencies.
   - Changes to API endpoints may not be reflected in documentation.

## Consolidation Plan

### Step 1: Define Documentation Standards

1. **Create Documentation Guidelines**:
   - Define a standard format for API documentation.
   - Specify required information for each endpoint.
   - Establish conventions for examples and error handling.

2. **Select Documentation Tools**:
   - Choose tools for generating API documentation.
   - Consider OpenAPI/Swagger for API specification.
   - Select tools for generating documentation from code comments.

### Step 2: Implement OpenAPI Specification

1. **Create OpenAPI Specification**:
   - Define the Nexlayer API using OpenAPI 3.0.
   - Include all endpoints, request/response schemas, and authentication.
   - Add examples and error responses.

2. **Validate OpenAPI Specification**:
   - Ensure the specification is valid and complete.
   - Verify that all endpoints are documented.
   - Check for consistency and accuracy.

### Step 3: Update Code Comments

1. **Enhance API Client Comments**:
   - Update comments in `pkg/core/api/client.go` to follow the documentation standard.
   - Ensure comments include all required information.
   - Add examples and error handling details.

2. **Add OpenAPI Annotations**:
   - Add OpenAPI annotations to code comments.
   - Ensure annotations match the OpenAPI specification.
   - Add validation for annotations during build.

### Step 4: Implement Documentation Generation

1. **Create Documentation Generator**:
   - Implement a tool for generating documentation from OpenAPI specification.
   - Add support for generating Markdown, HTML, and other formats.
   - Ensure generated documentation is consistent and complete.

2. **Integrate with Build Process**:
   - Add documentation generation to the build process.
   - Validate documentation during CI/CD.
   - Ensure documentation is always up-to-date.

### Step 5: Update Existing Documentation

1. **Update Command Documentation**:
   - Update command documentation to reference the API documentation.
   - Ensure consistent terminology and examples.
   - Add links to the API documentation.

2. **Create API Reference Guide**:
   - Create a comprehensive API reference guide.
   - Include examples, error handling, and authentication.
   - Add troubleshooting and best practices.

### Step 6: Implement Documentation Testing

1. **Add Documentation Tests**:
   - Implement tests for documentation accuracy.
   - Verify that examples work as expected.
   - Check for broken links and inconsistencies.

2. **Add API Contract Tests**:
   - Implement tests to verify that the API implementation matches the documentation.
   - Test error handling and edge cases.
   - Ensure backward compatibility.

## Implementation Details

### OpenAPI Specification Structure

The OpenAPI specification will be structured as follows:

```yaml
openapi: 3.0.0
info:
  title: Nexlayer API
  description: API for deploying and managing Nexlayer applications
  version: 1.0.0
servers:
  - url: https://api.nexlayer.dev
    description: Production server
  - url: https://api.staging.nexlayer.dev
    description: Staging server
paths:
  /startUserDeployment/{applicationID}:
    post:
      summary: Start a new deployment
      description: Starts a new deployment using a YAML configuration file
      parameters:
        - name: applicationID
          in: path
          description: Optional application ID
          required: false
          schema:
            type: string
      requestBody:
        content:
          text/x-yaml:
            schema:
              $ref: '#/components/schemas/NexlayerYAML'
      responses:
        '200':
          description: Deployment started successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/DeploymentResponse'
        '400':
          description: Invalid request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: Unauthorized
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
  # ... other endpoints ...
components:
  schemas:
    NexlayerYAML:
      type: object
      properties:
        application:
          $ref: '#/components/schemas/Application'
      required:
        - application
    # ... other schemas ...
```

### Code Comment Format

The code comments will follow this format:

```go
// StartDeployment starts a new deployment using a YAML configuration file.
//
// @Summary Start a new deployment
// @Description Starts a new deployment using a YAML configuration file
// @Tags deployment
// @Accept text/x-yaml
// @Produce json
// @Param applicationID path string false "Optional application ID"
// @Param config body NexlayerYAML true "YAML configuration"
// @Success 200 {object} DeploymentResponse "Deployment started successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /startUserDeployment/{applicationID} [post]
func (c *Client) StartDeployment(ctx context.Context, appID string, yamlFile string) (*schema.APIResponse[schema.DeploymentResponse], error) {
    // ... implementation ...
}
```

### Documentation Generator

The documentation generator will be implemented as follows:

```go
// Generator generates API documentation from OpenAPI specification.
type Generator struct {
    specPath string
    outputDir string
    formats []string
}

// NewGenerator creates a new documentation generator.
func NewGenerator(specPath string, outputDir string, formats []string) *Generator {
    return &Generator{
        specPath: specPath,
        outputDir: outputDir,
        formats: formats,
    }
}

// Generate generates documentation in the specified formats.
func (g *Generator) Generate() error {
    // ... implementation ...
}
```

## Documentation Formats

The documentation will be generated in the following formats:

1. **Markdown**:
   - For GitHub and other Markdown-based platforms.
   - Includes examples and code snippets.
   - Suitable for version control.

2. **HTML**:
   - For web-based documentation.
   - Includes interactive examples.
   - Suitable for publishing on the Nexlayer website.

3. **OpenAPI UI**:
   - Interactive API documentation.
   - Includes request/response examples.
   - Suitable for developers exploring the API.

## Testing Strategy

1. **Documentation Accuracy Tests**:
   - Verify that documentation matches the implementation.
   - Check for missing or outdated information.
   - Ensure examples work as expected.

2. **API Contract Tests**:
   - Test that the API implementation matches the documentation.
   - Verify error handling and edge cases.
   - Ensure backward compatibility.

3. **Documentation Generation Tests**:
   - Test that documentation is generated correctly.
   - Verify that all formats are generated.
   - Check for broken links and inconsistencies.

## Migration Strategy

1. **Phase 1: Define Standards and Create OpenAPI Specification**:
   - Define documentation standards.
   - Create the OpenAPI specification.
   - Validate the specification.

2. **Phase 2: Update Code Comments and Implement Generation**:
   - Update code comments to follow the standard.
   - Implement documentation generation.
   - Integrate with the build process.

3. **Phase 3: Update Existing Documentation and Add Tests**:
   - Update command documentation.
   - Create the API reference guide.
   - Add documentation tests.

4. **Phase 4: Deprecate Old Documentation**:
   - Mark old documentation as deprecated.
   - Add redirects to the new documentation.
   - Remove old documentation in a future release.

## Timeline

- **Week 1**: Define documentation standards and create OpenAPI specification.
- **Week 2**: Update code comments and implement documentation generation.
- **Week 3**: Update existing documentation and add tests.
- **Week 4**: Finalize documentation and integrate with the build process.

## Conclusion

Consolidating the API reference documentation will create a single source of truth for API documentation. This will improve maintainability, reduce the risk of inconsistencies, and provide a better developer experience. The migration strategy ensures backward compatibility while moving towards a more unified and robust documentation system. 