# Nexlayer CLI

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

**Nexlayer** is the fastest way to **deploy full-stack applications** with a single command.  
It automates **containerized, serverless, and full-stack deployments** without complex infrastructure setup.  

### **🔥 Why Use Nexlayer?**
✅ **One-command deploys** – `nexlayer deploy` auto-detects your stack.  
✅ **Built-in scaling** – Auto-scales with no manual config.  
✅ **Zero DevOps required** – Works out of the box.  
✅ **Instant rollbacks** – Deploy safely with built-in versioning.  
✅ **Live Watch Mode** – Auto-redeploy when code changes.  

---

# **⚡ Create Your First App**  

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

### **3️⃣ Deploy in Seconds**
```bash
nexlayer deploy
```
- Instantly deploys your app
- Generates build artifacts, provisions infrastructure, and handles CDN caching

### **4️⃣ Watch for Live Changes**
```bash
nexlayer watch
```
- Auto-redeploys when code changes
- Ideal for local development

## 🛠 Example: Deploying a Next.js App

Let's deploy a simple Next.js app with Nexlayer.

### 📂 Project Structure
```
myapp/
 ├── pages/
 │    ├── index.js
 │    ├── about.js
 ├── public/
 │    ├── logo.png
 ├── package.json
 ├── nexlayer.yaml
```

### 🔧 nexlayer.yaml Configuration
```yaml
name: myapp
runtime: node
build:
  command: npm install && npm run build
  output: .next
deploy:
  port: 3000
  env:
    NEXT_PUBLIC_API_URL: "https://api.example.com"
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
- 🔹 Website: [nexlayer.dev](https://nexlayer.dev)
- 🔹 Docs: [nexlayer.dev/docs](https://nexlayer.dev/docs)
- 🔹 Community: [Join Discord](https://discord.gg/nexlayer)
