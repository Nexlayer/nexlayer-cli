# Nexlayer Template Linter Plugin

A plugin for the Nexlayer CLI that validates YAML/JSON deployment templates against best practices and official Nexlayer schema guidelines.

## Features

- Validates template structure and required fields
- Checks for Kubernetes best practices
- Provides automatic fixes for common issues
- Supports both YAML and JSON templates

## Installation

1. Build the plugin:
```bash
go build -o nexlayer-lint
```

2. Move the binary to your Nexlayer plugins directory:
```bash
mv nexlayer-lint ~/.nexlayer/plugins/lint
```

## Usage

Basic validation:
```bash
nexlayer lint ./my-template.yaml
```

Validate and auto-fix issues:
```bash
nexlayer lint ./my-template.yaml --fix
```

## Example

The `example-template.yaml` file demonstrates a typical Nexlayer template with some common issues that the linter can detect and fix:

```yaml
name: my-app
version: 1.0.0
type: application
environment:
  stage: development
  ephemeral: true
resources:
  - name: web-service
    type: k8s/deployment
    properties:
      replicas: 3
      containers:
        - name: web
          image: nginx:latest
      # Missing labels - linter will catch this
  
  - name: database
    type: k8s/statefulset
    properties:
      replicas: 1
      containers:
        - name: postgres
          image: postgres:13
      labels:
        app: my-app  # Has proper labels
```

Running the linter on this template will:
1. Detect missing labels in the web-service resource
2. Offer to auto-fix the issue by adding appropriate labels
3. Validate that the database resource follows best practices

## Linting Rules

The linter checks for:

1. Required Fields:
   - Template name
   - Template version
   - Template type
   - Resource names and types

2. Kubernetes Best Practices:
   - Presence of required labels
   - Resource naming conventions
   - Basic configuration validation

3. Nexlayer Schema Guidelines:
   - Valid environment configuration
   - Proper resource type prefixes
   - Required property fields

## Auto-fix Capabilities

When run with `--fix`, the linter can automatically:
- Add missing labels
- Set default versions
- Apply naming conventions
- Add required fields with sensible defaults

## Contributing

Feel free to contribute additional linting rules or auto-fix capabilities by submitting a pull request!
