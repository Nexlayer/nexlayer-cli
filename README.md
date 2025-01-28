# Nexlayer CLI

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/Nexlayer/nexlayer-cli)](https://goreportcard.com/report/github.com/Nexlayer/nexlayer-cli)
[![GoDoc](https://godoc.org/github.com/Nexlayer/nexlayer-cli?status.svg)](https://godoc.org/github.com/Nexlayer/nexlayer-cli)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Deploy AI applications in seconds 

[Quick Start](#quick-start) • [Templates](#templates) • [Examples](#examples) • [Docs](https://docs.nexlayer.com)

</div>

## Quick Start

```bash
# Install
go install github.com/Nexlayer/nexlayer-cli@latest

# Initialize (auto-detects your stack)
nexlayer init myapp

# Deploy
nexlayer deploy
```

That's it! Your app is live in seconds 

## Templates

Choose your stack and start building:

### AI & LLM
```bash
# LangChain
nexlayer init myapp -t langchain-nextjs    # LangChain.js + Next.js
nexlayer init myapp -t langchain-fastapi   # LangChain Python + FastAPI
```

### Traditional
```bash
# Full-Stack
nexlayer init myapp -t mern    # MongoDB + Express + React + Node
nexlayer init myapp -t pern    # PostgreSQL + Express + React + Node
nexlayer init myapp -t mean    # MongoDB + Express + Angular + Node
```

## Examples

### LangChain Chat App
```yaml
# nexlayer.yaml
application:
  template:
    name: langchain-nextjs
    deploymentName: My Chat App
  pods:
    - type: nextjs
      exposeHttp: true
      name: app
      vars:
        - key: OPENAI_API_KEY
          value: your-key
        - key: LANGCHAIN_TRACING_V2
          value: "true"
```

### LangChain RAG App
```yaml
# nexlayer.yaml
application:
  template:
    name: langchain-fastapi
    deploymentName: My RAG App
  pods:
    - type: fastapi
      exposeHttp: true
      name: backend
      vars:
        - key: OPENAI_API_KEY
          value: your-key
        - key: PINECONE_API_KEY
          value: your-key
        - key: PINECONE_ENVIRONMENT
          value: gcp-starter
```

### MERN Stack App
```yaml
# nexlayer.yaml
application:
  template:
    name: mern
    deploymentName: My MERN App
  pods:
    - type: database
      exposeHttp: false
      name: mongodb
      vars:
        - key: MONGO_INITDB_DATABASE
          value: myapp
    - type: express
      exposeHttp: false
      name: backend
      vars:
        - key: MONGODB_URL
          value: DATABASE_CONNECTION_STRING
    - type: nginx
      exposeHttp: true
      name: frontend
      vars:
        - key: EXPRESS_URL
          value: BACKEND_CONNECTION_URL
```

## Template Configuration

Each Nexlayer deployment requires a YAML configuration file that defines your application structure. Here's how to configure it:

### Basic Structure
```yaml
application:
  template:
    name: my-app-stack          # Identifier for your app stack
    deploymentName: my-app      # Your deployment name
    registryLogin:              # Optional: for private registries
      username: user
      password: pass

  pods:                         # Define your app components
    - type: react              # Pod type (database/frontend/backend/etc)
      name: frontend           # Specific name for the pod
      tag: node:14-alpine      # Docker image
      privateTag: false        # Is it from a private registry?
      vars:                    # Environment variables
        - name: PORT
          value: "3000"
      exposeHttp: true        # Make pod accessible via HTTP
```

### Supported Pod Types
- **Database**: `postgres`, `mysql`, `neo4j`, `redis`, `mongodb`
- **Frontend**: `react`, `angular`, `vue`
- **Backend**: `django`, `fastapi`, `express`
- **Others**: `nginx`, `llm` (custom naming allowed)

### Environment Variables
Nexlayer automatically provides these environment variables to your pods:

| Variable | Description | Example |
|----------|-------------|---------|
| `PROXY_URL` | Your Nexlayer site URL | `https://your-site.alpha.nexlayer.ai` |
| `PROXY_DOMAIN` | Your Nexlayer site domain | `your-site.alpha.nexlayer.ai` |
| `DATABASE_HOST` | Database hostname | - |
| `DATABASE_CONNECTION_STRING` | Database connection string | `postgresql://user:pass@host:port/db` |
| `FRONTEND_CONNECTION_URL` | Frontend URL (with http://) | - |
| `BACKEND_CONNECTION_URL` | Backend URL (with http://) | - |
| `LLM_CONNECTION_URL` | LLM URL (with http://) | - |
| `FRONTEND_CONNECTION_DOMAIN` | Frontend domain (no prefix) | - |
| `BACKEND_CONNECTION_DOMAIN` | Backend domain (no prefix) | - |
| `LLM_CONNECTION_DOMAIN` | LLM domain (no prefix) | - |

### GitHub Actions Integration
Create `.github/workflows/docker-publish.yml`:

```yaml
name: Build and Push Docker Image

on:
  push:
    branches:
      - main

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
    - uses: actions/checkout@v2
    - uses: docker/login-action@v2
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    - run: echo "owner_lowercase=$(echo '${{ github.repository_owner }}' | tr '[:upper:]' '[:lower:]')" >> $GITHUB_ENV
    - uses: docker/build-push-action@v5
      with:
        context: .
        push: true
        tags: ghcr.io/${{ env.owner_lowercase }}/my-image-name:v0.0.1
```

## Features

- **Smart Detection**: Automatically detects your stack and configures everything
- **Simple Controls**: One command to initialize, one to deploy
- **Fast Cold Starts**: Sub-second startup times
- **Zero Config**: Sensible defaults for every stack
- **GPU Ready**: Built-in support for GPU acceleration
- **Cost Efficient**: Scale to zero when idle
- **Progress Feedback**: Visual progress indicators during operations
- **Error Handling**: Clear error messages and validation

## Plugins

Nexlayer supports plugins to extend its functionality. Plugins are Go shared libraries (.so files) that implement the Plugin interface.

### Using Plugins

```bash
# List installed plugins
nexlayer plugin list

# Run a plugin
nexlayer plugin run hello --name "John"

# Install a plugin
nexlayer plugin install ./my-plugin.so
```

### Creating Plugins

1. Create a new Go file for your plugin:

```go
package main

type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-plugin"
}

func (p *MyPlugin) Description() string {
    return "Description of what my plugin does"
}

func (p *MyPlugin) Run(opts map[string]interface{}) error {
    // Plugin logic here
    return nil
}

// Export the plugin
var Plugin MyPlugin
```

2. Build your plugin as a shared library:

```bash
go build -buildmode=plugin -o my-plugin.so my-plugin.go
```

3. Install your plugin:

```bash
nexlayer plugin install my-plugin.so
```

### Plugin Directory

Plugins are stored in `~/.nexlayer/plugins/`. Each plugin is a `.so` file that implements the Plugin interface.

### Plugin Interface

```go
type Plugin interface {
    // Name returns the name of the plugin
    Name() string
    
    // Description returns a description of what the plugin does
    Description() string
    
    // Run executes the plugin with the given options
    Run(opts map[string]interface{}) error
}

## Usage

```bash
# Deployment
nexlayer deploy          # Deploy your application
nexlayer status         # Check deployment status

# Configuration
nexlayer domain add     # Add custom domain

# AI-Powered Features
nexlayer init myapp     # Initialize a new app with AI-generated config
nexlayer ai detect      # Detect available AI assistants
nexlayer ai debug       # Get AI-powered deployment debugging
nexlayer ai scale       # AI-driven scaling recommendations
```

## AI Integration

Nexlayer CLI integrates with your IDE's AI capabilities to provide enhanced features:

### Automatic AI Detection
- Detects supported AI tools (GitHub Copilot, JetBrains AI, Cursor, Windsurf, Cline)
- Caches detection results in `~/.nexlayer/config.yaml`
- Runs automatically during installation or first `init`

```bash
$ nexlayer ai detect
✅ Detected AI Models:
   - GitHub Copilot (VS Code)
   - Cursor AI
```

### Smart YAML Generation
When using `nexlayer init`, the CLI:
- Analyzes your project structure
- Detects frameworks and dependencies
- Generates optimized deployment configuration

Example generated YAML:
```yaml
application:
  template:
    name: myapp
    deploymentName: myapp
    pods:
      - type: backend
        name: Node.js API
        tag: node:14
      - type: frontend
        name: React
        tag: nginx:latest
```

### AI-Powered Debugging
Debug deployment issues with AI assistance:
```bash
$ nexlayer ai debug --app myapp
❌ Deployment Error:
   - Issue: Missing environment variable `DATABASE_URL`
   - Suggested Fix: Add `DATABASE_URL` to your YAML under the `backend` pod

Suggested YAML Fix:
application:
  template:
    pods:
      - type: backend
        name: Node.js API
        vars:
          - key: DATABASE_URL
            value: mongodb://mongo-service
```

### Intelligent Scaling
Get AI-driven scaling recommendations:
```bash
$ nexlayer ai scale --app myapp
✅ Scaling Recommendation:
   - Current replicas: 2
   - Recommended replicas: 5 (based on traffic patterns)
```

## Testing

Run the test suite:

```bash
# Run all tests
./test/cli_test.sh

# Test specific functionality
nexlayer init myapp -t langchain-nextjs    # Test template initialization
nexlayer init myapp                        # Test auto-detection
```

The test suite covers:
- Command validation
- Template handling
- Project initialization
- Auto-detection
- Error scenarios
- Performance
- Concurrent operations

## Support
- [Documentation](https://docs.nexlayer.com)
- [GitHub Issues](https://github.com/Nexlayer/nexlayer-cli/issues)

## License

MIT
