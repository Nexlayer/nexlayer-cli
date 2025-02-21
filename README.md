# Nexlayer CLI

<div align="center">
  <img src="pkg/ui/assets/logo.svg" alt="Nexlayer Logo" width="400"/>
  <h1>Nexlayer CLI</h1>
  <p><strong>Deploy Full-Stack Applications in Seconds ⚡️</strong></p>
  <p>
    <a href="https://goreportcard.com/report/github.com/Nexlayer/nexlayer-cli">
      <img src="https://goreportcard.com/badge/github.com/Nexlayer/nexlayer-cli" alt="Go Report Card">
    </a>
    <a href="https://godoc.org/github.com/Nexlayer/nexlayer-cli?status.svg">
      <img src="https://godoc.org/github.com/Nexlayer/nexlayer-cli?status.svg" alt="GoDoc">
    </a>
    <a href="LICENSE">
      <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
    </a>
  </p>
</div>

## 🚀 Quick Start: Deploy Your Project in Seconds

### Installation Options

#### 1. Automated Installation (Recommended)
```bash
curl -sSL https://raw.githubusercontent.com/Nexlayer/nexlayer-cli/main/install.sh | bash
```
This script will:
- Check system requirements
- Install dependencies if needed
- Configure your PATH automatically
- Back up any existing installation
- Install the latest version of Nexlayer CLI

#### 2. Manual Installation
```bash
go install github.com/Nexlayer/nexlayer-cli@latest
```

### System Requirements
- Go 1.23.0 or higher
- Git (for development)
- 100MB free disk space

### Shell Configuration

The installer will automatically configure your shell. Supported shells:
- Bash (Linux: `~/.bashrc`, macOS: `~/.bash_profile`)
- Zsh (`~/.zshrc`)
- Fish (`~/.config/fish/config.fish`)

If you installed manually, add this to your shell configuration:
```bash
# For Bash/Zsh
export PATH=$PATH:~/go/bin

# For Fish
set -x PATH $PATH ~/go/bin
```

### Getting Started

1. **Initialize Your Project**
   ```bash
   nexlayer init
   ```
   This will:
   - Generate deployment configuration
   - Set up environment variables
   - Configure service dependencies

2. **Deploy Your Application**
   ```bash
   nexlayer deploy
   ```

That's it! Your app is now live. [Watch the demo →](https://nexlayer.dev/demo)

## 🎯 Why Nexlayer?

Nexlayer makes deploying full-stack applications fast, simple, and reliable. Here's why developers love it:

- **Smart Templates**: Production-ready configurations for any stack
- **One-Command Deploy**: No complex setup—just deploy
- **Live Watch**: Automatically redeploy on code changes
- **Real-Time Monitoring**: Track deployments with live logs
- **Built-in Security**: Automated security checks and best practices

## 💻 Core Commands

### Essential Commands
```bash
# Project Setup
nexlayer init              # Initialize your project
nexlayer init -i          # Interactive initialization

# Deployment
nexlayer deploy [appID]   # Deploy your application
nexlayer watch [appID]    # Watch for changes and auto-deploy

# Monitoring
nexlayer list            # List all deployments
nexlayer info <ns> <app> # Get deployment info

# Domain Management
nexlayer domain set      # Configure custom domain
```

### Development Commands
```bash
# Build and Test
make build              # Build optimized binary
make build-debug        # Build with debug symbols
make test              # Run tests with race detection
make coverage          # Generate test coverage report
make bench            # Run benchmarks

# Code Quality
make lint             # Run linters
make security         # Run security checks
make fmt              # Format code
make vet              # Run go vet

# Release
make release          # Create release builds
make docker          # Build Docker image
```

## 📝 Supported Project Types

Nexlayer supports a wide range of stacks out of the box:

### Web Frameworks
- `nextjs` - Next.js applications
- `react` - React applications
- `node` - Node.js/Express applications
- `python` - Python/FastAPI/Django applications
- `go` - Go applications

### Full-Stack Templates
- `mern` - MongoDB + Express + React + Node.js
- `pern` - PostgreSQL + Express + React + Node.js
- `mean` - MongoDB + Express + Angular + Node.js

More templates at [nexlayer.dev/templates](https://nexlayer.dev/templates)

## 👷 Development

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/Nexlayer/nexlayer-cli.git
cd nexlayer-cli

# Set up development tools
make setup

# Install dependencies
make deps

# Run tests and checks
make ci
```

### Development Features

- **Optimized Builds**: Static linking with trimpath and netgo tags
- **Security Checks**: Built-in gosec and dependency scanning
- **Performance Testing**: Benchmarking and race detection
- **Cross-Platform**: Builds for Linux, macOS (Intel/ARM), and Windows
- **CI/CD Ready**: Comprehensive test and build pipeline

## 📚 Documentation

### Core Documentation
- [YAML Reference](docs/reference/schemas/yaml/README.md) - How to write your `nexlayer.yaml` file
- [API Reference](docs/reference/api/README.md) - API endpoints used by the CLI

### Technical Reference
- YAML Schemas: [/docs/reference/schemas/yaml/](docs/reference/schemas/yaml/)
- API Schemas: [/docs/reference/schemas/api/](docs/reference/schemas/api/)

Full documentation at [nexlayer.dev/docs](https://nexlayer.dev/docs)

## 💪 Contributing

We love contributions! See our [Contributing Guide](CONTRIBUTING.md) for details on how to get involved.

## 📜 License

Nexlayer CLI is [MIT licensed](LICENSE).
