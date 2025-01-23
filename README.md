# Nexlayer CLI

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/Nexlayer/nexlayer-cli)](https://goreportcard.com/report/github.com/Nexlayer/nexlayer-cli)
[![GoDoc](https://godoc.org/github.com/Nexlayer/nexlayer-cli?status.svg)](https://godoc.org/github.com/Nexlayer/nexlayer-cli)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/Nexlayer/nexlayer-cli)](https://github.com/Nexlayer/nexlayer-cli/releases)

[Documentation](https://docs.nexlayer.com) • [API Reference](https://docs.nexlayer.com/api) • [Support](https://nexlayer.com/support)

</div>

🚀 Deploy, manage and scale full-stack applications in minutes with Nexlayer CLI. Built for developers who value simplicity, speed and flexibility without sacrificing power.

## Get started in 3 simple steps:


```bash
# 1. Install Nexlayer CLI
go install github.com/Nexlayer/nexlayer-cli@latest

# 2. Log in to your Nexlayer account
nexlayer login  # Opens a browser for quick login

# 3. Deploy your first app!
nexlayer wizard

```

That's it! The wizard will guide you through deployment setup. Your GitHub authentication is handled automatically through Nexlayer.

## Recent Updates

- Restructured API types into dedicated package for better organization
- Enhanced AI suggestion plugin with improved client implementation
- Updated template builder plugin with additional features
- Improved authentication handling and client testing
- Added comprehensive info command functionality

## Requirements

- Go version 1.21 or later
- Git (for version control)
- Docker (for container builds)

## Dependencies

The CLI uses the following major dependencies:

- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [patrickmn/go-cache](https://github.com/patrickmn/go-cache) - In-memory caching
- [charmbracelet/bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [briandowns/spinner](https://github.com/briandowns/spinner) - Terminal spinners
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML support

## Installation

```bash
go install github.com/Nexlayer/nexlayer-cli@latest
```

Make sure your `$GOPATH/bin` is in your system PATH to access the CLI globally.

## Workflow Overview

Nexlayer CLI helps you manage both the build and deployment processes of your application:

1. **Build Process** (via CI commands)
   - Automate container image builds
   - Push to your preferred container registry
   - Run tests and quality checks

2. **Deployment Process** (via Nexlayer platform)
   - Deploy your built images
   - Manage scaling and resources
   - Monitor application health

## Common Commands

```bash
# Build Process Commands
# ---------------------
# Set up GitHub Actions workflow for building images
nexlayer ci setup github-actions --stack mern --registry ghcr.io

# Customize build parameters
nexlayer ci customize github-actions --image-tag v1.0.0 --build-context ./frontend

# Manage your container images
nexlayer ci images list
nexlayer ci images logs --image-name my-app --tag latest

# Deployment Commands
# ------------------
# Deploy your application (from your YAML deployment template configuration)
nexlayer deploy my-app --template my-app.yaml

# Scale your deployment to handle more traffic
nexlayer scale my-app --replicas 3

# Check the status of your app in real-time
nexlayer status my-app

# View live logs to debug or monitor
nexlayer logs my-app --follow

```

## CI/CD Integration

### Container Registry Setup

The CLI supports multiple container registries for your CI/CD workflows:

- GitHub Container Registry (GHCR)
- Docker Hub
- Google Artifact Registry (GCR)
- Amazon Elastic Container Registry (ECR)
- JFrog Artifactory
- GitLab Container Registry

#### Setting Up Container Registry

Use the `ci setup github-actions` command with appropriate flags:

```bash
# GitHub Container Registry (default)
nexlayer ci setup github-actions --registry-type ghcr

# Docker Hub
nexlayer ci setup github-actions --registry-type dockerhub

# Google Artifact Registry
nexlayer ci setup github-actions \
  --registry-type gcr \
  --registry-region us-east1 \
  --registry-project my-project

# Amazon ECR
nexlayer ci setup github-actions \
  --registry-type ecr \
  --registry-region us-east-1 \
  --registry-project 123456789012

# JFrog Artifactory
nexlayer ci setup github-actions \
  --registry-type artifactory \
  --registry your-artifactory-registry.jfrog.io

# GitLab Container Registry
nexlayer ci setup github-actions --registry-type gitlab
```

#### Required Secrets

Depending on your chosen registry, you'll need to configure different secrets in your GitHub repository:

- **Docker Hub**:
  - `DOCKERHUB_USERNAME`: Your Docker Hub username
  - `DOCKERHUB_TOKEN`: Your Docker Hub access token

- **Google Artifact Registry**:
  - `GOOGLE_CREDENTIALS`: Your Google Cloud service account key

- **Amazon ECR**:
  - `AWS_ACCESS_KEY_ID`: Your AWS access key ID
  - `AWS_SECRET_ACCESS_KEY`: Your AWS secret access key

- **JFrog Artifactory**:
  - `ARTIFACTORY_SERVER_ID`: Your Artifactory server ID
  - `ARTIFACTORY_USERNAME`: Your Artifactory username
  - `ARTIFACTORY_PASSWORD`: Your Artifactory password/token

- **GitLab Container Registry**:
  - `GITLAB_USERNAME`: Your GitLab username
  - `GITLAB_PASSWORD`: Your GitLab personal access token

### Build Automation

The CLI helps you set up automated builds using GitHub Actions:

```bash
# Generate workflow file for building container images
nexlayer ci setup github-actions --stack mern --registry ghcr.io
```

This creates a workflow that:
- Builds your container image
- Runs tests
- Pushes to your container registry

### Container Image Management

Monitor and manage your container images:

```bash
# List all images in your registry
nexlayer ci images list

# View build logs
nexlayer ci images logs --image-name my-app --tag latest
```

### Deployment

Once your images are built and pushed, deploy them using Nexlayer:

```bash
# Deploy the latest version
nexlayer deploy my-app

# Scale your deployment
nexlayer scale my-app --replicas 3
```

## Application Management

The CLI provides commands to manage your Nexlayer applications:

### List Applications

List all your applications:

```bash
nexlayer app list
```

### Create Application

Create a new application:

```bash
nexlayer app create --name "my-app-name"
```

## Example: Full Workflow

Here's a typical workflow using Nexlayer CLI:

1. **Set Up Build Pipeline**
```bash
# Generate GitHub Actions workflow
nexlayer ci setup github-actions --stack mern --registry ghcr.io
```

2. **Customize Build Settings**
```bash
# Configure build parameters
nexlayer ci customize github-actions \
  --image-tag v1.0.0 \
  --build-context ./frontend
```

3. **Deploy Your Application**
```bash
# Deploy the built image
nexlayer deploy my-app \
  --image ghcr.io/your-org/your-app:v1.0.0 \
  --env production
```

4. **Configure Services**
```bash
# Set environment variables for frontend service
nexlayer service configure --app my-app --service frontend \
  --env API_URL=https://api.example.com \
  --env FEATURE_FLAGS='{"dark_mode":true}'

# Set environment variables for backend service
nexlayer service configure --app my-app --service backend \
  --env DB_URL=postgres://db.example.com:5432/mydb \
  --env REDIS_URL=redis://cache.example.com:6379
```

5. **Monitor and Scale**
```bash
# Check deployment status
nexlayer status my-app

# Scale if needed
nexlayer scale my-app --replicas 3

# Visualize service connections
nexlayer service visualize --app my-app --format svg --output services.svg
```

## Service Configuration

### Environment Variables

Configure environment variables for your services:

```bash
# Set single environment variable
nexlayer service configure --app my-app --service frontend \
  --env API_URL=https://api.example.com

# Set multiple environment variables
nexlayer service configure --app my-app --service backend \
  --env DB_URL=postgres://db:5432/mydb \
  --env REDIS_URL=redis://cache:6379 \
  --env LOG_LEVEL=debug
```

### Service Visualization

Generate visual diagrams of your service connections:

```bash
# Print ASCII diagram to terminal
nexlayer service visualize --app my-app

# Generate SVG diagram
nexlayer service visualize --app my-app \
  --format svg \
  --output services.svg

# Generate PNG diagram
nexlayer service visualize --app my-app \
  --format png \
  --output services.png
```

## AI Suggest Feature

Enhance your development workflow with AI-powered suggestions directly from your terminal. The AI suggest feature provides intelligent recommendations for optimizing your Nexlayer applications using the latest AI models.

### Setup

Ensure you have the following environment variables set:

- `OPENAI_API_KEY`: Your OpenAI API key for accessing GPT-4
- `ANTHROPIC_API_KEY`: Your Anthropic API key for accessing Claude

### Usage

Run the AI suggest feature with the following command:

```bash
nexlayer ai-suggest --model openai --docs /path/to/docs --templates /path/to/templates
```

- `--model`: Specify the AI model to use (`openai` or `claude`)
- `--docs`: Path to your documentation directory
- `--templates`: Path to your templates directory

### Features

- **Interactive UI**: Navigate through suggestions with a beautiful terminal interface.
- **Markdown Rendering**: View code snippets and explanations with proper formatting.
- **Fuzzy Search**: Quickly find relevant documentation using fuzzy search.

### Future Improvements

- Streaming responses from AI models
- History of previous queries
- Export functionality

For more details, visit our [Documentation](https://docs.nexlayer.com) or [Support](https://nexlayer.com/support).

## Real-World Examples

### 1. CI/CD Pipeline Integration

```yaml
# .github/workflows/deploy.yml
name: Deploy with Nexlayer
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Install Nexlayer CLI
        run: go install github.com/Nexlayer/nexlayer-cli@latest
      - name: Deploy to Production
        run: |
          nexlayer deploy my-app \
            --env production \
            --auto-approve
        env:
          NEXLAYER_AUTH_TOKEN: ${{ secrets.NEXLAYER_AUTH_TOKEN }}
```

### 2. Multi-Environment Setup

```bash
# Development deployment
nexlayer deploy my-app --env staging

# Production deployment with increased resources
nexlayer deploy my-app --env production \
  --replicas 3 \
  --cpu 2 \
  --memory 4Gi
```

### 3. Advanced Configuration

```yaml
# nexlayer.yaml
version: '1'
app:
  name: my-awesome-app
  template: node-typescript
  env:
    NODE_ENV: production
    API_URL: https://api.example.com
  resources:
    cpu: 1
    memory: 2Gi
  scaling:
    min: 2
    max: 5
    targetCPU: 70
```

## Environment Setup

The CLI requires only one environment variable:

```bash
NEXLAYER_AUTH_TOKEN="your_auth_token"    # Your authentication token from https://app.nexlayer.io/settings/tokens
```

This token is used to authenticate your CLI requests with the Nexlayer API. You can get your token by:
1. Logging into your Nexlayer account
2. Going to Settings → API Tokens
3. Creating a new token with appropriate permissions

For different environments:
- Use `--env staging` (default) to deploy to `https://app.staging.nexlayer.io`
- Use `--env production` to deploy to `https://app.nexlayer.io`

## Security Best Practices

### Authentication Token Security

1. **Never commit tokens to version control**
   - Store tokens in environment variables
   - Use secure secret management in CI/CD pipelines
   - Consider using tools like HashiCorp Vault or AWS Secrets Manager

2. **Token Best Practices**
   - Rotate tokens regularly
   - Use tokens with minimal required permissions
   - One token per environment/purpose
   - Revoke tokens immediately if exposed

3. **Local Development**
   - Use `.env` files (not committed to git)
   - Different tokens for different environments
   ```bash
   # .env.development
   NEXLAYER_AUTH_TOKEN=dev_token
   NEXLAYER_STAGING_API_URL=http://localhost:8080
   ```

### Secure CI/CD Integration

```yaml
# GitHub Actions example with secure token handling
name: Deploy with Nexlayer
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Install Nexlayer CLI
        run: go install github.com/Nexlayer/nexlayer-cli@latest
      - name: Deploy to Production
        env:
          # Use GitHub's secret management
          NEXLAYER_AUTH_TOKEN: ${{ secrets.NEXLAYER_AUTH_TOKEN }}
        run: nexlayer deploy my-app --env production
```

### File Security

Ensure these files are in your `.gitignore`:
- `.env` and `.env.*` files
- `*.pem`, `*.key`, `*.cert` files
- Local development configurations
- Build artifacts and logs

## Troubleshooting

### Common Error Messages

1. **Authentication Error**
   ```
   ❌ Authentication required
   💡 Quick fixes:
      • Make sure NEXLAYER_AUTH_TOKEN is set in your environment
      • Run 'nexlayer auth login' to authenticate
      • Visit https://app.nexlayer.io/settings/tokens to generate a token
   ```

2. **YAML Validation Error**
   ```
   ❌ Invalid template YAML
   💡 Quick fixes:
      • Check the YAML syntax
      • Ensure all required fields are present
      • Run 'nexlayer validate' to check template structure
   ```

3. **Deployment Error**
   ```
   ❌ Deployment failed
   💡 Quick fixes:
      • Check your network connection
      • Verify your authentication token
      • Run with --debug flag for more information
   ```

### Debug Mode

Add the `--debug` flag to any command for detailed output:

```bash
nexlayer deploy my-app --debug
```

## Best Practices

1. **Version Control**: Always commit your `nexlayer.yaml` configuration
2. **Environment Variables**: Use `.env` files for local development
3. **CI/CD**: Use environment-specific configurations
4. **Monitoring**: Regularly check `nexlayer status` and `nexlayer logs`

## 🌟 Hidden Gems

Discover some powerful features that make Nexlayer CLI even more magical:

### 🎨 Template-based Deployments
Deploy applications using pre-configured templates:
```bash
# Deploy using a custom template
nexlayer deploy my-app --template my-app.yaml

# Use different templates for different environments
nexlayer deploy my-app --template staging.yaml
```

### 🔍 Deployment Management
Keep track of your deployments:
```bash
# List all your applications
nexlayer list

# View application details
nexlayer info my-app

# Check deployment status
nexlayer status my-app
```

### ⚡️ CI/CD Integration
Automate your workflow:
```bash
# Set up GitHub Actions integration
nexlayer ci setup github

# Customize CI/CD pipeline
nexlayer ci customize --template custom-pipeline.yaml
```

### 🌐 Domain Management
Configure custom domains:
```bash
# Set custom domain for your app
nexlayer domain set my-app --domain app.example.com

# List configured domains
nexlayer domain list my-app
```

### 🔌 Extend with Plugins

Nexlayer CLI supports a powerful plugin system that lets you extend its functionality. All plugins are hosted on GitHub under the `nexlayer/plugin-*` organization.

#### Available Plugins
- `hello`: A simple example plugin to get started
  ```bash
  nexlayer plugin install hello
  nexlayer hello world  # Outputs: Hello, world!
  ```

- `lint`: Code linting and style checking
  ```bash
  nexlayer plugin install lint
  nexlayer lint ./...        # Check code
  nexlayer lint ./... --fix  # Auto-fix issues
  ```

- `template-builder`: AI-powered deployment template generator
  ```bash
  # Install the plugin
  nexlayer plugin install template-builder

  # Generate optimized template using AI (requires OPENAI_API_KEY or ANTHROPIC_API_KEY)
  nexlayer template:generate  # Analyzes your project and generates optimal deployment templates
  nexlayer template:generate --dry-run  # Preview the generated template
  
  # The AI will optimize for:
  # - Resource allocation
  # - Security best practices
  # - Scalability
  ```

#### Managing Plugins
```bash
# Install a plugin
nexlayer plugin install <plugin-name>

# List all installed plugins
nexlayer plugin list

# Remove a plugin (just delete its directory)
rm -rf ~/.nexlayer/plugins/<plugin-name>
```

#### Create Your Own Plugin
Creating a Nexlayer plugin is straightforward. Here's a quick guide:

1. Create a new repository named `plugin-<your-plugin-name>`
2. Structure your plugin like this:
   ```
   plugin-example/
   ├── main.go       # Your plugin's entry point
   ├── README.md     # Plugin documentation
   └── go.mod        # Go module file
   ```

3. Example plugin code (based on the hello plugin):
   ```go
   package main

   import (
       "fmt"
       "github.com/spf13/cobra"
   )

   var rootCmd = &cobra.Command{
       Use:   "hello",
       Short: "A hello world plugin",
       Long:  `A hello world plugin for Nexlayer CLI.`,
       Args:  cobra.MinimumNArgs(1),
       Run: func(cmd *cobra.Command, args []string) {
           name := args[0]
           fmt.Printf("Hello, %s!\n", name)
       },
   }

   func main() {
       rootCmd.Execute()
   }
   ```

4. Push to GitHub as `nexlayer/plugin-<your-plugin-name>`
5. Install with `nexlayer plugin install <your-plugin-name>`

Plugins are installed in `~/.nexlayer/plugins` and are automatically integrated into the CLI. Start with the `hello` plugin as a template for your own plugins!

## Support

- 📚 [Documentation](https://docs.nexlayer.com)
- 🐛 [Issue Tracker](https://github.com/Nexlayer/nexlayer-cli/issues)
- 📧 [Email Support](mailto:support@nexlayer.com)
