<div align="center">
  <img src="assets/logo.svg" alt="Nexlayer Logo" width="400"/>
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

# Create a new project
nexlayer init my-app

# Deploy your app
nexlayer deploy
```

That's it! Your app is live. [Watch the demo →](https://nexlayer.dev/demo)

## ✨ Features

- 🤖 **AI-Powered Detection** - Automatically analyze and configure your project
- 🎯 **Smart Templates** - Production-ready templates for any stack
- 🔄 **Live Sync** - Keep configuration in sync with project changes
- 🚀 **One-Command Deploy** - Deploy full-stack apps instantly
- 📊 **Real-Time Monitoring** - Live logs and deployment status
- 🔌 **Plugin System** - Extend functionality with custom plugins

## 📝 Templates

```bash
# Create a new project with an interactive template selector
nexlayer init my-app
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
nexlayer deploy

# View status and logs
nexlayer status
nexlayer logs -f [pod]

# Keep config in sync
nexlayer sync
```

Full documentation at [nexlayer.dev/docs](https://nexlayer.dev/docs)
## 👷 Development

```bash
# Clone the repository
git clone https://github.com/Nexlayer/nexlayer-cli.git
cd nexlayer-cli

# Install dependencies
make setup

# Run tests
make test
```

## 💪 Contributing

We love contributions! See our [Contributing Guide](CONTRIBUTING.md) for details.

## 📜 License

Nexlayer CLI is [MIT licensed](LICENSE).
