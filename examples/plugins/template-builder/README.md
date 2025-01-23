# Nexlayer Template Builder

An intelligent infrastructure template generator with AI-powered refinements and security scanning.

## Features

- 🤖 AI-powered template generation and refinement
- 🔍 Automatic stack detection for various technologies
- 🛡️ Comprehensive security scanning
- 💰 Infrastructure cost estimation
- 📦 Remote template registry integration
- 🔄 Template versioning and upgrades
- 🔧 Shell completion for enhanced productivity
- ⚙️ Flexible configuration management

## Template Structure

Nexlayer templates follow a specific YAML structure for deploying applications:

```yaml
application:
  template:
    name: "my-stack-name"              # Template stack identifier
    deploymentName: "My Application"    # Deployment name in Nexlayer
    registryLogin:                      # Optional: for private registries
      registry: ghcr.io
      username: <username>
      personalAccessToken: <token>
    pods:                              # Define application components
    - type: database                   # Pod type (database, llm, django, etc.)
      exposeHttp: false               # Whether to expose via HTTP
      name: mongoDB                    # Specific pod name
      tag: mongo:latest               # Docker image
      privateTag: false               # Is it from private registry?
      vars:                           # Environment variables
      - key: MONGO_INITDB_ROOT_USERNAME
        value: mongo
    - type: express
      exposeHttp: false
      name: backend
      tag: myapp-backend:v1.0
      privateTag: true
      vars:
      - key: MONGODB_URL
        value: DATABASE_CONNECTION_STRING
    - type: nginx
      exposeHttp: true
      name: frontend
      tag: myapp-frontend:v1.0
      privateTag: true
      vars:
      - key: EXPRESS_URL
        value: BACKEND_CONNECTION_URL
```

### Supported Pod Types
- **Database**: `postgres`, `mysql`, `neo4j`, `redis`, `mongodb`
- **Frontend**: `react`, `angular`, `vue`
- **Backend**: `django`, `fastapi`, `express`
- **Others**: `nginx`, `llm`

## Nexlayer-Provided Environment Variables

Nexlayer automatically injects these environment variables into your pods:

### Core Variables
- `PROXY_URL`: Full URL of your Nexlayer site (e.g., `https://your-site.alpha.nexlayer.ai`)
- `PROXY_DOMAIN`: Domain of your Nexlayer site (e.g., `your-site.alpha.nexlayer.ai`)

### Database Variables
- `DATABASE_HOST`: Database hostname
- `NEO4J_URI`: Neo4j database URI
- `DATABASE_CONNECTION_STRING`: Full database connection string

### Service Connection URLs
- `FRONTEND_CONNECTION_URL`: URL to frontend pod (with `http://` prefix)
- `BACKEND_CONNECTION_URL`: URL to backend pod (with `http://` prefix)
- `LLM_CONNECTION_URL`: URL to LLM pod (with `http://` prefix)

### Service Connection Domains
- `FRONTEND_CONNECTION_DOMAIN`: Frontend pod domain (without prefix)
- `BACKEND_CONNECTION_DOMAIN`: Backend pod domain (without prefix)
- `LLM_CONNECTION_DOMAIN`: LLM pod domain (without prefix)

## AI Model Integration

When using the template-builder with AI models (OpenAI or Claude), the following prompt ensures generated templates follow Nexlayer's structure:

```
You are generating a Nexlayer infrastructure template. Please follow these guidelines:

1. Use the standard Nexlayer pod-based YAML structure
2. Include required template fields: name, deploymentName, and optional registryLogin
3. Define pods with correct type, name, tag, and exposeHttp settings
4. Use Nexlayer-provided environment variables for service connections
5. Configure appropriate environment variables for each pod type
6. Follow pod naming conventions for databases, frontends, and backends
7. Set proper exposure settings for public-facing services
8. Include registry configuration for private images
9. Add descriptive comments for maintainability
10. Follow security best practices for sensitive variables

Available Nexlayer environment variables:
- PROXY_URL, PROXY_DOMAIN
- DATABASE_HOST, NEO4J_URI, DATABASE_CONNECTION_STRING
- FRONTEND_CONNECTION_URL, BACKEND_CONNECTION_URL, LLM_CONNECTION_URL
- FRONTEND_CONNECTION_DOMAIN, BACKEND_CONNECTION_DOMAIN, LLM_CONNECTION_DOMAIN
```

## Installation

### Using Go

```bash
go install github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2@latest
```

### Shell Completion

Enable shell completion for your preferred shell:

```bash
# Bash
nexlayer completion bash > /usr/local/etc/bash_completion.d/nexlayer

# Zsh
nexlayer completion zsh > "${fpath[1]}/_nexlayer"

# Fish
nexlayer completion fish > ~/.config/fish/completions/nexlayer.fish

# PowerShell
nexlayer completion powershell > nexlayer.ps1
```

## Configuration

Nexlayer uses a configuration file located at `~/.nexlayer/config.json`. You can configure:

- Registry URL
- Default template
- API keys
- Output format
- Verbosity

Example configuration:
```json
{
  "registry_url": "https://registry.nexlayer.dev",
  "default_template": "default",
  "api_keys": {
    "openai": "your-api-key"
  },
  "output_format": "yaml",
  "verbose": false
}
```

## Quick Start

1. Initialize a new template:
```bash
nexlayer init my-app
```

2. Generate a template from an existing project:
```bash
nexlayer generate ./my-project --output yaml
```

3. Upgrade a template with security scanning:
```bash
nexlayer upgrade my-app.yaml --security-scan --estimate-costs
```

4. Compare two templates:
```bash
nexlayer diff template-v1.yaml template-v2.yaml
```

## Command Reference

### Global Flags
- `--config`: Path to config file (default: ~/.nexlayer/config.json)
- `--verbose, -v`: Enable verbose output
- `--registry`: Template registry URL

### `nexlayer init [template-name]`
Initialize a new template with best practices and defaults.

**Flags:**
- `--type`: Template type (service|app|function)
- `--stack`: Technology stack (node|python|go)
- `--template`: Base template to use

**Example:**
```bash
nexlayer init my-service --type service --stack node
```

### `nexlayer generate [project-dir]`
Generate a template from an existing project.

**Flags:**
- `--output, -o`: Output format (yaml|json)
- `--dry-run`: Preview without writing
- `--ai-provider`: AI provider for refinement (openai|claude)
- `--exclude`: Patterns to exclude
- `--include-deps`: Include dependencies

**Example:**
```bash
nexlayer generate ./my-project -o yaml --ai-provider openai
```

### `nexlayer upgrade [template-file]`
Upgrade a template to the latest version.

**Flags:**
- `--security-scan`: Run security scan
- `--estimate-costs`: Estimate infrastructure costs
- `--force`: Force upgrade even with breaking changes
- `--backup`: Create backup before upgrading

**Example:**
```bash
nexlayer upgrade my-app.yaml --security-scan --backup
```

### `nexlayer diff [template1] [template2]`
Show differences between two templates.

**Flags:**
- `--format`: Diff format (unified|context|json)
- `--ignore-whitespace`: Ignore whitespace changes
- `--summary`: Show only summary of changes

**Example:**
```bash
nexlayer diff old.yaml new.yaml --format unified
```

## Template Registry

Share and reuse templates using the Nexlayer Registry:

```bash
# Publish a template
nexlayer publish my-app.yaml --version 1.0.0

# Download a template
nexlayer download my-app:1.0.0

# List available templates
nexlayer list --filter "type=service"
```

## Security Scanning

The security scanner performs comprehensive checks:

### Infrastructure Security
- Resource encryption configuration
- IAM roles and permissions
- Network security groups
- Public exposure analysis

### Application Security
- Exposed ports and protocols
- TLS configuration
- Secret management practices
- Environment variable analysis

### Compliance
- GDPR requirements
- SOC 2 controls
- PCI DSS requirements
- HIPAA compliance checks

## Cost Estimation

The cost estimator analyzes:

- Compute resources (CPU, memory)
- Storage requirements (volume types, sizes)
- Network traffic patterns
- Region-specific pricing
- Reserved instance opportunities
- Spot instance possibilities

## Error Handling

Nexlayer provides detailed error messages with:

- Error codes for programmatic handling
- Context-specific information
- Suggested solutions
- Debugging information in verbose mode

Example error output:
```
[TEMPLATE_INVALID] Invalid template structure: missing required field 'resources'
Context:
  file: my-template.yaml
  line: 25
  field: resources
Suggestion: Add a 'resources' section to define infrastructure components
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Run tests: `go test ./... -race -cover`
4. Submit a pull request

### Development Setup

1. Clone the repository:
```bash
git clone https://github.com/nexlayer/nexlayer-cli.git
cd nexlayer-cli
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
make test
```

4. Build locally:
```bash
make build
```

## Support

- 📖 [Documentation](https://docs.nexlayer.dev)
- 💬 [Discord Community](https://discord.gg/nexlayer)
- 📧 [Email Support](mailto:support@nexlayer.dev)
- 🐛 [Issue Tracker](https://github.com/nexlayer/nexlayer-cli/issues)

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
