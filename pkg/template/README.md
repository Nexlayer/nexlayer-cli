# Nexlayer Template Package

This package is the single source of truth for all Nexlayer template-related code.

## Template Structure

```yaml
version: v2
application:
  name: myapp
  url: https://myapp.nexlayer.dev
  registryLogin:
    registry: docker.io
    username: myuser
    personalAccessToken: token123
  pods:
    - name: web
      type: frontend
      image: myapp/web:latest
      vars:
        - key: REACT_APP_API_URL
          value: http://api:8080
      servicePorts:
        - 3000
```

## Supported Pod Types

### Frontend
- `frontend`: Generic frontend
- `react`: React.js application
- `angular`: Angular application
- `vue`: Vue.js application

### Backend
- `backend`: Generic backend
- `express`: Express.js application
- `django`: Django application
- `fastapi`: FastAPI application

### Database
- `database`: Generic database
- `mongodb`: MongoDB database
- `postgres`: PostgreSQL database
- `redis`: Redis database
- `neo4j`: Neo4j database

### Other
- `nginx`: NGINX web server/proxy
- `llm`: Large Language Model service

## Default Images

Each pod type has a default Docker image. See `defaults.go` for the complete mapping.

Examples:
- `postgres` → `docker.io/library/postgres:latest`
- `redis` → `docker.io/library/redis:7`
- `mongodb` → `docker.io/library/mongo:latest`

## Validation Rules

1. **General**
   - Version must be specified (v1 or v2)
   - Application name must be alphanumeric
   - At least one pod must be specified

2. **Pods**
   - Pod name must be alphanumeric
   - Pod type must be one of the supported types
   - Image must be a valid Docker image reference
   - Service ports must be unique across all pods

3. **Volumes**
   - Volume size must match pattern: `^\d+[KMGT]i?$` (e.g., "1Gi", "500Mi")
   - Mount paths must start with "/"

4. **Environment Variables**
   - Keys must be valid Unix environment variable names
   - Values are required but can be empty

## Usage Example

```go
import "github.com/Nexlayer/nexlayer-cli/pkg/template"

// Create a validator
validator := template.NewValidator()

// Create a template
yaml := &template.NexlayerYAML{
    Version: template.V2,
    Application: template.Application{
        Name: "myapp",
        Pods: []template.Pod{
            {
                Name:  "web",
                Type:  template.Frontend,
                Image: "myapp/web:latest",
            },
        },
    },
}

// Validate the template
if err := validator.Validate(yaml); err != nil {
    log.Fatal(err)
}
```

## Best Practices

1. Always use the validator before processing templates
2. Use default images and configurations when possible
3. Follow the standard template structure
4. Keep pod names short but descriptive
5. Use semantic versioning for images
