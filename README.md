# Nexlayer CLI

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/Nexlayer/nexlayer-cli)](https://goreportcard.com/report/github.com/Nexlayer/nexlayer-cli)
[![GoDoc](https://godoc.org/github.com/Nexlayer/nexlayer-cli?status.svg)](https://godoc.org/github.com/Nexlayer/nexlayer-cli)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/Nexlayer/nexlayer-cli)](https://github.com/Nexlayer/nexlayer-cli/releases)

[Documentation](https://docs.nexlayer.com) • [API Reference](https://docs.nexlayer.com/api) • [Support](https://nexlayer.com/support)

</div>

⚡️ Blazing-fast Kubernetes-powered CLI by Nexlayer for seamless deployment and scaling. Launch full-stack applications in seconds with enterprise-grade infrastructure.

## 🚀 Instant Value in 30 Seconds

```bash
# Install & deploy your first app
1. go install github.com/Nexlayer/nexlayer-cli@latest
2. nexlayer login
3. nexlayer wizard  # AI-powered deployment setup
```

## ✨ Why Developers Love It

- **Instant Compute**: Kubernetes-native design with sub-second cold starts
- **Zero Config**: AI-powered setup detects your stack and configures everything
- **Infinite Scale**: Auto-scaling from 0 to 1000s of instances in seconds
- **Developer Flow**: Git-native workflow with instant preview environments
- **Cost Efficient**: Pay only for actual compute time, scale to zero when idle

## Core Features

- 🎯 One-command deploys with automatic infrastructure provisioning
- 🔄 Real-time logs and metrics with built-in monitoring
- 🛡️ Enterprise-grade security with automatic SSL and secrets management
- 🌐 Kubernetes-powered compute with auto-scaling   
- 🏁 production-ready full-stack templates 
- 🤖 AI-powered suggestions and optimizations

## Requirements
- Go 1.21+
- Git
- Docker

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

## Deployment Examples

### Basic Deployment
```bash
# Deploy a simple application
nexlayer deploy myapp
```

### Using Private Container Images

When deploying applications that use private container images, you'll need to configure registry credentials in your deployment template. Here's a comprehensive example of a MERN stack deployment using private GitHub Container Registry (GHCR) images:

```yaml
application:
  template: 
    name: "mongodb-express-react-nodejs"
    deploymentName: "My MERN Stack"
    registryLogin:
      registry: ghcr.io
      username: <Github Username>
      personalAccessToken: <GitHub Read:Packages Personal Access Token>
    pods:
    - type: database
      exposeOn80: false
      name: mongoDB
      tag: ghcr.io/<Github Lowercase Username>/mern-mongo:v0.01
      privateTag: true
      vars:
      - key: MONGO_INITDB_ROOT_USERNAME
        value: mongo
      - key: MONGO_INITDB_ROOT_PASSWORD
        value: passw0rd
      - key: MONGO_INITDB_DATABASE
        value: todo
    - type: express
      exposeOn80: false
      name: express
      tag: ghcr.io/<Github Lowercase Username>/mern-express:v0.01
      privateTag: true
      ports:
      - name: express
        containerPort: 3000
        servicePort: 3000
    - type: nginx
      exposeOn80: true
      name: react
      tag: ghcr.io/<Github Lowercase Username>/mern-react:v0.01
      privateTag: true
      vars:
      - key: EXPRESS_URL
        value: BACKEND_CONNECTION_URL
      ports:
      - name: react
        containerPort: 80
        servicePort: 80
```

Key configuration points:
- `registryLogin`: Specifies credentials for private registry access
- `privateTag: true`: Indicates the image requires authentication
- `tag`: Full path to your private container image
- `exposeOn80`: Controls whether the pod should be exposed on port 80
- `vars`: Environment variables for container configuration
- `ports`: Container and service port mappings

To deploy using this template:
```bash
# Save the template as mern-stack.yaml
nexlayer deploy -f mern-stack.yaml
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

## AI-Powered Assistance

Nexlayer CLI includes AI-powered suggestions to help optimize your deployments. To enable AI suggestions, use the `--ai` flag with any command:

```bash
# Set up your OpenAI API key
export OPENAI_API_KEY=your_api_key_here

# Get AI suggestions for your commands
nexlayer deploy --ai
nexlayer configure --ai
```

The AI assistant will analyze your command and provide suggestions for:
- Resource optimization
- Security best practices
- Scaling strategies
- Configuration improvements

These suggestions are optional and you can choose to apply them or proceed with your original command.

## 🤖 AI Assistant

Add `--ai` to any command to get intelligent suggestions and improvements:

```bash
# Get AI suggestions while creating a template
nexlayer init my-app --ai

# Get deployment optimization suggestions
nexlayer deploy my-app --ai

# Get configuration recommendations
nexlayer configure my-app --ai
```

The AI will:
- Analyze your command context
- Suggest improvements
- Offer best practices
- Help troubleshoot issues

> Note: AI requires OpenAI or Anthropic API key. Set with `export OPENAI_API_KEY="your-key"` or configure in settings.

### Examples

```bash
# Initialize with AI suggestions
$ nexlayer init my-app --ai
✨ Creating new app "my-app"
🤖 AI Suggestions:
  • Add health checks for better reliability
  • Configure auto-scaling based on CPU usage
  • Set up monitoring endpoints
Apply these suggestions? [Y/n]

# Deploy with AI optimization
$ nexlayer deploy my-app --ai
🚀 Deploying "my-app"
🤖 AI Suggestions:
  • Increase replica count for high availability
  • Add resource limits to prevent overload
  • Enable SSL for security
Apply these suggestions? [Y/n]

# Configure with AI assistance
$ nexlayer configure my-app --ai
⚙️ Configuring "my-app"
🤖 AI Suggestions:
  • Set up environment-specific variables
  • Add backup strategy
  • Configure logging
Apply these suggestions? [Y/n]
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
