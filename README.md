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
nexlayer init

# Deploy
nexlayer deploy
```

That's it! Your app is live 

## Templates

Choose your stack and start building:

### AI & LLM
```bash
# LangChain
nexlayer init --template langchain-nextjs    # LangChain.js + Next.js
nexlayer init --template langchain-fastapi   # LangChain Python + FastAPI
```

### Traditional
```bash
# Full-Stack
nexlayer init --template mern    # MongoDB + Express + React + Node
nexlayer init --template pern    # PostgreSQL + Express + React + Node
nexlayer init --template mean    # MongoDB + Express + Angular + Node
```

## Examples

### LangChain Chat App
```yaml
# nexlayer.yaml
application:
  template: langchain-nextjs-mongodb
  deploymentName: My Chat App
  pods:
    - type: database
      name: mongodb
      vars:
        - key: MONGO_INITDB_DATABASE
          value: langchain
    - type: nextjs
      name: app
      vars:
        - key: OPENAI_API_KEY
          value: your-key
        - key: LANGCHAIN_TRACING_V2
          value: true
```

### LangChain RAG App
```yaml
# nexlayer.yaml
application:
  template: langchain-fastapi
  deploymentName: My RAG App
  pods:
    - type: fastapi
      name: backend
      vars:
        - key: OPENAI_API_KEY
          value: your-key
        - key: PINECONE_API_KEY
          value: your-key
        - key: PINECONE_ENVIRONMENT
          value: gcp-starter
```

## Features

- ** Smart Detection**: Automatically detects your stack and configures everything
- ** Simple Controls**: One command to initialize, one to deploy
- ** Fast Cold Starts**: Sub-second startup times
- ** Zero Config**: Sensible defaults for every stack
- ** GPU Ready**: Built-in support for GPU acceleration
- ** Cost Efficient**: Scale to zero when idle

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
