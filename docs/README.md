# Nexlayer CLI Documentation

Welcome to the Nexlayer CLI documentation! This directory contains comprehensive documentation for using and configuring the Nexlayer CLI.

## 📚 Documentation Structure

### Configuration
- [YAML Schema Documentation (v1.2)](configuration/yaml-reference.md)
  - Complete guide to the Nexlayer YAML format
  - Detailed examples and best practices
  - Troubleshooting common issues
- [Example Template](configuration/template.v2.yaml)
  - Full example with all configuration options
  - Extensively commented for clarity
- [JSON Schema](configuration/schema.v2.json)
  - Official JSON Schema for validation
  - Used by `nexlayer validate` command

### Quick Links
- [Main README](../README.md) - Quick start guide and feature overview
- [Contributing Guide](../CONTRIBUTING.md) - How to contribute to Nexlayer CLI
- [Official Documentation](https://nexlayer.dev/docs) - Online documentation and guides

## 🚀 Getting Started

The fastest way to get started is to:

1. Install the CLI:
   ```bash
   go install github.com/Nexlayer/nexlayer-cli@latest
   ```

2. Create a new project:
   ```bash
   nexlayer init my-app
   ```

3. Deploy your app:
   ```bash
   nexlayer deploy
   ```

## 🔍 Finding Help

- Use `nexlayer --help` to see all available commands
- Use `nexlayer [command] --help` for detailed command help
- Visit [nexlayer.dev/docs](https://nexlayer.dev/docs) for guides and tutorials
- Join our [Discord community](https://discord.gg/nexlayer) for support
