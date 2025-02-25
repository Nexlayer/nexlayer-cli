# Nexlayer YAML Schema: The Definitive Source of Truth

This package serves as the **DEFINITIVE**, single source of truth for the Nexlayer YAML schema. All other packages in the codebase that need to work with Nexlayer configuration should use the types defined here. This ensures consistency across the entire platform and prevents schema drift between different components.

## Role of this Package

The `pkg/schema` package has several important responsibilities:

1. **Type Definitions**: Centralized definition of all configuration types used in Nexlayer
2. **Validation Logic**: Comprehensive validation of configuration structure and values
3. **Schema Generation**: Utilities for generating valid configuration
4. **Consistency Enforcement**: Ensuring all parts of the codebase use the same schema

**DO NOT** define duplicate schema types in other packages. Instead, import and use the types defined here.

## Integration With Other Packages

Other packages in the codebase should interact with this package as follows:

- **Commands Package**: Should use the schema types directly for validation and generation
- **Template Package**: Should use schema types or provide explicit conversion to/from schema types
- **Validation Package**: Should be deprecated in favor of validation logic in this package
- **Core Types**: Should be aligned with schema types or provide explicit conversion

## Schema Overview

The Nexlayer YAML schema defines the structure for deploying applications on the Nexlayer platform. It specifies the required and optional fields, validation rules, and relationships between components.

## Package Structure

```
pkg/schema/
├── README.md           # This documentation (single source of truth)
├── types.go            # Core schema types and structures
├── validation.go       # Schema validation logic
├── errors.go           # Validation error handling
├── jsonschema.go       # JSON Schema utilities
├── generator.go        # Schema generation utilities
└── examples/           # Example templates
    └── standard.go     # Standard schema examples
```

## Schema Structure

### Top-Level Structure

```yaml
application:
  name: The name of the deployment
  url: Permanent domain URL (optional). No need to add this key if this is not going to be a permanent deployment.
  registryLogin:
    registry: The registry where private images are stored.
    username: Registry username.
    personalAccessToken: Read-only registry Personal Access Token.
  pods:
    - name: Pod name (must start with a lowercase letter and can include only alphanumeric characters, '-', '.')
      path: Path to render pod at (such as '/' for frontend). Only required for forward-facing pods.
      image: Docker image for the pod. 
        # For private images, use the following schema exactly as shown: '<% REGISTRY %>/some/path/image:tag'.
        # Images will be tagged as private if they include '<% REGISTRY %>', which will be replaced with the registry specified above.
      entrypoint: command to replace ENTRYPOINT of image
      command: command to replace CMD of image
      volumes:
        # Array of volumes to be mounted for this pod. Example:
        - name: Name of the volume (lowercase, alphanumeric, '-')
          size: 1Gi  # Required: Volume size (e.g., "1Gi", "500Mi").
          mountPath: /var/some/directory  # Required: Must start with '/'.
      secrets:
        # Array of secret files for this pod. Example:
        - name: Secret name (lowercase, alphanumeric, '-')
          data: Raw text or Base64-encoded string for the secret (e.g., JSON files should be encoded).
          mountPath: Mount path where the secret file will be stored (must start with '/').
          fileName: Name of the secret file (e.g., "secret-file.txt"). 
            # This will be available at "/var/secrets/my-secret-volume/secret-file.txt".
      vars:
        # Array of environment variables for this pod. Example:
        - key: ENV_VAR_NAME
          value: Value of the environment variable.
        # Can use <pod-name>.pod to reference other pods dynamically. Example:
        - key: API_URL
          value: http://express.pod:3000  # Where 'express' is the name of another pod.
        # Can use <% URL %> to reference the deployment's base URL dynamically. Example:
        - key: API_URL
          value: <% URL %>/api
    servicePorts:
      # Array of ports to expose for this pod. Example:
      - 3000  # Exposing port 3000.
  entrypoint: Custom container entrypoint (optional).
  command: Custom container command (optional).
```

## Detailed Field Specifications

### Application Configuration

| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|-----------------|
| `name` | string | Yes | Unique application name | Must be lowercase alphanumeric with optional hyphens or dots |
| `url` | string | No | Custom domain | Must be a valid URL format |
| `registryLogin` | object | No* | Private registry credentials | Required if using private images |

*Required if using private images with `<% REGISTRY %>` template variable.

### Registry Login Configuration

| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|-----------------|
| `registry` | string | Yes | Registry hostname | e.g., "docker.io/my-org" |
| `username` | string | Yes | Registry username | |
| `personalAccessToken` | string | Yes | Read-only PAT | |

### Pod Configuration

| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|-----------------|
| `name` | string | Yes | Pod name | Must start with lowercase letter, use alphanumeric characters, -, or . |
| `type` | string | No | Pod type | e.g., nextjs, react, node, python |
| `path` | string | No* | URL path for routing | Must start with "/" for forward-facing pods |
| `image` | string | Yes | Container image | Public: standard Docker format (e.g., nginx:latest)<br>Private: `<% REGISTRY %>/path/image:tag` |
| `entrypoint` | string | No | Container entrypoint | Overrides default entrypoint |
| `command` | string | No | Container command | Overrides default command |
| `annotations` | map[string]string | No | Custom Kubernetes annotations | |

*Required for forward-facing pods that need to be accessible via HTTP.

### Service Ports Configuration

| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|-----------------|
| `port` | integer | Yes | External port number | Must be between 1-65535 |
| `targetPort` | integer | Yes | Internal container port | Must be between 1-65535 |
| `name` | string | Yes | Port name | Must be lowercase alphanumeric with optional hyphens |
| `protocol` | string | No | Network protocol | Default is TCP, can be UDP |

### Environment Variables Configuration

| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|-----------------|
| `key` | string | Yes | Variable name | Must be a valid environment variable name |
| `value` | string | Yes | Variable value | Supports template variables:<br>- `<pod-name>.pod` for pod references<br>- `<% URL %>` for application URL |

### Volume Configuration

| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|-----------------|
| `name` | string | Yes | Volume name | Must be lowercase alphanumeric with optional hyphens |
| `path` | string | Yes | Mount path | Must start with "/" |
| `size` | string | Yes | Volume size | Must use valid units (Ki, Mi, Gi, Ti) |
| `type` | string | No | Volume type | Default is "standard" |
| `readOnly` | boolean | No | Read-only flag | Default is false |

### Secret Configuration

| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|-----------------|
| `name` | string | Yes | Secret name | Must be lowercase alphanumeric with optional hyphens |
| `data` | string | Yes | Secret data | Raw text or Base64-encoded string |
| `path` | string | Yes | Mount path | Must start with "/" |
| `fileName` | string | Yes | File name | Must be a valid filename |

## Validation Rules

### General Rules

- **Pod Names**: Must start with a lowercase letter, use alphanumeric characters, -, or . (e.g., "web-app").
- **Paths**: Must start with "/" for forward-facing pods (e.g., "/api").
- **Images**: 
  - Public: Standard Docker format (e.g., "postgres:latest").
  - Private: Use `<% REGISTRY %>/repository/image:tag` with registryLogin.
- **Service Ports**: Must be integers between 1 and 65535.
- **Volumes**: Paths must start with "/", sizes use Ki, Mi, Gi, Ti (e.g., "5Gi").
- **Secrets**: Mount paths must start with "/"; JSON data should be Base64-encoded.

### Template Variables

The schema supports the following template variables:

- `<% REGISTRY %>`: Replaced with the configured private registry URL.
- `<% URL %>`: Replaced with the application's public URL.
- `<pod-name>.pod`: Replaced with the internal DNS name for pod-to-pod communication.

## Complete Example

```yaml
application:
  name: "ai-fullstack-demo"
  url: "ai-demo.example.com"
  registryLogin:
    registry: "ghcr.io/myorg"
    username: "username"
    personalAccessToken: "pat_token"
  pods:
    - name: "frontend"
      type: "nextjs"
      path: "/"
      image: "<% REGISTRY %>/frontend:latest"
      servicePorts:
        - port: 3000
          targetPort: 3000
          name: "http"
      vars:
        - key: "API_URL"
          value: "http://backend.pod:8000"
        - key: "PUBLIC_URL"
          value: "<% URL %>"
      
    - name: "backend"
      type: "fastapi"
      path: "/api"
      image: "<% REGISTRY %>/backend:latest"
      servicePorts:
        - port: 8000
          targetPort: 8000
          name: "api"
      vars:
        - key: "DATABASE_URL"
          value: "postgresql://user:pass@db.pod:5432/mydb"
      volumes:
        - name: "uploads"
          path: "/app/uploads"
          size: "5Gi"
      secrets:
        - name: "api-key"
          data: "your-secret-api-key"
          path: "/app/secrets"
          fileName: "api-key.txt"
    
    - name: "db"
      image: "postgres:14"
      servicePorts:
        - port: 5432
          targetPort: 5432
          name: "postgres"
      vars:
        - key: "POSTGRES_USER"
          value: "user"
        - key: "POSTGRES_PASSWORD"
          value: "pass"
        - key: "POSTGRES_DB"
          value: "mydb"
      volumes:
        - name: "db-data"
          path: "/var/lib/postgresql/data"
          size: "10Gi"
```

## Usage in Code

### Loading and Validating

```go
import "github.com/Nexlayer/nexlayer-cli/pkg/schema"

// Load YAML
var config schema.NexlayerYAML
if err := yaml.Unmarshal(data, &config); err != nil {
    return err
}

// Validate
validator := schema.NewValidator(true)
if errors := validator.ValidateYAML(&config); len(errors) > 0 {
    // Handle validation errors
}
```

### Generating Configuration

```go
import "github.com/Nexlayer/nexlayer-cli/pkg/schema"

// Create generator
generator := schema.NewGenerator()

// Generate from project info
config, err := generator.GenerateFromProjectInfo("my-app", "nextjs", 3000)
if err != nil {
    return err
}
```

## Common Patterns

### Standard Web Application
```yaml
application:
  name: my-app
  url: my-app.nexlayer.dev
  registryLogin:
    registry: docker.io/my-org
    username: myuser
    personalAccessToken: mytoken
  pods:
    - name: web
      type: nextjs
      path: /
      image: <% REGISTRY %>/web:latest
      vars:
        - key: API_URL
          value: http://api.pod:8000
      servicePorts:
        - name: http
          port: 3000
          targetPort: 3000
    - name: api
      type: node
      path: /api
      image: <% REGISTRY %>/api:latest
      vars:
        - key: PORT
          value: "8000"
      servicePorts:
        - name: http
          port: 8000
          targetPort: 8000
```

### AI Application
```yaml
application:
  name: ai-app
  registryLogin:
    registry: docker.io/my-ai-org
    username: aiuser
    personalAccessToken: aitoken
  pods:
    - name: web
      type: langchain-nextjs
      path: /
      image: <% REGISTRY %>/web:latest
      vars:
        - key: OPENAI_API_KEY
          value: <% OPENAI_API_KEY %>
      annotations:
        ai.nexlayer.io/provider: openai
        ai.nexlayer.io/enabled: "true"
```

## Schema Evolution

This schema is versioned as v1.0. Future versions will maintain backward compatibility while adding new features. Any breaking changes will be clearly documented with migration guides.

## Related Packages

The schema implementation is spread across several packages:

- `pkg/schema`: Core schema types and validation logic (this package)
- `pkg/validation`: JSON Schema validation
- `pkg/vars`: Template variable processing

This document serves as the single source of truth that consolidates all this information. 