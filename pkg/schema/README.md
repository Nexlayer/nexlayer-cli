# Nexlayer YAML Schema Documentation

The Nexlayer YAML schema defines the structure for deploying applications on the Nexlayer platform. This document provides a comprehensive reference for the schema format, validation rules, and examples.

## Package Structure

```
pkg/schema/
├── README.md           # This documentation
├── types.go           # Core schema types and structures
├── validator.go       # Schema validation logic
├── generator.go       # Schema generation utilities
└── examples/          # Example templates
    └── standard.go    # Standard schema examples
```

## Schema Structure

### Top-Level Structure
```yaml
application:
  name: string       # Required: Application name
  url: string       # Optional: Custom domain
  registryLogin:    # Required if using private images
    registry: string
    username: string
    personalAccessToken: string
  pods:             # Required: List of pods
    - name: string
      type: string
      path: string
      image: string
      # ... other pod fields
```

### Pod Configuration
```yaml
pods:
  - name: string            # Required: Pod name (lowercase alphanumeric with hyphens)
    type: string           # Required: Pod type (e.g., nextjs, react, node, python)
    path: string          # Optional: Route path (e.g., /, /api)
    image: string         # Required: Container image
    entrypoint: string    # Optional: Container entrypoint
    command: string       # Optional: Container command
    volumes:             # Optional: List of persistent volumes
      - name: string
        path: string
        size: string
        type: string
        readOnly: boolean
    secrets:            # Optional: List of secrets
      - name: string
        data: string
        path: string
        fileName: string
    vars:              # Optional: Environment variables
      - key: string
        value: string
    servicePorts:      # Required: List of service ports
      - name: string
        port: integer
        targetPort: integer
        protocol: string
    annotations:       # Optional: Pod annotations
      key: string
```

## Usage

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

// Add database pod
if err := generator.AddPod(config, "postgres", 5432); err != nil {
    return err
}

// Add AI configurations
generator.AddAIConfigurations(config, "openai")
```

## Validation Rules

### Application Name
- Must start with a lowercase letter
- Can contain lowercase letters, numbers, hyphens, and dots
- Example: `my-app.v1`

### Pod Names
- Must start with a lowercase letter
- Can contain lowercase letters, numbers, and hyphens
- Must be unique within the application
- Example: `web-server`, `api-v1`

### Image Names
- For private images: Must use `<% REGISTRY %>` placeholder
  - Example: `<% REGISTRY %>/myapp/api:v1.0.0`
  - Requires `registryLogin` configuration
- For public images: Standard Docker image format
  - Example: `nginx:latest`, `postgres:14`

### Service Ports
- Port numbers must be between 1 and 65535
- Port names must be unique within a pod
- Common port names: `http`, `https`, `api`, `metrics`

### Volumes
- Volume names must be unique within a pod
- Paths must start with `/`
- Size must use valid units: Ki, Mi, Gi, Ti
- Example: `1Gi`, `500Mi`

### Environment Variables
- Pod references use `.pod` suffix
  - Example: `http://api.pod:8000`
- Support variable substitution
  - Example: `<% API_KEY %>`

## Examples

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

### Database Configuration
```yaml
pods:
  - name: db
    type: postgres
    image: postgres:latest
    volumes:
      - name: pg-data
        path: /var/lib/postgresql/data
        size: 5Gi
    vars:
      - key: POSTGRES_USER
        value: <% POSTGRES_USER %>
      - key: POSTGRES_PASSWORD
        value: <% POSTGRES_PASSWORD %>
    servicePorts:
      - name: postgres
        port: 5432
        targetPort: 5432
```

## Best Practices

1. **Security**
   - Use environment variables for sensitive data
   - Never commit secrets to version control
   - Use private registry for custom images

2. **Naming**
   - Use descriptive pod names
   - Follow consistent naming conventions
   - Use semantic versioning for images

3. **Configuration**
   - Group related services in the same application
   - Use annotations for metadata
   - Document environment variables

4. **Resources**
   - Specify appropriate volume sizes
   - Use readiness/liveness probes
   - Configure appropriate port numbers

## Common Patterns

1. **Web + API**
   ```yaml
   pods:
     - name: web
       path: /
       # ... frontend config
     - name: api
       path: /api
       # ... backend config
   ```

2. **Database + Cache**
   ```yaml
   pods:
     - name: db
       type: postgres
       # ... database config
     - name: cache
       type: redis
       # ... cache config
   ```

3. **AI Services**
   ```yaml
   pods:
     - name: llm
       type: ollama
       # ... LLM config
     - name: vector-db
       type: qdrant
       # ... vector DB config
   ```

## Error Messages

Common validation errors and their solutions:

1. **Invalid Pod Name**
   ```
   Error: pods[0].name: invalid pod name format
   Solution: Use lowercase letters, numbers, and hyphens
   ```

2. **Missing Registry Login**
   ```
   Error: Private image used without registry login
   Solution: Add registryLogin block with credentials
   ```

3. **Invalid Port**
   ```
   Error: pods[0].servicePorts[0].port: invalid port number
   Solution: Use port number between 1-65535
   ```

## Schema Evolution

The schema is versioned and evolves with the platform. Breaking changes are announced in advance and backward compatibility is maintained through migration tools.

For the latest updates and changes, visit [docs.nexlayer.io/schema](https://docs.nexlayer.io/schema). 