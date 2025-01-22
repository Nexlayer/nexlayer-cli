# Nexlayer CLI

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/Nexlayer/nexlayer-cli)](https://goreportcard.com/report/github.com/Nexlayer/nexlayer-cli)
[![GoDoc](https://godoc.org/github.com/Nexlayer/nexlayer-cli?status.svg)](https://godoc.org/github.com/Nexlayer/nexlayer-cli)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/Nexlayer/nexlayer-cli)](https://github.com/Nexlayer/nexlayer-cli/releases)

[Documentation](https://docs.nexlayer.com) • [API Reference](https://docs.nexlayer.com/api) • [Support](https://nexlayer.com/support)

</div>

Deploy and manage full-stack applications in minutes with Nexlayer CLI. Built for developers who value simplicity without sacrificing power.

```bash
# Deploy your first application
curl -sf https://get.nexlayer.com | sh && nexlayer init && nexlayer wizard
```

## Documentation

For a comprehensive guide to using Nexlayer CLI, see our [documentation](https://docs.nexlayer.com).

## Requirements

- Operating system: Linux, macOS, or Windows
- Go version 1.21 or later
- Git (for version control)
- Docker (for container builds)

## Installation

```bash
go install github.com/Nexlayer/nexlayer-cli@latest
```

Make sure your `$GOPATH/bin` is in your system PATH to access the CLI globally.

## Getting Started

After installation, verify the CLI is properly installed:

```bash
nexlayer --version
```

### Quick Start with Wizard

The fastest way to get started is using our interactive wizard:

```bash
nexlayer wizard
```

The wizard will guide you through:
- Project initialization
- Configuration setup
- Deployment options
- Plugin selection and setup

### Manual Setup

1. Initialize a new project:
```bash
nexlayer init
```

2. Configure your project:
```bash
nexlayer config set
```

3. Deploy your application:
```bash
nexlayer deploy
```

For more detailed information and advanced usage, see our [documentation](https://docs.nexlayer.com).

## Quickstart

1. Initialize Nexlayer CLI:
   ```bash
   nexlayer init
   ```

2. Create a new deployment:
   ```bash
   nexlayer wizard
   ```

3. Follow the interactive prompts to:
   - Name your application
   - Select your tech stack
   - Configure deployment settings

4. Deploy your application:
   ```bash
   nexlayer deploy -f deployment/deployment.yaml
   ```

For detailed examples and use cases, see our [quickstart guide](https://docs.nexlayer.com/cli/quickstart).

## Core Concepts

### Deployments

A deployment represents your application stack, including:
- Frontend application
- Backend services
- Database instances
- Environment configuration

```bash
# Create a new deployment
nexlayer wizard

# Deploy an existing configuration
nexlayer deploy -f deployment.yaml
```

### Templates

Pre-configured application stacks for common use cases:

```bash
# List available templates
nexlayer template list

# Deploy using a template
nexlayer deploy --template three-tier-app
```

### Scaling

Adjust resources and replicas for your applications:

```bash
# Scale application replicas
nexlayer scale --app myapp --replicas 3

# Update resource limits
nexlayer scale --app myapp --cpu 1000m --memory 1Gi
```

## Command Reference

### Global Flags

All commands accept these flags:

- `--profile`: Configuration profile to use
- `--output`: Output format (json, yaml, table)
- `--quiet`: Suppress output
- `--debug`: Enable debug logging

### Core Commands

#### `nexlayer init`

Initialize Nexlayer CLI configuration.

```bash
nexlayer init [flags]
```

#### `nexlayer wizard`

Interactive deployment configuration wizard.

```bash
nexlayer wizard [flags]
```

#### `nexlayer deploy`

Deploy an application stack.

```bash
nexlayer deploy [flags]

Flags:
  -f, --file string       Path to deployment configuration
  -t, --template string   Template to use
      --dry-run          Validate without deploying
      --wait             Wait for deployment completion
```

For a complete command reference, see our [CLI documentation](https://docs.nexlayer.com/cli/commands).

## Development and Contributing

### Running from source

```bash
git clone https://github.com/nexlayer/nexlayer-cli.git
cd nexlayer-cli
go build
```

### Running tests

```bash
go test ./...
```

For development guidelines and best practices, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Support

- 📚 [Documentation](https://docs.nexlayer.com)
- 💬 [Community Forums](https://discuss.nexlayer.com)
- 🐛 [Issue Tracker](https://github.com/nexlayer/nexlayer-cli/issues)
- 📧 [Email Support](mailto:support@nexlayer.com)

## License

The Nexlayer CLI is licensed under the [MIT License](LICENSE).
