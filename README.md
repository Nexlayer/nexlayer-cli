# Nexlayer CLI

> ⚠️ **Pre-Release Notice**: This project is currently in early development (pre-beta). The codebase is not yet ready for production use or forking. We expect to release beta v1 in Q2 2025. Until then, the repository will remain private and invite-only.

<div align="center">
  <img src="pkg/ui/assets/logo.svg" alt="Nexlayer Logo" width="400"/>
  <h1>Nexlayer CLI</h1>
  <p><strong>Deploy Full-Stack AI Applications in Seconds ⚡️</strong></p>
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

**Recommended**: Install with the automated script (supports all features)
```bash
curl -sSL https://raw.githubusercontent.com/Nexlayer/nexlayer-cli/main/install.sh | bash
```
- ✅ Configures shell environment automatically
- ✅ Verifies system requirements
- ✅ Supports both global and local installation
- ✅ Better project path handling and error reporting

**Alternative**: Install directly using Go (minimal installation)
```bash
go install github.com/Nexlayer/nexlayer-cli@latest
```
- ✅ Simple one-line installation
- ✅ Uses Go's standard package management
- ⚠️ Manual shell configuration may be needed
- ⚠️ Limited to current working directory

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
4. **nexlayer info <namespace> [appID]** – Get deployment details.  
   - Use `--verbose` flag for detailed information about pods, resources, and configuration.
   - Example: `nexlayer info my-namespace --verbose`
5. **nexlayer domain** – Manage custom domains.  
6. **nexlayer login** – Authenticate with Nexlayer.  
7. **nexlayer watch** – Monitor project changes and update configuration.  
8. **nexlayer feedback** – Send CLI feedback.  

### Watch Mode
The `watch` command runs in the foreground, actively monitoring your project for changes:

```bash
nexlayer watch
```

**Key Features:**
- Monitors project directory in real-time
- Auto-detects new dependencies, frameworks, and services
- Updates `nexlayer.yaml` automatically when changes are detected
- Handles Docker image updates and configuration changes
- Press `Ctrl+C` to stop watching

**Example Output:**
```
Watching for project changes...
Configuration will be updated when new components are detected.

Analyzing project changes...
📝 Configuration changes detected:
+ Added new Docker image: postgres:latest
+ Updated service configuration
+ Added database dependencies

Configuration updated successfully.
```

### Global Flags
```bash
-h, --help         Show help for commands
    --json         Output response in JSON format
    --verbose      Display detailed information (available for info command)
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

> **Note:** The definitive schema for nexlayer.yaml configuration is maintained in the [schema package](pkg/schema/README.md), which serves as the single source of truth for all YAML configurations.

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

## 📚 Documentation
- 📖 [YAML Schema](pkg/schema/README.md) – Single source of truth for `nexlayer.yaml` configuration
- 📡 [API Reference](docs/reference/api/README.md) – Manage deployments via API

## 💪 Contributing
We welcome contributions! Check out our [Contributing Guide](CONTRIBUTING.md) to get involved.

## 📜 License
Nexlayer CLI is [MIT licensed](LICENSE).

## 🚀 Ready to Deploy?
- 🔹 Website: [nexlayer.com](https://nexlayer.com)
- 🔹 Docs: [docs.nexlayer.com](https://docs.nexlayer.com)
- 🔹 Feedback: [Join discussion](https://github.com/Nexlayer/nexlayer-cli/issues)
