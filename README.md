# Nexlayer CLI

<div align="center">
  <img src="assets/logo.svg" alt="Nexlayer Logo" width="400"/>
  <h1>Deploy Full-Stack Applications in Seconds</h1>
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

---

## Table of Contents
- [Quick Start](#quick-start)
- [Getting Started with Nexlayer CLI](#getting-started-with-nexlayer-cli)
- [Features](#features)
- [Templates](#templates)
- [Commands Overview](#commands-overview)
- [Documentation & Support](#documentation--support)
- [Local Development](#local-development)
- [Contributing](#contributing)
- [License](#license)

---

## Quick Start

Install the Nexlayer CLI using Go:
Get up and running in **3 seconds** with these three simple commands:

### Install the CLI
```bash
go install github.com/Nexlayer/nexlayer-cli@latest
```
> **Tip:** Ensure `$GOPATH/bin` is in your `PATH` so that the nexlayer command is recognized.

### Initialize Your Project
**Create the Project Directory**
- Replace `my_app_name` with your project name:
  ```bash
  mkdir my_app_name
  cd my_app_name
  ```
# Create a new project
nexlayer init my-app

# Initialize an existing project
cd my-existing-app
nexlayer init
```
This command uses AI-powered detection to analyze your project and automatically generates a  `nexlayer.yaml` configuration file. This file defines your application stack, pods, and environment variables according to Nexlayer Cloud’s templating system—so you're ready for deployment.

### Deploy Your App
Once your project is initialized and the configuration file is in place, deploy your app with:

```bash
nexlayer deploy
```
Watch your full-stack AI app go live instantly!

> **Bonus:** To see a demo, check out our demo video.

---

## Getting Started with Nexlayer CLI

### Installation
Nexlayer CLI requires Go 1.18+.

**Video: Installation**

### Environment Setup
We recommend setting up a Go workspace for your project.

**Install Nexlayer CLI**
- **Prerequisites:** Ensure Go is installed and `$GOPATH/bin` is in your `PATH`.

**Create the Project Directory**
- Replace `my_app_name` with your project name:
  ```bash
  mkdir my_app_name
  cd my_app_name
  ```

**Install Nexlayer CLI**
- Run the following command:
  ```bash
  go install github.com/Nexlayer/nexlayer-cli@latest
  ```

### Initialize the Project
- Use the CLI to initialize your project:
  ```bash
  nexlayer init
  ```
  - Choose from predefined templates or start with a blank setup.

### Run the App
- Deploy and run your app in development mode:
  ```bash
  nexlayer deploy
  ```
  - Your app will be accessible at the specified endpoint.

**Logs and Debugging**
- Increase log verbosity with:
  ```bash
  nexlayer deploy --loglevel debug
  ```

**Hot Reload**
- Nexlayer supports hot reloading for real-time code changes.

---

## Features

- **Smart Project Templates:** Start with production-ready templates for full-stack or backend-only applications.
- **Intelligent Project Detection:** Automatically analyze existing projects and generate the perfect configuration.
- **Project Synchronization:** Keep your configuration in sync with project changes using `nexlayer sync`.
- **One-Command Deployment:** Deploy your app with a single command.
- **Real-Time Logs & Status:** Monitor your deployment status and view logs easily.
- **Custom Domain & Feedback:** Attach custom domains and send feedback to help us improve.
- **Plugin Support:** Extend Nexlayer CLI with custom plugins.

---

## Templates

Nexlayer offers production-ready templates to help you get started quickly:


### AI/LLM
- **langchain-nextjs**: LangChain.js with Next.js
- **langchain-fastapi**: LangChain Python with FastAPI
- **openai-node**: OpenAI with Express and React
- **openai-py**: OpenAI with FastAPI and Vue
- **llama-node**: Llama.cpp with Next.js
- **llama-py**: Llama.cpp with FastAPI
- **vertex-ai**: Google Vertex AI with Flask
- **huggingface**: Hugging Face with FastAPI
- **anthropic-py**: Anthropic Claude with FastAPI
- **anthropic-js**: Anthropic Claude with Next.js


### Traditional Full Stack App
- **MERN Stack**: Create a full-stack app with MongoDB, Express, React, and Node.js
- **MEAN Stack**: Build modern web applications with MongoDB, Express, Angular, and Node.js
- **MNFA Stack**: Create a powerful app with MongoDB, Neo4j, FastAPI, and Angular
- **MEVN Stack**: Build full-stack applications with MongoDB, Express, Vue.js, and Node.js
- **PERN Stack**: Build full-stack applications with PostgreSQL, Express, React, and Node.js

> **Tip:** Use the interactive template selector:
```bash
nexlayer init my-app
```
**Create the Project Directory**
- Replace `my_app_name` with your project name:
  ```bash
  mkdir my_app_name
  cd my_app_name
---

## Commands Overview

### Initialization
```bash
nexlayer init
```
Automatically generate a deployment template (`nexlayer.yaml`).

### Deployment
```bash
nexlayer deploy
```
Deploy your application using the generated template.

### Status & Logs
```bash
# View deployment status
nexlayer status

# View real-time logs
nexlayer logs -f [podName]

# Keep configuration in sync with project
nexlayer sync
```
Monitor your deployment, view real-time logs, and keep your configuration up to date.

### Domain Management
```bash
nexlayer domain add yourdomain.com
```
Add a custom domain to your app.

### Feedback
```bash
nexlayer feedback "Your feedback message"
```
Help us improve by sending your feedback.

### AI Assistance
```bash
nexlayer ai generate
nexlayer ai detect
nexlayer ai debug
nexlayer ai scale
```
Leverage AI for template generation, debugging, and scaling recommendations.

---

## Plugins

Nexlayer CLI supports a powerful plugin system that extends its functionality. Plugins can add new commands, provide AI-powered recommendations, and enhance your deployment workflow.

### Available Plugins

#### Smart Deployments Plugin
Provides AI-powered recommendations for optimizing your deployments:

```bash
# Get deployment optimization recommendations
nexlayer recommend deploy

# Get resource scaling recommendations
nexlayer recommend scale

# Get performance tuning suggestions
nexlayer recommend performance

# Run a pre-deployment audit
nexlayer recommend audit
```

Add the `--json` flag to any command to get machine-readable output.

### Creating Plugins

To create a new plugin:

1. Implement the `Plugin` interface in your Go code
2. Build your plugin as a shared object (.so file)
3. Place the .so file in the plugins directory

See our [Plugin Development Guide](#) for detailed instructions.

---

## Documentation & Support

- [Nexlayer Documentation](#)
- [GitHub Issues](#)

---

## Local Development

For testing locally, use Docker Compose:

### Initialize Your Project
```bash
nexlayer init myapp -t <template>
```

### Generate Docker Compose File
```bash
nexlayer compose generate
```

### Start Local Environment
```bash
nexlayer compose up
```

### View Logs
```bash
nexlayer compose logs -f [service-name]
```

### Stop Environment
```bash
nexlayer compose down
```

> **Local endpoints:** Frontend at `http://localhost:3000`, Backend at `http://localhost:3000`.

---

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

**Happy Deploying! 🚀**

Deploy full-stack AI-powered applications effortlessly with Nexlayer CLI!
