# Nexlayer CLI

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/Nexlayer/nexlayer-cli)](https://goreportcard.com/report/github.com/Nexlayer/nexlayer-cli)
[![GoDoc](https://godoc.org/github.com/Nexlayer/nexlayer-cli?status.svg)](https://godoc.org/github.com/Nexlayer/nexlayer-cli)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/Nexlayer/nexlayer-cli)](https://github.com/Nexlayer/nexlayer-cli/releases)

[Documentation](https://docs.nexlayer.com) • [API Reference](https://docs.nexlayer.com/api) • [Support](https://nexlayer.com/support)

</div>

🚀 Deploy and manage full-stack applications in minutes with Nexlayer CLI. Built for developers who value simplicity without sacrificing power.

## Quick Start (30 seconds)

```bash
# 1. Install the CLI
go install github.com/Nexlayer/nexlayer-cli@latest

# 2. Start the interactive wizard
nexlayer wizard
```

That's it! The wizard will guide you through deployment setup. Want more control? Check out the manual setup below.

## Requirements

- Go version 1.21 or later
- Git (for version control)
- Docker (for container builds)

## Installation

```bash
go install github.com/Nexlayer/nexlayer-cli@latest
```

Make sure your `$GOPATH/bin` is in your system PATH to access the CLI globally.

## Common Commands

```bash
# Start the interactive wizard (recommended for first-time users)
nexlayer wizard

# Initialize a new project
nexlayer init

# Deploy your application
nexlayer deploy my-app

# Check deployment status
nexlayer status my-app

# View logs
nexlayer logs my-app

# Scale your application
nexlayer scale my-app --replicas 3
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
