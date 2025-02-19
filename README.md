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
# Initialize your project with intelligent stack detection
nexlayer init

# Deploy your app with automatic validation
nexlayer deploy
```

That's it! Your app is live. [Watch the demo →](https://nexlayer.dev/demo)

### Intelligent Project Configuration

Run `nexlayer init` in your project directory to automatically configure it for deployment. Nexlayer will analyze your current directory and:
- Detect your tech stack and dependencies
- Configure appropriate container images
- Set up health checks and environment variables
- Validate your configuration against best practices

### YAML Schema Compliance

Nexlayer uses a standardized YAML schema for deployment templates. Key features include:
- **Private Registry Support**: Use `<% REGISTRY %>` placeholder for private images
- **Dynamic Pod References**: Reference other pods using `<pod-name>.pod` format
- **URL References**: Use `<% URL %>` to reference your deployment's base URL
- **Flexible Port Configuration**: Support for both simple and detailed port formats
- **Automatic Validation**: Built-in schema validation with helpful error messages

Example template:
```yaml
application:
  name: my-app
  url: my-app.nexlayer.dev
  registryLogin:
    registry: ghcr.io/my-org
    username: myuser
    personalAccessToken: pat_token
  pods:
    - name: frontend
      type: nextjs
      path: /
      image: <% REGISTRY %>/frontend:latest
      servicePorts:
        - 3000  # Simple port format
      vars:
        - key: API_URL
          value: http://api.pod:8000
    - name: api
      type: backend
      path: /api
      image: <% REGISTRY %>/api:latest
      servicePorts:
        - name: http
          port: 8000
          targetPort: 8000
      vars:
        - key: DATABASE_URL
          value: postgresql://postgres:postgres@postgres.pod:5432/app
```

### Development Mode

During development, you can use the watch command to automatically redeploy when files change:

```bash
# Start watching for changes
nexlayer watch
```

The watch command will monitor your project files and automatically trigger a redeployment whenever changes are detected.

## ✨ Features

- 🤖 **AI-Powered Detection** - Automatically analyze and configure your project
- 🎯 **Smart Templates** - Production-ready templates for any stack
- ✅ **Built-in Validation** - Ensure configurations meet best practices
- 🔄 **Live Sync** - Keep configuration in sync with project changes
- 🚀 **One-Command Deploy** - Deploy full-stack apps instantly
- 📊 **Real-Time Monitoring** - Live logs and deployment status
- 👀 **Live Watch Mode** - Auto-redeploy on file changes during development
- 🔌 **Plugin System** - Extend functionality with custom plugins

## 📝 Templates

```bash
# Initialize your project with intelligent stack detection
nexlayer init
```

### AI/LLM Templates
- `langchain-nextjs` - LangChain.js + Next.js
- `openai-node` - OpenAI + Express + React
- `llama-py` - Llama.cpp + FastAPI
- More at [nexlayer.dev/templates](https://nexlayer.dev/templates)

### Full-Stack Templates
- `mern` - MongoDB + Express + React + Node.js
- `pern` - PostgreSQL + Express + React + Node.js
- `mean` - MongoDB + Express + Angular + Node.js

## 💻 Commands

```bash
# Initialize a new or existing project
nexlayer init [name]

# Deploy your application
nexlayer deploy [appID] --file deployment.yaml

# View status and logs
nexlayer list [appID]        # List deployments
nexlayer info <namespace> <appID>  # Get deployment info

# Configure custom domain
nexlayer domain set <appID> --domain example.com

# AI Features
nexlayer ai generate <app-name>  # Generate deployment template
nexlayer ai detect              # Detect project type

# Utility Commands
nexlayer feedback              # Send feedback
nexlayer completion [shell]    # Generate shell completion scripts

# Shell completion
nexlayer completion bash > ~/.bash_completion
nexlayer completion zsh > "${fpath[1]}/_nexlayer"
nexlayer completion fish > ~/.config/fish/completions/nexlayer.fish
```

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

- `pkg/validation/` - YAML schema validation and component type checking
- `pkg/compose/` - Docker compose generation and component detection
- `pkg/core/` - Core functionality and API types
- `pkg/commands/` - CLI command implementations
```

## 💪 Contributing

We love contributions! See our [Contributing Guide](CONTRIBUTING.md) for details.

## 📜 License

Nexlayer CLI is [MIT licensed](LICENSE).
