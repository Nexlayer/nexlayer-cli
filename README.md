# Nexlayer CLI

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/Nexlayer/nexlayer-cli)](https://goreportcard.com/report/github.com/Nexlayer/nexlayer-cli)
[![GoDoc](https://godoc.org/github.com/Nexlayer/nexlayer-cli?status.svg)](https://godoc.org/github.com/Nexlayer/nexlayer-cli)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/Nexlayer/nexlayer-cli)](https://github.com/Nexlayer/nexlayer-cli/releases)

[Documentation](https://docs.nexlayer.com) • [API Reference](https://docs.nexlayer.com/api) • [Support](https://nexlayer.com/support)

</div>

🚀 Deploy, manage and scale full-stack applications in minutes with Nexlayer CLI. Built for developers who value simplicity, speed and flexibility without sacrificing power.

## Quick Start (30 seconds)

```bash
# 1. Install the CLI
go install github.com/Nexlayer/nexlayer-cli@latest

# 2. Login to Nexlayer (opens browser)
nexlayer login

# 3. Start the interactive wizard
nexlayer wizard
```

That's it! The wizard will guide you through deployment setup. Your GitHub authentication is handled automatically through Nexlayer.

## Requirements

- Go version 1.21 or later
- Git (for version control)
- Docker (for container builds)

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
# Deploy your built image
nexlayer deploy my-app

# Check deployment status
nexlayer status my-app

# Scale your deployment
nexlayer scale my-app --replicas 3

# View deployment logs
nexlayer logs my-app
```

## CI/CD Integration

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

## Support

- 📚 [Documentation](https://docs.nexlayer.com)
- 💬 [Discord Community](https://discord.gg/nexlayer)
- 🐛 [Issue Tracker](https://github.com/Nexlayer/nexlayer-cli/issues)
- 📧 [Email Support](mailto:support@nexlayer.com)
