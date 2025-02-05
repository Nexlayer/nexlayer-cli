# Nexlayer CLI Plugins

Welcome to the Nexlayer CLI Plugin System! Extend the capabilities of Nexlayer CLI by building your own plugins. Whether you want to provide AI-powered deployment insights, integrate new commands, or enhance your deployment workflow, our plugin system is designed to be simple, flexible, and powerful.

## Table of Contents

- [Introduction](#introduction)
- [Overview of the Plugin System](#overview-of-the-plugin-system)
- [Plugin Interface](#plugin-interface)
- [Smart Deployments Plugin Example](#smart-deployments-plugin-example)
- [How to Build Your Own Plugin](#how-to-build-your-own-plugin)
- [Usage Examples](#usage-examples)
- [Contributing](#contributing)
- [Additional Resources](#additional-resources)

## Introduction

The Nexlayer CLI supports plugins to let you extend its functionality without modifying the core code. Our plugin system is built for speed and simplicity. It allows you to:
- Dynamically load and initialize plugins at runtime.
- Add new CLI commands that integrate seamlessly with Nexlayer Cloud's deployment behavior.
- Provide structured output that is both human-readable and machine-parsable (JSON) for AI-powered editors like Cursor or Windsurf.

## Overview of the Plugin System

The plugin system consists of several key components:
- **Plugin Interface:** A contract that every plugin must implement (Name, Description, Version, Commands, Init, and Run).
- **Plugin Manager:** Handles scanning a directory for plugins, loading shared libraries (`.so` files), and aggregating any additional commands.
- **Dependencies Injection:** Plugins receive common dependencies (API client, Logger, and UI Manager) so they can interact with Nexlayer Cloud consistently.

This design ensures that plugins integrate perfectly with the Nexlayer CLI, following the same YAML templating and API behavior used for deployments.

## Plugin Interface

Every plugin must implement the following interface:

```go
type Plugin interface {
    // Name returns the plugin's name.
    Name() string
    // Description returns a description of what the plugin does.
    Description() string
    // Version returns the plugin version.
    Version() string
    // Commands returns any additional CLI commands provided by the plugin.
    Commands() []*cobra.Command
    // Init initializes the plugin with dependencies.
    Init(deps *PluginDependencies) error
    // Run executes the plugin with the given options.
    Run(opts map[string]interface{}) error
}
```

Plugins receive their dependencies via the PluginDependencies structure:

```go
type PluginDependencies struct {
    APIClient api.APIClient
    Logger    *observability.Logger
    UIManager ui.Manager
}
```

The plugin manager will automatically load plugins from the specified directory, look up the exported Plugin symbol, and call its Init method to wire in these dependencies.

## Smart Deployments Plugin Example

The Smart Deployments Plugin is a great example of how to extend Nexlayer CLI. It provides AI-powered recommendations for deployment optimizations such as scaling, performance tuning, and pre-deployment audits.

### Key Features

**Deployment Recommendations:**
- Suggests changes to optimize resource allocation, such as increasing CPU limits for backend pods.

**Scaling Advice:**
- Recommends scaling configurations based on current traffic and usage patterns.

**Performance Tuning:**
- Identifies potential bottlenecks in your deployment (e.g., suboptimal environment variable settings or misconfigured ports).

**Pre-deployment Audit:**
- Runs an audit of your deployment configuration before you deploy.

### Example Usage

```bash
# Get deployment optimization recommendations:
nexlayer recommend deploy --ai --deploy

# Get resource scaling recommendations:
nexlayer recommend scale --ai --scale

# Get performance tuning recommendations:
nexlayer recommend performance --ai --performance

# Run a full pre-deployment audit:
nexlayer recommend audit --ai --pre-deploy
```

When the `--json` flag is added, the plugin outputs its recommendations in structured JSON format, making it easy for AI-powered editors to process the results.

## How to Build Your Own Plugin

Building a plugin for Nexlayer CLI is straightforward. Follow these steps:

### 1. Implement the Plugin Interface

Create a new Go file (e.g., `myplugin.go`) and implement the Plugin interface. For example:

```go
package main

import (
    "context"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/Nexlayer/nexlayer-cli/pkg/observability"
    "github.com/Nexlayer/nexlayer-cli/pkg/core/api"
    "github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

// MyPlugin is a sample plugin implementation.
type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-plugin"
}

func (p *MyPlugin) Description() string {
    return "A sample plugin that demonstrates how to extend Nexlayer CLI."
}

func (p *MyPlugin) Version() string {
    return "1.0.0"
}

func (p *MyPlugin) Commands() []*cobra.Command {
    cmd := &cobra.Command{
        Use:   "hello",
        Short: "Prints a hello message",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello from MyPlugin!")
        },
    }
    return []*cobra.Command{cmd}
}

func (p *MyPlugin) Init(deps *PluginDependencies) error {
    // Optionally use deps.APIClient, deps.Logger, or deps.UIManager.
    return nil
}

func (p *MyPlugin) Run(opts map[string]interface{}) error {
    // Implement non-interactive behavior if needed.
    return nil
}

// Export the plugin symbol.
var Plugin MyPlugin
```

### 2. Compile as a Shared Library

Use the following command to compile your plugin:

```bash
go build -buildmode=plugin -o myplugin.so myplugin.go
```

### 3. Place the Plugin

Copy the resulting `myplugin.so` into the default plugins directory or specify its path when running Nexlayer CLI.

### 4. Test and Contribute

Run `nexlayer plugin list` to verify your plugin is loaded. Then, you can start using the commands provided by your plugin and share your work with the Nexlayer community!

## Usage Examples

**List All Plugins:**
```bash
nexlayer plugin list
```

**Run a Plugin Command:**
```bash
nexlayer plugin run my-plugin --name "Sample"
```

**Use the Smart Deployments Plugin:**
```bash
nexlayer recommend deploy --ai --deploy
```

## Contributing

We welcome contributions from the community! If you have a plugin idea or have built one that enhances Nexlayer CLI, please submit a pull request or open an issue in our GitHub repository.

When contributing your plugin:
1. Follow the Plugin Interface guidelines.
2. Provide clear documentation and usage examples.
3. Ensure your plugin outputs structured logs (with a `--json` flag) to support both human users and AI-powered code editors.

## Additional Resources

- [Nexlayer CLI Documentation](#)
- [GitHub Repository](#)
- [Observability and Logging Guidelines](#)

Happy developing! Extend Nexlayer CLI to make deploying full-stack AI-powered applications even easier!
