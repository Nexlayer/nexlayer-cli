# Nexlayer CLI

> ⚠️ **Pre-Release Notice**: This project is currently in early development (pre-beta). The codebase is not yet ready for production use or forking. We expect to release beta v1 in Q2 2025. Until then, the repository will remain private and invite-only.

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

---

## 🚀 What is Nexlayer?

**Nexlayer** is the fastest way to **deploy full-stack AI applications** with a single command.  
It automates **containerized full-stack AI deployments** on production-ready enterprise-grade kubernetes without complex setup or infrastructure management.

### 🔥 Why Use Nexlayer?
🚀 **Instant deployments. Infinite scale. Zero DevOps. All without Kubernetes complexity.**

✅ **Zero DevOps required** – Deploy without managing Kubernetes or infrastructure.  
✅ **One-command deploys** – `nexlayer deploy` gets your app live instantly.  
✅ **Smart project detection** – `nexlayer init` auto-configures your stack.  
✅ **Scales automatically** – Enterprise-grade auto-scaling, no config needed.  
✅ **Custom domains** – `nexlayer domain set` links your app to a domain in seconds.  
✅ **Simple monitoring** – `nexlayer info` provides instant deployment insights.  
✅ **True developer speed** – No YAML headaches, just focus on your code.  

🔥 **From local dev to internet scale in seconds—no infrastructure, no limits, no hassle.** 🚀

---

## ⚡ Quick Start

### **1️⃣ Install Nexlayer CLI**
```bash
curl -sSL https://raw.githubusercontent.com/Nexlayer/nexlayer-cli/main/install.sh | bash
```

### **2️⃣ Create and Initialize a Project**
```bash
mkdir myapp && cd myapp
nexlayer init
```
- Auto-detects your framework (Next.js, Python, etc.)
- Generates a `nexlayer.yaml` deployment file
- Sets up environment variables and dependencies

### **3️⃣ Deploy Your Application**
```bash
nexlayer deploy
```
- Instantly deploys your app
- Generates build artifacts, provisions infrastructure, and handles CDN caching

## 💻 Command Reference

### Core Commands
1. **nexlayer init** – Initialize a new project (auto-detects type).  
2. **nexlayer deploy** – Deploy an application (uses `nexlayer.yaml` if present).  
3. **nexlayer list** – List active deployments.  
4. **nexlayer info <namespace> <appID>** – Get deployment details.  
5. **nexlayer domain** – Manage custom domains.  
6. **nexlayer login** – Authenticate with Nexlayer.  
7. **nexlayer watch** – Monitor & auto-deploy changes.  
8. **nexlayer feedback** – Send CLI feedback.  

### Global Flags
```bash
-h, --help         Show help for commands
    --json         Output response in JSON format
```

## 🛠 Example: Deploying a Next.js App

Let's deploy a simple Next.js app with Nexlayer.
https://github.com/Nexlayer/hello-world-nextjs

### 📂 Project Structure
```
hello-world-nextjs/
├── app/                      # Next.js application
│   ├── pages/                # Next.js pages (routes)
│   │   ├── index.tsx         # Homepage
│   │   ├── about.tsx         # Example additional page
│   ├── public/               # Static assets (images, icons, etc.)
│   │   ├── logo.png          # Example asset
│   ├── package.json          # Node.js dependencies
│   ├── next.config.ts        # Next.js configuration
│   ├── tsconfig.json         # TypeScript configuration
├── nginx/                    # NGINX configuration (Reverse Proxy)
│   ├── default.conf          # NGINX site config
│   ├── nginx.conf            # Global NGINX settings
├── Dockerfile                # Defines the container image
├── nexlayer.yaml             # Nexlayer deployment configuration
├── .gitignore                # Git ignore file
├── README.md                 # Documentation
```

### 🔧 nexlayer.yaml Configuration
```yaml
application:
  name: "Hello World NextJS App"
  pods:
  - name: nextjs-nginx
    path: /
    image: ghcr.io/nexlayer/hello-world-nextjs:v0.0.1
    servicePorts:
    - 80
```

### 🚀 Deploy the App
```bash
nexlayer deploy
```
- Detects the framework automatically
- Builds and deploys the application
- Assigns a default domain (e.g., `myapp.nexlayer.app`)

### 🔍 How It Works
- Nexlayer detects `next.config.js` and automatically provisions a Next.js environment
- It builds the static site and deploys it on an optimized global CDN
- Rollbacks are instant if something goes wrong

## 💻 Command Reference

```bash
# Initialize a new project
nexlayer init                # Auto-detect project type
nexlayer init -i             # Interactive mode
nexlayer init --type react   # Initialize React project

# Deploy an application
nexlayer deploy              # Deploy using nexlayer.yaml
nexlayer deploy myapp        # Deploy specific application
nexlayer deploy -f config.yaml  # Deploy with a custom config

# Watch mode for auto-deployment
nexlayer watch               # Auto-redeploy on changes

# Monitoring
nexlayer list                # Show all deployments
nexlayer info myapp          # Show deployment details
nexlayer list --json         # Output results in JSON format

# Configure a custom domain
nexlayer domain set myapp --domain example.com

# Send feedback
nexlayer feedback            # Share feedback or report issues
```

## 📚 Documentation
- 📖 [YAML Reference](docs/reference/schemas/yaml/README.md) – Configure `nexlayer.yaml`
- 📡 [API Reference](docs/reference/api/README.md) – Manage deployments via API

## 💪 Contributing
We welcome contributions! Check out our [Contributing Guide](CONTRIBUTING.md) to get involved.

## 📜 License
Nexlayer CLI is [MIT licensed](LICENSE).

## 🚀 Ready to Deploy?
- 🔹 Website: [nexlayer.com](https://nexlayer.com)
- 🔹 Docs: [docs.nexlayer.com](https://docs.nexlayer.com)
- 🔹 Feedback: [Join discussion](https://github.com/Nexlayer/nexlayer-cli/issues)
