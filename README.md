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

## 🚀 Quick Start

```bash
# Install Nexlayer CLI
go install github.com/Nexlayer/nexlayer-cli@latest
```

### Setting up your PATH

When you run `go install`, it places the Nexlayer CLI executable in a directory called `~/go/bin`. However, your computer needs to know where to find this executable when you type `nexlayer` in the terminal. Here's how to set it up:

1. First, add this line to your shell configuration file:
   ```bash
   export PATH=$PATH:~/go/bin
   ```

2. The configuration file location depends on your shell:
   - For Bash: `~/.bashrc` or `~/.bash_profile`
   - For Zsh: `~/.zshrc`

3. After adding the line, either:
   - Restart your terminal, or
   - Run `source ~/.bashrc` (or `source ~/.zshrc` for Zsh)

Now you can run Nexlayer commands from any directory!

```bash
# Initialize your project (auto-detects type)
nexlayer init

# Deploy your app
nexlayer deploy
```

That's it! Your app is live. [Watch the demo →](https://nexlayer.dev/demo)

### Intelligent Project Configuration

Run `nexlayer init` in your project directory to automatically configure it for deployment. Nexlayer will:
- Detect your tech stack and dependencies
- Configure appropriate container images
- Set up health checks and environment variables
- Validate your configuration against best practices
- Automatically detect custom ports from configuration files

Features:
- 🔍 **Auto-Detection**: Automatically identifies your project type and configuration
- 🎯 **Smart Templates**: Production-ready templates for any stack
- ✅ **Built-in Validation**: Ensures configurations meet best practices
- 🔄 **Live Watch**: Auto-redeploy on file changes during development
- 🚀 **One-Command Deploy**: Deploy full-stack apps instantly
- 📊 **Real-Time Monitoring**: Live logs and deployment status

## 💻 Commands

```bash
# Project Initialization
nexlayer init              # Auto-detect and initialize project
nexlayer init -i          # Interactive initialization

# Deployment
nexlayer deploy [appID]   # Deploy your application
nexlayer watch [appID]    # Watch for changes and auto-deploy

# Status and Monitoring
nexlayer list            # List all deployments
nexlayer info <ns> <app> # Get deployment info

# Domain Management
nexlayer domain set      # Configure custom domain

# Utility Commands
nexlayer feedback       # Send feedback
nexlayer completion    # Generate shell completions

# Shell completion
nexlayer completion bash > ~/.bash_completion
nexlayer completion zsh > "${fpath[1]}/_nexlayer"
nexlayer completion fish > ~/.config/fish/completions/nexlayer.fish
```

## 📝 Templates

Nexlayer supports various project types out of the box:

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

## 📚 Documentation

### Core Documentation
- [YAML Reference](docs/reference/schemas/yaml/README.md) - How to write your `nexlayer.yaml` file
- [API Reference](docs/reference/api/README.md) - API endpoints used by the CLI

### Technical Reference
- YAML Schemas: [/docs/reference/schemas/yaml/](docs/reference/schemas/yaml/)
- API Schemas: [/docs/reference/schemas/api/](docs/reference/schemas/api/)

Full documentation at [nexlayer.dev/docs](https://nexlayer.dev/docs)

## 👷 Development

```bash
# Clone the repository
git clone https://github.com/Nexlayer/nexlayer-cli.git
cd nexlayer-cli

# Install dependencies
make setup

# Run tests and validation
make test

# Run specific test suites
go test ./pkg/validation -v  # Run validation tests
go test ./pkg/compose -v     # Run compose tests
```

### Code Organization

- `pkg/core/` - Core functionality and domain logic
  - `api/` - API client and types
  - `types/` - Core type definitions
  - `template/` - Template generation
- `pkg/commands/` - CLI command implementations
- `pkg/validation/` - YAML schema validation
- `pkg/detection/` - Project type detection

## 💪 Contributing

We love contributions! See our [Contributing Guide](CONTRIBUTING.md) for details.

## 📜 License

Nexlayer CLI is [MIT licensed](LICENSE).
