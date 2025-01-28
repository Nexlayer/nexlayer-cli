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

# Or specify a template
nexlayer init myapp -t langchain-nextjs

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
```

## Common Tasks

```bash
# Development
nexlayer dev              # Start local development
nexlayer test            # Run tests
nexlayer logs --follow   # Stream logs

# Deployment
nexlayer deploy          # Deploy to production
nexlayer rollback        # Instant rollback
nexlayer status         # Check status

# Configuration
nexlayer env set KEY=VALUE   # Set environment variable
nexlayer domain add example.com   # Add custom domain
nexlayer metrics              # View metrics
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

## GPU Support

Enable GPU acceleration with one line:

```bash
nexlayer deploy --gpu
```

Nexlayer automatically:
- Provisions GPU instances
- Configures CUDA environments
- Optimizes memory allocation
- Sets up monitoring

Available GPU types:
- NVIDIA T4 (Default)
- NVIDIA A100 (High-end)
- NVIDIA A10G (Mid-range)

## Support

- [Discord](https://discord.gg/nexlayer)
- [Documentation](https://docs.nexlayer.com)
- [GitHub Issues](https://github.com/Nexlayer/nexlayer-cli/issues)

## License

MIT
