# Nexlayer CLI

<div align="center">
  <img src="https://raw.githubusercontent.com/Nexlayer/nexlayer-cli/main/assets/logo.png" alt="Nexlayer Logo" width="200"/>
  <h1>Deploy Full-Stack AI-Powered Applications in Seconds</h1>
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
- [Features](#features)
- [Templates](#templates)
  - [AI/LLM](#aillm)
  - [Traditional Web Applications](#traditional-web-applications)
  - [Machine Learning](#machine-learning)
- [Commands Overview](#commands-overview)
- [Documentation & Support](#documentation--support)
- [Local Development](#local-development)
- [Contributing](#contributing)
- [License](#license)

---

## Quick Start

Get up and running in **3 seconds** with these three simple commands:

### Install the CLI
```bash
go install github.com/Nexlayer/nexlayer-cli@latest
```
> **Tip:** Ensure `$GOPATH/bin` is in your `PATH`.

### Initialize Your Project
```bash
nexlayer init
```
This command creates a `nexlayer.yaml` file in your project folder using AI-powered detection.

### Deploy Your App
```bash
nexlayer deploy
```
Watch your full-stack AI app go live instantly!

> **Bonus:** To see a demo, check out our demo video.

---

## Features

- **AI-Powered Template Generation:** Automatically detect your project and generate the perfect configuration with minimal effort.
- **One-Command Deployment:** Deploy your app with a single command.
- **Real-Time Logs & Status:** Monitor your deployment status and view logs easily.
- **Custom Domain & Feedback:** Attach custom domains and send feedback to help us improve.
- **Plugin Support:** Extend Nexlayer CLI with custom plugins.

---

## Templates

Nexlayer offers a variety of ready-to-use templates:

### AI/LLM
- `langchain-nextjs`: LangChain.js with Next.js
- `langchain-fastapi`: LangChain Python with FastAPI
- `openai-node`: OpenAI with Express and React
- `openai-py`: OpenAI with FastAPI and Vue
- `llama-node`: Llama.cpp with Next.js
- `llama-py`: Llama.cpp with FastAPI
- `vertex-ai`: Google Vertex AI with Flask
- `huggingface`: Hugging Face with FastAPI
- `anthropic-py`: Anthropic Claude with FastAPI
- `anthropic-js`: Anthropic Claude with Next.js

### Traditional Web Applications
- `mern`: MongoDB, Express, React, Node.js
- `mean`: MongoDB, Express, Angular, Node.js
- `mevn`: MongoDB, Express, Vue.js, Node.js
- `pern`: PostgreSQL, Express, React, Node.js

### Machine Learning
- `kubeflow`: ML pipeline with Kubeflow
- `mlflow`: MLflow with tracking server

> **Tip:** Use interactive mode:
```bash
nexlayer init my-project
```
Or directly specify a template:
```bash
nexlayer init my-project -t fastapi
```

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
nexlayer status
nexlayer logs -f [podName]
```
Monitor your deployment and view real-time logs.

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

We welcome contributions! Please see our [Contributing Guidelines](#) for more information.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

**Happy Deploying! 🚀**

Deploy full-stack AI-powered applications effortlessly with Nexlayer CLI!

